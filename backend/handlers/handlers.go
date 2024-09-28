package handlers

import (
	"fmt"
	"net/http"

	"github.com/vpreseault/hack-the-hill-2/backend/database"
)

func Root() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		response := `
		<html>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(response))
		if err != nil {
			http.Error(w, "Unable to write response", http.StatusInternalServerError)
		}
	}
}

func GetSessionInfo(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionIDPathValue := database.SessionID(r.PathValue("sessionId"))
		session, err := db.GetSessionInfo(sessionIDPathValue)
		if err != nil {
			http.Error(w, "Unable to get session info", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		template := `
		<html>
		<body>
			<h1>Session ID: %s</h1>
			<h2>Users: %v</h2>
			<h2>Timer: %v</h2>
		</body>
		</html>
		`

		response := fmt.Sprintf(template, sessionIDPathValue, session.Users, session.Timer)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(response))
		if err != nil {
			http.Error(w, "Unable to write response", http.StatusInternalServerError)
		}
	}
}
