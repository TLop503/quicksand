package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/tlop503/quicksand/docker_sdk"
)

type resp struct {
	OK        bool   `json:"ok"`
	IframeURL string `json:"iframeUrl,omitempty"`
	Error     string `json:"error,omitempty"`
}

var (
	currentName    string // name returned by docker_sdk.StartContainer
	currentHostURL string // e.g. http://localhost:5800
	imageFirefox   = "jlesage/firefox"
	imageTor       = "domistyle/tor-browser"
)

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// POST /api/start
// Optional JSON body: { "image": "jlesage/firefox" }
func startHandler(w http.ResponseWriter, r *http.Request) {
	type bodyReq struct {
		Image string `json:"image,omitempty"`
	}
	ctx := r.Context()

	// Parse optional body
	var b bodyReq
	_ = json.NewDecoder(r.Body).Decode(&b)

	// Defualt to Firefox if no image provided
	image := imageFirefox
	ctrName := "firefox_go"

	if b.Image != "" {
		image = b.Image
	}

	if image == "domistyle/tor-browser" {
		ctrName = "tor_go"
	}

	// Stop and remove any existing container first
	if currentName != "" {
		_ = docker_sdk.StopContainer(ctx, currentName)
		_ = docker_sdk.RemoveContainer(ctx, currentName, true)
		currentName = ""
		currentHostURL = ""
	}

	// Start new container
	ctx2, cancel := context.WithTimeout(ctx, 6*time.Minute)
	defer cancel()

	name, err := docker_sdk.StartContainer(image, ctx2, ctrName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, resp{OK: false, Error: err.Error()})
		return
	}

	iframe := "http://localhost:5800"
	currentName = name
	currentHostURL = iframe

	log.Printf("Started container %s (%s)\n", name, image)
	writeJSON(w, http.StatusOK, resp{OK: true, IframeURL: iframe})
}

// POST /api/stop
func stopHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if currentName == "" {
		writeJSON(w, http.StatusOK, resp{OK: true})
		return
	}

	_ = docker_sdk.StopContainer(ctx, currentName)
	_ = docker_sdk.RemoveContainer(ctx, currentName, false)

	log.Printf("Stopped container %s\n", currentName)
	currentName = ""
	currentHostURL = ""

	writeJSON(w, http.StatusOK, resp{OK: true})
}

// POST /api/restart
func restartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Stop current container if one exists
	if currentName != "" {
		_ = docker_sdk.StopContainer(ctx, currentName)
		_ = docker_sdk.RemoveContainer(ctx, currentName, true)
		currentName = ""
		currentHostURL = ""
	}

	// Restart default (Firefox)
	startHandler(w, r)
}

// POST /api/swap - Switch between Tor and Firefox
func swapHandler(w http.ResponseWriter, r *http.Request) {
	type bodyReq struct {
		To string `json:"to"` // "tor" or "firefox"
	}

	var b bodyReq
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		writeJSON(w, http.StatusBadRequest, resp{OK: false, Error: "Invalid JSON"})
		return
	}

	// Determine which image to use
	var image string
	switch b.To {
	case "tor":
		image = imageTor
	case "firefox":
		image = imageFirefox
	default:
		writeJSON(w, http.StatusBadRequest, resp{OK: false, Error: "Invalid browser type"})
		return
	}

	// Stop current container
	ctx := r.Context()
	if currentName != "" {
		_ = docker_sdk.StopContainer(ctx, currentName)
		_ = docker_sdk.RemoveContainer(ctx, currentName, true)
		currentName = ""
		currentHostURL = ""
	}

	// Start new container with selected image
	ctx2, cancel := context.WithTimeout(ctx, 6*time.Minute)
	defer cancel()

	name, err := docker_sdk.StartContainer(image, ctx2, b.To+"_go")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, resp{OK: false, Error: err.Error()})
		return
	}

	currentName = name
	currentHostURL = "http://localhost:5800"

	log.Printf("Swapped to %s container %s\n", b.To, name)
	writeJSON(w, http.StatusOK, resp{OK: true, IframeURL: currentHostURL})
}

// GET /api/health - Check if container is ready
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if currentName == "" {
		writeJSON(w, http.StatusServiceUnavailable, resp{OK: false, Error: "No container running"})
		return
	}

	// You could also check if port 5800 is actually responding
	// For now, just check if we have a container tracked
	writeJSON(w, http.StatusOK, resp{OK: true})
}

func main() {

	// Serve files from static folder
	http.Handle("/", http.FileServer(http.Dir("./Front-End")))

	mux := http.NewServeMux()

	mux.HandleFunc("/api/start", startHandler)
	mux.HandleFunc("/api/stop", stopHandler)
	mux.HandleFunc("/api/restart", restartHandler)
	mux.HandleFunc("/api/swap", swapHandler)
	mux.HandleFunc("/api/health", healthHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
