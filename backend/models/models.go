package models 

type Movie struct {
	ID 			int  	`json:"id"`
	Title 		string  `json:"title"`
	Description string 	`json:"description"`
	ReleaseDate int		`json:"release_date"`
	Genres 		[]Genre `json:"genres"`
}

type CreateMovieReq struct {
	Title 		string 	`json:"title"` 
	Description string 	`json:"description"`
	ReleaseDate int		`json:"release_date"`
	Genres 		[]Genre `json:"genres"`
}


type Genre struct {
	ID 			int  	`json:"id"`
	Genre 		string  `json:"genre_name"`
}


type CreateGenreReq struct {
	Genre 		string  `json:"genre_name"`
}

type PutGenreReq struct {
	ID 			int  	`json:"id"`
	Genre 		string  `json:"genre_name"`
}
