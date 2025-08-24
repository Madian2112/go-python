# Módulos y Paquetes en Go

## Introducción

Los módulos y paquetes son fundamentales para organizar y reutilizar código en Go. Permiten dividir programas grandes en componentes más pequeños y manejables, facilitando el mantenimiento y la colaboración. En esta sección, aprenderemos cómo crear, importar y utilizar módulos y paquetes en Go.

## Paquetes

Un paquete en Go es una colección de archivos fuente en el mismo directorio que comparten el mismo nombre de paquete. Los paquetes son la unidad básica de encapsulamiento y reutilización de código en Go.

### Declaración de paquetes

Cada archivo Go debe comenzar con una declaración de paquete que especifica a qué paquete pertenece el archivo:

```go
package nombre_del_paquete
```

Por convención, el nombre del paquete es el mismo que el último elemento del directorio que lo contiene, aunque esto no es obligatorio.

### Tipos de paquetes

1. **Paquete `main`**: Es un paquete especial que define un programa ejecutable. Debe contener una función `main()` que sirve como punto de entrada del programa.

```go
package main

import "fmt"

func main() {
    fmt.Println("Hola, mundo!")
}
```

2. **Paquetes de biblioteca**: Son paquetes diseñados para ser importados y utilizados por otros paquetes. No contienen una función `main()`.

```go
// En el archivo matematicas/operaciones.go
package matematicas

// Sumar devuelve la suma de dos números
func Sumar(a, b int) int {
    return a + b
}
```

### Visibilidad de identificadores

En Go, la visibilidad de un identificador (variable, función, tipo, etc.) fuera de su paquete se determina por la primera letra de su nombre:

- **Exportado (público)**: Si comienza con una letra mayúscula, es accesible desde otros paquetes.
- **No exportado (privado)**: Si comienza con una letra minúscula, solo es accesible dentro del paquete que lo define.

```go
package matematicas

// Sumar es exportado (accesible desde otros paquetes)
func Sumar(a, b int) int {
    return a + b
}

// multiplicar no es exportado (solo accesible dentro del paquete)
func multiplicar(a, b int) int {
    return a * b
}

// PI es una constante exportada
const PI = 3.14159265359

// factorInterno no es exportado
const factorInterno = 0.5
```

## Importación de paquetes

Para utilizar un paquete en Go, debes importarlo usando la declaración `import`:

### Importación simple

```go
import "fmt"

func main() {
    fmt.Println("Hola, mundo!")
}
```

### Importación múltiple

```go
import (
    "fmt"
    "math"
    "strings"
)

func main() {
    fmt.Println(math.Pi)
    fmt.Println(strings.ToUpper("hola"))
}
```

### Importación con alias

Puedes asignar un alias a un paquete importado para evitar conflictos de nombres o para abreviar nombres largos:

```go
import (
    "fmt"
    m "math"
    s "strings"
)

func main() {
    fmt.Println(m.Pi)          // En lugar de math.Pi
    fmt.Println(s.ToUpper("hola"))  // En lugar de strings.ToUpper
}
```

### Importación con punto

Puedes importar un paquete con un punto (`.`) para acceder a sus identificadores exportados sin calificarlos con el nombre del paquete:

```go
import (
    "fmt"
    . "math"
)

func main() {
    fmt.Println(Pi)  // En lugar de math.Pi
    fmt.Println(Sqrt(16))  // En lugar de math.Sqrt
}
```

> **Nota**: Esta forma de importación generalmente no se recomienda, ya que puede hacer que el código sea menos legible y causar conflictos de nombres.

### Importación en blanco

A veces, necesitas importar un paquete solo por sus efectos secundarios (como la inicialización), sin usar directamente ninguno de sus identificadores exportados. En este caso, puedes usar una importación en blanco con el guion bajo (`_`):

```go
import (
    "fmt"
    _ "github.com/lib/pq"  // Driver de PostgreSQL que se registra por sí mismo
)

func main() {
    // El paquete github.com/lib/pq se inicializa pero no se usa directamente
    fmt.Println("Conectando a la base de datos...")
}
```

## Organización de paquetes

### Estructura de directorios típica

Una estructura de directorios típica para un proyecto Go podría verse así:

```
miproyecto/
├── cmd/
│   └── miapp/
│       └── main.go         # Punto de entrada principal
├── internal/               # Código privado para este proyecto
│   ├── config/
│   │   └── config.go
│   └── database/
│       └── database.go
├── pkg/                    # Código que podría ser utilizado por otros proyectos
│   ├── matematicas/
│   │   └── operaciones.go
│   └── utilidades/
│       └── archivos.go
├── api/                    # Definiciones de API, protobuf, etc.
├── web/                    # Recursos web
├── configs/                # Archivos de configuración
├── scripts/                # Scripts para automatización
├── test/                   # Pruebas adicionales
├── docs/                   # Documentación
├── examples/               # Ejemplos de uso
├── third_party/            # Herramientas y código de terceros
├── go.mod                  # Definición del módulo
└── go.sum                  # Sumas de verificación de dependencias
```

### Convenciones de nombres

- Los nombres de paquetes deben ser cortos, concisos y descriptivos.
- Usa minúsculas para los nombres de paquetes.
- Evita guiones bajos o camelCase en los nombres de paquetes.
- El nombre del paquete suele ser el mismo que el último elemento de la ruta de importación.

## Módulos

A partir de Go 1.11, se introdujo el sistema de módulos como la solución oficial para la gestión de dependencias. Un módulo es una colección de paquetes relacionados que se versiona como una unidad.

### Creación de un módulo

Para crear un nuevo módulo, usa el comando `go mod init` seguido del nombre del módulo (generalmente la ruta de importación):

```bash
go mod init github.com/usuario/miproyecto
```

Esto creará un archivo `go.mod` que define el módulo y sus dependencias:

```
module github.com/usuario/miproyecto

go 1.16
```

### Gestión de dependencias

Cuando importas paquetes externos en tu código, Go automáticamente descarga y añade las dependencias a tu módulo cuando ejecutas comandos como `go build`, `go test` o `go mod tidy`.

```go
import (
    "fmt"
    "github.com/usuario/otroproyecto/pkg/utilidades"
)
```

Después de ejecutar `go mod tidy`, el archivo `go.mod` se actualizará con las dependencias:

```
module github.com/usuario/miproyecto

go 1.16

require github.com/usuario/otroproyecto v1.2.3
```

Además, se creará un archivo `go.sum` con las sumas de verificación de las dependencias para garantizar la integridad.

### Versionado de módulos

Go utiliza el versionado semántico (SemVer) para las versiones de los módulos. Las versiones tienen el formato `vMAJOR.MINOR.PATCH` (por ejemplo, `v1.2.3`).

Para publicar una nueva versión de tu módulo, simplemente crea una etiqueta Git con el formato adecuado:

```bash
git tag v1.0.0
git push origin v1.0.0
```

### Actualización de dependencias

Para actualizar una dependencia a la última versión:

```bash
go get -u github.com/usuario/otroproyecto
```

Para actualizar a una versión específica:

```bash
go get github.com/usuario/otroproyecto@v1.3.0
```

### Módulos y compatibilidad hacia atrás

En Go, la compatibilidad hacia atrás es muy importante. Según la especificación de módulos de Go:

- Las versiones con el mismo número MAJOR deben ser compatibles entre sí.
- Cuando se introducen cambios incompatibles, se debe incrementar el número MAJOR.

Para versiones 2 y superiores, el número de versión debe incluirse en la ruta de importación:

```go
import "github.com/usuario/otroproyecto/v2/pkg/utilidades"
```

## Ejemplo Práctico: Creación de un Módulo y Paquetes

Vamos a crear un módulo llamado `calculadora` con varios paquetes:

### Estructura del módulo

```
calculadora/
├── cmd/
│   └── calc/
│       └── main.go
├── pkg/
│   ├── basicas/
│   │   └── operaciones.go
│   ├── avanzadas/
│   │   └── operaciones.go
│   └── utilidades/
│       └── formato.go
├── go.mod
└── README.md
```

### Inicialización del módulo

```bash
mkdir -p calculadora
cd calculadora
go mod init github.com/usuario/calculadora
```

### Implementación de los paquetes

