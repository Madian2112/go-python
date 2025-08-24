# Testing Avanzado en Go

## Introducción

El testing es una parte fundamental del desarrollo de software profesional. Go proporciona un sólido framework de testing integrado en la biblioteca estándar a través del paquete `testing`. En este módulo, exploraremos técnicas avanzadas de testing en Go, incluyendo pruebas unitarias, pruebas de integración, benchmarking, y mocking.

## Repaso de Conceptos Básicos

### Estructura de Pruebas en Go

Las pruebas en Go se escriben como funciones que comienzan con `Test` seguido de un nombre que comienza con mayúscula, en archivos que terminan con `_test.go`.

```go
// archivo: calc_test.go
package calc

import "testing"

func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    if got != want {
        t.Errorf("Add(2, 3) = %d; want %d", got, want)
    }
}
```

### Ejecución de Pruebas

```bash
# Ejecutar todas las pruebas en el paquete actual
go test

# Ejecutar pruebas con salida detallada
go test -v

# Ejecutar una prueba específica
go test -run TestAdd
```

## Técnicas Avanzadas de Testing

### Table-Driven Tests

Los tests basados en tablas permiten ejecutar múltiples casos de prueba con el mismo código.

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -2, -3, -5},
        {"mixed", -2, 3, 1},
        {"zero", 0, 0, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.expected)
            }
        })
    }
}
```

### Subtests

Los subtests permiten agrupar casos de prueba relacionados y proporcionan un mejor control sobre la ejecución.

```go
func TestMath(t *testing.T) {
    t.Run("Addition", func(t *testing.T) {
        if Add(2, 3) != 5 {
            t.Error("Addition failed")
        }
    })
    
    t.Run("Subtraction", func(t *testing.T) {
        if Subtract(5, 3) != 2 {
            t.Error("Subtraction failed")
        }
    })
}
```

### Parallel Testing

Go permite ejecutar pruebas en paralelo para reducir el tiempo total de ejecución.

```go
func TestParallel(t *testing.T) {
    t.Run("group", func(t *testing.T) {
        for i := 0; i < 10; i++ {
            i := i // Capturar variable para la goroutine
            t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
                t.Parallel() // Marca esta prueba para ejecución paralela
                time.Sleep(100 * time.Millisecond) // Simular trabajo
                if i > 10 {
                    t.Error("i is too large")
                }
            })
        }
    })
}
```

### Setup y Teardown

Go no tiene funciones específicas para setup y teardown, pero puedes implementarlas usando funciones auxiliares o `t.Cleanup()`.

```go
func setupTest(t *testing.T) (*Database, func()) {
    db, err := NewDatabase("test.db")
    if err != nil {
        t.Fatalf("Failed to setup test: %v", err)
    }
    
    // Retornar una función de limpieza
    return db, func() {
        db.Close()
        os.Remove("test.db")
    }
}

func TestDatabase(t *testing.T) {
    db, cleanup := setupTest(t)
    defer cleanup() // Asegura que la limpieza se ejecute al final
    
    // Pruebas con la base de datos
    err := db.Insert("key", "value")
    if err != nil {
        t.Errorf("Insert failed: %v", err)
    }
}

// Usando t.Cleanup (Go 1.14+)
func TestWithCleanup(t *testing.T) {
    db, err := NewDatabase("test.db")
    if err != nil {
        t.Fatalf("Failed to setup test: %v", err)
    }
    
    t.Cleanup(func() {
        db.Close()
        os.Remove("test.db")
    })
    
    // Pruebas con la base de datos
}
```

### Helper Functions

Las funciones auxiliares pueden mejorar la legibilidad y reutilización del código de prueba.

```go
func assertNoError(t *testing.T, err error) {
    t.Helper() // Marca esta función como helper para mejorar los mensajes de error
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
}

func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("Got %v, want %v", got, want)
    }
}

func TestWithHelpers(t *testing.T) {
    result, err := SomeFunction()
    assertNoError(t, err)
    assertEqual(t, result, "expected value")
}
```

## Benchmarking

Go proporciona soporte integrado para benchmarking, permitiendo medir el rendimiento de tu código.

### Escribir Benchmarks

```go
func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20) // Función que queremos medir
    }
}

func BenchmarkFibonacciParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Fibonacci(20)
        }
    })
}
```

### Ejecutar Benchmarks

```bash
# Ejecutar todos los benchmarks
go test -bench=.

# Ejecutar un benchmark específico
go test -bench=BenchmarkFibonacci

# Mostrar asignaciones de memoria
go test -bench=. -benchmem
```

### Comparación de Benchmarks

Puedes usar herramientas como `benchstat` para comparar resultados de benchmarks entre diferentes versiones de tu código.

```bash
# Instalar benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# Guardar resultados de la versión actual
go test -bench=. -count=10 > old.txt

# Hacer cambios en el código

# Guardar resultados de la nueva versión
go test -bench=. -count=10 > new.txt

# Comparar resultados
benchstat old.txt new.txt
```

## Mocking y Testing de Interfaces

Go no tiene un framework de mocking integrado, pero su sistema de interfaces facilita la creación de mocks.

### Interfaces para Testing

```go
// Definir una interfaz para el servicio que queremos mockear
type UserRepository interface {
    GetUser(id string) (*User, error)
    SaveUser(user *User) error
}

// Implementación real
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) GetUser(id string) (*User, error) {
    // Implementación real que consulta la base de datos
}

// Mock para testing
type MockUserRepository struct {
    users map[string]*User
}

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{users: make(map[string]*User)}
}

func (m *MockUserRepository) GetUser(id string) (*User, error) {
    user, exists := m.users[id]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }
    return user, nil
}

