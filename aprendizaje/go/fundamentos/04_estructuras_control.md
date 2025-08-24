# Estructuras de Control en Go

## Introducción

Las estructuras de control determinan el flujo de ejecución de un programa. Go ofrece un conjunto conciso pero potente de estructuras de control que permiten tomar decisiones, repetir acciones y organizar el código de manera eficiente.

## Estructuras Condicionales

Las estructuras condicionales permiten ejecutar diferentes bloques de código dependiendo de si una condición es verdadera o falsa.

### if, else if, else

La estructura básica de una condición en Go es:

```go
if condicion {
    // Código a ejecutar si la condición es true
} else if otraCondicion {
    // Código a ejecutar si la condición anterior es false y esta es true
} else {
    // Código a ejecutar si todas las condiciones anteriores son false
}
```

Ejemplo:

```go
edad := 18

if edad < 18 {
    fmt.Println("Eres menor de edad")
} else if edad == 18 {
    fmt.Println("Acabas de cumplir la mayoría de edad")
} else {
    fmt.Println("Eres mayor de edad")
}
```

### Declaración de variables en if

Go permite declarar variables dentro de la condición `if`, limitando su alcance al bloque `if` y sus bloques `else` asociados:

```go
if edad := obtenerEdad(); edad < 18 {
    fmt.Println("Eres menor de edad")
} else {
    fmt.Println("Eres mayor de edad")
}
// La variable 'edad' no está disponible aquí
```

Esta característica es útil para limitar el alcance de variables que solo se necesitan en el contexto de la condición.

### Operadores de comparación en condiciones

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `==` | Igual a | `a == b` |
| `!=` | Diferente de | `a != b` |
| `>` | Mayor que | `a > b` |
| `<` | Menor que | `a < b` |
| `>=` | Mayor o igual que | `a >= b` |
| `<=` | Menor o igual que | `a <= b` |

### Operadores lógicos

Permiten combinar múltiples condiciones:

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `&&` | AND lógico (verdadero si ambas condiciones son verdaderas) | `a > 0 && b > 0` |
| `\|\|` | OR lógico (verdadero si al menos una condición es verdadera) | `a > 0 \|\| b > 0` |
| `!` | NOT lógico (invierte el valor de la condición) | `!(a > 0)` |

Ejemplo:

```go
edad := 25
ingresos := 1500

if edad > 18 && ingresos > 1000 {
    fmt.Println("Cumples los requisitos para el préstamo")
} else {
    fmt.Println("No cumples los requisitos")
}
```

### Evaluación de condiciones

En Go, las condiciones deben ser expresiones booleanas explícitas. A diferencia de otros lenguajes, Go no permite usar valores no booleanos como condiciones:

```go
// Esto NO es válido en Go
if valor {
    // ...
}

// Debe ser una expresión booleana explícita
if valor != 0 {
    // ...
}

if len(cadena) > 0 {
    // ...
}
```

## Estructura switch

La estructura `switch` permite seleccionar uno de varios bloques de código para ejecutar, basándose en el valor de una expresión.

### switch básico

```go
switch expresion {
case valor1:
    // Código a ejecutar si expresion == valor1
case valor2, valor3:
    // Código a ejecutar si expresion == valor2 o expresion == valor3
default:
    // Código a ejecutar si ningún caso coincide
}
```

Ejemplo:

```go
dia := "lunes"

switch dia {
case "lunes":
    fmt.Println("Inicio de semana")
case "martes", "miércoles", "jueves":
    fmt.Println("Mitad de semana")
case "viernes":
    fmt.Println("Fin de semana laboral")
case "sábado", "domingo":
    fmt.Println("Fin de semana")
default:
    fmt.Println("Día no válido")
}
```

### switch sin expresión

Go permite usar `switch` sin una expresión, lo que equivale a `switch true`. Esto permite escribir condiciones más complejas en cada caso:

```go
puntuacion := 85

switch {
case puntuacion >= 90:
    fmt.Println("A")
case puntuacion >= 80:
    fmt.Println("B")
case puntuacion >= 70:
    fmt.Println("C")
case puntuacion >= 60:
    fmt.Println("D")
default:
    fmt.Println("F")
}
```

