package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "log"
    "github.com/Quantifier58/gofr-account-rbac/internal/model"
    "github.com/Quantifier58/gofr-account-rbac/internal/service"
    "github.com/Quantifier58/gofr-account-rbac/internal/repository"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
    Service *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService) *UserHandler {
    return &UserHandler{Service: service}
}

// RegisterUser handles user registration requests
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req model.UserCreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("Failed to decode request: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    user, err := h.Service.Register(r.Context(), &req)
    if err != nil {
        log.Printf("Failed to register user: %v", err)
        
        // Check for validation errors
        if _, ok := err.(*model.ValidationError); ok {
            respondWithError(w, http.StatusBadRequest, err.Error())
            return
        }
        
        // Check for service errors (like user already exists)
        if _, ok := err.(*service.ServiceError); ok {
            respondWithError(w, http.StatusConflict, err.Error())
            return
        }
        
        respondWithError(w, http.StatusInternalServerError, "Failed to create user")
        return
    }
    
    respondWithJSON(w, http.StatusCreated, user.ToResponse())
}

// GetUser handles user retrieval requests
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Check query parameters for different lookup methods
    username := r.URL.Query().Get("username")
    email := r.URL.Query().Get("email")
    idStr := r.URL.Query().Get("id")
    
    var user *model.User
    var err error
    
    switch {
    case username != "":
        user, err = h.Service.GetByUsername(r.Context(), username)
    case email != "":
        user, err = h.Service.GetByEmail(r.Context(), email)
    case idStr != "":
        id, parseErr := strconv.Atoi(idStr)
        if parseErr != nil {
            respondWithError(w, http.StatusBadRequest, "Invalid user ID")
            return
        }
        user, err = h.Service.GetByID(r.Context(), id)
    default:
        respondWithError(w, http.StatusBadRequest, "Please provide username, email, or id parameter")
        return
    }
    
    if err != nil {
        log.Printf("Failed to get user: %v", err)
        
        // Check for service errors (like user not found)
        if _, ok := err.(*service.ServiceError); ok {
            respondWithError(w, http.StatusBadRequest, err.Error())
            return
        }
        
        // Check for repository errors (like user not found)
        if _, ok := err.(*repository.RepositoryError); ok {
            respondWithError(w, http.StatusNotFound, "User not found")
            return
        }
        
        respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user")
        return
    }
    
    respondWithJSON(w, http.StatusOK, user.ToResponse())
}

// HealthCheck handles health check requests
func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    response := map[string]string{
        "status":  "healthy",
        "service": "gofr-account-service",
    }
    
    respondWithJSON(w, http.StatusOK, response)
}

// Helper functions for consistent JSON responses
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Printf("Failed to encode JSON response: %v", err)
    }
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
    errorResponse := map[string]string{"error": message}
    respondWithJSON(w, statusCode, errorResponse)
}
