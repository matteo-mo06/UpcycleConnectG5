package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/sse"
)

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "streaming non supporté"})
		return
	}

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch := sse.Default.Subscribe(userID)
	defer sse.Default.Unsubscribe(userID, ch)

	for {
		select {
		case event := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", event)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
