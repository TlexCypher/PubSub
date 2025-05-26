package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	port = os.Getenv("PORT")
)

type Server interface {
	Run(context.Context, map[string]func(http.ResponseWriter, *http.Request)) error
}

type ApplicationServer struct{}

func NewApplicationServer() *ApplicationServer {
	return &ApplicationServer{}
}

func (s *ApplicationServer) Run(baseCtx context.Context, configs map[string]func(http.ResponseWriter, *http.Request)) error {
	eg, ctx := errgroup.WithContext(baseCtx)
	mux := http.NewServeMux()
	for pattern, handler := range configs {
		mux.HandleFunc(pattern, handler)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}

	eg.Go(func() error {
		log.Println("ApplicationServer: Starting application server...")
		if err := srv.ListenAndServe(); err != nil {
			return fmt.Errorf("http.Server.ListenAndServe failed: %w", err)
		}
		log.Println("ApplicationServer: Gracefully exited application server...")
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Println("Attempting to gracefully shutdown HTTP server...")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("http.Server.Shutdown failed: %v", err)
			return fmt.Errorf("http.Server.Shutdown failed: %w", err)
		}
		log.Println("HTTP server gracefully shutdown.")
		return nil
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("main http server run failed: %w", err)
	}
	log.Println("ApplicationServer.Run: errgroup finished successfully.")
	return nil
}
