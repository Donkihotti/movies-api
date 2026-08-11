package main 

import (
	"net/http"
	"log"

	"gitea.kood.tech/timdanielfiander/movies-api.git/api"
	"gitea.kood.tech/timdanielfiander/movies-api.git/internal"
	"gitea.kood.tech/timdanielfiander/movies-api.git/database"
)

func main() {

	db, err := database.Connect()
	if err != nil {
	log.Fatal(err) 
	}

	defer db.Close()

	err = database.RunMigrations(db)

	//"the chain of command"
	repo := internal.NewRepository(db) 
	service := internal.NewService(repo)
	actorRepo := internal.NewActorRepository(db)
	actorService := internal.NewActorService(actorRepo)
	handler := api.NewHandler(service, actorService)



	router := api.NewRouter(handler)


	


	log.Println("Server running at port :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Starting server failed", err)	
	}
}

