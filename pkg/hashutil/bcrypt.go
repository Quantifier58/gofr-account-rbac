package hashutil

import (
    "golang.org/x/crypto/bcrypt"
)

const (
    // DefaultCost is the default bcrypt cost
    DefaultCost = bcrypt.DefaultCost
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
    if password == "" {
        return "", NewHashError("password cannot be empty")
    }
    
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
    if err != nil {
        return "", NewHashError("failed to hash password: " + err.Error())
    }
    
    return string(bytes), nil
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
    if password == "" || hash == "" {
        return false
    }
    
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// HashError represents a hashing error
type HashError struct {
    Message string
}

func (e *HashError) Error() string {
    return e.Message
}

// NewHashError creates a new hash error
func NewHashError(message string) *HashError {
    return &HashError{Message: message}
}
