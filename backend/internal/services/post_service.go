package services

import (
	"errors"
	"strings"

	"tp06-testing/internal/models"
	"tp06-testing/internal/repository"
)

// PostService maneja la lógica de posts y comentarios
type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

// NewPostService crea una nueva instancia
func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// CreatePost crea un nuevo post
func (s *PostService) CreatePost(req *models.CreatePostRequest, userID int) (*models.Post, error) {
	// Validación 1: Título no puede estar vacío
	if strings.TrimSpace(req.Title) == "" {
		return nil, errors.New("el título es requerido")
	}

	// Validación 2: Título debe tener al menos 3 caracteres
	if len(strings.TrimSpace(req.Title)) < 3 {
		return nil, errors.New("el título debe tener al menos 3 caracteres")
	}

	// Validación 3: Contenido no puede estar vacío
	if strings.TrimSpace(req.Content) == "" {
		return nil, errors.New("el contenido es requerido")
	}

	// Validación 4: Usuario debe existir
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Crear el post
	post := &models.Post{
		Title:   strings.TrimSpace(req.Title),
		Content: strings.TrimSpace(req.Content),
		UserID:  userID,
	}

	err = s.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	// Agregar el username para la respuesta
	post.Username = user.Username

	return post, nil
}

// GetAllPosts obtiene todos los posts
func (s *PostService) GetAllPosts() ([]*models.Post, error) {
	posts, err := s.postRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Si no hay posts, devolver lista vacía (no es error)
	if posts == nil {
		return []*models.Post{}, nil
	}

	return posts, nil
}

// GetPostByID obtiene un post específico
func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	// Validación: ID debe ser positivo
	if id <= 0 {
		return nil, errors.New("id inválido")
	}

	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, errors.New("post no encontrado")
	}

	return post, nil
}

// DeletePost elimina un post (solo el autor puede hacerlo)
func (s *PostService) DeletePost(postID int, userID int) error {
	// Validación 1: Post debe existir
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post no encontrado")
	}

	// Validación 2: Solo el autor puede eliminar
	if post.UserID != userID {
		return errors.New("no tienes permiso para eliminar este post")
	}

	// Eliminar
	return s.postRepo.Delete(postID)
}

// CreateComment agrega un comentario a un post
func (s *PostService) CreateComment(postID int, req *models.CreateCommentRequest, userID int) (*models.Comment, error) {
	// Validación 1: Contenido no puede estar vacío
	if strings.TrimSpace(req.Content) == "" {
		return nil, errors.New("el contenido del comentario es requerido")
	}

	// Validación 2: Post debe existir
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post no encontrado")
	}

	// Validación 3: Usuario debe existir
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Crear comentario
	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: strings.TrimSpace(req.Content),
	}

	err = s.postRepo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	// Agregar username para la respuesta
	comment.Username = user.Username

	return comment, nil
}

// GetCommentsByPostID obtiene todos los comentarios de un post
func (s *PostService) GetCommentsByPostID(postID int) ([]*models.Comment, error) {
	// Validación: Post debe existir
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post no encontrado")
	}

	comments, err := s.postRepo.FindCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	// Si no hay comentarios, devolver lista vacía
	if comments == nil {
		return []*models.Comment{}, nil
	}

	return comments, nil
}

func (s *PostService) DeleteComment(postID int, commentID int, userID int) error {
	// Validación: Post debe existir
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post no encontrado")
	}

	// Validación: Usuario debe existir
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("usuario no encontrado")
	}

	// Eliminar comentario (solo el autor puede)
	return s.postRepo.DeleteComment(postID, commentID, userID)
}
