package services

import (
	"testing"

	"forum-app-ci-testing/internal/models"
	"forum-app-ci-testing/internal/services"
	"forum-app-ci-testing/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestRegister_Success verifies a user registers successfully
func TestRegister_Success(t *testing.T) {
	// ARRANGE: Preparar el mock y datos de prueba
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	// Configurar el mock: el email NO existe (devuelve nil)
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)

	// Configurar el mock: Create debe ejecutarse correctamente
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "123456",
		Username: "testuser",
	}

	// ACT: Execute the function under test
	user, err := authService.Register(req)

	// ASSERT: Verificar los resultados
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "testuser", user.Username)

	// Verify the mock's methods were called
	mockRepo.AssertExpectations(t)
}

// TestRegister_EmptyEmail verifies registration fails when the email is empty
func TestRegister_EmptyEmail(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "", // Empty email
		Password: "123456",
		Username: "testuser",
	}

	// ACT
	user, err := authService.Register(req)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email is required", err.Error())

	// Should NOT have called the DB because validation failed first
	mockRepo.AssertNotCalled(t, "FindByEmail")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestRegister_InvalidEmail verifies registration fails when the email is missing the @ symbol
func TestRegister_InvalidEmail(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "invalidemail", // Missing @
		Password: "123456",
		Username: "testuser",
	}

	// ACT
	user, err := authService.Register(req)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email must be valid", err.Error())
}

// TestRegister_PasswordTooShort verifies registration fails when the password is under 6 characters
func TestRegister_PasswordTooShort(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "123", // Too short
		Username: "testuser",
	}

	// ACT
	user, err := authService.Register(req)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "password must be at least 6 characters", err.Error())
}

// TestRegister_EmptyUsername verifies registration fails when the username is empty
func TestRegister_EmptyUsername(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "123456",
		Username: "", // Empty username
	}

	// ACT
	user, err := authService.Register(req)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "username is required", err.Error())
}

// TestRegister_DuplicateEmail verifies registration fails when the email is already registered
func TestRegister_DuplicateEmail(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	existingUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "existinguser",
	}

	// Configurar el mock: el email YA existe
	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "123456",
		Username: "testuser",
	}

	// ACT
	user, err := authService.Register(req)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email is already registered", err.Error())

	// NO debe llamar a Create porque el email ya existe
	mockRepo.AssertNotCalled(t, "Create")
}

// TestLogin_Success verifies a user logs in successfully
func TestLogin_Success(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	existingUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "123456",
		Username: "testuser",
	}

	// Configure the mock: the user exists
	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	creds := &models.Credentials{
		Email:    "test@example.com",
		Password: "123456",
	}

	// ACT
	user, err := authService.Login(creds)

	// ASSERT
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "testuser", user.Username)

	mockRepo.AssertExpectations(t)
}

// TestLogin_EmptyEmail verifies login fails when the email is empty
func TestLogin_EmptyEmail(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	creds := &models.Credentials{
		Email:    "",
		Password: "123456",
	}

	// ACT
	user, err := authService.Login(creds)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email is required", err.Error())

	mockRepo.AssertNotCalled(t, "FindByEmail")
}

// TestLogin_EmptyPassword verifies login fails when the password is empty
func TestLogin_EmptyPassword(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	creds := &models.Credentials{
		Email:    "test@example.com",
		Password: "",
	}

	// ACT
	user, err := authService.Login(creds)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "password is required", err.Error())
}

// TestLogin_UserNotFound verifies login fails when no user matches the given email
func TestLogin_UserNotFound(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	// Configure the mock: the user does NOT exist
	mockRepo.On("FindByEmail", "noexiste@example.com").Return(nil, nil)

	creds := &models.Credentials{
		Email:    "noexiste@example.com",
		Password: "123456",
	}

	// ACT
	user, err := authService.Login(creds)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid credentials", err.Error())

	mockRepo.AssertExpectations(t)
}

// TestLogin_IncorrectPassword verifies login fails when the password does not match
func TestLogin_IncorrectPassword(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	existingUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "123456",
		Username: "testuser",
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	creds := &models.Credentials{
		Email:    "test@example.com",
		Password: "wrongpassword", // Incorrect password
	}

	// ACT
	user, err := authService.Login(creds)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid credentials", err.Error())

	mockRepo.AssertExpectations(t)
}
