package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTService(secret string, ttl time.Duration) *JWTService {
	return &JWTService{secret: []byte(secret), ttl: ttl}
}

func (j *JWTService) Issue(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// Verify validates signature and expiry, returning the user id from the
// subject claim. Any failure (malformed, bad signature, expired) is
// collapsed into a single error so callers never leak which check failed.
func (j *JWTService) Verify(tokenString string) (int64, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid or expired token")
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid token subject")
	}
	return userID, nil
}

type contextKey int

const userIDContextKey contextKey = 0

func ContextWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(userIDContextKey).(int64)
	return v, ok
}
