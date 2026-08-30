package internal

import (
	"context"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"time"
)

type ActorService struct {
	repo *ActorRepository
}

func NewActorService(repo *ActorRepository) *ActorService {
	return &ActorService{
		repo: repo,
	}
}

//GET ALL ACTORS
func (as *ActorService) GetActors(name string) ([]models.Actor, error) {
	actors, err := as.repo.GetActors(name)
	if err != nil {
		return nil, err
	}
	return actors, nil
}

//GET ACTORS BY ID
func (as *ActorService) GetActorByID(id int) (models.Actor, error) {
	actor, err := as.repo.GetActorByID(id)
	if err != nil {
		return models.Actor{}, err
	}
	return actor, nil
}

//GET ACTOR BY NAME
func (as *ActorService) GetActorsByName(ctx context.Context, name string) ([]models.Actor, error) {
    actors, err := as.repo.GetActorsByName(ctx, name)
    if err != nil {
        return []models.Actor{}, err
    }
    return actors, nil
}

//GET ACTOR BY BIRTHDATE
func (as *ActorService) GetActorsByBirthdate(ctx context.Context, birthdate string) ([]models.Actor, error) {
    //validate time
    timeLayout := "2006-01-02"
    _, err := time.Parse(timeLayout, birthdate)
    if err != nil {
        return []models.Actor{}, errs.BadRequest
    }

    actors, err := as.repo.GetActorsByBirthdate(ctx, birthdate)
    if err != nil {
        return []models.Actor{}, err
    }
    return actors, nil
}

//POST ACTOR
func (as *ActorService) PostActor(ctx context.Context, req models.Actor) (models.Actor, error) {

    //Validate the input so that there are no empty fields.

    if req.Name == "" || req.BirthDate == "" {
        return models.Actor{}, errs.BadRequest 
    }

    //validate birthdate
    timeLayout := "2006-01-02"
    _, err := time.Parse(timeLayout, req.BirthDate)
    if err != nil {
        return models.Actor{}, errs.BadRequest
    }

    res, err := as.repo.PostActor(ctx, req)
    if err != nil {
        return models.Actor{}, err
    }
    return res, nil
}


//DELETE ACTOR
func (as *ActorService) DeleteActor(ctx context.Context, id int) error {

	err := as.repo.DeleteActor(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

//PATCH ACTOR
func (as *ActorService) PatchActor(ctx context.Context, req models.PatchActorReq, id int) error {

    if req.Name == nil && req.BirthDate == nil {
        return errs.BadRequest 
    }

    err := as.repo.PatchActor(ctx, req, id)
    if err != nil {
        return err
    }
    return nil
}

