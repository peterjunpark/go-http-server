package main

import (
	"log"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type: text/plain", "charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	const port = "8008"
	mux := http.NewServeMux()
	corsMux := corsMiddleware(mux)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.HandleFunc("/healthz", handleReadiness)

	log.Printf("🚀 Server listening on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
