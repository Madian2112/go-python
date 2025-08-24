# Patrones de Diseño en Go

## Introducción

Los patrones de diseño son soluciones probadas a problemas comunes en el diseño de software. Representan las mejores prácticas utilizadas por desarrolladores experimentados para resolver problemas recurrentes. En este documento, exploraremos cómo implementar patrones de diseño en Go, adaptándolos a las características y filosofía del lenguaje.

Go es un lenguaje con un enfoque minimalista y pragmático, por lo que la implementación de patrones de diseño en Go suele ser más sencilla y directa que en lenguajes orientados a objetos tradicionales. Sin embargo, esto no significa que los patrones no sean útiles en Go; simplemente se implementan de manera diferente, aprovechando las características únicas del lenguaje como interfaces, composición y funciones de primera clase.

## Patrones Creacionales

Los patrones creacionales se centran en mecanismos de creación de objetos, tratando de crear objetos de manera adecuada para cada situación.

### Factory Method

El patrón Factory Method define una interfaz para crear un objeto, pero deja que las subclases decidan qué clase instanciar.

```go
package main

import "fmt"

// Product define la interfaz para los objetos creados por el factory
type Product interface {
    Use() string
}

// ConcreteProductA implementa Product
type ConcreteProductA struct{}

func (p *ConcreteProductA) Use() string {
    return "Usando ConcreteProductA"
}

// ConcreteProductB implementa Product
type ConcreteProductB struct{}

func (p *ConcreteProductB) Use() string {
    return "Usando ConcreteProductB"
}

// Creator define la interfaz factory method
type Creator interface {
    CreateProduct() Product
}

// ConcreteCreatorA implementa Creator para crear ConcreteProductA
type ConcreteCreatorA struct{}

func (c *ConcreteCreatorA) CreateProduct() Product {
    return &ConcreteProductA{}
}

// ConcreteCreatorB implementa Creator para crear ConcreteProductB
type ConcreteCreatorB struct{}

func (c *ConcreteCreatorB) CreateProduct() Product {
    return &ConcreteProductB{}
}

func main() {
    creatorA := &ConcreteCreatorA{}
    productA := creatorA.CreateProduct()
    fmt.Println(productA.Use())
    
    creatorB := &ConcreteCreatorB{}
    productB := creatorB.CreateProduct()
    fmt.Println(productB.Use())
}
```

### Abstract Factory

El patrón Abstract Factory proporciona una interfaz para crear familias de objetos relacionados sin especificar sus clases concretas.

```go
package main

import "fmt"

// ProductA define una interfaz para un tipo de producto
type ProductA interface {
    UseA() string
}

// ProductB define una interfaz para otro tipo de producto
type ProductB interface {
    UseB() string
}

// ConcreteProductA1 implementa ProductA
type ConcreteProductA1 struct{}

func (p *ConcreteProductA1) UseA() string {
    return "Usando ProductA1"
}

// ConcreteProductA2 implementa ProductA
type ConcreteProductA2 struct{}

func (p *ConcreteProductA2) UseA() string {
    return "Usando ProductA2"
}

// ConcreteProductB1 implementa ProductB
type ConcreteProductB1 struct{}

func (p *ConcreteProductB1) UseB() string {
    return "Usando ProductB1"
}

// ConcreteProductB2 implementa ProductB
type ConcreteProductB2 struct{}

func (p *ConcreteProductB2) UseB() string {
    return "Usando ProductB2"
}

// AbstractFactory define la interfaz para crear familias de productos
type AbstractFactory interface {
    CreateProductA() ProductA
    CreateProductB() ProductB
}

// ConcreteFactory1 implementa AbstractFactory para la familia de productos 1
type ConcreteFactory1 struct{}

func (f *ConcreteFactory1) CreateProductA() ProductA {
    return &ConcreteProductA1{}
}

func (f *ConcreteFactory1) CreateProductB() ProductB {
    return &ConcreteProductB1{}
}

// ConcreteFactory2 implementa AbstractFactory para la familia de productos 2
type ConcreteFactory2 struct{}

func (f *ConcreteFactory2) CreateProductA() ProductA {
    return &ConcreteProductA2{}
}

func (f *ConcreteFactory2) CreateProductB() ProductB {
    return &ConcreteProductB2{}
}

func main() {
    // Usar la primera familia de productos
    factory1 := &ConcreteFactory1{}
    productA1 := factory1.CreateProductA()
    productB1 := factory1.CreateProductB()
    fmt.Println(productA1.UseA())
    fmt.Println(productB1.UseB())
    
    // Usar la segunda familia de productos
    factory2 := &ConcreteFactory2{}
    productA2 := factory2.CreateProductA()
    productB2 := factory2.CreateProductB()
    fmt.Println(productA2.UseA())
    fmt.Println(productB2.UseB())
}
```

### Builder

El patrón Builder separa la construcción de un objeto complejo de su representación, permitiendo el mismo proceso de construcción para crear diferentes representaciones.

```go
package main

import "fmt"

// Product representa el objeto complejo que estamos construyendo
type Product struct {
    PartA string
    PartB string
    PartC string
}

// Builder define la interfaz para construir partes del producto
type Builder interface {
    BuildPartA()
    BuildPartB()
    BuildPartC()
    GetResult() *Product
}

// ConcreteBuilder implementa Builder
type ConcreteBuilder struct {
    product *Product
}

func NewConcreteBuilder() *ConcreteBuilder {
    return &ConcreteBuilder{product: &Product{}}
}

func (b *ConcreteBuilder) BuildPartA() {
    b.product.PartA = "Parte A"
}

func (b *ConcreteBuilder) BuildPartB() {
    b.product.PartB = "Parte B"
}

func (b *ConcreteBuilder) BuildPartC() {
    b.product.PartC = "Parte C"
}

func (b *ConcreteBuilder) GetResult() *Product {
    return b.product
}

// Director controla el proceso de construcción
type Director struct {
    builder Builder
}

func NewDirector(builder Builder) *Director {
    return &Director{builder: builder}
}

func (d *Director) Construct() {
    d.builder.BuildPartA()
    d.builder.BuildPartB()
    d.builder.BuildPartC()
}

func main() {
    builder := NewConcreteBuilder()
    director := NewDirector(builder)
    
    director.Construct()
    product := builder.GetResult()
    
    fmt.Printf("Producto construido: %+v\n", product)
}
```

### Singleton

