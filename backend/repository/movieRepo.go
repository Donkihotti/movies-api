package repository

import (
	"context"
	"database/sql"
	"errors"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
	"strings"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

// POST A MOVIE
func (r *MovieRepository) PostMovie(ctx context.Context, req models.Movie) (models.Movie, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Movie{}, errs.ServerError
	}

	defer tx.Rollback()

	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO movies (title, description, release_date, duration)
	VALUES (?, ?, ?, ?)
	`,
		req.Title,
		req.Description,
		req.ReleaseDate,
		req.Duration,
	)
	if err != nil {
		return models.Movie{}, errs.ServerError
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Movie{}, err
	}

	req.ID = int(id)
	genreIds := req.Genres
	actorIds := req.Actors

	for _, genreId := range genreIds {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO movie_genres (movie_id, genre_id)
		VALUES (?, ?)
		`,
			req.ID,
			genreId,
		)
		if err != nil {
			return models.Movie{}, errs.ServerError
		}
	}

	for _, actorId := range actorIds {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO movie_actors (movie_id, actor_id)
		VALUES (?, ?)
		`,
			req.ID,
			actorId,
		)
		if err != nil {
			return models.Movie{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return models.Movie{}, errs.ServerError
	}

	return req, nil
}

// GET ALL MOVIES
func (r *MovieRepository) GetMovies(filters models.MovieFilters) ([]models.Movie, error) {

	movies := []models.Movie{}
	var conditions []string
	args := []any{}

	query := `SELECT m.id, m.title, m.description, m.release_date, m.duration FROM movies m`

	if filters.GenreID != nil {
		query += ` JOIN movie_genres gm ON gm.movie_id = m.id`
		conditions = append(conditions, "gm.genre_id = ?")
		args = append(args, *filters.GenreID)
	}

	if filters.ActorID != nil {
		query += ` JOIN movie_actors am ON am.movie_id = m.id`
		conditions = append(conditions, "am.actor_id = ?")
		args = append(args, *filters.ActorID)
	}

	if filters.ReleaseYear != nil {
		conditions = append(conditions, "m.release_date LIKE ?")
		args = append(args, "%"+*filters.ReleaseYear+"%")
	}

	//	if filters.Duration != nil {
	//		conditions = append(conditions, "m.duration = ?")
	//		args = append(args, "%" + *filters.Duration + "%")
	//	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, errs.ServerError
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.ReleaseDate,
			&movie.Duration,
		)
		if err != nil {
			return nil, errs.ServerError
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

// GET MOVIE BY ID
func (r *MovieRepository) GetMovieByID(id int) (models.Movie, error) {

	var movie models.Movie

	row := r.db.QueryRow(
		`SELECT id, title, description, release_date, duration
		FROM movies WHERE id = ?`, id)

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.ReleaseDate,
		&movie.Duration,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Movie{}, errs.NotFound
		}
		return models.Movie{}, errs.ServerError
	}
	return movie, nil
}

// PATCH MOVIE
func (r *MovieRepository) PatchMovie(ctx context.Context, movie models.PatchMovieReq, id int) error {

	res, err := r.db.ExecContext(
		ctx,
		`
    UPDATE movies SET title = COALESCE(?, title), description = COALESCE(?, description), release_date = COALESCE(?, release_date), duration = COALESCE(?, duration) WHERE id = ?
    `,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.Duration,
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

func (r *MovieRepository) DeleteForceMovie(ctx context.Context, id int) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
	return err
	}

	defer tx.Rollback()

	query := `DELETE FROM movie_genres WHERE movie_id = ?`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
	return err
	}

	query = `DELETE FROM movies WHERE id = ?`

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


	return tx.Commit()
}

func (r *MovieRepository) DeleteMovie(ctx context.Context, id int) error {

	query := `DELETE FROM movies WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errs.ServerError
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("Error fetching movies")
		return errs.ServerError
	}

	if rows == 0 {
		return errs.NotFound
	}

	return nil
}

// GET ACTORS FROM A GIVEN MOVIE
func (r *MovieRepository) MovieActors(ctx context.Context, id int) ([]models.Actor, error) {

	query := `SELECT a.id, a.name, a.birth_date FROM actors a JOIN movie_actors ma ON ma.actor_id = a.id WHERE ma.movie_id = ?`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return []models.Actor{}, errs.ServerError
	}
	defer rows.Close()
	actors := []models.Actor{}

	for rows.Next() {
		var actor models.Actor

		err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			return []models.Actor{}, errs.ServerError
		}
		actors = append(actors, actor)
	}
	return actors, nil
}
