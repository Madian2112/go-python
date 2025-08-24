# Pruebas de Rendimiento en Go

## Introducción

Las pruebas de rendimiento (benchmarking) son una parte esencial del desarrollo de software en Go, ya que permiten medir y optimizar el rendimiento de tu código. Go proporciona herramientas integradas para realizar benchmarks como parte de su paquete de testing estándar. En esta sección, exploraremos cómo realizar pruebas de rendimiento efectivas en Go.

## Conceptos Básicos de Benchmarking

Los benchmarks en Go son funciones que:

1. Se definen en archivos con sufijo `_test.go`
2. Comienzan con `Benchmark` (en lugar de `Test`)
3. Reciben un parámetro `*testing.B`
4. Ejecutan el código que se quiere medir dentro de un bucle controlado por el framework

## Escribiendo Benchmarks

### Estructura Básica

```go
// archivo: ejemplo_test.go
package ejemplo

import "testing"

func BenchmarkMiFuncion(b *testing.B) {
    // Código de preparación (no medido)
    
    // Reiniciar el temporizador si la preparación fue costosa
    b.ResetTimer()
    
    // El bucle de benchmark
    for i := 0; i < b.N; i++ {
        // Código que queremos medir
        MiFuncion()
    }
}
```

### Ejecutando Benchmarks

Para ejecutar los benchmarks, usamos el comando `go test` con la bandera `-bench`:

```bash
go test -bench=.
```

Esto ejecutará todos los benchmarks en el paquete actual. Para ejecutar un benchmark específico:

```bash
go test -bench=BenchmarkMiFuncion
```

### Ejemplo Completo

```go
// archivo: fibonacci_test.go
package fibonacci

import "testing"

// Función que queremos medir
func Fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return Fibonacci(n-1) + Fibonacci(n-2)
}

// Versión iterativa más eficiente
func FibonacciIterativo(n int) int {
    if n <= 1 {
        return n
    }
    
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}

// Benchmark para la versión recursiva
func BenchmarkFibonacciRecursivo10(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(10)
    }
}

// Benchmark para la versión iterativa
func BenchmarkFibonacciIterativo10(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FibonacciIterativo(10)
    }
}

// Benchmark con diferentes entradas
func BenchmarkFibonacciIterativo20(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FibonacciIterativo(20)
    }
}
```

Ejecutando estos benchmarks:

```bash
go test -bench=.
```

Podrías obtener resultados como:

```
goos: linux
goarch: amd64
pkg: fibonacci
BenchmarkFibonacciRecursivo10-8      200000              7325 ns/op
BenchmarkFibonacciIterativo10-8     20000000               83.0 ns/op
BenchmarkFibonacciIterativo20-8     10000000              136 ns/op
PASS
ok      fibonacci       5.106s
```

Esto muestra que la versión iterativa es significativamente más rápida que la recursiva.

## Técnicas Avanzadas de Benchmarking

### Medición de Asignaciones de Memoria

Puedes medir cuánta memoria asigna tu código usando la bandera `-benchmem`:

```bash
go test -bench=. -benchmem
```

Esto mostrará información adicional sobre asignaciones de memoria:

```
goos: linux
goarch: amd64
pkg: ejemplo
BenchmarkMiFuncion-8    1000000    1234 ns/op    1234 B/op    12 allocs/op
PASS
ok      ejemplo       1.234s
```

Donde:
- `1234 B/op` indica cuántos bytes se asignan por operación
- `12 allocs/op` indica cuántas asignaciones de memoria ocurren por operación

### Ejemplo de Medición de Memoria

```go
package ejemplo

import "testing"

func CreaSlice() []int {
    return make([]int, 1000)
}

func CreaSlicePreasignado() []int {
    slice := make([]int, 0, 1000)
    for i := 0; i < 1000; i++ {
        slice = append(slice, i)
    }
    return slice
}

func BenchmarkCreaSlice(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = CreaSlice()
    }
}

func BenchmarkCreaSlicePreasignado(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = CreaSlicePreasignado()
    }
}
```

Ejecutando con `-benchmem`:

```bash
go test -bench=. -benchmem
```

Podrías ver que `CreaSlicePreasignado` realiza más asignaciones individuales pero podría ser más eficiente en términos de memoria total.

### Benchmarks con Datos de Entrada

A menudo necesitamos medir el rendimiento con diferentes tamaños de entrada:

```go
func BenchmarkOrdenamiento(b *testing.B) {
    tamaños := []int{10, 100, 1000, 10000}
    
    for _, tamaño := range tamaños {
        b.Run(fmt.Sprintf("tamaño-%d", tamaño), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                // Crear datos de prueba
                datos := generarDatosAleatorios(tamaño)
                
                // Reiniciar el temporizador después de la generación
                b.ResetTimer()
                
                // Medir el ordenamiento
                OrdenarDatos(datos)
            }
        })
    }
}
```

Esto generará resultados separados para cada tamaño de entrada.

