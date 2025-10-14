package services

import (
	"testing"

	"tp06-testing/internal/models"
	"tp06-testing/internal/services"
	"tp06-testing/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestRegister_Success prueba el registro exitoso de un usuario
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

	// ACT: Ejecutar la función que estamos probando
	user, err := authService.Register(req)

	// ASSERT: Verificar los resultados
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "testuser", user.Username)

	// Verificar que se llamaron los métodos del mock
	mockRepo.AssertExpectations(t)
}

// TestRegister_EmailVacio prueba que falle con email vacío
func TestRegister_EmailVacio(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "", // Email vacío
		Password: "123456",
		Username: "testuser",
	}

	// ACT
	user, err := authService.Register(req)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "el email es requerido", err.Error())

	// NO debe haber llamado a la BD porque falló la validación antes
	mockRepo.AssertNotCalled(t, "FindByEmail")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestRegister_EmailInvalido prueba que falle con email sin @
func TestRegister_EmailInvalido(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "invalidemail", // Sin @
		Password: "123456",
		Username: "testuser",
	}

	// ACT
	user, err := authService.Register(req)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "el email debe ser válido", err.Error())
}

// TestRegister_PasswordCorto prueba que falle con password menor a 6 caracteres
func TestRegister_PasswordCorto(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "123", // Muy corto
		Username: "testuser",
	}

	// ACT
	user, err := authService.Register(req)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "la contraseña debe tener al menos 6 caracteres", err.Error())
}

// TestRegister_UsernameVacio prueba que falle con username vacío
func TestRegister_UsernameVacio(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "123456",
		Username: "", // Username vacío
	}

	// ACT
	user, err := authService.Register(req)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "el nombre de usuario es requerido", err.Error())
}

// TestRegister_EmailDuplicado prueba que falle si el email ya existe
func TestRegister_EmailDuplicado(t *testing.T) {
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
	assert.Equal(t, "el email ya está registrado", err.Error())

	// NO debe llamar a Create porque el email ya existe
	mockRepo.AssertNotCalled(t, "Create")
}

// TestLogin_Success prueba el login exitoso
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

	// Configurar el mock: el usuario existe
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

// TestLogin_EmailVacio prueba que falle con email vacío
func TestLogin_EmailVacio(t *testing.T) {
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
	assert.Equal(t, "el email es requerido", err.Error())

	mockRepo.AssertNotCalled(t, "FindByEmail")
}

// TestLogin_PasswordVacio prueba que falle con password vacío
func TestLogin_PasswordVacio(t *testing.T) {
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
	assert.Equal(t, "la contraseña es requerida", err.Error())
}

// TestLogin_UsuarioNoExiste prueba que falle si el usuario no existe
func TestLogin_UsuarioNoExiste(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	// Configurar el mock: el usuario NO existe
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
	assert.Equal(t, "credenciales inválidas", err.Error())

	mockRepo.AssertExpectations(t)
}

// TestLogin_PasswordIncorrecta prueba que falle con password incorrecta
func TestLogin_PasswordIncorrecta(t *testing.T) {
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
		Password: "wrongpassword", // Password incorrecta
	}

	// ACT
	user, err := authService.Login(creds)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "credenciales inválidas", err.Error())

	mockRepo.AssertExpectations(t)
}
