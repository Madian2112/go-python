# Ejemplos de Código: Python vs Go (Golang)

Este documento presenta ejemplos de código básicos en Python y Go para ilustrar las diferencias de sintaxis y características entre ambos lenguajes.

## 1. Hola Mundo

### Python

```python
# Hola mundo en Python
print("¡Hola, mundo!")
```

### Go

```go
// Hola mundo en Go
package main

import "fmt"

func main() {
    fmt.Println("¡Hola, mundo!")
}
```

## 2. Variables y Tipos de Datos

### Python

```python
# Variables en Python (tipado dinámico)
nombre = "Ana"  # String
edad = 30       # Integer
altura = 1.75   # Float
activo = True   # Boolean

# Python infiere el tipo automáticamente
print(f"Nombre: {nombre}, Tipo: {type(nombre)}")
print(f"Edad: {edad}, Tipo: {type(edad)}")
print(f"Altura: {altura}, Tipo: {type(altura)}")
print(f"Activo: {activo}, Tipo: {type(activo)}")

# Reasignación con diferente tipo (permitido en Python)
edad = "treinta"  # Ahora es un string
print(f"Edad: {edad}, Nuevo tipo: {type(edad)}")
```

### Go

```go
// Variables en Go (tipado estático)
package main

import "fmt"

func main() {
    // Declaración explícita de tipo
    var nombre string = "Ana"
    var edad int = 30
    var altura float64 = 1.75
    var activo bool = true
    
    // Declaración con inferencia de tipo
    nombre2 := "Carlos"  // Go infiere que es string
    
    fmt.Printf("Nombre: %s, Tipo: %T\n", nombre, nombre)
    fmt.Printf("Edad: %d, Tipo: %T\n", edad, edad)
    fmt.Printf("Altura: %.2f, Tipo: %T\n", altura, altura)
    fmt.Printf("Activo: %t, Tipo: %T\n", activo, activo)
    fmt.Printf("Nombre2: %s, Tipo: %T\n", nombre2, nombre2)
    
    // Esto NO es posible en Go:
    // edad = "treinta"  // Error: no se puede asignar string a int
}
```

## 3. Estructuras de Control

### Python

```python
# Condicionales en Python
edad = 18

if edad < 18:
    print("Menor de edad")
elif edad == 18:
    print("Justo en la mayoría de edad")
else:
    print("Mayor de edad")

# Bucle for en Python
print("Contando del 1 al 5:")
for i in range(1, 6):
    print(i)

# Bucle while en Python
print("Cuenta regresiva:")
contador = 3
while contador > 0:
    print(contador)
    contador -= 1
print("¡Despegue!")
```

### Go

```go
// Condicionales en Go
package main

import "fmt"

func main() {
    edad := 18
    
    if edad < 18 {
        fmt.Println("Menor de edad")
    } else if edad == 18 {
        fmt.Println("Justo en la mayoría de edad")
    } else {
        fmt.Println("Mayor de edad")
    }
    
    // Bucle for en Go (Go solo tiene el bucle for, pero con diferentes formas)
    fmt.Println("Contando del 1 al 5:")
    for i := 1; i <= 5; i++ {
        fmt.Println(i)
    }
    
    // Equivalente a while en Go
    fmt.Println("Cuenta regresiva:")
    contador := 3
    for contador > 0 {
        fmt.Println(contador)
        contador--
    }
    fmt.Println("¡Despegue!")
}
```

## 4. Funciones

### Python

```python
# Funciones en Python
def saludar(nombre):
    return f"Hola, {nombre}!"

# Función con valores por defecto
def presentar(nombre, profesion="programador"):
    return f"{nombre} es {profesion}."

# Función con número variable de argumentos
def sumar(*numeros):
    total = 0
    for num in numeros:
        total += num
    return total

# Uso de las funciones
print(saludar("María"))
print(presentar("Juan"))
print(presentar("Ana", "ingeniera"))
print(f"Suma: {sumar(1, 2, 3, 4, 5)}")
```

### Go

