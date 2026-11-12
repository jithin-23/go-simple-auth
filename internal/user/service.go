package user

import (
	"context"
	"strings"

	"simple-auth/internal/apperror"
	"simple-auth/internal/password"
)

// Service holds the business rules for users. It knows nothing about HTTP or GORM.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	email := normalizeEmail(req.Email)
	if err := s.ensureEmailFree(ctx, email); err != nil {
		return nil, err
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	role := req.Role
	if role == "" {
		role = RoleUser
	}

	u := &User{Name: strings.TrimSpace(req.Name), Email: email, PasswordHash: hash, Role: role}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) List(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id uint) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

// Update changes only the fields that were sent (non-nil pointers).
func (s *Service) Update(ctx context.Context, id uint, req UpdateUserRequest) (*User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		u.Name = strings.TrimSpace(*req.Name)
	}
	if req.Email != nil {
		email := normalizeEmail(*req.Email)
		if email != u.Email {
			if err := s.ensureEmailFree(ctx, email); err != nil {
				return nil, err
			}
			u.Email = email
		}
	}
	if req.Password != nil {
		hash, err := password.Hash(*req.Password)
		if err != nil {
			return nil, apperror.Internal(err)
		}
		u.PasswordHash = hash
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// ensureEmailFree returns a Conflict error if another user already has this email.
// (The unique index is the final safety net for two requests arriving at once.)
func (s *Service) ensureEmailFree(ctx context.Context, email string) error {
	_, err := s.repo.FindByEmail(ctx, email)
	switch {
	case err == nil:
		return apperror.Conflict("email already in use")
	case apperror.Is(err, apperror.CodeNotFound):
		return nil
	default:
		return err
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
