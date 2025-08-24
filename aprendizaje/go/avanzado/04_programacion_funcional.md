# Programación Funcional en Go

## Introducción

La programación funcional es un paradigma de programación que trata la computación como la evaluación de funciones matemáticas y evita el cambio de estado y los datos mutables. Aunque Go no es un lenguaje puramente funcional, ofrece características que permiten adoptar principios de programación funcional.

En este documento, exploraremos cómo aplicar conceptos de programación funcional en Go, sus beneficios y limitaciones.

## Conceptos Fundamentales

### Funciones como Ciudadanos de Primera Clase

En Go, las funciones son ciudadanos de primera clase, lo que significa que pueden ser:
- Asignadas a variables
- Pasadas como argumentos a otras funciones
- Retornadas desde otras funciones

```go
package main

import "fmt"

func main() {
    // Asignar función a variable
    add := func(a, b int) int {
        return a + b
    }
    
    fmt.Println(add(3, 4)) // 7
    
    // Pasar función como argumento
    result := applyOperation(3, 4, add)
    fmt.Println(result) // 7
    
    // Retornar función desde otra función
    multiply := getMultiplier()
    fmt.Println(multiply(3, 4)) // 12
}

func applyOperation(a, b int, operation func(int, int) int) int {
    return operation(a, b)
}

func getMultiplier() func(int, int) int {
    return func(a, b int) int {
        return a * b
    }
}
```

### Funciones Puras

Una función pura es aquella que:
1. Siempre produce el mismo resultado para los mismos argumentos
2. No tiene efectos secundarios (no modifica variables externas, no realiza I/O, etc.)

```go
// Función pura
func add(a, b int) int {
    return a + b
}

// Función impura (tiene efecto secundario)
func addAndModify(a, b int, result *int) {
    *result = a + b
}
```

### Inmutabilidad

La inmutabilidad es un principio clave en la programación funcional. Aunque Go no tiene soporte nativo para estructuras de datos inmutables, podemos simular inmutabilidad mediante convenciones:

```go
type Point struct {
    X, Y int
}

// En lugar de modificar el punto, creamos uno nuevo
func MovePoint(p Point, deltaX, deltaY int) Point {
    return Point{
        X: p.X + deltaX,
        Y: p.Y + deltaY,
    }
}

func main() {
    p1 := Point{10, 20}
    p2 := MovePoint(p1, 5, 5)
    
    // p1 sigue siendo {10, 20}
    // p2 es {15, 25}
}
```

## Técnicas de Programación Funcional en Go

### Funciones de Orden Superior

Las funciones de orden superior son funciones que toman otras funciones como argumentos o devuelven funciones como resultados.

```go
// Función de orden superior que devuelve una función
func multiply(factor int) func(int) int {
    return func(n int) int {
        return n * factor
    }
}

func main() {
    double := multiply(2)
    triple := multiply(3)
    
    fmt.Println(double(5)) // 10
    fmt.Println(triple(5)) // 15
}
```

### Closures

Los closures son funciones que capturan variables de su entorno léxico.

```go
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

func main() {
    c1 := counter()
    c2 := counter()
    
    fmt.Println(c1()) // 1
    fmt.Println(c1()) // 2
    fmt.Println(c2()) // 1 (c2 tiene su propio estado)
}
```

### Recursión

La recursión es una técnica común en programación funcional. Go soporta recursión, pero no tiene optimización de cola de recursión.

```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

// Versión con acumulador para simular recursión de cola
func factorialTail(n int, acc int) int {
    if n <= 1 {
        return acc
    }
    return factorialTail(n-1, n*acc)
}

func factorialWrapper(n int) int {
    return factorialTail(n, 1)
}
```

### Funciones Map, Filter y Reduce

Estas son operaciones fundamentales en programación funcional. Go no las proporciona en su biblioteca estándar, pero podemos implementarlas:

```go
// Map: aplica una función a cada elemento de un slice
func Map(vs []int, f func(int) int) []int {
    vsm := make([]int, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}

// Filter: filtra elementos según un predicado
func Filter(vs []int, f func(int) bool) []int {
    vsf := make([]int, 0)
    for _, v := range vs {
        if f(v) {
            vsf = append(vsf, v)
        }
    }
    return vsf
}

// Reduce: combina todos los elementos en un único valor
func Reduce(vs []int, f func(int, int) int, initial int) int {
    r := initial
    for _, v := range vs {
        r = f(r, v)
    }
    return r
}

func main() {
    nums := []int{1, 2, 3, 4, 5}
    
    // Duplicar cada número
    doubled := Map(nums, func(n int) int {
        return n * 2
    })
    fmt.Println(doubled) // [2 4 6 8 10]
    
    // Filtrar números pares
    evens := Filter(nums, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println(evens) // [2 4]
    
    // Sumar todos los números
    sum := Reduce(nums, func(acc, n int) int {
        return acc + n
    }, 0)
    fmt.Println(sum) // 15
}
```

