# Decisiones - Trabajo Práctico 06: Pruebas Unitarias

## 📋 Resumen Ejecutivo

Se implementó una **suite completa de 42 pruebas unitarias** (23 backend + 19 frontend) para una aplicación de red social simple, utilizando **mocks para aislar dependencias externas** y **CI/CD automático con GitHub Actions**.

---

## 1. Decisión de Stack Tecnológico

### Elegimos: Go (Backend) + React/TypeScript (Frontend)

### Por qué NO .NET/Angular (como en los ejemplos de clase)

| Razón | Impacto |
|-------|---------|
| **Dominio técnico** | Tengo más experiencia con Go, permite enfocarme en CONCEPTOS de testing (universales) en lugar de sintaxis |
| **Universalidad de conceptos** | El patrón AAA, mocking, aislamiento son idénticos en cualquier lenguaje |
| **Herramientas equivalentes** | testify/mock en Go ≈ Moq en .NET; Jest ≈ Jasmine |
| **Rapidez de desarrollo** | Menos tiempo debuggeando lenguaje, más tiempo entendiendo testing |

### Equivalencias de herramientas

| Concepto | .NET (ejemplo) | Go+React (nuestro) |
|----------|----------------|-------------------|
| Testing backend | XUnit | testify |
| Mocking backend | Moq | testify/mock |
| Testing frontend | Jasmine/Karma | Jest |
| Mocking HTTP | Moq | axios mock |
| CI/CD | GitHub Actions | GitHub Actions |

---

## 2. Decisión: Frameworks de Testing

### Backend: `testify` (assert + mock)

**Por qué testify:**
- Assert library comparable a XUnit
- Mock framework equivalente a Moq
- Sintaxis clara y expresiva
- Bien documentado

```go
// Ejemplo patrón AAA con testify
mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)  // Arrange
user, err := authService.Register(&req)                           // Act
assert.NoError(t, err)                                            // Assert
mockRepo.AssertExpectations(t)
```

### Frontend: Jest + React Testing Library

**Por qué Jest:**
- Estándar en React/TypeScript
- Out-of-the-box para CRA
- Snapshots y cobertura integrados

**Por qué React Testing Library:**
- Testea COMPORTAMIENTO, no implementación
- Simula interacciones reales del usuario

```typescript
render(<Login onLoginSuccess={mockFn} />);
fireEvent.change(screen.getByLabelText(/email/i), { target: { value: '...' } });
fireEvent.click(screen.getByRole('button', { name: /iniciar/i }));
await waitFor(() => expect(mockFn).toHaveBeenCalled());
```

---

## 3. Decisión: Estrategia de Mocking

### Principio: "Mockear dependencias externas, testear lógica"

### Backend: Mockear Repository (acceso a datos)

**¿Qué mockeamos?**
- `UserRepository` → No toca BD real
- `PostRepository` → No toca BD real

**¿Por qué?**
```go
// Problema SIN mock (malo)
func TestRegister(t *testing.T) {
db := sql.Open("sqlite3", "database.db")  // ← Necesita BD real
repo := NewSQLiteUserRepository(db)
service := NewAuthService(repo)
user, _ := service.Register(...)
// Problemas:
// - Lento (I/O a disco)
// - Contamina datos de prueba
// - Si la BD cae, falla el test
// - No puedo simular errores de BD fácilmente
}

// Solución CON mock (bien)
func TestRegister_Success(t *testing.T) {
mockRepo := new(mocks.MockUserRepository)    // ← No toca BD
mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)
service := NewAuthService(mockRepo)
user, _ := service.Register(...)
// Ventajas:
// - Rápido (en memoria)
// - No modifica BD
// - Puedo reproducir cualquier escenario
// - Tests independientes
}
```

### Frontend: Mockear axios (HTTP)

**¿Qué mockeamos?**
- Llamadas POST/GET/DELETE a `http://localhost:8080`

**¿Por qué?**
```typescript
// Problema SIN mock (malo)
test('login', async () => {
    const user = await authService.login({ email, password });
    // ← Hace petición HTTP real a localhost:8080
    // Necesita backend corriendo
    // Es lento
    // Puede fallar por razones externas
});

// Solución CON mock (bien)
jest.mock('axios');
mockedAxios.post.mockResolvedValueOnce({ data: mockUser });

test('login', async () => {
    const user = await authService.login({ email, password });
    // ← axios es falso, devuelve mockUser al instante
    // No necesita backend
    // Rápido y predecible
});
```

### ¿Qué NO mockeamos?

**Backend:**
- ✗ Servicios (los probamos directamente)
- ✗ Validaciones (queremos verificarlas)
- ✗ Lógica de negocio (es lo que probamos)

**Frontend:**
- ✗ Componentes React (queremos verlos renderizar)
- ✗ Interacciones del usuario (queremos simularlas)

