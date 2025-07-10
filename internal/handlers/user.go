package handlers

import (
    "context"
    "net/http"
    "strconv"

    "github.com/gofiber/fiber/v2"
    "github.com/Quantifier58/gofr-account-rbac/internal/db"
    "github.com/Quantifier58/gofr-account-rbac/internal/models"
    "github.com/Quantifier58/gofr-account-rbac/internal/utils"
)

type UserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func RegisterUser(c *fiber.Ctx) error {
    var req UserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }
    if req.Username == "" || req.Email == "" || req.Password == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Username, email and password are required"})
    }
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
    }
    user := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
    }
    ctx := context.Background()
    err = models.CreateUser(ctx, db.DB, user)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
    }
    return c.Status(http.StatusCreated).JSON(fiber.Map{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
        "created":  user.CreatedAt,
    })
}

func GetUser(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
    }
    ctx := context.Background()
    user, err := models.GetUserByID(ctx, db.DB, id)
    if err != nil {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }
    return c.JSON(fiber.Map{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
        "created":  user.CreatedAt,
    })
}
