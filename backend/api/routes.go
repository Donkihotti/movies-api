package api

import "net/http"

func NewRouter() *http.ServeMux {
	
	mux := http.NewServeMux()

	mux.HandleFunc("/", HomeHandler)

	return mux
}
