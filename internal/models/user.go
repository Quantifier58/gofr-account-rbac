package models

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
    ID           int       `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}

func CreateUser(ctx context.Context, db *pgxpool.Pool, user *User) error {
    query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at`
    row := db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash)
    return row.Scan(&user.ID, &user.CreatedAt)
}

func GetUserByID(ctx context.Context, db *pgxpool.Pool, id int) (*User, error) {
    query := `SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1`
    row := db.QueryRow(ctx, query, id)
    user := &User{}
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return user, nil
}
