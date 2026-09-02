package service

import (
	"context"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"gitea.kood.tech/timdanielfiander/movies-api.git/repository"
	"time"
	"log"
	"errors"
	"fmt"
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

func (as *ActorService) GetActorsByBirthdate(ctx context.Context, birthdate string) ([]models.Actor, error) {

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

func (as *ActorService) DeleteActor(ctx context.Context, id int) error {

	err := as.repo.DeleteActor(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

//done?
func (as *ActorService) PatchActor(ctx context.Context, req models.PatchActorReq, id int) errs.ErrorStruct {

	if req.Name == nil && req.BirthDate == nil {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("no new values set updating actor"),
		)
		log.Println(errStruct.Error())	
		return errStruct
	}

	_, err := time.Parse(timeLayout, *req.BirthDate)	
	if err != nil {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid birthdate: %v", *req.BirthDate),
		)
		return errStruct
	}

	errStruct := as.repo.PatchActor(ctx, req, id)
	if errStruct.ErrType != nil {
		return errStruct
	}
	return errs.ErrorStruct{} 
}







