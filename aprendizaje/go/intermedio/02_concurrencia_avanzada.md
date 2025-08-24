# Concurrencia Avanzada en Go

## Introducción

La concurrencia es una de las características más destacadas de Go. El lenguaje fue diseñado desde el principio con la concurrencia en mente, proporcionando primitivas simples pero poderosas: goroutines y canales. En este módulo, profundizaremos en técnicas avanzadas de concurrencia en Go, explorando patrones, mejores prácticas y herramientas para construir aplicaciones concurrentes robustas y eficientes.

## Repaso de Conceptos Básicos

### Goroutines

Las goroutines son funciones que se ejecutan concurrentemente con otras goroutines en el mismo espacio de direcciones. Son más ligeras que los hilos del sistema operativo.

```go
func main() {
    go func() {
        fmt.Println("Ejecutando en una goroutine")
    }()
    
    // Dar tiempo a que la goroutine se ejecute
    time.Sleep(100 * time.Millisecond)
}
```

### Canales

Los canales son conductos tipados a través de los cuales puedes enviar y recibir valores con el operador de canal `<-`.

```go
func main() {
    ch := make(chan string)
    
    go func() {
        ch <- "Hola desde la goroutine"
    }()
    
    msg := <-ch
    fmt.Println(msg)
}
```

## Patrones de Concurrencia Avanzados

### Worker Pools

Un patrón común es crear un grupo de workers que procesen tareas de un canal.

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d procesando trabajo %d\n", id, j)
        time.Sleep(time.Second) // Simulando trabajo
        results <- j * 2
    }
}

func main() {
    const numJobs = 5
    const numWorkers = 3
    
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    
    // Iniciar workers
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }
    
    // Enviar trabajos
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)
    
    // Recoger resultados
    for a := 1; a <= numJobs; a++ {
        <-results
    }
}
```

### Pipeline

El patrón pipeline consiste en una serie de etapas conectadas por canales, donde cada etapa es un grupo de goroutines que realizan un subconjunto de trabajo.

```go
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    // Configurar pipeline
    c := gen(2, 3, 4, 5)
    out := sq(c)
    
    // Consumir resultado
    for n := range out {
        fmt.Println(n) // 4, 9, 16, 25
    }
}
```

### Fan-out, Fan-in

Fan-out es el proceso de iniciar múltiples goroutines para manejar entradas de un canal, y fan-in es el proceso de combinar múltiples resultados en un solo canal.

```go
func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    
    // Función para copiar de un canal a otro
    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    
    wg.Add(len(cs))
    // Iniciar una goroutine para cada canal de entrada
    for _, c := range cs {
        go output(c)
    }
    
    // Iniciar una goroutine para cerrar el canal de salida una vez que
    // todas las goroutines de entrada hayan terminado
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

func main() {
    in := gen(2, 3, 4, 5, 6, 7, 8, 9, 10)
    
    // Distribuir el trabajo entre 3 instancias de sq
    c1 := sq(in)
    c2 := sq(in)
    c3 := sq(in)
    
    // Consumir el canal combinado
    for n := range merge(c1, c2, c3) {
        fmt.Println(n)
    }
}
```

### Context para Control de Cancelación

El paquete `context` proporciona una forma de pasar valores, señales de cancelación y plazos a través de la cadena de llamadas de API.

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Worker: Recibida señal de cancelación")
            return
        default:
            fmt.Println("Worker: Trabajando...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    // Crear un contexto con cancelación
    ctx, cancel := context.WithCancel(context.Background())
    
    // Iniciar worker
    go worker(ctx)
    
    // Trabajar por 2 segundos, luego cancelar
    time.Sleep(2 * time.Second)
    cancel()
    
    // Dar tiempo al worker para que procese la cancelación
    time.Sleep(500 * time.Millisecond)
}
```

### Timeouts y Deadlines

El paquete `context` también permite establecer timeouts y deadlines.