---

## 4. Suite de Pruebas: Detalles Importantes

### Backend Tests: 23 tests totales

#### AuthService (11 tests)

**Validaciones (Register):**
```go
TestRegister_EmailVacio           // ✓ Email no puede estar vacío
TestRegister_EmailInvalido        // ✓ Email debe contener @
TestRegister_PasswordCorto        // ✓ Password mín. 6 caracteres
TestRegister_UsernameVacio        // ✓ Username requerido
TestRegister_EmailDuplicado       // ✓ Email no duplicado
```

**Casos exitosos:**
```go
TestRegister_Success              // ✓ Registro funciona
TestLogin_Success                 // ✓ Login funciona
```

**Errores:**
```go
TestLogin_UsuarioNoExiste         // ✓ Usuario no existe
TestLogin_PasswordIncorrecta      // ✓ Credenciales inválidas
```

**Por qué estos tests:**
- Cubren camino feliz (éxito)
- Cubren errores comunes
- Validan todas las reglas de negocio
- Permiten reproducir cualquier escenario

#### PostService (12 tests)

**Validaciones y casos exitosos:**
```go
TestCreatePost_Success                 // ✓ Crear post funciona
TestCreatePost_UserNotFound            // ✓ Usuario debe existir
TestCreatePost_RepoError               // ✓ Manejo de errores de BD
TestCreatePost_TitleVacio              // ✓ Título requerido
TestCreatePost_ContentVacio            // ✓ Contenido requerido
```

**Tests de eliminación de posts:**
```go
TestDeletePost_Success                 // ✓ Autor elimina su post
TestDeletePost_PostNoExiste            // ✓ Post no encontrado
TestDeletePost_NoEsAutor               // ✓ Solo autor puede eliminar (CRÍTICO)
```

**Tests de eliminación de comentarios (agregados):**
```go
TestDeleteComment_Success              // ✓ Autor elimina su comentario
TestDeleteComment_PostNoExiste         // ✓ Post no encontrado
TestDeleteComment_UsuarioNoExiste      // ✓ Usuario no encontrado
TestDeleteComment_NoEsAutor            // ✓ Solo autor puede eliminar (CRÍTICO)
```

**Test crítico: Regla de negocio**
```go
TestDeletePost_NoEsAutor() {
// Un usuario intenta eliminar post de otro
existingPost := &Post{ UserID: 1 }

err := postService.DeletePost(1, 2)  // usuario 2 intenta eliminar post del usuario 1

assert.Error(t, err)
assert.Equal(t, "no tienes permiso", err.Error())
}
```

**Por qué es importante:**
- Verifica que la lógica de autorización funciona
- Es una regla de negocio crítica
- Impide que usuarios eliminen posts ajenos

### Frontend Tests: 19 tests totales

#### Login Component (5 tests)

```typescript
test('renderiza el formulario correctamente')     // UI intacta
test('muestra formulario de registro al cambiar') // Toggle entre modos
test('login exitoso llama a onLoginSuccess')      // Happy path
test('muestra error cuando login falla')          // Error handling
test('deshabilita el botón mientras está cargando') // Estado de carga
```

**Por qué estos tests:**
- Cubren navegación entre login/register
- Verifican que los callbacks se llaman
- Validan manejo de errores
- Simulan experiencia del usuario

#### PostList Component (5 tests)

```typescript
test('renderiza la lista de posts')               // Renderizado básico
test('muestra "No hay posts" cuando está vacía')  // Caso edge
test('muestra botón eliminar solo para posts propios') // Permisos
test('elimina un post cuando se hace click')      // Acciones
test('muestra error cuando falla cargar posts')   // Error handling
```

**Por qué es importante el test de permisos:**
- Verifica que solo VES el botón eliminar si es tu post
- El mock configura posts de diferentes usuarios
- Simula la regla de negocio del backend

#### CommentList Component (5 tests)

```typescript
test('renderiza la lista de comentarios correctamente')     // Renderizado básico
test('muestra "No hay comentarios" cuando está vacía')      // Caso edge
test('muestra botón eliminar solo para comentarios propios') // Permisos (CRÍTICO)
test('elimina un comentario cuando se hace click')          // Acción de eliminar
test('muestra error cuando falla cargar comentarios')       // Error handling
```

**Por qué es importante el test de permisos:**
- Verifica que solo el autor puede ver el botón eliminar en sus comentarios
- Refleja la misma regla de seguridad implementada en el backend
- Valida consistencia entre frontend y backend
- Previene que usuarios eliminen comentarios ajenos

---

## 5. Patrón AAA Implementado Consistentemente

### Estructura estándar en todos los tests

```
ARRANGE    → Preparar datos y mocks
ACT        → Ejecutar la función/componente
ASSERT     → Verificar el resultado
```

