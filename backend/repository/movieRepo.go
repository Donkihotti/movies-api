package repository

import (
	"context"
	"database/sql"
	"errors"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
	"strings"
	"fmt"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (r *MovieRepository) PostMovie(ctx context.Context, req models.Movie) (models.Movie, errs.ErrorStruct) {

	errStruct := errs.NewErrorStruct(
		errs.ServerError,
		errors.New("something went wrong creating a movie"),
	)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return models.Movie{}, errStruct 
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
		log.Println(err)
		return models.Movie{}, errStruct 
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return models.Movie{}, errStruct
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
			log.Println(err)
			return models.Movie{}, errStruct 
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
			log.Println(err)
			return models.Movie{}, errStruct 
		}
	}

	if err := tx.Commit(); err != nil {
			log.Println(err)
			return models.Movie{}, errStruct 
	}

	return req, errs.ErrorStruct{} 
}

func (r *MovieRepository) GetMovies(filters models.MovieFilters) ([]models.Movie, errs.ErrorStruct) {

	movies := []models.Movie{}
	var conditions []string
	args := []any{}

	errStruct := errs.NewErrorStruct(
		errs.ServerError,
		errors.New("something went wrong getting movies"),
	)

	query := `SELECT m.id, m.title, m.description, m.release_date, m.duration FROM movies m`

	if filters.GenreID != nil {

		var exists bool

		err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM genres WHERE id = ?)`,
		*filters.GenreID, 
		).Scan(&exists) 

		if err != nil {
		return nil, errStruct 
		}

		if !exists {
		errStruct.ErrType = errs.NotFound
		errStruct.ErrMsg = errors.New("movie does not exist")
		return nil, errStruct 
		}

		query += ` JOIN movie_genres gm ON gm.movie_id = m.id`
		conditions = append(conditions, "gm.genre_id = ?")
		args = append(args, *filters.GenreID)
	}

	if filters.ActorID != nil {

		var exists bool

		err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM actors WHERE id = ?)`,
		*filters.ActorID, 
		).Scan(&exists) 

		if err != nil {
		return nil, errStruct 
		}

		if !exists {
		errStruct.ErrType = errs.NotFound
		errStruct.ErrMsg = errors.New("actor not found")
		return nil, errStruct 
		}

		query += ` JOIN movie_actors am ON am.movie_id = m.id`
		conditions = append(conditions, "am.actor_id = ?")
		args = append(args, *filters.ActorID)
	}

	if filters.ReleaseYear != nil {
		conditions = append(conditions, "m.release_date LIKE ?")
		args = append(args, *filters.ReleaseYear+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Println(err)
		return nil, errStruct 
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
			log.Println(err)
			return nil, errStruct 
		}
		movies = append(movies, movie)
	}
	return movies, errs.ErrorStruct{} 
}

func (r *MovieRepository) GetMovieByID(id int) (models.Movie, errs.ErrorStruct) {

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
		log.Println(err)
		errStruct := errs.ErrorStruct{}
		if errors.Is(err, sql.ErrNoRows) {
			errStruct.ErrType = errs.NotFound
			errStruct.ErrMsg = fmt.Errorf("could not get movie with id: %v", id)
			return models.Movie{}, errStruct 
		}
		errStruct.ErrType = errs.ServerError
		errStruct.ErrMsg = errors.New("something went wrong getting movie") 
		return models.Movie{}, errStruct 
	}

	return movie, errs.ErrorStruct{} 
}

func (r *MovieRepository) PatchMovie(ctx context.Context, movie models.PatchMovieReq, id int) (models.Movie, errs.ErrorStruct) {

	errStruct := errs.NewErrorStruct(
		errs.ServerError,
		errors.New("something went wrong updating movies"),
	)

	var patchedMovie models.Movie
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

	err = r.db.QueryRow(
	`
	SELECT title, description, release_date, duration 
	FROM movies 
	WHERE id = ?`,
	id,
	).Scan(
	&patchedMovie.Title, 
	&patchedMovie.Description, 
	&patchedMovie.ReleaseDate, 
	&patchedMovie.Duration, 
	)

	if err != nil {
		log.Println(err)
		return models.Movie{}, errStruct 
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return models.Movie{}, errStruct 
	}

	if rows == 0 {
		log.Println(err)
		errStruct.ErrType = errs.NotFound
		errStruct.ErrMsg = fmt.Errorf("could not update movie with id: %v", id) 
		return models.Movie{}, errStruct 
	}
	return patchedMovie, errs.ErrorStruct{} 
}


func (r *MovieRepository) DeleteForceMovie(ctx context.Context, id int) errs.ErrorStruct {

	errStruct := errs.NewErrorStruct(
		errs.ServerError,
		errors.New("something went wrong deleting movie"),
	)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return errStruct 
	}

	defer tx.Rollback()

	query := `DELETE FROM movie_genres WHERE movie_id = ?`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return errStruct 
	}

	query = `DELETE FROM movies WHERE id = ?`

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
		errStruct = errs.NewErrorStruct(
			errs.NotFound,
			fmt.Errorf("could not delete movie with id: %v", id),
		)
		return errStruct 
	}
	
	 if err := tx.Commit(); err != nil {
		log.Println(err)
		return errStruct
	}

	return errs.ErrorStruct{} 
}

func (r *MovieRepository) DeleteMovie(ctx context.Context, id int) errs.ErrorStruct {

	errStruct := errs.NewErrorStruct(
		errs.ServerError,
		errors.New("something went wrong deleting a movie"),
	)

	query := `DELETE FROM movies WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, id)
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
			fmt.Errorf("could not delete movie with id: %v", id),
		)
		return errStruct 
	}

	return errs.ErrorStruct{} 
}

func (r *MovieRepository) MovieActors(ctx context.Context, id int) ([]models.Actor, errs.ErrorStruct) {

	errStruct := errs.NewErrorStruct(
		errs.ServerError,
		errors.New("something went wrong getting actors"),
	)

	query := `SELECT a.id, a.name, a.birth_date FROM actors a JOIN movie_actors ma ON ma.actor_id = a.id WHERE ma.movie_id = ?`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return []models.Actor{}, errStruct 
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
			log.Println(err)
			return []models.Actor{}, errStruct 
		}
		actors = append(actors, actor)
	}
	return actors, errs.ErrorStruct{} 
}
