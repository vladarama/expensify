package utils

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int64, jwtSecret string) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}

func ParseJWT(tokenStr, jwtSecret string) (int64, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtSecret), nil
    })
    if err != nil {
        return 0, err
    }
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userIDFloat, ok := claims["sub"].(float64)
        if !ok {
            return 0, errors.New("invalid token claims")
        }
        return int64(userIDFloat), nil
    }
    return 0, errors.New("invalid token")
}
