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
			year, err := strconv.Atoi(creationYear) // Convert to int
			if err != nil || artist.CreationDate != year {
				matches = false
			}
		}

		// Filter by first album year if provided
		if firstAlbumYear != "" {
			firstAlbumDateParts := strings.Split(artist.FirstAlbum, "-") // Assuming FirstAlbum is in "YYYY-MM-DD" format
			if len(firstAlbumDateParts) > 0 {
				firstAlbumYearInt, err := strconv.Atoi(firstAlbumDateParts[0]) // Extract year from FirstAlbum
				if err != nil {
					matches = false // If there's an error converting the year, do not match
				} else {
					providedFirstAlbumYearInt, err := strconv.Atoi(firstAlbumYear) // Convert provided first album year to int
					if err != nil || firstAlbumYearInt != providedFirstAlbumYearInt { // Compare extracted year with provided year
						matches = false
					}
				}
			} else {
				matches = false // If the format is unexpected, do not match
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
                matches = false // Only set matches to false if no member criteria are met.
            }
        }

        // Filter by locations if selected
        if len(selectedLocations) > 0 && selectedLocations[0] != "" { // Check if any location is selected
            locationMatch := false
            for _, artistLocation := range artist.Loca.Locations {
                for _, selectedLocation := range selectedLocations {
                    if artistLocation == selectedLocation {
                        locationMatch = true
                        break
                    }
                }
                if locationMatch {
                    break // Exit outer loop if a match is found
                }
            }
            if !locationMatch {
                matches = false // Only set matches to false if no location criteria are met.
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
	// Render the filter results template
	if err := renderTemplate(w, "filter.html", data); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
