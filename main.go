package main

import "net/http"

func main() {
	// create a new ServeMux (router)
	mux := http.NewServeMux()

	// serve files from the current directory at the root path
	mux.Handle("/", http.FileServer(http.Dir(".")))

	// create a Server struct with handler and addr
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
