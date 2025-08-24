# Manejo de Errores en Go

## Introducción

El manejo de errores es una parte fundamental de cualquier programa robusto. Go tiene un enfoque único para el manejo de errores que difiere significativamente de otros lenguajes de programación. En lugar de utilizar excepciones, Go utiliza valores de retorno explícitos para indicar y propagar errores. Este enfoque fomenta que los desarrolladores consideren y manejen los posibles errores de manera explícita, lo que conduce a un código más robusto y predecible.

## Fundamentos del Manejo de Errores en Go

### El Tipo Error

En Go, los errores son valores que implementan la interfaz `error`, que está definida en la biblioteca estándar de la siguiente manera:

```go
type error interface {
    Error() string
}
```

Cualquier tipo que implemente el método `Error()` que devuelve una cadena puede ser utilizado como un error en Go.

### Retorno de Errores

La convención en Go es que las funciones que pueden fallar devuelvan un valor de error como último valor de retorno:

```go
func Abrir(nombre string) (File, error) {
    // ...
}
```

El llamador debe verificar si el error es `nil` para determinar si la operación fue exitosa:

```go
f, err := Abrir("archivo.txt")
if err != nil {
    // Manejar el error
    return err
}
// Continuar con el procesamiento normal
```

### Creación de Errores

Go proporciona varias formas de crear errores:

#### Usando errors.New

La forma más simple de crear un error es usando la función `errors.New()` del paquete `errors`:

```go
package main

import (
    "errors"
    "fmt"
)

func main() {
    err := errors.New("algo salió mal")
    fmt.Println(err) // Imprime: algo salió mal
}
```

#### Usando fmt.Errorf

Para errores con formato, puedes usar `fmt.Errorf()`, que permite incluir valores variables en el mensaje de error:

```go
package main

import (
    "fmt"
)

func main() {
    nombre := "archivo.txt"
    err := fmt.Errorf("no se pudo abrir %s", nombre)
    fmt.Println(err) // Imprime: no se pudo abrir archivo.txt
}
```

#### Creando Tipos de Error Personalizados

Puedes crear tus propios tipos de error implementando la interfaz `error`:

```go
package main

import (
    "fmt"
)

// Definir un tipo de error personalizado
type MiError struct {
    Codigo int
    Mensaje string
}

// Implementar la interfaz error
func (e *MiError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Codigo, e.Mensaje)
}

func main() {
    err := &MiError{Codigo: 404, Mensaje: "Recurso no encontrado"}
    fmt.Println(err) // Imprime: Error 404: Recurso no encontrado
}
```

### Verificación de Errores

La forma estándar de verificar errores en Go es:

```go
resultado, err := algunaFuncion()
if err != nil {
    // Manejar el error
    return err // O manejarlo de otra manera
}
// Continuar con el procesamiento normal usando resultado
```

Este patrón es tan común en Go que se considera idiomático.

## Patrones de Manejo de Errores

### Propagación de Errores

Un patrón común es propagar los errores hacia arriba en la pila de llamadas:

```go
func ProcesarArchivo(nombre string) error {
    datos, err := LeerArchivo(nombre)
    if err != nil {
        return err // Propagar el error hacia arriba
    }
    
    resultado, err := ProcesarDatos(datos)
    if err != nil {
        return err // Propagar el error hacia arriba
    }
    
    err = GuardarResultado(resultado)
    if err != nil {
        return err // Propagar el error hacia arriba
    }
    
    return nil // Todo salió bien
}
```

### Enriquecimiento de Errores

A menudo es útil añadir contexto a los errores a medida que se propagan:

```go
func ProcesarArchivo(nombre string) error {
    datos, err := LeerArchivo(nombre)
    if err != nil {
        return fmt.Errorf("error al leer archivo %s: %w", nombre, err)
    }
    
    // ...
    
    return nil
}
```

El verbo `%w` en `fmt.Errorf()` (disponible desde Go 1.13) envuelve el error original, permitiendo que se recupere más tarde usando `errors.Unwrap()`.

### Errores Centinela

Los errores centinela son errores predefinidos que se pueden comparar directamente:

