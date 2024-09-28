package auth

import (
    "context"
    "errors"
    "net/http"
    "strings"

    "expense-tracker/pkg/utils"
)

type key int

const (
    UserIDKey key = iota
)

func JWTMiddleware(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header missing", http.StatusUnauthorized)
                return
            }

            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
                return
            }

            tokenString := parts[1]
            userID, err := utils.ParseJWT(tokenString, jwtSecret)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), UserIDKey, userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GenerateJWT(userID int64, jwtSecret string) (string, error) {
    return utils.GenerateJWT(userID, jwtSecret)
}

func ParseJWT(tokenStr, jwtSecret string) (int64, error) {
    return utils.ParseJWT(tokenStr, jwtSecret)
}

var ErrInvalidToken = errors.New("invalid token")
