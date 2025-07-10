package db

import (
    "context"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
)

var DB *pgxpool.Pool

func ConnectDB() (*pgxpool.Pool, error) {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        connStr = "postgres://gofr:gofrpass@localhost:5432/gofrdb?sslmode=disable"
    }
    pool, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        return nil, err
    }
    DB = pool
    return pool, nil
}

func RunMigrations() error {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        connStr = "postgres://gofr:gofrpass@localhost:5432/gofrdb?sslmode=disable"
    }
    driver, err := postgres.WithInstance(DB, &postgres.Config{})
    if err != nil {
        return err
    }
    m, err := migrate.NewWithDatabaseInstance(
        "file://internal/db/migrations",
        "postgres", driver)
    if err != nil {
        return err
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    return nil
}

func CloseDB() {
    if DB != nil {
        DB.Close()
    }
}