### Declaración de variables en switch

Al igual que con `if`, Go permite declarar variables dentro de la instrucción `switch`:

```go
switch hora := obtenerHora(); {
case hora < 12:
    fmt.Println("Buenos días")
case hora < 18:
    fmt.Println("Buenas tardes")
default:
    fmt.Println("Buenas noches")
}
```

### Comportamiento de fallthrough

A diferencia de otros lenguajes como C o Java, en Go cada caso termina automáticamente después de su ejecución (no es necesario usar `break`). Si deseas que la ejecución continúe con el siguiente caso, puedes usar la palabra clave `fallthrough`:

```go
n := 5
switch n {
case 5:
    fmt.Println("Es 5")
    fallthrough
case 4:
    fmt.Println("Es 4 o mayor")
    fallthrough
case 3:
    fmt.Println("Es 3 o mayor")
default:
    fmt.Println("Valor predeterminado")
}
// Salida:
// Es 5
// Es 4 o mayor
// Es 3 o mayor
```

**Nota**: `fallthrough` transfiere el control al siguiente caso sin evaluar su condición.

### Type switch

Go proporciona una variante especial de `switch` para determinar el tipo de una interfaz:

```go
var x interface{} = "Hola"

switch v := x.(type) {
case nil:
    fmt.Println("x es nil")
case int:
    fmt.Printf("x es un entero: %d\n", v)
case string:
    fmt.Printf("x es una cadena: %s\n", v)
case bool:
    fmt.Printf("x es un booleano: %t\n", v)
default:
    fmt.Printf("Tipo no manejado: %T\n", v)
}
```

## Estructuras de Repetición (Bucles)

Go simplifica las estructuras de repetición proporcionando solo el bucle `for`, pero con múltiples variantes que cubren todos los casos de uso.

### Bucle for estándar

La sintaxis estándar incluye inicialización, condición y post-instrucción:

```go
for inicializacion; condicion; post {
    // Código a ejecutar en cada iteración
}
```

Ejemplo:

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)  // Imprime: 0, 1, 2, 3, 4
}
```

### Bucle for como while

Omitiendo la inicialización y la post-instrucción, `for` funciona como un bucle `while`:

```go
i := 0
for i < 5 {
    fmt.Println(i)
    i++
}
```

### Bucle infinito

Si se omiten todas las partes, se crea un bucle infinito:

```go
for {
    // Código a ejecutar indefinidamente
    // Debe incluir una instrucción break para salir
    if condicionDeSalida {
        break
    }
}
```

### Iterando sobre colecciones

Go proporciona la construcción `for range` para iterar sobre arrays, slices, strings, maps y channels:

```go
// Iterando sobre un slice
numeros := []int{1, 2, 3, 4, 5}
for indice, valor := range numeros {
    fmt.Printf("numeros[%d] = %d\n", indice, valor)
}

// Iterando sobre un string (por runa/carácter)
for indice, caracter := range "¡Hola!" {
    fmt.Printf("%d: %c\n", indice, caracter)
}

// Iterando sobre un map
persona := map[string]string{
    "nombre": "Ana",
    "edad":   "25",
    "ciudad": "Madrid",
}
for clave, valor := range persona {
    fmt.Printf("%s: %s\n", clave, valor)
}
```

Si solo necesitas el índice o la clave, puedes omitir el valor usando el identificador `_`:

```go
for indice, _ := range numeros {
    fmt.Println(indice)
}
// O más conciso:
for indice := range numeros {
    fmt.Println(indice)
}
```

Si solo necesitas el valor, puedes omitir el índice:

```go
for _, valor := range numeros {
    fmt.Println(valor)
}
```

### Instrucciones break y continue

- `break`: Termina el bucle actual y continúa con la siguiente instrucción después del bucle.
- `continue`: Salta a la siguiente iteración del bucle, omitiendo el resto del código en la iteración actual.

Ejemplos:

```go
// Uso de break
for i := 0; i < 10; i++ {
    if i == 5 {
        break  // Sale del bucle cuando i es 5
    }
    fmt.Println(i)  // Imprime: 0, 1, 2, 3, 4
}

