package internal 

import (
	"database/sql"

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

func (r *GenreRepository) GetGenres() ([]models.Genre, error) {



	return nil, nil
}
