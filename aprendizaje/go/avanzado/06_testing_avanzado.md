# Testing Avanzado en Go

## Introducción

El testing es una parte fundamental del desarrollo de software de calidad. Go proporciona un sólido framework de testing en su biblioteca estándar, pero hay muchas técnicas avanzadas que pueden mejorar significativamente la calidad y cobertura de nuestras pruebas. En este documento, exploraremos técnicas avanzadas de testing en Go, incluyendo pruebas unitarias, de integración, de rendimiento, mocking y más.

## Fundamentos de Testing en Go

Antes de adentrarnos en técnicas avanzadas, repasemos brevemente los fundamentos del testing en Go:

### Estructura Básica de Tests

```go
package mypackage_test // Sufijo _test para separar el código de prueba

import (
    "testing"
    
    "mymodule/mypackage" // Importar el paquete a probar
)

func TestMyFunction(t *testing.T) {
    result := mypackage.MyFunction()
    expected := "expected result"
    
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

### Ejecutar Tests

```bash
# Ejecutar todos los tests en el directorio actual
go test ./...

# Ejecutar tests con verbose
go test -v ./...

# Ejecutar un test específico
go test -run TestMyFunction

# Ejecutar tests con cobertura
go test -cover ./...

# Generar informe de cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Técnicas Avanzadas de Testing

### Table-Driven Tests

Los table-driven tests permiten ejecutar múltiples casos de prueba con el mismo código, lo que reduce la duplicación y mejora la mantenibilidad:

```go
func TestCalculate(t *testing.T) {
    testCases := []struct {
        name     string
        input    int
        expected int
    }{
        {"positive number", 5, 10},
        {"zero", 0, 0},
        {"negative number", -5, -10},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := Calculate(tc.input)
            if result != tc.expected {
                t.Errorf("Expected %d, got %d", tc.expected, result)
            }
        })
    }
}
```

### Subtests y Setup/Teardown

Go permite organizar tests en subtests, lo que facilita la implementación de setup y teardown comunes:

```go
func TestDatabase(t *testing.T) {
    // Setup común
    db, err := setupTestDatabase()
    if err != nil {
        t.Fatalf("Failed to setup test database: %v", err)
    }
    defer db.Close() // Teardown común
    
    // Subtest 1
    t.Run("Insert", func(t *testing.T) {
        err := db.Insert("test_data")
        if err != nil {
            t.Errorf("Insert failed: %v", err)
        }
    })
    
    // Subtest 2
    t.Run("Query", func(t *testing.T) {
        result, err := db.Query("test_query")
        if err != nil {
            t.Errorf("Query failed: %v", err)
        }
        if result != "expected" {
            t.Errorf("Expected 'expected', got '%s'", result)
        }
    })
}
```

### Helpers de Testing

Los helpers de testing mejoran la legibilidad y reducen la duplicación en los tests:

```go
func TestComplex(t *testing.T) {
    // Helper para verificar errores
    assertNoError := func(t *testing.T, err error, msg string) {
        t.Helper() // Marca esta función como helper para mejorar los mensajes de error
        if err != nil {
            t.Fatalf("%s: %v", msg, err)
        }
    }
    
    // Helper para verificar igualdad
    assertEqual := func(t *testing.T, got, want interface{}, msg string) {
        t.Helper()
        if !reflect.DeepEqual(got, want) {
            t.Errorf("%s: got %v, want %v", msg, got, want)
        }
    }
    
    // Usar los helpers
    result, err := ComplexFunction()
    assertNoError(t, err, "ComplexFunction failed")
    assertEqual(t, result, "expected", "ComplexFunction result")
}
```

### Testing Paralelo

Go permite ejecutar tests en paralelo, lo que puede reducir significativamente el tiempo de ejecución:

```go
func TestParallel(t *testing.T) {
    // Definir subtests
    testCases := []struct {
        name  string
        input int
    }{
        {"case1", 1},
        {"case2", 2},
        {"case3", 3},
        // Muchos más casos...
    }
    
    for _, tc := range testCases {
        tc := tc // Capturar variable para ejecución paralela
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel() // Marcar este subtest para ejecución paralela
            
            // Test que puede tomar tiempo
            time.Sleep(100 * time.Millisecond)
            result := SlowFunction(tc.input)
            
            if result != tc.input*2 {
                t.Errorf("Expected %d, got %d", tc.input*2, result)
            }
        })
    }
}
```

### Mocking y Stubbing

#### Interfaces para Testing

Go facilita el mocking a través de interfaces:

```go
// En el código de producción
type DataStore interface {
    Get(id string) (string, error)
    Save(id, data string) error
}

type Service struct {
    store DataStore
}

func (s *Service) ProcessData(id string) (string, error) {
    data, err := s.store.Get(id)
    if err != nil {
        return "", err
    }
    
    processed := data + "_processed"
    err = s.store.Save(id, processed)
    if err != nil {
        return "", err
    }
    
    return processed, nil
}

// En el código de test
type MockDataStore struct {
    GetFunc  func(id string) (string, error)
    SaveFunc func(id, data string) error
}

func (m *MockDataStore) Get(id string) (string, error) {
    return m.GetFunc(id)
}

func (m *MockDataStore) Save(id, data string) error {
    return m.SaveFunc(id, data)
}

func TestService_ProcessData(t *testing.T) {
    // Configurar el mock
    mock := &MockDataStore{
        GetFunc: func(id string) (string, error) {
            if id == "valid" {
                return "data", nil
            }
            return "", errors.New("not found")
        },
        SaveFunc: func(id, data string) error {
            if id == "valid" && data == "data_processed" {
                return nil
            }
            return errors.New("save failed")
        },
    }
    
    service := &Service{store: mock}
    
    // Caso de éxito
    result, err := service.ProcessData("valid")
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if result != "data_processed" {
        t.Errorf("Expected 'data_processed', got '%s'", result)
    }
    
    // Caso de error
    _, err = service.ProcessData("invalid")
    if err == nil {
        t.Error("Expected error, got nil")
    }
}
```

#### Uso de Bibliotecas de Mocking

Aunque Go promueve el uso de interfaces simples para mocking, existen bibliotecas que pueden facilitar esta tarea:

```go
// Usando github.com/golang/mock/gomock
func TestWithMockGen(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockStore := NewMockDataStore(ctrl) // Generado por mockgen
    
    // Configurar expectativas
    mockStore.EXPECT().Get("valid").Return("data", nil)
    mockStore.EXPECT().Save("valid", "data_processed").Return(nil)
    
    service := &Service{store: mockStore}
    result, err := service.ProcessData("valid")
    
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if result != "data_processed" {
        t.Errorf("Expected 'data_processed', got '%s'", result)
    }
}
```

### Monkey Patching (con precaución)

El monkey patching permite reemplazar funciones en tiempo de ejecución, lo que puede ser útil para probar código que usa funciones globales o estáticas. Sin embargo, debe usarse con precaución ya que puede llevar a tests frágiles:

```go
// Usando github.com/bouk/monkey
func TestMonkeyPatching(t *testing.T) {
    // Guardar la función original para restaurarla después
    original := time.Now
    defer func() { time.Now = original }()
    
    // Reemplazar time.Now con una función mock
    time.Now = func() time.Time {
        return time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    }
    
    // Probar una función que usa time.Now
    result := FormatCurrentTime()
    expected := "2023-01-01"
    
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

## Testing de Integración

### Testing con Bases de Datos Reales

```go
func TestDatabaseIntegration(t *testing.T) {
    // Saltar si no estamos en modo integración
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Configurar base de datos de prueba
    db, err := sql.Open("postgres", "postgres://user:pass@localhost/testdb?sslmode=disable")
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    defer db.Close()
    
    // Limpiar datos antes de la prueba
    _, err = db.Exec("TRUNCATE TABLE users")
    if err != nil {
        t.Fatalf("Failed to truncate table: %v", err)
    }
    
    // Crear repositorio con DB real
    repo := NewUserRepository(db)
    
    // Probar operaciones
    user := User{Name: "Test User", Email: "test@example.com"}
    err = repo.Create(&user)
    if err != nil {
        t.Errorf("Failed to create user: %v", err)
    }
    
    // Verificar que el usuario fue creado
    retrieved, err := repo.GetByID(user.ID)
    if err != nil {
        t.Errorf("Failed to retrieve user: %v", err)
    }
    
    if retrieved.Name != user.Name || retrieved.Email != user.Email {
        t.Errorf("Retrieved user doesn't match: got %+v, want %+v", retrieved, user)
    }
}
```

### Testing con Contenedores

Usando `testcontainers-go` para crear contenedores Docker para tests:

```go
func TestWithContainer(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping container test in short mode")
    }
    
    // Crear contenedor PostgreSQL
    ctx := context.Background()
    req := testcontainers.ContainerRequest{
        Image:        "postgres:13",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_PASSWORD": "password",
            "POSTGRES_USER":     "user",
            "POSTGRES_DB":       "testdb",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections"),
    }
    
    postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("Failed to start container: %v", err)
    }
    defer postgresC.Terminate(ctx)
    
    // Obtener host y puerto mapeado
    host, err := postgresC.Host(ctx)
    if err != nil {
        t.Fatalf("Failed to get container host: %v", err)
    }
    
    port, err := postgresC.MappedPort(ctx, "5432")
    if err != nil {
        t.Fatalf("Failed to get container port: %v", err)
    }
    
    // Construir DSN y conectar
    dsn := fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, port.Port())
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // Ejecutar migraciones si es necesario
    // ...
    
    // Continuar con las pruebas usando la base de datos
    // ...
}
```

### Testing de APIs HTTP

```go
func TestHTTPAPI(t *testing.T) {
    // Configurar servidor de prueba
    handler := setupAPIHandler()
    server := httptest.NewServer(handler)
    defer server.Close()
    
    // Crear cliente HTTP
    client := &http.Client{}
    
    // Caso: Crear usuario
    userData := map[string]string{
        "name":  "Test User",
        "email": "test@example.com",
    }
    userJSON, _ := json.Marshal(userData)
    
    req, err := http.NewRequest("POST", server.URL+"/users", bytes.NewBuffer(userJSON))
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("Request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
    }
    
    // Leer y verificar respuesta
    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }
    
    userID, ok := response["id"].(string)
    if !ok || userID == "" {
        t.Errorf("Invalid or missing user ID in response: %v", response)
    }
    
    // Caso: Obtener usuario creado
    req, _ = http.NewRequest("GET", server.URL+"/users/"+userID, nil)
    resp, err = client.Do(req)
    if err != nil {
        t.Fatalf("Request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
    }
    
    var user map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }
    
    if user["name"] != userData["name"] || user["email"] != userData["email"] {
        t.Errorf("User data mismatch. Expected %v, got %v", userData, user)
    }
}
```

## Testing de Rendimiento

### Benchmarks

Go incluye soporte para benchmarks que miden el rendimiento del código:

```go
func BenchmarkMyFunction(b *testing.B) {
    // Preparación antes del benchmark
    data := prepareTestData()
    
    // Resetear el temporizador para no incluir la preparación
    b.ResetTimer()
    
    // Ejecutar la función b.N veces
    for i := 0; i < b.N; i++ {
        MyFunction(data)
    }
}

// Benchmark con diferentes tamaños de entrada
func BenchmarkMyFunction_Small(b *testing.B) {
    data := prepareTestData(10)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        MyFunction(data)
    }
}

