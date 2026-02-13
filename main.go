package main

import (
	"log"
	"net/http"
)

func main() {
	// create a new ServeMux (router)
	mux := http.NewServeMux()

	// serve files from the current directory at the root path
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	// set healthz path
	mux.HandleFunc("/healthz", healthzHandlerFunction)

	// create a Server struct with handler and addr
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Println("Serving on port: 8080")
	log.Fatal((server.ListenAndServe()))

	// server.ListenAndServe()
}

func healthzHandlerFunction(w http.ResponseWriter, r *http.Request) {
	// set content-type header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	// set the status code
	w.WriteHeader(http.StatusOK)

	// write the body text
	myBytes := []byte("OK")
	w.Write(myBytes)
}
