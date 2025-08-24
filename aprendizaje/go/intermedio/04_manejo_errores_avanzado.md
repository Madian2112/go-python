# Manejo Avanzado de Errores en Go

## Introducción

El manejo de errores es una parte fundamental de la programación en Go. A diferencia de muchos otros lenguajes que utilizan excepciones, Go adopta un enfoque explícito donde los errores son valores que se devuelven y se manejan de forma deliberada. En este módulo, exploraremos técnicas avanzadas de manejo de errores en Go, incluyendo la creación de errores personalizados, patrones de manejo de errores, y mejores prácticas para escribir código robusto y mantenible.

## Repaso de Conceptos Básicos

### Errores como Valores

En Go, los errores son valores que implementan la interfaz `error`:

```go
type error interface {
    Error() string
}
```

La forma más común de manejar errores es verificarlos después de llamar a una función:

```go
func main() {
    resultado, err := dividir(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Resultado:", resultado)
}

func dividir(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("división por cero")
    }
    return a / b, nil
}
```

## Creación de Errores Personalizados

### Errores Simples

La forma más básica de crear errores es usando `errors.New()` o `fmt.Errorf()`:

```go
import (
    "errors"
    "fmt"
)

func validarEdad(edad int) error {
    if edad < 0 {
        return errors.New("la edad no puede ser negativa")
    }
    if edad > 150 {
        return fmt.Errorf("edad %d demasiado alta", edad)
    }
    return nil
}
```

### Tipos de Error Personalizados

Para errores más complejos, puedes definir tipos personalizados que implementen la interfaz `error`:

```go
type ValidationError struct {
    Campo   string
    Mensaje string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("error de validación en campo %s: %s", e.Campo, e.Mensaje)
}

func validarUsuario(usuario string, edad int) error {
    if usuario == "" {
        return &ValidationError{Campo: "usuario", Mensaje: "no puede estar vacío"}
    }
    if edad < 0 {
        return &ValidationError{Campo: "edad", Mensaje: "no puede ser negativa"}
    }
    return nil
}

func main() {
    err := validarUsuario("", 25)
    if err != nil {
        fmt.Println(err) // error de validación en campo usuario: no puede estar vacío
        
        // Podemos verificar el tipo de error
        if valErr, ok := err.(*ValidationError); ok {
            fmt.Printf("Campo con error: %s\n", valErr.Campo)
        }
    }
}
```

### Errores con Datos Adicionales

Los errores personalizados pueden contener datos adicionales que ayuden a diagnosticar o manejar el problema:

```go
type QueryError struct {
    Query   string
    Message string
    Err     error // Error subyacente
}

func (e *QueryError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("error en consulta '%s': %s (causa: %v)", e.Query, e.Message, e.Err)
    }
    return fmt.Sprintf("error en consulta '%s': %s", e.Query, e.Message)
}

// Implementación de Unwrap para compatibilidad con errors.Is y errors.As
func (e *QueryError) Unwrap() error {
    return e.Err
}

func ejecutarConsulta(query string) error {
    // Simular un error de base de datos
    dbErr := errors.New("conexión rechazada")
    return &QueryError{
        Query:   query,
        Message: "falló la ejecución",
        Err:     dbErr,
    }
}
```

## Patrones de Manejo de Errores

### Sentinel Errors

Los "sentinel errors" son errores predefinidos que se pueden comparar directamente:

```go
var (
    ErrNotFound   = errors.New("recurso no encontrado")
    ErrPermission = errors.New("permiso denegado")
    ErrTimeout    = errors.New("tiempo de espera agotado")
)

func obtenerRecurso(id string) ([]byte, error) {
    // Simular diferentes errores
    if id == "" {
        return nil, ErrNotFound
    }
    if id == "privado" {
        return nil, ErrPermission
    }
    // ...
    return []byte("datos del recurso"), nil
}

func main() {
    datos, err := obtenerRecurso("")
    if err != nil {
        switch err {
        case ErrNotFound:
            fmt.Println("El recurso no existe")
        case ErrPermission:
            fmt.Println("No tienes permiso para acceder")
        default:
            fmt.Println("Error desconocido:", err)
        }
        return
    }
    fmt.Println("Datos:", string(datos))
}
```

### Error Wrapping (Go 1.13+)

A partir de Go 1.13, se introdujo la capacidad de "envolver" errores para mantener el contexto:

```go
import (
    "errors"
    "fmt"
)

func procesarArchivo(ruta string) error {
    datos, err := leerArchivo(ruta)
    if err != nil {
        return fmt.Errorf("error al procesar archivo %s: %w", ruta, err)
    }
    // Procesar datos...
    return nil
}

func leerArchivo(ruta string) ([]byte, error) {
    // Simular un error
    return nil, errors.New("archivo no encontrado")
}

func main() {
    err := procesarArchivo("/ruta/archivo.txt")
    if err != nil {
        fmt.Println(err) // error al procesar archivo /ruta/archivo.txt: archivo no encontrado
        
        // Podemos verificar si contiene un error específico
        if errors.Is(err, errors.New("archivo no encontrado")) {
            fmt.Println("El archivo no existe")
        }
    }
}
```

