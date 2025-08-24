# Pruebas en Go

## Introducción

Las pruebas son una parte fundamental del desarrollo de software, y Go proporciona un marco de pruebas integrado en su biblioteca estándar. Este marco facilita la escritura de pruebas unitarias, pruebas de rendimiento (benchmarks), ejemplos de código y más. En esta sección, exploraremos cómo escribir y ejecutar diferentes tipos de pruebas en Go, así como las mejores prácticas para garantizar la calidad del código.

## Pruebas Unitarias

Las pruebas unitarias en Go se escriben utilizando el paquete `testing`. Cada función de prueba debe comenzar con la palabra `Test` seguida de un nombre que comience con una letra mayúscula, y debe tomar un único parámetro de tipo `*testing.T`.

### Estructura Básica

```go
package mipack

import "testing"

func TestMiFuncion(t *testing.T) {
    resultado := MiFuncion()
    esperado := "valor esperado"
    
    if resultado != esperado {
        t.Errorf("MiFuncion() = %q, se esperaba %q", resultado, esperado)
    }
}
```

### Organización de Archivos de Prueba

Las pruebas en Go se colocan en archivos con el sufijo `_test.go` en el mismo paquete que el código que están probando. Por ejemplo, si tienes un archivo `calc.go`, sus pruebas estarían en `calc_test.go`.

```
├── calc.go
└── calc_test.go
```

### Ejecutando Pruebas

Para ejecutar todas las pruebas en el directorio actual:

```bash
go test
```

Para ejecutar pruebas con más detalles:

```bash
go test -v
```

Para ejecutar una prueba específica:

```bash
go test -run TestMiFuncion
```

El flag `-run` acepta expresiones regulares, por lo que puedes ejecutar múltiples pruebas que coincidan con un patrón.

### Métodos de Aserción

El paquete `testing` proporciona varios métodos para reportar fallos en las pruebas:

```go
func TestEjemplos(t *testing.T) {
    // Falla la prueba inmediatamente con el mensaje dado
    t.Fatal("Esta prueba ha fallado")
    
    // Reporta un error pero continúa la ejecución de la prueba
    t.Error("Algo salió mal")
    
    // Versiones formateadas de Fatal y Error
    t.Fatalf("Valor %d incorrecto", 42)
    t.Errorf("Se esperaba %q, se obtuvo %q", "a", "b")
    
    // Marca la prueba como fallida pero continúa la ejecución
    t.Fail()
    
    // Marca la prueba como fallida y detiene su ejecución
    t.FailNow()
    
    // Comprueba si la prueba ya ha fallado
    if t.Failed() {
        // Hacer algo si la prueba ha fallado
    }
}
```

### Subpruebas

Go 1.7 introdujo el concepto de subpruebas, que permiten agrupar casos de prueba relacionados:

```go
func TestSuma(t *testing.T) {
    casos := []struct {
        nombre   string
        a, b     int
        esperado int
    }{
        {"positivos", 2, 3, 5},
        {"negativo y positivo", -2, 3, 1},
        {"negativos", -2, -3, -5},
    }
    
    for _, caso := range casos {
        t.Run(caso.nombre, func(t *testing.T) {
            resultado := Suma(caso.a, caso.b)
            if resultado != caso.esperado {
                t.Errorf("Suma(%d, %d) = %d; se esperaba %d", 
                         caso.a, caso.b, resultado, caso.esperado)
            }
        })
    }
}
```

Las subpruebas pueden ejecutarse individualmente:

```bash
go test -run TestSuma/positivos
```

### Configuración y Limpieza

Para realizar configuraciones antes de las pruebas y limpiezas después, puedes usar las funciones `t.Cleanup` (Go 1.14+) o patrones de defer:

```go
func TestConConfiguracion(t *testing.T) {
    // Configuración
    recurso := configurarRecurso()
    
    // Limpieza al finalizar
    t.Cleanup(func() {
        liberarRecurso(recurso)
    })
    
    // Alternativa usando defer
    // defer liberarRecurso(recurso)
    
    // Prueba
    resultado := usarRecurso(recurso)
    if resultado != "esperado" {
        t.Errorf("resultado = %q, se esperaba %q", resultado, "esperado")
    }
}
```

### Pruebas en Paralelo

Puedes marcar pruebas para que se ejecuten en paralelo, lo que es útil para pruebas que no interfieren entre sí:

