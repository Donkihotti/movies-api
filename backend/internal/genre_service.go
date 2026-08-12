package internal 

import (
	"context"
	 
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

func(s *GenreService) GetGenres(ctx context.Context) ([]models.Genre, error) {
	
	genres, err := s.repo.GetGenres(ctx)
	if err != nil {
	return nil, err
	}
	
	return genres, nil
}

func(s *GenreService) PostGenre(ctx context.Context, req models.Genre) (models.Genre, error) {

	genre := models.Genre{
	Genre: req.Genre,
	}

	res, err := s.repo.PostGenre(ctx, genre)	
	if err != nil {
	return models.Genre{}, err
	}
	
	return res, nil
}

func(s *GenreService) DeleteGenre(ctx context.Context, id int) error {

	err := s.repo.DeleteGenre(ctx, id)
	if err != nil {
	return err
	}
	return nil
}


func (s *GenreService) PutGenre(ctx context.Context, newGenre models.Genre, id int) error {

	err := s.repo.PutGenre(ctx, newGenre, id)
	if err != nil {
	return err 
	}
	return nil

}
