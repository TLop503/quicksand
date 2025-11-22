package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tlop503/quicksand/docker_sdk"
	"github.com/tlop503/quicksand/web"
)

var (
	currentName    string // name returned by docker_sdk.StartContainer
	currentHostURL string // e.g. http://localhost:5800
	imageFirefox   = "jlesage/firefox"
	imageTor       = "domistyle/tor-browser"
	imageDefault   = imageFirefox
)

func catchSIGTERM() error {
	ctr := currentName
	fmt.Println("Sigterm caught!")
	req, err := http.NewRequest("POST", "http://localhost:8080/api/stop", bytes.NewBuffer(nil))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	log.Printf("%s stopped successfully.", ctr)
	return nil
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err := catchSIGTERM()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	// pull both containers in advance to save time later
	ctx := context.Background()
	err := docker_sdk.PullAll(ctx)
	if err != nil {
		log.Println(err)
	}

	mux := http.NewServeMux()

	// Serve Front-End files (CSS, JS)
	fs := http.FileServer(http.Dir("./Front-End"))
	mux.Handle("/Front-End/", http.StripPrefix("/Front-End/", fs))

	mux.HandleFunc("/api/start", web.StartHandler)
	mux.HandleFunc("/api/stop", web.StopHandler)
	mux.HandleFunc("/api/restart", web.RestartHandler)
	mux.HandleFunc("/api/swap", web.SwapHandler)
	mux.HandleFunc("/api/health", web.HealthHandler)

	// Serve HTML page at root
	mux.HandleFunc("/", web.HomeHandler)

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