El patrón Singleton garantiza que una clase tenga solo una instancia y proporciona un punto de acceso global a ella.

```go
package main

import (
    "fmt"
    "sync"
)

// Singleton representa la clase singleton
type Singleton struct {
    data string
}

var instance *Singleton
var once sync.Once

// GetInstance proporciona el punto de acceso global a la instancia singleton
func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{data: "Datos del singleton"}
        fmt.Println("Creando instancia singleton")
    })
    return instance
}

func (s *Singleton) GetData() string {
    return s.data
}

func main() {
    // Primera llamada, se crea la instancia
    singleton1 := GetInstance()
    fmt.Println(singleton1.GetData())
    
    // Segunda llamada, se reutiliza la instancia existente
    singleton2 := GetInstance()
    fmt.Println(singleton2.GetData())
    
    // Verificar que ambas variables apuntan a la misma instancia
    fmt.Printf("¿Misma instancia? %v\n", singleton1 == singleton2)
}
```

### Object Pool

El patrón Object Pool mantiene un conjunto de objetos inicializados listos para usar, en lugar de asignarlos y destruirlos bajo demanda.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Resource representa un recurso costoso de crear
type Resource struct {
    id int
}

func (r *Resource) Do(work int) {
    fmt.Printf("Recurso %d haciendo trabajo %d\n", r.id, work)
    time.Sleep(100 * time.Millisecond) // Simular trabajo
}

// Pool representa un pool de recursos
type Pool struct {
    resources []*Resource
    mutex     sync.Mutex
}

func NewPool(size int) *Pool {
    pool := &Pool{
        resources: make([]*Resource, size),
    }
    
    // Inicializar recursos
    for i := 0; i < size; i++ {
        pool.resources[i] = &Resource{id: i}
    }
    
    return pool
}

func (p *Pool) Acquire() *Resource {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    // Esperar hasta que haya un recurso disponible
    if len(p.resources) == 0 {
        fmt.Println("Esperando recurso disponible...")
        time.Sleep(100 * time.Millisecond)
        return p.Acquire()
    }
    
    // Obtener el último recurso
    resource := p.resources[len(p.resources)-1]
    p.resources = p.resources[:len(p.resources)-1]
    
    fmt.Printf("Adquirido recurso %d\n", resource.id)
    return resource
}

func (p *Pool) Release(r *Resource) {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    fmt.Printf("Liberando recurso %d\n", r.id)
    p.resources = append(p.resources, r)
}

func worker(id int, pool *Pool, wg *sync.WaitGroup) {
    defer wg.Done()
    
    // Adquirir recurso del pool
    resource := pool.Acquire()
    
    // Usar el recurso
    resource.Do(id)
    
    // Liberar el recurso de vuelta al pool
    pool.Release(resource)
}

func main() {
    // Crear un pool con 3 recursos
    pool := NewPool(3)
    
    // Crear 10 workers que compiten por los recursos
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go worker(i, pool, &wg)
    }
    
    wg.Wait()
    fmt.Println("Todos los trabajos completados")
}
```

## Patrones Estructurales

Los patrones estructurales se ocupan de cómo se componen las clases y los objetos para formar estructuras más grandes.

### Adapter

El patrón Adapter permite que interfaces incompatibles trabajen juntas, convirtiendo la interfaz de una clase en otra que el cliente espera.

```go
package main

import "fmt"

// Target define la interfaz que el cliente utiliza
type Target interface {
    Request() string
}

// Adaptee tiene una interfaz incompatible
type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string {
    return "Respuesta específica del Adaptee"
}

// Adapter adapta Adaptee a la interfaz Target
type Adapter struct {
    adaptee *Adaptee
}

func NewAdapter(adaptee *Adaptee) *Adapter {
    return &Adapter{adaptee: adaptee}
}

func (a *Adapter) Request() string {
    // Llamar al método del Adaptee y adaptarlo
    return fmt.Sprintf("Adapter: %s", a.adaptee.SpecificRequest())
}

// Cliente utiliza la interfaz Target
func Client(target Target) {
    fmt.Println(target.Request())
}

func main() {
    adaptee := &Adaptee{}
    adapter := NewAdapter(adaptee)
    
    Client(adapter)
}
```

### Bridge

El patrón Bridge separa una abstracción de su implementación para que ambas puedan variar independientemente.

```go
package main

import "fmt"

// Implementor define la interfaz para las clases de implementación
type Implementor interface {
    OperationImpl() string
}

// ConcreteImplementorA es una implementación concreta
type ConcreteImplementorA struct{}

func (c *ConcreteImplementorA) OperationImpl() string {
    return "Implementación A"
}

// ConcreteImplementorB es otra implementación concreta
type ConcreteImplementorB struct{}

func (c *ConcreteImplementorB) OperationImpl() string {
    return "Implementación B"
}

// Abstraction define la interfaz para la abstracción
type Abstraction interface {
    Operation() string
}

// RefinedAbstraction refina la interfaz definida por Abstraction
type RefinedAbstraction struct {
    implementor Implementor
}

func NewRefinedAbstraction(implementor Implementor) *RefinedAbstraction {
    return &RefinedAbstraction{implementor: implementor}
}

func (r *RefinedAbstraction) Operation() string {
    return fmt.Sprintf("Abstracción refinada con %s", r.implementor.OperationImpl())
}

func main() {
    implA := &ConcreteImplementorA{}
    implB := &ConcreteImplementorB{}
    
    abstractionA := NewRefinedAbstraction(implA)
    abstractionB := NewRefinedAbstraction(implB)
    
    fmt.Println(abstractionA.Operation())
    fmt.Println(abstractionB.Operation())
}
```

### Composite

El patrón Composite compone objetos en estructuras de árbol para representar jerarquías parte-todo, permitiendo a los clientes tratar objetos individuales y composiciones de objetos de manera uniforme.

```go
package main

import "fmt"

// Component define la interfaz para todos los objetos en la composición
type Component interface {
    Operation() string
    Add(Component)
    Remove(Component)
    GetChild(int) Component
}

// Leaf representa objetos hoja en la composición (sin hijos)
type Leaf struct {
    name string
}

func NewLeaf(name string) *Leaf {
    return &Leaf{name: name}
}

func (l *Leaf) Operation() string {
    return fmt.Sprintf("Leaf %s", l.name)
}

func (l *Leaf) Add(c Component) {
    // No hace nada, las hojas no pueden tener hijos
}

