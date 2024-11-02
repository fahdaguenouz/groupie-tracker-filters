package handlers

import (
	"groupie/models"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	artists, err := FetchArtists()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Extract unique locations from the artist data
	locationSet := make(map[string]struct{})
	for _, artist := range artists {
		locations := strings.Split(artist.Locations, ",") // Assuming locations are comma-separated
		for _, loc := range locations {
			locationSet[strings.TrimSpace(loc)] = struct{}{}
		}
	}
	// Create a slice of unique locations
	var locations []string
	for loc := range locationSet {
		locations = append(locations, loc)
	}

	data := struct {
		Artists   []models.Artist
		Locations []string
	}{
		Artists:   artists,
		Locations: locations,
	}
	// Render the main template
	if err := renderTemplate(w, "index.html", data); err != nil {
		log.Printf("Error rendering template: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) error {
	tmpl, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		// Return the error to be handled by the calling function
		return err
	}

	if err := tmpl.Execute(w, data); err != nil {
		// Return the error to be handled by the calling function
		return err
	}

	return nil
}
