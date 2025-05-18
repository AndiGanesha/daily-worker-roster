package handler

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"backend/auth"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var id int
		var name, role, hashed string
		err := db.QueryRow("SELECT id, name, role, password FROM users WHERE username = ?", req.Username).Scan(&id, &name, &role, &hashed)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		hashedInput := HashPassword(req.Password)
		if hashedInput != hashed {
			http.Error(w, "Wrong password", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(id, role, name)
		if err != nil {
			http.Error(w, "Token generation failed", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

func GetUser(r *http.Request) *auth.Claims {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == nil || user.Role != "admin" {
			http.Error(w, "Forbidden: admin only", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
