package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/devaraja-anu/eyo/server/internals/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbConfig struct {
	dsn         string
	maxIdleTime string
}

type cfg struct {
	port int
	db   dbConfig
}

type application struct {
	cfg     cfg
	logger  *log.Logger
	queries *db.Queries
}

func main() {

	logger := log.New(os.Stdout, "Log: ", log.Ldate|log.Ltime)

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		logger.Fatal("DATABASE_DSN not set")
	}

	cfg := cfg{
		port: 4000,
		db: dbConfig{
			maxIdleTime: "15m",
			dsn:         dsn,
		},
	}

	conn, err := OpenDB(cfg)

	if err != nil {
		logger.Fatal(err)
	}

	defer conn.Close()

	queries := db.New(conn)

	app := &application{
		cfg:     cfg,
		logger:  logger,
		queries: queries,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting server  on %v", app.cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
		return
	}

}

func OpenDB(cfg cfg) (*pgxpool.Pool, error) {
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConnIdleTime = duration

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil

}
