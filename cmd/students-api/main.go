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
	"github.com/saurabhdagwar/students-api/internal/http/handlers/student"
	"github.com/saurabhdagwar/students-api/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()
	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage Initialize")

	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetByID(storage))
	router.HandleFunc("GET /api/students/", student.GetList(storage))

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
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to Shutdown Server", slog.String("error", err.Error()))
	}
	slog.Info("Server Shutdown Successfully")
}