```go
package main

import (
    "errors"
    "fmt"
    "io/fs"
    "os"
)

// Definir errores centinela
var (
    ErrNoEncontrado = errors.New("recurso no encontrado")
    ErrPermisoDenegado = errors.New("permiso denegado")
)

func main() {
    // Usar errores centinela
    err := BuscarRecurso("algo")
    if err == ErrNoEncontrado {
        fmt.Println("El recurso no existe")
    } else if err == ErrPermisoDenegado {
        fmt.Println("No tienes permiso para acceder al recurso")
    } else if err != nil {
        fmt.Println("Ocurrió un error desconocido:", err)
    }
    
    // Errores centinela de la biblioteca estándar
    _, err = os.Open("archivo_inexistente.txt")
    if errors.Is(err, fs.ErrNotExist) {
        fmt.Println("El archivo no existe")
    }
}

func BuscarRecurso(id string) error {
    // Simulación
    if id == "secreto" {
        return ErrPermisoDenegado
    }
    return ErrNoEncontrado
}
```

### Comprobación de Tipos de Error

A veces necesitas verificar si un error es de un tipo específico:

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

// Definir un tipo de error personalizado
type NotFoundError struct {
    Nombre string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s no encontrado", e.Nombre)
}

func main() {
    err := BuscarUsuario("usuario123")
    
    // Usando type assertion
    if nfErr, ok := err.(*NotFoundError); ok {
        fmt.Printf("Error de no encontrado: %s\n", nfErr.Nombre)
    }
    
    // Usando errors.As (Go 1.13+)
    var nfErr *NotFoundError
    if errors.As(err, &nfErr) {
        fmt.Printf("Error de no encontrado: %s\n", nfErr.Nombre)
    }
}

func BuscarUsuario(id string) error {
    // Simulación
    return &NotFoundError{Nombre: id}
}
```

## Funciones de Manejo de Errores en Go 1.13+

Go 1.13 introdujo nuevas funciones para el manejo de errores en el paquete `errors`:

### errors.Is

Compara un error con un error centinela, incluso si el error está envuelto:

```go
package main

import (
    "errors"
    "fmt"
    "io/fs"
    "os"
)

func main() {
    // Abrir un archivo que no existe
    _, err := os.Open("archivo_inexistente.txt")
    
    // Envolver el error
    wrappedErr := fmt.Errorf("error al procesar: %w", err)
    
    // Verificar si el error original es fs.ErrNotExist
    if errors.Is(wrappedErr, fs.ErrNotExist) {
        fmt.Println("El archivo no existe")
    }
}
```

### errors.As

Comprueba si un error es de un tipo específico, incluso si está envuelto:

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

func main() {
    // Intentar abrir un archivo
    _, err := os.Open("/ruta/inexistente/archivo.txt")
    
    // Envolver el error
    wrappedErr := fmt.Errorf("error al procesar: %w", err)
    
    // Verificar si el error es de tipo *os.PathError
    var pathErr *os.PathError
    if errors.As(wrappedErr, &pathErr) {
        fmt.Printf("Error de ruta: op=%s, path=%s, err=%v\n", 
                  pathErr.Op, pathErr.Path, pathErr.Err)
    }
}
```

### errors.Unwrap

Desenvuelve un error para obtener el error original:

```go
package main

import (
    "errors"
    "fmt"
)

func main() {
    err1 := errors.New("error original")
    err2 := fmt.Errorf("error envuelto: %w", err1)
    
    // Desenvolver para obtener el error original
    originalErr := errors.Unwrap(err2)
    fmt.Println(originalErr) // Imprime: error original
    
    // Si intentamos desenvolver un error que no está envuelto, obtenemos nil
    fmt.Println(errors.Unwrap(err1)) // Imprime: <nil>
}
```

## Manejo de Múltiples Errores

### Combinación de Errores

A veces necesitas combinar múltiples errores en uno solo. Puedes crear tu propia implementación o usar paquetes de terceros como `github.com/hashicorp/go-multierror` o `github.com/pkg/errors`:

