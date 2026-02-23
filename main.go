package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gutek00714/chirpy---Boot.dev/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	// load .env into environment variables
	godotenv.Load()

	// get DB_URL from the environment
	dbURL := os.Getenv("DB_URL")

	// open a connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// create a new *database.Queries and store it in apiConfig struct
	dbQueries := database.New(db)

	// create a new ServeMux (router)
	mux := http.NewServeMux()

	// initialize apiConfig
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
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
	// mux.HandleFunc("POST /api/validate_chirp", validateHelperFunction)

	// set users path
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	// set chirps path
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)

	// create a Server struct with handler and addr
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Println("Serving on port: 8080")
	log.Fatal((server.ListenAndServe()))

	// server.ListenAndServe()
}
