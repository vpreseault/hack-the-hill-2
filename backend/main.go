package main

import (
	"log"
	"net/http"

	"github.com/vpreseault/hack-the-hill-2/backend/database"
	"github.com/vpreseault/hack-the-hill-2/backend/handlers"
)

func main() {
	const port = "8080"

	db, err := database.CreateDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// General
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", handlers.Root())

	// Session
	mux.HandleFunc("GET /sessions/{sessionID}", handlers.AddUserToSession(db))
	mux.HandleFunc("GET /api/sessions/{sessionID}", handlers.GetSessionInfo(db))
	mux.HandleFunc("POST /api/sessions", handlers.CreateSession(db))

	// Timer
	mux.HandleFunc("POST /api/sessions/{sessionID}/timer/start", handlers.StartTimer(db))
	mux.HandleFunc("POST /api/sessions/{sessionID}/timer/stop", handlers.StopTimer(db))

	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}