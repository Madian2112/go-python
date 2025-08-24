# Funciones en Go

## Introducción

Las funciones son bloques de código reutilizables que realizan una tarea específica. En Go, las funciones son ciudadanos de primera clase, lo que significa que pueden ser asignadas a variables, pasadas como argumentos y devueltas por otras funciones. Go proporciona una sintaxis clara y concisa para definir y utilizar funciones.

## Definición de Funciones

En Go, las funciones se definen con la palabra clave `func`, seguida del nombre de la función, parámetros entre paréntesis y tipo de retorno:

```go
func nombreFuncion(parametro1 tipo1, parametro2 tipo2) tipoRetorno {
    // Cuerpo de la función
    // Código a ejecutar
    return valor
}
```

Ejemplo simple:

```go
func saludar(nombre string) string {
    return "¡Hola, " + nombre + "!"
}

// Llamada a la función
mensaje := saludar("Ana")
fmt.Println(mensaje)  // Imprime: ¡Hola, Ana!
```

## Parámetros y Argumentos

### Parámetros del mismo tipo

Cuando varios parámetros consecutivos tienen el mismo tipo, puedes omitir el tipo excepto para el último:

```go
func sumar(a, b int) int {
    return a + b
}

// Equivalente a:
// func sumar(a int, b int) int {
//     return a + b
// }
```

### Parámetros variádicos

Go permite definir funciones que aceptan un número variable de argumentos usando el operador `...`:

```go
func sumar(numeros ...int) int {
    total := 0
    for _, numero := range numeros {
        total += numero
    }
    return total
}

// Llamadas a la función
sumar(1, 2)           // 3
sumar(1, 2, 3, 4, 5)  // 15

// También puedes pasar un slice
nums := []int{1, 2, 3, 4, 5}
sumar(nums...)  // 15
```

Los parámetros variádicos deben ser el último parámetro en la lista de parámetros.

## Valores de Retorno

### Retorno simple

Las funciones pueden devolver un valor usando la instrucción `return`:

```go
func cuadrado(x int) int {
    return x * x
}

resultado := cuadrado(5)  // 25
```

### Retorno múltiple

Una característica distintiva de Go es la capacidad de devolver múltiples valores desde una función:

```go
func dividir(dividendo, divisor int) (int, int) {
    cociente := dividendo / divisor
    resto := dividendo % divisor
    return cociente, resto
}

// Capturar ambos valores retornados
cociente, resto := dividir(10, 3)
fmt.Printf("10 ÷ 3 = %d con resto %d\n", cociente, resto)

// Si solo te interesa uno de los valores, puedes usar _ para ignorar el otro
cociente, _ = dividir(10, 3)  // Ignora el resto
```

### Valores de retorno con nombre

Go permite nombrar los valores de retorno, lo que los inicializa como variables locales y permite usar `return` sin argumentos:

```go
func dividir(dividendo, divisor int) (cociente int, resto int) {
    cociente = dividendo / divisor
    resto = dividendo % divisor
    return  // Devuelve automáticamente cociente y resto
}

// También puedes especificar los valores en el return
func dividir(dividendo, divisor int) (cociente int, resto int) {
    cociente = dividendo / divisor
    resto = dividendo % divisor
    return cociente, resto  // Explícito, pero redundante
}
```

Los valores de retorno con nombre mejoran la legibilidad del código y la documentación, especialmente cuando los tipos de retorno no son obvios.

## Funciones como Valores

En Go, las funciones son ciudadanos de primera clase, lo que significa que pueden ser tratadas como cualquier otro valor.

### Asignación a variables

```go
func saludar(nombre string) string {
    return "¡Hola, " + nombre + "!"
}

// Asignar la función a una variable
var miFuncion func(string) string
miFuncion = saludar

// Llamar a la función a través de la variable
fmt.Println(miFuncion("Ana"))  // Imprime: ¡Hola, Ana!
```

### Funciones como argumentos

```go
func aplicar(f func(int) int, valor int) int {
    return f(valor)
}

func cuadrado(x int) int {
    return x * x
}

func cubo(x int) int {
    return x * x * x
}

fmt.Println(aplicar(cuadrado, 5))  // 25
fmt.Println(aplicar(cubo, 5))      // 125
```

