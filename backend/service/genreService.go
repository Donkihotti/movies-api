package service

import (
	"context"
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

func (s *GenreService) GetGenres(ctx context.Context) ([]models.Genre, error) {

	genres, err := s.repo.GetGenres(ctx)
	if err != nil {
		return nil, err
	}

	return genres, nil
}

func (s *GenreService) GetGenre(ctx context.Context, id int) ([]models.Movie, error) {

	genre, err := s.repo.GetGenre(ctx, id)
	if err != nil {
		return nil, err
	}

	return genre, nil
}

func (s *GenreService) PostGenre(ctx context.Context, req models.Genre) (models.Genre, error) {

	genre := models.Genre{
		Genre: req.Genre,
	}

	if req.Genre == "" {
		return models.Genre{}, errs.BadRequest
	}

	res, err := s.repo.PostGenre(ctx, genre)
	if err != nil {
		return models.Genre{}, err
	}

	return res, nil
}

func (s *GenreService) DeleteGenre(ctx context.Context, id int) error {

	err := s.repo.DeleteGenre(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *GenreService) PatchGenre(ctx context.Context, newGenre models.Genre, id int) error {

	if newGenre.Genre == "" {
		return errs.BadRequest
	}

	err := s.repo.PatchGenre(ctx, newGenre, id)
	if err != nil {
		return err
	}
	return nil

}
