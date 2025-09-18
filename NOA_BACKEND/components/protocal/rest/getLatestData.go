package rest

import (
	"encoding/json"
	"net/http"

	"GOLANG_SERVER/components/db"
	schema "GOLANG_SERVER/components/schema"
)

// * get latest data
func HandleGetLatestData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		// * get data from request
		var req schema.GyroData
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// Get the data from the database
		data, err := db.GetGyroDataByDeviceAddressLatest(req.DeviceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the data into JSON
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
