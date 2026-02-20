package main

import (
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

	// set validate_chirp path
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
