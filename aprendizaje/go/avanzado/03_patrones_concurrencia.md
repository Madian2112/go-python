# Patrones Avanzados de Concurrencia en Go

## Introducción

Go es conocido por su excelente soporte para la concurrencia a través de goroutines y canales. Sin embargo, más allá de los conceptos básicos, existen patrones avanzados que permiten resolver problemas complejos de concurrencia de manera elegante y eficiente. En este módulo, exploraremos estos patrones avanzados, sus aplicaciones y cómo implementarlos correctamente.

La concurrencia avanzada en Go va más allá de simplemente lanzar goroutines y comunicarse a través de canales. Implica entender cómo estructurar programas concurrentes complejos, gestionar recursos compartidos, coordinar múltiples goroutines, y manejar errores en entornos concurrentes.

## Patrones de Orquestación

### El patrón Generator

El patrón Generator permite crear flujos de datos (streams) que pueden ser consumidos por otras partes del programa. Es especialmente útil para procesar grandes conjuntos de datos de manera eficiente.

```go
func generator(nums ...int) <-chan int {
    out := make(chan int)
    
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }    
    }()
    
    return out
}

// Uso
func main() {
    for n := range generator(1, 2, 3, 4, 5) {
        fmt.Println(n)
    }
}
```

### El patrón Fan-Out, Fan-In

Este patrón permite distribuir trabajo entre múltiples goroutines (fan-out) y luego combinar los resultados (fan-in). Es ideal para paralelizar tareas independientes.

```go
func fanOut(in <-chan int, n int) []<-chan int {
    channels := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        channels[i] = worker(in)
    }
    return channels
}

func worker(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            // Procesar n (por ejemplo, calcular el cuadrado)
            result := n * n
            out <- result
        }
    }()
    return out
}

func fanIn(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    multiplexed := make(chan int)
    
    // Función para reenviar valores de un canal al canal multiplexado
    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            multiplexed <- n
        }
    }
    
    // Establecer el contador de espera y lanzar goroutines
    wg.Add(len(channels))
    for _, c := range channels {
        go output(c)
    }
    
    // Cerrar el canal multiplexado cuando todos los canales de entrada estén cerrados
    go func() {
        wg.Wait()
        close(multiplexed)
    }()
    
    return multiplexed
}

// Uso
func main() {
    in := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    workers := fanOut(in, 3) // Distribuir trabajo entre 3 workers
    results := fanIn(workers...)
    
    for result := range results {
        fmt.Println(result)
    }
}
```

### El patrón Pipeline

El patrón Pipeline permite construir una serie de etapas de procesamiento donde la salida de una etapa es la entrada de la siguiente. Es útil para dividir tareas complejas en pasos más pequeños y manejables.

```go
func stage1(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * 2 // Duplicar
        }
    }()
    return out
}

func stage2(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n + 1 // Incrementar
        }
    }()
    return out
}

func stage3(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n // Elevar al cuadrado
        }
    }()
    return out
}

// Uso
func main() {
    in := generator(1, 2, 3, 4, 5)
    pipeline := stage3(stage2(stage1(in)))
    
    for result := range pipeline {
        fmt.Println(result)
    }
}
```

## Patrones de Cancelación y Timeout

### Cancelación con Context

Go proporciona el paquete `context` para manejar cancelaciones, timeouts y valores que atraviesan límites de API. Es especialmente útil para operaciones que pueden tardar mucho tiempo o que necesitan ser canceladas.

```go
func worker(ctx context.Context) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 0; ; i++ {
            select {
            case <-ctx.Done():
                fmt.Println("Worker cancelado:", ctx.Err())
                return
            case out <- i:
                time.Sleep(100 * time.Millisecond)
            }
        }
    }()
    return out
}

// Uso con cancelación manual
func manualCancellation() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // Asegurarse de que se llame a cancel
    
    results := worker(ctx)
    
    // Consumir algunos resultados
    for i := 0; i < 5; i++ {
        fmt.Println(<-results)
    }
    
    // Cancelar el worker
    cancel()
    
    // Dar tiempo para que el worker se cierre
    time.Sleep(200 * time.Millisecond)
}

// Uso con timeout
func timeoutCancellation() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    
    results := worker(ctx)
    
    // Consumir resultados hasta que se agote el tiempo
    for result := range results {
        fmt.Println(result)
    }
}
```

### Cancelación con Canales

