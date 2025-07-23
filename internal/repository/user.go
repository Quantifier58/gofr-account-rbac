package repository

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/Quantifier58/gofr-account-rbac/internal/model"
    _ "github.com/lib/pq"
)

// UserRepositoryInterface defines the contract for user repository operations
type UserRepositoryInterface interface {
    CreateUser(ctx context.Context, user *model.User) error
    GetUserByUsername(ctx context.Context, username string) (*model.User, error)
    GetUserByEmail(ctx context.Context, email string) (*model.User, error)
    GetUserByID(ctx context.Context, id int) (*model.User, error)
    CheckUserExists(ctx context.Context, username, email string) (bool, error)
}

// UserRepository handles database operations for users
type UserRepository struct {
    DB *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepositoryInterface {
    return &UserRepository{DB: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
    query := `
        INSERT INTO users (username, email, password_hash) 
        VALUES ($1, $2, $3) 
        RETURNING id, created_at, updated_at`
    
    err := r.DB.QueryRowContext(
        ctx, 
        query, 
        user.Username, 
        user.Email, 
        user.PasswordHash,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
    user := &model.User{}
    query := `
        SELECT id, username, email, password_hash, created_at, updated_at 
        FROM users 
        WHERE username = $1`
    
    err := r.DB.QueryRowContext(ctx, query, username).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, NewRepositoryError("user not found")
        }
        return nil, fmt.Errorf("failed to get user by username: %w", err)
    }
    
    return user, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    user := &model.User{}
    query := `
        SELECT id, username, email, password_hash, created_at, updated_at 
        FROM users 
        WHERE email = $1`
    
    err := r.DB.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, NewRepositoryError("user not found")
        }
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }
    
    return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
    user := &model.User{}
    query := `
        SELECT id, username, email, password_hash, created_at, updated_at 
        FROM users 
        WHERE id = $1`
    
    err := r.DB.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, NewRepositoryError("user not found")
        }
        return nil, fmt.Errorf("failed to get user by ID: %w", err)
    }
    
    return user, nil
}

// CheckUserExists checks if a user exists by username or email
func (r *UserRepository) CheckUserExists(ctx context.Context, username, email string) (bool, error) {
    var count int
    query := `SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2`
    
    err := r.DB.QueryRowContext(ctx, query, username, email).Scan(&count)
    if err != nil {
        return false, fmt.Errorf("failed to check user existence: %w", err)
    }
    
    return count > 0, nil
}

// RepositoryError represents a repository error
type RepositoryError struct {
    Message string
}

func (e *RepositoryError) Error() string {
    return e.Message
}

// NewRepositoryError creates a new repository error
func NewRepositoryError(message string) *RepositoryError {
    return &RepositoryError{Message: message}
}
