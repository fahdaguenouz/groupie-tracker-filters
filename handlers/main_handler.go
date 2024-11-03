package handlers

import (
	"groupie/controllers"
	"groupie/models"
	"html/template"
	"log"
	"net/http"
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

	// Fetch locations from the API
	var locations []string
	for _, artist := range artists {
		var loca models.Location
		err := controllers.FetchData(artist.Locations, &loca)
		if err == nil {
			locations = append(locations, loca.Locations...)
		}
	}

	// Remove duplicates from locations
	locationSet := make(map[string]struct{})
	for _, loc := range locations {
		locationSet[loc] = struct{}{}
	}

	var uniqueLocations []string
	for loc := range locationSet {
		uniqueLocations = append(uniqueLocations, loc)
	}

	data := struct {
		Artists      []models.Artist
		Locations    []string
	}{
		Artists:   artists,
		Locations: uniqueLocations,
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