Antes de la introducción del paquete `context`, la cancelación se manejaba típicamente con canales dedicados. Este enfoque aún puede ser útil en algunos casos.

```go
func workerWithChan(done <-chan struct{}) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 0; ; i++ {
            select {
            case <-done:
                fmt.Println("Worker cancelado con canal")
                return
            case out <- i:
                time.Sleep(100 * time.Millisecond)
            }
        }
    }()
    return out
}

// Uso
func main() {
    done := make(chan struct{})
    results := workerWithChan(done)
    
    // Consumir algunos resultados
    for i := 0; i < 5; i++ {
        fmt.Println(<-results)
    }
    
    // Cancelar el worker
    close(done)
    
    // Dar tiempo para que el worker se cierre
    time.Sleep(200 * time.Millisecond)
}
```

## Patrones de Sincronización

### El patrón Semáforo

Los semáforos permiten controlar el acceso a recursos limitados. En Go, se pueden implementar usando canales con buffer.

```go
type Semaphore chan struct{}

func NewSemaphore(maxConcurrency int) Semaphore {
    return make(Semaphore, maxConcurrency)
}

func (s Semaphore) Acquire() {
    s <- struct{}{}
}

func (s Semaphore) Release() {
    <-s
}

// Uso
func main() {
    sem := NewSemaphore(3) // Máximo 3 operaciones concurrentes
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            sem.Acquire()
            defer sem.Release()
            
            // Operación que requiere recursos limitados
            fmt.Printf("Worker %d iniciado\n", id)
            time.Sleep(1 * time.Second)
            fmt.Printf("Worker %d finalizado\n", id)
        }(i)
    }
    
    wg.Wait()
}
```

### El patrón Worker Pool

Un Worker Pool es un conjunto de workers que procesan tareas de una cola. Es útil para limitar la concurrencia y gestionar eficientemente los recursos.

```go
type Task struct {
    ID  int
    Job func() interface{}
}

type Result struct {
    TaskID int
    Value  interface{}
}

func workerPool(numWorkers int, tasks <-chan Task, results chan<- Result) {
    var wg sync.WaitGroup
    wg.Add(numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        go func(workerID int) {
            defer wg.Done()
            for task := range tasks {
                fmt.Printf("Worker %d procesando tarea %d\n", workerID, task.ID)
                value := task.Job()
                results <- Result{TaskID: task.ID, Value: value}
            }
        }(i)
    }
    
    // Esperar a que todos los workers terminen y cerrar el canal de resultados
    go func() {
        wg.Wait()
        close(results)
    }()
}

// Uso
func main() {
    numTasks := 20
    numWorkers := 5
    
    tasks := make(chan Task, numTasks)
    results := make(chan Result, numTasks)
    
    // Iniciar el worker pool
    workerPool(numWorkers, tasks, results)
    
    // Enviar tareas
    for i := 0; i < numTasks; i++ {
        taskID := i
        tasks <- Task{
            ID: taskID,
            Job: func() interface{} {
                // Simular trabajo
                time.Sleep(100 * time.Millisecond)
                return taskID * 2
            },
        }
    }
    close(tasks) // No más tareas
    
    // Recoger resultados
    for result := range results {
        fmt.Printf("Resultado de tarea %d: %v\n", result.TaskID, result.Value)
    }
}
```

## Patrones de Manejo de Errores

### Propagación de Errores en Goroutines

Manejar errores en código concurrente puede ser complicado. Un enfoque común es incluir el error en la estructura de resultado.

```go
type Result struct {
    Value interface{}
    Err   error
}

func worker(job int) <-chan Result {
    resultChan := make(chan Result, 1)
    go func() {
        defer close(resultChan)
        // Simular trabajo que puede fallar
        if job%3 == 0 {
            resultChan <- Result{Err: fmt.Errorf("error en trabajo %d", job)}
            return
        }
        // Trabajo exitoso
        time.Sleep(100 * time.Millisecond)
        resultChan <- Result{Value: job * 2}
    }()
    return resultChan
}

// Uso
func main() {
    for i := 0; i < 10; i++ {
        result := <-worker(i)
        if result.Err != nil {
            fmt.Printf("Error: %v\n", result.Err)
        } else {
            fmt.Printf("Valor: %v\n", result.Value)
        }
    }
}
```

### Errores en Pipelines

