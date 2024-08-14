package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type apiConfig struct {
	fileserverHits int
}

// middleware logic to add more code to handler
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(h)

}
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(
			`<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body>
		</html>`, cfg.fileserverHits)))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// resets the hits to 0
func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

type parameters struct {
	Body string `json:"body"`
}
type errors1 struct {
	Error string `json:"error"`
}
type CleanedResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	err1 := errors1{
		Error: msg,
	}
	respError, err := json.Marshal(err1)
	if err != nil {
		http.Error(w, "JSON Encoding Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(respError)

}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "JSON Encoding Error", http.StatusInternalServerError)
		return
	}
	if code != http.StatusInternalServerError {
		w.WriteHeader(code)
	}
	w.Write(response)
}
func validationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		p := parameters{}
		err := decoder.Decode(&p)
		if len(p.Body) > 140 {
			respondWithError(w, 400, "Chirp is too long")
			return
		}
		if err != nil {
			respondWithError(w, 500, "Something went wrong")
			return
		}
		data := strings.Split(p.Body, " ")
		profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
		for i, d := range data {
			temp := strings.ToLower(d)
			for _, p := range profaneWords {
				if temp == p {
					data[i] = "****"
				}
			}
		}
		res := strings.Join(data, " ")
		cleanedBody := CleanedResponse{
			CleanedBody: res,
		}
		respondWithJSON(w, http.StatusOK, cleanedBody)
		return
	}
}

// for validate a tweet like a chrip
// func validationHandler(w http.ResponseWriter, r *http.Request) {
// 	type parameters struct {
// 		Body string `json:"body"`
// 	}
// 	type resp struct {
// 		Error string `json:"error"`
// 	}
// 	type allOk struct {
// 		Valid bool `json:"valid"`
// 	}
// 	if r.Method == "POST" {
// 		decoder := json.NewDecoder(r.Body)
// 		p := parameters{}
// 		err := decoder.Decode(&p)
// 		if len(p.Body) > 140 {
// 			resperr := resp{
// 				Error: "Chirp is too long",
// 			}
// 			respError, _ := json.Marshal(resperr)
// 			w.WriteHeader(400)
// 			w.Write(respError)
// 			return

// 		}
// 		if err != nil {
// 			resperr := resp{
// 				Error: "Something went wrong",
// 			}
// 			respError, _ := json.Marshal(resperr)
// 			w.WriteHeader(500)
// 			w.Write(respError)
// 			return
// 		}
// 		allGood := allOk{
// 			Valid: true,
// 		}
// 		ans, _ := json.Marshal(allGood)
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(ans)
// 		return

// 	}

// }
func main() {
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}
	// separate handler for healthz
	readiness := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

	//main requests from the users

	//mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/api/healthz", readiness)
	mux.HandleFunc("/admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
	mux.HandleFunc("/api/validate_chirp", validationHandler)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
