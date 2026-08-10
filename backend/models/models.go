package models 

type Movie struct {
	ID 			int  	`json:"id"`
	Title 		string  `json:"title"`
	Description string 	`json:"description"`
	ReleaseDate int		`json:"release_date"`
}

type CreateMovieReq struct {
	Title 		string `json:"title"` 
	Description string `json:"description"`
	ReleaseDate int `json:"release_date"`
}

type Actor struct {
	ID 			int 	`json:"id"`
	Name		string	`json:"name"`
	BirthDate	string	`json:"birth_date"`
}

type CreateActorReq struct {
	Name		string	`json:"name"`
	BirthDate	string	`json:"birth_date"`
}
