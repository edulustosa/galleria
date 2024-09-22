package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/edulustosa/galleria/internal/api/handlers"
	"github.com/golang-jwt/jwt/v5"
)

func getTokenFromAuthorizationHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return "", errors.New("invalid token format")
	}

	return tokenString, nil
}

func verifyClaims(token *jwt.Token) (string, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("missing sub claim")
		}

		return userID, nil
	}

	return "", errors.New("invalid claims")
}

type ContextKey string

const UserIDKey ContextKey = "userID"

func JWTAuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := getTokenFromAuthorizationHeader(r)
			if err != nil {
				handlers.HandleError(
					w,
					http.StatusUnauthorized,
					handlers.Error{Message: err.Error()},
				)
				return
			}

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("invalid signing method: %v", t.Header["alg"])
				}
				return jwtKey, nil
			})
			if err != nil {
				handlers.HandleError(
					w,
					http.StatusUnauthorized,
					handlers.Error{Message: "invalid token", Details: err.Error()},
				)
				return
			}

			userID, err := verifyClaims(token)
			if err != nil {
				handlers.HandleError(
					w,
					http.StatusUnauthorized,
					handlers.Error{Message: "invalid token", Details: err.Error()},
				)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
