# Arquitectura del Proyecto - TP06 Pruebas Unitarias

## 🎯 Objetivo de esta guía
Esta guía te explica **QUÉ es cada parte del proyecto, POR QUÉ existe, y CÓMO se relaciona con las pruebas unitarias**. Úsala para preparar tu defensa oral.

---

## 📊 Visión General - Arquitectura en Capas

```
┌─────────────────────────────────────────────────┐
│           USUARIO (Navegador)                   │
└─────────────────────────────────────────────────┘
                      ↕ HTTP
┌─────────────────────────────────────────────────┐
│  FRONTEND (React + TypeScript)                  │
│  - Components (UI)                              │
│  - Services (HTTP)                              │
│  └─ Tests: Jest + React Testing Library         │
└─────────────────────────────────────────────────┘
                      ↕ REST API
┌─────────────────────────────────────────────────┐
│  BACKEND (Go)                                   │
│  - Handlers (HTTP)                              │
│  - Services (Lógica de Negocio) ← TESTEADO     │
│  - Repository (Acceso a Datos)                  │
│  └─ Tests: testify + mocks                      │
└─────────────────────────────────────────────────┘
                      ↕ SQL
┌─────────────────────────────────────────────────┐
│  DATABASE (SQLite)                              │
└─────────────────────────────────────────────────┘
```

---

## 🏗️ Backend - Estructura Detallada

### Directorio: `backend/`

```
backend/
├── cmd/api/
│   └── main.go                    ← Entry point (NO se testea)
├── internal/
│   ├── database/
│   │   └── database.go            ← Inicializa BD (NO se testea)
│   ├── models/
│   │   ├── user.go                ← Structs (NO se testea)
│   │   └── post.go
│   ├── repository/
│   │   ├── user_repository.go     ← Acceso a datos (SE MOCKEA)
│   │   └── post_repository.go
│   ├── services/
│   │   ├── auth_service.go        ← LÓGICA DE NEGOCIO (SE TESTEA) ✅
│   │   └── post_service.go        ← LÓGICA DE NEGOCIO (SE TESTEA) ✅
│   ├── handlers/
│   │   ├── auth_handler.go        ← HTTP controllers (NO se testea)
│   │   └── post_handler.go
│   └── router/
│       └── router.go               ← Rutas HTTP (NO se testea)
└── tests/
    ├── mocks/
    │   ├── user_repository_mock.go ← Mock FALSO del repository
    │   └── post_repository_mock.go
    └── services/
        ├── auth_service_test.go    ← 11 tests unitarios ✅
        └── post_service_test.go    ← 12 tests unitarios ✅
```

### ❓ Pregunta clave: ¿Por qué SOLO testeas services?

**Respuesta para defensa:**

> "En una arquitectura en capas, cada capa tiene una responsabilidad única:
>
> - **Repository**: Solo ejecuta SQL. No tiene lógica para testear. Se mockea.
> - **Handlers**: Solo recibe HTTP y llama al service. Lógica mínima.
> - **Services**: Contiene TODA la lógica de negocio: validaciones, reglas, permisos. **Aquí es donde están los bugs potenciales**, por eso es lo que testeo.
>
> Los tests unitarios prueban LÓGICA, no I/O. Repository y Handlers hacen I/O, Services hace lógica."

---

## 📱 Frontend - Estructura Detallada

### Directorio: `frontend/src/`

```
frontend/src/
├── index.tsx                      ← Entry point (NO se testea)
├── App.tsx                        ← Orquestador (NO se testea)
├── components/
│   ├── Login/
│   │   ├── Login.tsx              ← Componente (SE TESTEA) ✅
│   │   ├── Login.test.tsx         ← 5 tests
│   │   └── Login.css
│   ├── PostList/
│   │   ├── PostList.tsx           ← Componente (SE TESTEA) ✅
│   │   ├── PostList.test.tsx      ← 5 tests
│   │   └── PostList.css
│   ├── CommentList/
│   │   ├── CommentList.tsx        ← Componente (SE TESTEA) ✅
│   │   ├── CommentList.test.tsx   ← 5 tests
│   │   └── CommentList.css
│   ├── CreatePost/
│   │   ├── CreatePost.tsx         ← Sin tests (solo presentación)
│   │   └── CreatePost.css
│   ├── CommentForm/               ← Sin tests (solo presentación)
│   └── PostDetail/                ← Sin tests (solo presentación)
├── services/
│   ├── authService.ts             ← HTTP client (SE TESTEA) ✅
│   ├── authService.test.ts        ← 4 tests
│   └── postService.ts             ← HTTP client (parcialmente testeado)
├── __mocks__/
│   └── axios.ts                   ← Mock FALSO de axios
└── types/
    └── index.ts                   ← TypeScript types (NO se testea)
```

