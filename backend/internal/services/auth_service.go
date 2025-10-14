package services

import (
	"errors"
	"strings"

	"tp06-testing/internal/models"
	"tp06-testing/internal/repository"
)

// AuthService maneja la lógica de autenticación
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService crea una nueva instancia
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register registra un nuevo usuario
// Aquí validamos las reglas de negocio
func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Validación 1: Email no puede estar vacío
	if strings.TrimSpace(req.Email) == "" {
		return nil, errors.New("el email es requerido")
	}

	// Validación 2: Email debe contener @
	if !strings.Contains(req.Email, "@") {
		return nil, errors.New("el email debe ser válido")
	}

	// Validación 3: Password debe tener al menos 6 caracteres
	if len(req.Password) < 6 {
		return nil, errors.New("la contraseña debe tener al menos 6 caracteres")
	}

	// Validación 4: Username no puede estar vacío
	if strings.TrimSpace(req.Username) == "" {
		return nil, errors.New("el nombre de usuario es requerido")
	}

	// Validación 5: Verificar que el email no esté registrado
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("el email ya está registrado")
	}

	// Crear el usuario
	user := &models.User{
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: req.Password, // En producción: hashear con bcrypt
		Username: strings.TrimSpace(req.Username),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login autentica un usuario
func (s *AuthService) Login(creds *models.Credentials) (*models.User, error) {
	// Validación 1: Email no puede estar vacío
	if strings.TrimSpace(creds.Email) == "" {
		return nil, errors.New("el email es requerido")
	}

	// Validación 2: Password no puede estar vacío
	if creds.Password == "" {
		return nil, errors.New("la contraseña es requerida")
	}

	// Buscar usuario por email
	user, err := s.userRepo.FindByEmail(strings.ToLower(strings.TrimSpace(creds.Email)))
	if err != nil {
		return nil, err
	}

	// Validación 3: Usuario debe existir
	if user == nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Validación 4: Password debe coincidir
	// En producción: usar bcrypt.CompareHashAndPassword
	if user.Password != creds.Password {
		return nil, errors.New("credenciales inválidas")
	}

	return user, nil
}