```go
func TestParalelo1(t *testing.T) {
    t.Parallel()
    // Prueba que puede ejecutarse en paralelo
}

func TestParalelo2(t *testing.T) {
    t.Parallel()
    // Otra prueba que puede ejecutarse en paralelo
}
```

## Pruebas de Rendimiento (Benchmarks)

Los benchmarks miden el rendimiento de tu código. Se definen como funciones que comienzan con `Benchmark` y toman un parámetro `*testing.B`.

```go
func BenchmarkMiFuncion(b *testing.B) {
    // Resetear el temporizador para excluir el tiempo de configuración
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        MiFuncion()
    }
}
```

Para ejecutar benchmarks:

```bash
go test -bench=.
```

Para ejecutar un benchmark específico:

```bash
go test -bench=BenchmarkMiFuncion
```

### Comparando Implementaciones

Los benchmarks son útiles para comparar diferentes implementaciones:

```go
func BenchmarkImplementacion1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Implementacion1()
    }
}

func BenchmarkImplementacion2(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Implementacion2()
    }
}
```

### Benchmarks con Datos de Entrada

Puedes parametrizar benchmarks para probar con diferentes tamaños de entrada:

```go
func BenchmarkOrdenar(b *testing.B) {
    tamaños := []int{100, 1000, 10000}
    
    for _, n := range tamaños {
        b.Run(fmt.Sprintf("tamaño-%d", n), func(b *testing.B) {
            datos := generarDatos(n)
            b.ResetTimer()
            
            for i := 0; i < b.N; i++ {
                // Crear una copia para no medir la ordenación de datos ya ordenados
                datosCopia := make([]int, len(datos))
                copy(datosCopia, datos)
                Ordenar(datosCopia)
            }
        })
    }
}
```

### Medición de Memoria

Puedes medir el uso de memoria en tus benchmarks:

```go
func BenchmarkMemoria(b *testing.B) {
    b.ReportAllocs() // Reportar todas las asignaciones de memoria
    
    for i := 0; i < b.N; i++ {
        _ = FuncionQueAsignaMemoria()
    }
}
```

## Ejemplos

Los ejemplos en Go sirven como documentación ejecutable y pruebas. Se definen como funciones que comienzan con `Example`.

```go
func ExampleSuma() {
    resultado := Suma(1, 2)
    fmt.Println(resultado)
    // Output: 3
}

func ExamplePersona_NombreCompleto() {
    p := Persona{Nombre: "Juan", Apellido: "Pérez"}
    fmt.Println(p.NombreCompleto())
    // Output: Juan Pérez
}
```

Los ejemplos se ejecutan como pruebas y verifican que la salida coincida con el comentario `// Output:`. También aparecen en la documentación generada por `godoc`.

## Pruebas de Tabla

Las pruebas de tabla son un patrón común en Go para probar múltiples casos de entrada y salida:

```go
func TestMultiplicar(t *testing.T) {
    casos := []struct {
        nombre     string
        a, b       int
        esperado   int
        debefallar bool
    }{
        {"positivos", 2, 3, 6, false},
        {"cero", 0, 5, 0, false},
        {"negativos", -2, -3, 6, false},
        {"overflow", math.MaxInt64, 2, 0, true},
    }
    
    for _, caso := range casos {
        t.Run(caso.nombre, func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil && !caso.debefallar {
                    t.Errorf("Multiplicar(%d, %d) causó un panic inesperado: %v", 
                             caso.a, caso.b, r)
                } else if r == nil && caso.debefallar {
                    t.Errorf("Multiplicar(%d, %d) debería haber fallado", 
                             caso.a, caso.b)
                }
            }()
            
            resultado := Multiplicar(caso.a, caso.b)
            if !caso.debefallar && resultado != caso.esperado {
                t.Errorf("Multiplicar(%d, %d) = %d; se esperaba %d", 
                         caso.a, caso.b, resultado, caso.esperado)
            }
        })
    }
}
```

## Mocks y Stubs

Go no incluye un framework de mocking en su biblioteca estándar, pero puedes crear mocks fácilmente usando interfaces:

