package internal

import (
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"context"
)

type ActorService struct {
	repo *ActorRepository
}

func NewActorService(repo *ActorRepository) *ActorService {
	return &ActorService{
		repo: repo,
	}
} 

//functions...


func(as *ActorService) GetActors() ([]models.Actor, error) {
	actors, err := as.repo.GetActors()
	if err != nil {
	return nil, err
	}
	return actors, nil
}

func(as *ActorService) GetActorByID(id int) (models.Actor, error) {
	actor, err := as.repo.GetActorByID(id)
	if err != nil {
	return models.Actor{}, err
	}
	return actor, nil
}

func(as *ActorService) PostActor(ctx context.Context, req models.CreateActorReq) (models.Actor, error) {

	actor := models.Actor{
		Name: req.Name,
		BirthDate: req.BirthDate,		
	}

	res, err := as.repo.PostActor(ctx, actor)
	if err != nil {
	return models.Actor{}, err
	}
	return res, nil
}
