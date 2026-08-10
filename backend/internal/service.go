package internal

import (
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"context"
)

type Service struct {
	repo *Repository	
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}


func(s *Service) GetMovies() ([]models.Movie, error) {
	movies, err := s.repo.GetMovies()	
	if err != nil {
	return nil, err
	}
	return movies, nil
}

func(s *Service) PostMovie(ctx context.Context, req models.CreateMovieReq) (models.Movie, error) {

	movie := models.Movie{
	Title: req.Title, 
	Description: req.Description, 
	ReleaseDate: req.ReleaseDate,
	}
	
	res, err := s.repo.PostMovie(ctx, movie)
	if err != nil {
	return models.Movie{}, err
	}
	return res, nil
}

func(s *Service) GetMovieByID(id int) (models.Movie, error) {
	movie, err := s.repo.GetMovieByID(id)
	if err != nil {
	return models.Movie{}, err
	}
	return movie, nil
}