```go
// Funciones en Go
package main

import "fmt"

// Función simple
func saludar(nombre string) string {
    return fmt.Sprintf("Hola, %s!", nombre)
}

// Función con múltiples parámetros
func presentar(nombre string, profesion string) string {
    return fmt.Sprintf("%s es %s.", nombre, profesion)
}

// Función con múltiples valores de retorno
func dividir(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("no se puede dividir por cero")
    }
    return a / b, nil
}

func main() {
    fmt.Println(saludar("María"))
    fmt.Println(presentar("Juan", "programador"))
    fmt.Println(presentar("Ana", "ingeniera"))
    
    // Manejo de múltiples valores de retorno
    resultado, err := dividir(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("10 / 2 = %.2f\n", resultado)
    }
    
    resultado, err = dividir(5, 0)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("5 / 0 = %.2f\n", resultado)
    }
}
```

## 5. Estructuras de Datos

### Python

```python
# Listas en Python
frutas = ["manzana", "banana", "cereza"]
print(f"Frutas: {frutas}")
frutas.append("naranja")
print(f"Después de agregar: {frutas}")
print(f"Primera fruta: {frutas[0]}")

# Diccionarios en Python
persona = {
    "nombre": "Laura",
    "edad": 28,
    "ciudad": "Madrid"
}
print(f"Persona: {persona}")
print(f"Nombre: {persona['nombre']}")
persona["profesion"] = "desarrolladora"
print(f"Después de agregar profesión: {persona}")

# Conjuntos en Python
colores = {"rojo", "verde", "azul"}
print(f"Colores: {colores}")
colores.add("amarillo")
print(f"Después de agregar: {colores}")
```

### Go

```go
// Estructuras de datos en Go
package main

import "fmt"

func main() {
    // Arrays en Go (tamaño fijo)
    var colores [3]string
    colores[0] = "rojo"
    colores[1] = "verde"
    colores[2] = "azul"
    fmt.Println("Array de colores:", colores)
    
    // Slices en Go (tamaño dinámico)
    frutas := []string{"manzana", "banana", "cereza"}
    fmt.Println("Frutas:", frutas)
    frutas = append(frutas, "naranja")
    fmt.Println("Después de agregar:", frutas)
    fmt.Println("Primera fruta:", frutas[0])
    
    // Maps en Go (equivalente a diccionarios)
    persona := map[string]interface{}{
        "nombre": "Laura",
        "edad":   28,
        "ciudad": "Madrid",
    }
    fmt.Println("Persona:", persona)
    fmt.Println("Nombre:", persona["nombre"])
    persona["profesion"] = "desarrolladora"
    fmt.Println("Después de agregar profesión:", persona)
}
```

## 6. Concurrencia

### Python

```python
# Concurrencia en Python con threading
import threading
import time

def tarea(nombre, segundos):
    print(f"Tarea {nombre} iniciada")
    time.sleep(segundos)  # Simula trabajo
    print(f"Tarea {nombre} completada después de {segundos} segundos")

# Crear y ejecutar hilos
hilo1 = threading.Thread(target=tarea, args=("A", 2))
hilo2 = threading.Thread(target=tarea, args=("B", 1))

print("Iniciando hilos...")
hilo1.start()
hilo2.start()

# Esperar a que ambos hilos terminen
hilo1.join()
hilo2.join()
print("Todos los hilos han terminado")

# Concurrencia en Python con asyncio (Python 3.5+)
import asyncio

async def tarea_asincrona(nombre, segundos):
    print(f"Tarea asíncrona {nombre} iniciada")
    await asyncio.sleep(segundos)  # Pausa sin bloquear
    print(f"Tarea asíncrona {nombre} completada después de {segundos} segundos")

async def main():
    print("Iniciando tareas asíncronas...")
    # Ejecutar tareas concurrentemente
    await asyncio.gather(
        tarea_asincrona("X", 2),
        tarea_asincrona("Y", 1)
    )
    print("Todas las tareas asíncronas han terminado")

# En Python 3.7+
# asyncio.run(main())
```

### Go