En pipelines, los errores deben propagarse a través de las etapas. Una forma de hacerlo es incluir el error en la estructura que se pasa por el canal.

```go
type Item struct {
    Value int
    Err   error
}

func stage1(in <-chan int) <-chan Item {
    out := make(chan Item)
    go func() {
        defer close(out)
        for n := range in {
            if n < 0 {
                out <- Item{Err: fmt.Errorf("número negativo: %d", n)}
                continue
            }
            out <- Item{Value: n * 2}
        }
    }()
    return out
}

func stage2(in <-chan Item) <-chan Item {
    out := make(chan Item)
    go func() {
        defer close(out)
        for item := range in {
            // Propagar error si existe
            if item.Err != nil {
                out <- item
                continue
            }
            
            // Procesar valor
            if item.Value > 10 {
                out <- Item{Err: fmt.Errorf("valor demasiado grande: %d", item.Value)}
                continue
            }
            out <- Item{Value: item.Value + 1}
        }
    }()
    return out
}

// Uso
func main() {
    in := generator(1, 2, -3, 4, 6)
    pipeline := stage2(stage1(in))
    
    for item := range pipeline {
        if item.Err != nil {
            fmt.Printf("Error: %v\n", item.Err)
        } else {
            fmt.Printf("Valor: %v\n", item.Value)
        }
    }
}
```

## Patrones de Composición

### Composición de Contextos

Los contextos pueden componerse para combinar diferentes comportamientos, como timeouts y valores.

```go
func composedContexts() {
    // Contexto base
    rootCtx := context.Background()
    
    // Añadir un timeout
    timeoutCtx, cancelTimeout := context.WithTimeout(rootCtx, 5*time.Second)
    defer cancelTimeout()
    
    // Añadir un valor
    valueCtx := context.WithValue(timeoutCtx, "userID", 123)
    
    // Añadir una cancelación manual
    finalCtx, cancelFinal := context.WithCancel(valueCtx)
    defer cancelFinal()
    
    // Usar el contexto compuesto
    go func(ctx context.Context) {
        select {
        case <-ctx.Done():
            fmt.Println("Contexto cancelado:", ctx.Err())
            return
        case <-time.After(1 * time.Second):
            // Acceder al valor
            userID := ctx.Value("userID")
            fmt.Println("UserID del contexto:", userID)
            
            // Cancelar manualmente
            cancelFinal()
        }
    }(finalCtx)
    
    // Esperar a que la goroutine termine
    time.Sleep(2 * time.Second)
}
```

### Composición de Canales

Los canales pueden componerse para crear comportamientos complejos, como multiplexación, demultiplexación, y más.

```go
// Multiplexar múltiples canales en uno
func merge(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(channels))
    
    for _, c := range channels {
        go func(ch <-chan int) {
            defer wg.Done()
            for n := range ch {
                out <- n
            }
        }(c)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

// Dividir un canal en múltiples canales según una condición
func split(in <-chan int, predicate func(int) bool) (trueCh, falseCh <-chan int) {
    outTrue := make(chan int)
    outFalse := make(chan int)
    
    go func() {
        defer close(outTrue)
        defer close(outFalse)
        
        for n := range in {
            if predicate(n) {
                outTrue <- n
            } else {
                outFalse <- n
            }
        }
    }()
    
    return outTrue, outFalse
}

// Uso
func main() {
    // Crear tres generadores
    gen1 := generator(1, 4, 7)
    gen2 := generator(2, 5, 8)
    gen3 := generator(3, 6, 9)
    
    // Multiplexar los generadores
    merged := merge(gen1, gen2, gen3)
    
    // Dividir según si es par o impar
    isEven := func(n int) bool { return n%2 == 0 }
    evenCh, oddCh := split(merged, isEven)
    
    // Procesar pares e impares por separado
    go func() {
        for n := range evenCh {
            fmt.Printf("Par: %d\n", n)
        }
    }()
    
    for n := range oddCh {
        fmt.Printf("Impar: %d\n", n)
    }
}
```

## Patrones de Rendimiento

### Buffering y Batching

El buffering y el batching pueden mejorar significativamente el rendimiento al reducir la sobrecarga de comunicación.

