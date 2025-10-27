# Implementación de Funcionalidades - TP06

## 🎯 Objetivo de esta guía
Mapear cada requisito del TP con su implementación, tests y evidencias. Úsala para demostrar que cumpliste TODO.

---

## ✅ Requisito 1: Configuración del Entorno de Testing

### ¿Qué pedía el TP?
- Configurar frameworks de testing apropiados
- Configurar mocking frameworks

### ¿Qué hiciste?

#### Backend (Go)
```bash
# Instalaste (en go.mod):
github.com/stretchr/testify v1.11.1

# Qué incluye:
- testify/assert: Aserciones legibles
- testify/mock: Framework de mocking
```

**Dónde verlo:**
- Archivo: `backend/go.mod`
- Línea: `github.com/stretchr/testify v1.11.1`

**Cómo demostrarlo:**
```bash
cd backend
cat go.mod | grep testify
```

#### Frontend (React)
```bash
# Instalaste (en package.json):
@testing-library/react
@testing-library/jest-dom
jest (incluido en react-scripts)
```

**Dónde verlo:**
- Archivo: `frontend/package.json`
- Sección: `devDependencies`

**Cómo demostrarlo:**
```bash
cd frontend
cat package.json | grep testing-library
```

### ¿Por qué elegiste estas herramientas?

**Para la defensa:**

> "**Backend - testify:**
> - Es el estándar de facto en Go para testing
> - Tiene assert (aserciones legibles) y mock (mocking)
> - Equivalente a NUnit (assert) + Moq (mock) en .NET
>
> **Frontend - Jest + React Testing Library:**
> - Jest viene incluido con Create React App
> - React Testing Library promueve testear COMPORTAMIENTO del usuario, no detalles de implementación
> - Equivalente a Jasmine/Karma en Angular"

---

## ✅ Requisito 2: Implementación de Pruebas Unitarias

### ¿Qué pedía el TP?
- Crear tests para lógica de negocio en backend
- Implementar tests para componentes en frontend
- Utilizar patrón AAA

### ¿Qué hiciste?

#### Backend: 23 tests

| Archivo | Tests | Qué testea |
|---------|-------|------------|
| `auth_service_test.go` | 11 | Register (6) + Login (5) |
| `post_service_test.go` | 12 | CreatePost (5) + DeletePost (3) + DeleteComment (4) |

**Ejemplo de test con patrón AAA:**

```go
func TestRegister_EmailVacio(t *testing.T) {
    // ARRANGE: Preparar
    mockRepo := new(mocks.MockUserRepository)
    authService := services.NewAuthService(mockRepo)
    req := &models.RegisterRequest{
        Email:    "",  // Email vacío para probar validación
        Password: "123456",
        Username: "testuser",
    }

    // ACT: Ejecutar
    user, err := authService.Register(req)

    // ASSERT: Verificar
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Equal(t, "el email es requerido", err.Error())
    
    // Verificar que NO llamó a la BD (falló antes)
    mockRepo.AssertNotCalled(t, "FindByEmail")
}
```

**Dónde ejecutarlo:**
```bash
cd backend
go test ./tests/services/... -v
```

**Captura esperada:**
```
=== RUN   TestRegister_EmailVacio
--- PASS: TestRegister_EmailVacio (0.00s)
```

#### Frontend: 19 tests

| Archivo | Tests | Qué testea |
|---------|-------|------------|
| `Login.test.tsx` | 5 | Renderizado, validaciones, estados |
| `PostList.test.tsx` | 5 | Listado, eliminación, permisos |
| `CommentList.test.tsx` | 5 | Listado, eliminación, permisos |
| `authService.test.ts` | 4 | Llamadas HTTP mockeadas |

**Ejemplo de test con patrón AAA:**

```typescript
test('login exitoso llama a onLoginSuccess', async () => {
    // ARRANGE: Preparar mock y datos
    const mockUser = { id: 1, email: 'test@example.com', ... };
    mockedAxios.post.mockResolvedValueOnce({ data: mockUser });
    const mockOnLoginSuccess = jest.fn();
    render(<Login onLoginSuccess={mockOnLoginSuccess} />);
    
    // ACT: Simular acciones del usuario
    fireEvent.change(screen.getByLabelText(/email/i), 
        { target: { value: 'test@example.com' } });
    fireEvent.click(screen.getByRole('button', { name: /iniciar/i }));
    
    // ASSERT: Verificar comportamiento
    await waitFor(() => {
        expect(mockOnLoginSuccess).toHaveBeenCalledWith(mockUser);
    });
});
```

**Dónde ejecutarlo:**
```bash
cd frontend
npm test -- --watchAll=false
```

**Captura esperada:**
```
PASS  src/components/Login/Login.test.tsx
  ✓ login exitoso llama a onLoginSuccess (13ms)
```

---

## ✅ Requisito 3: Testing Avanzado

### ¿Qué pedía el TP?
- Crear mocks para dependencias externas
- Tests para manejo de excepciones
- Tests para casos edge y validaciones

### ¿Qué hiciste?

#### 3.1 Mocks para Dependencias Externas

**Backend - Mock del Repository:**

```go
// Archivo: tests/mocks/user_repository_mock.go
type MockUserRepository struct {
    mock.Mock  // ← Hereda capacidad de mockear
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}
```

**Por qué es un mock:**
> "No toca la BD real. Devuelve lo que yo configure en el test. Así puedo simular que el usuario existe, no existe, o que la BD falla."

**Frontend - Mock de axios:**

