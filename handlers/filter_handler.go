package handlers

import (
	"groupie/models"
	"net/http"
	"strconv"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	// Fetch all artists
	artists, err := FetchArtists()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Get filter values from the form
	creationYear := r.FormValue("creationDate")
	firstAlbumYear := r.FormValue("firstAlbumDate")
	selectedMembers := r.Form["members"]
	selectedLocation := r.FormValue("locations")

	// Prepare to store filtered results
	var filteredArtists []models.Artist

	for _, artist := range artists {
		matches := true // Assume it matches unless proven otherwise

		// Filter by creation year if provided
		if creationYear != "" {
			year, err := strconv.Atoi(creationYear[:4]) // Extract year from the date
			if err != nil || artist.CreationDate != year {
				matches = false
			}
		}

		// Filter by first album year if provided
		if firstAlbumYear != "" {
			year, err := strconv.Atoi(firstAlbumYear) // Convert range value to int
			if err != nil || artist.FirstAlbum != strconv.Itoa(year) {
				matches = false
			}
		}

		// Filter by number of members if any checkboxes are selected
		if len(selectedMembers) > 0 {
			membersCount := len(artist.Members)
			memberMatches := false
			for _, member := range selectedMembers {
				if member == "5" && membersCount >= 5 {
					memberMatches = true
					break
				} else if strconv.Itoa(membersCount) == member {
					memberMatches = true
					break
				}
			}
			if !memberMatches {
				matches = false
			}
		}

		// Filter by location if selected
		if selectedLocation != "" && selectedLocation != "Select a location" {
			locationMatch := false
			for _, artistLocation := range artist.Loca.Locations {
				if artistLocation == selectedLocation {
					locationMatch = true
					break
				}
			}
			if !locationMatch {
				matches = false
			}
		}

		// If it passes all filters, add to results
		if matches {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// Create a template data structure
	data := struct {
		Artists []models.Artist
	}{
		Artists: filteredArtists,
	}

	// Render the filter results template
	if err := renderTemplate(w, "filter.html", data); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
