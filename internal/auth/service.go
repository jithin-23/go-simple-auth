package auth

import (
	"context"

	"simple-auth/internal/apperror"
	"simple-auth/internal/password"
	"simple-auth/internal/token"
	"simple-auth/internal/user"
)

// Same message for "no such email" and "wrong password", so attackers can't
// find out which emails are registered.
const invalidCredentials = "invalid email or password"

// dummyHash is compared against when the email doesn't exist, so a failed
// login takes the same time either way (another hint attackers could use).
var dummyHash, _ = password.Hash("dummy-password-for-timing")

// Service handles signup, login and "who am I" (≈ AuthService in NestJS).
type Service struct {
	users  *user.Service
	tokens *token.Manager
}

func NewService(users *user.Service, tokens *token.Manager) *Service {
	return &Service{users: users, tokens: tokens}
}

// Signup creates a regular user. The role is always "user": nobody can sign up as admin.
func (s *Service) Signup(ctx context.Context, req SignupRequest) (*user.User, error) {
	return s.users.Create(ctx, user.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     user.RoleUser,
	})
}

// Login checks the credentials and returns a signed access token.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	u, err := s.users.GetByEmail(ctx, req.Email)
	if apperror.Is(err, apperror.CodeNotFound) {
		password.Matches(dummyHash, req.Password) // keep timing the same as a real check
		return nil, apperror.Unauthorized(invalidCredentials)
	}
	if err != nil {
		return nil, err
	}

	if !password.Matches(u.PasswordHash, req.Password) {
		return nil, apperror.Unauthorized(invalidCredentials)
	}

	accessToken, err := s.tokens.Generate(u.ID, u.Email, string(u.Role))
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return &TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.tokens.TTL().Seconds()),
	}, nil
}

// Me returns the current user's latest data from the database.
func (s *Service) Me(ctx context.Context, userID uint) (*user.User, error) {
	u, err := s.users.GetByID(ctx, userID)
	if apperror.Is(err, apperror.CodeNotFound) {
		// Token is valid but the user was deleted after it was issued
		return nil, apperror.Unauthorized("user no longer exists")
	}
	return u, err
}