### Composición de Funciones

La composición de funciones es una técnica para combinar múltiples funciones en una sola.

```go
func compose(f, g func(int) int) func(int) int {
    return func(x int) int {
        return f(g(x))
    }
}

func main() {
    addOne := func(x int) int { return x + 1 }
    double := func(x int) int { return x * 2 }
    
    // (x + 1) * 2
    doubleAfterAddOne := compose(double, addOne)
    
    // 2 * (x + 1)
    addOneAfterDouble := compose(addOne, double)
    
    fmt.Println(doubleAfterAddOne(3)) // (3 + 1) * 2 = 8
    fmt.Println(addOneAfterDouble(3)) // 2 * 3 + 1 = 7
}
```

## Patrones Funcionales Avanzados

### Option Pattern (Maybe Monad)

El patrón Option es útil para manejar valores que pueden ser nulos.

```go
type Option[T any] struct {
    value T
    valid bool
}

func Some[T any](value T) Option[T] {
    return Option[T]{value: value, valid: true}
}

func None[T any]() Option[T] {
    return Option[T]{valid: false}
}

func (o Option[T]) IsSome() bool {
    return o.valid
}

func (o Option[T]) IsNone() bool {
    return !o.valid
}

func (o Option[T]) Unwrap() (T, bool) {
    return o.value, o.valid
}

func (o Option[T]) Map(f func(T) T) Option[T] {
    if !o.valid {
        return o
    }
    return Some(f(o.value))
}

func (o Option[T]) FlatMap(f func(T) Option[T]) Option[T] {
    if !o.valid {
        return o
    }
    return f(o.value)
}

func main() {
    // Ejemplo de uso
    findUser := func(id int) Option[string] {
        users := map[int]string{
            1: "Alice",
            2: "Bob",
        }
        if name, ok := users[id]; ok {
            return Some(name)
        }
        return None[string]()
    }
    
    // Buscar usuario y transformar
    result := findUser(1).Map(func(name string) string {
        return "User: " + name
    })
    
    if value, ok := result.Unwrap(); ok {
        fmt.Println(value) // "User: Alice"
    } else {
        fmt.Println("User not found")
    }
    
    // Buscar usuario inexistente
    result = findUser(3).Map(func(name string) string {
        return "User: " + name
    })
    
    if value, ok := result.Unwrap(); ok {
        fmt.Println(value)
    } else {
        fmt.Println("User not found") // "User not found"
    }
}
```

### Result Pattern (Either Monad)

El patrón Result es útil para manejar operaciones que pueden fallar.

```go
type Result[T any, E any] struct {
    value T
    err   E
    isOk  bool
}

func Ok[T any, E any](value T) Result[T, E] {
    return Result[T, E]{value: value, isOk: true}
}

func Err[T any, E any](err E) Result[T, E] {
    return Result[T, E]{err: err, isOk: false}
}

func (r Result[T, E]) IsOk() bool {
    return r.isOk
}

func (r Result[T, E]) IsErr() bool {
    return !r.isOk
}

func (r Result[T, E]) Unwrap() (T, E, bool) {
    return r.value, r.err, r.isOk
}

func (r Result[T, E]) Map(f func(T) T) Result[T, E] {
    if !r.isOk {
        return r
    }
    return Ok[T, E](f(r.value))
}

func (r Result[T, E]) FlatMap(f func(T) Result[T, E]) Result[T, E] {
    if !r.isOk {
        return r
    }
    return f(r.value)
}

func main() {
    // Ejemplo de uso
    divide := func(a, b int) Result[int, string] {
        if b == 0 {
            return Err[int, string]("division by zero")
        }
        return Ok[int, string](a / b)
    }
    
    // División exitosa
    result := divide(10, 2).Map(func(n int) int {
        return n * 2
    })
    
    if value, _, ok := result.Unwrap(); ok {
        fmt.Println(value) // 10
    } else {
        fmt.Println("Error:")
    }
    
    // División por cero
    result = divide(10, 0).Map(func(n int) int {
        return n * 2
    })
    
    if value, err, ok := result.Unwrap(); ok {
        fmt.Println(value)
    } else {
        fmt.Println("Error:", err) // "Error: division by zero"
    }
}
```

### Pipeline Pattern

El patrón Pipeline permite encadenar operaciones de manera funcional.