### Benchmarks Paralelos

Para medir el rendimiento en código concurrente:

```go
func BenchmarkConcurrente(b *testing.B) {
    // Indicar que este benchmark debe ejecutarse con múltiples goroutines
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Código a medir en paralelo
            MiFuncionConcurrente()
        }
    })
}
```

Puedes controlar el número de goroutines con la bandera `-cpu`:

```bash
go test -bench=. -cpu=1,2,4,8
```

Esto ejecutará los benchmarks con 1, 2, 4 y 8 goroutines, permitiéndote ver cómo escala tu código.

## Comparación de Benchmarks

Go proporciona una herramienta llamada `benchstat` para comparar resultados de benchmarks:

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

Para comparar dos versiones de tu código:

```bash
# Ejecutar benchmarks de la versión actual y guardar resultados
go test -bench=. -count=5 > old.txt

# Hacer cambios en el código

# Ejecutar benchmarks de la nueva versión
go test -bench=. -count=5 > new.txt

# Comparar resultados
benchstat old.txt new.txt
```

Esto mostrará un análisis estadístico de las diferencias de rendimiento.

## Perfilado (Profiling)

Go ofrece herramientas de perfilado para análisis más detallados:

### CPU Profiling

```bash
go test -bench=. -cpuprofile=cpu.prof
```

### Memory Profiling

```bash
go test -bench=. -memprofile=mem.prof
```

### Block Profiling (para detectar contención en concurrencia)

```bash
go test -bench=. -blockprofile=block.prof
```

### Analizando Perfiles

Puedes analizar los perfiles con `pprof`:

```bash
go tool pprof cpu.prof
```

O generar visualizaciones:

```bash
go tool pprof -http=:8080 cpu.prof
```

Esto abrirá una interfaz web con gráficos y análisis detallados.

## Buenas Prácticas

1. **Benchmarks Estables**: Asegúrate de que tus benchmarks sean estables y reproducibles. Evita dependencias externas como E/S de disco o red.

2. **Múltiples Ejecuciones**: Usa `-count=N` para ejecutar cada benchmark múltiples veces y obtener resultados más confiables.

3. **Tamaños Realistas**: Usa tamaños de entrada que reflejen casos de uso reales.

4. **Aislar lo que Mides**: Asegúrate de medir solo el código que te interesa, usando `b.ResetTimer()` después de la preparación.

5. **Evitar Optimizaciones del Compilador**: A veces el compilador puede eliminar código que no tiene efectos observables. Asegúrate de que los resultados de tus cálculos se utilicen de alguna manera.

6. **Comparar Alternativas**: Siempre compara diferentes implementaciones para encontrar la más eficiente.

7. **Medir Memoria**: No solo te enfoques en la velocidad; la eficiencia de memoria es igualmente importante.

8. **Documentar Condiciones**: Documenta las condiciones bajo las cuales se ejecutaron los benchmarks (hardware, versión de Go, etc.).

## Ejemplo Práctico: Optimización de JSON

Vamos a comparar diferentes formas de trabajar con JSON en Go:

```go
package json_ejemplo

import (
    "encoding/json"
    "testing"
)

type Persona struct {
    Nombre   string `json:"nombre"`
    Edad     int    `json:"edad"`
    Direccion string `json:"direccion"`
    Email    string `json:"email"`
    Telefono string `json:"telefono"`
}

// Datos de ejemplo
var jsonData = []byte(`{"nombre":"Juan Pérez","edad":30,"direccion":"Calle Principal 123","email":"juan@ejemplo.com","telefono":"555-1234"}`)

// Usando Marshal/Unmarshal estándar
func BenchmarkJSONUnmarshal(b *testing.B) {
    var p Persona
    for i := 0; i < b.N; i++ {
        _ = json.Unmarshal(jsonData, &p)
    }
}

func BenchmarkJSONMarshal(b *testing.B) {
    p := Persona{
        Nombre:    "Juan Pérez",
        Edad:      30,
        Direccion: "Calle Principal 123",
        Email:     "juan@ejemplo.com",
        Telefono:  "555-1234",
    }
    
    for i := 0; i < b.N; i++ {
        _, _ = json.Marshal(p)
    }
}

// Usando Decoder/Encoder
func BenchmarkJSONDecoder(b *testing.B) {
    var p Persona
    for i := 0; i < b.N; i++ {
        decoder := json.NewDecoder(bytes.NewReader(jsonData))
        _ = decoder.Decode(&p)
    }
}

func BenchmarkJSONEncoder(b *testing.B) {
    p := Persona{
        Nombre:    "Juan Pérez",
        Edad:      30,
        Direccion: "Calle Principal 123",
        Email:     "juan@ejemplo.com",
        Telefono:  "555-1234",
    }
    
    for i := 0; i < b.N; i++ {
        var buf bytes.Buffer
        encoder := json.NewEncoder(&buf)
        _ = encoder.Encode(p)
    }
}
```