// Uso de continue
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue  // Salta a la siguiente iteración si i es par
    }
    fmt.Println(i)  // Imprime: 1, 3, 5, 7, 9
}
```

### Etiquetas y bucles anidados

Go permite etiquetar bucles y usar `break` o `continue` con una etiqueta específica, lo que es útil para controlar bucles anidados:

```go
exterior:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if i*j >= 3 {
            fmt.Println("Saliendo del bucle exterior")
            break exterior  // Sale de ambos bucles
        }
        fmt.Printf("i=%d, j=%d\n", i, j)
    }
}
```

## Instrucción goto

Go incluye la instrucción `goto` para saltar a una etiqueta dentro de la misma función. Aunque su uso generalmente no se recomienda, puede ser útil en ciertos casos:

```go
func ejemploGoto() {
    i := 0

incremento:
    fmt.Println(i)
    i++
    if i < 5 {
        goto incremento  // Salta de vuelta a la etiqueta
    }
}
```

**Nota**: El uso excesivo de `goto` puede hacer que el código sea difícil de leer y mantener. En la mayoría de los casos, es mejor usar estructuras de control estándar.

## Instrucción defer

Aunque no es estrictamente una estructura de control, `defer` es una característica importante de Go que afecta el flujo de ejecución. La instrucción `defer` programa una llamada a función para que se ejecute justo antes de que la función actual retorne:

```go
func ejemploDefer() {
    defer fmt.Println("Esto se imprime al final")
    fmt.Println("Esto se imprime primero")
}
```

Las instrucciones `defer` se ejecutan en orden LIFO (último en entrar, primero en salir):

```go
func ejemploDeferMultiple() {
    for i := 0; i < 5; i++ {
        defer fmt.Println(i)  // Se imprimirá: 4, 3, 2, 1, 0
    }
}
```

`defer` es especialmente útil para limpiar recursos como archivos, conexiones de red o bloqueos:

```go
func procesarArchivo(nombre string) error {
    archivo, err := os.Open(nombre)
    if err != nil {
        return err
    }
    defer archivo.Close()  // Se cerrará automáticamente al salir de la función
    
    // Procesar el archivo...
    return nil
}
```

## Ejemplos Prácticos

### Calculadora simple con menú

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    for {
        fmt.Println("\nCalculadora Simple")
        fmt.Println("1. Suma")
        fmt.Println("2. Resta")
        fmt.Println("3. Multiplicación")
        fmt.Println("4. División")
        fmt.Println("5. Salir")
        
        var opcion int
        fmt.Print("Seleccione una opción (1-5): ")
        fmt.Scan(&opcion)
        
        if opcion == 5 {
            fmt.Println("¡Hasta luego!")
            break
        }
        
        if opcion < 1 || opcion > 4 {
            fmt.Println("Opción no válida. Intente de nuevo.")
            continue
        }
        
        var num1, num2 float64
        fmt.Print("Ingrese el primer número: ")
        fmt.Scan(&num1)
        fmt.Print("Ingrese el segundo número: ")
        fmt.Scan(&num2)
        
        switch opcion {
        case 1:
            fmt.Printf("Resultado: %.2f\n", num1+num2)
        case 2:
            fmt.Printf("Resultado: %.2f\n", num1-num2)
        case 3:
            fmt.Printf("Resultado: %.2f\n", num1*num2)
        case 4:
            if num2 == 0 {
                fmt.Println("Error: No se puede dividir por cero")
            } else {
                fmt.Printf("Resultado: %.2f\n", num1/num2)
            }
        }
    }
}
```

### Verificación de número primo

```go
package main

import (
    "fmt"
)

func esPrimo(numero int) bool {
    if numero <= 1 {
        return false
    }
    if numero <= 3 {
        return true
    }
    if numero%2 == 0 || numero%3 == 0 {
        return false
    }
    
    for i := 5; i*i <= numero; i += 6 {
        if numero%i == 0 || numero%(i+2) == 0 {
            return false
        }
    }
    
    return true
}

func main() {
    for {
        var numero int
        fmt.Print("Ingrese un número (0 para salir): ")
        fmt.Scan(&numero)
        
        if numero == 0 {
            fmt.Println("¡Hasta luego!")
            break
        }
        
        if esPrimo(numero) {
            fmt.Printf("%d es un número primo\n", numero)
        } else {
            fmt.Printf("%d no es un número primo\n", numero)
        }
    }
}
```

