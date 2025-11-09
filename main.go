package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/7minutech/chripy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
	Email     string    `json:"email"`
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

	if apiCfg.platform != "dev" {
		msg := "must be in dev platform"
		respondWithError(w, http.StatusForbidden, msg, fmt.Errorf("error: trying reset while not platform is not dev"))
		return
	}

	if err := apiCfg.dbQueries.DeleteUsers(r.Context()); err != nil {
		msg := "could not delete users"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}
}

func (apiCfg *apiConfig) handerUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	var params = parameters{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		msg := "could not decode request body"
		respondWithError(w, http.StatusBadRequest, msg, err)
		return
	}

	user, err := apiCfg.dbQueries.CreateUser(r.Context(), params.Email)

	if err != nil {
		msg := "could not create user"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	resp := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func main() {
	godotenv.Load(".env")

	platform := os.Getenv("PLATFORM")
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
		platform:  platform,
	}

	mux := http.NewServeMux()

	handlerFile := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handlerFile))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetric)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handerUser)
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
