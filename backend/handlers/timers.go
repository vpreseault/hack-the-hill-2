package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vpreseault/hack-the-hill-2/backend/database"
)

func StartTimer(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {		
		decoder := json.NewDecoder(r.Body)
		timer := database.Timer{}
		err := decoder.Decode(&timer)
		if err != nil {
			http.Error(w, "Unable to decode timer", http.StatusBadRequest)
			return
		}
		
		sessionID := r.PathValue("sessionID")
		log.Printf("Starting timer for session: %s", sessionID)

		err = db.StartTimer(sessionID, timer)
		if err != nil {
			http.Error(w, "Unable to start timer", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func StopTimer(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}