### Funciones anónimas

Go permite definir funciones sin nombre (anónimas), útiles para funciones que solo se usan en un lugar:

```go
// Función anónima asignada a una variable
cuadrado := func(x int) int {
    return x * x
}
fmt.Println(cuadrado(5))  // 25

// Función anónima ejecutada inmediatamente
resultado := func(x int) int {
    return x * x
}(5)
fmt.Println(resultado)  // 25

// Función anónima como argumento
numeros := []int{1, 2, 3, 4, 5}
sort.Slice(numeros, func(i, j int) bool {
    return numeros[i] < numeros[j]  // Ordenar ascendentemente
})
```

### Closures

Las funciones en Go pueden capturar y acceder a variables de su ámbito circundante, formando un closure:

```go
func crearContador() func() int {
    contador := 0
    return func() int {
        contador++  // Accede y modifica la variable contador del ámbito exterior
        return contador
    }
}

contador := crearContador()
fmt.Println(contador())  // 1
fmt.Println(contador())  // 2
fmt.Println(contador())  // 3

// Cada llamada a crearContador crea un nuevo contador independiente
contador2 := crearContador()
fmt.Println(contador2())  // 1
fmt.Println(contador())   // 4 (el primer contador sigue en 4)
```

## Recursión

La recursión es cuando una función se llama a sí misma. Es útil para problemas que pueden dividirse en subproblemas más pequeños del mismo tipo:

```go
func factorial(n uint) uint {
    if n == 0 || n == 1 {
        return 1
    }
    return n * factorial(n-1)
}

fmt.Println(factorial(5))  // 120 (5 * 4 * 3 * 2 * 1)
```

Consideraciones sobre la recursión en Go:
- Debe tener un caso base para evitar recursión infinita
- Go no optimiza la recursión de cola (tail recursion)
- Para problemas grandes, una solución iterativa puede ser más eficiente

## Manejo de Errores

Go utiliza valores de retorno múltiples para manejar errores, en lugar de excepciones:

```go
func dividir(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("división por cero")
    }
    return a / b, nil
}

// Uso de la función con manejo de errores
resultado, err := dividir(10, 0)
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println("Resultado:", resultado)
}
```

### Creación de errores personalizados

Puedes crear errores personalizados de varias formas:

```go
// Usando errors.New
err1 := errors.New("error simple")

// Usando fmt.Errorf (permite formateo)
err2 := fmt.Errorf("error en archivo %s: %v", "datos.txt", "no encontrado")

// Definiendo un tipo de error personalizado
type MiError struct {
    Codigo int
    Mensaje string
}

func (e *MiError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Codigo, e.Mensaje)
}

// Uso del error personalizado
func operacion() error {
    return &MiError{404, "Recurso no encontrado"}
}
```

## Defer

La instrucción `defer` programa una llamada a función para que se ejecute justo antes de que la función actual retorne. Es útil para limpiar recursos:

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

Características importantes de `defer`:

1. Las funciones diferidas se ejecutan en orden LIFO (último en entrar, primero en salir)

```go
func ejemplo() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
}
// Imprime: 3, 2, 1
```

2. Los argumentos de las funciones diferidas se evalúan inmediatamente

```go
func ejemplo() {
    i := 1
    defer fmt.Println(i)  // Captura i=1
    i = 2
}
// Imprime: 1 (no 2)
```

3. Las funciones diferidas pueden acceder y modificar variables con nombre de la función que las contiene

```go
func ejemplo() (resultado int) {
    defer func() { resultado *= 2 }()  // Modifica el valor de retorno
    return 5
}
// Devuelve: 10
```

## Panic y Recover

Go proporciona mecanismos para situaciones excepcionales: `panic` y `recover`.

### Panic

`panic` detiene la ejecución normal de la función actual y comienza a desenrollar la pila, ejecutando cualquier función diferida:

```go
func dividir(a, b int) int {
    if b == 0 {
        panic("división por cero")
    }
    return a / b
}
```

### Recover

`recover` permite a un programa capturar un pánico y continuar la ejecución normal:

```go
func seguro() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recuperado de:", r)
        }
    }()
    
    dividir(10, 0)  // Esto causará un pánico
    fmt.Println("Esta línea nunca se ejecutará")
}

func main() {
    seguro()
    fmt.Println("Programa continúa normalmente")
}
```

**Nota**: `panic` y `recover` son similares a las excepciones en otros lenguajes, pero en Go se recomienda usarlos solo para errores irrecuperables, no para el flujo de control normal.

## Métodos

Los métodos en Go son funciones asociadas a un tipo específico. Se definen con un "receptor" que aparece entre la palabra clave `func` y el nombre del método:

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

func main() {
    r := Rectangulo{Ancho: 10, Alto: 5}
    fmt.Println("Área:", r.Area())  // 50
    
    r.Escalar(2)
    fmt.Println("Nuevo ancho:", r.Ancho)  // 20
    fmt.Println("Nuevo alto:", r.Alto)    // 10
    fmt.Println("Nueva área:", r.Area())  // 200
}
```

### Receptor de valor vs. receptor de puntero

- **Receptor de valor** (`func (r Rectangulo) ...`): Recibe una copia del valor. No puede modificar el valor original.
- **Receptor de puntero** (`func (r *Rectangulo) ...`): Recibe un puntero al valor. Puede modificar el valor original.

Cuándo usar cada uno:

- Usa receptores de puntero cuando necesites modificar el receptor o cuando el receptor sea una estructura grande (para evitar copias innecesarias).
- Usa receptores de valor para tipos inmutables (como tipos numéricos) o cuando quieras enfatizar que el método no modifica el receptor.

## Ejemplos Prácticos

### Calculadora con funciones

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

// Funciones de operaciones
func suma(a, b float64) float64 {
    return a + b
}

func resta(a, b float64) float64 {
    return a - b
}

func multiplicacion(a, b float64) float64 {
    return a * b
}

func division(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("error: división por cero")
    }
    return a / b, nil
}

// Tipo para representar una operación
type operacion struct {
    nombre string
    fn     func(float64, float64) (float64, error)
}

func main() {
    // Definir operaciones disponibles
    operaciones := map[string]operacion{
        "1": {"Suma", func(a, b float64) (float64, error) { return suma(a, b), nil }},
        "2": {"Resta", func(a, b float64) (float64, error) { return resta(a, b), nil }},
        "3": {"Multiplicación", func(a, b float64) (float64, error) { return multiplicacion(a, b), nil }},
        "4": {"División", division},
    }

    for {
        fmt.Println("\nCalculadora Simple")
        for key, op := range operaciones {
            fmt.Printf("%s. %s\n", key, op.nombre)
        }
        fmt.Println("5. Salir")

        var opcion string
        fmt.Print("Seleccione una opción (1-5): ")
        fmt.Scanln(&opcion)

        if opcion == "5" {
            fmt.Println("¡Hasta luego!")
            break
        }

        op, existe := operaciones[opcion]
        if !existe {
            fmt.Println("Opción no válida. Intente de nuevo.")
            continue
        }

        var num1, num2 float64
        var err error

        fmt.Print("Ingrese el primer número: ")
        _, err = fmt.Scanln(&num1)
        if err != nil {
            fmt.Println("Error: Ingrese un número válido")
            continue
        }

        fmt.Print("Ingrese el segundo número: ")
        _, err = fmt.Scanln(&num2)
        if err != nil {
            fmt.Println("Error: Ingrese un número válido")
            continue
        }

        resultado, err := op.fn(num1, num2)
        if err != nil {
            fmt.Println(err)
        } else {
            fmt.Printf("Resultado de %s: %.2f\n", op.nombre, resultado)
        }
    }
}
```

### Generador de contraseñas

