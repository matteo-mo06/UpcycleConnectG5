package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"upcycle_connect-api/internal/db"
)

func GetStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats, err := db.GetStats()
	if err != nil {
		fmt.Println("GetStats error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch stats"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stats)
}
