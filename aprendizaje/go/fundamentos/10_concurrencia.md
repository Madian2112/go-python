# Concurrencia en Go

## Introducción

La concurrencia es una de las características más distintivas y poderosas de Go. Diseñada para simplificar la programación concurrente, Go ofrece primitivas de concurrencia ligeras y eficientes que permiten escribir programas que aprovechan al máximo los sistemas multicore modernos. En esta sección, exploraremos los conceptos fundamentales de la concurrencia en Go, incluyendo goroutines, canales, patrones de concurrencia comunes y técnicas para manejar la sincronización.

## Goroutines

Las goroutines son funciones que se ejecutan concurrentemente con otras goroutines en el mismo espacio de direcciones. Son extremadamente ligeras, con un tamaño de pila inicial de solo unos pocos kilobytes, lo que permite crear miles o incluso millones de goroutines en una sola aplicación.

### Creación de Goroutines

Para iniciar una goroutine, simplemente antepone la palabra clave `go` a una llamada de función:

```go
package main

import (
    "fmt"
    "time"
)

func decir(mensaje string) {
    for i := 0; i < 5; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Println(mensaje)
    }
}

func main() {
    // Ejecutar decir("mundo") en una nueva goroutine
    go decir("mundo")
    
    // decir("hola") se ejecuta en la goroutine principal
    decir("hola")
}
```

En este ejemplo, `decir("mundo")` se ejecuta concurrentemente con `decir("hola")`. La función `main` y la goroutine se ejecutan en el mismo espacio de direcciones, por lo que ambas pueden acceder a las mismas variables.

### Características de las Goroutines

1. **Ligereza**: Las goroutines son mucho más ligeras que los hilos del sistema operativo. Puedes crear miles de goroutines sin problemas.

2. **Multiplexación**: El runtime de Go multiplexa goroutines en un número menor de hilos del sistema operativo, lo que reduce la sobrecarga de cambio de contexto.

3. **Comunicación**: Las goroutines se comunican a través de canales, lo que facilita la sincronización y evita problemas comunes de la programación concurrente.

4. **Inicio rápido**: Iniciar una goroutine es tan simple como anteponer `go` a una llamada de función.

### Goroutines y Funciones Anónimas

Es común usar funciones anónimas para crear goroutines:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Goroutine con función anónima
    go func() {
        for i := 0; i < 5; i++ {
            time.Sleep(100 * time.Millisecond)
            fmt.Println("mundo")
        }
    }()
    
    // Goroutine con función anónima y parámetros
    go func(mensaje string) {
        for i := 0; i < 5; i++ {
            time.Sleep(150 * time.Millisecond)
            fmt.Println(mensaje)
        }
    }("goroutine")
    
    // Código en la goroutine principal
    for i := 0; i < 5; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Println("hola")
    }
}
```

### Esperar a que las Goroutines Terminen

Un problema común es que el programa principal puede terminar antes de que las goroutines hayan completado su ejecución. Hay varias formas de esperar a que las goroutines terminen:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func trabajador(id int, wg *sync.WaitGroup) {
    defer wg.Done() // Notificar que hemos terminado al salir
    
    fmt.Printf("Trabajador %d iniciando\n", id)
    time.Sleep(time.Second) // Simular trabajo
    fmt.Printf("Trabajador %d terminado\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    // Lanzar varios trabajadores
    for i := 1; i <= 5; i++ {
        wg.Add(1) // Incrementar el contador
        go trabajador(i, &wg)
    }
    
    // Esperar a que todos los trabajadores terminen
    wg.Wait()
    
    fmt.Println("Todos los trabajadores han terminado")
}
```

En este ejemplo, usamos `sync.WaitGroup` para esperar a que todas las goroutines terminen antes de que el programa principal continúe.

## Canales

Los canales son un tipo de dato en Go que proporciona un mecanismo para que las goroutines se comuniquen y sincronicen su ejecución. Los canales implementan el principio de "no compartir memoria para comunicarse; en su lugar, comunicarse para compartir memoria".

### Creación y Uso Básico de Canales

```go
package main

import "fmt"

func main() {
    // Crear un canal de enteros
    ch := make(chan int)
    
    // Enviar un valor al canal en una goroutine
    go func() {
        ch <- 42 // Enviar 42 al canal
    }()
    
    // Recibir el valor del canal
    valor := <-ch
    fmt.Println(valor) // Imprime: 42
}
```

