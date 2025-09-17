package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"GOLANG_SERVER/configs/env"
	"GOLANG_SERVER/internal/adapters/db"
	schema "GOLANG_SERVER/pkg/schema"
)

func HandleCleanData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the password from client
	if r.Method == "POST" {
		var req schema.PasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// * check password
		if req.Password == req.CFP {
			if req.Password == env.GetEnv("PASSWORD") {
				// * clean data
				if _, err := db.CleanData(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{"message": "Data cleaned!"}`)
			} else {
				http.Error(w, "Invalid password", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Password doesn't match", http.StatusUnauthorized)
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
