package handlers

import (
	"encoding/json"
	"net/http"

	"tp06-testing/internal/models"
	"tp06-testing/internal/services"
)

// AuthHandler maneja las peticiones HTTP de autenticación
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler crea una nueva instancia
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register maneja POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Decodificar el body JSON
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Llamar al servicio
	user, err := h.authService.Register(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Responder con el usuario creado
	respondWithJSON(w, http.StatusCreated, user)
}

// Login maneja POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Decodificar el body JSON
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Llamar al servicio
	user, err := h.authService.Login(&creds)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Responder con el usuario autenticado
	respondWithJSON(w, http.StatusOK, user)
}

// Funciones auxiliares para responder JSON

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
