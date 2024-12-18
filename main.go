package main

import (
	"fmt"
	"groupie/handlers"
	"net/http"
)

func main() {

	http.HandleFunc("/",handlers.MainHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/assets/", handlers.AssetsHandler)
	http.HandleFunc("/artists/", handlers.ArtistDetailHandler)
	http.HandleFunc("/filter", handlers.FilterHandler)

	fmt.Println("Server is running at http://localhost:3000")
	
	
	err:=http.ListenAndServe(":3001",nil)
	if err!=nil{
        fmt.Println("Error starting server: ",err)
    }

}