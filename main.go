package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tlop503/quicksand/docker_sdk"
	"github.com/tlop503/quicksand/web"
)

func catchSIGTERM() error {
	ctr := docker_sdk.CurrentCtr
	fmt.Println(" Sigterm caught!")
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
	fs := http.FileServer(http.Dir("./front_end"))
	mux.Handle("/front_end/", http.StripPrefix("/front_end/", fs))

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
	err = srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server closed. Goodbye!")
	} else if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
