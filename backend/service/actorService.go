package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"gitea.kood.tech/timdanielfiander/movies-api.git/repository"
)

type ActorService struct {
	repo *repository.ActorRepository
}

func NewActorService(repo *repository.ActorRepository) *ActorService {
	return &ActorService{
		repo: repo,
	}
}

const timeLayout = "2006-01-02"

func (as *ActorService) GetActors(name string) ([]models.Actor, errs.ErrorStruct) {
	actors, errStruct := as.repo.GetActors(name)
	if errStruct.ErrType != nil {
		return nil, errStruct
	}
	return actors, errs.ErrorStruct{}
}

func (as *ActorService) GetActorByID(id int) (models.Actor, errs.ErrorStruct) {
	actor, errStruct := as.repo.GetActorByID(id)
	if errStruct.ErrType != nil {
		return models.Actor{}, errStruct
	}
	return actor, errs.ErrorStruct{}
}

func (as *ActorService) GetActorsByName(ctx context.Context, name string) ([]models.Actor, errs.ErrorStruct) {

	actors, errStruct := as.repo.GetActorsByName(ctx, name)
	if errStruct.ErrType != nil {
		return []models.Actor{}, errStruct
	}
	return actors, errs.ErrorStruct{}
}

func (as *ActorService) GetActorsByBirthdate(ctx context.Context, birthdate string) ([]models.Actor, errs.ErrorStruct) {

	_, err := time.Parse(timeLayout, birthdate)
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid birthdate: %v", birthdate),
		)
		return []models.Actor{}, errStruct
	}

	actors, errStruct := as.repo.GetActorsByBirthdate(ctx, birthdate)
	if errStruct.ErrType != nil {
		return []models.Actor{}, errStruct
	}
	return actors, errs.ErrorStruct{}
}

func (as *ActorService) PostActor(ctx context.Context, req models.Actor) (models.Actor, errs.ErrorStruct) {

	if req.Name == "" || req.BirthDate == "" {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("empty fields are not allowed when creating an actor"),
		)
		log.Println(errStruct.Error())
		return models.Actor{}, errStruct
	}

	_, err := time.Parse(timeLayout, req.BirthDate)
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid time: %v", req.BirthDate),
		)
		return models.Actor{}, errStruct
	}

	res, errStruct := as.repo.PostActor(ctx, req)
	if errStruct.ErrType != nil {
		return models.Actor{}, errStruct
	}
	return res, errs.ErrorStruct{}
}

func (as *ActorService) DeleteActor(ctx context.Context, id int, force bool) errs.ErrorStruct {

	if force {
	errStruct := as.repo.DeleteForceActor(ctx, id)
	if errStruct.ErrType != nil {
		return errStruct
	}
	return errs.ErrorStruct{}
	}

	errStruct := as.repo.DeleteActor(ctx, id)
	if errStruct.ErrType != nil {
		return errStruct
	}
	return errs.ErrorStruct{}
}

func (as *ActorService) PatchActor(ctx context.Context, req models.PatchActorReq, id int) (models.Actor, errs.ErrorStruct) {

	if req.Name == nil && req.BirthDate == nil {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("no new values set updating actor"),
		)
		log.Println(errStruct.Error())
		return models.Actor{}, errStruct
	}

	if req.BirthDate != nil {
		_, err := time.Parse(timeLayout, *req.BirthDate)
		if err != nil {
			errStruct := errs.NewErrorStruct(
				errs.BadRequest,
				fmt.Errorf("invalid birthdate: %v", *req.BirthDate),
			)
			return models.Actor{}, errStruct
		}
	}

	patchedMovie, errStruct := as.repo.PatchActor(ctx, req, id)
	if errStruct.ErrType != nil {
		return models.Actor{}, errStruct
	}
	return patchedMovie, errs.ErrorStruct{}
}
