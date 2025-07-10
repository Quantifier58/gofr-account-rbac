package main

import (
    "log"

    "github.com/gofiber/fiber/v2"
    "github.com/Quantifier58/gofr-account-rbac/internal/db"
    "github.com/Quantifier58/gofr-account-rbac/internal/handlers"
)

func main() {
    pool, err := db.ConnectDB()
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }
    defer pool.Close()

    if err := db.RunMigrations(); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    app := fiber.New()
    app.Post("/users", handlers.RegisterUser)
    app.Get("/users/:id", handlers.GetUser)

    log.Println("Starting server on :8080")
    if err := app.Listen(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
