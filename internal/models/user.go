package models

// User base model
type User struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password,omitempty"`
}

// Sanitize password
func (u *User) SanitizePassword() {
	u.Password = ""
}
