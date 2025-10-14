package repository

import (
	"database/sql"

	"tp06-testing/internal/models"
)

// UserRepository define las operaciones sobre usuarios
// INTERFACE: permite crear mocks fácilmente para testing
type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id int) (*models.User, error)
}

// SQLiteUserRepository implementa UserRepository usando SQLite
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository crea una nueva instancia
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// Create inserta un nuevo usuario en la base de datos
func (r *SQLiteUserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password, username, created_at)
		VALUES (?, ?, ?, datetime('now'))
	`
	result, err := r.db.Exec(query, user.Email, user.Password, user.Username)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

// FindByEmail busca un usuario por email
func (r *SQLiteUserRepository) FindByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, username, created_at FROM users WHERE email = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Usuario no encontrado (no es error)
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID busca un usuario por ID
func (r *SQLiteUserRepository) FindByID(id int) (*models.User, error) {
	query := `SELECT id, email, password, username, created_at FROM users WHERE id = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
