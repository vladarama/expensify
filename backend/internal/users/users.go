package users

import (
    "database/sql"
    "errors"
    "time"

    "expense-tracker/pkg/utils"
)

type User struct {
    ID           int64     `json:"id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

func RegisterUser(db *sql.DB, name, email, password string) (*User, error) {
    // Check if the email is already in use
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("email already registered")
    }

    passwordHash, err := utils.HashPassword(password)
    if err != nil {
        return nil, err
    }

    user := &User{
        Name:         name,
        Email:        email,
        PasswordHash: passwordHash,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    err = db.QueryRow(
        "INSERT INTO users (name, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
        user.Name, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
    ).Scan(&user.ID)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func AuthenticateUser(db *sql.DB, email, password string) (*User, error) {
    user := &User{}
    err := db.QueryRow(
        "SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email=$1",
        email,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("invalid email or password")
        }
        return nil, err
    }

    if !utils.CheckPasswordHash(password, user.PasswordHash) {
        return nil, errors.New("invalid email or password")
    }

    return user, nil
}
