# Programación Genérica en Go

## Introducción

La programación genérica es una característica introducida oficialmente en Go 1.18 que permite escribir código que funciona con diferentes tipos de datos sin sacrificar la seguridad de tipos. Antes de Go 1.18, los desarrolladores tenían que recurrir a interfaces vacías (`interface{}`) y aserciones de tipo, lo que sacrificaba la seguridad de tipos en tiempo de compilación. Con la introducción de los genéricos, Go ofrece una forma más segura y elegante de escribir código reutilizable.

En este módulo, exploraremos en profundidad la programación genérica en Go, sus casos de uso, patrones comunes y mejores prácticas.

## Fundamentos de la Programación Genérica

### Sintaxis Básica

La sintaxis básica para definir funciones y tipos genéricos en Go utiliza parámetros de tipo entre corchetes `[]`:

```go
// Función genérica
func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// Tipo genérico
type Pair[T, U any] struct {
    First  T
    Second U
}
```

### Restricciones de Tipo

Las restricciones de tipo definen qué tipos pueden utilizarse como argumentos de tipo. Go proporciona algunas restricciones predefinidas en el paquete `constraints`:

```go
import "golang.org/x/exp/constraints"

// Función que acepta cualquier tipo numérico
func Sum[T constraints.Integer | constraints.Float](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}
```

### Definiendo Restricciones Personalizadas

Puedes definir tus propias restricciones utilizando interfaces:

```go
// Restricción personalizada para tipos que pueden ser convertidos a string
type Stringer interface {
    String() string
}

// Restricción que combina múltiples requisitos
type Number interface {
    constraints.Integer | constraints.Float
}

type Printable interface {
    Stringer
    fmt.Formatter
}
```

### Restricción `any`

La restricción `any` es un alias para `interface{}` y permite cualquier tipo:

```go
func PrintAny[T any](value T) {
    fmt.Println(value)
}
```

### Restricción `comparable`

La restricción `comparable` permite tipos que pueden ser comparados con `==` y `!=`:

```go
func Contains[T comparable](slice []T, value T) bool {
    for _, item := range slice {
        if item == value {
            return true
        }
    }
    return false
}
```

## Estructuras de Datos Genéricas

### Implementación de una Pila (Stack) Genérica

```go
package stack

// Stack es una implementación genérica de una pila
type Stack[T any] struct {
    items []T
}

// NewStack crea una nueva pila vacía
func NewStack[T any]() *Stack[T] {
    return &Stack[T]{items: make([]T, 0)}
}

// Push añade un elemento a la pila
func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

// Pop elimina y devuelve el elemento superior de la pila
func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    
    index := len(s.items) - 1
    item := s.items[index]
    s.items = s.items[:index]
    return item, true
}

// Peek devuelve el elemento superior sin eliminarlo
func (s *Stack[T]) Peek() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    
    return s.items[len(s.items)-1], true
}

// IsEmpty verifica si la pila está vacía
func (s *Stack[T]) IsEmpty() bool {
    return len(s.items) == 0
}

// Size devuelve el número de elementos en la pila
func (s *Stack[T]) Size() int {
    return len(s.items)
}
```

### Implementación de un Árbol Binario de Búsqueda Genérico

```go
package bst

import "golang.org/x/exp/constraints"

// Node representa un nodo en el árbol binario de búsqueda
type Node[T constraints.Ordered] struct {
    Value T
    Left  *Node[T]
    Right *Node[T]
}

// BST es un árbol binario de búsqueda genérico
type BST[T constraints.Ordered] struct {
    Root *Node[T]
}

// NewBST crea un nuevo árbol binario de búsqueda vacío
func NewBST[T constraints.Ordered]() *BST[T] {
    return &BST[T]{Root: nil}
}

// Insert añade un valor al árbol
func (bst *BST[T]) Insert(value T) {
    bst.Root = insert(bst.Root, value)
}

// insert es una función auxiliar recursiva para insertar un valor
func insert[T constraints.Ordered](node *Node[T], value T) *Node[T] {
    if node == nil {
        return &Node[T]{Value: value}
    }
    
    if value < node.Value {
        node.Left = insert(node.Left, value)
    } else if value > node.Value {
        node.Right = insert(node.Right, value)
    }
    
    return node
}

// Search busca un valor en el árbol
func (bst *BST[T]) Search(value T) bool {
    return search(bst.Root, value)
}

// search es una función auxiliar recursiva para buscar un valor
func search[T constraints.Ordered](node *Node[T], value T) bool {
    if node == nil {
        return false
    }
    
    if value == node.Value {
        return true
    }
    
    if value < node.Value {
        return search(node.Left, value)
    }
    
    return search(node.Right, value)
}

// InOrderTraversal realiza un recorrido en orden del árbol
func (bst *BST[T]) InOrderTraversal(visit func(T)) {
    inOrder(bst.Root, visit)
}

// inOrder es una función auxiliar recursiva para el recorrido en orden
func inOrder[T constraints.Ordered](node *Node[T], visit func(T)) {
    if node != nil {
        inOrder(node.Left, visit)
        visit(node.Value)
        inOrder(node.Right, visit)
    }
}
```

