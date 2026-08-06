package internal 

import (

	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,		
	}
}

func (r *Repository) GetMovies() ([]models.Movie, error) {

	var movies []models.Movie
	rows, err := r.db.Query(
	`
	SELECT id, title, description, release_date
	FROM movies
	`, 
	)
	if err != nil {
	return nil, err
	}
	defer rows.Close()

	for rows.Next() {
	var movie models.Movie
	err := rows.Scan(
	&movie.ID, 
	&movie.Title, 
	&movie.Description, 
	&movie.ReleaseDate,
	) 
	if err != nil { 
	return nil, err
	}

	movies = append(movies, movie)
	}

	return movies, nil
}
