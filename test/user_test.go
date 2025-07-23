package test

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/Quantifier58/gofr-account-rbac/internal/handler"
    "github.com/Quantifier58/gofr-account-rbac/internal/model"
    "github.com/Quantifier58/gofr-account-rbac/internal/repository"
    "github.com/Quantifier58/gofr-account-rbac/internal/service"
    "github.com/Quantifier58/gofr-account-rbac/pkg/hashutil"
)

// MockUserRepository implements repository.UserRepositoryInterface for testing
type MockUserRepository struct {
    users     map[string]*model.User
    usersByID map[int]*model.User
    idCounter int
}

// Ensure MockUserRepository implements the interface at compile time
var _ repository.UserRepositoryInterface = (*MockUserRepository)(nil)

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users:     make(map[string]*model.User),
        usersByID: make(map[int]*model.User),
        idCounter: 0,
    }
}

// CreateUser creates a new user in the mock repository
func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) error {
    // Check if user already exists
    if _, exists := m.users[user.Username]; exists {
        return &repository.RepositoryError{Message: "user already exists"}
    }
    
    // Check if email already exists
    for _, existingUser := range m.users {
        if existingUser.Email == user.Email {
            return &repository.RepositoryError{Message: "email already exists"}
        }
    }
    
    m.idCounter++
    user.ID = m.idCounter
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    m.users[user.Username] = user
    m.usersByID[user.ID] = user
    
    return nil
}

// GetUserByUsername retrieves a user by username from the mock repository
func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
    if user, exists := m.users[username]; exists {
        return user, nil
    }
    return nil, &repository.RepositoryError{Message: "user not found"}
}

// GetUserByEmail retrieves a user by email from the mock repository
func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    for _, user := range m.users {
        if user.Email == email {
            return user, nil
        }
    }
    return nil, &repository.RepositoryError{Message: "user not found"}
}

// GetUserByID retrieves a user by ID from the mock repository
func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
    if user, exists := m.usersByID[id]; exists {
        return user, nil
    }
    return nil, &repository.RepositoryError{Message: "user not found"}
}

// CheckUserExists checks if a user exists in the mock repository
func (m *MockUserRepository) CheckUserExists(ctx context.Context, username, email string) (bool, error) {
    // Check username
    if _, exists := m.users[username]; exists {
        return true, nil
    }
    
    // Check email
    for _, user := range m.users {
        if user.Email == email {
            return true, nil
        }
    }
    
    return false, nil
}

// Test user registration
func TestRegisterUser(t *testing.T) {
    // Setup
    mockRepo := NewMockUserRepository()
    userService := service.NewUserService(mockRepo)
    userHandler := handler.NewUserHandler(userService)
    
    tests := []struct {
        name           string
        payload        map[string]string
        expectedStatus int
        expectedError  string
    }{
        {
            name: "successful registration",
            payload: map[string]string{
                "username": "testuser",
                "email":    "test@example.com",
                "password": "password123",
            },
            expectedStatus: http.StatusCreated,
        },
        {
            name: "missing username",
            payload: map[string]string{
                "email":    "test@example.com",
                "password": "password123",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "username is required",
        },
        {
            name: "missing email",
            payload: map[string]string{
                "username": "testuser",
                "password": "password123",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "email is required",
        },
        {
            name: "short password",
            payload: map[string]string{
                "username": "testuser",
                "email":    "test@example.com",
                "password": "123",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "password must be at least 6 characters",
        },
        {
            name: "empty username",
            payload: map[string]string{
                "username": "",
                "email":    "test@example.com",
                "password": "password123",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "username is required",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create request
            payloadBytes, _ := json.Marshal(tt.payload)
            req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payloadBytes))
            req.Header.Set("Content-Type", "application/json")
            
            // Create response recorder
            w := httptest.NewRecorder()
            
            // Call handler
            userHandler.RegisterUser(w, req)
            
            // Check status code
            if w.Code != tt.expectedStatus {
                t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
            }
            
            // Check error message if expected
            if tt.expectedError != "" {
                var response map[string]string
                if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
                    t.Fatalf("Failed to decode error response: %v", err)
                }
                if response["error"] != tt.expectedError {
                    t.Errorf("Expected error '%s', got '%s'", tt.expectedError, response["error"])
                }
            }
            
            // Check successful response
            if tt.expectedStatus == http.StatusCreated {
                var response model.UserResponse
                if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
                    t.Fatalf("Failed to decode success response: %v", err)
                }
                if response.Username != tt.payload["username"] {
                    t.Errorf("Expected username '%s', got '%s'", tt.payload["username"], response.Username)
                }
                if response.Email != tt.payload["email"] {
                    t.Errorf("Expected email '%s', got '%s'", tt.payload["email"], response.Email)
                }
                if response.ID == 0 {
                    t.Error("Expected non-zero user ID")
                }
            }
        })
    }
}

