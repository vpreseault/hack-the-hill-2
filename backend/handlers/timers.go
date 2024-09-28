package handlers

import (
	"net/http"

	"github.com/vpreseault/hack-the-hill-2/backend/database"
)

func StartTimer(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func StopTimer(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}