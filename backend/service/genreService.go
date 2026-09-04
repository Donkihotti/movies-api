package service

import (
	"context"
	"errors"
	"log"

	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"gitea.kood.tech/timdanielfiander/movies-api.git/repository"
)

type GenreService struct {
	repo *repository.GenreRepository
}

func NewGenreService(repo *repository.GenreRepository) *GenreService {
	return &GenreService{
		repo: repo,
	}
}

func (s *GenreService) GetGenres(ctx context.Context) ([]models.Genre, errs.ErrorStruct) {

	genres, errStruct := s.repo.GetGenres(ctx)
	if errStruct.ErrType != nil {
		return nil, errStruct
	}

	return genres, errs.ErrorStruct{}
}

func (s *GenreService) GetGenre(ctx context.Context, id int) (models.Genre, errs.ErrorStruct) {

	genre, errStruct := s.repo.GetGenre(ctx, id)
	if errStruct.ErrType != nil {
		return models.Genre{}, errStruct
	}
	return genre, errs.ErrorStruct{}
}

func (s *GenreService) PostGenre(ctx context.Context, req models.Genre) (models.Genre, errs.ErrorStruct) {

	if req.Genre == "" {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("empty fields not allowed when posting genre"),
		)
		return models.Genre{}, errStruct
	}

	res, errStruct := s.repo.PostGenre(ctx, req)
	if errStruct.ErrType != nil {
		return models.Genre{}, errStruct
	}

	return res, errs.ErrorStruct{}
}

func (s *GenreService) DeleteGenre(ctx context.Context, id int, force bool) errs.ErrorStruct {

	if force {
	errStruct := s.repo.DeleteForceGenre(ctx, id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		return errStruct
	}
	return errStruct
	}

	errStruct := s.repo.DeleteGenre(ctx, id)
	if errStruct.ErrType != nil {
		return errStruct 
	}
	return errStruct
}

func (s *GenreService) PatchGenre(ctx context.Context, newGenre models.Genre, id int) (models.Genre, errs.ErrorStruct) {

	if newGenre.Genre == "" {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("cannot update genre with empty field"),
		)
		return models.Genre{}, errStruct 
	}

	patchedGenre, errStruct := s.repo.PatchGenre(ctx, newGenre, id)
	if errStruct.ErrType != nil {
		log.Println(patchedGenre)
		return models.Genre{}, errStruct 
	}
	return patchedGenre, errs.ErrorStruct{} 
}
