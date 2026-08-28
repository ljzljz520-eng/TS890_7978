package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"petcolor/internal/media"
	"petcolor/internal/preview"
	"petcolor/internal/service"
	"petcolor/internal/storage"
	"petcolor/internal/web"
	"syscall"
	"time"
)

type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now().UTC()
}

func main() {
	if err := run(); err != nil {
		log.Printf("petgrade stopped: %v", err)
		os.Exit(1)
	}
}

func run() error {
	address := flag.String("addr", "127.0.0.1:8080", "HTTP listen address")
	dataRoot := flag.String("data", "./petgrade-data", "temporary data directory")
	flag.Parse()

	if err := os.MkdirAll(*dataRoot, 0755); err != nil {
		return fmt.Errorf("create data root: %w", err)
	}
	store, err := storage.Open(filepath.Join(*dataRoot, "petgrade.db"))
	if err != nil {
		return err
	}
	defer store.Close()

	uploads, err := media.NewUploadManager(filepath.Join(*dataRoot, "media"), 512*1024*1024)
	if err != nil {
		return err
	}
	planner, err := preview.NewPlanner(960, 540, 12)
	if err != nil {
		return err
	}
	application, err := service.New(store, uploads, planner, systemClock{})
	if err != nil {
		return err
	}
	handler, err := web.NewServer(application)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              *address,
		Handler:           handler.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-shutdownSignal.Done()
		shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownContext)
	}()

	log.Printf("pet color workbench listening on http://%s", *address)
	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
