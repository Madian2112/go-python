# Estructuras de Datos en Go

## Introducción

Las estructuras de datos son formas de organizar y almacenar datos para que puedan ser accedidos y modificados eficientemente. Go proporciona varias estructuras de datos incorporadas, cada una con sus propias características y casos de uso. En esta sección, exploraremos las principales estructuras de datos en Go y cómo utilizarlas efectivamente.

## Arrays

Los arrays en Go son colecciones de elementos del mismo tipo con un tamaño fijo definido en tiempo de compilación.

### Declaración y inicialización

```go
// Declaración con tamaño y tipo
var numeros [5]int

// Inicialización con valores
numeros = [5]int{1, 2, 3, 4, 5}

// Declaración e inicialización en una línea
colores := [3]string{"rojo", "verde", "azul"}

// Inicialización con tamaño inferido
dias := [...]string{"lunes", "martes", "miércoles", "jueves", "viernes"}
// El compilador cuenta los elementos (5)

// Inicialización con índices específicos
valores := [5]int{0: 100, 2: 200, 4: 300}
// Equivale a [100, 0, 200, 0, 300]
```

### Acceso y modificación

```go
numeros := [5]int{1, 2, 3, 4, 5}

// Acceso por índice (comienza en 0)
fmt.Println(numeros[0])  // 1
fmt.Println(numeros[2])  // 3

// Modificación
numeros[1] = 20
fmt.Println(numeros)  // [1 20 3 4 5]

// Obtener longitud
fmt.Println(len(numeros))  // 5
```

### Iteración

```go
numeros := [5]int{1, 2, 3, 4, 5}

// Usando for tradicional
for i := 0; i < len(numeros); i++ {
    fmt.Printf("numeros[%d] = %d\n", i, numeros[i])
}

// Usando range (índice y valor)
for i, valor := range numeros {
    fmt.Printf("numeros[%d] = %d\n", i, valor)
}

// Usando range (solo valor, ignorando índice)
for _, valor := range numeros {
    fmt.Println(valor)
}
```

### Arrays multidimensionales

```go
// Declaración de matriz 3x3
var matriz [3][3]int

// Inicialización
matriz = [3][3]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}

// Acceso a elementos
fmt.Println(matriz[1][2])  // 6 (fila 1, columna 2)

// Iteración
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        fmt.Printf("%d ", matriz[i][j])
    }
    fmt.Println()
}
```

### Características importantes

1. **Tamaño fijo**: El tamaño del array es parte de su tipo y no puede cambiar.
2. **Valor vs. referencia**: Los arrays son valores, no referencias. Cuando se asigna o pasa un array, se copia todo el array.

```go
a := [3]int{1, 2, 3}
b := a  // Copia completa
b[0] = 10
fmt.Println(a)  // [1 2 3] (no cambia)
fmt.Println(b)  // [10 2 3]
```

3. **Comparación**: Los arrays del mismo tipo pueden compararse directamente con `==` y `!=`.

```go
a := [3]int{1, 2, 3}
b := [3]int{1, 2, 3}
c := [3]int{1, 2, 4}

fmt.Println(a == b)  // true
fmt.Println(a == c)  // false
```

## Slices

Los slices son segmentos flexibles y dinámicos de arrays. A diferencia de los arrays, los slices no tienen un tamaño fijo.

### Declaración e inicialización

```go
// Slice vacío
var numeros []int

// Inicialización con valores
numeros = []int{1, 2, 3, 4, 5}

// Declaración e inicialización en una línea
colores := []string{"rojo", "verde", "azul"}

// Creación con make (longitud y capacidad)
// make([]T, longitud, capacidad)
numeros = make([]int, 5)       // longitud 5, capacidad 5
numeros = make([]int, 3, 5)    // longitud 3, capacidad 5

// Slice a partir de un array
arr := [5]int{1, 2, 3, 4, 5}
slice1 := arr[1:4]  // [2, 3, 4] (índices 1, 2, 3)
slice2 := arr[:3]   // [1, 2, 3] (índices 0, 1, 2)
slice3 := arr[2:]   // [3, 4, 5] (índices 2, 3, 4)
slice4 := arr[:]    // [1, 2, 3, 4, 5] (todos los elementos)
```

### Longitud y capacidad

