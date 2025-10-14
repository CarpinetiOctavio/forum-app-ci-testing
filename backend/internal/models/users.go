package models

import "time"

// User representa un usuario del sistema
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // No se serializa en JSON (por seguridad)
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// Credentials se usa para login
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest se usa para registro
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}
