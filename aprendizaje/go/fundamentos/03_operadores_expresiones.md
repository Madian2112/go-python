# Operadores y Expresiones en Go

## Introducción

Los operadores en Go son símbolos que indican al compilador que realice operaciones específicas sobre operandos (variables y valores). Las expresiones son combinaciones de operadores y operandos que se evalúan para producir un resultado.

## Operadores Aritméticos

Realizan operaciones matemáticas básicas:

| Operador | Descripción | Ejemplo | Resultado |
|----------|-------------|---------|----------|
| `+` | Suma | `5 + 3` | `8` |
| `-` | Resta | `5 - 3` | `2` |
| `*` | Multiplicación | `5 * 3` | `15` |
| `/` | División | `5 / 3` | `1` (con enteros) o `1.6666...` (con flotantes) |
| `%` | Módulo (resto) | `5 % 3` | `2` |

```go
// Ejemplos de operadores aritméticos
a := 10
b := 3

suma := a + b          // 13
resta := a - b         // 7
multiplicacion := a * b  // 30

// División con enteros (trunca el resultado)
divisionEntera := a / b   // 3

// División con flotantes (mantiene decimales)
divisionFloat := float64(a) / float64(b)  // 3.3333...

// Módulo (solo para enteros)
modulo := a % b        // 1
```

### Comportamiento de la división

En Go, el comportamiento de la división depende del tipo de los operandos:

- Si ambos operandos son enteros, el resultado es un entero (se trunca la parte decimal)
- Si al menos uno de los operandos es de punto flotante, el resultado es de punto flotante

```go
fmt.Println(5 / 2)       // 2 (división entera)
fmt.Println(5.0 / 2)     // 2.5 (división flotante)
fmt.Println(5 / 2.0)     // 2.5 (división flotante)
```

## Operadores de Asignación

Asignan valores a variables:

| Operador | Descripción | Ejemplo | Equivalente a |
|----------|-------------|---------|---------------|
| `=` | Asignación simple | `x = 5` | `x = 5` |
| `+=` | Suma y asignación | `x += 3` | `x = x + 3` |
| `-=` | Resta y asignación | `x -= 3` | `x = x - 3` |
| `*=` | Multiplicación y asignación | `x *= 3` | `x = x * 3` |
| `/=` | División y asignación | `x /= 3` | `x = x / 3` |
| `%=` | Módulo y asignación | `x %= 3` | `x = x % 3` |

```go
// Ejemplos de operadores de asignación
x := 10

x += 5   // x ahora es 15
x -= 3   // x ahora es 12
x *= 2   // x ahora es 24
x /= 4   // x ahora es 6
x %= 4   // x ahora es 2
```

## Operadores de Incremento y Decremento

Go proporciona operadores para incrementar y decrementar variables en una unidad:

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `++` | Incremento | `x++` |
| `--` | Decremento | `x--` |

```go
// Ejemplos de incremento y decremento
x := 5
x++    // x ahora es 6
x--    // x ahora es 5 de nuevo
```

**Nota importante**: En Go, `++` y `--` son sentencias, no expresiones. Esto significa que no pueden ser usadas en asignaciones o como parte de expresiones más complejas:

```go
// Esto NO es válido en Go
y := x++    // Error de compilación
z := ++x    // Error de compilación
if (x++) > 5 { // Error de compilación
    // ...
}
```

## Operadores de Comparación

Comparan valores y devuelven un resultado booleano (true o false):

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `==` | Igual a | `a == b` |
| `!=` | Diferente de | `a != b` |
| `>` | Mayor que | `a > b` |
| `<` | Menor que | `a < b` |
| `>=` | Mayor o igual que | `a >= b` |
| `<=` | Menor o igual que | `a <= b` |

```go
// Ejemplos de operadores de comparación
a := 10
b := 5

fmt.Println(a == b)  // false
fmt.Println(a != b)  // true
fmt.Println(a > b)   // true
fmt.Println(a < b)   // false
fmt.Println(a >= b)  // true
fmt.Println(a <= b)  // false

// También funcionan con strings (comparación lexicográfica)
fmt.Println("apple" < "banana")  // true
```