// Test get user
func TestGetUser(t *testing.T) {
    // Setup
    mockRepo := NewMockUserRepository()
    userService := service.NewUserService(mockRepo)
    userHandler := handler.NewUserHandler(userService)
    
    // Create a test user
    hashedPassword, err := hashutil.HashPassword("password123")
    if err != nil {
        t.Fatalf("Failed to hash password: %v", err)
    }
    
    testUser := &model.User{
        Username:     "testuser",
        Email:        "test@example.com",
        PasswordHash: hashedPassword,
    }
    
    if err := mockRepo.CreateUser(context.Background(), testUser); err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
    
    tests := []struct {
        name           string
        queryParam     string
        queryValue     string
        expectedStatus int
        expectUser     bool
    }{
        {
            name:           "get by username",
            queryParam:     "username",
            queryValue:     "testuser",
            expectedStatus: http.StatusOK,
            expectUser:     true,
        },
        {
            name:           "get by email",
            queryParam:     "email",
            queryValue:     "test@example.com",
            expectedStatus: http.StatusOK,
            expectUser:     true,
        },
        {
            name:           "get by id",
            queryParam:     "id",
            queryValue:     "1",
            expectedStatus: http.StatusOK,
            expectUser:     true,
        },
        {
            name:           "user not found by username",
            queryParam:     "username",
            queryValue:     "nonexistent",
            expectedStatus: http.StatusNotFound,
            expectUser:     false,
        },
        {
            name:           "user not found by email",
            queryParam:     "email",
            queryValue:     "nonexistent@example.com",
            expectedStatus: http.StatusNotFound,
            expectUser:     false,
        },
        {
            name:           "missing query parameter",
            queryParam:     "",
            queryValue:     "",
            expectedStatus: http.StatusBadRequest,
            expectUser:     false,
        },
        {
            name:           "invalid user id",
            queryParam:     "id",
            queryValue:     "invalid",
            expectedStatus: http.StatusBadRequest,
            expectUser:     false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create request
            url := "/user"
            if tt.queryParam != "" && tt.queryValue != "" {
                url += "?" + tt.queryParam + "=" + tt.queryValue
            }
            req := httptest.NewRequest(http.MethodGet, url, nil)
            
            // Create response recorder
            w := httptest.NewRecorder()
            
            // Call handler
            userHandler.GetUser(w, req)
            
            // Check status code
            if w.Code != tt.expectedStatus {
                t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
            }
            
            // Check user response if expected
            if tt.expectUser {
                var response model.UserResponse
                if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
                    t.Fatalf("Failed to decode user response: %v", err)
                }
                if response.Username != "testuser" {
                    t.Errorf("Expected username 'testuser', got '%s'", response.Username)
                }
                if response.Email != "test@example.com" {
                    t.Errorf("Expected email 'test@example.com', got '%s'", response.Email)
                }
                if response.ID != 1 {
                    t.Errorf("Expected user ID 1, got %d", response.ID)
                }
            }
        })
    }
}

// Test password hashing utility
func TestPasswordHashing(t *testing.T) {
    tests := []struct {
        name     string
        password string
        wantErr  bool
    }{
        {
            name:     "valid password",
            password: "testpassword123",
            wantErr:  false,
        },
        {
            name:     "empty password",
            password: "",
            wantErr:  true,
        },
        {
            name:     "long password",
            password: "thisissuperlongpasswordtotestlimits1234567890",
            wantErr:  false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Hash password
            hash, err := hashutil.HashPassword(tt.password)
            
            if tt.wantErr {
                if err == nil {
                    t.Error("Expected error but got none")
                }
                return
            }
            
            if err != nil {
                t.Fatalf("Failed to hash password: %v", err)
            }
            
            if hash == "" {
                t.Error("Expected non-empty hash")
            }
            
            // Verify correct password
            if !hashutil.CheckPasswordHash(tt.password, hash) {
                t.Error("Password hash verification failed for correct password")
            }
            
            // Verify incorrect password
            if hashutil.CheckPasswordHash("wrongpassword", hash) {
                t.Error("Password hash verification passed for incorrect password")
            }
            
            // Test empty password check
            if hashutil.CheckPasswordHash("", hash) {
                t.Error("Password hash verification passed for empty password")
            }
            
            // Test empty hash check
            if hashutil.CheckPasswordHash(tt.password, "") {
                t.Error("Password hash verification passed for empty hash")
            }
        })
    }
}