### Canales con Buffer

Por defecto, los canales son no bufferizados, lo que significa que solo pueden contener un valor a la vez y bloquean hasta que ese valor es recibido. Los canales con buffer tienen una capacidad definida y solo bloquean cuando el buffer está lleno (para envíos) o vacío (para recepciones).

```go
package main

import "fmt"

func main() {
    // Crear un canal con buffer de tamaño 3
    ch := make(chan int, 3)
    
    // Enviar valores sin bloquear
    ch <- 1
    ch <- 2
    ch <- 3
    
    // Recibir valores
    fmt.Println(<-ch) // 1
    fmt.Println(<-ch) // 2
    fmt.Println(<-ch) // 3
}
```

### Cierre de Canales

Un canal puede ser cerrado para indicar que no se enviarán más valores. Los receptores pueden comprobar si un canal está cerrado:

```go
package main

import "fmt"

func generador(ch chan int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch) // Cerrar el canal
}

func main() {
    ch := make(chan int)
    go generador(ch)
    
    // Método 1: Comprobar si el canal está cerrado
    for {
        valor, ok := <-ch
        if !ok {
            break // Canal cerrado
        }
        fmt.Println(valor)
    }
    
    // Método 2: Usar range (más conciso)
    ch = make(chan int)
    go generador(ch)
    
    for valor := range ch {
        fmt.Println(valor)
    }
}
```

### Canales Unidireccionales

Los canales pueden ser restringidos para solo enviar o solo recibir valores, lo que ayuda a clarificar la intención y prevenir errores:

```go
package main

import "fmt"

// solo puede enviar a ch
func productor(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)
}

// solo puede recibir de ch
func consumidor(ch <-chan int, done chan<- bool) {
    for valor := range ch {
        fmt.Println(valor)
    }
    done <- true
}

func main() {
    ch := make(chan int)
    done := make(chan bool)
    
    go productor(ch)
    go consumidor(ch, done)
    
    <-done // Esperar a que el consumidor termine
}
```

### Select

La declaración `select` permite esperar en múltiples operaciones de canal. Es similar a `switch`, pero para canales:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    // Goroutine que envía a ch1 después de 1 segundo
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "uno"
    }()
    
    // Goroutine que envía a ch2 después de 2 segundos
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "dos"
    }()
    
    // Esperar en ambos canales
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Recibido", msg1)
        case msg2 := <-ch2:
            fmt.Println("Recibido", msg2)
        }
    }
}
```

### Timeout con Select

`select` también puede incluir un caso de timeout usando `time.After`:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "respuesta"
    }()
    
    select {
    case res := <-ch:
        fmt.Println("Recibido:", res)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout después de 1 segundo")
    }
}
```

### Caso Default en Select

Un caso `default` en `select` se ejecuta inmediatamente si ningún otro caso está listo:

```go
package main

import "fmt"

func main() {
    ch := make(chan string)
    
    select {
    case msg := <-ch:
        fmt.Println("Recibido mensaje:", msg)
    default:
        fmt.Println("No hay mensajes disponibles")
    }
}
```

Esto hace que `select` sea no bloqueante.

## Sincronización

Aunque los canales son la forma preferida de sincronización en Go, el paquete `sync` proporciona primitivas de sincronización tradicionales para casos donde los canales no son la mejor opción.

### Mutex

Un mutex (exclusión mutua) protege datos de accesos concurrentes:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var contador int
    var mu sync.Mutex
    var wg sync.WaitGroup
    
    // Lanzar 1000 goroutines que incrementan el contador
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            // Proteger el acceso al contador
            mu.Lock()
            contador++
            mu.Unlock()
        }()
    }
    
    // Esperar a que todas las goroutines terminen
    wg.Wait()
    
    fmt.Println("Contador final:", contador)
}
```

### RWMutex

Un `RWMutex` (mutex de lectura/escritura) permite múltiples lectores o un solo escritor:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var datos map[string]string
    var rwmu sync.RWMutex
    var wg sync.WaitGroup
    
    // Inicializar el mapa
    datos = make(map[string]string)
    
    // Goroutine escritora
    wg.Add(1)
    go func() {
        defer wg.Done()
        
        for i := 0; i < 5; i++ {
            rwmu.Lock() // Bloqueo exclusivo para escritura
            datos[fmt.Sprintf("clave%d", i)] = fmt.Sprintf("valor%d", i)
            fmt.Println("Escritor: añadida clave", i)
            rwmu.Unlock()
            
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    // Múltiples goroutines lectoras
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := 0; j < 10; j++ {
                rwmu.RLock() // Bloqueo compartido para lectura
                fmt.Printf("Lector %d: datos = %v\n", id, datos)
                rwmu.RUnlock()
                
                time.Sleep(50 * time.Millisecond)
            }
        }(i)
    }
    
    // Esperar a que todas las goroutines terminen
    wg.Wait()
}
```

