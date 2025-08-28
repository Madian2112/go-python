# Variables y Tipos de Datos en Go

## Variables

En Go, las variables son espacios de memoria que almacenan valores. A diferencia de lenguajes dinámicos como Python, Go es un lenguaje de tipado estático, lo que significa que cada variable tiene un tipo específico que se determina en tiempo de compilación.

### Declaración de variables

Hay varias formas de declarar variables en Go:

#### 1. Declaración con la palabra clave `var`

```go
// Declaración básica
var nombre string
nombre = "Gopher"

// Declaración con inicialización
var edad int = 25

// Múltiples variables del mismo tipo
var x, y, z int

// Múltiples variables con inicialización
var ancho, alto int = 100, 200

// Múltiples variables de diferentes tipos
var (
    nombre string = "Gopher"
    edad   int    = 25
    activo bool   = true
)
```

#### 2. Declaración corta con operador `:=`

Esta forma solo funciona dentro de funciones, no a nivel de paquete:

```go
// El tipo se infiere del valor asignado
nombre := "Gopher"
edad := 25
activo := true

// Múltiples variables en una línea
nombre, edad := "Gopher", 25
```

### Reglas para nombrar variables

- Deben comenzar con una letra o un guion bajo (_)
- Pueden contener letras, números y guiones bajos
- Son sensibles a mayúsculas y minúsculas
- No pueden ser palabras reservadas (como `if`, `for`, `func`, etc.)

### Convenciones de nomenclatura

- Usa camelCase para variables locales: `nombreUsuario`
- Usa PascalCase para variables exportadas (visibles fuera del paquete): `NombreUsuario`
- Usa nombres descriptivos
- Las constantes suelen usar PascalCase o todo mayúsculas: `PI` o `MaxConexiones`

### Visibilidad de variables

En Go, la visibilidad se determina por la primera letra del nombre:

- Si comienza con mayúscula, es exportada (visible fuera del paquete)
- Si comienza con minúscula, es privada (solo visible dentro del paquete)

```go
var Nombre string // Exportada, visible fuera del paquete
var edad int      // Privada, solo visible dentro del paquete
```

## Tipos de datos básicos

Go tiene varios tipos de datos incorporados:

### Números enteros

#### Enteros con signo

```go
var a int    // Tamaño dependiente de la plataforma (32 o 64 bits)
var b int8   // 8 bits (-128 a 127)
var c int16  // 16 bits (-32,768 a 32,767)
var d int32  // 32 bits (-2,147,483,648 a 2,147,483,647)
var e int64  // 64 bits (-9,223,372,036,854,775,808 a 9,223,372,036,854,775,807)
```

#### Enteros sin signo

```go
var a uint    // Tamaño dependiente de la plataforma (32 o 64 bits)
var b uint8   // 8 bits (0 a 255)
var c uint16  // 16 bits (0 a 65,535)
var d uint32  // 32 bits (0 a 4,294,967,295)
var e uint64  // 64 bits (0 a 18,446,744,073,709,551,615)
var f byte    // Alias para uint8, comúnmente usado para datos binarios
var g rune    // Alias para int32, representa un punto de código Unicode
```

### Números de punto flotante

```go
var a float32  // IEEE-754 de 32 bits
var b float64  // IEEE-754 de 64 bits (mayor precisión, recomendado)
```

### Números complejos

```go
var a complex64   // Parte real e imaginaria de float32
var b complex128  // Parte real e imaginaria de float64

// Ejemplos
c1 := 1 + 2i      // complex128
c2 := complex(3, 4) // complex128, equivalente a 3 + 4i
```

### Booleanos

```go
var a bool = true
var b bool = false
```

### Cadenas de texto (string)

En Go, las cadenas son inmutables y contienen bytes UTF-8:

```go
var nombre string = "Gopher"

// Cadenas multilínea con comillas invertidas (backticks)
descripcion := `Esta es una cadena
que ocupa varias
líneas.`
```

#### Operaciones con cadenas

```go
// Concatenación
nombreCompleto := nombre + " " + apellido

// Longitud (en bytes, no en caracteres Unicode)
longitud := len(nombre)

// Acceso a bytes individuales (no caracteres Unicode)
primeraLetra := nombre[0] // tipo byte (uint8)

// Subcadenas
primerosTres := nombre[0:3] // "Gop"
```

### Valor cero

En Go, las variables tienen un "valor cero" predeterminado cuando se declaran sin inicializar:

- Numéricos: `0`
- Booleanos: `false`
- Cadenas: `""` (cadena vacía)
- Punteros, funciones, interfaces, slices, canales, mapas: `nil`

```go
var i int     // i = 0
var f float64 // f = 0.0
var b bool    // b = false
var s string  // s = ""
```

## Constantes

Las constantes son valores que no pueden cambiar durante la ejecución del programa:

```go
const Pi = 3.14159
const (
    Domingo  = 0
    Lunes    = 1
    Martes   = 2
    Miercoles = 3
    Jueves   = 4
    Viernes  = 5
    Sabado   = 6
)

// Constantes no tipadas
const x = 42       // x es una constante no tipada
const y float64 = 42 // y es una constante tipada (float64)
```

### iota

Go proporciona `iota` para crear constantes enumeradas:

```go
const (
    Domingo  = iota // 0
    Lunes           // 1
    Martes          // 2
    Miercoles       // 3
    Jueves          // 4
    Viernes         // 5
    Sabado          // 6
)

// Uso más complejo de iota
const (
    KB = 1 << (10 * iota) // 1 << (10 * 0) = 1
    MB                     // 1 << (10 * 1) = 1024
    GB                     // 1 << (10 * 2) = 1048576
    TB                     // 1 << (10 * 3) = 1073741824
)
```

