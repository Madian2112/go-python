# Interfaces Avanzadas en Go

## Introducción

Las interfaces son uno de los conceptos más poderosos en Go. A diferencia de otros lenguajes orientados a objetos, Go implementa un sistema de tipos basado en la composición en lugar de la herencia, y las interfaces juegan un papel fundamental en este enfoque.

En este módulo, profundizaremos en el uso avanzado de interfaces, explorando patrones y técnicas que te permitirán escribir código más flexible, mantenible y reutilizable.

## Repaso de Interfaces Básicas

Antes de adentrarnos en conceptos avanzados, recordemos los fundamentos:

```go
// Definición de una interfaz
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Implementación implícita
type File struct {
    // campos
}

// File implementa Reader
func (f *File) Read(p []byte) (n int, err error) {
    // implementación
    return len(p), nil
}
```

En Go, las interfaces se implementan implícitamente. No hay una palabra clave `implements` como en otros lenguajes. Si un tipo tiene todos los métodos declarados por una interfaz, entonces ese tipo implementa la interfaz.

## Interfaces Vacías y Type Assertions

La interfaz vacía `interface{}` (o `any` en Go 1.18+) no especifica ningún método, por lo que todos los tipos implementan esta interfaz.

```go
func PrintAny(v interface{}) {
    fmt.Println(v)
}
```

Para recuperar el valor concreto de una interfaz, usamos type assertions o type switches:

```go
// Type assertion
value, ok := v.(string)
if ok {
    fmt.Println("Es una cadena:", value)
} else {
    fmt.Println("No es una cadena")
}

// Type switch
switch val := v.(type) {
case string:
    fmt.Println("Es una cadena:", val)
case int:
    fmt.Println("Es un entero:", val)
default:
    fmt.Println("Tipo desconocido")
}
```

## Composición de Interfaces

Una práctica poderosa en Go es la composición de interfaces más pequeñas para formar interfaces más complejas:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Composición de interfaces
type ReadWriter interface {
    Reader
    Writer
}
```

Este enfoque promueve la creación de interfaces pequeñas y enfocadas que pueden combinarse según sea necesario.

## Interfaces como Contratos

Las interfaces en Go actúan como contratos entre diferentes partes del código. Definen comportamientos esperados sin especificar implementaciones concretas.

```go
type PaymentProcessor interface {
    Process(amount float64) error
    Refund(transactionID string, amount float64) error
}

// Diferentes implementaciones
type CreditCardProcessor struct {}
type PayPalProcessor struct {}
type CryptoProcessor struct {}

// Cada una implementa PaymentProcessor a su manera
func (cc *CreditCardProcessor) Process(amount float64) error {
    // Implementación específica para tarjetas de crédito
    return nil
}

func (cc *CreditCardProcessor) Refund(transactionID string, amount float64) error {
    // Implementación de reembolso para tarjetas de crédito
    return nil
}

// Implementaciones similares para PayPalProcessor y CryptoProcessor
```

## Interfaces Paramétricas (Go 1.18+)

Con la introducción de generics en Go 1.18, podemos crear interfaces paramétricas:

```go
type Comparable[T any] interface {
    CompareTo(other T) int
}

type SortedCollection[T Comparable[T]] struct {
    items []T
}

func (sc *SortedCollection[T]) Add(item T) {
    // Lógica para insertar manteniendo el orden
}
```

## Patrones Avanzados con Interfaces

### Patrón Decorator

El patrón Decorator permite añadir comportamiento a objetos individuales dinámicamente sin afectar el comportamiento de otros objetos de la misma clase.

```go
type DataSource interface {
    Read() ([]byte, error)
    Write(data []byte) error
}

type FileDataSource struct {
    filename string
}

func (fds *FileDataSource) Read() ([]byte, error) {
    // Leer desde archivo
    return ioutil.ReadFile(fds.filename)
}

func (fds *FileDataSource) Write(data []byte) error {
    // Escribir en archivo
    return ioutil.WriteFile(fds.filename, data, 0644)
}

// Decorator para encriptar datos
type EncryptionDecorator struct {
    source DataSource
    key    []byte
}

func (ed *EncryptionDecorator) Read() ([]byte, error) {
    data, err := ed.source.Read()
    if err != nil {
        return nil, err
    }
    // Desencriptar datos
    return decrypt(data, ed.key), nil
}

func (ed *EncryptionDecorator) Write(data []byte) error {
    // Encriptar datos
    encryptedData := encrypt(data, ed.key)
    return ed.source.Write(encryptedData)
}

// Uso
func main() {
    source := &FileDataSource{filename: "data.txt"}
    encryptedSource := &EncryptionDecorator{
        source: source,
        key:    []byte("secret-key"),
    }
    
    // Ahora podemos usar encryptedSource como un DataSource normal
    // pero con encriptación automática
    data := []byte("datos sensibles")
    encryptedSource.Write(data)
    
    readData, _ := encryptedSource.Read()
    fmt.Println(string(readData)) // Imprime "datos sensibles"
}
```

### Patrón Adapter

El patrón Adapter permite que interfaces incompatibles trabajen juntas.

```go
// Interfaz existente
type LegacyPrinter interface {
    Print(s string) string
}

// Implementación existente
type OldPrinter struct{}

func (op *OldPrinter) Print(s string) string {
    return "Legacy Printer: " + s
}

// Nueva interfaz
type ModernPrinter interface {
    PrintModern(s string) string
}

// Adaptador
type PrinterAdapter struct {
    OldPrinter LegacyPrinter
}

func (pa *PrinterAdapter) PrintModern(s string) string {
    return pa.OldPrinter.Print(s)
}

// Cliente que espera ModernPrinter
func ClientCode(printer ModernPrinter) {
    fmt.Println(printer.PrintModern("Hello World!"))
}

func main() {
    // Usando el adaptador
    oldPrinter := &OldPrinter{}
    adapter := &PrinterAdapter{OldPrinter: oldPrinter}
    
    ClientCode(adapter) // Imprime "Legacy Printer: Hello World!"
}
```

## Mejores Prácticas

1. **Interfaces pequeñas**: Prefiere interfaces con pocos métodos, idealmente uno solo. Esto facilita la implementación y promueve la composición.

2. **Define interfaces donde las uses**: Define las interfaces en el paquete que las usa, no en el paquete que las implementa.

3. **Acepta interfaces, devuelve estructuras**: Al diseñar funciones y métodos, acepta interfaces como parámetros para mayor flexibilidad, pero devuelve tipos concretos.

4. **Usa interfaces para desacoplar componentes**: Las interfaces permiten que diferentes partes de tu código se comuniquen sin conocerse directamente.

5. **Evita interfaces vacías cuando sea posible**: Aunque `interface{}` es flexible, pierde la seguridad de tipos en tiempo de compilación.

## Ejercicios Prácticos

1. Implementa el patrón Strategy usando interfaces para diferentes algoritmos de ordenamiento.

2. Crea un sistema de plugins usando interfaces para permitir extensiones de funcionalidad.

3. Diseña un middleware HTTP usando composición de interfaces.

4. Implementa un sistema de caché con diferentes backends (memoria, disco, redis) usando una interfaz común.

## Conclusión

Las interfaces en Go proporcionan un mecanismo poderoso para crear código flexible y modular. Al dominar los conceptos avanzados y patrones de diseño basados en interfaces, podrás escribir código más mantenible, testeable y adaptable a cambios futuros.

Recuerda que el verdadero poder de las interfaces en Go viene de su simplicidad y de cómo facilitan la composición sobre la herencia. Aprovecha estas características para crear sistemas bien diseñados y desacoplados.