## Patrones Avanzados con Genéricos

### Funciones de Orden Superior

```go
// Map aplica una función a cada elemento de un slice y devuelve un nuevo slice
func Map[T, U any](slice []T, f func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = f(v)
    }
    return result
}

// Filter filtra elementos de un slice según un predicado
func Filter[T any](slice []T, predicate func(T) bool) []T {
    var result []T
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce aplica una función de reducción a un slice
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
    result := initial
    for _, v := range slice {
        result = reducer(result, v)
    }
    return result
}

// Ejemplo de uso
func Example() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Duplicar cada número
    doubled := Map(numbers, func(n int) int {
        return n * 2
    })
    
    // Filtrar números pares
    evens := Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    
    // Sumar todos los números
    sum := Reduce(numbers, 0, func(acc, n int) int {
        return acc + n
    })
    
    fmt.Println(doubled) // [2 4 6 8 10]
    fmt.Println(evens)    // [2 4]
    fmt.Println(sum)      // 15
}
```

### Opcionales Genéricos

```go
// Optional representa un valor que puede estar presente o ausente
type Optional[T any] struct {
    value T
    present bool
}

// NewOptional crea un nuevo Optional con un valor presente
func NewOptional[T any](value T) Optional[T] {
    return Optional[T]{value: value, present: true}
}

// Empty crea un Optional vacío
func Empty[T any]() Optional[T] {
    var zero T
    return Optional[T]{value: zero, present: false}
}

// IsPresent verifica si el valor está presente
func (o Optional[T]) IsPresent() bool {
    return o.present
}

// Get devuelve el valor si está presente, o panic si no lo está
func (o Optional[T]) Get() T {
    if !o.present {
        panic("Optional: valor no presente")
    }
    return o.value
}

// OrElse devuelve el valor si está presente, o el valor alternativo si no lo está
func (o Optional[T]) OrElse(other T) T {
    if o.present {
        return o.value
    }
    return other
}

// OrElseGet devuelve el valor si está presente, o el resultado de la función si no lo está
func (o Optional[T]) OrElseGet(supplier func() T) T {
    if o.present {
        return o.value
    }
    return supplier()
}

// IfPresent ejecuta una función si el valor está presente
func (o Optional[T]) IfPresent(consumer func(T)) {
    if o.present {
        consumer(o.value)
    }
}

// Map transforma el valor si está presente
func (o Optional[T]) Map[U any](mapper func(T) U) Optional[U] {
    if !o.present {
        return Empty[U]()
    }
    return NewOptional(mapper(o.value))
}
```

### Result Genérico para Manejo de Errores

