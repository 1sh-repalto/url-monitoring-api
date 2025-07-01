package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/1sh-repalto/url-monitoring-api/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret = os.Getenv("access_secret")
	refreshSecret = os.Getenv("refresh_secret")
)

type contextKey string

const userIDKey contextKey = "userID"

func GenerateAccessToken(userID int) (string, error) {
	claims := model.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func GenerateRefreshToken(userID int) (string, error) {
	claims := model.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func ValidateAccessToken(accessToken string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func ValidateRefreshToken(refresToken string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(refresToken, &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*model.JWTClaims)

	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Missing access token", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateAccessToken(strings.TrimSpace(cookie.Value))
		if err != nil {
			http.Error(w, "Invalid access token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}