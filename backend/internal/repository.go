package internal 

import (
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"database/sql"
	"context"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,		
	}
}

func (r *Repository) PostMovie(ctx context.Context, req models.CreateMovieReq) (models.Movie, error) {
	
	err := r.db.QueryRowContext(

	)

	if err != nil {
	return models.Movie{}, err
	}
	return movie, nil
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

func (r *Repository) GetMovieByID(id int) (models.Movie, error) {

	var movie models.Movie
	
	row := r.db.QueryRow(
		`SELECT id, title, description, release_date
		FROM movies WHERE id = ?`, id)

	err := row.Scan(
	&movie.ID,
	&movie.Title,
	&movie.Description,
	&movie.ReleaseDate,
	)
	if err != nil {
		return models.Movie{}, err
	}
	return movie, nil
}


