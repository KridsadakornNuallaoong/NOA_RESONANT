package rest

import (
	"encoding/json"
	"net/http"

	"GOLANG_SERVER/components/db"
)

func HandleGetDeviceAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var userDetail map[string]string
	if err := json.NewDecoder(r.Body).Decode(&userDetail); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := userDetail["userID"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// ดึงข้อมูลอุปกรณ์จากฐานข้อมูล
	devices, err := db.GetDeviceAddress(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ส่งข้อมูลกลับในรูปแบบ JSON
	if err := json.NewEncoder(w).Encode(devices); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