```go
// Canal con buffer para reducir bloqueos
func bufferedGenerator(nums ...int) <-chan int {
    out := make(chan int, len(nums)) // Buffer del tamaño de los datos
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

// Procesamiento por lotes para reducir la sobrecarga
func batchProcessor(in <-chan int, batchSize int) <-chan []int {
    out := make(chan []int)
    go func() {
        defer close(out)
        batch := make([]int, 0, batchSize)
        
        for n := range in {
            batch = append(batch, n)
            
            // Enviar lote cuando está lleno
            if len(batch) == batchSize {
                out <- batch
                batch = make([]int, 0, batchSize)
            }
        }
        
        // Enviar el último lote parcial si existe
        if len(batch) > 0 {
            out <- batch
        }
    }()
    return out
}

// Uso
func main() {
    // Generar números del 1 al 100
    nums := make([]int, 100)
    for i := range nums {
        nums[i] = i + 1
    }
    
    in := bufferedGenerator(nums...)
    batches := batchProcessor(in, 10)
    
    for batch := range batches {
        fmt.Printf("Procesando lote: %v\n", batch)
        // Simular procesamiento por lotes
        time.Sleep(100 * time.Millisecond)
    }
}
```

### Paralelismo Controlado

Controlar el grado de paralelismo es crucial para optimizar el rendimiento sin agotar los recursos del sistema.

```go
func parallelProcessor(in <-chan int, numWorkers int) <-chan int {
    out := make(chan int)
    
    // Crear un grupo de workers
    var wg sync.WaitGroup
    wg.Add(numWorkers)
    
    // Función worker
    worker := func() {
        defer wg.Done()
        for n := range in {
            // Simular procesamiento intensivo
            result := processItem(n)
            out <- result
        }
    }
    
    // Iniciar workers
    for i := 0; i < numWorkers; i++ {
        go worker()
    }
    
    // Cerrar canal de salida cuando todos los workers terminen
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

func processItem(n int) int {
    // Simular trabajo intensivo en CPU
    time.Sleep(100 * time.Millisecond)
    return n * n
}

// Uso
func main() {
    // Número óptimo de workers basado en CPUs disponibles
    numCPU := runtime.NumCPU()
    fmt.Printf("Usando %d workers (CPUs disponibles)\n", numCPU)
    
    // Generar trabajo
    in := make(chan int, 100)
    go func() {
        defer close(in)
        for i := 1; i <= 100; i++ {
            in <- i
        }
    }()
    
    // Procesar en paralelo
    results := parallelProcessor(in, numCPU)
    
    // Recoger resultados
    count := 0
    startTime := time.Now()
    for range results {
        count++
    }
    elapsed := time.Since(startTime)
    
    fmt.Printf("Procesados %d items en %v\n", count, elapsed)
}
```

## Patrones de Coordinación

### Barreras de Sincronización

Las barreras permiten que múltiples goroutines esperen en un punto específico hasta que todas lleguen allí.

```go
type Barrier struct {
    count    int
    mutex    sync.Mutex
    cond     *sync.Cond
    threshold int
}

func NewBarrier(threshold int) *Barrier {
    b := &Barrier{
        threshold: threshold,
    }
    b.cond = sync.NewCond(&b.mutex)
    return b
}

func (b *Barrier) Wait() {
    b.mutex.Lock()
    defer b.mutex.Unlock()
    
    b.count++
    if b.count == b.threshold {
        // Último en llegar, resetear y despertar a todos
        b.count = 0
        b.cond.Broadcast()
        return
    }
    
    // Esperar a que todos lleguen
    b.cond.Wait()
}

// Uso
func main() {
    numWorkers := 5
    barrier := NewBarrier(numWorkers)
    var wg sync.WaitGroup
    wg.Add(numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        go func(id int) {
            defer wg.Done()
            
            // Fase 1
            fmt.Printf("Worker %d completando fase 1\n", id)
            time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
            
            fmt.Printf("Worker %d esperando en barrera 1\n", id)
            barrier.Wait()
            fmt.Printf("Worker %d pasó barrera 1\n", id)
            
            // Fase 2
            fmt.Printf("Worker %d completando fase 2\n", id)
            time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
            
            fmt.Printf("Worker %d esperando en barrera 2\n", id)
            barrier.Wait()
            fmt.Printf("Worker %d pasó barrera 2\n", id)
            
            // Fase 3
            fmt.Printf("Worker %d completando fase 3\n", id)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("Todos los workers completaron todas las fases")
}
```

### Rendezvous

El patrón Rendezvous permite que dos goroutines se sincronicen en un punto específico antes de continuar.

