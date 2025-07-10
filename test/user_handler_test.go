package test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gofiber/fiber/v2"
    "github.com/Quantifier58/gofr-account-rbac/internal/handlers"
)

func setupApp() *fiber.App {
    app := fiber.New()
    app.Post("/users", handlers.RegisterUser)
    app.Get("/users/:id", handlers.GetUser)
    return app
}

func TestRegisterUser(t *testing.T) {
    app := setupApp()
    user := map[string]string{
        "username": "testuser",
        "email": "test@example.com",
        "password": "password123",
    }
    body, _ := json.Marshal(user)
    req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(req)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != 201 {
        t.Errorf("Expected status 201, got %d", resp.StatusCode)
    }
}

func TestGetUser(t *testing.T) {
    app := setupApp()
    req := httptest.NewRequest("GET", "/users/1", nil)
    resp, err := app.Test(req)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != 200 && resp.StatusCode != 404 {
        t.Errorf("Expected status 200 or 404, got %d", resp.StatusCode)
    }
}