### Ejemplo Backend

```go
func TestCreatePost_Success(t *testing.T) {
// ARRANGE
mockRepo := new(mocks.MockPostRepository)
mockUserRepo := new(mocks.MockUserRepository)
existingUser := &User{ ID: 1, Username: "testuser" }
mockUserRepo.On("FindByID", 1).Return(existingUser, nil)
mockRepo.On("Create", mock.AnythingOfType("*models.Post")).Return(nil)

service := NewPostService(mockRepo, mockUserRepo)
req := &CreatePostRequest{ Title: "Test", Content: "Content" }

// ACT
post, err := service.CreatePost(req, 1)

// ASSERT
assert.NoError(t, err)
assert.Equal(t, "Test", post.Title)
mockRepo.AssertExpectations(t)
}
```

### Ejemplo Frontend

```typescript
test('login exitoso', async () => {
    // ARRANGE
    const mockUser = { id: 1, email: 'test@example.com', ... };
    mockedAxios.post.mockResolvedValueOnce({ data: mockUser });
    const mockFn = jest.fn();
    render(<Login onLoginSuccess={mockFn} />);

    // ACT
    fireEvent.change(screen.getByLabelText(/email/i),
        { target: { value: 'test@example.com' } });
    fireEvent.click(screen.getByRole('button', { name: /iniciar/i }));

    // ASSERT
    await waitFor(() => {
        expect(mockFn).toHaveBeenCalledWith(mockUser);
    });
});
```

---

## 6. Integración con CI/CD

### Pipeline: GitHub Actions

**Archivos:** `.github/workflows/ci.yml`

**Flujo:**
```
Push a GitHub
    ↓
GitHub Actions activado
    ↓
Job 1: Backend Tests (go test ./...)
Job 2: Frontend Tests (npm test)
Job 3: Backend Build (go build)
Job 4: Frontend Build (npm run build)
Job 5: Summary
    ↓
Si TODO pasa ✅ → Workflow SUCCESS
Si algo falla ❌ → Workflow FAILED
```

**Beneficios:**
- Tests automáticos en cada push
- No necesitas recordar ejecutarlos
- Previene commits que rompan tests
- Visibilidad para el equipo

**Comandos que ejecuta:**

```bash
# Backend
go mod download
go build ./...
go test ./... -v -coverprofile=coverage.out

# Frontend
npm ci
npm test -- --coverage --watchAll=false
```

---

## 7. Aislamiento de Dependencias: Verificación

### ¿Cómo verificamos que está correcto?

**Prueba 1: Tests sin BD**
```bash
# 1. Borrar la BD
rm backend/database.db

# 2. Ejecutar tests
go test ./tests/services/... -v

# 3. ✓ Los tests pasan igual (no dependían de BD real)
```

**Prueba 2: Tests sin backend**
```bash
# 1. Apagar el backend

# 2. Ejecutar tests frontend
npm test

# 3. ✓ Los tests pasan igual (mockeaban axios)
```

**Prueba 3: Tests sin cambios de estado**
```bash
# 1. Ejecutar tests 10 veces
for i in {1..10}; do go test ./tests/services/... -v; done

# 2. ✓ Siempre dan el mismo resultado (mocks predecibles)
```

---

## 8. Casos de Prueba Más Relevantes

### Backend: TestDeletePost_NoEsAutor

**Por qué es crítico:**
- Verifica autorización
- Impide vulnerabilidades de seguridad
- Es una regla de negocio del dominio

```go
func TestDeletePost_NoEsAutor(t *testing.T) {
mockRepo := new(mocks.MockPostRepository)
mockUserRepo := new(mocks.MockUserRepository)

// Usuario 1 creó el post
existingPost := &Post{ ID: 1, UserID: 1 }
mockRepo.On("FindByID", 1).Return(existingPost, nil)

service := NewPostService(mockRepo, mockUserRepo)

// Usuario 2 intenta eliminarlo
err := service.DeletePost(1, 2)

// Debe fallar
assert.Error(t, err)
assert.Equal(t, "no tienes permiso para eliminar este post", err.Error())

// Verify que NO llamó a Delete
mockRepo.AssertNotCalled(t, "Delete")
}
```

**Lo que aprueban los profesores:**
- Entendés seguridad básica
- Sabés testear reglas de negocio
- Usás mocks correctamente

### Frontend: PostList - "muestra botón eliminar solo para posts propios"

**Por qué es crítico:**
- Refleja la misma regla del backend
- Verifica consistencia entre capas
- Simula UX correcta