```go
type Rendezvous struct {
    first  chan struct{}
    second chan struct{}
}

func NewRendezvous() *Rendezvous {
    return &Rendezvous{
        first:  make(chan struct{}),
        second: make(chan struct{}),
    }
}

func (r *Rendezvous) FirstArrives() {
    // Señalar llegada
    close(r.first)
    // Esperar a que el segundo llegue
    <-r.second
}

func (r *Rendezvous) SecondArrives() {
    // Esperar a que el primero llegue
    <-r.first
    // Señalar llegada
    close(r.second)
}

// Uso
func main() {
    r := NewRendezvous()
    
    go func() {
        fmt.Println("Goroutine 1: Iniciando trabajo")
        time.Sleep(2 * time.Second)
        fmt.Println("Goroutine 1: Llegando al punto de encuentro")
        r.FirstArrives()
        fmt.Println("Goroutine 1: Continuando después del encuentro")
    }()
    
    go func() {
        fmt.Println("Goroutine 2: Iniciando trabajo")
        time.Sleep(1 * time.Second)
        fmt.Println("Goroutine 2: Llegando al punto de encuentro")
        r.SecondArrives()
        fmt.Println("Goroutine 2: Continuando después del encuentro")
    }()
    
    // Esperar a que ambas goroutines terminen
    time.Sleep(3 * time.Second)
}
```

## Mejores Prácticas

### Evitar Goroutine Leaks

Una goroutine leak ocurre cuando una goroutine nunca termina, lo que puede llevar a un consumo excesivo de memoria. Es importante asegurarse de que todas las goroutines terminen correctamente.

```go
// Mal: Potencial goroutine leak
func leakyFunction() {
    ch := make(chan int)
    go func() {
        // Esta goroutine nunca terminará si nadie lee del canal
        for i := 0; ; i++ {
            ch <- i
        }
    }()
    
    // Solo leemos una vez y luego descartamos el canal
    fmt.Println(<-ch)
    // La goroutine queda bloqueada para siempre
}

// Bien: Uso de contexto para cancelación
func nonLeakyFunction() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // Asegura que la goroutine termine
    
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; ; i++ {
            select {
            case <-ctx.Done():
                return
            case ch <- i:
                // Continuar enviando
            }
        }
    }()
    
    // Leer una vez
    fmt.Println(<-ch)
    // Al salir de la función, cancel() será llamado y la goroutine terminará
}
```

### Manejo de Recursos

Es importante gestionar adecuadamente los recursos en entornos concurrentes, especialmente cuando se comparten entre múltiples goroutines.

```go
type ResourceManager struct {
    resources map[string]*Resource
    mutex     sync.RWMutex
}

type Resource struct {
    data []byte
    // otros campos
}

func NewResourceManager() *ResourceManager {
    return &ResourceManager{
        resources: make(map[string]*Resource),
    }
}

func (rm *ResourceManager) Get(id string) (*Resource, bool) {
    rm.mutex.RLock()
    defer rm.mutex.RUnlock()
    
    resource, exists := rm.resources[id]
    return resource, exists
}

func (rm *ResourceManager) Add(id string, resource *Resource) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()
    
    rm.resources[id] = resource
}

func (rm *ResourceManager) Remove(id string) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()
    
    delete(rm.resources, id)
}

// Uso
func main() {
    rm := NewResourceManager()
    var wg sync.WaitGroup
    
    // Múltiples goroutines accediendo a recursos compartidos
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            // Crear y añadir recurso
            resourceID := fmt.Sprintf("resource-%d", id)
            rm.Add(resourceID, &Resource{data: []byte(fmt.Sprintf("data-%d", id))})
            
            // Simular trabajo
            time.Sleep(100 * time.Millisecond)
            
            // Obtener recurso
            if resource, exists := rm.Get(resourceID); exists {
                fmt.Printf("Recurso %s: %s\n", resourceID, resource.data)
            }
            
            // Eliminar recurso
            rm.Remove(resourceID)
        }(i)
    }
    
    wg.Wait()
}
```

### Pruebas de Código Concurrente

Probar código concurrente puede ser desafiante. Aquí hay algunas técnicas para hacerlo de manera efectiva.

