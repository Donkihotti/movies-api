package internal

import (
	"context"
	"database/sql"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"fmt"
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

	//TODO we need also to create a birthDate validator here. The date that the user inputs in this program needs to be a valid date. 
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

//TODO: DELETE actors, PATCH actors, GET actors in a specific movie, GET all actors by its specific movie id. Retrieve actor by filtering their name. 
func (r *ActorRepository) DeleteActor(ctx context.Context, id int) error {

	query := `DELETE FROM actors WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
	return err	
	}
	
	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err	
	}
	
	if rows == 0 {
		return sql.ErrNoRows
	}	

	return nil
}

func (r *ActorRepository) PatchActor(ctx context.Context, actor models.Actor, id int) error {
	
	query := `UPDATE actors SET name = ?, birth_date = ? WHERE id = ?`
	row, err := r.db.ExecContext(ctx, query, actor.Name, actor.BirthDate, id)
	if err != nil {
		fmt.Println(err)
		return err 
	}
	rows, err := row.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
