Los slices tienen dos propiedades importantes:
- **Longitud**: Número de elementos en el slice (accesible con `len()`)
- **Capacidad**: Número máximo de elementos que puede contener sin reasignar memoria (accesible con `cap()`)

```go
numeros := make([]int, 3, 5)
fmt.Println(len(numeros))  // 3
fmt.Println(cap(numeros))  // 5

// Slice a partir de un array
arr := [5]int{1, 2, 3, 4, 5}
slice := arr[1:4]
fmt.Println(len(slice))  // 3
fmt.Println(cap(slice))  // 4 (desde el índice 1 hasta el final del array)
```

### Modificación de slices

```go
// Añadir elementos
numeros := []int{1, 2, 3}
numeros = append(números, 4)        // [1, 2, 3, 4]
numeros = append(números, 5, 6, 7)  // [1, 2, 3, 4, 5, 6, 7]

// Añadir un slice a otro slice
otros := []int{8, 9, 10}
numeros = append(números, otros...)  // [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

// Modificar elementos
numeros[2] = 30
fmt.Println(numeros)  // [1, 2, 30, 4, 5, 6, 7, 8, 9, 10]

// Eliminar elementos (no hay función directa)
// Para eliminar el elemento en el índice 2:
i := 2
numeros = append(numeros[:i], numeros[i+1:]...)
fmt.Println(numeros)  // [1, 2, 4, 5, 6, 7, 8, 9, 10]
```

### Slices y referencias

A diferencia de los arrays, los slices son referencias a arrays subyacentes. Cuando se modifica un slice, se modifica el array subyacente, lo que puede afectar a otros slices que comparten el mismo array.

```go
arr := [5]int{1, 2, 3, 4, 5}
slice1 := arr[1:4]  // [2, 3, 4]
slice2 := arr[2:5]  // [3, 4, 5]

slice1[1] = 30      // Modifica arr[2]
fmt.Println(arr)    // [1, 2, 30, 4, 5]
fmt.Println(slice1) // [2, 30, 4]
fmt.Println(slice2) // [30, 4, 5]
```

### Copia de slices

Para crear una copia independiente de un slice, se puede usar la función `copy()`:

```go
src := []int{1, 2, 3, 4, 5}
dst := make([]int, len(src))
copiados := copy(dst, src)  // Devuelve el número de elementos copiados

fmt.Println(dst)      // [1, 2, 3, 4, 5]
fmt.Println(copiados) // 5

// Modificar dst no afecta a src
dst[0] = 10
fmt.Println(src)  // [1, 2, 3, 4, 5]
fmt.Println(dst)  // [10, 2, 3, 4, 5]
```

### Slices multidimensionales

```go
// Declaración
var matriz [][]int

// Inicialización
matriz = [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}

// Creación dinámica
filas, columnas := 3, 4
matriz = make([][]int, filas)
for i := range matriz {
    matriz[i] = make([]int, columnas)
}

// Acceso y modificación
matriz[1][2] = 42
fmt.Println(matriz[1][2])  // 42
```

## Maps

Los maps en Go son colecciones no ordenadas de pares clave-valor, similares a los diccionarios o tablas hash en otros lenguajes.

### Declaración e inicialización

```go
// Map vacío
var edades map[string]int

// Inicialización con make
edades = make(map[string]int)

// Declaración e inicialización en una línea
edades := make(map[string]int)

// Inicialización con valores
edades := map[string]int{
    "Ana":    30,
    "Carlos": 25,
    "Elena":  28,
}

// Map vacío literal
edades := map[string]int{}
```

### Operaciones con maps

```go
// Añadir o actualizar elementos
edades := map[string]int{
    "Ana":    30,
    "Carlos": 25,
}

edades["Elena"] = 28  // Añadir
edades["Ana"] = 31    // Actualizar

// Acceso a elementos
fmt.Println(edades["Carlos"])  // 25

// Verificar existencia
edad, existe := edades["David"]
if existe {
    fmt.Printf("David tiene %d años\n", edad)
} else {
    fmt.Println("David no está en el map")
}

// Eliminar elementos
delete(edades, "Carlos")
fmt.Println(edades)  // map[Ana:31 Elena:28]

// Longitud
fmt.Println(len(edades))  // 2
```

