package internal

import (
	"context"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"time"
	"strconv"
)

type MovieService struct {
	repo *MovieRepository
}


func NewMovieService(repo *MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

var timeLayout = "2006-01-02"

//GET ALL MOVIES
func (s *MovieService) GetMovies(filters models.MovieFilters) ([]models.Movie, error) {

	if filters.ReleaseYear != nil {
		if _, err := strconv.Atoi(*filters.ReleaseYear); err != nil {
			return []models.Movie{}, errs.BadRequest
		}
	}
	
	movies, err := s.repo.GetMovies(filters)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

//GET MOVIE BY ID
func (s *MovieService) GetMovieByID(id int) (models.Movie, error) {
	movie, err := s.repo.GetMovieByID(id)
	if err != nil {
		return models.Movie{}, err
	}
	return movie, nil
}


func (s *MovieService) PostMovie(ctx context.Context, req models.Movie) (models.Movie, error) {

    if req.Title == "" || req.Description == "" || req.ReleaseDate == "" {
        return models.Movie{}, errs.BadRequest 
    }

    releaseDate := req.ReleaseDate
    _, err := time.Parse(timeLayout, releaseDate)
    if err != nil {
        return models.Movie{}, errs.BadRequest
    }

    movie, err := s.repo.PostMovie(ctx, req)
    if err != nil {
        return models.Movie{}, err
    }

    return movie, nil
}


//PATCH MOVIE
func (s *MovieService) PatchMovie(ctx context.Context, req models.PatchMovieReq, id int) error {

    if req.Title == nil && req.Description == nil && req.ReleaseDate == nil {
        return errs.BadRequest 
    }

	if req.ReleaseDate != nil {
		releaseDate := *req.ReleaseDate
		_, err := time.Parse(timeLayout, releaseDate)
		if err != nil {
			return errs.BadRequest
		}
	}

    err := s.repo.PatchMovie(ctx, req, id)
    if err != nil {
        return err
    }
    return nil
}

//DELETE MOVIE
func (s *MovieService) DeleteMovie(ctx context.Context, id int) error {
    if err := s.repo.DeleteMovie(ctx, id); err != nil {
        return err
    }
    return nil
}


