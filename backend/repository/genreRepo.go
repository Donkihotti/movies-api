package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
)

type GenreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{
		db: db,
	}
}

// GET ALL GENRES
func (r *GenreRepository) GetGenres(ctx context.Context) ([]models.Genre, errs.ErrorStruct) {

	genres := []models.Genre{}
	rows, err := r.db.Query("SELECT id, genre_name FROM genres")
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong fetching genres"),
		)
		return nil, errStruct
	}
	defer rows.Close()

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(
			&genre.ID,
			&genre.Genre,
		)
		if err != nil {
			log.Println(err)
			errStruct := errs.NewErrorStruct(
				errs.ServerError,
				errors.New("something went wrong fetching genres"),
			)
			return nil, errStruct
		}
		genres = append(genres, genre)
	}
	return genres, errs.ErrorStruct{}
}

// GET GENRE BY ID
func (r *GenreRepository) GetGenre(ctx context.Context, id int) (models.Genre, errs.ErrorStruct) {

	genre := models.Genre{}

	query := `SELECT id, genre_name FROM genres WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&genre.ID, &genre.Genre)

	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			errStruct := errs.NewErrorStruct(
				errs.NotFound,
				fmt.Errorf("no matching genres with id: %v", id),
			)
			return models.Genre{}, errStruct
		}
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			fmt.Errorf("something went wrong getting genre by id"),
		)
		return models.Genre{}, errStruct
	}
	return genre, errs.ErrorStruct{}
}

// still needs changing
// POST GENRE
func (r *GenreRepository) PostGenre(ctx context.Context, req models.Genre) (models.Genre, errs.ErrorStruct) {

	res, err := r.db.ExecContext(ctx, "INSERT INTO genres (genre_name) VALUES (?)", req.Genre)
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong creating a genre"),
		)
		return models.Genre{}, errStruct
	}
	id, err := res.LastInsertId()
	if err != nil {
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong creating a genre"),
		)
		return models.Genre{}, errStruct
	}
	req.ID = int(id)
	return req, errs.ErrorStruct{}
}

// DELETE GENRE
func (r *GenreRepository) DeleteGenre(ctx context.Context, id int) errs.ErrorStruct {

	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM genres WHERE id = ?`,
		id,
	)
	if err != nil {
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong deleting a genre"),
		)
		return errStruct 
	}

	rows, err := res.RowsAffected()

	if err != nil {
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong deleting a genre"),
		)
		return errStruct 
	}

	if rows == 0 {
			errStruct := errs.NewErrorStruct(
				errs.NotFound,
				fmt.Errorf("no matching genres with id: %v", id),
			)
		return errStruct
	}
	return errs.ErrorStruct{} 
}

// PATCH GENRE
func (r *GenreRepository) PatchGenre(ctx context.Context, newGenre models.Genre, id int) errs.ErrorStruct {

	res, err := r.db.ExecContext(
		ctx,
		`UPDATE genres SET genre_name = ? WHERE id = ?`,
		newGenre.Genre,
		id,
	)
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong updating a genre"),
		)
		return errStruct 
	}

	rows, err := res.RowsAffected()
	if err != nil {
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong updating a genre"),
		)
		return errStruct 
	}
	if rows == 0 {
		errStruct := errs.NewErrorStruct(
			errs.NotFound,
			fmt.Errorf("genre not found with id: %v", id), 
		)
		return errStruct 
	}
	return errs.ErrorStruct{}
}
