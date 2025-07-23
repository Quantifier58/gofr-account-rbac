package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/Quantifier58/gofr-account-rbac/internal/handler"
    "github.com/Quantifier58/gofr-account-rbac/internal/repository"
    "github.com/Quantifier58/gofr-account-rbac/internal/service"
    _ "github.com/lib/pq"
)

func main() {
    // Load configuration from environment variables
    config := loadConfig()
    
    // Connect to database
    db, err := connectDB(config)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // Initialize repositories - returns UserRepositoryInterface
    userRepo := repository.NewUserRepository(db)
    
    // Initialize services
    userService := service.NewUserService(userRepo)
    
    // Initialize handlers
    userHandler := handler.NewUserHandler(userService)
    
    // Setup routes
    setupRoutes(userHandler)
    
    // Start server
    port := config.AppPort
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Config holds application configuration
type Config struct {
    DBHost     string
    DBUser     string
    DBPassword string
    DBName     string
    DBPort     string
    AppPort    string
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBUser:     getEnv("DB_USER", "gouser"),
        DBPassword: getEnv("DB_PASSWORD", "pass"),
        DBName:     getEnv("DB_NAME", "accounts"),
        DBPort:     getEnv("DB_PORT", "5432"),
        AppPort:    getEnv("APP_PORT", "8080"),
    }
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

// connectDB establishes database connection
func connectDB(config *Config) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.DBHost,
        config.DBPort,
        config.DBUser,
        config.DBPassword,
        config.DBName,
    )
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("Successfully connected to database")
    return db, nil
}

// setupRoutes configures HTTP routes
func setupRoutes(userHandler *handler.UserHandler) {
    http.HandleFunc("/health", userHandler.HealthCheck)
    http.HandleFunc("/register", userHandler.RegisterUser)
    http.HandleFunc("/user", userHandler.GetUser)
    
    // Add CORS middleware for all routes
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        http.NotFound(w, r)
    })
}
