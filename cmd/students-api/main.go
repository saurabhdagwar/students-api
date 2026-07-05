package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saurabhdagwar/students-api/internal/config"
)

func main() {
	cfg := config.MustLoad()
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Students API"))
	})
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server Started", slog.String("address", cfg.Addr))
	fmt.Printf("Server Started %s", cfg.HTTPServer.Addr)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done
	slog.Info("Shutting down the server")
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to Shutdown Server", slog.String("error", err.Error()))
	}
	slog.Info("Server Shutdown Successfully")
}
