package handler

import (
	"backend/auth"
	"backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Register admin-related routes
func RegisterAdminRoutes(db *sql.DB) {
	http.Handle("/admin/requests", CORS(auth.JWTMiddleware(AdminOnly(GetPendingRequests(db)))))
	http.Handle("/requests/", CORS(auth.JWTMiddleware(AdminOnly(HandleApproval(db)))))

	http.Handle("/admin/shifts", CORS(auth.JWTMiddleware(AdminOnly(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			CreateShift(db)(w, r)
		}
	}))))

	http.Handle("/admin/shifts/", CORS(auth.JWTMiddleware(AdminOnly(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			UpdateShift(db)(w, r)
		case "DELETE":
			DeleteShift(db)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))
}

func GetPendingRequests(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
        SELECT r.id, u.name, s.date, s.start, s.end, s.role, s.location, r.status
        FROM requests r
        JOIN users u ON r.user_id = u.id
        JOIN shifts s ON r.shift_id = s.id
        WHERE r.status = 'pending'`)
		if err != nil {
			http.Error(w, "Query failed", http.StatusInternalServerError)
			return
		}

		type AdminRequestView struct {
			RequestID int    `json:"request_id"`
			UserName  string `json:"user_name"`
			ShiftDate string `json:"date"`
			Start     string `json:"start"`
			End       string `json:"end"`
			Role      string `json:"role"`
			Location  string `json:"location"`
			Status    string `json:"status"`
		}

		var results []AdminRequestView
		for rows.Next() {
			var r AdminRequestView
			rows.Scan(&r.RequestID, &r.UserName, &r.ShiftDate, &r.Start, &r.End, &r.Role, &r.Location, &r.Status)
			results = append(results, r)
		}

		json.NewEncoder(w).Encode(results)
	}
}

func HandleApproval(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/requests/")
		var status string

		if strings.HasSuffix(path, "/approve") {
			path = strings.TrimSuffix(path, "/approve")
			status = "approved"
		} else if strings.HasSuffix(path, "/reject") {
			path = strings.TrimSuffix(path, "/reject")
			status = "rejected"
		} else {
			http.Error(w, "Unknown action", http.StatusBadRequest)
			return
		}

		reqID, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid request ID", http.StatusBadRequest)
			return
		}

		// Conflict check only if approving
		if status == "approved" {
			var shiftID int
			err = db.QueryRow("SELECT shift_id FROM requests WHERE id = ?", reqID).Scan(&shiftID)
			if err != nil {
				http.Error(w, "Request not found", http.StatusNotFound)
				return
			}

			var taken int
			db.QueryRow("SELECT COUNT(*) FROM requests WHERE shift_id = ? AND status = 'approved'", shiftID).Scan(&taken)
			if taken > 0 {
				http.Error(w, "Shift already approved for another worker", http.StatusConflict)
				return
			}
		}

		_, err = db.Exec("UPDATE requests SET status = ? WHERE id = ?", status, reqID)
		if err != nil {
			http.Error(w, "Update failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request " + status))
	}
}

func CreateShift(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		var s models.Shift
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO shifts (date, start, end, role, location) VALUES (?, ?, ?, ?, ?)",
			s.Date, s.Start, s.End, s.Role, s.Location)
		if err != nil {
			http.Error(w, "Insert failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateShift(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/admin/shifts/")
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var s models.Shift
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(`UPDATE shifts SET date = ?, start = ?, end = ?, role = ?, location = ? WHERE id = ?`,
			s.Date, s.Start, s.End, s.Role, s.Location, id)
		if err != nil {
			http.Error(w, "Update failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteShift(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/admin/shifts/")
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM shifts WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Delete failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