```go
package main

import (
    "crypto/rand"
    "fmt"
    "math/big"
    "strings"
)

// Configuración para la generación de contraseñas
type ConfigContrasena struct {
    Longitud         int
    IncluirMinusculas bool
    IncluirMayusculas bool
    IncluirNumeros    bool
    IncluirEspeciales bool
}

// Genera una contraseña aleatoria según la configuración
func generarContrasena(config ConfigContrasena) (string, error) {
    // Verificar que al menos un conjunto de caracteres esté habilitado
    if !config.IncluirMinusculas && !config.IncluirMayusculas &&
       !config.IncluirNumeros && !config.IncluirEspeciales {
        return "", fmt.Errorf("error: debe habilitar al menos un conjunto de caracteres")
    }

    // Verificar longitud mínima
    if config.Longitud < 4 {
        return "", fmt.Errorf("error: la longitud mínima es 4")
    }

    // Definir conjuntos de caracteres
    const (
        minusculas = "abcdefghijklmnopqrstuvwxyz"
        mayusculas = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
        numeros    = "0123456789"
        especiales = "!@#$%^&*()-_=+[]{}|;:,.<>?/"
    )

    // Construir el conjunto de caracteres disponibles
    var caracteresDisponibles strings.Builder
    if config.IncluirMinusculas {
        caracteresDisponibles.WriteString(minusculas)
    }
    if config.IncluirMayusculas {
        caracteresDisponibles.WriteString(mayusculas)
    }
    if config.IncluirNumeros {
        caracteresDisponibles.WriteString(numeros)
    }
    if config.IncluirEspeciales {
        caracteresDisponibles.WriteString(especiales)
    }

    caracteres := caracteresDisponibles.String()
    maxIndex := big.NewInt(int64(len(caracteres)))

    // Generar la contraseña
    var contrasena strings.Builder
    contrasena.Grow(config.Longitud)

    // Asegurar que al menos un carácter de cada tipo esté presente si se solicita
    if config.IncluirMinusculas && config.Longitud > 0 {
        idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(minusculas))))
        contrasena.WriteByte(minusculas[idx.Int64()])
    }
    if config.IncluirMayusculas && contrasena.Len() < config.Longitud {
        idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(mayusculas))))
        contrasena.WriteByte(mayusculas[idx.Int64()])
    }
    if config.IncluirNumeros && contrasena.Len() < config.Longitud {
        idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(numeros))))
        contrasena.WriteByte(numeros[idx.Int64()])
    }
    if config.IncluirEspeciales && contrasena.Len() < config.Longitud {
        idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(especiales))))
        contrasena.WriteByte(especiales[idx.Int64()])
    }

    // Completar con caracteres aleatorios
    for contrasena.Len() < config.Longitud {
        idx, err := rand.Int(rand.Reader, maxIndex)
        if err != nil {
            return "", fmt.Errorf("error al generar número aleatorio: %v", err)
        }
        contrasena.WriteByte(caracteres[idx.Int64()])
    }

    // Convertir a slice para poder mezclar
    contrasenaMezclada := []rune(contrasena.String())

    // Mezclar los caracteres (Fisher-Yates shuffle)
    for i := len(contrasenaMezclada) - 1; i > 0; i-- {
        j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
        contrasenaMezclada[i], contrasenaMezclada[j.Int64()] = contrasenaMezclada[j.Int64()], contrasenaMezclada[i]
    }

    return string(contrasenaMezclada), nil
}

func main() {
    fmt.Println("Generador de Contraseñas\n")

    config := ConfigContrasena{
        Longitud:         12,  // Valor predeterminado
        IncluirMinusculas: true,
        IncluirMayusculas: true,
        IncluirNumeros:    true,
        IncluirEspeciales: true,
    }

    // Solicitar configuración al usuario
    fmt.Print("Longitud de la contraseña (mínimo 4): ")
    fmt.Scanln(&config.Longitud)
    if config.Longitud < 4 {
        fmt.Println("La longitud mínima es 4. Se usará 4 como longitud.")
        config.Longitud = 4
    }

    var respuesta string

    fmt.Print("¿Incluir minúsculas? (s/n): ")
    fmt.Scanln(&respuesta)
    config.IncluirMinusculas = strings.ToLower(respuesta) == "s"

    fmt.Print("¿Incluir mayúsculas? (s/n): ")
    fmt.Scanln(&respuesta)
    config.IncluirMayusculas = strings.ToLower(respuesta) == "s"

    fmt.Print("¿Incluir números? (s/n): ")
    fmt.Scanln(&respuesta)
    config.IncluirNumeros = strings.ToLower(respuesta) == "s"

    fmt.Print("¿Incluir caracteres especiales? (s/n): ")
    fmt.Scanln(&respuesta)
    config.IncluirEspeciales = strings.ToLower(respuesta) == "s"

    contrasena, err := generarContrasena(config)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("\nContraseña generada: %s\n", contrasena)
}
```

