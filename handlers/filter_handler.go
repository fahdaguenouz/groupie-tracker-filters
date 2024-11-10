package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"groupie/models"
)

// Helper function to reverse the user's input date format (YYYY-MM-DD -> DD-MM-YYYY)
func reverseDateFormat(date string) (string, error) {
	// Split the date into its components (YYYY, MM, DD)
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid date format")
	}

	// Reverse the format to DD-MM-YYYY
	reversedDate := fmt.Sprintf("%s-%s-%s", parts[2], parts[1], parts[0])

	return reversedDate, nil
}

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
	firstAlbumDate := r.FormValue("firstAlbumDate") // New filter for full date of first album
	members := r.Form["members"]  
	locations := r.Form["locations"]
	fmt.Println("User input for first album date:", firstAlbumDate)
	fmt.Println("User selected members filter:", members)
	fmt.Println("User selected locations filter:", locations)

	var filteredArtists []models.Artist

	// Filter artists by creation year if provided
	if creationYear != "" && creationYear != "1900" {
		// Parse creationYear to an integer
		year, err := strconv.Atoi(creationYear)
		if err != nil {
			// If the creationYear is not a valid integer, return an error
			ErrorHandler(w, r, http.StatusBadRequest)
			return
		}

		// Filter by creation year
		for _, artist := range artists {
			fmt.Printf("Artist: %s, CreationDate: %d, User's creation year input: %d\n", artist.Name, artist.CreationDate, year)

			if artist.CreationDate == year {
				filteredArtists = append(filteredArtists, artist)
			}
		}

		// If no artists match the creation year filter, print a debug message
		if len(filteredArtists) == 0 {
			fmt.Println("No artists found after filtering by creation year:", year)
		}
	} else {
		// If no creationYear filter is applied, use all artists
		filteredArtists = artists
	}

	// If there is a first album date filter, apply it after creation year filtering
	if firstAlbumDate != "" {
		// Reverse the user's date input (YYYY-MM-DD) to DD-MM-YYYY format
		reversedUserDate, err := reverseDateFormat(firstAlbumDate)
		if err != nil {
			fmt.Println("Error reversing user date:", err)
			ErrorHandler(w, r, http.StatusBadRequest)
			return
		}

		// Print the reversed user date for debugging
		fmt.Println("Reversed user date:", reversedUserDate)

		var tempArtists []models.Artist
		fmt.Println("Filtered artists before first album date filtering:", len(filteredArtists))

		// Loop through the filtered list and apply first album date filtering
		for _, artist := range filteredArtists {
			// Check if the FirstAlbum field is empty or contains an invalid value
			if artist.FirstAlbum == "" {
				fmt.Println("Empty FirstAlbum for artist:", artist.Name)
				continue
			}

			// Compare the reversed user date with the artist's FirstAlbum
			if artist.FirstAlbum == reversedUserDate {
				tempArtists = append(tempArtists, artist)
			}
		}

		// Update filteredArtists with the result of first album date filtering
		filteredArtists = tempArtists

		// If no artists match the first album date, we need to debug this
		if len(filteredArtists) == 0 {
			fmt.Println("No artists found after filtering by first album date")
		}
	}
	// If there is a members filter, apply it after the first album date filter
	if len(members) > 0 {
		var tempArtists []models.Artist
		for _, artist := range filteredArtists {
			// Check the number of members for the artist
			memberCount := len(artist.Members)

			// Apply the filter for number of members selected by the user
			for _, member := range members {
				// Convert the member count from string to integer
				memCount, err := strconv.Atoi(member)
				if err != nil {
					fmt.Println("Error converting member count:", err)
					continue
				}

				// If the artist's member count matches, add to the filtered list
				if (memCount == 5 && memberCount >= 5) || (memCount < 5 && memberCount == memCount) {
					tempArtists = append(tempArtists, artist)
					break
				}
			}
		}

		// Update filteredArtists with the result of the member count filtering
		filteredArtists = tempArtists

		// If no artists match the member count, print a debug message
		if len(filteredArtists) == 0 {
			fmt.Println("No artists found after filtering by number of members")
		}
	}

// Apply location filter (user selected location)
if len(locations) > 0 {
	var tempArtists []models.Artist
	selectedLocation := locations[0] // Only one location selected by the user

	// Print the selected location for debugging
	fmt.Println("Selected location:", selectedLocation)

	// Loop through all filtered artists
	for _, artist := range filteredArtists {
		if err := GetForeigenData(&artist); err != nil {
			fmt.Println("Error fetching foreign data:", err)
			return // Skip this artist if there's an error
		}
		// Print artist's locations for debugging
		fmt.Println("Artist:", artist.Name, "Locations:", artist.Loca.Locations)

		// Check if the artist has the selected location in their Loca.Locations slice
		for _, artistLocation := range artist.Loca.Locations {
			// Trim spaces for accurate comparison
			artistLocation = strings.TrimSpace(artistLocation)
			selectedLocation = strings.TrimSpace(selectedLocation)

			// Print the locations being compared for debugging
			fmt.Println("Comparing with location:", artistLocation, "==", selectedLocation)

			// If the artist's location matches the selected location, add to tempArtists
			if artistLocation == selectedLocation {
				tempArtists = append(tempArtists, artist)
				break // No need to check further locations for this artist
			}
		}
	}

	// Update filteredArtists with the result of location filtering
	filteredArtists = tempArtists

	// Debug output to check if any artists were found after filtering by location
	if len(filteredArtists) == 0 {
		fmt.Println("No artists found after filtering by location")
	} else {
		fmt.Println("Artists found after location filtering:", len(filteredArtists))
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
