package models

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"release_date"`
	Genres      []int   `json:"genres"`
}

type PatchMovieReq struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	ReleaseDate *string	 `json:"release_date"`
	Genres      []Genre  `json:"genres"`
}

type Genre struct {
	ID    int    `json:"id"`
	Genre string `json:"genre_name"`
}

type CreateGenreReq struct {
	Genre string `json:"genre_name"`
}

type Actor struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

type PatchActorReq struct {
	Name      *string `json:"name"`
	BirthDate *string `json:"birth_date"`
}

type MovieFilters struct {
	ActorID	   *string
	GenreID    *string
}
