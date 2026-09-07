package service

import (
	"context"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"gitea.kood.tech/timdanielfiander/movies-api.git/repository"
	"strconv"
	"time"
	"log"
	"fmt"
	"errors"
)

type MovieService struct {
	repo *repository.MovieRepository
}

func NewMovieService(repo *repository.MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func (s *MovieService) GetMovies(filters models.MovieFilters) ([]models.Movie, errs.ErrorStruct) {

	if filters.ReleaseYear != nil {
		if _, err := strconv.Atoi(*filters.ReleaseYear); err != nil {
			log.Println(err)
			errStruct := errs.NewErrorStruct(
				errs.BadRequest,
				fmt.Errorf("invalid releaseyear: %v", *filters.ReleaseYear),
			)
			return []models.Movie{}, errStruct 
		}
	}

	movies, errStruct := s.repo.GetMovies(filters)
	if errStruct.ErrType != nil {
		return nil, errStruct 
	}
	return movies, errs.ErrorStruct{} 
}

func (s *MovieService) GetMovieByID(id int) (models.MovieReq, errs.ErrorStruct) {

	movie, errStruct := s.repo.GetMovieByID(id)
	if errStruct.ErrType != nil {
		return models.MovieReq{}, errStruct
	}
	return movie, errs.ErrorStruct{} 
}

func (s *MovieService) PostMovie(ctx context.Context, req models.Movie) (models.Movie, errs.ErrorStruct) {

	if req.Title == "" || req.Description == "" || req.ReleaseDate == "" || req.Duration == "" {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("cannot create movie with empty fields"),
		)
		return models.Movie{}, errStruct 
	}

	releaseDate := req.ReleaseDate
	_, err := time.Parse(timeLayout, releaseDate)
	if err != nil {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid release date: %v", req.ReleaseDate),
		)
		return models.Movie{}, errStruct 
	}

	movie, errStruct := s.repo.PostMovie(ctx, req)
	if errStruct.ErrType != nil {
		return models.Movie{}, errStruct
	}

	return movie, errs.ErrorStruct{}
}

func (s *MovieService) PatchMovie(ctx context.Context, req models.PatchMovieReq, id int) (models.Movie, errs.ErrorStruct) {

	if req.Title == nil && req.Description == nil && req.ReleaseDate == nil && req.Duration == nil {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("cannot update movie with no values"),
		)
		return models.Movie{}, errStruct 
	}

	if req.ReleaseDate != nil {
		_, err := time.Parse(timeLayout, *req.ReleaseDate)
		if err != nil {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid releasedate: %v", *req.ReleaseDate),
		)
		return models.Movie{}, errStruct 
		}
	}

	movie, errStruct := s.repo.PatchMovie(ctx, req, id)
	if errStruct.ErrType != nil {
		return models.Movie{}, errStruct 
	}

	return movie, errs.ErrorStruct{} 
}


func (s *MovieService) DeleteMovie(ctx context.Context, id int, force bool) errs.ErrorStruct {

	if force {
	errStruct := s.repo.DeleteForceMovie(ctx, id) 
	if errStruct.ErrType != nil {
	return errStruct 
	}
	return errs.ErrorStruct{} 
	} 

	if errStruct := s.repo.DeleteMovie(ctx, id); errStruct.ErrType != nil {
		return errStruct 
	}
	return errs.ErrorStruct{} 
}

func (s *MovieService) MovieActors(ctx context.Context, id int) ([]models.Actor, errs.ErrorStruct) {

	if id <= 0 {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid id: %v", id),
		)
		return []models.Actor{}, errStruct 
	}

	actors, errStruct := s.repo.MovieActors(ctx, id)
	if errStruct.ErrType != nil {
		return []models.Actor{}, errStruct 
	}
	if len(actors) == 0 {
		errStruct := errs.NewErrorStruct(
			errs.NotFound,
			errors.New("no actors found"),
		)	
		return actors, errStruct 
	}
	return actors, errs.ErrorStruct{} 
}
