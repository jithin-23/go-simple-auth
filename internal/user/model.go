package user

import "gorm.io/gorm"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User maps to the "users" table (created by migrations/000001_create_users).
// gorm.Model adds ID, CreatedAt, UpdatedAt and DeletedAt (soft delete),
// the same as AbstractEntity in the NestJS app.
type User struct {
	gorm.Model
	Name         string
	Email        string
	PasswordHash string `json:"-"`            // never serialized to JSON
	Role         Role   `gorm:"default:user"` // let the DB default apply when empty
}
