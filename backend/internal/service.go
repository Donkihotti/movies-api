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
