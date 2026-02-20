package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

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
	// fmt.Fprintf(w, "Hits: %v", cfg.fileserverHits.Load())
	ret := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(ret))
}

func (cfg *apiConfig) resetHandlerFunction(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}