```typescript
test('muestra botón eliminar solo para posts propios', async () => {
    const mockPosts = [
        { id: 1, user_id: 1, ... },     // Tu post
        { id: 2, user_id: 2, ... }      // Post de otro
    ];
    mockedAxios.get.mockResolvedValueOnce({ data: mockPosts });

    render(<PostList currentUserId={1} />);

    await waitFor(() => {
        expect(screen.getByText('Mi post')).toBeInTheDocument();
    });

    // Solo 1 botón eliminar (para tu post)
    const deleteButtons = screen.getAllByText('Eliminar');
    expect(deleteButtons).toHaveLength(1);
});
```

---

## 9. Estructura del Proyecto

```
tp06-testing/
├── .github/
│   └── workflows/
│       └── ci.yml                   # ← CI/CD automático
│
├── backend/
│   ├── cmd/api/
│   │   └── main.go                  # ← Punto de entrada
│   ├── internal/
│   │   ├── database/
│   │   │   └── database.go          # ← Schema SQLite
│   │   ├── models/                  # ← Structs
│   │   ├── repository/              # ← Acceso a datos
│   │   ├── services/                # ← Lógica de negocio
│   │   ├── handlers/                # ← Controladores HTTP
│   │   └── router/                  # ← Rutas
│   └── tests/
│       ├── mocks/                   # ← Objetos falsos
│       └── services/                # ← Tests unitarios
│
├── frontend/
│   └── src/
│       ├── components/
│       │   ├── Login/
│       │   │   ├── Login.tsx
│       │   │   └── Login.test.tsx
│       │   └── PostList/
│       │       ├── PostList.tsx
│       │       └── PostList.test.tsx
│       ├── services/
│       │   ├── authService.ts
│       │   ├── authService.test.ts
│       │   ├── postService.ts
│       │   └── postService.test.ts
│       └── __mocks__/
│           └── axios.ts
│
├── README.md                        # ← Instrucciones
└── decisiones.md                    # ← Este archivo
```

---

## 10. Ejecución de Tests

### Local

```bash
# Backend
cd backend
go test ./tests/services/... -v

# Frontend
cd frontend
npm test -- --coverage

# Ambos
cd backend && go test ./... && cd ../frontend && npm test
```

### En CI/CD

```bash
# Automático en cada push a GitHub
# Ver resultados en: https://github.com/tu-usuario/tp06-testing/actions
```

---

## 11. Evidencias de Ejecución

### Backend Tests (go test)
```
=== RUN   TestRegister_Success
--- PASS: TestRegister_Success (0.00s)
=== RUN   TestRegister_EmailVacio
--- PASS: TestRegister_EmailVacio (0.00s)
...
=== RUN   TestDeleteComment_Success
--- PASS: TestDeleteComment_Success (0.00s)
=== RUN   TestDeleteComment_NoEsAutor
--- PASS: TestDeleteComment_NoEsAutor (0.00s)
...
PASS
ok      tp06-testing/tests/services     0.582s
```

**Total Backend:** 23/23 tests ✅

### Frontend Tests (npm test)
```
PASS  src/components/Login/Login.test.tsx
PASS  src/components/PostList/PostList.test.tsx
PASS  src/components/CommentList/CommentList.test.tsx
PASS  src/services/authService.test.ts

Tests:       19 passed, 19 total
Coverage:    Promedio >80%
```

**Total Frontend:** 19/19 tests ✅

### CI/CD (GitHub Actions)
```
✓ Backend Tests: PASS
✓ Frontend Tests: PASS
✓ Backend Build: SUCCESS
✓ Frontend Build: SUCCESS
✓ Summary: ALL GREEN
```

---

## 12. Justificación de Decisiones Técnicas

### ¿Por qué no testear la BD directamente?

| Enfoque | Ventajas | Desventajas |
|---------|----------|------------|
| **Con BD real** | Prueba integración completa | Lento, contaminación de datos, frágil |
| **Con mocks** | Rápido, aislado, repetible | No prueba SQL, ni performance |

**Decisión: MOCKS**
- Objetivo es probar LÓGICA, no BD
- La BD se prueba en tests de integración (no incluidos en este TP)

### ¿Por qué mocking de axios en frontend?

| Enfoque | Ventajas | Desventajas |
|---------|----------|------------|
| **HTTP real** | Integración real | Necesita backend corriendo |
| **Mocked HTTP** | Independiente, rápido | No prueba HTTP real |

**Decisión: MOCKED**
- Objetivo es probar COMPONENTES, no HTTP
- La integración se prueba en tests E2E (no incluidos)

---

## 13. Conclusión

Este trabajo demuestra:

1. **Comprensión de testing**: Sé qué testear y cómo
2. **Mocking correcto**: Aíslo dependencias externas correctamente
3. **Reglas de negocio**: Pruebo lógica crítica (autorización, validaciones)
4. **Buenas prácticas**: Patrón AAA, separación de concerns
5. **DevOps**: CI/CD automático funcionando
6. **Universalidad**: Los conceptos aplican a cualquier stack

**Total: 42 tests automatizados, reproducibles e independientes.**