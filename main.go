package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/7minutech/chripy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func (apiCfg *apiConfig) handlerMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(
		"<html><body>"+
			"<h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p>"+
			"</body></html>", apiCfg.fileserverHits.Load())
	w.Write([]byte(body))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cfg.fileserverHits.Add(1)

		next.ServeHTTP(w, r)
	})
}

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	apiCfg.fileserverHits.Store(0)

	w.WriteHeader(http.StatusOK)
}

func main() {
	godotenv.Load(".env")

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("failed to open data base")
	}

	queries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	var apiCfg = apiConfig{
		dbQueries: queries,
	}

	mux := http.NewServeMux()

	handlerFile := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handlerFile))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetric)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