```go
// Interfaz que queremos mockear
type Almacen interface {
    Guardar(clave string, valor interface{}) error
    Obtener(clave string) (interface{}, error)
}

// Implementación real
type AlmacenRedis struct {
    // campos para conexión a Redis
}

func (a *AlmacenRedis) Guardar(clave string, valor interface{}) error {
    // Implementación real
    return nil
}

func (a *AlmacenRedis) Obtener(clave string) (interface{}, error) {
    // Implementación real
    return nil, nil
}

// Mock para pruebas
type AlmacenMock struct {
    GuardarFn    func(clave string, valor interface{}) error
    ObtenerFn    func(clave string) (interface{}, error)
    GuardarCalls []struct {
        Clave string
        Valor interface{}
    }
    ObtenerCalls []struct {
        Clave string
    }
}

func (m *AlmacenMock) Guardar(clave string, valor interface{}) error {
    m.GuardarCalls = append(m.GuardarCalls, struct {
        Clave string
        Valor interface{}
    }{clave, valor})
    return m.GuardarFn(clave, valor)
}

func (m *AlmacenMock) Obtener(clave string) (interface{}, error) {
    m.ObtenerCalls = append(m.ObtenerCalls, struct {
        Clave string
    }{clave})
    return m.ObtenerFn(clave)
}

// Función que usa el almacén
func GuardarUsuario(a Almacen, usuario Usuario) error {
    return a.Guardar(usuario.ID, usuario)
}

// Prueba usando el mock
func TestGuardarUsuario(t *testing.T) {
    mock := &AlmacenMock{
        GuardarFn: func(clave string, valor interface{}) error {
            return nil // Simular éxito
        },
    }
    
    usuario := Usuario{ID: "123", Nombre: "Juan"}
    err := GuardarUsuario(mock, usuario)
    
    if err != nil {
        t.Errorf("Se esperaba nil, se obtuvo %v", err)
    }
    
    if len(mock.GuardarCalls) != 1 {
        t.Fatalf("Se esperaba 1 llamada a Guardar, se obtuvieron %d", len(mock.GuardarCalls))
    }
    
    if mock.GuardarCalls[0].Clave != "123" {
        t.Errorf("Se esperaba clave '123', se obtuvo '%s'", mock.GuardarCalls[0].Clave)
    }
}
```

## Pruebas de Integración

Las pruebas de integración verifican que diferentes partes de tu sistema funcionen juntas correctamente. En Go, no hay una distinción formal entre pruebas unitarias y de integración, pero puedes usar tags de compilación para separar pruebas que requieren recursos externos.

```go
// +build integration

package mipack

import "testing"

func TestIntegracionConBaseDeDatos(t *testing.T) {
    // Prueba que requiere una base de datos real
}
```

Para ejecutar solo pruebas de integración:

```bash
go test -tags=integration
```

## Cobertura de Código

Go incluye herramientas para medir la cobertura de código de tus pruebas:

```bash
go test -cover
```

Para generar un perfil de cobertura y visualizarlo:

```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Esto abrirá un navegador con una visualización de qué partes de tu código están cubiertas por pruebas.

## Fuzzing

A partir de Go 1.18, el fuzzing está integrado en la herramienta de pruebas estándar. El fuzzing es una técnica que proporciona entradas aleatorias a tu código para encontrar casos extremos y vulnerabilidades.

```go
func FuzzReverse(f *testing.F) {
    testcases := []string{"Hello, world", " ", "!12345"}
    for _, tc := range testcases {
        f.Add(tc) // Proporcionar casos semilla
    }
    
    f.Fuzz(func(t *testing.T, orig string) {
        rev := Reverse(orig)
        doubleRev := Reverse(rev)
        if orig != doubleRev {
            t.Errorf("Reverse(Reverse(%q)) = %q, se esperaba %q", orig, doubleRev, orig)
        }
        if utf8.ValidString(orig) && !utf8.ValidString(rev) {
            t.Errorf("Reverse produjo una cadena UTF-8 inválida %q", rev)
        }
    })
}
```

Para ejecutar pruebas de fuzzing:

```bash
go test -fuzz=FuzzReverse
```

## Herramientas de Terceros

Aunque la biblioteca estándar de Go proporciona herramientas sólidas para pruebas, hay varias bibliotecas de terceros que pueden ser útiles:

### Testify

[Testify](https://github.com/stretchr/testify) proporciona funciones de aserción más expresivas y un framework de mocking:

```go
func TestConTestify(t *testing.T) {
    resultado := Suma(2, 3)
    
    assert.Equal(t, 5, resultado, "Suma(2, 3) debería ser 5")
    assert.NotEqual(t, 0, resultado, "Suma(2, 3) no debería ser 0")
    assert.True(t, resultado > 0, "Suma(2, 3) debería ser positivo")
    
    // Aserciones que detienen la prueba si fallan
    require.Equal(t, 5, resultado, "Suma(2, 3) debería ser 5")
}
```

### GoMock

[GoMock](https://github.com/golang/mock) es un framework de mocking que genera mocks a partir de interfaces:

```go
//go:generate mockgen -source=almacen.go -destination=mock_almacen.go -package=mipack

func TestConGoMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockAlmacen := NewMockAlmacen(ctrl)
    mockAlmacen.EXPECT().Guardar("123", gomock.Any()).Return(nil)
    
    usuario := Usuario{ID: "123", Nombre: "Juan"}
    err := GuardarUsuario(mockAlmacen, usuario)
    
    assert.NoError(t, err)
}
```

### Ginkgo y Gomega

[Ginkgo](https://github.com/onsi/ginkgo) y [Gomega](https://github.com/onsi/gomega) proporcionan un framework de pruebas BDD (Behavior-Driven Development):

```go
var _ = Describe("Calculadora", func() {
    Describe("Suma", func() {
        It("suma correctamente dos números positivos", func() {
            Expect(Suma(2, 3)).To(Equal(5))
        })
        
        It("suma correctamente un número positivo y uno negativo", func() {
            Expect(Suma(2, -3)).To(Equal(-1))
        })
    })
})
```

## Ejemplo Práctico: Pruebas para una API REST

Vamos a ver un ejemplo completo de pruebas para una API REST simple:

```go
// api.go
package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
)

type Usuario struct {
    ID    int    `json:"id"`
    Nombre string `json:"nombre"`
    Email  string `json:"email"`
}

type Almacen interface {
    ObtenerUsuario(id int) (Usuario, error)
    GuardarUsuario(usuario Usuario) error
    ListarUsuarios() ([]Usuario, error)
}

type API struct {
    almacen Almacen
}

func NuevaAPI(almacen Almacen) *API {
    return &API{almacen: almacen}
}

func (a *API) ManejadorUsuario(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        partes := strings.Split(r.URL.Path, "/")
        if len(partes) < 3 {
            // Listar todos los usuarios
            usuarios, err := a.almacen.ListarUsuarios()
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            json.NewEncoder(w).Encode(usuarios)
            return
        }
        
        // Obtener un usuario específico
        id, err := strconv.Atoi(partes[2])
        if err != nil {
            http.Error(w, "ID inválido", http.StatusBadRequest)
            return
        }
        
        usuario, err := a.almacen.ObtenerUsuario(id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(usuario)
        
    case http.MethodPost:
        var usuario Usuario
        if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        if err := a.almacen.GuardarUsuario(usuario); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(usuario)
        
    default:
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}
```

Ahora, las pruebas para esta API:

```go
// api_test.go
package api

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"
)

// Mock del almacén
type AlmacenMock struct {
    usuarios map[int]Usuario
}

func NuevoAlmacenMock() *AlmacenMock {
    return &AlmacenMock{
        usuarios: make(map[int]Usuario),
    }
}

func (a *AlmacenMock) ObtenerUsuario(id int) (Usuario, error) {
    usuario, existe := a.usuarios[id]
    if !existe {
        return Usuario{}, errors.New("usuario no encontrado")
    }
    return usuario, nil
}

func (a *AlmacenMock) GuardarUsuario(usuario Usuario) error {
    a.usuarios[usuario.ID] = usuario
    return nil
}

func (a *AlmacenMock) ListarUsuarios() ([]Usuario, error) {
    usuarios := make([]Usuario, 0, len(a.usuarios))
    for _, u := range a.usuarios {
        usuarios = append(usuarios, u)
    }
    return usuarios, nil
}

func TestObtenerUsuario(t *testing.T) {
    // Crear el mock y la API
    almacen := NuevoAlmacenMock()
    almacen.usuarios[1] = Usuario{ID: 1, Nombre: "Juan", Email: "juan@example.com"}
    api := NuevaAPI(almacen)
    
    // Crear una solicitud HTTP de prueba
    req := httptest.NewRequest("GET", "/usuarios/1", nil)
    w := httptest.NewRecorder()
    
    // Llamar al manejador
    api.ManejadorUsuario(w, req)
    
    // Verificar la respuesta
    if w.Code != http.StatusOK {
        t.Errorf("Se esperaba código de estado %d, se obtuvo %d", http.StatusOK, w.Code)
    }
    
    var usuario Usuario
    if err := json.NewDecoder(w.Body).Decode(&usuario); err != nil {
        t.Fatalf("Error al decodificar respuesta: %v", err)
    }
    
    if usuario.ID != 1 || usuario.Nombre != "Juan" || usuario.Email != "juan@example.com" {
        t.Errorf("Usuario incorrecto: %+v", usuario)
    }
}

