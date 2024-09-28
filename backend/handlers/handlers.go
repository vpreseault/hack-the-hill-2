package handlers

import (
	"fmt"
	"net/http"
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

func Session() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionIDPathValue := r.PathValue("sessionId")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		template := `
		<html>
		<body>
			<h1>Session ID: %s</h1>
		</body>
		</html>
		`

		response := fmt.Sprintf(template, sessionIDPathValue)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(response))
		if err != nil {
			http.Error(w, "Unable to write response", http.StatusInternalServerError)
		}
	}
}