```go
type Pipeline[T any] struct {
    value T
}

func NewPipeline[T any](value T) Pipeline[T] {
    return Pipeline[T]{value: value}
}

func (p Pipeline[T]) Pipe(f func(T) T) Pipeline[T] {
    return Pipeline[T]{value: f(p.value)}
}

func (p Pipeline[T]) Value() T {
    return p.value
}

func main() {
    // Ejemplo de uso
    result := NewPipeline(5).
        Pipe(func(n int) int { return n * 2 }).      // 10
        Pipe(func(n int) int { return n + 1 }).      // 11
        Pipe(func(n int) int { return n * n }).      // 121
        Value()
    
    fmt.Println(result) // 121
}
```

## Limitaciones y Consideraciones

### Rendimiento

Algunas técnicas funcionales pueden tener impacto en el rendimiento:

1. **Creación de objetos**: La inmutabilidad implica crear nuevos objetos en lugar de modificar los existentes.
2. **Recursión**: Go no tiene optimización de cola de recursión, lo que puede llevar a desbordamiento de pila.
3. **Asignación de memoria**: Las closures capturan variables, lo que puede aumentar la presión sobre el recolector de basura.

### Legibilidad y Mantenibilidad

La programación funcional puede mejorar la legibilidad y mantenibilidad del código:

1. **Código más declarativo**: Enfoque en qué hacer, no en cómo hacerlo.
2. **Menos efectos secundarios**: Código más predecible y fácil de probar.
3. **Composición**: Construcción de programas complejos a partir de piezas simples.

Sin embargo, el uso excesivo de técnicas funcionales en Go puede hacer que el código sea menos idiomático y más difícil de entender para desarrolladores acostumbrados al estilo imperativo de Go.

### Equilibrio con el Estilo Idiomático de Go

Go es pragmático y favorece la simplicidad. Al aplicar programación funcional en Go, es importante encontrar un equilibrio:

1. **Usar técnicas funcionales cuando aporten claridad**: Map, Filter, Reduce para operaciones de colecciones.
2. **Evitar abstracciones excesivas**: Go prefiere código explícito sobre abstracciones complejas.
3. **Respetar las convenciones de Go**: Seguir las prácticas recomendadas de la comunidad.

## Ejercicios Prácticos

1. **Implementar una biblioteca de utilidades funcionales**:
   - Crear implementaciones genéricas de Map, Filter, Reduce.
   - Implementar funciones para composición y currying.
   - Añadir documentación y pruebas.

2. **Refactorizar código imperativo a estilo funcional**:
   - Tomar un programa existente con muchos efectos secundarios.
   - Refactorizarlo para usar principios funcionales.
   - Comparar legibilidad, mantenibilidad y rendimiento.

3. **Implementar un procesador de datos en pipeline**:
   - Crear un pipeline para procesar datos de un archivo.
   - Usar composición de funciones para transformar los datos.
   - Implementar manejo de errores funcional.

4. **Crear una mini-biblioteca de estructuras de datos inmutables**:
   - Implementar versiones inmutables de listas, mapas, etc.
   - Proporcionar operaciones funcionales para estas estructuras.
   - Comparar con las estructuras de datos mutables de Go.

5. **Desarrollar un mini-framework para manejo de efectos**:
   - Implementar mónadas para IO, estado, etc.
   - Crear un sistema para componer efectos.
   - Demostrar su uso en una aplicación real.

## Conclusiones

La programación funcional en Go ofrece herramientas valiosas para escribir código más claro, modular y fácil de probar. Aunque Go no es un lenguaje puramente funcional, muchos principios funcionales pueden aplicarse con éxito.

Las técnicas funcionales son especialmente útiles para:

1. **Transformación de datos**: Operaciones en colecciones.
2. **Manejo de errores**: Patrones como Option y Result.
3. **Composición de comportamientos**: Funciones de orden superior y composición.

Sin embargo, es importante encontrar un equilibrio entre los principios funcionales y el estilo idiomático de Go. La clave está en usar técnicas funcionales cuando mejoren la claridad y mantenibilidad del código, sin sacrificar la simplicidad y pragmatismo que caracterizan a Go.

## Referencias

1. Rob Pike. (2012). Go at Google: Language Design in the Service of Software Engineering.
2. Alan A. A. Donovan & Brian W. Kernighan. (2015). The Go Programming Language. Addison-Wesley Professional.
3. Mat Ryer. (2020). Go Programming Blueprints. Packt Publishing.
4. John Carmack. (2013). Functional Programming in C++.
5. Bartosz Milewski. (2018). Category Theory for Programmers.
6. Eric Elliott. (2016). Composing Software: An Exploration of Functional Programming and Object Composition in JavaScript.
7. Scott Wlaschin. (2018). Domain Modeling Made Functional: Tackle Software Complexity with Domain-Driven Design and F#.