### ❓ Pregunta clave: ¿Por qué algunos componentes NO tienen tests?

**Respuesta para defensa:**

> "Prioricé testear componentes con LÓGICA:
>
> - **Login**: Tiene validaciones, manejo de errores, estados (loading/error)
> - **PostList/CommentList**: Tienen lógica de permisos (solo autor ve botón eliminar)
>
> Los componentes sin tests (CreatePost, CommentForm, PostDetail) son mayormente presentacionales, solo reciben props y los muestran. Sin lógica compleja para testear. Para el TP7 agregaré tests para estos."

---

## 🔑 Conceptos Clave para la Defensa

### 1. ¿Qué es una Prueba Unitaria?

**Definición académica:**
> Una prueba unitaria testea UNA unidad de código (función/método/componente) de forma AISLADA, sin depender de sistemas externos (BD, APIs, filesystem).

**En tu proyecto:**
- Backend: Una función del service (ej: `Register()`)
- Frontend: Un componente React (ej: `<Login />`)

### 2. ¿Qué es un Mock?

**Definición académica:**
> Un mock es un objeto FALSO que simula el comportamiento de uno real, permitiendo controlar sus respuestas en los tests.

**En tu proyecto:**
- Backend: `MockUserRepository` simula el repository sin tocar la BD
- Frontend: Mock de axios simula llamadas HTTP sin tocar el servidor

**Analogía para explicar:**
> "Es como actuar en una obra de teatro. El actor (service) interactúa con utilería falsa (mock repository) en lugar de objetos reales. Así podemos repetir la escena sin consecuencias reales."

### 3. Patrón AAA (Arrange, Act, Assert)

**Estructura de TODOS tus tests:**

```go
// ARRANGE: Preparar datos y mocks
mockRepo := new(mocks.MockUserRepository)
mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)

// ACT: Ejecutar la función que estamos probando
user, err := authService.Register(&req)

// ASSERT: Verificar el resultado
assert.NoError(t, err)
assert.Equal(t, "test@example.com", user.Email)
```

**Por qué este patrón:**
> "AAA hace que los tests sean legibles y mantenibles. Cualquier desarrollador puede entender: qué preparo, qué ejecuto, qué espero."

---

## 🎯 Separación de Responsabilidades

### Backend - Capas y sus roles

| Capa | Responsabilidad | ¿Se testea? | ¿Por qué? |
|------|----------------|-------------|-----------|
| **Models** | Definir estructuras de datos | ❌ NO | Solo structs, sin lógica |
| **Repository** | Ejecutar SQL | ❌ SE MOCKEA | I/O, no lógica de negocio |
| **Services** | Lógica de negocio | ✅ SÍ | Aquí está la lógica crítica |
| **Handlers** | Recibir HTTP | ❌ NO | Solo delega al service |
| **Router** | Definir rutas | ❌ NO | Configuración, sin lógica |
| **Database** | Inicializar BD | ❌ NO | Solo setup |

### Frontend - Componentes y su testeo

| Tipo | Ejemplo | ¿Se testea? | ¿Por qué? |
|------|---------|-------------|-----------|
| **Entry Point** | index.tsx, App.tsx | ❌ NO | Solo orquestación |
| **Componentes con lógica** | Login, PostList | ✅ SÍ | Validaciones, estados, permisos |
| **Componentes presentacionales** | CreatePost, CommentForm | ⏳ OPCIONAL | Poca lógica, solo UI |
| **Services** | authService | ✅ SÍ | Lógica de HTTP |
| **Types** | index.ts | ❌ NO | Solo definiciones TypeScript |

---

## 🔄 Flujo de una Prueba Unitaria

### Ejemplo: Test de Register en Backend

```
1. PREPARACIÓN (Arrange)
   ├─ Crear mock del repository
   ├─ Configurar qué debe devolver el mock
   └─ Crear datos de prueba (email, password, username)

2. EJECUCIÓN (Act)
   └─ Llamar a authService.Register(datos)

3. VERIFICACIÓN (Assert)
   ├─ ¿Devolvió un usuario?
   ├─ ¿El email es correcto?
   ├─ ¿No hay error?
   └─ ¿Se llamó al repository correctamente?

4. RESULTADO
   └─ Test PASS ✅ o FAIL ❌
```

