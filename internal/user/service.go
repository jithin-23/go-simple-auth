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

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.FindByEmail(ctx, normalizeEmail(email))
}

// Update changes only the fields that were sent (non-nil pointers).
// Regular users can only update themselves; admins can update anyone.
func (s *Service) Update(ctx context.Context, actor Actor, id uint, req UpdateUserRequest) (*User, error) {
	if !actor.IsAdmin() && actor.ID != id {
		return nil, apperror.Forbidden("you can only update your own profile")
	}

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

// ChangeRole promotes or demotes a user (admin only). The last admin can't be demoted.
func (s *Service) ChangeRole(ctx context.Context, actor Actor, id uint, role Role) (*User, error) {
	if !actor.IsAdmin() {
		return nil, apperror.Forbidden("only admins can change roles")
	}

	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u.Role == role {
		return u, nil // nothing to change
	}

	if u.Role == RoleAdmin {
		if err := s.ensureNotLastAdmin(ctx, "demote"); err != nil {
			return nil, err
		}
	}

	u.Role = role
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// Delete soft-deletes a user. The last admin can't be deleted.
func (s *Service) Delete(ctx context.Context, id uint) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if u.Role == RoleAdmin {
		if err := s.ensureNotLastAdmin(ctx, "delete"); err != nil {
			return err
		}
	}

	return s.repo.Delete(ctx, id)
}

// ensureNotLastAdmin stops an action that would leave the system with no admin.
func (s *Service) ensureNotLastAdmin(ctx context.Context, action string) error {
	admins, err := s.repo.CountByRole(ctx, RoleAdmin)
	if err != nil {
		return err
	}
	if admins <= 1 {
		return apperror.Conflict("cannot " + action + " the last admin")
	}
	return nil
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
