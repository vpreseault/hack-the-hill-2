package handlers

import (
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