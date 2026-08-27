package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"errors"
	"time"
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
    var actors []models.Actor
    var rows *sql.Rows
    var err error

    if name == "" {
        query := "SELECT id, name, birth_date FROM actors"
        rows, err = r.db.Query(query)
        if err != nil {
            log.Println(err)
            return actors, err
        }
    } else {
        name = fmt.Sprintf("%%%s%%", name)
        query := "SELECT id, name, birth_date FROM actors WHERE name LIKE ?" 
        rows, err = r.db.Query(query, name)
        if err != nil {
            log.Println(err)
            return actors, err
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

//POST AN ACTOR
func (r *ActorRepository) PostActor(ctx context.Context, req models.Actor) (models.Actor, error) {

	// res is an sql.Result, an interface that has a method such as the
	// res.LastInsertId() method.

    if req.Name == "" || req.BirthDate == "" {
        message := "cannot post actor with empty struct fields"
        log.Println(message)
        return models.Actor{}, errors.New(message)
    }   

    timeLayout := "2006-01-02"
    birthdate := req.BirthDate
    _, err := time.Parse(timeLayout, birthdate)
    if err != nil {
        log.Println("invalid birthdate: ", birthdate)
        return models.Actor{}, err 
    }   	

	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO actors (name, birth_date) VALUES (?, ?)`,
		req.Name,
		req.BirthDate,
	)

	if err != nil {
		return models.Actor{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Actor{}, err
	}

	req.ID = int(id)
	return req, nil

}

//DELETE AN ACTOR
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

//PATCH AN ACTOR
func (r *ActorRepository) PatchActor(ctx context.Context, req models.PatchActorReq, id int) error {

    query := `UPDATE actors SET name = COALESCE(?, name), birth_date = COALESCE(?, birth_date) WHERE id = ?`
    row, err := r.db.ExecContext(ctx, query, req.Name, req.BirthDate, id)
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

//GET ACTORS BY NAME
func (r *ActorRepository) GetActorsByName(ctx context.Context, name string) ([]models.Actor, error) {

    var actors []models.Actor

    name = fmt.Sprintf("%%%s%%", name)

    query := `SELECT id, name, birth_date FROM actors WHERE name LIKE ?`

    rows, err := r.db.QueryContext(ctx, query, name)
    if err != nil {
        log.Println(err)
        return actors, err
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
            return []models.Actor{}, err
        }
        actors = append(actors, actor)
    }

    return actors, nil

}

//GET ACTORS BY BIRTHDATE
func (r *ActorRepository) GetActorsByBirthdate(ctx context.Context, birthdate string) ([]models.Actor, error) {

    var actors []models.Actor

    query := `SELECT id, name, birth_date FROM actors WHERE birth_date = ?`

    rows, err := r.db.QueryContext(ctx, query, birthdate)
    if err != nil {
        log.Println(err)
        return []models.Actor{}, err
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
                return []models.Actor{}, err
        }
        actors = append(actors, actor)
    }

    return actors, nil

}