### Generador de patrones

```go
package main

import (
    "fmt"
    "strings"
)

func imprimirPatron(filas int) {
    // Patrón de triángulo
    fmt.Println("\nPatrón 1:")
    for i := 1; i <= filas; i++ {
        fmt.Println(strings.Repeat("*", i))
    }
    
    // Patrón de pirámide
    fmt.Println("\nPatrón 2:")
    for i := 1; i <= filas; i++ {
        espacios := filas - i
        estrellas := 2*i - 1
        fmt.Println(strings.Repeat(" ", espacios) + strings.Repeat("*", estrellas))
    }
    
    // Patrón de diamante
    fmt.Println("\nPatrón 3:")
    for i := 1; i <= filas; i++ {
        espacios := filas - i
        estrellas := 2*i - 1
        fmt.Println(strings.Repeat(" ", espacios) + strings.Repeat("*", estrellas))
    }
    
    for i := filas - 1; i >= 1; i-- {
        espacios := filas - i
        estrellas := 2*i - 1
        fmt.Println(strings.Repeat(" ", espacios) + strings.Repeat("*", estrellas))
    }
}

func main() {
    var filas int
    fmt.Print("Ingrese el número de filas para los patrones: ")
    fmt.Scan(&filas)
    
    if filas <= 0 {
        fmt.Println("Por favor, ingrese un número positivo")
        return
    }
    
    imprimirPatron(filas)
}
```

## Buenas Prácticas

1. **Usa llaves en la misma línea**: En Go, las llaves de apertura deben estar en la misma línea que la declaración de la estructura de control, no en la siguiente línea.

2. **Evita bucles anidados profundos**: Los bucles anidados pueden hacer que el código sea difícil de leer y mantener. Intenta limitar la anidación a 2-3 niveles.

3. **Usa `switch` en lugar de múltiples `if-else`**: Cuando tengas múltiples condiciones que comparan la misma variable, `switch` suele ser más legible.

4. **Aprovecha la declaración de variables en `if` y `switch`**: Limita el alcance de las variables cuando sea posible.

5. **Usa `defer` para limpiar recursos**: Asegúrate de liberar recursos como archivos, conexiones de red o bloqueos con `defer`.

6. **Evita `goto`**: Aunque Go incluye `goto`, su uso generalmente no se recomienda. Prefiere estructuras de control estándar.

7. **Usa el identificador `_` para valores no utilizados**: Cuando iteres con `range` y no necesites alguno de los valores, usa `_` para ignorarlo.

8. **Verifica condiciones de borde**: Asegúrate de que tus bucles y condiciones manejen correctamente los casos límite.

9. **Prefiere `for range` para iterar sobre colecciones**: Es más conciso y menos propenso a errores que los bucles tradicionales.

10. **Recuerda que no hay bucle `do-while`**: Si necesitas un bucle que se ejecute al menos una vez, deberás implementarlo con un bucle `for` y una condición de salida.

## Recursos Adicionales

- [Especificación del lenguaje Go - Sentencias](https://golang.org/ref/spec#Statements)
- [Tour of Go - Control de flujo](https://tour.golang.org/flowcontrol/1)
- [Effective Go - Control Structures](https://golang.org/doc/effective_go#control-structures)
- [Go by Example - If/Else](https://gobyexample.com/if-else)
- [Go by Example - Switch](https://gobyexample.com/switch)
- [Go by Example - For](https://gobyexample.com/for)
- [Go by Example - Range](https://gobyexample.com/range)
- [Go by Example - Defer](https://gobyexample.com/defer)

---

En la siguiente sección, aprenderemos sobre funciones en Go, que nos permitirán organizar y reutilizar nuestro código de manera más eficiente.