func (l *Leaf) Remove(c Component) {
    // No hace nada, las hojas no pueden tener hijos
}

func (l *Leaf) GetChild(i int) Component {
    return nil // Las hojas no tienen hijos
}

// Composite representa objetos compuestos que pueden tener hijos
type Composite struct {
    name     string
    children []Component
}

func NewComposite(name string) *Composite {
    return &Composite{
        name:     name,
        children: []Component{},
    }
}

func (c *Composite) Operation() string {
    result := fmt.Sprintf("Composite %s [\n", c.name)
    for _, child := range c.children {
        result += child.Operation() + "\n"
    }
    result += "]"
    return result
}

func (c *Composite) Add(component Component) {
    c.children = append(c.children, component)
}

func (c *Composite) Remove(component Component) {
    // Implementación simplificada, en la práctica necesitaríamos identificar el componente
    if len(c.children) > 0 {
        c.children = c.children[:len(c.children)-1]
    }
}

func (c *Composite) GetChild(i int) Component {
    if i >= 0 && i < len(c.children) {
        return c.children[i]
    }
    return nil
}

func main() {
    // Crear estructura de árbol
    root := NewComposite("root")
    branch1 := NewComposite("branch1")
    branch2 := NewComposite("branch2")
    leaf1 := NewLeaf("leaf1")
    leaf2 := NewLeaf("leaf2")
    leaf3 := NewLeaf("leaf3")
    
    root.Add(branch1)
    root.Add(branch2)
    branch1.Add(leaf1)
    branch1.Add(leaf2)
    branch2.Add(leaf3)
    
    // Mostrar la estructura completa
    fmt.Println(root.Operation())
}
```

### Decorator

El patrón Decorator añade responsabilidades adicionales a un objeto dinámicamente, proporcionando una alternativa flexible a la herencia para extender la funcionalidad.

```go
package main

import "fmt"

// Component define la interfaz para los objetos que pueden tener responsabilidades añadidas
type Component interface {
    Operation() string
}

// ConcreteComponent implementa la interfaz Component
type ConcreteComponent struct{}

func (c *ConcreteComponent) Operation() string {
    return "ConcreteComponent"
}

// Decorator mantiene una referencia al componente y define la interfaz conforme a Component
type Decorator struct {
    component Component
}

func (d *Decorator) Operation() string {
    if d.component != nil {
        return d.component.Operation()
    }
    return ""
}

// ConcreteDecoratorA añade responsabilidades al componente
type ConcreteDecoratorA struct {
    Decorator
}

func NewConcreteDecoratorA(c Component) *ConcreteDecoratorA {
    return &ConcreteDecoratorA{Decorator{component: c}}
}

func (d *ConcreteDecoratorA) Operation() string {
    return fmt.Sprintf("ConcreteDecoratorA(%s)", d.Decorator.Operation())
}

// ConcreteDecoratorB añade responsabilidades al componente
type ConcreteDecoratorB struct {
    Decorator
}

func NewConcreteDecoratorB(c Component) *ConcreteDecoratorB {
    return &ConcreteDecoratorB{Decorator{component: c}}
}

func (d *ConcreteDecoratorB) Operation() string {
    return fmt.Sprintf("ConcreteDecoratorB(%s)", d.Decorator.Operation())
}

func main() {
    // Crear un componente simple
    component := &ConcreteComponent{}
    fmt.Println("1. Componente simple:")
    fmt.Println(component.Operation())
    
    // Decorar el componente con DecoratorA
    decoratorA := NewConcreteDecoratorA(component)
    fmt.Println("\n2. Componente decorado con A:")
    fmt.Println(decoratorA.Operation())
    
    // Decorar el componente con DecoratorB
    decoratorB := NewConcreteDecoratorB(component)
    fmt.Println("\n3. Componente decorado con B:")
    fmt.Println(decoratorB.Operation())
    
    // Decorar el componente con ambos decoradores
    decoratorBA := NewConcreteDecoratorB(decoratorA)
    fmt.Println("\n4. Componente decorado con B y A:")
    fmt.Println(decoratorBA.Operation())
}
```

### Facade

El patrón Facade proporciona una interfaz unificada a un conjunto de interfaces en un subsistema, definiendo una interfaz de nivel superior que hace que el subsistema sea más fácil de usar.

```go
package main

import "fmt"

// Subsistema 1
type SubsystemA struct{}

func (s *SubsystemA) OperationA() string {
    return "Subsistema A: Listo!\n"
}

// Subsistema 2
type SubsystemB struct{}

func (s *SubsystemB) OperationB() string {
    return "Subsistema B: Listo!\n"
}

// Subsistema 3
type SubsystemC struct{}

func (s *SubsystemC) OperationC() string {
    return "Subsistema C: Listo!\n"
}

// Facade proporciona una interfaz unificada a un conjunto de interfaces en un subsistema
type Facade struct {
    subsystemA *SubsystemA
    subsystemB *SubsystemB
    subsystemC *SubsystemC
}

func NewFacade() *Facade {
    return &Facade{
        subsystemA: &SubsystemA{},
        subsystemB: &SubsystemB{},
        subsystemC: &SubsystemC{},
    }
}

func (f *Facade) Operation() string {
    result := "Facade inicializa subsistemas:\n"
    result += f.subsystemA.OperationA()
    result += f.subsystemB.OperationB()
    result += f.subsystemC.OperationC()
    return result
}

func main() {
    facade := NewFacade()
    fmt.Println(facade.Operation())
}
```

### Proxy

El patrón Proxy proporciona un sustituto o marcador de posición para otro objeto para controlar el acceso a él.

```go
package main

import "fmt"

// Subject define la interfaz común para RealSubject y Proxy
type Subject interface {
    Request() string
}

// RealSubject define el objeto real que el proxy representa
type RealSubject struct{}

func (s *RealSubject) Request() string {
    return "RealSubject: Manejando solicitud"
}

// Proxy mantiene una referencia al RealSubject
type Proxy struct {
    realSubject *RealSubject
}

func NewProxy() *Proxy {
    return &Proxy{}
}

func (p *Proxy) Request() string {
    // Inicialización lazy del RealSubject
    if p.realSubject == nil {
        fmt.Println("Proxy: Creando e inicializando RealSubject")
        p.realSubject = &RealSubject{}
    }
    
    fmt.Println("Proxy: Pre-procesando solicitud")
    result := p.realSubject.Request()
    fmt.Println("Proxy: Post-procesando solicitud")
    
    return result
}

