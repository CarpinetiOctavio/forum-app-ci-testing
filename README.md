# TP06 - Pruebas Unitarias

**Materia:** Ingeniería de Software 3  
**Alumno:** Octavio Carpineti - Kevin Massholder  
**Año:** 2025

Mini red social con suite completa de pruebas unitarias (42 tests), mocking de dependencias externas y CI/CD automático.

---

## 📋 Tabla de Contenidos

- [Tecnologías](#tecnologías)
- [Arquitectura](#arquitectura)
- [Funcionalidades](#funcionalidades)
- [Prerequisitos](#prerequisitos)
- [Instalación](#instalación)
- [Ejecución](#ejecución)
- [Testing](#testing)
- [CI/CD](#cicd)
- [Estructura del Proyecto](#estructura-del-proyecto)

---

## 🛠️ Tecnologías

### Backend
- **Go 1.21+**
- **SQLite** (base de datos)
- **Gorilla Mux** (routing)
- **testify** (testing + mocking)

### Frontend
- **React 18** con **TypeScript**
- **Axios** (HTTP client)
- **Jest** + **React Testing Library** (testing)

### DevOps
- **GitHub Actions** (CI/CD)
- **Docker** (futuro TP08)

---

## 🏗️ Arquitectura

```
┌─────────────┐      HTTP      ┌──────────────┐
│   Frontend  │ ←────────────→ │   Backend    │
│   (React)   │                │     (Go)     │
└─────────────┘                └──────────────┘
                                       ↓
                               ┌──────────────┐
                               │   SQLite DB  │
                               └──────────────┘
```

### Capas del Backend

```
Handlers     ← Controladores HTTP (reciben requests)
    ↓
Services     ← Lógica de negocio (validaciones, reglas)
    ↓
Repository   ← Acceso a datos (SQL)
    ↓
Database     ← SQLite
```

### Testing Strategy

```
PRODUCCIÓN                      TESTING
─────────────────────────────────────────────────
Repository (SQLite)       →     Mock Repository
HTTP (axios)              →     Mock axios
                          
Resultado: Tests rápidos, aislados y reproducibles
```

---

## ✨ Funcionalidades

### Autenticación
- ✅ Registro de usuarios
- ✅ Login con email/password
- ✅ Validaciones (email, password, username)

### Posts
- ✅ Crear post (título + contenido)
- ✅ Listar todos los posts
- ✅ Ver detalle de un post
- ✅ Eliminar post (solo el autor)

### Comentarios
- ✅ Agregar comentario a un post
- ✅ Listar comentarios
- ✅ Eliminar comentario (solo el autor)

### Reglas de Negocio (testeadas)
- 🔒 Solo el autor puede eliminar su post
- 🔒 Solo el autor puede eliminar su comentario
- ✉️ Email debe ser válido y único
- 🔑 Password mínimo 6 caracteres
- 📝 Título de post mínimo 3 caracteres

---

## 📦 Prerequisitos

### Instalación de Herramientas

#### Go (Backend)
```bash
# Verificar instalación
go version  # Debe ser 1.21+

# Si no está instalado: https://go.dev/dl/
```

#### Node.js (Frontend)
```bash
# Verificar instalación
node --version  # Debe ser 18+
npm --version

# Si no está instalado: https://nodejs.org/
```

#### Git
```bash
git --version

# Si no está instalado: https://git-scm.com/
```

---

## 🚀 Instalación

### 1. Clonar el repositorio

```bash
git clone https://github.com/TU-USUARIO/tp06-testing.git
cd tp06-testing
```

### 2. Instalar dependencias del Backend

```bash
cd backend
go mod download
```

### 3. Instalar dependencias del Frontend

```bash
cd ../frontend
npm install
```

---

## ▶️ Ejecución

### Opción A: Ejecutar Backend y Frontend por separado

#### Terminal 1 - Backend
```bash
cd backend
go run cmd/api/main.go
```

Deberías ver:
```
Base de datos inicializada correctamente
🚀 Servidor corriendo en http://localhost:8080
```

#### Terminal 2 - Frontend
```bash
cd frontend
npm start
```

Se abrirá automáticamente en: `http://localhost:3000`

### Opción B: Script para ejecutar ambos (Linux/Mac)

```bash
# Crear script
cat > run.sh << 'EOF'
#!/bin/bash
cd backend && go run cmd/api/main.go &
BACKEND_PID=$!
cd ../frontend && npm start
kill $BACKEND_PID
EOF

chmod +x run.sh
./run.sh
```

---

## 🧪 Testing

### Backend Tests (Go)

```bash
cd backend

# Ejecutar todos los tests
go test ./tests/services/... -v

# Con cobertura
go test ./tests/services/... -v -cover

# Solo un test específico
go test ./tests/services/ -v -run TestRegister_Success
```

**Resultado esperado:**
```
=== RUN   TestRegister_Success
--- PASS: TestRegister_Success (0.00s)
...
PASS
ok      tp06-testing/tests/services     0.582s
```

**Total: 23 tests** ✅

### Frontend Tests (React)

```bash
cd frontend

# Ejecutar todos los tests
npm test

# Con cobertura
npm test -- --coverage

# Sin modo watch
npm test -- --watchAll=false
```

**Resultado esperado:**
```
PASS  src/components/Login/Login.test.tsx
PASS  src/components/PostList/PostList.test.tsx
PASS  src/components/CommentList/CommentList.test.tsx
PASS  src/services/authService.test.ts

Test Suites: 4 passed, 4 total
Tests:       19 passed, 19 total
```

**Total: 19 tests** ✅

### Ejecutar TODOS los tests (Backend + Frontend)

```bash
# Desde la raíz del proyecto
cd backend && go test ./... && cd ../frontend && npm test -- --watchAll=false
```

---

## 🔄 CI/CD

### GitHub Actions

El proyecto incluye un pipeline de CI/CD que se ejecuta automáticamente en cada push.

**Archivo:** `.github/workflows/ci.yml`

**Workflow:**
1. ✅ **Backend Tests** - Ejecuta `go test`
2. ✅ **Frontend Tests** - Ejecuta `npm test`
3. ✅ **Backend Build** - Compila con `go build`
4. ✅ **Frontend Build** - Compila con `npm run build`
5. ✅ **Summary** - Resumen final

**Ver resultados:**
1. Ir a: `https://github.com/TU-USUARIO/tp06-testing/actions`
2. Seleccionar el workflow más reciente
3. Ver logs detallados de cada job

---

## 📁 Estructura del Proyecto

```
tp06-testing/
├── .github/
│   └── workflows/
│       └── ci.yml                   # Pipeline CI/CD
│
├── backend/
│   ├── cmd/api/
│   │   └── main.go                  # Punto de entrada
│   ├── internal/
│   │   ├── database/
│   │   │   └── database.go          # Inicialización SQLite
│   │   ├── models/                  # Structs (User, Post, Comment)
│   │   ├── repository/              # Acceso a datos
│   │   │   ├── user_repository.go
│   │   │   └── post_repository.go
│   │   ├── services/                # Lógica de negocio
│   │   │   ├── auth_service.go
│   │   │   └── post_service.go
│   │   ├── handlers/                # Controladores HTTP
│   │   │   ├── auth_handler.go
│   │   │   └── post_handler.go
│   │   └── router/
│   │       └── router.go            # Rutas
│   ├── tests/
│   │   ├── mocks/                   # Repositorios mockeados
│   │   │   ├── user_repository_mock.go
│   │   │   └── post_repository_mock.go
│   │   └── services/                # Tests unitarios
│   │       ├── auth_service_test.go
│   │       └── post_service_test.go
│   ├── go.mod
│   └── database.db                  # SQLite (generado automáticamente)
│
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Login/
│   │   │   │   ├── Login.tsx
│   │   │   │   ├── Login.test.tsx
│   │   │   │   └── Login.css
│   │   │   ├── PostList/
│   │   │   │   ├── PostList.tsx
│   │   │   │   ├── PostList.test.tsx
│   │   │   │   └── PostList.css
│   │   │   ├── CreatePost/
│   │   │   ├── CommentList/
│   │   │   │   ├── CommentList.tsx
│   │   │   │   ├── CommentList.test.tsx
│   │   │   │   └── CommentList.css
│   │   │   ├── CommentForm/
│   │   │   └── PostDetail/
│   │   ├── services/
│   │   │   ├── authService.ts
│   │   │   ├── authService.test.ts
│   │   │   └── postService.ts
│   │   ├── __mocks__/
│   │   │   └── axios.ts             # Mock de HTTP
│   │   ├── types/
│   │   │   └── index.ts             # TypeScript types
│   │   ├── App.tsx
│   │   └── setupTests.ts
│   ├── package.json
│   └── tsconfig.json
│
├── README.md                        # Este archivo
└── decisiones.md                    # Documentación técnica
```

---

## 📊 Cobertura de Tests

### Backend (23 tests)

| Componente  | Tests |                 Descripción                     |
|-------------|-------|-------------------------------------------------|
| AuthService | 11    | Register (6), Login (5)                         |
| PostService | 12    | CreatePost (5), DeletePost (3), DeleteComment(4)|

### Frontend (19 tests)

| Componente | Tests |            Descripción             |
|------------|-------|------------------------------------|
| Login      | 5     | Renderizado, validaciones, estados |
| PostList   | 5     | Renderizado, eliminación, permisos |
| CommentList| 5     | Renderizado, eliminación, permisos |
| authService| 4     | Login/Register con mocks HTTP      |

**Total: 42 tests automatizados** ✅

---

## 🎯 Conceptos Implementados

### Testing
- ✅ **Pruebas Unitarias** (backend + frontend)
- ✅ **Patrón AAA** (Arrange, Act, Assert)
- ✅ **Mocking** (Repository + HTTP)
- ✅ **Aislamiento** de dependencias
- ✅ **Casos edge** y validaciones

### Arquitectura
- ✅ **Separación de concerns** (capas)
- ✅ **Dependency Injection** (interfaces)
- ✅ **Repository Pattern**
- ✅ **RESTful API**

### DevOps
- ✅ **CI/CD** con GitHub Actions
- ✅ **Automatización** de tests
- ✅ **Build automático**

---

## 🔍 Comandos Útiles

### Backend
```bash
# Compilar
go build ./...

# Tests
go test ./...

# Tests con detalle
go test ./tests/services/... -v

# Limpiar base de datos
rm backend/database.db
```

### Frontend
```bash
# Desarrollo
npm start

# Tests
npm test

# Build producción
npm run build

# Limpiar node_modules
rm -rf node_modules && npm install
```

### Git
```bash
# Status
git status

# Commit
git add .
git commit -m "mensaje"

# Push
git push origin main
```

---

## 📚 Documentación Adicional

- **[decisiones.md](./decisiones.md)** - Decisiones técnicas y justificaciones
- **[backend/tests/desc.md](./backend/tests/desc.md)** - Explicación de tests backend
- **[backend/internal/database/desc.md](./backend/internal/database/desc.md)** - Explicación de base de datos
- **[frontend/src/services/desc.md](./frontend/src/services/desc.md)** - Explicación de servicios HTTP

---

## 🐛 Troubleshooting

### El backend no arranca
```bash
# Verificar que no esté corriendo en otro lado
lsof -i :8080
kill -9 PID_DEL_PROCESO

# Verificar dependencias
cd backend
go mod tidy
```

### El frontend no arranca
```bash
# Reinstalar dependencias
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Los tests fallan
```bash
# Backend: Verificar que no dependa de BD
rm backend/database.db
go test ./tests/services/... -v  # Deben pasar igual

# Frontend: Limpiar cache de Jest
npm test -- --clearCache
npm test
```

### CORS errors
Verificar que el backend tenga el middleware CORS configurado en `router/router.go`

---

## 👥 Autores:
**Carpineti Octavio - Kevin Massholder**  
Ingenieria en sistemas de informacion - UCC
Materia: Ingeniería de Software 3  
Año: 2025

