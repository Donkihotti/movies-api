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