func main() {
    // Usar el proxy para acceder al RealSubject
    proxy := NewProxy()
    fmt.Println(proxy.Request())
    
    // Segunda llamada, el RealSubject ya está creado
    fmt.Println("\nSegunda llamada:")
    fmt.Println(proxy.Request())
}
```

## Patrones de Comportamiento

Los patrones de comportamiento se ocupan de la comunicación entre objetos y cómo se distribuyen las responsabilidades.

### Chain of Responsibility

El patrón Chain of Responsibility evita acoplar el emisor de una petición a su receptor, dando a más de un objeto la oportunidad de manejar la petición.

```go
package main

import "fmt"

// Handler define la interfaz para manejar solicitudes
type Handler interface {
    SetNext(Handler) Handler
    Handle(request string) string
}

// BaseHandler proporciona funcionalidad común para los manejadores concretos
type BaseHandler struct {
    nextHandler Handler
}

func (h *BaseHandler) SetNext(handler Handler) Handler {
    h.nextHandler = handler
    return handler
}

func (h *BaseHandler) Handle(request string) string {
    if h.nextHandler != nil {
        return h.nextHandler.Handle(request)
    }
    return ""
}

// ConcreteHandlerA maneja solicitudes específicas
type ConcreteHandlerA struct {
    BaseHandler
}

func (h *ConcreteHandlerA) Handle(request string) string {
    if request == "A" {
        return fmt.Sprintf("ConcreteHandlerA: Manejando %s", request)
    }
    fmt.Printf("ConcreteHandlerA: Pasando %s al siguiente manejador\n", request)
    return h.BaseHandler.Handle(request)
}

// ConcreteHandlerB maneja solicitudes específicas
type ConcreteHandlerB struct {
    BaseHandler
}

func (h *ConcreteHandlerB) Handle(request string) string {
    if request == "B" {
        return fmt.Sprintf("ConcreteHandlerB: Manejando %s", request)
    }
    fmt.Printf("ConcreteHandlerB: Pasando %s al siguiente manejador\n", request)
    return h.BaseHandler.Handle(request)
}

// ConcreteHandlerC maneja solicitudes específicas
type ConcreteHandlerC struct {
    BaseHandler
}

func (h *ConcreteHandlerC) Handle(request string) string {
    if request == "C" {
        return fmt.Sprintf("ConcreteHandlerC: Manejando %s", request)
    }
    fmt.Printf("ConcreteHandlerC: Pasando %s al siguiente manejador\n", request)
    return h.BaseHandler.Handle(request)
}

func main() {
    // Configurar la cadena
    handlerA := &ConcreteHandlerA{}
    handlerB := &ConcreteHandlerB{}
    handlerC := &ConcreteHandlerC{}
    
    handlerA.SetNext(handlerB).SetNext(handlerC)
    
    // Probar la cadena con diferentes solicitudes
    fmt.Println(handlerA.Handle("A"))
    fmt.Println(handlerA.Handle("B"))
    fmt.Println(handlerA.Handle("C"))
    fmt.Println(handlerA.Handle("D"))
}
```

### Command

El patrón Command encapsula una solicitud como un objeto, permitiendo parametrizar clientes con diferentes solicitudes, encolar o registrar solicitudes, y soportar operaciones que se pueden deshacer.

```go
package main

import "fmt"

// Command define la interfaz para ejecutar una operación
type Command interface {
    Execute() string
    Undo() string
}

// Receiver sabe cómo realizar las operaciones asociadas con una solicitud
type Receiver struct {
    name string
}

func NewReceiver(name string) *Receiver {
    return &Receiver{name: name}
}

func (r *Receiver) Action(action string) string {
    return fmt.Sprintf("Receiver '%s': %s", r.name, action)
}

// ConcreteCommand implementa Command y define la vinculación entre un Receiver y una acción
type ConcreteCommand struct {
    receiver *Receiver
    action   string
}

func NewConcreteCommand(receiver *Receiver, action string) *ConcreteCommand {
    return &ConcreteCommand{
        receiver: receiver,
        action:   action,
    }
}

func (c *ConcreteCommand) Execute() string {
    return c.receiver.Action(c.action)
}

func (c *ConcreteCommand) Undo() string {
    return c.receiver.Action("Deshacer " + c.action)
}

// Invoker solicita al comando que ejecute la petición
type Invoker struct {
    commands []Command
    history  []Command
}

func NewInvoker() *Invoker {
    return &Invoker{
        commands: make([]Command, 0),
        history:  make([]Command, 0),
    }
}

func (i *Invoker) AddCommand(command Command) {
    i.commands = append(i.commands, command)
}

func (i *Invoker) ExecuteCommands() []string {
    results := make([]string, 0, len(i.commands))
    
    for _, command := range i.commands {
        result := command.Execute()
        i.history = append(i.history, command)
        results = append(results, result)
    }
    
    // Limpiar la lista de comandos pendientes
    i.commands = make([]Command, 0)
    
    return results
}

func (i *Invoker) UndoLastCommand() string {
    if len(i.history) == 0 {
        return "No hay comandos para deshacer"
    }
    
    lastIndex := len(i.history) - 1
    lastCommand := i.history[lastIndex]
    i.history = i.history[:lastIndex]
    
    return lastCommand.Undo()
}

func main() {
    // Configurar
    receiver1 := NewReceiver("Luz")
    receiver2 := NewReceiver("Música")
    
    turnOnLightCommand := NewConcreteCommand(receiver1, "Encender")
    turnOffLightCommand := NewConcreteCommand(receiver1, "Apagar")
    turnUpMusicCommand := NewConcreteCommand(receiver2, "Subir volumen")
    
    invoker := NewInvoker()
    
    // Agregar comandos y ejecutarlos
    invoker.AddCommand(turnOnLightCommand)
    invoker.AddCommand(turnUpMusicCommand)
    results := invoker.ExecuteCommands()
    
    fmt.Println("Resultados de ejecución:")
    for _, result := range results {
        fmt.Println(result)
    }
    
    // Deshacer el último comando
    fmt.Println("\nDeshacer:")
    fmt.Println(invoker.UndoLastCommand())
    
    // Ejecutar otro comando
    invoker.AddCommand(turnOffLightCommand)
    results = invoker.ExecuteCommands()
    
    fmt.Println("\nResultados de ejecución adicional:")
    for _, result := range results {
        fmt.Println(result)
    }
    
    // Deshacer de nuevo
    fmt.Println("\nDeshacer de nuevo:")
    fmt.Println(invoker.UndoLastCommand())
}
```

### Iterator

El patrón Iterator proporciona una forma de acceder a los elementos de un objeto agregado secuencialmente sin exponer su representación subyacente.

```go
package main

