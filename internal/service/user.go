package service

import (
    "context"
    "fmt"
    "strings"
    "github.com/Quantifier58/gofr-account-rbac/internal/model"
    "github.com/Quantifier58/gofr-account-rbac/internal/repository"
    "github.com/Quantifier58/gofr-account-rbac/pkg/hashutil"
)

// UserService handles business logic for user operations
type UserService struct {
    Repo repository.UserRepositoryInterface
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepositoryInterface) *UserService {
    return &UserService{Repo: repo}
}

// Register creates a new user account
func (s *UserService) Register(ctx context.Context, req *model.UserCreateRequest) (*model.User, error) {
    // Validate request
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Normalize input
    req.Username = strings.TrimSpace(strings.ToLower(req.Username))
    req.Email = strings.TrimSpace(strings.ToLower(req.Email))
    
    // Check if user already exists
    exists, err := s.Repo.CheckUserExists(ctx, req.Username, req.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to check user existence: %w", err)
    }
    if exists {
        return nil, NewServiceError("user with this username or email already exists")
    }
    
    // Hash password
    hashedPassword, err := hashutil.HashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Create user model
    user := &model.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
    }
    
    // Save to database
    err = s.Repo.CreateUser(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}

// GetByUsername retrieves a user by username
func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
    if username == "" {
        return nil, NewServiceError("username is required")
    }
    
    username = strings.TrimSpace(strings.ToLower(username))
    
    user, err := s.Repo.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    if email == "" {
        return nil, NewServiceError("email is required")
    }
    
    email = strings.TrimSpace(strings.ToLower(email))
    
    user, err := s.Repo.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
    if id <= 0 {
        return nil, NewServiceError("invalid user ID")
    }
    
    user, err := s.Repo.GetUserByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

// ValidateCredentials validates user credentials for login
func (s *UserService) ValidateCredentials(ctx context.Context, username, password string) (*model.User, error) {
    if username == "" || password == "" {
        return nil, NewServiceError("username and password are required")
    }
    
    user, err := s.GetByUsername(ctx, username)
    if err != nil {
        return nil, NewServiceError("invalid credentials")
    }
    
    if !hashutil.CheckPasswordHash(password, user.PasswordHash) {
        return nil, NewServiceError("invalid credentials")
    }
    
    return user, nil
}

// ServiceError represents a service layer error
type ServiceError struct {
    Message string
}

func (e *ServiceError) Error() string {
    return e.Message
}

// NewServiceError creates a new service error
func NewServiceError(message string) *ServiceError {
    return &ServiceError{Message: message}
}