```go
package main

import (
    "fmt"
    "strings"
)

// Un tipo simple para combinar errores
type MultiError struct {
    Errores []error
}

func (m *MultiError) Error() string {
    mensajes := make([]string, len(m.Errores))
    for i, err := range m.Errores {
        mensajes[i] = err.Error()
    }
    return strings.Join(mensajes, "; ")
}

func (m *MultiError) Add(err error) {
    if err != nil {
        m.Errores = append(m.Errores, err)
    }
}

func (m *MultiError) HasErrors() bool {
    return len(m.Errores) > 0
}

func main() {
    var errores MultiError
    
    // Simular algunas operaciones que pueden fallar
    err1 := ProcesarItem("item1")
    errores.Add(err1)
    
    err2 := ProcesarItem("item2")
    errores.Add(err2)
    
    err3 := ProcesarItem("item3")
    errores.Add(err3)
    
    // Verificar si hubo errores
    if errores.HasErrors() {
        fmt.Println("Ocurrieron errores:")
        fmt.Println(errores.Error())
    } else {
        fmt.Println("Todo salió bien")
    }
}

func ProcesarItem(id string) error {
    // Simulación: fallar para item2
    if id == "item2" {
        return fmt.Errorf("error al procesar %s", id)
    }
    return nil
}
```

## Panic y Recover

Aunque Go favorece el manejo de errores mediante valores de retorno, también proporciona mecanismos para situaciones excepcionales: `panic` y `recover`.

### Panic

`panic` es similar a lanzar una excepción en otros lenguajes. Detiene la ejecución normal de la función actual y comienza a desenrollar la pila, ejecutando cualquier función `defer` en el camino:

```go
package main

import (
    "fmt"
)

func main() {
    fmt.Println("Inicio")
    funcionQuePuedeEntrarEnPanico()
    fmt.Println("Fin") // Esta línea no se ejecutará si hay un panic
}

func funcionQuePuedeEntrarEnPanico() {
    defer fmt.Println("Esto se ejecutará incluso si hay un panic")
    
    // Causar un panic
    panic("¡Algo muy malo ocurrió!")
    
    // El código después de un panic no se ejecuta
    fmt.Println("Esta línea nunca se ejecutará")
}
```

### Recover

`recover` permite a un programa capturar un `panic` y continuar la ejecución normal. Solo funciona cuando se llama desde dentro de una función `defer`:

```go
package main

import (
    "fmt"
)

func main() {
    fmt.Println("Inicio")
    funcionQuePuedeEntrarEnPanico()
    fmt.Println("Fin") // Esta línea se ejecutará porque recuperamos del panic
}

func funcionQuePuedeEntrarEnPanico() {
    // Configurar una función defer para recuperarse de un posible panic
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recuperado del panic:", r)
            // Aquí podrías hacer limpieza o registrar el error
        }
    }()
    
    // Causar un panic
    panic("¡Algo muy malo ocurrió!")
    
    // El código después de un panic no se ejecuta
    fmt.Println("Esta línea nunca se ejecutará")
}
```

### Cuándo Usar Panic y Recover

En Go, `panic` y `recover` no se utilizan para el manejo normal de errores. Se reservan para situaciones verdaderamente excepcionales:

1. **Errores de inicialización**: Cuando un programa no puede iniciar correctamente.
2. **Errores irrecuperables**: Situaciones donde no tiene sentido continuar.
3. **Errores de programación**: Como acceder a un índice fuera de rango o desreferenciar un puntero nulo.

La biblioteca estándar de Go utiliza `panic` en situaciones donde un error indica un bug en el programa, no una condición de error esperada.

## Buenas Prácticas para el Manejo de Errores

### 1. Verificar Todos los Errores

No ignores los errores. Siempre verifica si `err != nil` y maneja el error adecuadamente:

```go
archivo, err := os.Open("archivo.txt")
if err != nil {
    // Manejar el error
    return err
}
// Usar archivo...
```

### 2. Proporcionar Contexto

Añade contexto a los errores para facilitar la depuración:

```go
// Mal
return err

// Bien
return fmt.Errorf("error al procesar archivo %s: %w", nombreArchivo, err)
```

### 3. Usar Errores Centinela con Moderación

Los errores centinela son útiles, pero no abuses de ellos. Úsalos solo para condiciones de error que los llamadores necesitan distinguir:

```go
// Definir solo los errores que realmente necesitan ser distinguidos
var (
    ErrNoEncontrado = errors.New("recurso no encontrado")
    ErrPermisoDenegado = errors.New("permiso denegado")
)
```

### 4. Preferir Tipos de Error para Información Adicional

Cuando necesites incluir información adicional en un error, usa tipos de error personalizados en lugar de simplemente formatear un mensaje:

```go
// Mejor que un simple mensaje de error
type NotFoundError struct {
    Recurso string
    ID      string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s con ID %s no encontrado", e.Recurso, e.ID)
}
```

### 5. Manejar Errores Solo Una Vez

Evita manejar el mismo error múltiples veces. Decide si vas a registrar el error, devolverlo, o ambos, pero sé consistente:

```go
// Mal: manejar el mismo error múltiples veces
if err != nil {
    log.Printf("Error: %v", err) // Registrar
    fmt.Fprintf(os.Stderr, "Error: %v\n", err) // Imprimir en stderr
    return err // Devolver
}

// Mejor: decidir cómo manejar el error
if err != nil {
    log.Printf("Error: %v", err) // Registrar para depuración
    return fmt.Errorf("operación fallida: %w", err) // Devolver con contexto
}
```

### 6. Usar `defer` para Limpieza

Utiliza `defer` para asegurarte de que los recursos se limpien correctamente, incluso cuando ocurren errores:

```go
archivo, err := os.Open("archivo.txt")
if err != nil {
    return err
}
defer archivo.Close() // Se ejecutará al salir de la función

// Resto del código...
```

### 7. Evitar Panic en Código de Biblioteca

En código de biblioteca, evita usar `panic`. Las bibliotecas deben devolver errores para que los llamadores puedan decidir cómo manejarlos:

```go
// Mal (para una biblioteca)
func ProcesarDatos(datos []byte) []byte {
    if len(datos) == 0 {
        panic("datos vacíos")
    }
    // ...
}

// Bien
func ProcesarDatos(datos []byte) ([]byte, error) {
    if len(datos) == 0 {
        return nil, errors.New("datos vacíos")
    }
    // ...
}
```

### 8. Usar Errores Envueltos (Go 1.13+)

Utiliza `fmt.Errorf()` con `%w` para envolver errores y mantener la cadena de errores:

```go
func ProcesarArchivo(nombre string) error {
    datos, err := LeerArchivo(nombre)
    if err != nil {
        return fmt.Errorf("error al leer %s: %w", nombre, err)
    }
    // ...
}
```

### 9. Ser Consistente en el Estilo de Manejo de Errores

Sé consistente en cómo manejas los errores en todo tu código. Esto hace que el código sea más predecible y fácil de mantener.

## Ejemplo Práctico: Validador de Datos

Vamos a crear un ejemplo práctico que demuestra varias técnicas de manejo de errores en Go:

```go
package main

import (
    "errors"
    "fmt"
    "strings"
)

// Definir tipos de error personalizados
type ValidationError struct {
    Campo   string
    Mensaje string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("error de validación en campo '%s': %s", e.Campo, e.Mensaje)
}

type MultiValidationError struct {
    Errores []*ValidationError
}

func (e *MultiValidationError) Error() string {
    if len(e.Errores) == 1 {
        return e.Errores[0].Error()
    }
    
    mensajes := make([]string, len(e.Errores))
    for i, err := range e.Errores {
        mensajes[i] = err.Error()
    }
    return fmt.Sprintf("%d errores de validación: %s", 
                       len(e.Errores), strings.Join(mensajes, "; "))
}

func (e *MultiValidationError) Add(campo, mensaje string) {
    e.Errores = append(e.Errores, &ValidationError{Campo: campo, Mensaje: mensaje})
}

func (e *MultiValidationError) HasErrors() bool {
    return len(e.Errores) > 0
}

// Estructura de datos a validar
type Usuario struct {
    ID        string
    Nombre    string
    Email     string
    Edad      int
    Contraseña string
}

// Función de validación
func ValidarUsuario(usuario *Usuario) error {
    var errores MultiValidationError
    
    // Validar ID
    if usuario.ID == "" {
        errores.Add("ID", "no puede estar vacío")
    } else if len(usuario.ID) < 3 {
        errores.Add("ID", "debe tener al menos 3 caracteres")
    }
    
    // Validar Nombre
    if usuario.Nombre == "" {
        errores.Add("Nombre", "no puede estar vacío")
    }
    
    // Validar Email
    if usuario.Email == "" {
        errores.Add("Email", "no puede estar vacío")
    } else if !strings.Contains(usuario.Email, "@") {
        errores.Add("Email", "formato inválido")
    }
    
    // Validar Edad
    if usuario.Edad < 18 {
        errores.Add("Edad", "debe ser mayor de 18 años")
    } else if usuario.Edad > 120 {
        errores.Add("Edad", "valor no válido")
    }
    
    // Validar Contraseña
    if len(usuario.Contraseña) < 8 {
        errores.Add("Contraseña", "debe tener al menos 8 caracteres")
    }
    
    // Devolver errores o nil
    if errores.HasErrors() {
        return &errores
    }
    return nil
}

// Función que usa la validación
func CrearUsuario(usuario *Usuario) error {
    // Validar usuario
    err := ValidarUsuario(usuario)
    if err != nil {
        return fmt.Errorf("error al crear usuario: %w", err)
    }
    
    // Simular verificación de ID único
    if usuario.ID == "admin" {
        return errors.New("el ID 'admin' está reservado")
    }
    
    // Simular creación exitosa
    fmt.Printf("Usuario creado: %s (%s)\n", usuario.Nombre, usuario.Email)
    return nil
}

func main() {
    // Caso 1: Usuario válido
    usuario1 := &Usuario{
        ID:        "user123",
        Nombre:    "Juan Pérez",
        Email:     "juan@ejemplo.com",
        Edad:      30,
        Contraseña: "contraseña123",
    }
    
    err := CrearUsuario(usuario1)
    if err != nil {
        fmt.Println("Error:", err)
    }
    
    // Caso 2: Usuario con múltiples errores de validación
    usuario2 := &Usuario{
        ID:        "u",
        Nombre:    "",
        Email:     "correo-invalido",
        Edad:      15,
        Contraseña: "123",
    }
    
    err = CrearUsuario(usuario2)
    if err != nil {
        fmt.Println("Error:", err)
        
        // Verificar si es un error de validación múltiple
        var multiErr *MultiValidationError
        if errors.As(err, &multiErr) {
            fmt.Printf("Se encontraron %d errores de validación\n", len(multiErr.Errores))
            for i, valErr := range multiErr.Errores {
                fmt.Printf("  %d. Campo: %s, Problema: %s\n", 
                          i+1, valErr.Campo, valErr.Mensaje)
            }
        }
    }
    
    // Caso 3: ID reservado
    usuario3 := &Usuario{
        ID:        "admin",
        Nombre:    "Administrador",
        Email:     "admin@sistema.com",
        Edad:      35,
        Contraseña: "admin12345",
    }
    
    err = CrearUsuario(usuario3)
    if err != nil {
        fmt.Println("Error:", err)
        
        // Verificar si contiene un mensaje específico
        if strings.Contains(err.Error(), "reservado") {
            fmt.Println("Por favor, elija un ID diferente")
        }
    }
}
```

Este ejemplo demuestra:

1. Tipos de error personalizados (`ValidationError` y `MultiValidationError`)
2. Acumulación de múltiples errores
3. Enriquecimiento de errores con `fmt.Errorf` y `%w`
4. Verificación de tipos de error con `errors.As`
5. Manejo de diferentes tipos de errores

## Recursos Adicionales

- [Documentación del paquete errors](https://golang.org/pkg/errors/)
- [Effective Go - Errors](https://golang.org/doc/effective_go#errors)
- [Go Blog - Error Handling and Go](https://blog.golang.org/error-handling-and-go)
- [Go Blog - Working with Errors in Go 1.13](https://blog.golang.org/go1.13-errors)
- [Dave Cheney - Errors are values](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
- [Rob Pike - Errors are values](https://blog.golang.org/errors-are-values)

---

En la siguiente sección, exploraremos la concurrencia en Go, que es una de las características más poderosas y distintivas del lenguaje.