import "fmt"

// Iterator define la interfaz para acceder y recorrer elementos
type Iterator interface {
    HasNext() bool
    Next() interface{}
}

// Aggregate define la interfaz para crear un Iterator
type Aggregate interface {
    CreateIterator() Iterator
}

// ConcreteAggregate implementa Aggregate
type ConcreteAggregate struct {
    items []interface{}
}

func NewConcreteAggregate() *ConcreteAggregate {
    return &ConcreteAggregate{
        items: make([]interface{}, 0),
    }
}

func (a *ConcreteAggregate) Add(item interface{}) {
    a.items = append(a.items, item)
}

func (a *ConcreteAggregate) CreateIterator() Iterator {
    return &ConcreteIterator{
        aggregate: a,
        index:     0,
    }
}

// ConcreteIterator implementa Iterator
type ConcreteIterator struct {
    aggregate *ConcreteAggregate
    index     int
}

func (i *ConcreteIterator) HasNext() bool {
    return i.index < len(i.aggregate.items)
}

func (i *ConcreteIterator) Next() interface{} {
    if i.HasNext() {
        item := i.aggregate.items[i.index]
        i.index++
        return item
    }
    return nil
}

func main() {
    // Crear y llenar el agregado
    aggregate := NewConcreteAggregate()
    aggregate.Add("Item 1")
    aggregate.Add("Item 2")
    aggregate.Add("Item 3")
    aggregate.Add("Item 4")
    
    // Crear iterador y recorrer elementos
    iterator := aggregate.CreateIterator()
    
    fmt.Println("Iterando sobre elementos:")
    for iterator.HasNext() {
        item := iterator.Next()
        fmt.Println(item)
    }
}
```

### Observer

El patrón Observer define una dependencia uno-a-muchos entre objetos, de modo que cuando un objeto cambia de estado, todos sus dependientes son notificados y actualizados automáticamente.

```go
package main

import "fmt"

// Observer define la interfaz para los objetos que deben ser notificados de cambios
type Observer interface {
    Update(subject Subject)
}

// Subject define la interfaz para los objetos que pueden ser observados
type Subject interface {
    Attach(observer Observer)
    Detach(observer Observer)
    Notify()
    GetState() int
    SetState(state int)
}

// ConcreteSubject implementa Subject
type ConcreteSubject struct {
    observers []Observer
    state     int
}

func NewConcreteSubject() *ConcreteSubject {
    return &ConcreteSubject{
        observers: make([]Observer, 0),
    }
}

func (s *ConcreteSubject) Attach(observer Observer) {
    s.observers = append(s.observers, observer)
}

func (s *ConcreteSubject) Detach(observer Observer) {
    for i, obs := range s.observers {
        if obs == observer {
            s.observers = append(s.observers[:i], s.observers[i+1:]...)
            break
        }
    }
}

func (s *ConcreteSubject) Notify() {
    for _, observer := range s.observers {
        observer.Update(s)
    }
}

func (s *ConcreteSubject) GetState() int {
    return s.state
}

func (s *ConcreteSubject) SetState(state int) {
    s.state = state
    s.Notify()
}

// ConcreteObserverA implementa Observer
type ConcreteObserverA struct {
    name string
}

func NewConcreteObserverA(name string) *ConcreteObserverA {
    return &ConcreteObserverA{name: name}
}

func (o *ConcreteObserverA) Update(subject Subject) {
    fmt.Printf("ConcreteObserverA %s: Reaccionando al cambio de estado a %d\n", o.name, subject.GetState())
}

// ConcreteObserverB implementa Observer
type ConcreteObserverB struct {
    name string
}

func NewConcreteObserverB(name string) *ConcreteObserverB {
    return &ConcreteObserverB{name: name}
}

func (o *ConcreteObserverB) Update(subject Subject) {
    fmt.Printf("ConcreteObserverB %s: Reaccionando al cambio de estado a %d\n", o.name, subject.GetState())
}

func main() {
    // Crear sujeto y observadores
    subject := NewConcreteSubject()
    
    observerA1 := NewConcreteObserverA("A1")
    observerA2 := NewConcreteObserverA("A2")
    observerB := NewConcreteObserverB("B")
    
    // Registrar observadores
    subject.Attach(observerA1)
    subject.Attach(observerA2)
    subject.Attach(observerB)
    
    // Cambiar estado y notificar a los observadores
    fmt.Println("Cambiando estado a 10:")
    subject.SetState(10)
    
    // Eliminar un observador
    subject.Detach(observerA2)
    
    // Cambiar estado de nuevo
    fmt.Println("\nCambiando estado a 20:")
    subject.SetState(20)
}
```

### Strategy

El patrón Strategy define una familia de algoritmos, encapsula cada uno de ellos y los hace intercambiables. Strategy permite que el algoritmo varíe independientemente de los clientes que lo utilizan.

```go
package main

import "fmt"

// Strategy define la interfaz común a todos los algoritmos soportados
type Strategy interface {
    Execute(a, b int) int
}

// ConcreteStrategyAdd implementa la operación de suma
type ConcreteStrategyAdd struct{}

func (s *ConcreteStrategyAdd) Execute(a, b int) int {
    return a + b
}

// ConcreteStrategySubtract implementa la operación de resta
type ConcreteStrategySubtract struct{}

func (s *ConcreteStrategySubtract) Execute(a, b int) int {
    return a - b
}

// ConcreteStrategyMultiply implementa la operación de multiplicación
type ConcreteStrategyMultiply struct{}

func (s *ConcreteStrategyMultiply) Execute(a, b int) int {
    return a * b
}

// Context utiliza una Strategy para realizar una operación
type Context struct {
    strategy Strategy
}

func NewContext(strategy Strategy) *Context {
    return &Context{strategy: strategy}
}

func (c *Context) SetStrategy(strategy Strategy) {
    c.strategy = strategy
}

func (c *Context) ExecuteStrategy(a, b int) int {
    return c.strategy.Execute(a, b)
}

