package internal

import (
	"context"
	"strings"
	"database/sql"
	"log"
	"errors"
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

func (r *MovieRepository) PostMovie(ctx context.Context, req models.Movie) (models.Movie, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Movie{}, errs.ServerError
	}

	defer tx.Rollback()

	res, err := tx.ExecContext(
	ctx, 
	`INSERT INTO movies (title, description, release_date)
	VALUES (?, ?, ?)
	`,
	req.Title,
	req.Description,
	req.ReleaseDate,
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

//GET ALL MOVIES
func (r *MovieRepository) GetMovies(filters models.MovieFilters) ([]models.Movie, error){

	movies := []models.Movie{}
	var conditions []string
	args := []any{}

	query := `SELECT m.id, m.title, m.description, m.release_date FROM movies m`

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
		args = append(args, "%" + *filters.ReleaseYear + "%")
	}

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
		)
		if err != nil {
			return nil, errs.ServerError
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.Movie{}, errs.NotFound
		}
		return models.Movie{}, errs.ServerError
	}
	return movie, nil
}

//PATCH MOVIE
func (r *MovieRepository) PatchMovie(ctx context.Context, movie models.PatchMovieReq, id int) error {

	
    res, err := r.db.ExecContext(
        ctx,
        `
    UPDATE movies SET title = COALESCE(?, title), description = COALESCE(?, description), release_date = COALESCE(?, release_date) WHERE id = ?
    `,
        movie.Title,		//voidaan antaa pointerit sellaisenaan, sql osaa tulkita niitä itse
        movie.Description,
        movie.ReleaseDate,
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

//DELETE MOVIE
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