### Once

`sync.Once` garantiza que una función se ejecute exactamente una vez, incluso si es llamada desde múltiples goroutines:

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var once sync.Once
    var wg sync.WaitGroup
    
    // Función que queremos ejecutar solo una vez
    inicializar := func() {
        fmt.Println("Inicialización realizada")
    }
    
    // Lanzar múltiples goroutines que intentan inicializar
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            fmt.Printf("Goroutine %d intentando inicializar\n", id)
            once.Do(inicializar) // Solo se ejecutará una vez
            fmt.Printf("Goroutine %d terminada\n", id)
        }(i)
    }
    
    wg.Wait()
}
```

### Cond

`sync.Cond` implementa una variable de condición, útil para señalizar a múltiples goroutines:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var mu sync.Mutex
    cond := sync.NewCond(&mu)
    
    // Variable que queremos monitorear
    listo := false
    
    // Lanzar goroutines que esperan la señal
    for i := 0; i < 3; i++ {
        go func(id int) {
            // Adquirir el lock
            mu.Lock()
            
            // Esperar mientras la condición no se cumpla
            for !listo {
                fmt.Printf("Goroutine %d: esperando...\n", id)
                cond.Wait() // Libera el mutex y se bloquea hasta que recibe Signal o Broadcast
            }
            
            // La condición se cumple
            fmt.Printf("Goroutine %d: ¡condición cumplida!\n", id)
            mu.Unlock()
        }(i)
    }
    
    // Esperar un poco y luego señalizar
    time.Sleep(time.Second)
    
    mu.Lock()
    listo = true
    cond.Broadcast() // Despertar a todas las goroutines que esperan
    mu.Unlock()
    
    // Dar tiempo para que las goroutines impriman
    time.Sleep(time.Second)
}
```

## Patrones de Concurrencia

Go facilita la implementación de patrones de concurrencia comunes. Veamos algunos de los más útiles:

### Patrón Worker Pool

Un grupo de trabajadores que procesan tareas de un canal:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Tarea a realizar
type Tarea struct {
    ID      int
    Trabajo int // Simula algún trabajo (por ejemplo, tiempo en ms)
}

// Trabajador que procesa tareas
func trabajador(id int, tareas <-chan Tarea, resultados chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for tarea := range tareas {
        fmt.Printf("Trabajador %d iniciando tarea %d\n", id, tarea.ID)
        time.Sleep(time.Duration(tarea.Trabajo) * time.Millisecond) // Simular trabajo
        resultado := fmt.Sprintf("Tarea %d completada por trabajador %d", tarea.ID, id)
        resultados <- resultado
    }
}

func main() {
    numTrabajadores := 3
    numTareas := 10
    
    // Crear canales
    tareas := make(chan Tarea, numTareas)
    resultados := make(chan string, numTareas)
    
    // Iniciar trabajadores
    var wg sync.WaitGroup
    for w := 1; w <= numTrabajadores; w++ {
        wg.Add(1)
        go trabajador(w, tareas, resultados, &wg)
    }
    
    // Enviar tareas
    for t := 1; t <= numTareas; t++ {
        tareas <- Tarea{ID: t, Trabajo: 100 + t*50} // Tareas con diferentes duraciones
    }
    close(tareas) // No más tareas
    
    // Esperar a que los trabajadores terminen en una goroutine separada
    go func() {
        wg.Wait()
        close(resultados) // Cerrar el canal de resultados cuando todos los trabajadores terminen
    }()
    
    // Recoger resultados
    for resultado := range resultados {
        fmt.Println(resultado)
    }
}
```

### Patrón Pipeline

Una serie de etapas conectadas por canales, donde cada etapa procesa los datos y los pasa a la siguiente:

```go
package main

