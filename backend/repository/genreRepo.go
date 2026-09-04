package repository

import (
	"context"
	"database/sql"
	 
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


func (r *GenreRepository) GetGenres(ctx context.Context) ([]models.Genre, error) {

	genres := []models.Genre{}
	rows, err := r.db.Query(
		`
	SELECT id, name		
	FROM genres
	`,
	)
	if err != nil {
		return nil, errs.ServerError
	}
	defer rows.Close()

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(
			&genre.ID,
			&genre.Genre,
		)
		if err != nil {
			return nil, errs.ServerError
		}
		genres = append(genres, genre)
	}

	return genres, nil
}


func (r *GenreRepository) GetGenre(ctx context.Context, id int) (models.Genre, error) {

	var genre models.Genre
	rows, err := r.db.Query(
	`
	SELECT id, name		
	FROM genres
	WHERE id = ? 
	`, 
	id,
	)
	if err != nil {
	return models.Genre{}, err
	}
	defer rows.Close()

	for rows.Next() {
	err := rows.Scan(
	&genre.ID,
	&genre.Genre,
	)
	if err != nil {
	return models.Genre{}, err 
	}
	
	}
	
	return genre, nil
}


func (r *GenreRepository) PostGenre(ctx context.Context, genre models.Genre) (models.Genre, error) {

	res, err := r.db.ExecContext(
		ctx,
		`
	INSERT INTO genres (name) VALUES (?)
	`,
		genre.Genre,
	)
	if err != nil {
		return models.Genre{}, errs.ServerError
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Genre{}, errs.ServerError
	}
	genre.ID = int(id)
	return genre, nil
}


func (r *GenreRepository) DeleteForceGenre(ctx context.Context, id int) error {

	tx, err := r.db.BeginTx(ctx, nil) 
	if err != nil {
	return err
	}

	defer tx.Rollback()
	
	query := `DELETE FROM movie_genres WHERE genre_id = ?`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
	return err
	}

	query = `DELETE FROM genres WHERE id = ?`

	res, err := tx.ExecContext(ctx, query, id) 
	if err != nil {
	return err
	}

	rows, err := res.RowsAffected() 
	if err != nil {
	return err
	}

	if rows == 0 {
	return sql.ErrNoRows
	}

	return nil
}


func (r *GenreRepository) DeleteGenre(ctx context.Context, id int) error {

	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM genres WHERE id = ?`,
		id,
	)
	if err != nil {
		return errs.ServerError
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return errs.ServerError
	}

	if rows == 0 {
		return errs.NotFound
	}
	return nil
}


func (r *GenreRepository) PatchGenre(ctx context.Context, newGenre models.Genre, id int) (models.Genre, error) {

	res, err := r.db.ExecContext(
		ctx,
		`UPDATE genres SET name = ? WHERE id = ?`,
		newGenre.Genre,
		id,
	)
	
	if err != nil {
		return models.Genre{}, errs.ServerError
	}

	var patchedGenre models.Genre
	err = r.db.QueryRow(
	`SELECT name
	FROM genres 
	WHERE id = ?`,
	id,
	).Scan(
	&patchedGenre.Genre,
	)

	if err != nil {
	return models.Genre{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return models.Genre{}, errs.ServerError
	}
	if rows == 0 {
		return models.Genre{}, errs.NotFound
	}

	return patchedGenre, nil

}
