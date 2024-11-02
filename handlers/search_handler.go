package handlers

import (
	"groupie/models"
	"net/http"
	"strconv"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get search query
	query := r.FormValue("q")
	queryLower := strings.ToLower(query)

	// Get filter values
	creationDateStr := r.FormValue("creationDate")
	firstAlbumYearFromStr := r.FormValue("firstAlbumYearFrom")
	firstAlbumYearToStr := r.FormValue("firstAlbumYearTo")
	selectedMemberCounts := r.Form["memberCount"]
	selectedLocations := r.Form["concertLocations"]

	var filteredArtists []models.Artist
	for _, artist := range artists {
		matches := true

		// Search query check
		if !strings.Contains(strings.ToLower(artist.Name), queryLower) {
			matches = false
		}

		// Creation date filter
		if creationDateStr != "" {
			creationDate, _ := strconv.Atoi(creationDateStr)
			if artist.CreationDate != creationDate {
				matches = false
			}
		}

		// First album year range filter
		if firstAlbumYearFromStr != "" || firstAlbumYearToStr != "" {
			// Assuming the firstAlbum contains the year; you might need to adjust this logic
			firstAlbumYear := extractYearFromFirstAlbum(artist.FirstAlbum)

			if firstAlbumYearFromStr != "" {
				firstAlbumYearFrom, _ := strconv.Atoi(firstAlbumYearFromStr)
				if firstAlbumYear < firstAlbumYearFrom {
					matches = false
				}
			}

			if firstAlbumYearToStr != "" {
				firstAlbumYearTo, _ := strconv.Atoi(firstAlbumYearToStr)
				if firstAlbumYear > firstAlbumYearTo {
					matches = false
				}
			}
		}

		// Member count filter
		if len(selectedMemberCounts) > 0 {
			matchedCount := false
			for _, count := range selectedMemberCounts {
				memberCount, err := strconv.Atoi(count)
				if err == nil && len(artist.Members) == memberCount {
					matchedCount = true
					break
				}
			}
			if !matchedCount {
				matches = false
			}
		}

		// Concert locations filter
		if len(selectedLocations) > 0 {
			matchedLocation := false
			locations := strings.Split(artist.Locations, ",") // Split locations into a slice
			for _, loc := range selectedLocations {
				for _, artistLoc := range locations {
					if strings.TrimSpace(artistLoc) == loc {
						matchedLocation = true
						break
					}
				}
				if matchedLocation {
					break
				}
			}
			if !matchedLocation {
				matches = false
			}
		}

		if matches {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// Create a template data structure
	data := struct {
		Artists []models.Artist
		Query   string
	}{
		Artists: filteredArtists,
		Query:   query,
	}

	// Render the search results template
	if err := renderTemplate(w, "search.html", data); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
}

// Helper function to extract the year from the first album string
func extractYearFromFirstAlbum(firstAlbum string) int {
	// Logic to extract the year from the firstAlbum string
	// Here, we'll assume the year is embedded in the string in a known format
	year, _ := strconv.Atoi(firstAlbum) // Placeholder logic
	return year
}
