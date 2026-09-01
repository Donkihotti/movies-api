package main

import (
	"gitea.kood.tech/timdanielfiander/movies-api.git/api"
	"gitea.kood.tech/timdanielfiander/movies-api.git/database"
	"gitea.kood.tech/timdanielfiander/movies-api.git/repository"
	"gitea.kood.tech/timdanielfiander/movies-api.git/service"
	"log"
	"net/http"
)

func main() {

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = database.RunMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	movieRepo := repository.NewMovieRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	actorRepo := repository.NewActorRepository(db)

	movieService := service.NewMovieService(movieRepo)
	genreService := service.NewGenreService(genreRepo)
	actorService := service.NewActorService(actorRepo)

	handler := api.NewHandler(movieService, genreService, actorService)

	router := api.NewRouter(handler)

	log.Println("Server running at port :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Starting server failed", err)
	}
}