```go
// Función a probar
func ProcessItems(items []int) <-chan int {
    results := make(chan int)
    
    go func() {
        defer close(results)
        for _, item := range items {
            results <- item * 2
        }
    }()
    
    return results
}

// Prueba
func TestProcessItems(t *testing.T) {
    // Datos de prueba
    items := []int{1, 2, 3, 4, 5}
    expected := []int{2, 4, 6, 8, 10}
    
    // Llamar a la función
    resultChan := ProcessItems(items)
    
    // Recoger resultados
    var results []int
    for r := range resultChan {
        results = append(results, r)
    }
    
    // Verificar resultados
    if len(results) != len(expected) {
        t.Errorf("Expected %d results, got %d", len(expected), len(results))
    }
    
    // Verificar cada resultado
    for i, v := range results {
        if v != expected[i] {
            t.Errorf("Result %d: expected %d, got %d", i, expected[i], v)
        }
    }
}

// Prueba con timeout
func TestProcessItemsWithTimeout(t *testing.T) {
    items := []int{1, 2, 3, 4, 5}
    resultChan := ProcessItems(items)
    
    // Usar un timeout para evitar que la prueba se bloquee
    timeout := time.After(1 * time.Second)
    count := 0
    
    for count < len(items) {
        select {
        case _, ok := <-resultChan:
            if !ok {
                t.Fatal("Channel closed prematurely")
            }
            count++
        case <-timeout:
            t.Fatal("Test timed out")
        }
    }
    
    // Verificar que el canal se cierre correctamente
    select {
    case _, ok := <-resultChan:
        if ok {
            t.Error("Channel should be closed")
        }
    case <-timeout:
        t.Fatal("Test timed out waiting for channel to close")
    }
}
```

## Ejercicios Prácticos

1. **Implementar un sistema de caché concurrente**:
   - Crear una caché que permita múltiples lecturas concurrentes pero bloquee durante escrituras.
   - Implementar expiración de entradas.
   - Añadir estadísticas de hit/miss.

2. **Crear un servidor de chat con salas**:
   - Implementar un servidor que maneje múltiples clientes concurrentemente.
   - Permitir que los clientes se unan a diferentes salas de chat.
   - Implementar broadcast de mensajes dentro de cada sala.

3. **Desarrollar un crawler web concurrente**:
   - Implementar un crawler que visite URLs concurrentemente.
   - Limitar el número máximo de goroutines activas.
   - Implementar detección de ciclos para evitar visitar la misma URL múltiples veces.
   - Añadir timeouts para evitar bloqueos en sitios lentos.

4. **Implementar un sistema de procesamiento de datos en pipeline**:
   - Crear varias etapas de procesamiento (lectura, transformación, filtrado, agregación, escritura).
   - Conectar las etapas mediante canales.
   - Implementar manejo de errores a lo largo del pipeline.
   - Añadir capacidad para cancelar todo el pipeline.

5. **Crear un pool de conexiones a base de datos**:
   - Implementar un pool que gestione un número limitado de conexiones.
   - Permitir que múltiples goroutines soliciten y liberen conexiones.
   - Implementar timeouts para solicitudes de conexión.
   - Añadir health checks para las conexiones.

## Conclusiones

Los patrones avanzados de concurrencia en Go proporcionan herramientas poderosas para resolver problemas complejos de manera elegante y eficiente. Al dominar estos patrones, los desarrolladores pueden crear aplicaciones concurrentes robustas, escalables y de alto rendimiento.

Recuerda que la concurrencia no es siempre la respuesta. A veces, un enfoque secuencial simple puede ser más claro y menos propenso a errores. Usa la concurrencia cuando realmente aporta beneficios, como mejorar el rendimiento, la capacidad de respuesta o la estructura del código.

Finalmente, siempre presta atención a los posibles problemas de concurrencia como race conditions, deadlocks y goroutine leaks. Usa las herramientas proporcionadas por Go, como el race detector, para identificar y solucionar estos problemas.

## Referencias

1. Donovan, Alan A. A., & Kernighan, Brian W. (2015). The Go Programming Language. Addison-Wesley Professional.
2. Pike, Rob. (2012). Go Concurrency Patterns. Google I/O.
3. Cox-Buday, Katherine. (2017). Concurrency in Go: Tools and Techniques for Developers. O'Reilly Media.
4. Go Blog: https://blog.golang.org/pipelines
5. Go Blog: https://blog.golang.org/context
6. Go Documentation: https://golang.org/pkg/context/
7. Go Documentation: https://golang.org/pkg/sync/