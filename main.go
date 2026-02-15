package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	// create a new ServeMux (router)
	mux := http.NewServeMux()

	// initialize apiConfig
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// serve files from the current directory at the root path
	// mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	// set healthz path
	mux.HandleFunc("GET /healthz", healthzHandlerFunction)

	// set metrics path
	mux.HandleFunc("GET /metrics", apiCfg.metricsHandlerFunction)

	// set reset path
	mux.HandleFunc("POST /reset", apiCfg.resetHandlerFunction)

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

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandlerFunction(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hits: %v", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetHandlerFunction(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}
