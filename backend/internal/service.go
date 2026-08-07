package internal

import (
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
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

func(s *Service) GetMovieByID(id int) (models.Movie, error) {
	movie, err := s.repo.GetMovieByID(id)
	if err != nil {
	// return a struct because nil does not work on empty structs
	return models.Movie{}, err
	}
	return movie, nil
}