## Operadores Lógicos

Combina expresiones booleanas:

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `&&` | AND lógico (true si ambos operandos son true) | `a && b` |
| `\|\|` | OR lógico (true si al menos un operando es true) | `a \|\| b` |
| `!` | NOT lógico (invierte el valor booleano) | `!a` |

```go
// Ejemplos de operadores lógicos
x := 5
y := 10

fmt.Println(x > 0 && y > 0)    // true (ambas condiciones son verdaderas)
fmt.Println(x > 10 || y > 5)   // true (la segunda condición es verdadera)
fmt.Println(!(x > 10))         // true (x > 10 es falso, ! lo invierte)

// Cortocircuito en operadores lógicos
// && : si el primer operando es false, el segundo no se evalúa
// || : si el primer operando es true, el segundo no se evalúa
```

### Tabla de verdad para operadores lógicos

| a | b | a && b | a \|\| b | !a |
|---|---|--------|---------|----|
| true | true | true | true | false |
| true | false | false | true | false |
| false | true | false | true | true |
| false | false | false | false | true |

## Operadores Bit a Bit

Realizan operaciones a nivel de bits:

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `&` | AND bit a bit | `a & b` |
| `\|` | OR bit a bit | `a \| b` |
| `^` | XOR bit a bit | `a ^ b` |
| `&^` | AND NOT bit a bit (bit clear) | `a &^ b` |
| `<<` | Desplazamiento a la izquierda | `a << n` |
| `>>` | Desplazamiento a la derecha | `a >> n` |

```go
// Ejemplos de operadores bit a bit
a := 60  // 0011 1100 en binario
b := 13  // 0000 1101 en binario

fmt.Println(a & b)    // 12 (0000 1100)
fmt.Println(a | b)    // 61 (0011 1101)
fmt.Println(a ^ b)    // 49 (0011 0001)
fmt.Println(a &^ b)   // 48 (0011 0000) - bits en a pero no en b
fmt.Println(a << 2)   // 240 (1111 0000)
fmt.Println(a >> 2)   // 15 (0000 1111)
```

### Operador AND NOT (&^)

Este operador es único en Go y se conoce como "bit clear". Para cada bit:
- Si el bit en b es 0, se mantiene el bit de a
- Si el bit en b es 1, el resultado es 0

Es útil para desactivar bits específicos.

## Operadores de Dirección

Go proporciona operadores para trabajar con punteros:

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `&` | Obtiene la dirección de memoria de una variable | `&a` |
| `*` | Desreferencia un puntero (accede al valor) | `*p` |

```go
// Ejemplos de operadores de dirección
x := 10
p := &x    // p contiene la dirección de memoria de x

fmt.Println(p)     // Imprime algo como 0xc000018030 (dirección de memoria)
fmt.Println(*p)    // 10 (valor almacenado en esa dirección)

*p = 20            // Modifica el valor de x a través del puntero
fmt.Println(x)     // 20
```

## Operadores de Canal

Go incluye operadores específicos para trabajar con canales (channels):

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `<-` | Enviar/recibir datos de un canal | `ch <- v` o `v := <-ch` |

```go
// Ejemplos de operadores de canal
ch := make(chan int)

// Enviar valor al canal (en una goroutine)
go func() { ch <- 42 }()

// Recibir valor del canal
valor := <-ch
fmt.Println(valor)  // 42
```

## Precedencia de Operadores

Los operadores en Go se evalúan en el siguiente orden de precedencia (de mayor a menor):

1. `()` (paréntesis)
2. `*` (desreferencia de puntero), `&` (dirección), `!` (NOT lógico), `+` (unario), `-` (unario), `^` (XOR bit a bit unario)
3. `*` (multiplicación), `/` (división), `%` (módulo), `<<` (desplazamiento izquierda), `>>` (desplazamiento derecha), `&` (AND bit a bit), `&^` (AND NOT bit a bit)
4. `+` (suma), `-` (resta), `|` (OR bit a bit), `^` (XOR bit a bit)
5. `==` (igual), `!=` (diferente), `<` (menor), `<=` (menor o igual), `>` (mayor), `>=` (mayor o igual)
6. `&&` (AND lógico)
7. `||` (OR lógico)