func TestObtenerUsuarioNoExistente(t *testing.T) {
    almacen := NuevoAlmacenMock()
    api := NuevaAPI(almacen)
    
    req := httptest.NewRequest("GET", "/usuarios/999", nil)
    w := httptest.NewRecorder()
    
    api.ManejadorUsuario(w, req)
    
    if w.Code != http.StatusNotFound {
        t.Errorf("Se esperaba código de estado %d, se obtuvo %d", http.StatusNotFound, w.Code)
    }
}

func TestGuardarUsuario(t *testing.T) {
    almacen := NuevoAlmacenMock()
    api := NuevaAPI(almacen)
    
    usuario := Usuario{ID: 2, Nombre: "María", Email: "maria@example.com"}
    body, _ := json.Marshal(usuario)
    
    req := httptest.NewRequest("POST", "/usuarios", bytes.NewReader(body))
    w := httptest.NewRecorder()
    
    api.ManejadorUsuario(w, req)
    
    if w.Code != http.StatusCreated {
        t.Errorf("Se esperaba código de estado %d, se obtuvo %d", http.StatusCreated, w.Code)
    }
    
    // Verificar que el usuario se guardó
    guardado, err := almacen.ObtenerUsuario(2)
    if err != nil {
        t.Fatalf("Error al obtener usuario guardado: %v", err)
    }
    
    if guardado.ID != 2 || guardado.Nombre != "María" || guardado.Email != "maria@example.com" {
        t.Errorf("Usuario guardado incorrecto: %+v", guardado)
    }
}

func TestListarUsuarios(t *testing.T) {
    almacen := NuevoAlmacenMock()
    almacen.usuarios[1] = Usuario{ID: 1, Nombre: "Juan", Email: "juan@example.com"}
    almacen.usuarios[2] = Usuario{ID: 2, Nombre: "María", Email: "maria@example.com"}
    api := NuevaAPI(almacen)
    
    req := httptest.NewRequest("GET", "/usuarios", nil)
    w := httptest.NewRecorder()
    
    api.ManejadorUsuario(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Se esperaba código de estado %d, se obtuvo %d", http.StatusOK, w.Code)
    }
    
    var usuarios []Usuario
    if err := json.NewDecoder(w.Body).Decode(&usuarios); err != nil {
        t.Fatalf("Error al decodificar respuesta: %v", err)
    }
    
    if len(usuarios) != 2 {
        t.Errorf("Se esperaban 2 usuarios, se obtuvieron %d", len(usuarios))
    }
}

func TestMetodoNoPermitido(t *testing.T) {
    almacen := NuevoAlmacenMock()
    api := NuevaAPI(almacen)
    
    req := httptest.NewRequest("DELETE", "/usuarios/1", nil)
    w := httptest.NewRecorder()
    
    api.ManejadorUsuario(w, req)
    
    if w.Code != http.StatusMethodNotAllowed {
        t.Errorf("Se esperaba código de estado %d, se obtuvo %d", http.StatusMethodNotAllowed, w.Code)
    }
}

