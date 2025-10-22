package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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
	const filepathRoot = "."
	const port = "8080"

	var apiCfg = apiConfig{}

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

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	const maxChripLength int = 140

	var params parameters

	type returnError struct {
		Error string `json:"error"`
	}

	type returnVal struct {
		Valid bool `json:"valid"`
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		resp, err := json.Marshal(returnError{Error: "Something went wrong"})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	if len(params.Body) > maxChripLength {
		resp, err := json.Marshal(returnError{Error: "Chirp is too long"})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	resp, err := json.Marshal(returnVal{Valid: true})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
