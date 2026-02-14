package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/devaraja-anu/eyo/server/internals/data"
	"github.com/devaraja-anu/eyo/server/internals/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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
		logger  *logger.Logger
		models 	data.Models
	}


func main() {

	defaultLogger := log.New(os.Stdout, "MAIN LOG: ", log.Ldate|log.Ltime)


	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	godotenv.Load("../.env")

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		defaultLogger.Fatal("DB_URL not set")
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
		defaultLogger.Fatal(err)
	}

	defer conn.Close()

	models := data.NewModels(conn)

	app := &application{
		cfg:     cfg,
		logger:  logger.New(infoLog, errorLog),
		models: models,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	defaultLogger.Printf("starting server  on %v", app.cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		defaultLogger.Fatal(err)
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