### Función para calcular estadísticas

```go
package main

import (
    "fmt"
    "math"
    "sort"
    "strings"
    "strconv"
)

// Estructura para almacenar estadísticas
type Estadisticas struct {
    Cantidad          int
    Minimo            float64
    Maximo            float64
    Suma              float64
    Media             float64
    Mediana           float64
    Moda              []float64
    DesviacionEstandar float64
    Rango             float64
    Q1                float64
    Q3                float64
    RangoIntercuartil float64
}

// Calcula estadísticas para un slice de números
func calcularEstadisticas(numeros []float64) (Estadisticas, error) {
    if len(numeros) == 0 {
        return Estadisticas{}, fmt.Errorf("error: la lista está vacía")
    }

    // Crear una copia ordenada
    ordenados := make([]float64, len(numeros))
    copy(ordenados, numeros)
    sort.Float64s(ordenados)

    n := len(ordenados)
    minimo := ordenados[0]
    maximo := ordenados[n-1]

    // Calcular suma y media
    suma := 0.0
    for _, num := range ordenados {
        suma += num
    }
    media := suma / float64(n)

    // Calcular mediana
    var mediana float64
    if n%2 == 0 { // Si hay un número par de elementos
        mediana = (ordenados[n/2-1] + ordenados[n/2]) / 2
    } else { // Si hay un número impar de elementos
        mediana = ordenados[n/2]
    }

    // Calcular moda (valor más frecuente)
    frecuencias := make(map[float64]int)
    for _, num := range ordenados {
        frecuencias[num]++
    }

    maxFrecuencia := 0
    for _, freq := range frecuencias {
        if freq > maxFrecuencia {
            maxFrecuencia = freq
        }
    }

    var moda []float64
    for num, freq := range frecuencias {
        if freq == maxFrecuencia {
            moda = append(moda, num)
        }
    }

    // Calcular desviación estándar
    sumaCuadradosDiff := 0.0
    for _, num := range numeros {
        diff := num - media
        sumaCuadradosDiff += diff * diff
    }
    desviacionEstandar := math.Sqrt(sumaCuadradosDiff / float64(n))

    // Calcular rango y cuartiles
    rango := maximo - minimo

    // Calcular Q1 (primer cuartil)
    q1Pos := n / 4
    var q1 float64
    if n%4 == 0 {
        q1 = (ordenados[q1Pos-1] + ordenados[q1Pos]) / 2
    } else {
        q1 = ordenados[q1Pos]
    }

    // Calcular Q3 (tercer cuartil)
    q3Pos := 3 * n / 4
    var q3 float64
    if n%4 == 0 {
        q3 = (ordenados[q3Pos-1] + ordenados[q3Pos]) / 2
    } else {
        q3 = ordenados[q3Pos]
    }

    rangoIntercuartil := q3 - q1

    return Estadisticas{
        Cantidad:          n,
        Minimo:            minimo,
        Maximo:            maximo,
        Suma:              suma,
        Media:             media,
        Mediana:           mediana,
        Moda:              moda,
        DesviacionEstandar: desviacionEstandar,
        Rango:             rango,
        Q1:                q1,
        Q3:                q3,
        RangoIntercuartil: rangoIntercuartil,
    }, nil
}

func main() {
    fmt.Println("Calculadora de Estadísticas\n")

    fmt.Print("Ingrese una lista de números separados por espacios: ")
    var entrada string
    fmt.Scanln(&entrada)

    // Dividir la entrada en tokens
    tokens := strings.Fields(entrada)
    if len(tokens) == 0 {
        fmt.Println("No se ingresaron números.")
        return
    }

    // Convertir tokens a números
    numeros := make([]float64, 0, len(tokens))
    for _, token := range tokens {
        num, err := strconv.ParseFloat(token, 64)
        if err != nil {
            fmt.Printf("Error al convertir '%s' a número: %v\n", token, err)
            return
        }
        numeros = append(numeros, num)
    }

    // Calcular estadísticas
    estadisticas, err := calcularEstadisticas(numeros)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Mostrar resultados
    fmt.Println("\nEstadísticas:")
    fmt.Printf("Cantidad: %d\n", estadisticas.Cantidad)
    fmt.Printf("Mínimo: %.2f\n", estadisticas.Minimo)
    fmt.Printf("Máximo: %.2f\n", estadisticas.Maximo)
    fmt.Printf("Suma: %.2f\n", estadisticas.Suma)
    fmt.Printf("Media: %.2f\n", estadisticas.Media)
    fmt.Printf("Mediana: %.2f\n", estadisticas.Mediana)

    fmt.Print("Moda: ")
    if len(estadisticas.Moda) == 1 {
        fmt.Printf("%.2f\n", estadisticas.Moda[0])
    } else {
        fmt.Print("[ ")
        for i, m := range estadisticas.Moda {
            if i > 0 {
                fmt.Print(", ")
            }
            fmt.Printf("%.2f", m)
        }
        fmt.Println(" ]")
    }

    fmt.Printf("Desviación estándar: %.2f\n", estadisticas.DesviacionEstandar)
    fmt.Printf("Rango: %.2f\n", estadisticas.Rango)
    fmt.Printf("Q1: %.2f\n", estadisticas.Q1)
    fmt.Printf("Q3: %.2f\n", estadisticas.Q3)
    fmt.Printf("Rango intercuartil: %.2f\n", estadisticas.RangoIntercuartil)
}
```