// Test user service layer
func TestUserService(t *testing.T) {
    mockRepo := NewMockUserRepository()
    userService := service.NewUserService(mockRepo)
    ctx := context.Background()
    
    t.Run("successful user registration", func(t *testing.T) {
        req := &model.UserCreateRequest{
            Username: "testuser",
            Email:    "test@example.com",
            Password: "password123",
        }
        
        user, err := userService.Register(ctx, req)
        if err != nil {
            t.Fatalf("Failed to register user: %v", err)
        }
        
        if user.Username != "testuser" {
            t.Errorf("Expected username 'testuser', got '%s'", user.Username)
        }
        
        if user.Email != "test@example.com" {
            t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
        }
        
        if user.ID == 0 {
            t.Error("Expected non-zero user ID")
        }
        
        if user.PasswordHash == "" {
            t.Error("Expected non-empty password hash")
        }
    })
    
    t.Run("duplicate user registration", func(t *testing.T) {
        req := &model.UserCreateRequest{
            Username: "testuser", // Same as above
            Email:    "different@example.com",
            Password: "password123",
        }
        
        _, err := userService.Register(ctx, req)
        if err == nil {
            t.Error("Expected error for duplicate username")
        }
    })
    
    t.Run("duplicate email registration", func(t *testing.T) {
        req := &model.UserCreateRequest{
            Username: "differentuser",
            Email:    "test@example.com", // Same as above
            Password: "password123",
        }
        
        _, err := userService.Register(ctx, req)
        if err == nil {
            t.Error("Expected error for duplicate email")
        }
    })
    
    t.Run("get user by username", func(t *testing.T) {
        retrievedUser, err := userService.GetByUsername(ctx, "testuser")
        if err != nil {
            t.Fatalf("Failed to get user by username: %v", err)
        }
        
        if retrievedUser.Username != "testuser" {
            t.Errorf("Expected username 'testuser', got '%s'", retrievedUser.Username)
        }
    })
    
    t.Run("get user by email", func(t *testing.T) {
        retrievedUser, err := userService.GetByEmail(ctx, "test@example.com")
        if err != nil {
            t.Fatalf("Failed to get user by email: %v", err)
        }
        
        if retrievedUser.Email != "test@example.com" {
            t.Errorf("Expected email 'test@example.com', got '%s'", retrievedUser.Email)
        }
    })
    
    t.Run("get user by id", func(t *testing.T) {
        retrievedUser, err := userService.GetByID(ctx, 1)
        if err != nil {
            t.Fatalf("Failed to get user by ID: %v", err)
        }
        
        if retrievedUser.ID != 1 {
            t.Errorf("Expected user ID 1, got %d", retrievedUser.ID)
        }
    })
    
    t.Run("validate credentials", func(t *testing.T) {
        user, err := userService.ValidateCredentials(ctx, "testuser", "password123")
        if err != nil {
            t.Fatalf("Failed to validate credentials: %v", err)
        }
        
        if user.Username != "testuser" {
            t.Errorf("Expected username 'testuser', got '%s'", user.Username)
        }
    })
    
    t.Run("invalid credentials", func(t *testing.T) {
        _, err := userService.ValidateCredentials(ctx, "testuser", "wrongpassword")
        if err == nil {
            t.Error("Expected error for invalid credentials")
        }
    })
}

// Test health check endpoint
func TestHealthCheck(t *testing.T) {
    mockRepo := NewMockUserRepository()
    userService := service.NewUserService(mockRepo)
    userHandler := handler.NewUserHandler(userService)
    
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    w := httptest.NewRecorder()
    
    userHandler.HealthCheck(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
    }
    
    var response map[string]string
    if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
        t.Fatalf("Failed to decode health check response: %v", err)
    }
    
    if response["status"] != "healthy" {
        t.Errorf("Expected status 'healthy', got '%s'", response["status"])
    }
    
    if response["service"] != "gofr-account-service" {
        t.Errorf("Expected service 'gofr-account-service', got '%s'", response["service"])
    }
}