```go
// Result representa el resultado de una operación que puede fallar
type Result[T any] struct {
    value T
    err   error
}

// Success crea un Result exitoso con un valor
func Success[T any](value T) Result[T] {
    return Result[T]{value: value, err: nil}
}

// Failure crea un Result fallido con un error
func Failure[T any](err error) Result[T] {
    var zero T
    return Result[T]{value: zero, err: err}
}

// IsSuccess verifica si el resultado es exitoso
func (r Result[T]) IsSuccess() bool {
    return r.err == nil
}

// Value devuelve el valor si el resultado es exitoso, o panic si no lo es
func (r Result[T]) Value() T {
    if r.err != nil {
        panic(fmt.Sprintf("Result: intentando acceder al valor de un resultado fallido: %v", r.err))
    }
    return r.value
}

// Error devuelve el error si el resultado es fallido, o nil si es exitoso
func (r Result[T]) Error() error {
    return r.err
}

// Map transforma el valor si el resultado es exitoso
func (r Result[T]) Map[U any](mapper func(T) U) Result[U] {
    if r.err != nil {
        return Failure[U](r.err)
    }
    return Success(mapper(r.value))
}

// FlatMap transforma el valor si el resultado es exitoso y devuelve un nuevo Result
func (r Result[T]) FlatMap[U any](mapper func(T) Result[U]) Result[U] {
    if r.err != nil {
        return Failure[U](r.err)
    }
    return mapper(r.value)
}

// Ejemplo de uso
func divide(a, b int) Result[int] {
    if b == 0 {
        return Failure[int](errors.New("división por cero"))
    }
    return Success(a / b)
}

func Example() {
    result := divide(10, 2)
    if result.IsSuccess() {
        fmt.Println("Resultado:", result.Value())
    } else {
        fmt.Println("Error:", result.Error())
    }
    
    // Encadenamiento de operaciones
    result2 := divide(10, 2).
        Map(func(v int) int { return v * 2 }).
        FlatMap(func(v int) Result[int] { return divide(v, 0) })
    
    if result2.IsSuccess() {
        fmt.Println("Resultado:", result2.Value())
    } else {
        fmt.Println("Error:", result2.Error())
    }
}
```

## Rendimiento y Consideraciones

### Monomorización vs. Diccionarios

Go utiliza un enfoque híbrido para implementar genéricos:

1. **Monomorización**: Para tipos simples como `int`, `string`, etc., el compilador genera código específico para cada instanciación de tipo.
2. **Diccionarios**: Para tipos más complejos, Go utiliza diccionarios de tipos que contienen información sobre cómo operar con un tipo específico.

Esto puede tener implicaciones de rendimiento y tamaño del binario:

```go
// Puede generar código específico para cada tipo
func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// Uso con diferentes tipos
func Example() {
    fmt.Println(Max(10, 20))       // int
    fmt.Println(Max(10.5, 20.5))    // float64
    fmt.Println(Max("abc", "def"))  // string
}
```

### Impacto en el Tiempo de Compilación

El uso extensivo de genéricos puede aumentar el tiempo de compilación debido a la necesidad de realizar más análisis de tipos y generación de código.

### Impacto en el Tamaño del Binario

La monomorización puede aumentar el tamaño del binario, ya que se genera código específico para cada instanciación de tipo.

### Impacto en el Rendimiento en Tiempo de Ejecución

En general, el código genérico bien diseñado debería tener un rendimiento similar al código no genérico equivalente, pero hay casos donde puede haber diferencias:

```go
package main

import (
    "fmt"
    "time"
    "golang.org/x/exp/constraints"
)

// Versión genérica
func SumGeneric[T constraints.Integer | constraints.Float](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}

// Versión específica para int
func SumInt(values []int) int {
    var sum int
    for _, v := range values {
        sum += v
    }
    return sum
}

func main() {
    // Crear un slice grande para la prueba
    const size = 10_000_000
    ints := make([]int, size)
    for i := range ints {
        ints[i] = i
    }
    
    // Medir tiempo para la versión genérica
    start := time.Now()
    resultGeneric := SumGeneric(ints)
    durationGeneric := time.Since(start)
    
    // Medir tiempo para la versión específica
    start = time.Now()
    resultInt := SumInt(ints)
    durationInt := time.Since(start)
    
    fmt.Printf("Resultado genérico: %d, Tiempo: %v\n", resultGeneric, durationGeneric)
    fmt.Printf("Resultado específico: %d, Tiempo: %v\n", resultInt, durationInt)
    fmt.Printf("Diferencia: %.2f%%\n", float64(durationGeneric)/float64(durationInt)*100-100)
}
```

## Mejores Prácticas

### Cuándo Usar Genéricos

1. **Estructuras de datos reutilizables**: Colecciones, contenedores, etc.
2. **Algoritmos genéricos**: Ordenamiento, búsqueda, etc.
3. **Funciones de utilidad**: Map, Filter, Reduce, etc.
4. **Patrones de diseño**: Optional, Result, etc.

### Cuándo Evitar Genéricos

1. **Cuando la interfaz es suficiente**: Si solo necesitas un comportamiento específico, una interfaz puede ser más apropiada.
2. **Código simple**: Para código simple que no necesita reutilización con diferentes tipos.
3. **Rendimiento crítico**: En casos donde el rendimiento es crítico y necesitas optimizaciones específicas para un tipo.