**Sin mock:**
```
authService.Register()
    ↓
repository.Create() ← Va a BD REAL
    ↓
SQLite escribe en disco
    ↓
❌ Lento, dependiente, contamina datos
```

**Con mock:**
```
authService.Register()
    ↓
mockRepository.Create() ← FALSO, devuelve lo configurado
    ↓
✅ Rápido, independiente, predecible
```

---

## 📊 Cobertura (Coverage)

### ¿Qué es Coverage?

**Definición:**
> Porcentaje de líneas de código que se ejecutan cuando corren los tests.

### Tu Coverage Actual

```
Backend (services):  54.1%
Frontend (components): 56.4%
```

### ¿Por qué NO es 100%?

**Respuesta para defensa:**

> "El coverage mide líneas ejecutadas, no cantidad de tests. Tengo 42 tests, pero hay funciones que NO testeé:
>
> **Backend:**
> - GetAllPosts() - 0%
> - GetPostByID() - 0%
> - CreateComment() - 0%
> - GetCommentsByPostID() - 0%
>
> **Frontend:**
> - CommentForm - 0%
> - CreatePost - 0%
> - PostDetail - 0%
>
> Las funciones que SÍ testeé tienen 85-90% de coverage. El promedio es 55% porque hay código sin tests. Para el TP7 subiré a 70%."

### ¿Por qué NO testear todo?

**Respuesta académica:**

> "100% de coverage NO garantiza 0 bugs. Es un equilibrio entre:
>
> - **Costo** (tiempo de escribir tests)
> - **Beneficio** (bugs evitados)
>
> Prioricé testear LÓGICA CRÍTICA:
> - Validaciones (email, password)
> - Reglas de negocio (solo autor puede eliminar)
> - Operaciones de escritura (Create, Delete)
>
> Las funciones de lectura (Get) son simples mapeos de BD, con poco riesgo de bugs."

---

## 🎓 Resumen para Defensa Oral

### Preguntas y Respuestas Clave

**P: ¿Por qué Go en lugar de .NET?**
> "Tengo más dominio de Go. Los conceptos de testing son universales: mocking, AAA, aislamiento. Lo importante es entender los principios, no la sintaxis."

**P: ¿Cómo decidiste qué mockear?**
> "Mockeé dependencias EXTERNAS: Repository (BD) y axios (HTTP). La regla es: si hace I/O (input/output), se mockea. Si tiene lógica, se testea."

**P: ¿Cómo sabés que tus tests prueban lo correcto?**
> "Tengo tests positivos (caso éxito) y negativos (errores). Por ejemplo, `TestDeletePost_NoEsAutor` verifica que SOLO el autor puede eliminar. Eso es una regla de negocio real. Si el test pasa, la regla funciona."

**P: ¿Por qué 55% de coverage y no 100%?**
> "Prioricé lógica crítica. Las funciones testeadas tienen 85-90% coverage individual. El promedio es 55% porque hay funciones de lectura simples sin tests. Para TP7 subiré a 70%+ agregando tests para esas funciones."

**P: ¿Qué es el patrón AAA?**
> "Arrange-Act-Assert. Preparo datos y mocks (Arrange), ejecuto la función (Act), verifico el resultado (Assert). Todos mis tests siguen este patrón. Facilita la lectura y mantenimiento."

---

## 📚 Términos Clave para Usar en la Defensa

- **Prueba Unitaria**: Test de una unidad aislada
- **Mock**: Objeto falso que simula uno real
- **Stub**: Mock que devuelve valores predefinidos
- **Aislamiento**: Probar sin dependencias externas
- **Coverage**: % de código ejecutado por tests
- **Patrón AAA**: Arrange, Act, Assert
- **Lógica de negocio**: Reglas y validaciones del dominio
- **I/O**: Input/Output (BD, HTTP, archivos)
- **Dependency Injection**: Pasar dependencias al constructor
- **Interface**: Contrato que permite crear mocks

---

## ✅ Checklist de Comprensión

Antes de la defensa, asegurate de poder explicar:

- [ ] ¿Qué es cada directorio del proyecto?
- [ ] ¿Por qué SOLO testeas services en backend?
- [ ] ¿Qué es un mock y por qué lo usás?
- [ ] ¿Cómo funciona el patrón AAA?
- [ ] ¿Por qué tu coverage es 55% y no 100%?
- [ ] ¿Qué funciones tienen 0% de coverage y por qué?
- [ ] ¿Cómo sabés que tus tests prueban lo correcto?
- [ ] ¿Qué es dependency injection y cómo lo usás?