func main() {
    // Crear contexto con estrategia inicial
    context := NewContext(&ConcreteStrategyAdd{})
    
    // Ejecutar estrategia
    result := context.ExecuteStrategy(5, 3)
    fmt.Printf("5 + 3 = %d\n", result)
    
    // Cambiar estrategia y ejecutar
    context.SetStrategy(&ConcreteStrategySubtract{})
    result = context.ExecuteStrategy(5, 3)
    fmt.Printf("5 - 3 = %d\n", result)
    
    // Cambiar estrategia de nuevo y ejecutar
    context.SetStrategy(&ConcreteStrategyMultiply{})
    result = context.ExecuteStrategy(5, 3)
    fmt.Printf("5 * 3 = %d\n", result)
}
```

### Template Method

El patrón Template Method define el esqueleto de un algoritmo en una operación, aplazando algunos pasos a las subclases. Template Method permite que las subclases redefinan ciertos pasos de un algoritmo sin cambiar la estructura del algoritmo.

```go
package main

import "fmt"

// AbstractClass define la interfaz para el template method
type AbstractClass interface {
    TemplateMethod()
    PrimitiveOperation1() string
    PrimitiveOperation2() string
    Hook() bool
}

// AbstractClassImpl proporciona la implementación por defecto
type AbstractClassImpl struct {
    AbstractClass
}

func (a *AbstractClassImpl) TemplateMethod() {
    fmt.Println("Inicio del template method")
    fmt.Println(a.PrimitiveOperation1())
    fmt.Println(a.PrimitiveOperation2())
    
    if a.Hook() {
        fmt.Println("Hook activado")
    }
    
    fmt.Println("Fin del template method")
}

func (a *AbstractClassImpl) Hook() bool {
    return true // Implementación por defecto
}

// ConcreteClassA implementa las operaciones primitivas
type ConcreteClassA struct {
    AbstractClassImpl
}

func NewConcreteClassA() *ConcreteClassA {
    concrete := &ConcreteClassA{}
    concrete.AbstractClass = concrete // Inyectar la implementación concreta
    return concrete
}

func (c *ConcreteClassA) PrimitiveOperation1() string {
    return "ConcreteClassA: Implementación de la operación 1"
}

func (c *ConcreteClassA) PrimitiveOperation2() string {
    return "ConcreteClassA: Implementación de la operación 2"
}

// ConcreteClassB implementa las operaciones primitivas y sobrescribe el hook
type ConcreteClassB struct {
    AbstractClassImpl
}

func NewConcreteClassB() *ConcreteClassB {
    concrete := &ConcreteClassB{}
    concrete.AbstractClass = concrete // Inyectar la implementación concreta
    return concrete
}

func (c *ConcreteClassB) PrimitiveOperation1() string {
    return "ConcreteClassB: Implementación de la operación 1"
}

func (c *ConcreteClassB) PrimitiveOperation2() string {
    return "ConcreteClassB: Implementación de la operación 2"
}

func (c *ConcreteClassB) Hook() bool {
    return false // Sobrescribir el hook
}

func main() {
    fmt.Println("Ejecutando ConcreteClassA:")
    concreteA := NewConcreteClassA()
    concreteA.TemplateMethod()
    
    fmt.Println("\nEjecutando ConcreteClassB:")
    concreteB := NewConcreteClassB()
    concreteB.TemplateMethod()
}
```

## Patrones Funcionales en Go

Go es un lenguaje que soporta programación funcional, lo que permite implementar patrones funcionales que pueden ser más simples y expresivos que sus equivalentes orientados a objetos.

### Función de Orden Superior

Las funciones de orden superior toman otras funciones como argumentos o devuelven funciones como resultados.

```go
package main

import "fmt"

// Función que toma otra función como argumento
func applyOperation(a, b int, operation func(int, int) int) int {
    return operation(a, b)
}

// Función que devuelve otra función
func createMultiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

func main() {
    // Definir operaciones
    add := func(a, b int) int { return a + b }
    subtract := func(a, b int) int { return a - b }
    multiply := func(a, b int) int { return a * b }
    
    // Usar función de orden superior
    fmt.Printf("5 + 3 = %d\n", applyOperation(5, 3, add))
    fmt.Printf("5 - 3 = %d\n", applyOperation(5, 3, subtract))
    fmt.Printf("5 * 3 = %d\n", applyOperation(5, 3, multiply))
    
    // Usar función que devuelve otra función
    double := createMultiplier(2)
    triple := createMultiplier(3)
    
    fmt.Printf("El doble de 5 es %d\n", double(5))
    fmt.Printf("El triple de 5 es %d\n", triple(5))
}
```

### Decorador Funcional

El patrón decorador se puede implementar de manera más simple usando funciones de orden superior.

```go
package main

import (
    "fmt"
    "time"
)

// Función original
func greet(name string) string {
    return fmt.Sprintf("Hola, %s!", name)
}

// Decorador que mide el tiempo de ejecución
func withTiming(f func(string) string) func(string) string {
    return func(name string) string {
        start := time.Now()
        result := f(name)
        elapsed := time.Since(start)
        fmt.Printf("La función tomó %s en ejecutarse\n", elapsed)
        return result
    }
}

// Decorador que añade logging
func withLogging(f func(string) string) func(string) string {
    return func(name string) string {
        fmt.Printf("Llamando a función con argumento: %s\n", name)
        result := f(name)
        fmt.Printf("Función retornó: %s\n", result)
        return result
    }
}

func main() {
    // Decorar la función greet con timing
    timedGreet := withTiming(greet)
    fmt.Println(timedGreet("Juan"))
    
    // Decorar la función greet con logging y timing (composición de decoradores)
    loggedAndTimedGreet := withTiming(withLogging(greet))
    fmt.Println(loggedAndTimedGreet("María"))
}
```

### Pipeline

El patrón pipeline permite procesar datos a través de una serie de etapas, donde la salida de una etapa es la entrada de la siguiente.

```go
package main

import (
    "fmt"
    "strings"
)

// Etapas del pipeline
func generateNumbers(max int) []int {
    numbers := make([]int, max)
    for i := 0; i < max; i++ {
        numbers[i] = i + 1
    }
    return numbers
}

func filterEven(numbers []int) []int {
    var result []int
    for _, n := range numbers {
        if n%2 == 0 {
            result = append(result, n)
        }
    }
    return result
}

func square(numbers []int) []int {
    result := make([]int, len(numbers))
    for i, n := range numbers {
        result[i] = n * n
    }
    return result
}

