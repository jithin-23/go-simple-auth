package user

import "time"

// CreateUserRequest is the JSON body for creating a user.
// `binding` tags are validated by Gin (like class-validator decorators).
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"` // bcrypt only uses the first 72 bytes
	Role     Role   `json:"role" binding:"omitempty,oneof=user admin"`
}

// UpdateUserRequest is the JSON body for PATCH /users/:id.
// Pointer fields: nil means "not sent, don't change it".
type UpdateUserRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=2,max=100"`
	Email    *string `json:"email" binding:"omitempty,email,max=255"`
	Password *string `json:"password" binding:"omitempty,min=8,max=72"`
}

// UserResponse is what the API returns. It never contains the password hash.
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToResponse converts a User model into the API response shape.
func ToResponse(u *User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToResponseList converts a slice of users.
func ToResponseList(users []User) []UserResponse {
	out := make([]UserResponse, 0, len(users))
	for i := range users {
		out = append(out, ToResponse(&users[i]))
	}
	return out
}

// UpdateRoleRequest is the JSON body for PATCH /users/:id/role (admin only).
type UpdateRoleRequest struct {
	Role Role `json:"role" binding:"required,oneof=user admin"`
}
