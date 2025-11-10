// package main

// import (
// 	"context"
// 	"fmt"

// 	"github.com/tlop503/quicksand/docker_sdk"
// )

// func main() {
// 	ctx := context.Background()

// 	containerName, err := docker_sdk.StartContainer("jlesage/firefox", ctx, "firefox_go")
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = docker_sdk.StopContainer(ctx, containerName)
// 	if err != nil {
// 		fmt.Printf("Error stopping container: %v\n", err)
// 	} else {
// 		fmt.Println("Container stopped successfully")
// 	}

// 	// Remove the container (using force=true in case stop failed)
// 	err = docker_sdk.RemoveContainer(ctx, containerName, false)
// 	if err != nil {
// 		fmt.Printf("Error removing container: %v\n", err)
// 	} else {
// 		fmt.Println("Container removed successfully")
// 	}
// }

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/tlop503/quicksand/docker_sdk"
)

type resp struct {
	OK        bool   `json:"ok"`
	IframeURL string `json:"iframeUrl,omitempty"`
	Error     string `json:"error,omitempty"`
}

var (
	mu             sync.Mutex
	currentName    string // container name returned by docker_sdk.StartContainer
	currentHostURL string // e.g. http://localhost:5800
	imageFirefox   = "jlesage/firefox"
	imageTor       = "jlesage/tor-browser"
)

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// accept optional {"image":"jlesage/firefox"} or {"image":"tor"}
func startHandler(w http.ResponseWriter, r *http.Request) {
	type bodyReq struct {
		Image string `json:"image,omitempty"`
	}
	ctx := r.Context()

	// parse optional body
	var b bodyReq
	_ = json.NewDecoder(r.Body).Decode(&b)
	image := imageFirefox
	if b.Image != "" {
		image = b.Image
	}

	// ensure only one active container at a time
	mu.Lock()
	// if one running, stop/remove first
	prev := currentName
	mu.Unlock()

	if prev != "" {
		// best-effort stop/remove
		_ = docker_sdk.StopContainer(ctx, prev)
		_ = docker_sdk.RemoveContainer(ctx, prev, true)
		// reset
		mu.Lock()
		currentName = ""
		currentHostURL = ""
		mu.Unlock()
	}

	// start with timeout context
	ctx2, cancel := context.WithTimeout(ctx, 6*time.Minute)
	defer cancel()

	name, hostPort, err := docker_sdk.StartContainer(image, ctx2, "firefox_go")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, resp{OK: false, Error: err.Error()})
		return
	}

	if hostPort == "" {
		writeJSON(w, http.StatusInternalServerError, resp{OK: false, Error: "Could not determine host port for container"})
		return
	}

	iframe := "http://localhost:" + hostPort

	mu.Lock()
	currentName = name
	currentHostURL = iframe
	mu.Unlock()

	writeJSON(w, http.StatusOK, resp{OK: true, IframeURL: iframe})
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mu.Lock()
	name := currentName
	mu.Unlock()

	if name == "" {
		writeJSON(w, http.StatusOK, resp{OK: true})
		return
	}

	_ = docker_sdk.StopContainer(ctx, name)
	_ = docker_sdk.RemoveContainer(ctx, name, true)

	mu.Lock()
	currentName = ""
	currentHostURL = ""
	mu.Unlock()

	writeJSON(w, http.StatusOK, resp{OK: true})
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	// simpler: call stop then start (start uses default firefox)
	stopHandler(w, r)
	startHandler(w, r)
}

func main() {
	mux := http.NewServeMux()
	// API used by script.js (see uploaded script.js for calls). :contentReference[oaicite:2]{index=2}
	mux.HandleFunc("/api/start", startHandler)
	mux.HandleFunc("/api/stop", stopHandler)
	mux.HandleFunc("/api/restart", restartHandler)

	// static files (if you serve a UI)
	// mux.Handle("/", http.FileServer(http.Dir("./web")))

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
