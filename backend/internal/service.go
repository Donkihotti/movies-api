package internal

import (
	"context"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
)

type MovieService struct {
	repo *MovieRepository
}

func NewMovieService(repo *MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func (s *MovieService) GetMovies() ([]models.Movie, error) {
	movies, err := s.repo.GetMovies()
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (s *MovieService) PostMovie(ctx context.Context, req models.CreateMovieReq) (models.Movie, error) {

	movie := models.Movie{
		Title:       req.Title,
		Description: req.Description,
		ReleaseDate: req.ReleaseDate,
	}

	res, err := s.repo.PostMovie(ctx, movie)
	if err != nil {
		return models.Movie{}, err
	}
	return res, nil
}

func (s *MovieService) GetMovieByID(id int) (models.Movie, error) {
	movie, err := s.repo.GetMovieByID(id)
	if err != nil {
		return models.Movie{}, err
	}
	return movie, nil
}

func (s *MovieService) PatchMovie(ctx context.Context, req models.Movie, id int) error {

	err := s.repo.PatchMovie(ctx, req, id)
	if err != nil {
		return err
	}
	return nil
}
