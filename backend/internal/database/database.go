package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB inicializa la base de datos SQLite
func InitDB(filepath string) (*sql.DB, error) {
	// Abrir conexión a SQLite
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	// Verificar que la conexión funcione
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Crear las tablas
	if err = createTables(db); err != nil {
		return nil, err
	}

	log.Println("Base de datos inicializada correctamente")
	return db, nil
}

// createTables crea el schema de la base de datos
func createTables(db *sql.DB) error {
	schema := `
	-- Tabla de usuarios
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		username TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Tabla de posts
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- Tabla de comentarios
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- Índices para mejorar rendimiento
	CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
	CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
	CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
	`

	_, err := db.Exec(schema)
	return err
}
