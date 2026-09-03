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
func (r *ActorRepository) GetActors(name string) ([]models.Actor, error) {

	actors := []models.Actor{}
	var rows *sql.Rows
	var err error

	if name == "" {
		query := "SELECT id, name, birth_date FROM actors"
		rows, err = r.db.Query(query)
		if err != nil {
			return actors, errs.ServerError
		}
	} else {
		name = fmt.Sprintf("%%%s%%", name)
		query := "SELECT id, name, birth_date FROM actors WHERE name LIKE ?"
		rows, err = r.db.Query(query, name)
		if err != nil {
			return actors, errs.ServerError
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
			return nil, errs.ServerError
		}
		actors = append(actors, actor)
	}

	return actors, nil
}

// get actor by id
func (r *ActorRepository) GetActorByID(id int) (models.Actor, error) {
	var actor models.Actor

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
			return models.Actor{}, errs.NotFound
		}
		return models.Actor{}, errs.ServerError
	}

	return actor, nil
}

// POST AN ACTOR
func (r *ActorRepository) PostActor(ctx context.Context, req models.Actor) (models.Actor, error) {

	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO actors (name, birth_date) VALUES (?, ?)`,
		req.Name,
		req.BirthDate,
	)

	if err != nil {
		return models.Actor{}, errs.ServerError
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Actor{}, errs.ServerError
	}

	req.ID = int(id)
	return req, nil

}

func (r *ActorRepository) DeleteForceActor(ctx context.Context, id int) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
	return err
	}

	defer tx.Rollback()
	query := `DELETE FROM movie_actors WHERE actor_id = ?`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
	return err
	}

	query = `DELETE FROM actors WHERE id = ?`

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

func (r *ActorRepository) DeleteActor(ctx context.Context, id int) error {

	query := `DELETE FROM actors WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, id)
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

// PATCH AN ACTOR
func (r *ActorRepository) PatchActor(ctx context.Context, req models.PatchActorReq, id int) error {

	query := `UPDATE actors SET name = COALESCE(?, name), birth_date = COALESCE(?, birth_date) WHERE id = ?`
	row, err := r.db.ExecContext(ctx, query, req.Name, req.BirthDate, id)
	if err != nil {
		return errs.ServerError
	}
	rows, err := row.RowsAffected()
	if err != nil {
		return errs.ServerError
	}

	if rows == 0 {
		return errs.NotFound
	}
	return nil
}

// GET ACTORS BY NAME
func (r *ActorRepository) GetActorsByName(ctx context.Context, name string) ([]models.Actor, error) {

	actors := []models.Actor{}

	name = fmt.Sprintf("%%%s%%", name)

	query := `SELECT id, name, birth_date FROM actors WHERE name LIKE ?`

	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return actors, errs.ServerError
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
			return []models.Actor{}, errs.ServerError
		}
		actors = append(actors, actor)
	}

	return actors, nil

}

// GET ACTORS BY BIRTHDATE
func (r *ActorRepository) GetActorsByBirthdate(ctx context.Context, birthdate string) ([]models.Actor, error) {

	actors := []models.Actor{}

	query := `SELECT id, name, birth_date FROM actors WHERE birth_date = ?`

	rows, err := r.db.QueryContext(ctx, query, birthdate)
	if err != nil {
		return []models.Actor{}, errs.ServerError
	}

	for rows.Next() {
		var actor models.Actor

		err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			log.Println(err)
			return []models.Actor{}, errs.ServerError
		}
		actors = append(actors, actor)
	}

	return actors, nil

}