### Iteración

```go
edades := map[string]int{
    "Ana":    30,
    "Carlos": 25,
    "Elena":  28,
}

// Iteración sobre claves y valores
for nombre, edad := range edades {
    fmt.Printf("%s tiene %d años\n", nombre, edad)
}

// Iteración solo sobre claves
for nombre := range edades {
    fmt.Println(nombre)
}
```

### Características importantes

1. **Orden no garantizado**: El orden de iteración en un map no está garantizado y puede cambiar entre ejecuciones.

2. **Referencias**: Los maps son referencias, no valores. Cuando se asigna o pasa un map, se pasa una referencia al mismo map.

```go
m1 := map[string]int{"a": 1, "b": 2}
m2 := m1  // m2 es una referencia al mismo map
m2["a"] = 10
fmt.Println(m1)  // map[a:10 b:2] (m1 también cambia)
```

3. **Comparación**: Los maps no pueden compararse directamente con `==`. Solo pueden compararse con `nil`.

```go
m1 := map[string]int{"a": 1}
fmt.Println(m1 == nil)  // false

var m2 map[string]int
fmt.Println(m2 == nil)  // true

// m1 == m2  // Error de compilación
```

4. **Maps como conjuntos**: Los maps pueden usarse para implementar conjuntos, usando las claves como elementos y los valores como marcadores.

```go
// Implementación de un conjunto
conjunto := map[string]bool{
    "manzana": true,
    "banana":  true,
    "cereza":  true,
}

// Verificar pertenencia
fmt.Println(conjunto["manzana"])  // true
fmt.Println(conjunto["pera"])    // false (valor cero de bool)

// Añadir elemento
conjunto["pera"] = true

// Eliminar elemento
delete(conjunto, "banana")

// Alternativa más eficiente en memoria
conjuntoVacio := map[string]struct{}{
    "manzana": {},
    "banana":  {},
    "cereza":  {},
}
conjuntoVacio["pera"] = struct{}{}
```

## Structs

Las structs en Go son tipos de datos compuestos que agrupan variables (campos) bajo un solo nombre.

### Declaración y definición

```go
// Definición de una struct
type Persona struct {
    Nombre   string
    Edad     int
    Direccion string
    Email    string
}

// Struct anónima (sin definir un tipo)
var empleado struct {
    ID     int
    Nombre  string
    Activo  bool
}
```

### Inicialización

```go
// Inicialización con valores en orden
p1 := Persona{"Ana", 30, "Calle Principal 123", "ana@example.com"}

// Inicialización con nombres de campo (recomendado)
p2 := Persona{
    Nombre:    "Carlos",
    Edad:      25,
    Direccion: "Avenida Central 456",
    Email:     "carlos@example.com",
}

// Inicialización con valores cero
var p3 Persona  // Todos los campos tienen su valor cero

// Inicialización con new (devuelve un puntero)
p4 := new(Persona)  // &Persona{"" 0 "" ""}
```

### Acceso y modificación de campos

```go
p := Persona{"Ana", 30, "Calle Principal 123", "ana@example.com"}

// Acceso a campos
fmt.Println(p.Nombre)  // Ana
fmt.Println(p.Edad)    // 30

// Modificación de campos
p.Edad = 31
p.Email = "ana.garcia@example.com"
fmt.Println(p)  // {Ana 31 Calle Principal 123 ana.garcia@example.com}
```

### Structs anidadas

```go
type Direccion struct {
    Calle  string
    Ciudad string
    CP     string
}

type Persona struct {
    Nombre    string
    Edad      int
    Direccion Direccion  // Struct anidada
}

// Inicialización
p := Persona{
    Nombre: "Ana",
    Edad:   30,
    Direccion: Direccion{
        Calle:  "Calle Principal",
        Ciudad: "Madrid",
        CP:     "28001",
    },
}

// Acceso a campos anidados
fmt.Println(p.Direccion.Ciudad)  // Madrid
```

### Campos promocionados

Go permite "promocionar" campos de structs anidadas para acceder a ellos directamente:

```go
type Persona struct {
    Nombre string
    Edad   int
    Direccion  // Struct anidada sin nombre de campo
}

p := Persona{
    Nombre: "Ana",
    Edad:   30,
    Direccion: Direccion{
        Calle:  "Calle Principal",
        Ciudad: "Madrid",
        CP:     "28001",
    },
}

// Acceso directo a campos promocionados
fmt.Println(p.Ciudad)  // Madrid (en lugar de p.Direccion.Ciudad)
```

### Structs y punteros

```go
p := Persona{"Ana", 30, "Calle Principal 123", "ana@example.com"}

// Crear un puntero a la struct
pp := &p

// Acceso a campos a través del puntero
// Go permite usar la notación de punto directamente con punteros a structs
fmt.Println(pp.Nombre)  // Ana (equivalente a (*pp).Nombre)

// Modificación a través del puntero
pp.Edad = 31
fmt.Println(p.Edad)  // 31 (p también cambia)
```

### Métodos de structs

Go permite asociar métodos a structs mediante receptores:

```go
type Rectangulo struct {
    Ancho, Alto float64
}

// Método con receptor de valor
func (r Rectangulo) Area() float64 {
    return r.Ancho * r.Alto
}

// Método con receptor de puntero
func (r *Rectangulo) Escalar(factor float64) {
    r.Ancho *= factor
    r.Alto *= factor
}

// Uso
r := Rectangulo{10, 5}
fmt.Println(r.Area())  // 50

r.Escalar(2)
fmt.Println(r)         // {20 10}
fmt.Println(r.Area())  // 200
```

### Comparación de structs

Las structs pueden compararse con `==` si todos sus campos son comparables:

```go
type Persona struct {
    Nombre string
    Edad   int
}

p1 := Persona{"Ana", 30}
p2 := Persona{"Ana", 30}
p3 := Persona{"Carlos", 25}

fmt.Println(p1 == p2)  // true
fmt.Println(p1 == p3)  // false
```

Sin embargo, las structs que contienen campos no comparables (como slices o maps) no pueden compararse directamente.

## Interfaces

Las interfaces en Go definen un conjunto de métodos que un tipo debe implementar. Son una forma de lograr polimorfismo.

### Definición e implementación

```go
// Definición de una interfaz
type Forma interface {
    Area() float64
    Perimetro() float64
}

// Implementación para Rectangulo
type Rectangulo struct {
    Ancho, Alto float64
}

func (r Rectangulo) Area() float64 {
    return r.Ancho * r.Alto
}

func (r Rectangulo) Perimetro() float64 {
    return 2 * (r.Ancho + r.Alto)
}

// Implementación para Circulo
type Circulo struct {
    Radio float64
}

func (c Circulo) Area() float64 {
    return math.Pi * c.Radio * c.Radio
}

func (c Circulo) Perimetro() float64 {
    return 2 * math.Pi * c.Radio
}
```

### Uso de interfaces

```go
// Función que acepta cualquier tipo que implemente la interfaz Forma
func ImprimirInfo(f Forma) {
    fmt.Printf("Área: %.2f, Perímetro: %.2f\n", f.Area(), f.Perimetro())
}

// Uso
r := Rectangulo{10, 5}
c := Circulo{7}

ImprimirInfo(r)  // Área: 50.00, Perímetro: 30.00
ImprimirInfo(c)  // Área: 153.94, Perímetro: 43.98

// Slice de interfaces
formas := []Forma{
    Rectangulo{10, 5},
    Circulo{7},
    Rectangulo{3, 4},
}

for _, forma := range formas {
    ImprimirInfo(forma)
}
```

### Interfaces vacías

La interfaz vacía `interface{}` (o `any` en Go 1.18+) no tiene métodos y puede contener valores de cualquier tipo:

```go
// Función que acepta cualquier tipo
func Imprimir(v interface{}) {
    fmt.Printf("Valor: %v, Tipo: %T\n", v, v)
}

// Uso
Imprimir(42)            // Valor: 42, Tipo: int
Imprimir("hola")        // Valor: hola, Tipo: string
Imprimir(Rectangulo{10, 5})  // Valor: {10 5}, Tipo: main.Rectangulo
```

### Type assertions y type switches

Para recuperar el valor concreto de una interfaz, se usan type assertions y type switches:

```go
// Type assertion
func Procesar(v interface{}) {
    // Intentar obtener un string
    str, ok := v.(string)
    if ok {
        fmt.Printf("Es un string: %s\n", str)
    } else {
        fmt.Println("No es un string")
    }
}

// Type switch
func Clasificar(v interface{}) string {
    switch x := v.(type) {
    case nil:
        return "nil"
    case int, int32, int64:
        return "entero"
    case float32, float64:
        return "flotante"
    case string:
        return "string"
    case Rectangulo:
        return "rectángulo"
    case Forma:
        return "forma"
    default:
        return fmt.Sprintf("tipo desconocido: %T", x)
    }
}
```

## Ejemplos Prácticos

### Implementación de una pila (stack)

```go
package main

import (
    "errors"
    "fmt"
)

// Pila es una implementación de una pila (LIFO) usando un slice
type Pila struct {
    elementos []interface{}
}

// Nuevo crea una nueva pila vacía
func NuevaPila() *Pila {
    return &Pila{elementos: make([]interface{}, 0)}
}

// Apilar añade un elemento al tope de la pila
func (p *Pila) Apilar(elemento interface{}) {
    p.elementos = append(p.elementos, elemento)
}

// Desapilar elimina y devuelve el elemento del tope de la pila
func (p *Pila) Desapilar() (interface{}, error) {
    if p.EstaVacia() {
        return nil, errors.New("la pila está vacía")
    }
    
    indice := len(p.elementos) - 1
    elemento := p.elementos[indice]
    p.elementos = p.elementos[:indice]
    return elemento, nil
}

// Tope devuelve el elemento del tope sin eliminarlo
func (p *Pila) Tope() (interface{}, error) {
    if p.EstaVacia() {
        return nil, errors.New("la pila está vacía")
    }
    
    return p.elementos[len(p.elementos)-1], nil
}

// EstaVacia verifica si la pila está vacía
func (p *Pila) EstaVacia() bool {
    return len(p.elementos) == 0
}

// Tamaño devuelve el número de elementos en la pila
func (p *Pila) Tamaño() int {
    return len(p.elementos)
}

func main() {
    pila := NuevaPila()
    
    // Apilar elementos
    pila.Apilar(1)
    pila.Apilar("dos")
    pila.Apilar(3.14)
    
    fmt.Printf("Tamaño de la pila: %d\n", pila.Tamaño())
    
    // Ver el tope
    if tope, err := pila.Tope(); err == nil {
        fmt.Printf("Elemento en el tope: %v\n", tope)
    }
    
    // Desapilar elementos
    for !pila.EstaVacia() {
        elemento, _ := pila.Desapilar()
        fmt.Printf("Desapilado: %v\n", elemento)
    }
    
    // Intentar desapilar de una pila vacía
    if _, err := pila.Desapilar(); err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Implementación de una cola (queue)

```go
package main

import (
    "errors"
    "fmt"
)

// Cola es una implementación de una cola (FIFO) usando un slice
type Cola struct {
    elementos []interface{}
}

// NuevaCola crea una nueva cola vacía
func NuevaCola() *Cola {
    return &Cola{elementos: make([]interface{}, 0)}
}

// Encolar añade un elemento al final de la cola
func (c *Cola) Encolar(elemento interface{}) {
    c.elementos = append(c.elementos, elemento)
}

// Desencolar elimina y devuelve el elemento del frente de la cola
func (c *Cola) Desencolar() (interface{}, error) {
    if c.EstaVacia() {
        return nil, errors.New("la cola está vacía")
    }
    
    elemento := c.elementos[0]
    c.elementos = c.elementos[1:]
    return elemento, nil
}

// Frente devuelve el elemento del frente sin eliminarlo
func (c *Cola) Frente() (interface{}, error) {
    if c.EstaVacia() {
        return nil, errors.New("la cola está vacía")
    }
    
    return c.elementos[0], nil
}

// EstaVacia verifica si la cola está vacía
func (c *Cola) EstaVacia() bool {
    return len(c.elementos) == 0
}

// Tamaño devuelve el número de elementos en la cola
func (c *Cola) Tamaño() int {
    return len(c.elementos)
}

