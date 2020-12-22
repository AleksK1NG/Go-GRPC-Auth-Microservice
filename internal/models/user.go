package models

import "github.com/google/uuid"

// User base model
type User struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar"`
	Password  string    `json:"password,omitempty"`
}

// Sanitize password
func (u *User) SanitizePassword() {
	u.Password = ""
}
