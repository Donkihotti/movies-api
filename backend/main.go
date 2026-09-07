package main

import (
	"log"
	"net/http"
	"os"

	"gitea.kood.tech/timdanielfiander/movies-api.git/api"
	"gitea.kood.tech/timdanielfiander/movies-api.git/database"
	"gitea.kood.tech/timdanielfiander/movies-api.git/repository"
	"gitea.kood.tech/timdanielfiander/movies-api.git/service"
	"gitea.kood.tech/timdanielfiander/movies-api.git/api/handlers"
)

func main() {

	args := os.Args
	var populateDb bool

	if len(args) > 1 && args[1] == "-db" {
	populateDb = true	
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = database.RunMigrations(db, populateDb)
	if err != nil {
		log.Fatal(err)
	}

	movieRepo := repository.NewMovieRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	actorRepo := repository.NewActorRepository(db)

	movieService := service.NewMovieService(movieRepo)
	genreService := service.NewGenreService(genreRepo)
	actorService := service.NewActorService(actorRepo)

	movieHandler := &handlers.MovieHandler{
		Service: movieService,
	}
	genreHandler := &handlers.GenreHandler{
		Service: genreService,
	}
	actorHandler := &handlers.ActorHandler{
		Service: actorService,
	}

	router := api.NewRouter(
	movieHandler, 
	genreHandler, 
	actorHandler,
	)

	log.Println("Server running at port :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Starting server failed", err)
	}
}
