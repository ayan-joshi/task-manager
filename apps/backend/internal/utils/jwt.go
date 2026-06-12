package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"taskmanager/internal/domain"
)

// ErrInvalidToken is returned when a token fails parsing or validation.
var ErrInvalidToken = errors.New("invalid or expired token")

// Claims is the JWT payload. It embeds the registered claims (exp, iat, sub)
// and adds the user role so authorization can be performed without an extra DB
// lookup on every request.
type Claims struct {
	UserID uuid.UUID   `json:"uid"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager issues and verifies signed JSON Web Tokens.
type JWTManager struct {
	secret      []byte
	expiryHours int
}

// NewJWTManager constructs a JWTManager.
func NewJWTManager(secret string, expiryHours int) *JWTManager {
	return &JWTManager{secret: []byte(secret), expiryHours: expiryHours}
}

// Generate signs a token for the given user.
func (m *JWTManager) Generate(userID uuid.UUID, role domain.Role) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.expiryHours) * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Verify parses and validates a token string, returning its claims. It enforces
// the HMAC signing method to prevent algorithm-confusion attacks.
func (m *JWTManager) Verify(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