```go
// pkg/basicas/operaciones.go
package basicas

// Sumar devuelve la suma de dos números
func Sumar(a, b float64) float64 {
    return a + b
}

// Restar devuelve la diferencia entre dos números
func Restar(a, b float64) float64 {
    return a - b
}

// Multiplicar devuelve el producto de dos números
func Multiplicar(a, b float64) float64 {
    return a * b
}

// Dividir devuelve el cociente de dos números
// Devuelve un error si el divisor es cero
func Dividir(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("división por cero")
    }
    return a / b, nil
}
```

```go
// pkg/avanzadas/operaciones.go
package avanzadas

import "math"

// Potencia devuelve a elevado a la potencia b
func Potencia(a, b float64) float64 {
    return math.Pow(a, b)
}

// RaizCuadrada devuelve la raíz cuadrada de un número
// Devuelve un error si el número es negativo
func RaizCuadrada(a float64) (float64, error) {
    if a < 0 {
        return 0, fmt.Errorf("no se puede calcular la raíz cuadrada de un número negativo")
    }
    return math.Sqrt(a), nil
}

// Factorial devuelve el factorial de un número entero
func Factorial(n int) int {
    if n <= 0 {
        return 1
    }
    return n * Factorial(n-1)
}
```

```go
// pkg/utilidades/formato.go
package utilidades

import "fmt"

// FormatearResultado formatea un resultado numérico con un mensaje
func FormatearResultado(operacion string, resultado float64) string {
    return fmt.Sprintf("El resultado de %s es: %.2f", operacion, resultado)
}

// FormatearError formatea un mensaje de error
func FormatearError(operacion string, err error) string {
    return fmt.Sprintf("Error en %s: %v", operacion, err)
}
```

```go
// cmd/calc/main.go
package main

import (
    "fmt"
    "os"
    "strconv"

    "github.com/usuario/calculadora/pkg/basicas"
    "github.com/usuario/calculadora/pkg/avanzadas"
    "github.com/usuario/calculadora/pkg/utilidades"
)

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Uso: calc <operación> <num1> <num2>")
        fmt.Println("Operaciones disponibles: sumar, restar, multiplicar, dividir, potencia, raiz")
        os.Exit(1)
    }

    operacion := os.Args[1]
    num1, err := strconv.ParseFloat(os.Args[2], 64)
    if err != nil {
        fmt.Println("Error: el primer número no es válido")
        os.Exit(1)
    }

    num2, err := strconv.ParseFloat(os.Args[3], 64)
    if err != nil {
        fmt.Println("Error: el segundo número no es válido")
        os.Exit(1)
    }

    switch operacion {
    case "sumar":
        resultado := basicas.Sumar(num1, num2)
        fmt.Println(utilidades.FormatearResultado("suma", resultado))
    case "restar":
        resultado := basicas.Restar(num1, num2)
        fmt.Println(utilidades.FormatearResultado("resta", resultado))
    case "multiplicar":
        resultado := basicas.Multiplicar(num1, num2)
        fmt.Println(utilidades.FormatearResultado("multiplicación", resultado))
    case "dividir":
        resultado, err := basicas.Dividir(num1, num2)
        if err != nil {
            fmt.Println(utilidades.FormatearError("división", err))
            os.Exit(1)
        }
        fmt.Println(utilidades.FormatearResultado("división", resultado))
    case "potencia":
        resultado := avanzadas.Potencia(num1, num2)
        fmt.Println(utilidades.FormatearResultado("potencia", resultado))
    case "raiz":
        // En este caso, num2 se ignora
        resultado, err := avanzadas.RaizCuadrada(num1)
        if err != nil {
            fmt.Println(utilidades.FormatearError("raíz cuadrada", err))
            os.Exit(1)
        }
        fmt.Println(utilidades.FormatearResultado("raíz cuadrada", resultado))
    default:
        fmt.Println("Operación no reconocida")
        os.Exit(1)
    }
}
```

### Compilación y ejecución

```bash
go build -o calc ./cmd/calc
./calc sumar 10 5
# El resultado de suma es: 15.00
./calc potencia 2 3
# El resultado de potencia es: 8.00
```

## Patrones Comunes

### Patrón de inicialización

Go ejecuta la función `init()` de cada paquete antes de ejecutar la función `main()`. Puedes tener múltiples funciones `init()` en un paquete, y se ejecutarán en el orden en que aparecen en el archivo.

```go
package database

import "fmt"

var db *Database

func init() {
    fmt.Println("Inicializando la base de datos...")
    db = &Database{}
    // Configuración inicial
}

type Database struct {
    // campos
}

func GetDB() *Database {
    return db
}
```

