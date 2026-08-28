package internal

import (
	"context"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"time"
	"errors"
	"log"
)

type MovieService struct {
	repo *MovieRepository
}


func NewMovieService(repo *MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

//GET ALL MOVIES
func (s *MovieService) GetMovies(filters models.MovieFilters) ([]models.Movie, error) {

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

//POST MOVIE
func (s *MovieService) PostMovie(ctx context.Context, req models.Movie) (models.Movie, error) {

    if req.Title == "" || req.Description == "" || req.ReleaseDate == "" {
        message := "cannot post movie with empty struct fields"
        log.Println(message)
        return models.Movie{}, errors.New(message)
    }

    timeLayout := "2006-01-02"
    releaseDate := req.ReleaseDate
    _, err := time.Parse(timeLayout, releaseDate)
    if err != nil {
        log.Println("invalid release date: ", releaseDate)
        return models.Movie{}, err
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
        message := "cannot accept empty struct fields in patch movie method"
        log.Println(message)
        return errors.New(message)
    }

    timeLayout := "2006-01-02"
    releaseDate := *req.ReleaseDate
    _, err := time.Parse(timeLayout, releaseDate)
    if err != nil {
        log.Println("invalid release date: ", releaseDate)
        return err
    }

    err = s.repo.PatchMovie(ctx, req, id)
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