## Conversión de tipos

Go no realiza conversiones implícitas entre tipos. Debes convertir explícitamente:

```go
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)

// Conversión entre string y números
s := "100"
i, err := strconv.Atoi(s) // string a int
if err != nil {
    // manejar error
}

s2 := strconv.Itoa(i) // int a string

// Conversiones más específicas
f, err := strconv.ParseFloat("3.14", 64) // string a float64
b, err := strconv.ParseBool("true")      // string a bool
```

## Tipos compuestos

### Arrays

Un array tiene un tamaño fijo definido en tiempo de compilación:

```go
// Declaración
var a [5]int // Array de 5 enteros, inicializado con ceros

// Inicialización
b := [3]int{1, 2, 3}

// Tamaño inferido
c := [...]int{1, 2, 3, 4} // Array de 4 enteros

// Acceso a elementos
b[0] = 10
fmt.Println(b[1]) // 2

// Longitud
longitud := len(b) // 3
```

### Slices

Un slice es una vista dinámica de un array:

```go
// Declaración
var s []int // Slice vacío (nil)

// Creación desde un array
a := [5]int{1, 2, 3, 4, 5}
s := a[1:4] // s contiene {2, 3, 4}

// Creación directa
s := []int{1, 2, 3} // Slice con 3 elementos

// Creación con make
s := make([]int, 5)    // Longitud 5, capacidad 5
s := make([]int, 3, 5) // Longitud 3, capacidad 5

// Añadir elementos
s = append(s, 6, 7, 8)

// Copiar slices
destino := make([]int, len(origen))
copy(destino, origen)
```

### Maps

Un map es una colección de pares clave-valor:

```go
// Declaración
var m map[string]int // Map vacío (nil)

// Creación
m := make(map[string]int)

// Asignacion de Valores
m["Juan"] = 1

// Inicialización
mad := map[string]int{
    "uno":  1,
    "dos":  2,
    "tres": 3,
}

// Acceso y modificación
m["cuatro"] = 4
valor := m["uno"] // 1

// Verificar existencia
valor, existe := m["cinco"]
if !existe {
    fmt.Println("La clave 'cinco' no existe")
}

// Eliminar elemento
delete(m, "dos")
```

### Structs

Un struct es una colección de campos:

```go
// Definición
type Persona struct {
    Nombre   string
    Edad     int
    Direccion string
}

// Creación
var p Persona
p.Nombre = "Gopher"
p.Edad = 25

// Inicialización
p := Persona{Nombre: "Gopher", Edad: 25}

// Inicialización posicional (no recomendada)
p := Persona{"Gopher", 25, "Calle Go 123"}

// Struct anónimo
usuario := struct {
    Nombre string
    Admin  bool
}{
    Nombre: "Admin",
    Admin:  true,
}
```

## Ejemplo práctico

Veamos un ejemplo que utiliza diferentes tipos de datos:

```go
package main

import (
    "fmt"
    "math"
)

// Definición de constantes
const (
    Pi     = 3.14159
    Saludo = "Hola, "
)

// Definición de struct
type Circulo struct {
    Radio float64
}

// Método para calcular el área
func (c Circulo) Area() float64 {
    return Pi * math.Pow(c.Radio, 2)
}

func main() {
    // Variables básicas
    var nombre string = "Gopher"
    edad := 25
    
    // Imprimir usando variables
    fmt.Println(Saludo + nombre)
    fmt.Printf("%s tiene %d años\n", nombre, edad)
    
    // Arrays y slices
    numeros := [5]int{1, 2, 3, 4, 5}
    subNumeros := numeros[1:4] // [2, 3, 4]
    subNumeros = append(subNumeros, 6)
    
    fmt.Println("Array completo:", numeros)
    fmt.Println("Slice:", subNumeros)
    
    // Maps
    capitales := map[string]string{
        "España":  "Madrid",
        "Francia": "París",
        "Italia":  "Roma",
    }
    
    capitales["Portugal"] = "Lisboa"
    
    for pais, capital := range capitales {
        fmt.Printf("La capital de %s es %s\n", pais, capital)
    }
    
    // Structs
    circulo := Circulo{Radio: 5}
    fmt.Printf("Un círculo con radio %.2f tiene un área de %.2f\n", 
              circulo.Radio, circulo.Area())
}
```

## Buenas prácticas

- Usa la forma corta `:=` dentro de funciones para mayor concisión
- Usa nombres descriptivos para tus variables
- Sigue las convenciones de nomenclatura (camelCase para variables locales, PascalCase para exportadas)
- Inicializa mapas con `make()` antes de usarlos
- Usa structs para agrupar datos relacionados
- Aprovecha el tipado estático para detectar errores en tiempo de compilación
- Usa constantes para valores que no cambian
- Prefiere slices sobre arrays para la mayoría de los casos

## Recursos adicionales

- [Especificación del lenguaje Go - Tipos](https://golang.org/ref/spec#Types)
- [Effective Go - Variables](https://golang.org/doc/effective_go#variables)
- [Go by Example - Variables](https://gobyexample.com/variables)
- [Go by Example - Constants](https://gobyexample.com/constants)
- [Tour of Go - Tipos básicos](https://tour.golang.org/basics/11)

---

En la siguiente sección, aprenderemos sobre operadores y expresiones en Go.