```go
// Ejemplos de precedencia
resultado := 2 + 3 * 4    // 14 (no 20, porque * tiene mayor precedencia que +)
resultado = (2 + 3) * 4  // 20 (los paréntesis cambian la precedencia)
```

## Expresiones

Las expresiones son combinaciones de valores, variables y operadores que se evalúan para producir un resultado.

### Tipos de expresiones

1. **Expresiones aritméticas**: Realizan cálculos matemáticos
   ```go
   resultado := (a + b) * c / d
   ```

2. **Expresiones relacionales**: Comparan valores
   ```go
   esMayor := edad >= 18
   ```

3. **Expresiones lógicas**: Combinan condiciones
   ```go
   puedeVotar := edad >= 18 && esCiudadano
   ```

4. **Expresiones de asignación**: Asignan valores a variables
   ```go
   x, y = y, x  // Intercambio de valores
   ```

## Ejemplos prácticos

### Calculadora simple

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    if len(os.Args) != 4 {
        fmt.Println("Uso: calculadora <número1> <operación> <número2>")
        fmt.Println("Operaciones soportadas: +, -, *, /, %")
        os.Exit(1)
    }

    num1, err := strconv.ParseFloat(os.Args[1], 64)
    if err != nil {
        fmt.Println("Error: El primer argumento no es un número válido")
        os.Exit(1)
    }

    operacion := os.Args[2]

    num2, err := strconv.ParseFloat(os.Args[3], 64)
    if err != nil {
        fmt.Println("Error: El tercer argumento no es un número válido")
        os.Exit(1)
    }

    var resultado float64

    switch operacion {
    case "+":
        resultado = num1 + num2
    case "-":
        resultado = num1 - num2
    case "*":
        resultado = num1 * num2
    case "/":
        if num2 == 0 {
            fmt.Println("Error: División por cero")
            os.Exit(1)
        }
        resultado = num1 / num2
    case "%":
        resultado = float64(int(num1) % int(num2))
    default:
        fmt.Println("Error: Operación no soportada")
        os.Exit(1)
    }

    fmt.Printf("%.2f %s %.2f = %.2f\n", num1, operacion, num2, resultado)
}
```

### Verificación de año bisiesto

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Uso: bisiesto <año>")
        os.Exit(1)
    }

    año, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Println("Error: El argumento no es un número válido")
        os.Exit(1)
    }

    // Un año es bisiesto si es divisible por 4,
    // excepto los divisibles por 100 que no son divisibles por 400
    esBisiesto := (año%4 == 0 && año%100 != 0) || (año%400 == 0)

    if esBisiesto {
        fmt.Printf("%d es un año bisiesto\n", año)
    } else {
        fmt.Printf("%d no es un año bisiesto\n", año)
    }
}
```

## Buenas prácticas

- Usa paréntesis para clarificar el orden de evaluación, incluso cuando no son estrictamente necesarios
- Evita expresiones demasiado complejas; divide en partes más simples si es necesario
- Usa variables con nombres descriptivos para almacenar resultados intermedios
- Recuerda que `++` y `--` son sentencias, no expresiones en Go
- Ten cuidado con la división por cero
- Aprovecha el intercambio de valores con asignación múltiple (`x, y = y, x`)
- Usa el operador `&^` (bit clear) cuando necesites desactivar bits específicos

## Recursos adicionales

- [Especificación del lenguaje Go - Expresiones](https://golang.org/ref/spec#Expressions)
- [Especificación del lenguaje Go - Operadores](https://golang.org/ref/spec#Operators)
- [Go by Example - Operadores](https://gobyexample.com/operators)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Playground](https://play.golang.org/) - Prueba expresiones en línea

---

En la siguiente sección, aprenderemos sobre estructuras de control en Go (condicionales y bucles).