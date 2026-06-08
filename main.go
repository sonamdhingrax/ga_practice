// Command ga_practice is a tiny HTTP service used to practice GitHub Actions.
//
// It exposes two endpoints:
//
//	GET /         -> JSON build metadata (version, commit, build date)
//	GET /healthz  -> liveness probe, always returns 200 "ok"
//
// The build metadata is injected at compile time via -ldflags "-X main.var=...",
// which is how our CI pipeline stamps the binary with the git SHA and tag.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// These values are overridden at build time with:
//
//	go build -ldflags "-X main.version=v1.2.3 -X main.commit=abc123 -X main.date=2026-06-08"
//
// When built without ldflags (e.g. `go run .`) they fall back to "dev".
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// buildInfo is the JSON payload returned by the root handler.
type buildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Service string `json:"service"`
}

// rootHandler reports the build metadata so we can confirm which image is live.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(buildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
		Service: "ga_practice",
	})
}

// healthHandler is a trivial liveness probe used by ECS/ALB health checks.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// newServer wires up the routes. Kept separate from main() so tests can exercise
// the same mux without binding a real port.
func newServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthHandler)
	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           newServer(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("ga_practice %s (commit %s, built %s) listening on :%s", version, commit, date, port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Errorf("server failed: %w", err))
	}
}
