package handlers

import (
	"fmt"
	"net/http"

	"github.com/vpreseault/hack-the-hill-2/backend/cookies"
	"github.com/vpreseault/hack-the-hill-2/backend/database"
)

func GetSessionInfo(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.PathValue("sessionID")
		session, err := db.GetSessionInfo(sessionID)
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

		response := fmt.Sprintf(template, sessionID, session.Users, session.Timer)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(response))
		if err != nil {
			http.Error(w, "Unable to write response", http.StatusInternalServerError)
		}
	}
}

func CreateSession(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            http.Error(w, "Unable to parse form data", http.StatusBadRequest)
            return
        }

        hostName := r.FormValue("username")
        
        if hostName == "" {
            http.Error(w, "Username is required", http.StatusBadRequest)
            return
        }

		cookie := &http.Cookie{
			Name:  "user",
			Value: hostName,
			Path:     "/",
		}
		http.SetCookie(w, cookie)

		sessionIDCookie, err := r.Cookie("sessionID")
		if err == nil && sessionIDCookie.Value != "" {
			sessionID := sessionIDCookie.Value
			// user already has a session, delete cookie and redirect to session page
			sessionIDCookie.Value = ""
			http.SetCookie(w, sessionIDCookie)

			w.Header().Set("HX-Redirect", "/"+sessionID)
			w.WriteHeader(http.StatusOK)
			return
		}

		sessionID, err := db.CreateSession(hostName)
		if err != nil {
			http.Error(w, "Unable to create session", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/"+sessionID)
		w.WriteHeader(http.StatusOK)
	}
}

func AddUserToSession(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.PathValue("sessionID")

		user, err := cookies.GetUserNameFromCookie(r)
		if err != nil {
			// no user name, redirect to root
			saveSessionIDInCookie(w, sessionID)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		err = db.AddUserToSession(sessionID, user)
		if err != nil {
			http.Error(w, "Unable to add user to session", http.StatusInternalServerError)
			return
		}

		renderTemplate(w, "timer")
	}
}