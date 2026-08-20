package internal

import (
	"context"
	"database/sql"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (r *MovieRepository) PostMovie(ctx context.Context, movie models.Movie) (models.Movie, error) {

	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO movies (title, description, release_date) VALUES (?, ?, ?)`,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
	)

	if err != nil {
		return models.Movie{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Movie{}, err
	}

	movie.ID = int(id)
	return movie, nil
}

func (r *MovieRepository) GetMovies() ([]models.Movie, error) {

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

func (r *MovieRepository) GetMovieByID(id int) (models.Movie, error) {

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

func (r *MovieRepository) PatchMovie(ctx context.Context, movie models.Movie, id int) error {

	res, err := r.db.ExecContext(
		ctx,
		`
	UPDATE movies SET title = ?, description = ?, release_date = ? WHERE id = ?
	`,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		id,
	)
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
