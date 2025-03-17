package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vnlab/makeshop-payment/src/domain/entities"
)

// JWTService provides JWT token generation and validation
type JWTService struct {
	secretKey     string
	tokenDuration time.Duration
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID   string        `json:"user_id"`
	Username string        `json:"username"`
	Role     entities.Role `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWTService
func NewJWTService() *JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_jwt_secret_key_change_in_production"
	}

	// Default token duration is 24 hours
	tokenDuration := 24 * time.Hour

	return &JWTService{
		secretKey:     secret,
		tokenDuration: tokenDuration,
	}
}

// GenerateToken generates a new JWT token for a user
func (s *JWTService) GenerateToken(user *entities.User) (string, error) {
	if user == nil {
		return "", errors.New("user is nil")
	}

	claims := TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken validates the provided token string and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ExtractUserIDFromToken extracts user ID from a token string
func (s *JWTService) ExtractUserIDFromToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
