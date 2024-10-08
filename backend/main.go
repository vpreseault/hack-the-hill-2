package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vpreseault/hack-the-hill-2/backend/database"
	"github.com/vpreseault/hack-the-hill-2/backend/handlers"
	"github.com/vpreseault/hack-the-hill-2/backend/sockets"
)

func main() {
	const port = "8080"

	db, err := database.CreateDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiRouter := http.NewServeMux()

	// General
	apiRouter.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static"))))
	apiRouter.HandleFunc("/", handlers.Root())

	// Session
	apiRouter.HandleFunc("GET /sessions/{sessionID}", handlers.AddUserToSession(db))
	apiRouter.HandleFunc("GET /api/sessions/{sessionID}", handlers.GetSessionInfo(db))
	apiRouter.HandleFunc("POST /api/sessions", handlers.CreateSession(db))

	// Timer
	apiRouter.HandleFunc("POST /api/sessions/{sessionID}/timer/start", handlers.StartTimer(db))
	apiRouter.HandleFunc("POST /api/sessions/{sessionID}/timer/stop", handlers.StopTimer(db))

	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: apiRouter,
	}


	hub := sockets.NewHub()
    go hub.Run()

    // Create a new router for the WebSocket server
    wsRouter := mux.NewRouter()
    wsRouter.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        sockets.ServeWs(hub, w, r)
    })

    // Start the main API server
    go func() {
		log.Fatal(srv.ListenAndServe())
    }()

    // Start the WebSocket server
	websocketServer := &http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: wsRouter,
	}

	log.Fatal(websocketServer.ListenAndServe())
}