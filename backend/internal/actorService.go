package internal

import (
	"context"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
)

type ActorService struct {
	repo *ActorRepository
}

func NewActorService(repo *ActorRepository) *ActorService {
	return &ActorService{
		repo: repo,
	}
}

func (as *ActorService) GetActors() ([]models.Actor, error) {
	actors, err := as.repo.GetActors()
	if err != nil {
		return nil, err
	}
	return actors, nil
}

func (as *ActorService) GetActorByID(id int) (models.Actor, error) {
	actor, err := as.repo.GetActorByID(id)
	if err != nil {
		return models.Actor{}, err
	}
	return actor, nil
}

func (as *ActorService) PostActor(ctx context.Context, req models.CreateActorReq) (models.Actor, error) {

	actor := models.Actor{
		Name:      req.Name,
		BirthDate: req.BirthDate,
	}

	res, err := as.repo.PostActor(ctx, actor)
	if err != nil {
		return models.Actor{}, err
	}
	return res, nil
}

func (as *ActorService) DeleteActor(ctx context.Context, id int) error {

	err := as.repo.DeleteActor(ctx, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (as *ActorService) PatchActor(ctx context.Context, actor models.Actor, id int) error {
	err := as.repo.PatchActor(ctx, actor, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
