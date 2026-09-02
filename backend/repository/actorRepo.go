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

type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{
		db: db,
	}
}

// GET ALL ACTORS
// done
func (r *ActorRepository) GetActors(name string) ([]models.Actor, errs.ErrorStruct) {

	actors := []models.Actor{}
	var rows *sql.Rows
	var errStruct errs.ErrorStruct
	var err error

	if name == "" {
		query := "SELECT id, name, birth_date FROM actors"
		rows, err = r.db.Query(query)
		if err != nil {
			log.Println(err)
			errStruct.ErrType = errs.ServerError
			errStruct.ErrMsg = errors.New("could not fetch actors")
			return actors, errStruct
		}
	} else {
		name = fmt.Sprintf("%%%s%%", name)
		query := "SELECT id, name, birth_date FROM actors WHERE name LIKE ?"
		rows, err = r.db.Query(query, name)
		if err != nil {
			log.Println(err)
			errStruct.ErrType = errs.ServerError
			errStruct.ErrMsg = errors.New("could not fetch actors")
			return actors, errStruct
		}
	}

	defer rows.Close()

	for rows.Next() {
		var actor models.Actor
		err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			log.Println(err)
			errStruct.ErrType = errs.ServerError
			errStruct.ErrMsg = errors.New("could not fetch actors")
			return nil, errStruct
		}
		actors = append(actors, actor)
	}

	return actors, errs.ErrorStruct{}
}

// GET ACTOR BY ID
// done
func (r *ActorRepository) GetActorByID(id int) (models.Actor, errs.ErrorStruct) {
	var actor models.Actor
	var errStruct errs.ErrorStruct

	row := r.db.QueryRow(
		`SELECT id, name, birth_date
		FROM actors WHERE id = ?`, id)

	err := row.Scan(
		&actor.ID,
		&actor.Name,
		&actor.BirthDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errStruct.ErrType = errs.NotFound
			errStruct.ErrMsg = fmt.Errorf("actor not found with id: %v", id)
			return models.Actor{}, errStruct
		}
		log.Println(err)
		return models.Actor{}, errs.ErrorStruct{
			ErrType: errs.ServerError,
			ErrMsg:  errors.New("something went wrong fetching actor by id"),
		}
	}

	return actor, errs.ErrorStruct{}
}

// POST AN ACTOR
// done
func (r *ActorRepository) PostActor(ctx context.Context, req models.Actor) (models.Actor, errs.ErrorStruct) {

	errStruct := errs.ErrorStruct{}

	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO actors (name, birth_date) VALUES (?, ?)`,
		req.Name,
		req.BirthDate,
	)

	if err != nil {
		log.Println(err)
		errStruct.ErrType = errs.ServerError
		errStruct.ErrMsg = errors.New("something went wrong creating an actor")
		return models.Actor{}, errStruct
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		errStruct.ErrType = errs.ServerError
		errStruct.ErrMsg = errors.New("something went wrong creating an actor")
		return models.Actor{}, errStruct
	}

	req.ID = int(id)
	return req, errs.ErrorStruct{}

}

// DELETE AN ACTOR
func (r *ActorRepository) DeleteActor(ctx context.Context, id int) errs.ErrorStruct {

	query := `DELETE FROM actors WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong deleting actor"),
		)
		return errStruct
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong deleting actor"),
		)
		return errStruct
	}

	if rows == 0 {
		errStruct := errs.NewErrorStruct(
			errs.NotFound,
			fmt.Errorf("no matching actors with id: %v", id),
		)

		return errStruct
	}

	return errs.ErrorStruct{}
}

// PATCH AN ACTOR
func (r *ActorRepository) PatchActor(ctx context.Context, req models.PatchActorReq, id int) errs.ErrorStruct {

	query := `UPDATE actors SET name = COALESCE(?, name), birth_date = COALESCE(?, birth_date) WHERE id = ?`
	row, err := r.db.ExecContext(ctx, query, req.Name, req.BirthDate, id)
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("error updating actor"),
		)
		return errStruct
	}

	rows, err := row.RowsAffected()
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("error updating actor"),
		)
		return errStruct
	}

	if rows == 0 {
		errStruct := errs.NewErrorStruct(
			errs.NotFound,
			fmt.Errorf("actor with id: %v not found", id),
		)
		return errStruct
	}

	return errs.ErrorStruct{}
}

// GET ACTORS BY NAME
// done
func (r *ActorRepository) GetActorsByName(ctx context.Context, Name string) ([]models.Actor, errs.ErrorStruct) {

	errStruct := errs.ErrorStruct{}
	actors := []models.Actor{}
	name := fmt.Sprintf("%%%s%%", Name)
	query := `SELECT id, name, birth_date FROM actors WHERE name LIKE ?`

	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		log.Println(err)
		errStruct.ErrType = errs.ServerError
		errStruct.ErrMsg = fmt.Errorf("error fetching actor by name: %v", Name)
		return actors, errStruct
	}
	defer rows.Close()

	for rows.Next() {
		var actor models.Actor
		err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			log.Println(err)
			errStruct.ErrType = errs.ServerError
			errStruct.ErrMsg = fmt.Errorf("error fetching actor by name: %v", Name)
			return []models.Actor{}, errStruct
		}
		actors = append(actors, actor)
	}
	return actors, errs.ErrorStruct{}
}

// GET ACTORS BY BIRTHDATE
func (r *ActorRepository) GetActorsByBirthdate(ctx context.Context, birthdate string) ([]models.Actor, errs.ErrorStruct) {

	actors := []models.Actor{}
	query := `SELECT id, name, birth_date FROM actors WHERE birth_date = ?`

	rows, err := r.db.QueryContext(ctx, query, birthdate)
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.ServerError,
			errors.New("something went wrong getting actor by birthdate"),
		)
		return []models.Actor{}, errStruct
	}
	defer rows.Close()

	for rows.Next() {
		var actor models.Actor
		err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			log.Println(err)
			errStruct := errs.NewErrorStruct(
				errs.ServerError,
				errors.New("something went wrong getting actor by birthday"),
			)
			return []models.Actor{}, errStruct
		}
		actors = append(actors, actor)
	}

	return actors, errs.ErrorStruct{}

}
