package handler

import (
	"backend/auth"
	"backend/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// Register all employee-related routes
func RegisterEmployeeRoutes(db *sql.DB) {
	http.Handle("/shifts", CORS(auth.JWTMiddleware(GetAvailableShifts(db))))
	http.Handle("/myshifts", CORS(auth.JWTMiddleware(GetAssignedShifts(db))))
	http.Handle("/requests", CORS(auth.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			RequestShift(db)(w, r)
		} else {
			GetMyRequests(db)(w, r)
		}
	})))
}

func GetAvailableShifts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
        SELECT s.id, s.date, s.start, s.end, s.role, s.location
        FROM shifts s
        WHERE s.id NOT IN (
            SELECT shift_id FROM requests WHERE status = 'approved'
        )`)
		if err != nil {
			http.Error(w, "Query failed", http.StatusInternalServerError)
			return
		}

		var shifts []models.Shift
		for rows.Next() {
			var s models.Shift
			rows.Scan(&s.ID, &s.Date, &s.Start, &s.End, &s.Role, &s.Location)
			shifts = append(shifts, s)
		}

		json.NewEncoder(w).Encode(shifts)
	}
}

func GetAssignedShifts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
        SELECT s.id, s.date, s.start, s.end, s.role, s.location
        FROM requests r
        JOIN shifts s ON r.shift_id = s.id
        WHERE r.user_id = ? AND r.status = 'approved'`, user.UserID)
		if err != nil {
			http.Error(w, "Query failed", http.StatusInternalServerError)
			return
		}

		var shifts []models.Shift
		for rows.Next() {
			var s models.Shift
			rows.Scan(&s.ID, &s.Date, &s.Start, &s.End, &s.Role, &s.Location)
			shifts = append(shifts, s)
		}

		json.NewEncoder(w).Encode(shifts)
	}
}

func RequestShift(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req models.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var date, start, end string
		err := db.QueryRow("SELECT date, start, end FROM shifts WHERE id = ?", req.ShiftID).Scan(&date, &start, &end)
		if err != nil {
			log.Println(err)
			http.Error(w, "Shift not found", http.StatusNotFound)
			return
		}

		// Check if shift already assigned
		var assigned int
		db.QueryRow("SELECT COUNT(*) FROM requests WHERE shift_id = ? AND status = 'approved'", req.ShiftID).Scan(&assigned)
		if assigned > 0 {
			log.Println(err)
			http.Error(w, "Shift already taken", http.StatusConflict)
			return
		}

		// Overlap check (same day & overlap time)
		rows, err := db.Query(`
			SELECT s.start, s.end FROM requests r
			JOIN shifts s ON r.shift_id = s.id
			WHERE r.user_id = ? AND s.date = ? AND r.status IN ('pending', 'approved')`,
			user.UserID, date)
		if err != nil {
			log.Println(err)
			http.Error(w, "Overlap check failed", http.StatusInternalServerError)
			return
		}

		// 1. Check if already has a shift on this date
		var dayCount int
		err = db.QueryRow(`
			SELECT COUNT(*) FROM requests r
			JOIN shifts s ON r.shift_id = s.id
			WHERE r.user_id = ? AND s.date = ? AND r.status IN ('pending', 'approved')`,
			user.UserID, date).Scan(&dayCount)

		if dayCount >= 1 {
			log.Println(err)
			http.Error(w, "Already has a shift on this day", http.StatusConflict)
			return
		}

		// 2. Check total shifts this week
		// Get year + ISO week number of the shift being requested
		var weekNum int
		err = db.QueryRow(`SELECT strftime('%W', ?)`, date).Scan(&weekNum)
		if err != nil {
			log.Println(err)
			http.Error(w, "Week calculation failed", http.StatusInternalServerError)
			return
		}

		rows, err = db.Query(`
			SELECT s.date FROM requests r
			JOIN shifts s ON r.shift_id = s.id
			WHERE r.user_id = ? AND r.status IN ('pending', 'approved')`, user.UserID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Week limit check failed", http.StatusInternalServerError)
			return
		}

		shiftCount := 0
		for rows.Next() {
			var d string
			rows.Scan(&d)
			var w int
			db.QueryRow(`SELECT strftime('%W', ?)`, d).Scan(&w)
			if w == weekNum {
				shiftCount++
			}
		}
		if shiftCount >= 5 {
			http.Error(w, "Exceeds 5 shifts per week", http.StatusConflict)
			return
		}

		for rows.Next() {
			var es, ee string
			rows.Scan(&es, &ee)
			if start < ee && end > es {
				http.Error(w, "Overlapping shift exists", http.StatusConflict)
				return
			}
		}

		_, err = db.Exec("INSERT INTO requests (user_id, shift_id, status) VALUES (?, ?, 'pending')",
			user.UserID, req.ShiftID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to insert request", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func GetMyRequests(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
        SELECT id, shift_id, user_id, status FROM requests
        WHERE user_id = ?`, user.UserID)
		if err != nil {
			http.Error(w, "Query failed", http.StatusInternalServerError)
			return
		}

		var requests []models.Request
		for rows.Next() {
			var req models.Request
			rows.Scan(&req.ID, &req.ShiftID, &req.UserID, &req.Status)
			requests = append(requests, req)
		}

		json.NewEncoder(w).Encode(requests)
	}
}
