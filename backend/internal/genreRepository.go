package internal

import (
	"context"
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

func (r *GenreRepository) GetGenre(ctx context.Context, id int) ([]models.Movie, error) {

	var movies []models.Movie

	query := `SELECT m.id, m.title, m.description, m.release_date
	FROM movies AS m
	JOIN movie_genres AS mg ON mg.movie_id = m.id
	JOIN genres AS g ON g.id = mg.genre_id
	WHERE mg.genre_id = ? 
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
	var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate); err != nil {
		return nil, err
		}
	movies = append(movies, movie)
	}

	return movies, nil
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

func (r *GenreRepository) DeleteGenre(ctx context.Context, id int) error {

	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM genres WHERE id = ?`,
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

func (r *GenreRepository) PutGenre(ctx context.Context, newGenre models.Genre, id int) error {

	res, err := r.db.ExecContext(
		ctx,
		`UPDATE genres SET genre_name = ? WHERE id = ?`,
		newGenre.Genre,
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
