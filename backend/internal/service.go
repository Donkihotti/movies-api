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
	
	movie, err := s.repo.PostMovie(ctx, req)
	if err != nil {
	return nil, err
	}
	return movie, nil
}
