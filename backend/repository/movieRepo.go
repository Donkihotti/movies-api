package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"fmt"
	
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
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

func (r *MovieRepository) GetMovies(filters models.MovieFilters) ([]models.MovieReq, errs.ErrorStruct) {

	movies := []models.MovieReq{}
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
		errStruct.ErrMsg = errors.New("genre not found")
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
		var movie models.MovieReq
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

		genreRows, err := r.db.Query(
		`
		SELECT g.id, g.name
		FROM genres g 
		JOIN movie_genres mg ON mg.genre_id = g.id 
		WHERE mg.movie_id = ? 
		`, movie.ID)

		for genreRows.Next() {
		var moviegenre models.Genre
		err := genreRows.Scan(&moviegenre.ID, &moviegenre.Genre)
		if err != nil {
		return movies, errs.ErrorStruct{}
		}
		movie.Genres = append(movie.Genres, moviegenre)
		}

		actorRows, err := r.db.Query(
		`
		SELECT a.id, a.name, a.birth_date
		FROM actors a 
		JOIN movie_actors ag ON ag.actor_id = a.id 
		WHERE ag.movie_id = ? 
		`, movie.ID)

		for actorRows.Next() {
		var actor models.Actor
		err := actorRows.Scan(&actor.ID, &actor.Name, &actor.BirthDate)
		if err != nil {
		return movies, errs.ErrorStruct{}
		}
		movie.Actors = append(movie.Actors, actor)
		}

	movies = append(movies, movie)
	}
	return movies, errs.ErrorStruct{} 
}

func (r *MovieRepository) GetMovieByID(id int) (models.MovieReq, errs.ErrorStruct) {

	var movie models.MovieReq

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
			return models.MovieReq{}, errStruct 
		}
		errStruct.ErrType = errs.ServerError
		errStruct.ErrMsg = errors.New("something went wrong getting movie") 
		return models.MovieReq{}, errStruct 
	}

	rows, err := r.db.Query(
	`SELECT g.id, g.name
	FROM genres g 
	JOIN movie_genres mg ON mg.genre_id = g.id
	WHERE mg.movie_id = ? 
	`,
	id,
	) 
	if err != nil {
	errStruct := errs.ErrorStruct{}
	return models.MovieReq{}, errStruct
	}

	for rows.Next() {
	var genre models.Genre	
	err := rows.Scan(
	&genre.ID,
	&genre.Genre,
	)

	if err != nil {
		errStruct := errs.ErrorStruct{}
		return models.MovieReq{}, errStruct 
	}
	movie.Genres = append(movie.Genres, genre)
	}

	rows, err = r.db.Query(
	`SELECT a.id, a.name, a.birth_date
	FROM actors a 	
	JOIN movie_actors ag ON ag.actor_id = a.id
	WHERE ag.movie_id = ? 
	`,
	id,
	)
	if err != nil {
		errStruct := errs.ErrorStruct{}
		return models.MovieReq{}, errStruct 
	}

	for rows.Next() {
	var actor models.Actor
	err = rows.Scan(
	&actor.ID,
	&actor.Name,
	&actor.BirthDate,
	)
	if err != nil {
		errStruct := errs.ErrorStruct{}
		return models.MovieReq{}, errStruct 
	}
	movie.Actors = append(movie.Actors, actor)
	}

	return movie, errs.ErrorStruct{} 
}

func (r *MovieRepository) PatchMovie(ctx context.Context, movie models.PatchMovieReq, id int) (models.Movie, errs.ErrorStruct) {

	errStruct := errs.NewErrorStruct(
		errs.ServerError,
		errors.New("something went wrong updating movies"),
	)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
	return models.Movie{}, errStruct
	}

	defer tx.Rollback()

	patchedMovie := models.Movie{}

	query := `UPDATE movies SET title = COALESCE(?, title),
			  description = COALESCE(?, description),
			  release_date = COALESCE(?, release_date),
			  duration = COALESCE(?, duration)
			  WHERE id = ?
			  RETURNING id, title, description, release_date, duration`

	row := tx.QueryRowContext(ctx, query, 
				movie.Title, 
				movie.Description, 
				movie.ReleaseDate, 
				movie.Duration, id)

	err = row.Scan(&patchedMovie.ID,
					&patchedMovie.Title,
					&patchedMovie.Description,
					&patchedMovie.ReleaseDate,
					&patchedMovie.Duration,
					)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			errStruct.ErrType = errs.NotFound
			errStruct.ErrMsg = fmt.Errorf("Movie not found with id: %v", id)
			return models.Movie{}, errStruct
		}
		return models.Movie{}, errStruct
	}

	if len(movie.Actors) > 0 {
		
		_, err := tx.ExecContext(
		ctx,
		`DELETE from movie_actors WHERE movie_id = ?`,
		id)
		if err != nil {
		return models.Movie{}, errStruct
		}

		for _, actorId:= range movie.Actors {
	
		_, err := tx.ExecContext(
		ctx, 
		`INSERT INTO movie_actors (movie_id, actor_id) 
		VALUES (?,?)
		`, 
		id, actorId)
		if err != nil {
		errStruct := errs.ErrorStruct{
		errs.BadRequest,
		errors.New("actor id doesn't exist in database"),
		}
		return models.Movie{}, errStruct
		}
		}
	}

	rows, err := tx.QueryContext(
	ctx, 
	`
	SELECT a.id
	FROM actors a
	JOIN movie_actors mg ON mg.actor_id = a.id
	WHERE mg.movie_id = ? 
	`, id)

	if err != nil {
	return models.Movie{}, errStruct
	}

	for rows.Next() {
	var actorId int
	
	err := rows.Scan(
	&actorId,
	) 
	if err != nil {
	return models.Movie{}, errStruct
	}
	patchedMovie.Actors = append(patchedMovie.Actors, actorId)
	}

	if err := rows.Err(); err != nil {
		return models.Movie{}, errStruct
	}

	if len(movie.Genres) > 0 {
		
		_, err := tx.ExecContext(
		ctx,
		`DELETE from movie_genres WHERE movie_id = ?`,
		id)
		if err != nil {
		return models.Movie{}, errStruct
		}

		for _, genreId := range movie.Genres {
	
		_, err := tx.ExecContext(
		ctx, 
		`INSERT INTO movie_genres (movie_id, genre_id) 
		VALUES (?,?)
		`, 
		id, genreId)
		if err != nil {
		errStruct := errs.ErrorStruct{
		errs.BadRequest,
		errors.New("genre id doesn't exist in database"),
		}
		return models.Movie{}, errStruct
		}
		}
	}

	rows, err = tx.QueryContext(
	ctx, 
	`
	SELECT g.id
	FROM genres g
	JOIN movie_genres mg ON mg.genre_id = g.id
	WHERE mg.movie_id = ? 
	`, id)

	if err != nil {
	return models.Movie{}, errStruct
	}

	for rows.Next() {
	var genreId int
	
	err := rows.Scan(
	&genreId,
	) 
	if err != nil {
	return models.Movie{}, errStruct
	}
	patchedMovie.Genres = append(patchedMovie.Genres, genreId)
	}

	if err := rows.Err(); err != nil {
		return models.Movie{}, errStruct
	}

	err = tx.Commit() 
	if err != nil {
		return models.Movie{}, errStruct
	}

	return patchedMovie, errs.ErrorStruct{}
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
