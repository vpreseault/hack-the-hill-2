package main

import (
	"log"
	"net/http"

	"github.com/vpreseault/hack-the-hill-2/backend/database"
	"github.com/vpreseault/hack-the-hill-2/backend/handlers"
)

func main() {
	const port = "8080"

	_, err := database.CreateDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Root())
	mux.HandleFunc("/{sessionId}", handlers.Session())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Backend server running on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}