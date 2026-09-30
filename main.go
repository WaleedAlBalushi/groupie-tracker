package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistPageData struct {
	Artist   Artist
	Concerts map[string][]string
}

func fetchArtists() ([]Artist, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("artists API returned status %d", resp.StatusCode)
	}

	var artists []Artist

	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}

	return artists, nil
}

func fetchRelation(id int) (Relation, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf(
		"https://groupietrackers.herokuapp.com/api/relation/%d",
		id,
	)

	resp, err := client.Get(url)
	if err != nil {
		return Relation{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Relation{}, fmt.Errorf("relation API returned status %d", resp.StatusCode)
	}

	var relation Relation

	if err := json.NewDecoder(resp.Body).Decode(&relation); err != nil {
		return Relation{}, err
	}

	return relation, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	artists, err := fetchArtists()
	if err != nil {
		log.Println("Error fetching artists:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		log.Println("Template loading error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var output bytes.Buffer

	if err := tmpl.Execute(&output, artists); err != nil {
		log.Println("Template execution error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := output.WriteTo(w); err != nil {
		log.Println("Response writing error:", err)
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artists, err := fetchArtists()
	if err != nil {
		log.Println("Could not fetch artists:", err)
		http.Error(w, "Unable to load artist", http.StatusBadGateway)
		return
	}

	var selectedArtist Artist
	found := false

	for _, artist := range artists {
		if artist.ID == id {
			selectedArtist = artist
			found = true
			break
		}
	}

	if !found {
		http.NotFound(w, r)
		return
	}

	relation, err := fetchRelation(id)
	if err != nil {
		log.Println("Could not fetch relation:", err)
		http.Error(w, "Unable to load artist details", http.StatusBadGateway)
		return
	}

	data := ArtistPageData{
		Artist:   selectedArtist,
		Concerts: relation.DatesLocations,
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		log.Println("Template loading error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var output bytes.Buffer

	if err := tmpl.Execute(&output, data); err != nil {
		log.Println("Template rendering error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := output.WriteTo(w); err != nil {
		log.Println("Response writing error:", err)
	}
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/artist", artistHandler)

	http.HandleFunc("/", homeHandler)

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
