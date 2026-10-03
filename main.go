package main

import (
	"log"
	"net/http"

	"groupie-tracker/tracker"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/artist", tracker.ArtistHandler)
	http.HandleFunc("/", tracker.HomeHandler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
