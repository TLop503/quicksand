package web

import (
	"context"
	"encoding/json"
	"github.com/tlop503/quicksand/docker_sdk"
	"log"
	"net/http"
	"time"
)

// resp struct represents HTTP responses w/ field for the iframe used for VNC
type resp struct {
	OK        bool   `json:"ok"`
	IframeURL string `json:"iframeUrl,omitempty"`
	Error     string `json:"error,omitempty"`
}

// writeJSON is a helper for returning JSON http responses
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// HomeHandler serves the root HTML page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the HTML file
	http.ServeFile(w, r, "front_end/index.html")
}

/*
StartHandler calls the "start" api for docker containers
POST /api/start
Optional JSON body: { "image": "publisher/name" }
*/
func StartHandler(w http.ResponseWriter, r *http.Request) {

	// Parse optional request body
	type bodyReq struct {
		Image string `json:"image,omitempty"`
	}
	var b bodyReq
	_ = json.NewDecoder(r.Body).Decode(&b)

	// Parse body for image, and default to firefox if not provided
	var image, ctrName string
	if b.Image == docker_sdk.TorImage {
		image = b.Image
		ctrName = "tor_quicksand"
	} else {
		image = docker_sdk.FireFoxImage
		ctrName = "firefox_quicksand"
	}

	// Stop and remove any existing container before starting new instance
	StopHandler(w, r)

	// Start new container
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 6*time.Minute)
	defer cancel()
	name, err := docker_sdk.StartContainer(image, ctx, ctrName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, resp{OK: false, Error: err.Error()})
		return
	}

	// update globals to store state
	// TODO: does url ever change?
	docker_sdk.CurrentCtr = name

	log.Printf("Started container %s (%s)\n", name, image)
	writeJSON(w, http.StatusOK, resp{OK: true, IframeURL: docker_sdk.CurrentHostURL})
}

/*
StopHandler stops and removes the current container
POST /api/stop
*/
func StopHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if docker_sdk.CurrentCtr == "" {
		writeJSON(w, http.StatusOK, resp{OK: true})
		return
	}

	_ = docker_sdk.StopContainer(ctx, docker_sdk.CurrentCtr)
	_ = docker_sdk.RemoveContainer(ctx, docker_sdk.CurrentCtr, false)

	log.Printf("Stopped container %s\n", docker_sdk.CurrentCtr)
	docker_sdk.CurrentCtr = ""
	docker_sdk.CurrentHostURL = ""

	writeJSON(w, http.StatusOK, resp{OK: true})
}

/*
RestartHandler just calls stop and then start
POST /api/restart
*/
func RestartHandler(w http.ResponseWriter, r *http.Request) {
	StopHandler(w, r)
	StartHandler(w, r)
}

/*
SwapHandler toggles active containers between Tor and Firefox
POST /api/restart
*/
func SwapHandler(w http.ResponseWriter, r *http.Request) {
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
		image = docker_sdk.TorImage
	case "firefox":
		image = docker_sdk.FireFoxImage
	default:
		writeJSON(w, http.StatusBadRequest, resp{OK: false, Error: "Invalid browser type"})
		return
	}

	// Stop current container
	StopHandler(w, r)

	// Start new container with selected image
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 6*time.Minute)
	defer cancel()
	name, err := docker_sdk.StartContainer(image, ctx, b.To+"_quicksand")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, resp{OK: false, Error: err.Error()})
		return
	}

	docker_sdk.CurrentCtr = name

	log.Printf("Swapped to %s container %s\n", b.To, name)
	writeJSON(w, http.StatusOK, resp{OK: true, IframeURL: docker_sdk.CurrentHostURL})
}

/*
HealthHandler returns if the intended current container is running
// GET /api/health
*/
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if docker_sdk.CurrentCtr == "" {
		writeJSON(w, http.StatusServiceUnavailable, resp{OK: false, Error: "No container running"})
		return
	}
	// TODO: expand ad hoc to verify VNC port 5800 is active
	writeJSON(w, http.StatusOK, resp{OK: true})
}
