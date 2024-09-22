package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edulustosa/galleria/internal/api/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,
		os.Kill,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println("To infinity and beyond!")
}

func run(ctx context.Context) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	jwtKey := os.Getenv("JWT_SECRET")
	srv := router.NewServer(pool, jwtKey)
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      srv,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	defer shutdown(httpServer)

	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("server is listening on %s\n", httpServer.Addr)
		errChan <- httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func shutdown(srv *http.Server) {
	const timeout = 30 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}
}