func sum(numbers []int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

// Función de composición para crear un pipeline
func pipe(functions ...interface{}) interface{} {
    return func(input interface{}) interface{} {
        result := input
        for _, f := range functions {
            switch fn := f.(type) {
            case func([]int) []int:
                result = fn(result.([]int))
            case func([]int) int:
                result = fn(result.([]int))
            }
        }
        return result
    }
}

func main() {
    // Crear y ejecutar el pipeline
    pipeline := pipe(filterEven, square, sum)
    result := pipeline.(func(interface{}) interface{})(generateNumbers(10))
    
    fmt.Printf("La suma de los cuadrados de los números pares del 1 al 10 es: %d\n", result)
    
    // Otro ejemplo con strings
    process := func(texts []string) []string {
        // Pipeline de procesamiento de texto
        toUpper := func(texts []string) []string {
            result := make([]string, len(texts))
            for i, t := range texts {
                result[i] = strings.ToUpper(t)
            }
            return result
        }
        
        addPrefix := func(texts []string) []string {
            result := make([]string, len(texts))
            for i, t := range texts {
                result[i] = "Prefijo: " + t
            }
            return result
        }
        
        // Ejecutar el pipeline manualmente
        return addPrefix(toUpper(texts))
    }
    
    texts := []string{"hola", "mundo", "go"}
    processedTexts := process(texts)
    
    fmt.Println("Textos procesados:")
    for _, t := range processedTexts {
        fmt.Println(t)
    }
}
```

### Option Pattern

El patrón Option permite configurar objetos con opciones opcionales de manera flexible y legible.

```go
package main

import "fmt"

// Server representa un servidor HTTP
type Server struct {
    host    string
    port    int
    timeout int // en segundos
    maxConn int
    tls     bool
}

// ServerOption define una función que modifica un Server
type ServerOption func(*Server)

// WithPort establece el puerto del servidor
func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

// WithTimeout establece el timeout del servidor
func WithTimeout(timeout int) ServerOption {
    return func(s *Server) {
        s.timeout = timeout
    }
}

// WithMaxConn establece el número máximo de conexiones
func WithMaxConn(maxConn int) ServerOption {
    return func(s *Server) {
        s.maxConn = maxConn
    }
}

// WithTLS habilita o deshabilita TLS
func WithTLS(tls bool) ServerOption {
    return func(s *Server) {
        s.tls = tls
    }
}

// NewServer crea un nuevo servidor con opciones por defecto
func NewServer(host string, options ...ServerOption) *Server {
    // Valores por defecto
    server := &Server{
        host:    host,
        port:    8080,
        timeout: 30,
        maxConn: 100,
        tls:     false,
    }
    
    // Aplicar opciones
    for _, option := range options {
        option(server)
    }
    
    return server
}

func main() {
    // Crear servidor con valores por defecto
    server1 := NewServer("localhost")
    fmt.Printf("Servidor 1: %+v\n", server1)
    
    // Crear servidor con opciones personalizadas
    server2 := NewServer(
        "api.example.com",
        WithPort(443),
        WithTimeout(60),
        WithTLS(true),
    )
    fmt.Printf("Servidor 2: %+v\n", server2)
    
    // Crear servidor con solo algunas opciones
    server3 := NewServer(
        "admin.example.com",
        WithMaxConn(1000),
    )
    fmt.Printf("Servidor 3: %+v\n", server3)
}
```

## Patrones Concurrentes en Go

Go tiene soporte nativo para concurrencia a través de goroutines y canales, lo que permite implementar patrones concurrentes de manera elegante.

### Worker Pool

El patrón Worker Pool utiliza un número fijo de workers para procesar tareas de una cola.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Task representa una tarea a procesar
type Task struct {
    ID      int
    Data    string
    Process func(string) string
}

// Result representa el resultado de procesar una tarea
type Result struct {
    TaskID int
    Output string
}

// WorkerPool implementa un pool de workers
func WorkerPool(numWorkers int, tasks <-chan Task, results chan<- Result) {
    var wg sync.WaitGroup
    
    // Iniciar workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go worker(i, tasks, results, &wg)
    }
    
    // Esperar a que todos los workers terminen
    wg.Wait()
    close(results)
}

// worker procesa tareas de la cola
func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for task := range tasks {
        fmt.Printf("Worker %d procesando tarea %d\n", id, task.ID)
        output := task.Process(task.Data)
        
        // Simular trabajo
        time.Sleep(100 * time.Millisecond)
        
        results <- Result{
            TaskID: task.ID,
            Output: output,
        }
    }
}

func main() {
    // Crear canales para tareas y resultados
    tasks := make(chan Task, 10)
    results := make(chan Result, 10)
    
    // Iniciar worker pool con 3 workers
    go WorkerPool(3, tasks, results)
    
    // Enviar tareas
    numTasks := 10
    for i := 0; i < numTasks; i++ {
        tasks <- Task{
            ID:   i,
            Data: fmt.Sprintf("Datos de la tarea %d", i),
            Process: func(data string) string {
                return fmt.Sprintf("Procesado: %s", data)
            },
        }
    }
    close(tasks)
    
    // Recoger resultados
    for i := 0; i < numTasks; i++ {
        result := <-results
        fmt.Printf("Resultado de la tarea %d: %s\n", result.TaskID, result.Output)
    }
}
```

### Fan-Out, Fan-In

El patrón Fan-Out, Fan-In distribuye trabajo entre múltiples goroutines y luego combina los resultados.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// generator genera valores y los envía al canal de salida
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// square recibe valores, los eleva al cuadrado y los envía al canal de salida
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            fmt.Printf("Calculando cuadrado de %d\n", n)
            time.Sleep(100 * time.Millisecond) // Simular trabajo
            out <- n * n
        }
        close(out)
    }()
    return out
}

// fanOut distribuye el trabajo entre múltiples goroutines
func fanOut(in <-chan int, numWorkers int) []<-chan int {
    workers := make([]<-chan int, numWorkers)
    for i := 0; i < numWorkers; i++ {
        workers[i] = square(in)
    }
    return workers
}

