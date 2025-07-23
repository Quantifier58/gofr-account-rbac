package model

import (
    "time"
   //  "encoding/json"
)

// User represents a user in the system
type User struct {
    ID           int       `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"` // Hidden from JSON
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserCreateRequest represents the request payload for user creation
type UserCreateRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email,max=100"`
    Password string `json:"password" validate:"required,min=6,max=100"`
}

// UserResponse represents the response payload for user operations
type UserResponse struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts User to UserResponse (excludes sensitive data)
func (u *User) ToResponse() *UserResponse {
    return &UserResponse{
        ID:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}

// ValidateCreateRequest validates user creation request
func (r *UserCreateRequest) Validate() error {
    if r.Username == "" {
        return NewValidationError("username is required")
    }
    if len(r.Username) < 3 || len(r.Username) > 50 {
        return NewValidationError("username must be between 3 and 50 characters")
    }
    if r.Email == "" {
        return NewValidationError("email is required")
    }
    if r.Password == "" {
        return NewValidationError("password is required")
    }
    if len(r.Password) < 6 {
        return NewValidationError("password must be at least 6 characters")
    }
    return nil
}

// ValidationError represents a validation error
type ValidationError struct {
    Message string `json:"message"`
}

func (e *ValidationError) Error() string {
    return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
    return &ValidationError{Message: message}
}