func BenchmarkMyFunction_Medium(b *testing.B) {
    data := prepareTestData(100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        MyFunction(data)
    }
}

func BenchmarkMyFunction_Large(b *testing.B) {
    data := prepareTestData(1000)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        MyFunction(data)
    }
}
```

Ejecutar benchmarks:

```bash
go test -bench=. -benchmem
```

### Profiling

Go proporciona herramientas para perfilar el rendimiento de las aplicaciones:

```go
func TestMain(m *testing.M) {
    // Configurar profiling
    if *cpuProfile != "" {
        f, err := os.Create(*cpuProfile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    
    // Ejecutar tests
    code := m.Run()
    
    // Generar perfil de memoria si se solicita
    if *memProfile != "" {
        f, err := os.Create(*memProfile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.WriteHeapProfile(f)
        f.Close()
    }
    
    os.Exit(code)
}
```

Ejecutar tests con profiling:

```bash
go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=.

# Analizar perfiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

### Benchmarks Comparativos

```go
func BenchmarkAlgorithmA(b *testing.B) {
    for i := 0; i < b.N; i++ {
        AlgorithmA(testData)
    }
}

func BenchmarkAlgorithmB(b *testing.B) {
    for i := 0; i < b.N; i++ {
        AlgorithmB(testData)
    }
}

// Comparar diferentes implementaciones con diferentes tamaños
func BenchmarkCompare(b *testing.B) {
    benchmarks := []struct {
        name      string
        size      int
        algorithm func([]int) int
    }{
        {"AlgorithmA/Small", 10, AlgorithmA},
        {"AlgorithmB/Small", 10, AlgorithmB},
        {"AlgorithmA/Medium", 100, AlgorithmA},
        {"AlgorithmB/Medium", 100, AlgorithmB},
        {"AlgorithmA/Large", 1000, AlgorithmA},
        {"AlgorithmB/Large", 1000, AlgorithmB},
    }
    
    for _, bm := range benchmarks {
        b.Run(bm.name, func(b *testing.B) {
            data := generateData(bm.size)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                bm.algorithm(data)
            }
        })
    }
}
```

## Testing de Concurrencia

### Race Detector

Go incluye un detector de condiciones de carrera que puede identificar problemas de concurrencia:

```bash
go test -race ./...
```

### Testing de Código Concurrente

```go
func TestConcurrentMap(t *testing.T) {
    cm := NewConcurrentMap()
    const iterations = 1000
    const keys = 100
    
    // Función para escribir valores
    writer := func(wg *sync.WaitGroup, id int) {
        defer wg.Done()
        for i := 0; i < iterations; i++ {
            key := fmt.Sprintf("key%d", i%keys)
            value := fmt.Sprintf("value%d-%d", id, i)
            cm.Set(key, value)
        }
    }
    
    // Función para leer valores
    reader := func(wg *sync.WaitGroup, id int) {
        defer wg.Done()
        for i := 0; i < iterations; i++ {
            key := fmt.Sprintf("key%d", i%keys)
            cm.Get(key)
        }
    }
    
    // Iniciar goroutines
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(2)
        go writer(&wg, i)
        go reader(&wg, i)
    }
    
    // Esperar a que todas terminen
    wg.Wait()
    
    // Verificar que el mapa tenga el número correcto de entradas
    if size := cm.Size(); size > keys {
        t.Errorf("Expected at most %d entries, got %d", keys, size)
    }
}
```

### Testing de Timeouts

```go
func TestWithTimeout(t *testing.T) {
    // Crear un contexto con timeout
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()
    
    // Canal para recibir el resultado
    resultCh := make(chan string)
    errCh := make(chan error)
    
    // Ejecutar función en goroutine
    go func() {
        result, err := FunctionThatMightTakeTooLong()
        if err != nil {
            errCh <- err
            return
        }
        resultCh <- result
    }()
    
    // Esperar resultado o timeout
    select {
    case result := <-resultCh:
        // Verificar resultado
        if result != "expected" {
            t.Errorf("Expected 'expected', got '%s'", result)
        }
    case err := <-errCh:
        t.Errorf("Function returned error: %v", err)
    case <-ctx.Done():
        t.Errorf("Function took too long to complete")
    }
}
```

## Fuzzing

Go 1.18 introdujo soporte nativo para fuzzing, que genera automáticamente entradas para encontrar casos extremos y vulnerabilidades:

```go
func FuzzParseJSON(f *testing.F) {
    // Proporcionar casos semilla
    f.Add(`{"name":"John","age":30}`)
    f.Add(`{}`)
    
    // Función de fuzzing
    f.Fuzz(func(t *testing.T, data string) {
        var result map[string]interface{}
        err := json.Unmarshal([]byte(data), &result)
        
        // No verificamos si hay error, solo que no haya pánico
        if err == nil {
            // Verificar que podemos volver a serializar
            _, err = json.Marshal(result)
            if err != nil {
                t.Errorf("Failed to marshal parsed JSON: %v", err)
            }
        }
    })
}
```

Ejecutar fuzzing:

```bash
go test -fuzz=FuzzParseJSON -fuzztime=30s
```

## Cobertura de Código

### Generación de Informes de Cobertura

```bash
# Generar informe de cobertura
go test -coverprofile=coverage.out ./...

# Ver informe en HTML
go tool cover -html=coverage.out

# Ver informe en terminal
go tool cover -func=coverage.out
```

### Configuración de Umbrales de Cobertura

Puedes crear un script para verificar que la cobertura cumpla con un umbral mínimo:

```go
package main

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func main() {
    // Ejecutar tests con cobertura
    cmd := exec.Command("go", "test", "-coverprofile=coverage.out", "./...")
    if err := cmd.Run(); err != nil {
        fmt.Printf("Tests failed: %v\n", err)
        os.Exit(1)
    }
    
    // Obtener porcentaje de cobertura
    out, err := exec.Command("go", "tool", "cover", "-func=coverage.out").Output()
    if err != nil {
        fmt.Printf("Failed to get coverage: %v\n", err)
        os.Exit(1)
    }
    
    // Extraer porcentaje total
    lines := strings.Split(string(out), "\n")
    lastLine := lines[len(lines)-2] // La última línea es vacía
    parts := strings.Fields(lastLine)
    coverageStr := parts[len(parts)-1]
    coverageStr = strings.TrimSuffix(coverageStr, "%")
    
    coverage, err := strconv.ParseFloat(coverageStr, 64)
    if err != nil {
        fmt.Printf("Failed to parse coverage: %v\n", err)
        os.Exit(1)
    }
    
    // Verificar umbral
    threshold := 80.0 // 80%
    if coverage < threshold {
        fmt.Printf("Coverage %.2f%% is below threshold %.2f%%\n", coverage, threshold)
        os.Exit(1)
    }
    
    fmt.Printf("Coverage %.2f%% meets threshold %.2f%%\n", coverage, threshold)
}
```

## Mejores Prácticas

### Organización de Tests

1. **Estructura de Directorios**:
   - Tests en el mismo paquete: `mypackage/myfile_test.go`
   - Tests en paquete separado: `mypackage/myfile_test.go` con `package mypackage_test`
   - Tests de integración: `tests/integration/`

2. **Convenciones de Nomenclatura**:
   - Tests: `TestXxx`
   - Benchmarks: `BenchmarkXxx`
   - Ejemplos: `ExampleXxx`
   - Fuzzing: `FuzzXxx`

### Principios de Testing

1. **Tests Independientes**: Cada test debe poder ejecutarse de forma aislada.

2. **Tests Deterministas**: Los tests deben producir el mismo resultado en cada ejecución.

3. **Tests Rápidos**: Los tests unitarios deben ser rápidos para facilitar la ejecución frecuente.

4. **Tests Legibles**: Los tests deben ser fáciles de entender y mantener.

5. **Tests Completos**: Los tests deben cubrir casos normales, casos extremos y casos de error.

### Patrones de Testing

1. **Given-When-Then**:
   ```go
   func TestSomething(t *testing.T) {
       // Given
       input := "test input"
       expected := "expected output"
       
       // When
       result := ProcessInput(input)
       
       // Then
       if result != expected {
           t.Errorf("Expected %s, got %s", expected, result)
       }
   }
   ```

2. **Arrange-Act-Assert**:
   ```go
   func TestSomething(t *testing.T) {
       // Arrange
       service := NewService()
       input := "test input"
       expected := "expected output"
       
       // Act
       result := service.Process(input)
       
       // Assert
       if result != expected {
           t.Errorf("Expected %s, got %s", expected, result)
       }
   }
   ```

3. **Test Fixtures**:
   ```go
   type fixture struct {
       service *Service
       repo    *MockRepository
   }
   
   func setup() *fixture {
       repo := NewMockRepository()
       service := NewService(repo)
       return &fixture{service, repo}
   }
   
   func TestServiceProcess(t *testing.T) {
       f := setup()
       
       f.repo.On("Get", "123").Return("data", nil)
       
       result, err := f.service.Process("123")
       
       if err != nil {
           t.Errorf("Expected no error, got %v", err)
       }
       if result != "processed data" {
           t.Errorf("Expected 'processed data', got '%s'", result)
       }
       
       f.repo.AssertExpectations(t)
   }
   ```

## Herramientas y Bibliotecas

### Bibliotecas de Aserciones

1. **testify**:
   ```go
   import (
       "testing"
       "github.com/stretchr/testify/assert"
       "github.com/stretchr/testify/require"
   )
   
   func TestWithAssertions(t *testing.T) {
       // assert permite continuar después de una falla
       result := Calculate(10)
       assert.Equal(t, 20, result, "Calculate(10) should return 20")
       
       // require detiene la ejecución después de una falla
       user, err := GetUser(123)
       require.NoError(t, err, "GetUser should not return error")
       assert.Equal(t, "John", user.Name, "User name should be John")
   }
   ```

2. **go-cmp**:
   ```go
   import (
       "testing"
       "github.com/google/go-cmp/cmp"
   )
   
   func TestWithCmp(t *testing.T) {
       got := ComplexStruct{Name: "test", Values: []int{1, 2, 3}}
       want := ComplexStruct{Name: "test", Values: []int{1, 2, 3}}
       
       if diff := cmp.Diff(want, got); diff != "" {
           t.Errorf("mismatch (-want +got):\n%s", diff)
       }
   }
   ```

### Bibliotecas de Mocking

1. **gomock**:
   ```go
   // Generar mocks: mockgen -source=interfaces.go -destination=mock_interfaces.go -package=mocks
   
   import (
       "testing"
       "github.com/golang/mock/gomock"
       "myapp/mocks" // Paquete generado por mockgen
   )
   
   func TestWithGomock(t *testing.T) {
       ctrl := gomock.NewController(t)
       defer ctrl.Finish()
       
       mockRepo := mocks.NewMockRepository(ctrl)
       mockRepo.EXPECT().Get("123").Return("data", nil)
       
       service := NewService(mockRepo)
       result, err := service.Process("123")
       
       if err != nil {
           t.Errorf("Expected no error, got %v", err)
       }
       if result != "processed data" {
           t.Errorf("Expected 'processed data', got '%s'", result)
       }
   }
   ```

2. **testify/mock**:
   ```go
   import (
       "testing"
       "github.com/stretchr/testify/mock"
   )
   
   type MockRepository struct {
       mock.Mock
   }
   
   func (m *MockRepository) Get(id string) (string, error) {
       args := m.Called(id)
       return args.String(0), args.Error(1)
   }
   
   func TestWithTestifyMock(t *testing.T) {
       mockRepo := new(MockRepository)
       mockRepo.On("Get", "123").Return("data", nil)
       
       service := NewService(mockRepo)
       result, err := service.Process("123")
       
       if err != nil {
           t.Errorf("Expected no error, got %v", err)
       }
       if result != "processed data" {
           t.Errorf("Expected 'processed data', got '%s'", result)
       }
       
       mockRepo.AssertExpectations(t)
   }
   ```

### Herramientas de Cobertura y Calidad

1. **gocov**:
   ```bash
   go get github.com/axw/gocov/gocov
   gocov test ./... | gocov report
   ```

2. **golangci-lint**:
   ```bash
   golangci-lint run
   ```

3. **SonarQube**:
   ```bash
   sonar-scanner \
     -Dsonar.projectKey=myproject \
     -Dsonar.sources=. \
     -Dsonar.go.coverage.reportPaths=coverage.out
   ```

## Ejercicios Prácticos

1. **Implementar Tests para una API REST**:
   - Crear una API REST simple con operaciones CRUD.
   - Implementar tests unitarios para cada handler.
   - Implementar tests de integración usando `httptest`.
   - Medir y mejorar la cobertura de código.

2. **Testing de Concurrencia**:
   - Implementar un worker pool concurrente.
   - Escribir tests que verifiquen su correcto funcionamiento bajo carga.
   - Usar el race detector para identificar posibles condiciones de carrera.

3. **Benchmarking y Optimización**:
   - Implementar dos algoritmos diferentes para resolver un problema.
   - Escribir benchmarks para comparar su rendimiento.
   - Optimizar el algoritmo más lento basándose en los resultados del profiling.

4. **Mocking de Dependencias Externas**:
   - Crear un servicio que dependa de una API externa.
   - Implementar mocks para la API externa.
   - Escribir tests que verifiquen el comportamiento del servicio en diferentes escenarios.

5. **Fuzzing de Parsers**:
   - Implementar un parser para un formato específico (JSON, CSV, etc.).
   - Escribir tests de fuzzing para encontrar casos extremos.
   - Corregir los problemas encontrados por el fuzzing.

## Conclusiones

El testing avanzado en Go va más allá de los tests unitarios básicos, abarcando técnicas como mocking, benchmarking, fuzzing y testing de integración. Estas técnicas permiten verificar no solo la corrección funcional del código, sino también su rendimiento, seguridad y comportamiento en escenarios complejos.

Las herramientas y bibliotecas del ecosistema Go facilitan la implementación de estas técnicas avanzadas, permitiendo crear suites de tests completas y efectivas. Al combinar estas técnicas con buenas prácticas de organización y principios de testing, podemos asegurar la calidad y robustez de nuestras aplicaciones Go.

Recuerda que el objetivo final del testing no es solo encontrar errores, sino prevenir su introducción en primer lugar. Un enfoque de testing completo y bien diseñado es una inversión que paga dividendos a lo largo de todo el ciclo de vida del software.

## Referencias

1. Go Testing Package Documentation: https://golang.org/pkg/testing/
2. Go Blog - Fuzzing: https://go.dev/blog/fuzz-beta
3. Go Blog - Profiling Go Programs: https://blog.golang.org/pprof
4. Mitchell Hashimoto. (2017). Advanced Testing with Go. GopherCon.
5. Dave Cheney. (2019). High Performance Go Workshop. https://dave.cheney.net/high-performance-go-workshop/gophercon-2019.html
6. Mat Ryer. (2015). Writing Table Driven Tests in Go. https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
7. Peter Bourgon. (2017). Go best practices, six years in. https://peter.bourgon.org/go-best-practices-2016/
8. Testify Documentation: https://github.com/stretchr/testify
9. GoMock Documentation: https://github.com/golang/mock
10. Testcontainers-go Documentation: https://github.com/testcontainers/testcontainers-go