// fanIn combina múltiples canales en uno solo
func fanIn(channels []<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    
    // Función para recopilar resultados de un canal
    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    
    wg.Add(len(channels))
    for _, c := range channels {
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
    // Generar números
    in := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    
    // Distribuir trabajo (fan-out)
    numWorkers := 3
    workers := fanOut(in, numWorkers)
    
    // Combinar resultados (fan-in)
    results := fanIn(workers)
    
    // Recoger resultados
    for result := range results {
        fmt.Println(result)
    }
}
```

### Future/Promise

El patrón Future/Promise permite ejecutar una tarea de forma asíncrona y obtener su resultado más tarde.

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

// Future representa una operación asíncrona
type Future struct {
    result chan interface{}
    err    chan error
}

// NewFuture crea un nuevo Future
func NewFuture() *Future {
    return &Future{
        result: make(chan interface{}, 1),
        err:    make(chan error, 1),
    }
}

// Get espera y devuelve el resultado del Future
func (f *Future) Get() (interface{}, error) {
    select {
    case result := <-f.result:
        return result, nil
    case err := <-f.err:
        return nil, err
    }
}

// Complete completa el Future con un resultado
func (f *Future) Complete(result interface{}) {
    f.result <- result
    close(f.result)
    close(f.err)
}

// CompleteWithError completa el Future con un error
func (f *Future) CompleteWithError(err error) {
    f.err <- err
    close(f.result)
    close(f.err)
}

// AsyncTask ejecuta una tarea de forma asíncrona y devuelve un Future
func AsyncTask(task func() (interface{}, error)) *Future {
    future := NewFuture()
    
    go func() {
        result, err := task()
        if err != nil {
            future.CompleteWithError(err)
        } else {
            future.Complete(result)
        }
    }()
    
    return future
}

func main() {
    // Definir una tarea que toma tiempo
    task := func() (interface{}, error) {
        // Simular trabajo
        fmt.Println("Iniciando tarea asíncrona...")
        time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
        
        // Generar un resultado aleatorio
        result := rand.Intn(100)
        fmt.Printf("Tarea completada con resultado: %d\n", result)
        
        return result, nil
    }
    
    // Ejecutar la tarea de forma asíncrona
    fmt.Println("Lanzando tarea...")
    future := AsyncTask(task)
    
    // Hacer otras cosas mientras la tarea se ejecuta
    fmt.Println("Haciendo otras cosas mientras la tarea se ejecuta...")
    time.Sleep(500 * time.Millisecond)
    
    // Obtener el resultado cuando esté listo
    fmt.Println("Esperando resultado...")
    result, err := future.Get()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Resultado obtenido: %d\n", result.(int))
    }
}
```

### Semaphore

El patrón Semaphore limita el número de goroutines que pueden acceder a un recurso o ejecutarse concurrentemente.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Semaphore implementa un semáforo usando un canal
type Semaphore chan struct{}

// NewSemaphore crea un nuevo semáforo con el número máximo de permisos
func NewSemaphore(maxConcurrency int) Semaphore {
    return make(Semaphore, maxConcurrency)
}

// Acquire adquiere un permiso del semáforo
func (s Semaphore) Acquire() {
    s <- struct{}{}
}

// Release libera un permiso al semáforo
func (s Semaphore) Release() {
    <-s
}

func main() {
    // Crear un semáforo que permite 3 operaciones concurrentes
    sem := NewSemaphore(3)
    var wg sync.WaitGroup
    
    // Ejecutar 10 tareas que compiten por el semáforo
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            fmt.Printf("Tarea %d: intentando adquirir el semáforo\n", id)
            sem.Acquire()
            defer sem.Release()
            
            fmt.Printf("Tarea %d: semáforo adquirido, ejecutando...\n", id)
            time.Sleep(1 * time.Second) // Simular trabajo
            fmt.Printf("Tarea %d: completada\n", id)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("Todas las tareas completadas")
}
```

## Ejercicios Prácticos

### Ejercicio 1: Implementar un Sistema de Plugins usando Factory Method

Crea un sistema de plugins donde diferentes tipos de plugins pueden ser registrados y creados dinámicamente.

```go
package main

import "fmt"

// Plugin define la interfaz para todos los plugins
type Plugin interface {
    Name() string
    Execute() string
}

// PluginFactory define la interfaz para crear plugins
type PluginFactory interface {
    CreatePlugin() Plugin
}

// Implementa diferentes tipos de plugins y sus factories
// Luego crea un registro de plugins que permita registrar y crear plugins por nombre
```

### Ejercicio 2: Implementar un Pipeline de Procesamiento de Datos

Crea un pipeline que procese una lista de números, aplicando diferentes transformaciones en cada etapa.

```go
package main

import "fmt"

// Implementa funciones para cada etapa del pipeline:
// 1. Generar números
// 2. Filtrar números pares
// 3. Elevar al cuadrado
// 4. Calcular la suma
// Conecta estas etapas usando canales
```

### Ejercicio 3: Implementar un Pool de Conexiones usando Object Pool

Crea un pool de conexiones a una base de datos simulada, limitando el número máximo de conexiones.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Connection simula una conexión a una base de datos
type Connection struct {
    id int
}

// Implementa un pool de conexiones que permita adquirir y liberar conexiones
// Asegúrate de que el pool tenga un límite máximo de conexiones
```

## Conclusiones

Los patrones de diseño son herramientas valiosas en el desarrollo de software, pero deben utilizarse con criterio. En Go, muchos patrones tradicionales pueden simplificarse gracias a las características del lenguaje como interfaces, composición y funciones de primera clase.

Al implementar patrones de diseño en Go, es importante seguir la filosofía del lenguaje:

1. **Simplicidad**: Preferir soluciones simples y directas.
2. **Pragmatismo**: Usar patrones solo cuando realmente resuelven un problema.
3. **Composición sobre herencia**: Aprovechar la composición de interfaces y structs.
4. **Concurrencia**: Utilizar goroutines y canales para patrones concurrentes.

Recuerda que el objetivo final de los patrones de diseño es crear código mantenible, flexible y comprensible. No uses patrones solo por usarlos; asegúrate de que realmente mejoran tu código.

## Referencias

1. Gamma, E., Helm, R., Johnson, R., & Vlissides, J. (1994). Design Patterns: Elements of Reusable Object-Oriented Software.
2. Pike, R. (2012). Go Concurrency Patterns.
3. Butcher, M., & Farina, M. (2017). Go in Practice: Includes 70 Techniques.
4. Gerrand, A. (2013). Go Concurrency Patterns: Pipelines and cancellation.
5. Donovan, A. A., & Kernighan, B. W. (2015). The Go Programming Language.
6. Bodner, J. (2019). Learning Go: An Idiomatic Approach to Real-World Go Programming.