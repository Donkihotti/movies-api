package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
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

func (r *GenreRepository) GetGenres(ctx context.Context) ([]models.Genre, errs.ErrorStruct) {

	errStruct := errs.NewErrorStruct(
	errs.ServerError,
	errors.New("something went wrong fetching genres"),
	)

	genres := []models.Genre{}

	rows, err := r.db.Query("SELECT id, name FROM genres")
	if err != nil {
		log.Println(err)
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
			return nil, errStruct
		}
		genres = append(genres, genre)
	}
	return genres, errs.ErrorStruct{}
}


func (r *GenreRepository) GetGenreByID(ctx context.Context, id int) (models.Genre, errs.ErrorStruct) {

	var genre models.Genre
	query := `SELECT id, name FROM genres WHERE id = ?`
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
	}
	
	return genre, errs.ErrorStruct{}
}


func (r *GenreRepository) PostGenre(ctx context.Context, req models.Genre) (models.Genre, errs.ErrorStruct) {

	errStruct := errs.NewErrorStruct(
	errs.ServerError,
	errors.New("something went wrong creating a genre"),
	)

	res, err := r.db.ExecContext(
	ctx,
	`INSERT INTO genres (name) VALUES (?)`,
	req.Genre,
	)

	if err != nil {
		log.Println(err)
		return models.Genre{}, errStruct
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return models.Genre{}, errStruct
	}
	req.ID = int(id)
	return req, errs.ErrorStruct{}
}


func (r *GenreRepository) DeleteForceGenre(ctx context.Context, id int) errs.ErrorStruct {

	errStruct := errs.NewErrorStruct(
		errs.ServerError,
		errors.New("something went wrong deleting a genre"),
	)
	
	tx, err := r.db.BeginTx(ctx, nil) 
	if err != nil {
		log.Println(err)
		return errStruct 
	}

	defer tx.Rollback()
	
	query := `DELETE FROM movie_genres WHERE genre_id = ?`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return errStruct 
	}

	query = `DELETE FROM genres WHERE id = ?`

	res, err := tx.ExecContext(ctx, query, id) 
	if err != nil {
		log.Println(err)
		return errStruct 
	}

	rows, err := res.RowsAffected() 
	if err != nil {
		log.Println(err)
		return errStruct 
	}

	if rows == 0 {
		errStruct := errs.NewErrorStruct(
			errs.NotFound,
			fmt.Errorf("could not delete genre with id: %v", id),
		)
		return errStruct 
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong deleting a genre"),
		)
	}

	return errs.ErrorStruct{} 
}


func (r *GenreRepository) DeleteGenre(ctx context.Context, id int) errs.ErrorStruct {

	errStruct := errs.NewErrorStruct(
		errs.ServerError,
		errors.New("something went wrong deleting a genre"),
	)

	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM genres WHERE id = ?`,
		id,
	)
	if err != nil {
		log.Println(err)
		return errStruct 
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return errStruct 
	}

	if rows == 0 {
			errStruct = errs.NewErrorStruct(
				errs.NotFound,
				fmt.Errorf("no matching genres with id: %v", id),
			)
		return errStruct
	}
	return errs.ErrorStruct{} 
}


func (r *GenreRepository) PatchGenre(ctx context.Context, newGenre models.Genre, id int) (models.Genre, errs.ErrorStruct) {

	errStruct := errs.NewErrorStruct(
		errs.ServerError,
		errors.New("something went wrong updating a genre"),
	)

	patchedGenre := models.Genre{}
	
	query := "UPDATE genres SET name = ? WHERE id = ? RETURNING id, name"
	err := r.db.QueryRowContext(ctx, query, newGenre.Genre, id).Scan(&patchedGenre.ID, &patchedGenre.Genre)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			errStruct.ErrType = errs.NotFound
			errStruct.ErrMsg = fmt.Errorf("Genre does not exist with id %v", id)
			return models.Genre{}, errStruct
		}
		return models.Genre{}, errStruct
	}
	return patchedGenre, errs.ErrorStruct{}
}
