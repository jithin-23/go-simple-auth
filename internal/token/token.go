// Package token creates and verifies JWT access tokens (≈ JwtService in NestJS).
package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const issuer = "simple-auth"

// Claims is the data stored inside the token (≈ JwtPayload in NestJS).
type Claims struct {
	UserID uint   `json:"uid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Manager signs and parses tokens with one secret and lifetime.
type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{secret: []byte(secret), ttl: ttl}
}

// TTL is how long a newly generated token stays valid.
func (m *Manager) TTL() time.Duration { return m.ttl }

// Generate returns a signed token for the given user.
func (m *Manager) Generate(userID uint, email, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

// Parse verifies the signature, algorithm, issuer and expiry, then returns the claims.
func (m *Manager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), // reject "alg: none" etc.
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	if claims.UserID == 0 {
		return nil, errors.New("parse token: missing user id")
	}
	return claims, nil
}