```go
func slowOperation(ctx context.Context) (string, error) {
    // Simulando una operación lenta
    select {
    case <-time.After(2 * time.Second):
        return "Operación completada", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func main() {
    // Crear un contexto con timeout de 1 segundo
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    result, err := slowOperation(ctx)
    if err != nil {
        fmt.Println("Error:", err) // Imprime: Error: context deadline exceeded
    } else {
        fmt.Println("Resultado:", result)
    }
}
```

## Sincronización Avanzada

### sync.Mutex y sync.RWMutex

`sync.Mutex` proporciona exclusión mutua, mientras que `sync.RWMutex` permite múltiples lectores o un único escritor.

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

// Ejemplo con RWMutex
type SafeMap struct {
    mu sync.RWMutex
    data map[string]string
}

func (m *SafeMap) Get(key string) (string, bool) {
    m.mu.RLock() // Bloqueo de lectura (múltiples lectores permitidos)
    defer m.mu.RUnlock()
    val, ok := m.data[key]
    return val, ok
}

func (m *SafeMap) Set(key, value string) {
    m.mu.Lock() // Bloqueo de escritura (exclusivo)
    defer m.mu.Unlock()
    m.data[key] = value
}
```

### sync.Once

`sync.Once` garantiza que una función se ejecute solo una vez, incluso si se llama desde múltiples goroutines.

```go
var instance *Singleton
var once sync.Once

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

### sync.Cond

`sync.Cond` implementa una variable de condición, un punto de espera para goroutines que esperan una condición particular.

```go
type Queue struct {
    cond *sync.Cond
    data []interface{}
    capacity int
}

func NewQueue(capacity int) *Queue {
    return &Queue{
        cond: sync.NewCond(&sync.Mutex{}),
        data: make([]interface{}, 0, capacity),
        capacity: capacity,
    }
}

func (q *Queue) Put(item interface{}) {
    q.cond.L.Lock()
    defer q.cond.L.Unlock()
    
    // Esperar si la cola está llena
    for len(q.data) == q.capacity {
        q.cond.Wait()
    }
    
    // Añadir el elemento
    q.data = append(q.data, item)
    
    // Señalar a los consumidores que hay un nuevo elemento
    q.cond.Signal()
}

func (q *Queue) Get() interface{} {
    q.cond.L.Lock()
    defer q.cond.L.Unlock()
    
    // Esperar si la cola está vacía
    for len(q.data) == 0 {
        q.cond.Wait()
    }
    
    // Obtener el primer elemento
    item := q.data[0]
    q.data = q.data[1:]
    
    // Señalar a los productores que hay espacio disponible
    q.cond.Signal()
    
    return item
}
```

### sync.Pool

`sync.Pool` es un conjunto de objetos temporales que pueden ser guardados y reutilizados, reduciendo la presión sobre el recolector de basura.

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func processRequest(data []byte) {
    // Obtener un buffer del pool
    buffer := bufferPool.Get().(*bytes.Buffer)
    buffer.Reset() // Limpiar el buffer para reutilizarlo
    
    // Usar el buffer
    buffer.Write(data)
    process(buffer.Bytes())
    
    // Devolver el buffer al pool
    bufferPool.Put(buffer)
}
```

## Patrones de Cancelación y Timeout

### Cancelación con Canales

Antes de la introducción del paquete `context`, la cancelación se manejaba típicamente con canales dedicados.

```go
func worker(done <-chan struct{}) {
    for {
        select {
        case <-done:
            fmt.Println("Worker: Cancelado")
            return
        default:
            fmt.Println("Worker: Trabajando...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    done := make(chan struct{})
    
    go worker(done)
    
    time.Sleep(2 * time.Second)
    close(done) // Señal de cancelación
    
    time.Sleep(500 * time.Millisecond)
}
```

### Timeout con Select

```go
func fetchWithTimeout(url string, timeout time.Duration) ([]byte, error) {
    client := http.Client{
        Timeout: timeout,
    }
    
    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    return ioutil.ReadAll(resp.Body)
}

// Alternativa con select y canales
func fetchWithTimeoutChannel(url string, timeout time.Duration) ([]byte, error) {
    result := make(chan []byte, 1)
    errc := make(chan error, 1)
    
    go func() {
        resp, err := http.Get(url)
        if err != nil {
            errc <- err
            return
        }
        defer resp.Body.Close()
        
        data, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            errc <- err
            return
        }
        
        result <- data
    }()
    
    select {
    case data := <-result:
        return data, nil
    case err := <-errc:
        return nil, err
    case <-time.After(timeout):
        return nil, fmt.Errorf("timeout after %v", timeout)
    }
}
```

## Errores en Programación Concurrente

### Race Conditions

Las race conditions ocurren cuando múltiples goroutines acceden a los mismos datos sin sincronización adecuada.

```go
// Ejemplo de race condition
func main() {
    counter := 0
    
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            counter++ // Race condition: múltiples goroutines modifican counter sin sincronización
            wg.Done()
        }()
    }
    
    wg.Wait()
    fmt.Println("Counter:", counter) // El resultado puede no ser 1000
}