```typescript
// Archivo: src/__mocks__/axios.ts
const axiosMock = {
  post: jest.fn(() => Promise.resolve({ data: {} })),
  get: jest.fn(() => Promise.resolve({ data: {} })),
  delete: jest.fn(() => Promise.resolve({ data: {} })),
};
```

**Por qué es un mock:**
> "No hace peticiones HTTP reales. Devuelve respuestas configuradas en el test. Así puedo simular éxito, error 404, error 500, sin necesitar el servidor."

#### 3.2 Tests de Manejo de Excepciones

**Ejemplo - Login con credenciales inválidas:**

```go
func TestLogin_PasswordIncorrecta(t *testing.T) {
    mockRepo := new(mocks.MockUserRepository)
    authService := services.NewAuthService(mockRepo)

    existingUser := &models.User{
        Email:    "test@example.com",
        Password: "123456",  // Password correcta
    }
    mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

    creds := &models.Credentials{
        Email:    "test@example.com",
        Password: "wrongpassword",  // ← Password INCORRECTA
    }

    user, err := authService.Login(creds)

    // Debe fallar con error específico
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Equal(t, "credenciales inválidas", err.Error())
}
```

**Dónde ejecutarlo:**
```bash
cd backend
go test ./tests/services/... -v -run TestLogin_PasswordIncorrecta
```

#### 3.3 Tests de Casos Edge y Validaciones

**Casos edge testeados:**

| Test | Caso Edge | Por qué es importante |
|------|-----------|----------------------|
| `TestRegister_EmailVacio` | Email = "" | Evita registros sin email |
| `TestRegister_EmailInvalido` | Email sin @ | Valida formato |
| `TestRegister_PasswordCorto` | Password = "123" | Seguridad mínima |
| `TestLogin_UsuarioNoExiste` | Usuario no en BD | Manejo de 404 |
| `TestDeletePost_NoEsAutor` | Usuario 2 elimina post del 1 | **CRÍTICO: Seguridad** |

**Test más importante - Regla de seguridad:**

```go
func TestDeletePost_NoEsAutor(t *testing.T) {
    mockRepo := new(mocks.MockPostRepository)
    mockUserRepo := new(mocks.MockUserRepository)
    postService := services.NewPostService(mockRepo, mockUserRepo)

    existingPost := &Post{ ID: 1, UserID: 1 }  // Post del usuario 1
    mockRepo.On("FindByID", 1).Return(existingPost, nil)

    // Usuario 2 intenta eliminar
    err := postService.DeletePost(1, 2)

    // DEBE FALLAR
    assert.Error(t, err)
    assert.Equal(t, "no tienes permiso", err.Error())
    mockRepo.AssertNotCalled(t, "Delete")  // ← NO debe eliminar
}
```

**Por qué es crítico:**
> "Este test verifica una regla de SEGURIDAD. Si falla, cualquier usuario podría eliminar posts de otros. Es el tipo de bug que causa incidentes de seguridad graves en producción."

---

## ✅ Requisito 4: Integración con CI/CD

### ¿Qué pedía el TP?
- Configurar ejecución automática de tests en pipeline

### ¿Qué hiciste?

**Archivo: `.github/workflows/ci.yml`**

```yaml
jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run backend tests
        working-directory: ./backend
        run: |
          go mod download
          go test ./tests/services/... -v -cover -coverpkg=./internal/services
  
  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Run frontend tests
        working-directory: ./frontend
        run: |
          npm ci
          npm test -- --coverage --watchAll=false
```

**Qué hace:**
1. En cada `git push` a GitHub
2. GitHub Actions se activa automáticamente
3. Ejecuta tests de backend (Go)
4. Ejecuta tests de frontend (React)
5. Si TODO pasa ✅ → OK para mergear
6. Si algo falla ❌ → Bloquea el merge

**Dónde verlo:**
- GitHub → Tu repositorio → Tab "Actions"
- Cada push aparece como un "workflow run"

**Captura esperada:**

```
✅ backend-tests (23 passed)
✅ frontend-tests (19 passed)
✅ All checks have passed
```

**Por qué es importante:**
> "Sin CI/CD, los tests son opcionales. Con CI/CD, son obligatorios. Nadie puede hacer push sin que los tests pasen. Esto previene bugs en producción."

---

## ✅ Requisito 5: Evidencias y Documentación

### ¿Qué pedía el TP?
- Capturas de ejecución de tests
- Documentar estrategia de testing en `decisiones.md`

### ¿Qué hiciste?

#### 5.1 Capturas de Ejecución

**Backend - Ejecutar y capturar:**
```bash
cd backend
go test ./tests/services/... -v -cover -coverpkg=./internal/services
```

**Captura esperada:**
```
=== RUN   TestRegister_Success
--- PASS: TestRegister_Success (0.00s)
...
PASS
coverage: 54.1% of statements in ./internal/services
ok      tp06-testing/tests/services     0.582s
```

**Frontend - Ejecutar y capturar:**
```bash
cd frontend
npm test -- --coverage --watchAll=false
```

**Captura esperada:**
```
PASS  src/components/Login/Login.test.tsx
PASS  src/components/PostList/PostList.test.tsx
...
Test Suites: 4 passed, 4 total
Tests:       19 passed, 19 total
Coverage:    56.4%
```

#### 5.2 Documentación en decisiones.md

**Archivo: `decisiones.md`**

Contiene:
- ✅ Stack tecnológico y justificación
- ✅ Frameworks de testing elegidos
- ✅ Estrategia de mocking
- ✅ Casos de prueba relevantes explicados
- ✅ Coverage actual y justificación
- ✅ Evidencias de ejecución

**Secciones clave:**
1