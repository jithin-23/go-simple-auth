package user

import (
	"context"

	"gorm.io/gorm"
)

// Repository is everything the service needs from storage.
// The service depends on this interface, not on GORM, so tests can use a fake.
type Repository interface {
	Create(ctx context.Context, u *User) error
	FindAll(ctx context.Context) ([]User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uint) error
}

// GormRepository implements Repository using GORM + Postgres.
type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *GormRepository) FindAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.WithContext(ctx).Order("id").Find(&users).Error
	return users, err
}

func (r *GormRepository) FindByID(ctx context.Context, id uint) (*User, error) {
	var u User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err // gorm.ErrRecordNotFound if missing
	}
	return &u, nil
}

func (r *GormRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Update saves all fields of an existing user (load it first, change fields, then call Update).
func (r *GormRepository) Update(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *GormRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&User{}, id) // soft delete: sets deleted_at
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Compile-time check that GormRepository satisfies Repository.
var _ Repository = (*GormRepository)(nil)
