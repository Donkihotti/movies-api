package internal

import (
	"context"
	"database/sql"

	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
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
func (r *ActorRepository) GetActors() ([]models.Actor, error) {
	var actors []models.Actor
	rows, err := r.db.Query(
		`
	SELECT id, name, birth_date
	FROM actors
	`,
	)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		actors = append(actors, actor)
	}

	return actors, nil
}

// GET AN INDIVIDUAL ACTOR
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
		return models.Actor{}, err
	}

	return actor, nil
}

//POST AN INDIVIDUAL ACTOR

func (r *ActorRepository) PostActor(ctx context.Context, actor models.Actor) (models.Actor, error) {

	// res is an sql.Result, an interface that has a method such as the
	// res.LastInsertId() method. 
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO actors (name, birth_date) VALUES (?, ?)`,
		actor.Name,
		actor.BirthDate,
	)

	if err != nil {
		return models.Actor{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Actor{}, err
	}

	actor.ID = int(id)
	return actor, nil

}