### errors.Is y errors.As

Las funciones `errors.Is` y `errors.As` (introducidas en Go 1.13) facilitan la comprobación de errores en cadenas de errores envueltos:

```go
func main() {
    err := procesarDatos("datos.txt")
    if err != nil {
        // Verificar si el error o alguno de sus errores envueltos es ErrNotFound
        if errors.Is(err, ErrNotFound) {
            fmt.Println("El recurso no se encontró")
        }
        
        // Extraer un tipo específico de error de la cadena
        var valErr *ValidationError
        if errors.As(err, &valErr) {
            fmt.Printf("Error de validación en campo: %s\n", valErr.Campo)
        }
    }
}
```

### Manejo de Múltiples Errores

En algunos casos, necesitamos recopilar múltiples errores antes de devolverlos:

```go
type MultiError struct {
    Errors []error
}

func (m *MultiError) Error() string {
    if len(m.Errors) == 0 {
        return ""
    }
    
    errMsgs := make([]string, len(m.Errors))
    for i, err := range m.Errors {
        errMsgs[i] = err.Error()
    }
    
    return fmt.Sprintf("%d errores: %s", len(m.Errors), strings.Join(errMsgs, "; "))
}

func validarFormulario(nombre, email, telefono string) error {
    var multiErr MultiError
    
    if nombre == "" {
        multiErr.Errors = append(multiErr.Errors, errors.New("nombre requerido"))
    }
    
    if email == "" {
        multiErr.Errors = append(multiErr.Errors, errors.New("email requerido"))
    } else if !strings.Contains(email, "@") {
        multiErr.Errors = append(multiErr.Errors, errors.New("email inválido"))
    }
    
    if telefono != "" && len(telefono) < 10 {
        multiErr.Errors = append(multiErr.Errors, errors.New("teléfono demasiado corto"))
    }
    
    if len(multiErr.Errors) > 0 {
        return &multiErr
    }
    
    return nil
}
```

### Uso de Bibliotecas de Terceros

Existen bibliotecas que extienden las capacidades de manejo de errores en Go:

```go
// Ejemplo con github.com/pkg/errors
import (
    "github.com/pkg/errors"
)

func procesarDatos(ruta string) error {
    datos, err := leerArchivo(ruta)
    if err != nil {
        return errors.Wrap(err, "falló la lectura del archivo")
    }
    
    err = analizarDatos(datos)
    if err != nil {
        return errors.Wrap(err, "falló el análisis de datos")
    }
    
    return nil
}

func main() {
    err := procesarDatos("datos.txt")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        fmt.Printf("Stack trace:\n%+v\n", err) // Imprime el stack trace completo
    }
}
```

## Patrones Avanzados

### Errores como Estado

En algunos casos, los errores pueden representar un estado específico más que una condición de fallo:

```go
type Result struct {
    Value int
    Err   error
}

func (r Result) String() string {
    if r.Err != nil {
        return r.Err.Error()
    }
    return fmt.Sprintf("%d", r.Value)
}

func dividir(a, b int) Result {
    if b == 0 {
        return Result{Err: errors.New("división por cero")}
    }
    return Result{Value: a / b}
}

func main() {
    r1 := dividir(10, 2)
    r2 := dividir(10, 0)
    
    fmt.Println("Resultado 1:", r1) // Resultado 1: 5
    fmt.Println("Resultado 2:", r2) // Resultado 2: división por cero
}
```

### Errores con Comportamiento

Los errores pueden implementar interfaces adicionales para proporcionar comportamiento específico:

```go
type temporary interface {
    Temporary() bool
}

type timeoutError struct {
    mensaje string
    timeout bool
}

func (e *timeoutError) Error() string {
    return e.mensaje
}

func (e *timeoutError) Temporary() bool {
    return e.timeout
}

func operacionRed() error {
    // Simular un error temporal
    return &timeoutError{mensaje: "tiempo de espera agotado", timeout: true}
}

func main() {
    err := operacionRed()
    if err != nil {
        if temp, ok := err.(temporary); ok && temp.Temporary() {
            fmt.Println("Error temporal, reintentando...")
            // Lógica de reintento
        } else {
            fmt.Println("Error permanente, abortando")
        }
    }
}
```

### Panic y Recover

Aunque Go favorece el manejo explícito de errores, `panic` y `recover` pueden ser útiles en situaciones específicas:

```go
func procesarDatos(datos []byte) (resultado string, err error) {
    // Usar defer y recover para convertir panics en errores
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic durante el procesamiento: %v", r)
        }
    }()
    
    // Código que podría causar panic
    if len(datos) == 0 {
        panic("datos vacíos")
    }
    
    // Procesamiento normal
    resultado = string(datos)
    return resultado, nil
}

func main() {
    resultado, err := procesarDatos([]byte{})
    if err != nil {
        fmt.Println("Error:", err) // Error: panic durante el procesamiento: datos vacíos
        return
    }
    fmt.Println("Resultado:", resultado)
}
```

## Mejores Prácticas

### 1. Ser Explícito y Descriptivo

Los mensajes de error deben ser claros y proporcionar contexto suficiente:

```go
// Malo
return errors.New("falló")

// Bueno
return fmt.Errorf("falló la conexión a la base de datos %s: %w", dbName, err)
```

### 2. Envolver Errores para Mantener Contexto

```go
func procesarUsuario(id string) error {
    usuario, err := obtenerUsuario(id)
    if err != nil {
        return fmt.Errorf("error al obtener usuario %s: %w", id, err)
    }
    
    err = validarUsuario(usuario)
    if err != nil {
        return fmt.Errorf("error al validar usuario %s: %w", id, err)
    }
    
    return nil
}
```

### 3. Manejar Errores Solo Una Vez

Evita manejar el mismo error múltiples veces:

```go
// Malo
if err != nil {
    log.Printf("Error: %v", err)
    return err
}

// Bueno
if err != nil {
    return fmt.Errorf("falló la operación: %w", err)
}
```

### 4. Usar Tipos de Error Específicos para Casos Comunes

```go
var (
    ErrNotFound = errors.New("recurso no encontrado")
    ErrTimeout  = errors.New("tiempo de espera agotado")
)

func obtenerRecurso(id string) ([]byte, error) {
    // Implementación...
    if recursoNoExiste(id) {
        return nil, ErrNotFound
    }
    // ...
}
```

### 5. Evitar Errores Anidados Demasiado Profundos

```go
// Evitar esto
if err != nil {
    return fmt.Errorf("error A: %w", fmt.Errorf("error B: %w", fmt.Errorf("error C: %w", err)))
}

// Mejor
if err != nil {
    err = fmt.Errorf("error C: %w", err)
    err = fmt.Errorf("error B: %w", err)
    return fmt.Errorf("error A: %w", err)
}
```

### 6. Usar errors.Is y errors.As para Verificar Errores

```go
// En lugar de comparación directa
if err == ErrNotFound { /* ... */ }

// Usar errors.Is para verificar en toda la cadena de errores
if errors.Is(err, ErrNotFound) { /* ... */ }

// Para tipos de error personalizados
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Printf("Campo con error: %s\n", valErr.Campo)
}
```

### 7. Documentar el Comportamiento de Error

```go
// ObtenerUsuario recupera un usuario por su ID.
// Devuelve ErrNotFound si el usuario no existe.
func ObtenerUsuario(id string) (*Usuario, error) {
    // Implementación...
}
```

### 8. Considerar el Nivel de Detalle Apropiado

```go
// Para APIs internas, proporcionar detalles completos
if err != nil {
    return fmt.Errorf("error al conectar a la base de datos %s en %s:%d: %w", 
                      dbName, dbHost, dbPort, err)
}

// Para APIs públicas, limitar la información sensible
if err != nil {
    return errors.New("error de conexión a la base de datos")
}
```

## Ejercicios Prácticos

1. Implementa un tipo de error personalizado para validación de formularios que pueda contener múltiples errores de campo.

2. Crea una función que utilice `errors.Is` y `errors.As` para manejar diferentes tipos de errores en una cadena de errores envueltos.

3. Implementa un patrón de reintentos para operaciones que pueden fallar temporalmente, utilizando un error que implemente una interfaz `Temporary()`.

4. Diseña un sistema de logging de errores que capture el contexto completo de los errores, incluyendo la cadena de llamadas.

5. Refactoriza un código existente para mejorar su manejo de errores siguiendo las mejores prácticas.

## Conclusión

El manejo de errores en Go es explícito y deliberado, lo que fomenta un código más robusto y predecible. Aunque puede parecer verboso al principio, este enfoque ayuda a evitar errores silenciosos y facilita el diagnóstico de problemas.

Las técnicas avanzadas como errores personalizados, error wrapping, y patrones como `errors.Is` y `errors.As` proporcionan herramientas poderosas para crear sistemas con un manejo de errores sofisticado y mantenible.

Recuerda que el objetivo del manejo de errores no es solo informar que algo salió mal, sino proporcionar información suficiente para entender qué falló, por qué falló, y potencialmente cómo solucionarlo.