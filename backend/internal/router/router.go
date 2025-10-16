package router

import (
	"net/http"

	"tp06-testing/internal/handlers"

	"github.com/gorilla/mux"
)

// Setup configura todas las rutas de la aplicación
func Setup(authHandler *handlers.AuthHandler, postHandler *handlers.PostHandler) *mux.Router {
	router := mux.NewRouter()

	// Middleware CORS
	router.Use(corsMiddleware)

	// Rutas de autenticación
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Rutas de posts
	router.HandleFunc("/api/posts", postHandler.GetAllPosts).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/posts", postHandler.CreatePost).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/posts/{id}", postHandler.GetPostByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/posts/{id}", postHandler.DeletePost).Methods("DELETE", "OPTIONS")

	// Rutas de comentarios
	router.HandleFunc("/api/posts/{id}/comments", postHandler.GetComments).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/posts/{id}/comments", postHandler.CreateComment).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/posts/{postId}/comments/{commentId}", postHandler.DeleteComment).Methods("DELETE", "OPTIONS")

	return router
}

// corsMiddleware permite peticiones desde el frontend
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar headers CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")

		// Si es una petición OPTIONS (preflight), responder inmediatamente
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