func main() {
    cola := NuevaCola()
    
    // Encolar elementos
    cola.Encolar("primero")
    cola.Encolar(2)
    cola.Encolar(true)
    
    fmt.Printf("Tamaño de la cola: %d\n", cola.Tamaño())
    
    // Ver el frente
    if frente, err := cola.Frente(); err == nil {
        fmt.Printf("Elemento en el frente: %v\n", frente)
    }
    
    // Desencolar elementos
    for !cola.EstaVacia() {
        elemento, _ := cola.Desencolar()
        fmt.Printf("Desencolado: %v\n", elemento)
    }
    
    // Intentar desencolar de una cola vacía
    if _, err := cola.Desencolar(); err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Implementación de un conjunto (set)

```go
package main

import (
    "fmt"
    "strings"
)

// Conjunto implementa un conjunto usando un map
type Conjunto struct {
    elementos map[string]struct{}
}

// NuevoConjunto crea un nuevo conjunto vacío
func NuevoConjunto() *Conjunto {
    return &Conjunto{elementos: make(map[string]struct{})}
}

// NuevoConjuntoDesde crea un conjunto a partir de una lista de elementos
func NuevoConjuntoDesde(elementos ...string) *Conjunto {
    c := NuevoConjunto()
    for _, elemento := range elementos {
        c.Agregar(elemento)
    }
    return c
}

// Agregar añade un elemento al conjunto
func (c *Conjunto) Agregar(elemento string) {
    c.elementos[elemento] = struct{}{}
}

// Eliminar elimina un elemento del conjunto
func (c *Conjunto) Eliminar(elemento string) {
    delete(c.elementos, elemento)
}

// Contiene verifica si un elemento está en el conjunto
func (c *Conjunto) Contiene(elemento string) bool {
    _, existe := c.elementos[elemento]
    return existe
}

// Tamaño devuelve el número de elementos en el conjunto
func (c *Conjunto) Tamaño() int {
    return len(c.elementos)
}

// Elementos devuelve una slice con todos los elementos del conjunto
func (c *Conjunto) Elementos() []string {
    elementos := make([]string, 0, len(c.elementos))
    for elemento := range c.elementos {
        elementos = append(elementos, elemento)
    }
    return elementos
}

// Union devuelve un nuevo conjunto que es la unión de este conjunto con otro
func (c *Conjunto) Union(otro *Conjunto) *Conjunto {
    resultado := NuevoConjunto()
    
    // Añadir elementos de este conjunto
    for elemento := range c.elementos {
        resultado.Agregar(elemento)
    }
    
    // Añadir elementos del otro conjunto
    for elemento := range otro.elementos {
        resultado.Agregar(elemento)
    }
    
    return resultado
}

// Interseccion devuelve un nuevo conjunto que es la intersección de este conjunto con otro
func (c *Conjunto) Interseccion(otro *Conjunto) *Conjunto {
    resultado := NuevoConjunto()
    
    // Añadir elementos que están en ambos conjuntos
    for elemento := range c.elementos {
        if otro.Contiene(elemento) {
            resultado.Agregar(elemento)
        }
    }
    
    return resultado
}

// Diferencia devuelve un nuevo conjunto con los elementos que están en este conjunto pero no en el otro
func (c *Conjunto) Diferencia(otro *Conjunto) *Conjunto {
    resultado := NuevoConjunto()
    
    // Añadir elementos que están en este conjunto pero no en el otro
    for elemento := range c.elementos {
        if !otro.Contiene(elemento) {
            resultado.Agregar(elemento)
        }
    }
    
    return resultado
}

// String devuelve una representación en string del conjunto
func (c *Conjunto) String() string {
    elementos := c.Elementos()
    return "{" + strings.Join(elementos, ", ") + "}"
}

func main() {
    // Crear conjuntos
    frutas1 := NuevoConjuntoDesde("manzana", "banana", "cereza", "dátil")
    frutas2 := NuevoConjuntoDesde("banana", "cereza", "kiwi", "limón")
    
    fmt.Printf("Conjunto 1: %s\n", frutas1)
    fmt.Printf("Conjunto 2: %s\n", frutas2)
    
    // Operaciones con conjuntos
    union := frutas1.Union(frutas2)
    interseccion := frutas1.Interseccion(frutas2)
    diferencia1 := frutas1.Diferencia(frutas2)
    diferencia2 := frutas2.Diferencia(frutas1)
    
    fmt.Printf("Unión: %s\n", union)
    fmt.Printf("Intersección: %s\n", interseccion)
    fmt.Printf("Diferencia (1-2): %s\n", diferencia1)
    fmt.Printf("Diferencia (2-1): %s\n", diferencia2)
    
    // Verificar pertenencia
    fmt.Printf("¿'manzana' está en el conjunto 1? %t\n", frutas1.Contiene("manzana"))
    fmt.Printf("¿'kiwi' está en el conjunto 1? %t\n", frutas1.Contiene("kiwi"))
    
    // Modificar conjunto
    frutas1.Agregar("uva")
    frutas1.Eliminar("dátil")
    fmt.Printf("Conjunto 1 modificado: %s\n", frutas1)
}
```

## Buenas Prácticas

1. **Elegir la estructura adecuada**:
   - Usa arrays cuando el tamaño es fijo y conocido en tiempo de compilación.
   - Usa slices para colecciones dinámicas que pueden crecer o reducirse.
   - Usa maps cuando necesites asociar valores con claves.
   - Usa structs para agrupar datos relacionados.

2. **Capacidad inicial**:
   - Cuando sea posible, inicializa slices y maps con una capacidad aproximada para evitar reasignaciones frecuentes.
   ```go
   // Si sabes que necesitarás aproximadamente 100 elementos
   slice := make([]int, 0, 100)
   mapa := make(map[string]int, 100)
   ```

3. **Copias vs. referencias**:
   - Recuerda que arrays son valores (se copian) mientras que slices y maps son referencias.
   - Si necesitas una copia independiente de un slice, usa `copy()`.

4. **Inmutabilidad**:
   - Para datos inmutables, considera usar arrays o structs sin punteros.
   - Para implementar conjuntos inmutables, puedes usar maps con valores vacíos (`struct{}{}`)

5. **Eficiencia**:
   - Para operaciones frecuentes al principio de un slice, considera usar un slice con índices invertidos o una implementación de cola de doble extremo.
   - Para eliminar elementos de un slice sin preservar el orden, puedes usar el "truco de intercambio con el último elemento":
   ```go
   // Eliminar el elemento en el índice i sin preservar el orden
   slice[i] = slice[len(slice)-1]
   slice = slice[:len(slice)-1]
   ```

6. **Interfaces**:
   - Diseña interfaces pequeñas y específicas.
   - Implementa interfaces implícitamente (sin declarar explícitamente que un tipo implementa una interfaz).

7. **Structs**:
   - Usa nombres de campo descriptivos.
   - Organiza los campos para minimizar el padding y optimizar el uso de memoria.
   - Usa comentarios para documentar el propósito de cada campo.

8. **Métodos**:
   - Usa receptores de puntero cuando necesites modificar el receptor o cuando el receptor sea grande.
   - Usa receptores de valor para tipos pequeños e inmutables.

9. **Manejo de errores**:
   - Verifica siempre los errores devueltos por funciones.
   - Usa type assertions con la forma de dos valores (`valor, ok := x.(Tipo)`) para evitar pánicos.

10. **Concurrencia**:
    - Ten cuidado con el acceso concurrente a maps y slices, ya que no son seguros para concurrencia.
    - Usa `sync.Mutex` o canales para sincronizar el acceso a datos compartidos.

## Recursos Adicionales

- [Especificación del lenguaje Go - Tipos](https://golang.org/ref/spec#Types)
- [Tour of Go - Estructuras de datos](https://tour.golang.org/moretypes/1)
- [Effective Go - Slices](https://golang.org/doc/effective_go#slices)
- [Effective Go - Maps](https://golang.org/doc/effective_go#maps)
- [Go by Example - Arrays](https://gobyexample.com/arrays)
- [Go by Example - Slices](https://gobyexample.com/slices)
- [Go by Example - Maps](https://gobyexample.com/maps)
- [Go by Example - Structs](https://gobyexample.com/structs)
- [Go by Example - Interfaces](https://gobyexample.com/interfaces)
- [Go Slices: usage and internals](https://blog.golang.org/slices-intro)

---

En la siguiente sección, exploraremos la programación orientada a objetos en Go, que nos permitirá crear estructuras de datos personalizadas y organizar nuestro código de manera más modular y reutilizable.