// Solución con mutex
func main() {
    counter := 0
    var mu sync.Mutex
    
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            mu.Lock()
            counter++
            mu.Unlock()
            wg.Done()
        }()
    }
    
    wg.Wait()
    fmt.Println("Counter:", counter) // Ahora siempre será 1000
}
```

### Deadlocks

Un deadlock ocurre cuando dos o más goroutines se bloquean mutuamente, esperando que la otra libere un recurso.

```go
func main() {
    ch := make(chan int) // Canal sin buffer
    
    // Deadlock: nadie está recibiendo del canal
    ch <- 1
    
    fmt.Println("Este código nunca se ejecutará")
}

// Solución: usar un canal con buffer o asegurar que alguien reciba
func main() {
    ch := make(chan int, 1) // Canal con buffer de tamaño 1
    
    ch <- 1 // Ahora no bloquea porque hay espacio en el buffer
    
    fmt.Println("Mensaje enviado")
}
```

### Detección de Race Conditions

Go incluye un detector de race conditions que puedes activar con la flag `-race`.

```bash
go run -race main.go
go test -race ./...
```

## Mejores Prácticas

1. **Evita goroutines huérfanas**: Asegúrate de que todas las goroutines terminen correctamente, especialmente en caso de error.

2. **Usa WaitGroups para esperar goroutines**: `sync.WaitGroup` proporciona un mecanismo para esperar a que un grupo de goroutines termine.

3. **Prefiere canales a locks cuando sea posible**: Los canales promueven un estilo de programación más seguro y claro.

4. **Usa el paquete context para cancelación**: Proporciona un mecanismo estándar para propagar cancelación a través de la cadena de llamadas.

5. **Limita el número de goroutines**: Demasiadas goroutines pueden agotar los recursos del sistema.

6. **Usa buffered channels cuando sea apropiado**: Los canales con buffer pueden evitar bloqueos innecesarios.

7. **Cierra los canales desde el lado del remitente**: Solo el remitente debe cerrar un canal, nunca el receptor.

8. **Usa select con default para operaciones no bloqueantes**: Permite realizar otras tareas cuando un canal no está listo.

## Ejercicios Prácticos

1. Implementa un sistema de rate limiting usando goroutines y canales.

2. Crea un pool de conexiones a base de datos con un número máximo de conexiones concurrentes.

3. Implementa un patrón de pipeline para procesar archivos grandes en paralelo.

4. Diseña un sistema de caché concurrente con expiración de elementos.

## Conclusión

La concurrencia en Go es una herramienta poderosa que, cuando se usa correctamente, puede mejorar significativamente el rendimiento y la capacidad de respuesta de tus aplicaciones. Sin embargo, la programación concurrente también introduce desafíos únicos como race conditions, deadlocks y leaks de goroutines.

Dominar los patrones avanzados de concurrencia y las herramientas de sincronización te permitirá escribir código concurrente robusto y eficiente. Recuerda que la simplicidad es clave: usa la concurrencia cuando realmente aporta beneficios y mantén tu diseño lo más simple posible.