package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"groupie/models"
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
	creationYear := r.FormValue("creationYear")
	firstAlbumYear := r.FormValue("firstAlbumYear")
	selectedMembers := r.Form["members"]
	selectedLocations := r.Form["locations"]

	// Prepare to store filtered results
	var filteredArtists []models.Artist

	for _, artist := range artists {
		matches := true // Assume it matches unless proven otherwise

		// Filter by creation year if provided
		if creationYear != "" {
			if year, err := strconv.Atoi(creationYear); err == nil && artist.CreationDate != year {
				matches = false
			}
		}

		// Filter by first album year if provided
		if firstAlbumYear != "" {
			firstAlbumDateParts := strings.Split(artist.FirstAlbum, "-") // Assuming FirstAlbum is in "YYYY-MM-DD" format
			if len(firstAlbumDateParts) > 0 {
				if firstAlbumYearInt, err := strconv.Atoi(firstAlbumDateParts[0]); err == nil {
					if providedFirstAlbumYearInt, err := strconv.Atoi(firstAlbumYear); err == nil && firstAlbumYearInt != providedFirstAlbumYearInt {
						matches = false
					}
				} else {
					matches = false
				}
			} else {
				matches = false
			}
		}

		// Filter by number of members if provided
		if matches && len(selectedMembers) > 0 {
			membersCount := len(artist.Members)
			memberMatches := false
			for _, member := range selectedMembers {
				if member == "5" && membersCount >= 5 || strconv.Itoa(membersCount) == member {
					memberMatches = true
					break
				}
			}
			if !memberMatches {
				matches = false
			}
		}

		// Filter by locations if provided
		if matches && len(selectedLocations) > 0 && selectedLocations[0] != "" {
			locationMatch := false
			for _, artistLocation := range artist.Loca.Locations {
				for _, selectedLocation := range selectedLocations {
					if artistLocation == selectedLocation {
						locationMatch = true
						break
					}
				}
				if locationMatch {
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

	// Create a template data structure for rendering results.
	data := struct {
		Artists []models.Artist
	}{
		Artists: filteredArtists,
	}

	// Render the filter results template.
	if err := renderTemplate(w, "filter.html", data); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}