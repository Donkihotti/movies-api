package internal 

import (
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
)

type GenreService struct {
	repo *GenreRepository	
}

func NewGenreService(repo *GenreRepository) *GenreService {
	return &GenreService{
		repo: repo,
	}
}

func(s *GenreService) getGenres() ([]models.Genre, error) {
	
	genres, err := s.repo.GetGenres()
	if err != nil {
	return nil, err
	}
	
	return genres, nil
}
