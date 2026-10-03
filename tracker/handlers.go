package tracker

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderError(
			w,
			http.StatusNotFound,
			"Page not found",
			"The page you requested does not exist.",
		)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)

		renderError(
			w,
			http.StatusMethodNotAllowed,
			"Method not allowed",
			"Only GET requests are allowed.",
		)
		return
	}

	artists, err := fetchArtists()
	if err != nil {
		log.Println("Error fetching artists:", err)

		renderError(
			w,
			http.StatusBadGateway,
			"Service unavailable",
			"Artist data could not be loaded right now.",
		)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	artists = filterArtists(artists, query)

	data := HomePageData{
		Artists: artists,
		Query:   query,
		Total:   len(artists),
	}

	renderTemplate(w, "templates/index.html", data)
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)

		renderError(
			w,
			http.StatusMethodNotAllowed,
			"Method not allowed",
			"Only GET requests are allowed.",
		)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		renderError(
			w,
			http.StatusBadRequest,
			"Invalid artist",
			"The artist ID is missing or invalid.",
		)
		return
	}

	artists, err := fetchArtists()
	if err != nil {
		log.Println("Error fetching artists:", err)

		renderError(
			w,
			http.StatusBadGateway,
			"Service unavailable",
			"Artist data could not be loaded right now.",
		)
		return
	}

	artist, found := findArtist(artists, id)

	if !found {
		renderError(
			w,
			http.StatusNotFound,
			"Artist not found",
			"The requested artist does not exist.",
		)
		return
	}

	relation, err := fetchRelation(id)
	if err != nil {
		log.Println("Error fetching relation:", err)

		renderError(
			w,
			http.StatusBadGateway,
			"Service unavailable",
			"Concert data could not be loaded right now.",
		)
		return
	}

	data := ArtistPageData{
		Artist:   artist,
		Concerts: relation.DatesLocations,
	}

	renderTemplate(w, "templates/artist.html", data)
}

func filterArtists(artists []Artist, query string) []Artist {
	if query == "" {
		return artists
	}

	query = strings.ToLower(query)

	var result []Artist

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), query) {
			result = append(result, artist)
		}
	}

	return result
}

func findArtist(artists []Artist, id int) (Artist, bool) {
	for _, artist := range artists {
		if artist.ID == id {
			return artist, true
		}
	}

	return Artist{}, false
}

func renderTemplate(w http.ResponseWriter, file string, data any) {
	tmpl, err := template.ParseFiles(file)

	if err != nil {
		log.Println("Template loading error:", err)

		renderError(
			w,
			http.StatusInternalServerError,
			"Internal server error",
			"The page could not be loaded.",
		)
		return
	}

	var output bytes.Buffer

	if err := tmpl.Execute(&output, data); err != nil {
		log.Println("Template execution error:", err)

		renderError(
			w,
			http.StatusInternalServerError,
			"Internal server error",
			"The page could not be loaded.",
		)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if _, err := output.WriteTo(w); err != nil {
		log.Println("Response writing error:", err)
	}
}

func renderError(w http.ResponseWriter, code int, title string, message string) {
	tmpl, err := template.ParseFiles("templates/error.html")

	if err != nil {
		http.Error(w, http.StatusText(code), code)
		return
	}

	var output bytes.Buffer

	data := ErrorPageData{
		Code:    code,
		Title:   title,
		Message: message,
	}

	if err := tmpl.Execute(&output, data); err != nil {
		http.Error(w, http.StatusText(code), code)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)

	if _, err := output.WriteTo(w); err != nil {
		log.Println("Error response writing error:", err)
	}
}