### Patrón de fábrica

En lugar de exponer directamente los constructores, puedes proporcionar funciones de fábrica para crear instancias de tipos:

```go
package config

type Config struct {
    // campos privados
    dbURL      string
    apiKey     string
    maxRetries int
}

// NewConfig crea una nueva configuración con valores predeterminados
func NewConfig() *Config {
    return &Config{
        dbURL:      "localhost:5432",
        apiKey:     "",
        maxRetries: 3,
    }
}

// SetDBURL establece la URL de la base de datos
func (c *Config) SetDBURL(url string) *Config {
    c.dbURL = url
    return c
}

// SetAPIKey establece la clave de API
func (c *Config) SetAPIKey(key string) *Config {
    c.apiKey = key
    return c
}

// SetMaxRetries establece el número máximo de reintentos
func (c *Config) SetMaxRetries(retries int) *Config {
    c.maxRetries = retries
    return c
}

// DBURL devuelve la URL de la base de datos
func (c *Config) DBURL() string {
    return c.dbURL
}

// APIKey devuelve la clave de API
func (c *Config) APIKey() string {
    return c.apiKey
}

// MaxRetries devuelve el número máximo de reintentos
func (c *Config) MaxRetries() int {
    return c.maxRetries
}
```

### Patrón de opciones funcionales

Este patrón permite configurar objetos de manera flexible y legible:

```go
package server

import "time"

type Server struct {
    addr         string
    port         int
    readTimeout  time.Duration
    writeTimeout time.Duration
    maxConns     int
}

type Option func(*Server)

func WithAddr(addr string) Option {
    return func(s *Server) {
        s.addr = addr
    }
}

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func WithReadTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.readTimeout = timeout
    }
}

func WithWriteTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.writeTimeout = timeout
    }
}

func WithMaxConnections(maxConns int) Option {
    return func(s *Server) {
        s.maxConns = maxConns
    }
}

func NewServer(options ...Option) *Server {
    // Valores predeterminados
    server := &Server{
        addr:         "localhost",
        port:         8080,
        readTimeout:  30 * time.Second,
        writeTimeout: 30 * time.Second,
        maxConns:     1000,
    }

    // Aplicar opciones
    for _, option := range options {
        option(server)
    }

    return server
}

// Uso:
// server := NewServer(
//     WithAddr("0.0.0.0"),
//     WithPort(9000),
//     WithReadTimeout(60 * time.Second),
// )
```

## Buenas Prácticas

1. **Organización clara**: Organiza tu código en paquetes lógicos basados en funcionalidad, no en tipos.

2. **Nombres descriptivos**: Usa nombres claros y descriptivos para tus paquetes. Evita nombres genéricos como `util` o `common`.

3. **Cohesión**: Cada paquete debe tener un propósito claro y cohesivo. Si un paquete hace demasiadas cosas, considera dividirlo.

4. **Visibilidad mínima**: Exporta solo lo que sea necesario. Mantén los detalles de implementación privados.

5. **Evita la circularidad**: Evita las dependencias circulares entre paquetes, ya que pueden causar problemas de diseño y compilación.

6. **Documentación**: Documenta tus paquetes y funciones exportadas con comentarios que comiencen con el nombre del elemento.

7. **Pruebas**: Coloca las pruebas en el mismo paquete que el código que están probando, con el sufijo `_test.go`.

8. **Versionado**: Sigue las reglas de versionado semántico para tus módulos.

9. **Estructura estándar**: Sigue las convenciones de estructura de proyectos de la comunidad Go.

10. **Interfaces pequeñas**: Diseña interfaces pequeñas y específicas que se centren en comportamientos concretos.

## Recursos Adicionales

- [Cómo escribir código Go](https://golang.org/doc/code.html)
- [Effective Go - Paquetes](https://golang.org/doc/effective_go#package-names)
- [Referencia del comando go mod](https://golang.org/ref/mod)
- [Guía de módulos de Go](https://blog.golang.org/using-go-modules)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Modules by Example](https://github.com/go-modules-by-example/index)

---

En la siguiente sección, exploraremos el manejo de errores en Go, que es un aspecto fundamental del lenguaje y sigue un enfoque diferente al de muchos otros lenguajes de programación.