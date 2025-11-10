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
	imageTor       = "jlesage/tor-browser"
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
	image := imageFirefox
	if b.Image != "" {
		image = b.Image
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

	name, err := docker_sdk.StartContainer(image, ctx2, "firefox_go")
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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/start", startHandler)
	mux.HandleFunc("/api/stop", stopHandler)
	mux.HandleFunc("/api/restart", restartHandler)

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
