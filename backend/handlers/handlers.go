package handlers

import (
	"net/http"
)

func Root() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/static/index.html")
	}
}