Ejecutando estos benchmarks con `-benchmem` podríamos descubrir que `Decoder`/`Encoder` pueden ser más eficientes en términos de asignaciones de memoria, especialmente para JSON grandes o streams de datos.

## Optimización Basada en Benchmarks

Veamos un ejemplo de cómo usar benchmarks para guiar la optimización:

```go
package strings_ejemplo

import (
    "strings"
    "testing"
)

// Versión inicial: concatenación con +
func ConcatenaStrings(elementos []string) string {
    resultado := ""
    for _, elem := range elementos {
        resultado += elem
    }
    return resultado
}

// Versión optimizada: usando strings.Builder
func ConcatenaStringsOptimizado(elementos []string) string {
    var builder strings.Builder
    for _, elem := range elementos {
        builder.WriteString(elem)
    }
    return builder.String()
}

// Versión con capacidad preasignada
func ConcatenaStringsPreasignado(elementos []string) string {
    // Calcular la longitud total
    longitud := 0
    for _, elem := range elementos {
        longitud += len(elem)
    }
    
    // Preasignar capacidad
    var builder strings.Builder
    builder.Grow(longitud)
    
    for _, elem := range elementos {
        builder.WriteString(elem)
    }
    return builder.String()
}

// Benchmarks
func BenchmarkConcatenacion(b *testing.B) {
    elementos := []string{"Hola", ", ", "mundo", "!", " ", "Esto", " ", "es", " ", "una", " ", "prueba", "."}
    
    b.Run("Concatenacion+", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = ConcatenaStrings(elementos)
        }
    })
    
    b.Run("Builder", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = ConcatenaStringsOptimizado(elementos)
        }
    })
    
    b.Run("BuilderPreasignado", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = ConcatenaStringsPreasignado(elementos)
        }
    })
}
```

Ejecutando este benchmark con `-benchmem` probablemente mostraría:

1. La concatenación con `+` es significativamente más lenta y realiza muchas más asignaciones de memoria
2. `strings.Builder` es mucho más eficiente
3. Preasignar la capacidad ofrece una mejora adicional al reducir las reasignaciones

## Microbenchmarks vs. Benchmarks Realistas

Es importante entender la diferencia entre microbenchmarks (que miden funciones aisladas) y benchmarks más realistas (que miden el rendimiento de sistemas completos):

### Microbenchmarks

Útiles para:
- Comparar algoritmos alternativos
- Optimizar funciones críticas
- Entender el costo de operaciones específicas

Limitaciones:
- Pueden no reflejar el rendimiento en un sistema real
- El compilador puede optimizar código aislado de formas que no ocurrirían en un sistema completo
- No capturan efectos de interacción entre componentes

### Benchmarks Realistas

Para benchmarks más realistas:

```go
func BenchmarkSistemaCompleto(b *testing.B) {
    // Configurar un entorno realista
    servidor := ConfigurarServidor()
    cliente := ConfigurarCliente()
    
    // Preparar datos de prueba realistas
    datos := PrepararDatosPrueba()
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        // Simular una operación completa del sistema
        cliente.EnviarSolicitud(servidor, datos)
    }
}
```

## Herramientas Adicionales

### go-torch

`go-torch` genera visualizaciones de llama (flame graphs) a partir de perfiles de CPU:

```bash
go-torch -f cpu.prof
```

### benchcmp (obsoleto, usar benchstat)

`benchcmp` era una herramienta para comparar resultados de benchmarks, pero ha sido reemplazada por `benchstat`.

### trace

Go incluye una herramienta de trazado para analizar la ejecución de programas:

```bash
go test -bench=. -trace=trace.out
go tool trace trace.out
```

Esto abre una interfaz web con visualizaciones detalladas de la ejecución, incluyendo goroutines, GC, y eventos del sistema.

## Conclusión

Las pruebas de rendimiento son una parte esencial del desarrollo en Go, especialmente para aplicaciones donde el rendimiento es crítico. Go proporciona herramientas excelentes para medir y optimizar el rendimiento, desde benchmarks simples hasta análisis detallados con perfilado.

Recuerda que la optimización prematura puede ser contraproducente. Primero escribe código claro y correcto, luego usa benchmarks para identificar cuellos de botella reales, y finalmente optimiza esas partes específicas.

## Recursos Adicionales

- [Documentación oficial de testing en Go](https://golang.org/pkg/testing/)
- [Perfilado en Go](https://blog.golang.org/pprof)
- [Herramientas de análisis de rendimiento en Go](https://golang.org/doc/diagnostics.html)
- [Optimización de rendimiento en Go](https://github.com/dgryski/go-perfbook)
- [Patrones de alto rendimiento en Go](https://dave.cheney.net/high-performance-go-workshop/gophercon-2019.html)

---

Con esto concluimos nuestra exploración de las pruebas de rendimiento en Go. Estas herramientas y técnicas te permitirán medir y optimizar el rendimiento de tus aplicaciones Go de manera efectiva.