import (
    "fmt"
)

// Generador: produce enteros del 1 al n
func generador(n int) <-chan int {
    out := make(chan int)
    go func() {
        for i := 1; i <= n; i++ {
            out <- i
        }
        close(out)
    }()
    return out
}

// Cuadrado: calcula el cuadrado de cada número recibido
func cuadrado(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Filtro: solo deja pasar números pares
func filtrarPares(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            if n%2 == 0 {
                out <- n
            }
        }
        close(out)
    }()
    return out
}

func main() {
    // Construir el pipeline
    numeros := generador(10)           // Genera: 1, 2, 3, ..., 10
    cuadrados := cuadrado(numeros)     // Genera: 1, 4, 9, ..., 100
    pares := filtrarPares(cuadrados)   // Genera: 4, 16, 36, 64, 100
    
    // Consumir los resultados
    for n := range pares {
        fmt.Println(n)
    }
}
```

### Patrón Fan-Out, Fan-In

Distribuir trabajo entre múltiples goroutines (fan-out) y luego combinar los resultados (fan-in):

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Generador de trabajo
func generador(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// Procesador (puede ser lento)
func procesador(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            // Simular procesamiento que toma tiempo
            time.Sleep(100 * time.Millisecond)
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Fan-In: combina múltiples canales en uno solo
func combinar(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    
    // Función para copiar de un canal a otro
    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            out <- n
        }
    }
    
    wg.Add(len(cs))
    // Iniciar una goroutine para cada canal de entrada
    for _, c := range cs {
        go output(c)
    }
    
    // Cerrar el canal de salida cuando todos los canales de entrada estén cerrados
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

func main() {
    // Fuente de datos
    origen := generador(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    
    // Fan-Out: distribuir el trabajo entre 3 procesadores
    proc1 := procesador(origen)
    proc2 := procesador(origen)
    proc3 := procesador(origen)
    
    // Fan-In: combinar los resultados
    for resultado := range combinar(proc1, proc2, proc3) {
        fmt.Println(resultado)
    }
}
```

### Patrón de Cancelación

Cancelar operaciones en curso usando un canal de cancelación:

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// Tarea que puede ser cancelada
func tarea(ctx context.Context) <-chan string {
    resultados := make(chan string)
    
    go func() {
        defer close(resultados)
        
        // Simular una tarea larga dividida en pasos
        for i := 1; i <= 10; i++ {
            // Comprobar si debemos cancelar
            select {
            case <-ctx.Done():
                resultados <- fmt.Sprintf("Tarea cancelada en el paso %d: %v", i, ctx.Err())
                return
            case <-time.After(200 * time.Millisecond):
                // Continuar con el siguiente paso
                resultados <- fmt.Sprintf("Paso %d completado", i)
            }
        }
        
        resultados <- "Tarea completada con éxito"
    }()
    
    return resultados
}

func main() {
    // Crear un contexto con cancelación
    ctx, cancelar := context.WithCancel(context.Background())
    
    // Iniciar la tarea
    resultados := tarea(ctx)
    
    // Cancelar después de 1 segundo
    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("Solicitando cancelación...")
        cancelar()
    }()
    
    // Recoger resultados
    for resultado := range resultados {
        fmt.Println(resultado)
    }
}
```

### Patrón de Timeout

Establecer un límite de tiempo para operaciones:

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// Operación que puede tardar mucho tiempo
func operacionLenta(ctx context.Context) (string, error) {
    // Crear un canal para el resultado
    resultadoCh := make(chan string)
    
    // Realizar la operación en una goroutine
    go func() {
        // Simular trabajo que tarda 2 segundos
        time.Sleep(2 * time.Second)
        resultadoCh <- "Operación completada"
    }()
    
    // Esperar el resultado o timeout
    select {
    case resultado := <-resultadoCh:
        return resultado, nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func main() {
    // Crear un contexto con timeout de 1 segundo
    ctx, cancelar := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancelar() // Asegurar que se liberen recursos
    
    // Intentar realizar la operación
    resultado, err := operacionLenta(ctx)
    if err != nil {
        fmt.Println("Error:", err) // Probablemente context.DeadlineExceeded
    } else {
        fmt.Println(resultado)
    }
}
```