func (m *MockUserRepository) SaveUser(user *User) error {
    m.users[user.ID] = user
    return nil
}
```

### Testing con Mocks

```go
func TestUserService(t *testing.T) {
    // Crear un mock del repositorio
    mockRepo := NewMockUserRepository()
    
    // Preparar datos de prueba
    mockUser := &User{ID: "123", Name: "Test User"}
    mockRepo.SaveUser(mockUser)
    
    // Crear el servicio con el mock
    service := NewUserService(mockRepo)
    
    // Probar el servicio
    user, err := service.GetUserByID("123")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if user.Name != "Test User" {
        t.Errorf("Expected user name 'Test User', got '%s'", user.Name)
    }
}
```

### Usando Bibliotecas de Mocking

Aunque Go promueve la simplicidad, existen bibliotecas de mocking que pueden ser útiles para casos complejos.

```go
// Usando github.com/golang/mock/gomock
func TestUserServiceWithGoMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockUser := &User{ID: "123", Name: "Test User"}
    
    // Crear mock y definir expectativas
    mockRepo := mocks.NewMockUserRepository(ctrl)
    mockRepo.EXPECT().
        GetUser("123").
        Return(mockUser, nil).
        Times(1)
    
    // Crear servicio con mock
    service := NewUserService(mockRepo)
    
    // Probar el servicio
    user, err := service.GetUserByID("123")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if user.Name != "Test User" {
        t.Errorf("Expected user name 'Test User', got '%s'", user.Name)
    }
}
```

## Testing HTTP Handlers

Go proporciona herramientas para probar handlers HTTP sin necesidad de un servidor real.

```go
func TestHandler(t *testing.T) {
    // Crear un request de prueba
    req, err := http.NewRequest("GET", "/user/123", nil)
    if err != nil {
        t.Fatal(err)
    }
    
    // Crear un ResponseRecorder para registrar la respuesta
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(UserHandler)
    
    // Ejecutar el handler
    handler.ServeHTTP(rr, req)
    
    // Verificar el código de estado
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    
    // Verificar el cuerpo de la respuesta
    expected := `{"id":"123","name":"Test User"}`
    if rr.Body.String() != expected {
        t.Errorf("Handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}
```

## Testing de Integración

Las pruebas de integración verifican que diferentes componentes del sistema funcionen correctamente juntos.

### Configuración para Pruebas de Integración

```go
func setupIntegrationTest(t *testing.T) (*sql.DB, func()) {
    // Usar una base de datos real para pruebas de integración
    db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/testdb?sslmode=disable")
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    
    // Ejecutar migraciones o configurar datos de prueba
    if err := runMigrations(db); err != nil {
        db.Close()
        t.Fatalf("Failed to run migrations: %v", err)
    }
    
    return db, func() {
        // Limpiar datos de prueba
        cleanupDatabase(db)
        db.Close()
    }
}

func TestIntegration(t *testing.T) {
    // Saltar si no estamos ejecutando pruebas de integración
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    db, cleanup := setupIntegrationTest(t)
    defer cleanup()
    
    // Crear repositorios y servicios reales
    userRepo := NewPostgresUserRepository(db)
    userService := NewUserService(userRepo)
    
    // Ejecutar pruebas de integración
    user, err := userService.CreateUser("Test User", "test@example.com")
    if err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }
    
    // Verificar que el usuario se guardó correctamente
    savedUser, err := userService.GetUserByID(user.ID)
    if err != nil {
        t.Fatalf("Failed to get user: %v", err)
    }
    if savedUser.Email != "test@example.com" {
        t.Errorf("Expected email 'test@example.com', got '%s'", savedUser.Email)
    }
}
```

### Ejecutar Pruebas de Integración

```bash
# Ejecutar todas las pruebas, incluyendo integración
go test ./...

# Ejecutar solo pruebas unitarias (rápidas)
go test -short ./...
```

## Cobertura de Código

Go proporciona herramientas para medir la cobertura de código de tus pruebas.

```bash
# Ejecutar pruebas con cobertura
go test -cover

# Generar perfil de cobertura
go test -coverprofile=coverage.out

# Ver informe de cobertura en el navegador
go tool cover -html=coverage.out
```

## Fuzzing

A partir de Go 1.18, el fuzzing está integrado en la herramienta de testing estándar. El fuzzing genera automáticamente entradas para tus funciones para encontrar casos extremos y vulnerabilidades.

```go
func FuzzReverse(f *testing.F) {
    testcases := []string{"Hello, world", " ", "!12345"}
    for _, tc := range testcases {
        f.Add(tc) // Proporcionar casos de semilla
    }
    
    f.Fuzz(func(t *testing.T, orig string) {
        rev := Reverse(orig)
        doubleRev := Reverse(rev)
        if orig != doubleRev {
            t.Errorf("Before: %q, after: %q", orig, doubleRev)
        }
        if utf8.ValidString(orig) && !utf8.ValidString(rev) {
            t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
        }
    })
}
```

## Mejores Prácticas

1. **Escribe tests primero (TDD)**: Considera escribir las pruebas antes de implementar la funcionalidad.

2. **Mantén las pruebas simples y enfocadas**: Cada prueba debe verificar una sola cosa.

3. **Usa table-driven tests**: Facilitan añadir nuevos casos de prueba.

4. **Separa pruebas unitarias e integración**: Las pruebas unitarias deben ser rápidas y no depender de servicios externos.

5. **Usa mocks para dependencias externas**: Evita depender de servicios externos en pruebas unitarias.

6. **Ejecuta pruebas regularmente**: Integra las pruebas en tu flujo de trabajo de desarrollo.

7. **Mantén alta cobertura de código**: Apunta a una cobertura de al menos 80%.

8. **Usa benchmarks para código crítico**: Identifica y optimiza cuellos de botella.

9. **Documenta casos de prueba complejos**: Explica qué estás probando y por qué.

10. **Evita pruebas frágiles**: Las pruebas no deben romperse con cambios menores en la implementación.

## Ejercicios Prácticos

1. Implementa un conjunto de pruebas para una API REST usando table-driven tests.

2. Crea mocks para un servicio que interactúa con una base de datos externa.

3. Escribe benchmarks para comparar diferentes implementaciones de un algoritmo.

4. Implementa pruebas de integración para un sistema que interactúa con múltiples servicios.

5. Usa fuzzing para encontrar casos extremos en una función de procesamiento de texto.

## Conclusión

El testing es una parte esencial del desarrollo de software profesional. Go proporciona un conjunto robusto de herramientas para escribir y ejecutar pruebas, benchmarks y análisis de cobertura. Dominar estas herramientas te permitirá escribir código más confiable y mantenible.

Recuerda que el objetivo del testing no es solo encontrar errores, sino también documentar el comportamiento esperado del código y facilitar futuros cambios. Un buen conjunto de pruebas te da la confianza para refactorizar y mejorar tu código sin miedo a romper la funcionalidad existente.