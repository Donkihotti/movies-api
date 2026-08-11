package internal 

import (
	"database/sql"
	"context"

	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
)

type GenreRepository struct {
	db *sql.DB	
}

func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{
	db: db,
	}
}

func (r *GenreRepository) GetGenres(ctx context.Context) ([]models.Genre, error) {

	var genres []models.Genre
	rows, err := r.db.Query(
	`
	SELECT id, genre_name		
	FROM genres
	`, 
	)
	if err != nil {
	return nil, err
	}
	defer rows.Close()

	for rows.Next() {
	var genre models.Genre
	err := rows.Scan(
	&genre.ID, 
	&genre.Genre,
	)
	if err != nil {
	return nil, err 
	}
	genres = append(genres, genre)
	}

	return genres, nil
}


func (r *GenreRepository) PostGenre(ctx context.Context, genre models.Genre) (models.Genre, error) {
	
	res, err := r.db.ExecContext(
	ctx, 
	`
	INSERT INTO genres (genre_name) VALUES (?)
	`, 
	genre.Genre,
	)
	if err != nil { 
	return models.Genre{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
	return models.Genre{}, err
	}
	genre.ID = int(id)
	return genre, nil	
}