## Context

El paquete `context` proporciona una forma de pasar plazos, señales de cancelación y otros valores a través de la API y entre procesos. Es especialmente útil para controlar goroutines.

### Creación y Uso Básico

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func operacion(ctx context.Context) {
    // Crear un canal para simular trabajo
    trabajoCh := make(chan struct{})
    
    // Realizar trabajo en una goroutine
    go func() {
        // Simular trabajo que tarda 2 segundos
        time.Sleep(2 * time.Second)
        trabajoCh <- struct{}{}
    }()
    
    // Esperar a que el trabajo termine o el contexto sea cancelado
    select {
    case <-trabajoCh:
        fmt.Println("Operación completada con éxito")
    case <-ctx.Done():
        fmt.Println("Operación cancelada:", ctx.Err())
    }
}

func main() {
    // Contexto con timeout
    ctx, cancelar := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancelar()
    
    operacion(ctx)
}
```

### Valores en el Contexto

```go
package main

import (
    "context"
    "fmt"
)

// Clave para el valor en el contexto
type clave string

func procesar(ctx context.Context) {
    // Obtener valor del contexto
    if userID, ok := ctx.Value(clave("userID")).(string); ok {
        fmt.Println("Procesando para usuario:", userID)
    } else {
        fmt.Println("Usuario no especificado")
    }
    
    // Obtener otro valor
    if traceID, ok := ctx.Value(clave("traceID")).(string); ok {
        fmt.Println("ID de seguimiento:", traceID)
    }
}

func main() {
    // Contexto base
    ctx := context.Background()
    
    // Añadir valores al contexto
    ctx = context.WithValue(ctx, clave("userID"), "user123")
    ctx = context.WithValue(ctx, clave("traceID"), "trace456")
    
    procesar(ctx)
}
```

### Contexto con Cancelación

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func operacionLarga(ctx context.Context) <-chan string {
    resultados := make(chan string)
    
    go func() {
        defer close(resultados)
        
        // Simular operación larga con múltiples pasos
        for i := 1; i <= 5; i++ {
            // Verificar cancelación antes de cada paso
            select {
            case <-ctx.Done():
                resultados <- fmt.Sprintf("Cancelado en paso %d: %v", i, ctx.Err())
                return
            default:
                // Continuar con el trabajo
            }
            
            // Simular trabajo
            time.Sleep(200 * time.Millisecond)
            resultados <- fmt.Sprintf("Paso %d completado", i)
        }
        
        resultados <- "Operación completada con éxito"
    }()
    
    return resultados
}

func main() {
    // Crear contexto con cancelación
    ctx, cancelar := context.WithCancel(context.Background())
    
    // Iniciar operación
    resultados := operacionLarga(ctx)
    
    // Cancelar después de 0.5 segundos
    go func() {
        time.Sleep(500 * time.Millisecond)
        fmt.Println("Solicitando cancelación...")
        cancelar()
    }()
    
    // Recoger resultados
    for resultado := range resultados {
        fmt.Println(resultado)
    }
}
```

## Ejemplo Práctico: Servidor Web Concurrente

Vamos a crear un servidor web simple que utiliza concurrencia para manejar múltiples solicitudes simultáneamente y realiza operaciones en segundo plano:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "sync"
    "time"
)

// Estructura para almacenar datos
type Almacen struct {
    mu    sync.RWMutex
    datos map[string]string
}

// Crear un nuevo almacén
func NuevoAlmacen() *Almacen {
    return &Almacen{
        datos: make(map[string]string),
    }
}

// Obtener un valor
func (a *Almacen) Obtener(clave string) (string, bool) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    valor, existe := a.datos[clave]
    return valor, existe
}

// Establecer un valor
func (a *Almacen) Establecer(clave, valor string) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.datos[clave] = valor
}

// Eliminar una clave
func (a *Almacen) Eliminar(clave string) {
    a.mu.Lock()
    defer a.mu.Unlock()
    delete(a.datos, clave)
}

// Obtener todas las claves
func (a *Almacen) Claves() []string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    claves := make([]string, 0, len(a.datos))
    for k := range a.datos {
        claves = append(claves, k)
    }
    return claves
}