### Diseño de API Genéricas

1. **Mantén las restricciones simples**: Usa restricciones predefinidas cuando sea posible.
2. **Documenta las restricciones**: Explica qué tipos son compatibles con tus funciones genéricas.
3. **Proporciona ejemplos de uso**: Muestra cómo usar tus funciones genéricas con diferentes tipos.
4. **Considera la ergonomía**: Diseña APIs que sean fáciles de usar y entender.

### Nombrado de Parámetros de Tipo

1. **Usa nombres descriptivos para parámetros de tipo complejos**:

```go
// Bueno: nombres descriptivos
type Graph[Node comparable, Edge any] struct {
    // ...
}

// Menos claro: nombres genéricos
type Graph[T comparable, U any] struct {
    // ...
}
```

2. **Usa nombres cortos para parámetros de tipo simples**:

```go
// Bueno: nombres cortos para tipos simples
func Map[T, U any](slice []T, f func(T) U) []U {
    // ...
}
```

### Pruebas de Código Genérico

1. **Prueba con múltiples tipos**: Asegúrate de probar tu código genérico con diferentes tipos para verificar que funciona correctamente.
2. **Prueba los límites de las restricciones**: Verifica que tu código funciona con todos los tipos permitidos por tus restricciones.
3. **Prueba casos extremos**: Prueba con valores extremos, como slices vacíos, valores nulos, etc.

```go
func TestStack(t *testing.T) {
    // Probar con int
    stackInt := NewStack[int]()
    stackInt.Push(1)
    stackInt.Push(2)
    value, ok := stackInt.Pop()
    if !ok || value != 2 {
        t.Errorf("Expected 2, got %v, ok: %v", value, ok)
    }
    
    // Probar con string
    stackString := NewStack[string]()
    stackString.Push("hello")
    stackString.Push("world")
    value2, ok := stackString.Pop()
    if !ok || value2 != "world" {
        t.Errorf("Expected 'world', got %v, ok: %v", value2, ok)
    }
    
    // Probar con struct personalizado
    type Person struct {
        Name string
        Age  int
    }
    stackPerson := NewStack[Person]()
    stackPerson.Push(Person{Name: "Alice", Age: 30})
    stackPerson.Push(Person{Name: "Bob", Age: 25})
    value3, ok := stackPerson.Pop()
    if !ok || value3.Name != "Bob" {
        t.Errorf("Expected Person{Name: 'Bob'}, got %v, ok: %v", value3, ok)
    }
}
```

## Ejercicios Prácticos

1. **Implementa un conjunto (Set) genérico**:
   - Debe permitir añadir, eliminar y verificar la existencia de elementos.
   - Debe funcionar con cualquier tipo comparable.

2. **Implementa una cola (Queue) genérica**:
   - Debe permitir encolar (enqueue) y desencolar (dequeue) elementos.
   - Debe funcionar con cualquier tipo.

3. **Implementa una función `GroupBy` genérica**:
   - Debe agrupar elementos de un slice según una función clave.
   - Debe devolver un mapa donde las claves son los resultados de la función y los valores son slices de elementos.

4. **Implementa un cache genérico con expiración**:
   - Debe permitir almacenar valores con una clave y un tiempo de expiración.
   - Debe eliminar automáticamente los valores expirados.

5. **Implementa un algoritmo de ordenamiento genérico (por ejemplo, QuickSort o MergeSort)**:
   - Debe funcionar con cualquier tipo ordenable.
   - Debe permitir personalizar la función de comparación.

## Conclusión

La programación genérica en Go proporciona una forma poderosa de escribir código reutilizable y seguro en cuanto a tipos. Aunque es una característica relativamente nueva en Go, ya ha demostrado ser valiosa para muchos casos de uso, como estructuras de datos, algoritmos y patrones de diseño.

Al utilizar genéricos de manera efectiva, puedes escribir código más limpio, más seguro y más mantenible, evitando la duplicación y mejorando la seguridad de tipos. Sin embargo, es importante utilizarlos con moderación y considerar las implicaciones de rendimiento y complejidad.

Con la práctica y la experiencia, aprenderás a identificar cuándo los genéricos son la herramienta adecuada para el trabajo y cómo diseñar APIs genéricas efectivas y fáciles de usar.