// Prueba de tabla para diferentes escenarios
func TestManejadorUsuario(t *testing.T) {
    casos := []struct {
        nombre      string
        metodo      string
        ruta        string
        cuerpo      interface{}
        configurar  func(*AlmacenMock)
        codigoEsperado int
        verificar   func(*testing.T, *httptest.ResponseRecorder, *AlmacenMock)
    }{
        {
            nombre:      "obtener usuario existente",
            metodo:      http.MethodGet,
            ruta:        "/usuarios/1",
            configurar:  func(a *AlmacenMock) {
                a.usuarios[1] = Usuario{ID: 1, Nombre: "Juan", Email: "juan@example.com"}
            },
            codigoEsperado: http.StatusOK,
            verificar:   func(t *testing.T, w *httptest.ResponseRecorder, a *AlmacenMock) {
                var usuario Usuario
                json.NewDecoder(w.Body).Decode(&usuario)
                if usuario.ID != 1 || usuario.Nombre != "Juan" {
                    t.Errorf("Usuario incorrecto: %+v", usuario)
                }
            },
        },
        {
            nombre:      "obtener usuario no existente",
            metodo:      http.MethodGet,
            ruta:        "/usuarios/999",
            configurar:  func(a *AlmacenMock) {},
            codigoEsperado: http.StatusNotFound,
            verificar:   func(t *testing.T, w *httptest.ResponseRecorder, a *AlmacenMock) {},
        },
        {
            nombre:      "crear usuario",
            metodo:      http.MethodPost,
            ruta:        "/usuarios",
            cuerpo:      Usuario{ID: 3, Nombre: "Pedro", Email: "pedro@example.com"},
            configurar:  func(a *AlmacenMock) {},
            codigoEsperado: http.StatusCreated,
            verificar:   func(t *testing.T, w *httptest.ResponseRecorder, a *AlmacenMock) {
                usuario, err := a.ObtenerUsuario(3)
                if err != nil || usuario.Nombre != "Pedro" {
                    t.Errorf("Usuario no guardado correctamente: %v, %+v", err, usuario)
                }
            },
        },
    }
    
    for _, caso := range casos {
        t.Run(caso.nombre, func(t *testing.T) {
            almacen := NuevoAlmacenMock()
            caso.configurar(almacen)
            api := NuevaAPI(almacen)
            
            var body *bytes.Reader
            if caso.cuerpo != nil {
                datos, _ := json.Marshal(caso.cuerpo)
                body = bytes.NewReader(datos)
            }
            
            req := httptest.NewRequest(caso.metodo, caso.ruta, body)
            w := httptest.NewRecorder()
            
            api.ManejadorUsuario(w, req)
            
            if w.Code != caso.codigoEsperado {
                t.Errorf("Se esperaba código de estado %d, se obtuvo %d", 
                         caso.codigoEsperado, w.Code)
            }
            
            caso.verificar(t, w, almacen)
        })
    }
}
```

## Buenas Prácticas

1. **Mantén las pruebas simples y enfocadas**: Cada prueba debe verificar una sola funcionalidad o comportamiento.

2. **Usa nombres descriptivos**: Los nombres de las funciones de prueba deben describir claramente qué están probando y bajo qué condiciones.

3. **Organiza las pruebas en subpruebas**: Usa `t.Run()` para agrupar pruebas relacionadas y facilitar la ejecución selectiva.

4. **Usa pruebas de tabla**: Para probar múltiples casos de entrada y salida con el mismo código.

5. **Evita dependencias externas**: Usa mocks o stubs para aislar el código que estás probando.

6. **Limpia después de las pruebas**: Usa `t.Cleanup()` o `defer` para asegurarte de que los recursos se liberen correctamente.

7. **Escribe pruebas antes de implementar**: Considera el desarrollo guiado por pruebas (TDD) para clarificar los requisitos y diseñar mejores interfaces.

8. **Mantén las pruebas rápidas**: Las pruebas lentas desalientan su ejecución frecuente. Usa mocks para operaciones costosas.

9. **Verifica los casos límite**: Prueba valores extremos, entradas vacías, errores y otros casos especiales.

10. **Ejecuta las pruebas con frecuencia**: Integra las pruebas en tu flujo de trabajo de desarrollo y en tu pipeline de CI/CD.

11. **Mantén un alto nivel de cobertura**: Apunta a una cobertura de código alta, pero recuerda que la calidad de las pruebas es más importante que la cantidad.

12. **Documenta con ejemplos**: Usa `Example` funciones para proporcionar documentación ejecutable.

## Recursos Adicionales

- [Documentación oficial de testing en Go](https://golang.org/pkg/testing/)
- [Go by Example: Testing](https://gobyexample.com/testing)
- [Go by Example: HTTP Testing](https://gobyexample.com/http-server-testing)
- [Testify - Toolkit for Go testing](https://github.com/stretchr/testify)
- [GoMock - Mocking framework for Go](https://github.com/golang/mock)
- [Ginkgo - BDD Testing Framework for Go](https://github.com/onsi/ginkgo)
- [Go Testing Patterns](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Advanced Testing in Go](https://www.youtube.com/watch?v=8hQG7QlcLBk) (Video)
- [Testing Web Applications in Go](https://markphelps.me/posts/testing-web-apps-in-go/)

---

En la siguiente sección, exploraremos la programación orientada a objetos en Go, incluyendo interfaces, composición, y patrones de diseño comunes.