// Simular una operación que tarda tiempo
func operacionLenta(ctx context.Context, duracion time.Duration) (string, error) {
    select {
    case <-time.After(duracion):
        return fmt.Sprintf("Operación completada después de %v", duracion), nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func main() {
    // Crear almacén compartido
    almacen := NuevoAlmacen()
    
    // Inicializar con algunos datos
    almacen.Establecer("clave1", "valor1")
    almacen.Establecer("clave2", "valor2")
    
    // Canal para tareas en segundo plano
    tareasCh := make(chan func())
    
    // Iniciar trabajadores en segundo plano
    const numTrabajadores = 3
    var wg sync.WaitGroup
    wg.Add(numTrabajadores)
    
    for i := 0; i < numTrabajadores; i++ {
        go func(id int) {
            defer wg.Done()
            for tarea := range tareasCh {
                log.Printf("Trabajador %d ejecutando tarea", id)
                tarea()
            }
            log.Printf("Trabajador %d terminando", id)
        }(i)
    }
    
    // Manejador para obtener un valor
    http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
        clave := r.URL.Query().Get("key")
        if clave == "" {
            http.Error(w, "Parámetro 'key' requerido", http.StatusBadRequest)
            return
        }
        
        // Usar un contexto con timeout
        ctx, cancelar := context.WithTimeout(r.Context(), 500*time.Millisecond)
        defer cancelar()
        
        // Simular una operación que puede tardar
        resultado := make(chan struct {
            valor  string
            existe bool
            err    error
        })
        
        go func() {
            // Simular latencia aleatoria
            latencia := time.Duration(rand.Intn(1000)) * time.Millisecond
            time.Sleep(latencia)
            
            valor, existe := almacen.Obtener(clave)
            resultado <- struct {
                valor  string
                existe bool
                err    error
            }{valor, existe, nil}
        }()
        
        // Esperar resultado o timeout
        select {
        case res := <-resultado:
            if !res.existe {
                http.Error(w, "Clave no encontrada", http.StatusNotFound)
                return
            }
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]string{"key": clave, "value": res.valor})
        
        case <-ctx.Done():
            http.Error(w, "Timeout al obtener el valor", http.StatusRequestTimeout)
        }
    })
    
    // Manejador para establecer un valor
    http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
            return
        }
        
        var datos struct {
            Clave string `json:"key"`
            Valor string `json:"value"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
            http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
            return
        }
        
        if datos.Clave == "" {
            http.Error(w, "Clave no puede estar vacía", http.StatusBadRequest)
            return
        }
        
        // Establecer el valor
        almacen.Establecer(datos.Clave, datos.Valor)
        
        // Programar una tarea en segundo plano
        tareasCh <- func() {
            log.Printf("Procesando en segundo plano para clave: %s", datos.Clave)
            time.Sleep(2 * time.Second) // Simular trabajo
            log.Printf("Procesamiento en segundo plano completado para clave: %s", datos.Clave)
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"status": "success"})
    })
    
    // Manejador para listar todas las claves
    http.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
        claves := almacen.Claves()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string][]string{"keys": claves})
    })
    
    // Manejador para operación lenta con posibilidad de cancelación
    http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
        // Usar el contexto de la solicitud para manejar cancelación
        ctx := r.Context()
        
        // Simular una operación que tarda entre 1 y 3 segundos
        duracion := time.Duration(1+rand.Intn(2)) * time.Second
        log.Printf("Iniciando operación lenta de %v", duracion)
        
        resultado, err := operacionLenta(ctx, duracion)
        if err != nil {
            http.Error(w, fmt.Sprintf("Operación cancelada: %v", err), http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"result": resultado})
    })
    
    // Iniciar servidor
    puerto := ":8080"
    servidor := &http.Server{
        Addr: puerto,
    }
    
    // Manejar señal de cierre
    go func() {
        // En una aplicación real, aquí manejaríamos señales del sistema operativo
        time.Sleep(30 * time.Second) // Simular tiempo de ejecución
        log.Println("Iniciando apagado del servidor...")
        
        // Crear contexto con timeout para el apagado
        ctx, cancelar := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancelar()
        
        // Intentar apagar el servidor de forma ordenada
        if err := servidor.Shutdown(ctx); err != nil {
            log.Printf("Error al apagar el servidor: %v\n", err)
        }
        
        // Cerrar el canal de tareas y esperar a que los trabajadores terminen
        close(tareasCh)
        wg.Wait()
        log.Println("Todos los trabajadores han terminado")
    }()
    
    log.Printf("Servidor iniciado en http://localhost%s\n", puerto)
    if err := servidor.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("Error al iniciar el servidor: %v\n", err)
    }
    log.Println("Servidor detenido correctamente")
}
```

Este ejemplo demuestra:

1. Uso de `sync.RWMutex` para proteger acceso concurrente a datos compartidos
2. Uso de goroutines para manejar solicitudes HTTP concurrentemente
3. Uso de canales para comunicación entre goroutines
4. Uso de `context` para manejar timeouts y cancelación
5. Implementación de un worker pool para tareas en segundo plano
6. Apagado ordenado del servidor

## Buenas Prácticas

1. **Evitar Goroutines Huérfanas**: Asegúrate de que todas las goroutines terminen correctamente. Usa `context` o canales de cancelación para señalizar cuando deben detenerse.

2. **Usar Canales para Comunicación**: Sigue el principio "no compartir memoria para comunicarse; en su lugar, comunicarse para compartir memoria".

3. **Proteger Datos Compartidos**: Si necesitas compartir datos entre goroutines, usa primitivas de sincronización como `sync.Mutex` o `sync.RWMutex`.

4. **Preferir Concurrencia Simple**: Mantén tus patrones de concurrencia lo más simples posible. La complejidad innecesaria puede llevar a errores difíciles de depurar.

5. **Usar `context` para Cancelación**: El paquete `context` proporciona una forma estándar de manejar cancelación y timeouts.

6. **Limitar el Número de Goroutines**: Aunque las goroutines son ligeras, crear demasiadas puede agotar los recursos del sistema. Usa worker pools para limitar su número.

7. **Manejar Errores en Goroutines**: Los errores en goroutines no se propagan automáticamente. Usa canales para comunicar errores de vuelta a la goroutine principal.

8. **Cerrar Canales Correctamente**: Solo el remitente debe cerrar un canal, nunca el receptor. Cerrar un canal es una forma de señalizar que no habrá más valores.

9. **Usar `select` con Caso Default para Operaciones No Bloqueantes**: Incluye un caso `default` en `select` cuando necesites una operación no bloqueante.

10. **Usar `sync.WaitGroup` para Esperar Goroutines**: Cuando necesites esperar a que un grupo de goroutines termine, usa `sync.WaitGroup`.

## Recursos Adicionales

- [Effective Go - Concurrency](https://golang.org/doc/effective_go#concurrency)
- [Go Blog - Share Memory By Communicating](https://blog.golang.org/share-memory-by-communicating)
- [Go Blog - Go Concurrency Patterns: Pipelines and Cancellation](https://blog.golang.org/pipelines)
- [Go Blog - Go Concurrency Patterns: Context](https://blog.golang.org/context)
- [Go Blog - Go Concurrency Patterns: Timing out, moving on](https://blog.golang.org/go-concurrency-patterns-timing-out-and)
- [Go by Example - Goroutines](https://gobyexample.com/goroutines)
- [Go by Example - Channels](https://gobyexample.com/channels)
- [Go by Example - Select](https://gobyexample.com/select)
- [Go by Example - Timeouts](https://gobyexample.com/timeouts)
- [Go by Example - Non-Blocking Channel Operations](https://gobyexample.com/non-blocking-channel-operations)
- [Go by Example - Closing Channels](https://gobyexample.com/closing-channels)
- [Go by Example - Range over Channels](https://gobyexample.com/range-over-channels)
- [Go by Example - Worker Pools](https://gobyexample.com/worker-pools)
- [Go by Example - WaitGroups](https://gobyexample.com/waitgroups)
- [Go by Example - Rate Limiting](https://gobyexample.com/rate-limiting)
- [Go by Example - Atomic Counters](https://gobyexample.com/atomic-counters)
- [Go by Example - Mutexes](https://gobyexample.com/mutexes)

---

En la siguiente sección, exploraremos las pruebas en Go, incluyendo pruebas unitarias, de integración, de rendimiento y herramientas para garantizar la calidad del código.