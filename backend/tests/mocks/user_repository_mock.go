package mocks

import (
	"tp06-testing/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository es un mock del UserRepository
type MockUserRepository struct {
	mock.Mock
}

// Create simula la creación de un usuario
func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// FindByEmail simula la búsqueda por email
func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)

	// Si se configuró para devolver nil (usuario no encontrado)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}

// FindByID simula la búsqueda por ID
func (m *MockUserRepository) FindByID(id int) (*models.User, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}
