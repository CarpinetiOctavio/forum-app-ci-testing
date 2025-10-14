package main

import (
	"log"
	"net/http"

	"tp06-testing/internal/database"
	"tp06-testing/internal/handlers"
	"tp06-testing/internal/repository"
	"tp06-testing/internal/router"
	"tp06-testing/internal/services"
)

func main() {
	// Inicializar base de datos
	db, err := database.InitDB("./database.db")
	if err != nil {
		log.Fatal("Error al inicializar la base de datos:", err)
	}
	defer db.Close()

	// Crear repositorios
	userRepo := repository.NewSQLiteUserRepository(db)
	postRepo := repository.NewSQLitePostRepository(db)

	// Crear servicios
	authService := services.NewAuthService(userRepo)
	postService := services.NewPostService(postRepo, userRepo)

	// Crear handlers
	authHandler := handlers.NewAuthHandler(authService)
	postHandler := handlers.NewPostHandler(postService)

	// Configurar rutas
	r := router.Setup(authHandler, postHandler)

	// Iniciar servidor
	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