```go
// Concurrencia en Go con goroutines
package main

import (
    "fmt"
    "time"
)

func tarea(nombre string, segundos int) {
    fmt.Printf("Tarea %s iniciada\n", nombre)
    time.Sleep(time.Duration(segundos) * time.Second) // Simula trabajo
    fmt.Printf("Tarea %s completada después de %d segundos\n", nombre, segundos)
}

func main() {
    fmt.Println("Iniciando goroutines...")
    
    // Lanzar goroutines (funciones concurrentes)
    go tarea("A", 2)
    go tarea("B", 1)
    
    // Esperar para ver los resultados (en código real usaríamos canales o sync.WaitGroup)
    time.Sleep(3 * time.Second)
    fmt.Println("Programa principal terminado")
    
    // Ejemplo con canales para sincronización
    fmt.Println("\nEjemplo con canales:")
    done := make(chan bool)
    
    go func() {
        tarea("C", 2)
        done <- true
    }()
    
    go func() {
        tarea("D", 1)
        done <- true
    }()
    
    // Esperar a que ambas goroutines terminen
    <-done
    <-done
    fmt.Println("Todas las goroutines han terminado")
}
```

## 7. Manejo de Errores

### Python

```python
# Manejo de excepciones en Python
def dividir(a, b):
    try:
        resultado = a / b
        return resultado
    except ZeroDivisionError:
        print("Error: No se puede dividir por cero")
        return None
    except TypeError:
        print("Error: Tipo de datos incorrecto")
        return None
    finally:
        print("Operación de división finalizada")

# Probar diferentes casos
print(f"10 / 2 = {dividir(10, 2)}")
print(f"5 / 0 = {dividir(5, 0)}")
print(f"'5' / 2 = {dividir('5', 2)}")

# Crear y lanzar excepciones personalizadas
class EdadInvalidaError(Exception):
    pass

def verificar_edad(edad):
    if edad < 0:
        raise EdadInvalidaError("La edad no puede ser negativa")
    if edad > 120:
        raise EdadInvalidaError("La edad parece demasiado alta")
    return f"Edad válida: {edad}"

try:
    print(verificar_edad(25))
    print(verificar_edad(-5))
except EdadInvalidaError as e:
    print(f"Error de validación: {e}")
```

### Go

```go
// Manejo de errores en Go
package main

import (
    "errors"
    "fmt"
)

// Función que puede devolver un error
func dividir(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("no se puede dividir por cero")
    }
    return a / b, nil
}

// Error personalizado
type EdadInvalidaError struct {
    Edad    int
    Mensaje string
}

func (e *EdadInvalidaError) Error() string {
    return fmt.Sprintf("edad inválida %d: %s", e.Edad, e.Mensaje)
}

func verificarEdad(edad int) (string, error) {
    if edad < 0 {
        return "", &EdadInvalidaError{edad, "la edad no puede ser negativa"}
    }
    if edad > 120 {
        return "", &EdadInvalidaError{edad, "la edad parece demasiado alta"}
    }
    return fmt.Sprintf("Edad válida: %d", edad), nil
}

func main() {
    // Probar división
    resultado, err := dividir(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("10 / 2 = %.2f\n", resultado)
    }
    
    resultado, err = dividir(5, 0)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("5 / 0 = %.2f\n", resultado)
    }
    
    // Probar verificación de edad
    mensaje, err := verificarEdad(25)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println(mensaje)
    }
    
    mensaje, err = verificarEdad(-5)
    if err != nil {
        fmt.Println("Error:", err)
        
        // Type assertion para verificar el tipo de error
        if edadErr, ok := err.(*EdadInvalidaError); ok {
            fmt.Printf("Detalles adicionales - Edad intentada: %d\n", edadErr.Edad)
        }
    } else {
        fmt.Println(mensaje)
    }
}
```

Estos ejemplos muestran las diferencias fundamentales en la sintaxis y el enfoque de programación entre Python y Go. Python tiende a ser más conciso y flexible, mientras que Go es más explícito y estructurado, con un fuerte enfoque en la eficiencia y la concurrencia.