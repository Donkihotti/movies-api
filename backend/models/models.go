package models

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	Duration    string `json:"duration"`
	Genres      []int  `json:"genres"`
	Actors      []int  `json:"actors"`
}

type PatchMovieReq struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	ReleaseDate *string `json:"release_date"`
	Duration    *string `json:"duration"`
	Genres      []Genre `json:"genres"`
	Actors      []int   `json:"actors"`
}

type Genre struct {
	ID    int    `json:"id"`
	Genre string `json:"name"`
}

type CreateGenreReq struct {
	Genre string `json:"name"`
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
	ActorID     *string
	GenreID     *string
	ReleaseYear *string
	Duration    *string
}
