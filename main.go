package main

import (
	"encoding/json"
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
	mux.HandleFunc("GET /api/healthz", healthzHandlerFunction)

	// set metrics path
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandlerFunction)

	// set reset path
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandlerFunction)

	// set validate_chirp paht
	mux.HandleFunc("POST /api/validate_chirp", validateHelperFunction)

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

func validateHelperFunction(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	// check if body is 140 or less characters
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		respondWithJSON(w, 200, validResponse{Valid: true})
	}
}

// helper functions
type validResponse struct {
	Valid bool `json:"valid"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	resp := errorResponse{Error: msg}
	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
