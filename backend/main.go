package main 

import (
	"net/http"
	"log"

	"gitea.kood.tech/timdanielfiander/movies-api.git/api"
)

func main() {
	
	router := api.NewRouter()

	log.Println("Server running at port :8080")
	err := http.ListenAndServe(":8080", router)
		if err != nil {
			log.Fatal("Starting server failed", err)	
		}

}


