package mocks

import (
	"tp06-testing/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockPostRepository es un mock del PostRepository
type MockPostRepository struct {
	mock.Mock
}

// Create simula la creación de un post
func (m *MockPostRepository) Create(post *models.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

// FindAll simula obtener todos los posts
func (m *MockPostRepository) FindAll() ([]*models.Post, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Post), args.Error(1)
}

// FindByID simula buscar un post por ID
func (m *MockPostRepository) FindByID(id int) (*models.Post, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Post), args.Error(1)
}

// Delete simula eliminar un post
func (m *MockPostRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// CreateComment simula crear un comentario
func (m *MockPostRepository) CreateComment(comment *models.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}

// FindCommentsByPostID simula obtener comentarios de un post
func (m *MockPostRepository) FindCommentsByPostID(postID int) ([]*models.Comment, error) {
	args := m.Called(postID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Comment), args.Error(1)
}

// DeleteComment simula eliminar un comentario
func (m *MockPostRepository) DeleteComment(postID int, commentID int, userID int) error {
	args := m.Called(postID, commentID, userID)
	return args.Error(0)
}