## Buenas Prácticas

1. **Nombres descriptivos**: Usa nombres de funciones que describan claramente lo que hacen. En Go, se usa camelCase para funciones no exportadas y PascalCase para funciones exportadas.

2. **Funciones pequeñas y específicas**: Cada función debe hacer una sola cosa y hacerla bien. Divide funciones grandes en funciones más pequeñas y específicas.

3. **Manejo de errores explícito**: Devuelve errores explícitamente en lugar de usar panic. Verifica siempre los errores devueltos por las funciones.

4. **Usa defer para limpiar recursos**: Utiliza defer para asegurarte de que los recursos se liberen correctamente, incluso en caso de error.

5. **Comentarios útiles**: Comenta el propósito de la función, no cómo funciona (a menos que sea complejo). Los comentarios para funciones exportadas deben comenzar con el nombre de la función.

```go
// CalcularArea devuelve el área de un rectángulo.
// El área se calcula multiplicando el ancho por el alto.
func CalcularArea(ancho, alto float64) float64 {
    return ancho * alto
}
```

6. **Evita variables globales**: Prefiere pasar estado como parámetros en lugar de usar variables globales.

7. **Usa receptores de puntero cuando sea necesario**: Usa receptores de puntero cuando necesites modificar el receptor o cuando el receptor sea una estructura grande.

8. **Retorno temprano para manejo de errores**: Verifica errores al principio de la función y retorna temprano si hay errores.

```go
func procesarDatos(datos []int) (int, error) {
    if len(datos) == 0 {
        return 0, errors.New("datos vacíos")
    }
    // Procesar datos...
    return resultado, nil
}
```

9. **Usa valores de retorno con nombre cuando mejoren la legibilidad**: Especialmente útil cuando devuelves múltiples valores del mismo tipo.

10. **Evita interfaces vacías cuando sea posible**: Prefiere tipos concretos o interfaces específicas sobre `interface{}`.

## Recursos Adicionales

- [Especificación del lenguaje Go - Funciones](https://golang.org/ref/spec#Function_declarations)
- [Tour of Go - Funciones](https://tour.golang.org/basics/4)
- [Effective Go - Funciones](https://golang.org/doc/effective_go#functions)
- [Go by Example - Funciones](https://gobyexample.com/functions)
- [Go by Example - Funciones Variádicas](https://gobyexample.com/variadic-functions)
- [Go by Example - Closures](https://gobyexample.com/closures)
- [Go by Example - Recursión](https://gobyexample.com/recursion)
- [Go by Example - Defer](https://gobyexample.com/defer)
- [Go by Example - Panic](https://gobyexample.com/panic)
- [Go by Example - Recover](https://gobyexample.com/recover)

---

En la siguiente sección, aprenderemos sobre estructuras de datos en Go, que nos permitirán organizar y manipular datos de manera eficiente.