# Patrones de Diseño en Go

## Introducción

Los patrones de diseño son soluciones probadas para problemas comunes en el diseño de software. Representan las mejores prácticas utilizadas por desarrolladores experimentados. En Go, algunos patrones de diseño tradicionales se implementan de manera diferente debido a las características del lenguaje, como la ausencia de herencia de clases y la presencia de interfaces implícitas, goroutines y canales.

En esta sección, exploraremos cómo se implementan varios patrones de diseño en Go, adaptados a su filosofía y características únicas.

## Patrones Creacionales

Los patrones creacionales se centran en mecanismos de creación de objetos, tratando de crear objetos de manera adecuada para cada situación.

### Factory Method

El patrón Factory Method define una interfaz para crear un objeto, pero deja que las subclases decidan qué clase instanciar.

```go
package main

import "fmt"

// Producto define la interfaz para los objetos creados por el factory
type Producto interface {
	Operacion() string
}

// ProductoConcreto1 implementa la interfaz Producto
type ProductoConcreto1 struct{}

func (p *ProductoConcreto1) Operacion() string {
	return "Resultado de la operación del ProductoConcreto1"
}

// ProductoConcreto2 implementa la interfaz Producto
type ProductoConcreto2 struct{}

func (p *ProductoConcreto2) Operacion() string {
	return "Resultado de la operación del ProductoConcreto2"
}

// CreadorProducto define la interfaz factory
type CreadorProducto interface {
	CrearProducto() Producto
}

// CreadorConcreto1 implementa CreadorProducto para crear ProductoConcreto1
type CreadorConcreto1 struct{}

func (c *CreadorConcreto1) CrearProducto() Producto {
	return &ProductoConcreto1{}
}

// CreadorConcreto2 implementa CreadorProducto para crear ProductoConcreto2
type CreadorConcreto2 struct{}

func (c *CreadorConcreto2) CrearProducto() Producto {
	return &ProductoConcreto2{}
}

func main() {
	creadores := []CreadorProducto{
		&CreadorConcreto1{},
		&CreadorConcreto2{},
	}

	for _, creador := range creadores {
		producto := creador.CrearProducto()
		fmt.Println(producto.Operacion())
	}
}
```

### Abstract Factory

El patrón Abstract Factory proporciona una interfaz para crear familias de objetos relacionados sin especificar sus clases concretas.

```go
package main

import "fmt"

// ProductoA define una interfaz para un tipo de producto
type ProductoA interface {
	OperacionA() string
}

// ProductoB define una interfaz para otro tipo de producto
type ProductoB interface {
	OperacionB() string
	OtraOperacionB(a ProductoA) string
}

// ProductoA1 implementa ProductoA
type ProductoA1 struct{}

func (p *ProductoA1) OperacionA() string {
	return "Producto A1"
}

// ProductoA2 implementa ProductoA
type ProductoA2 struct{}

func (p *ProductoA2) OperacionA() string {
	return "Producto A2"
}

// ProductoB1 implementa ProductoB
type ProductoB1 struct{}

func (p *ProductoB1) OperacionB() string {
	return "Producto B1"
}

func (p *ProductoB1) OtraOperacionB(a ProductoA) string {
	resultado := a.OperacionA()
	return fmt.Sprintf("Resultado de B1 colaborando con (%s)", resultado)
}

// ProductoB2 implementa ProductoB
type ProductoB2 struct{}

func (p *ProductoB2) OperacionB() string {
	return "Producto B2"
}

func (p *ProductoB2) OtraOperacionB(a ProductoA) string {
	resultado := a.OperacionA()
	return fmt.Sprintf("Resultado de B2 colaborando con (%s)", resultado)
}

// FabricaAbstracta define la interfaz para crear productos
type FabricaAbstracta interface {
	CrearProductoA() ProductoA
	CrearProductoB() ProductoB
}

// FabricaConcreta1 implementa FabricaAbstracta para crear ProductoA1 y ProductoB1
type FabricaConcreta1 struct{}

func (f *FabricaConcreta1) CrearProductoA() ProductoA {
	return &ProductoA1{}
}

func (f *FabricaConcreta1) CrearProductoB() ProductoB {
	return &ProductoB1{}
}

// FabricaConcreta2 implementa FabricaAbstracta para crear ProductoA2 y ProductoB2
type FabricaConcreta2 struct{}

func (f *FabricaConcreta2) CrearProductoA() ProductoA {
	return &ProductoA2{}
}

func (f *FabricaConcreta2) CrearProductoB() ProductoB {
	return &ProductoB2{}
}

// Cliente utiliza solo interfaces declaradas por FabricaAbstracta y productos
func Cliente(fabrica FabricaAbstracta) {
	productoA := fabrica.CrearProductoA()
	productoB := fabrica.CrearProductoB()

	fmt.Println(productoB.OperacionB())
	fmt.Println(productoB.OtraOperacionB(productoA))
}

func main() {
	fmt.Println("Cliente: Probando código cliente con el primer tipo de fábrica")
	Cliente(&FabricaConcreta1{})

	fmt.Println("\nCliente: Probando el mismo código cliente con el segundo tipo de fábrica")
	Cliente(&FabricaConcreta2{})
}
```

### Builder

El patrón Builder separa la construcción de un objeto complejo de su representación, permitiendo el mismo proceso de construcción para crear diferentes representaciones.

```go
package main

import "fmt"

// Producto representa el objeto complejo que se está construyendo
type Producto struct {
	ParteA string
	ParteB string
	ParteC string
}

// Builder define la interfaz para construir partes del producto
type Builder interface {
	ConstruirParteA()
	ConstruirParteB()
	ConstruirParteC()
	ObtenerProducto() *Producto
}

// BuilderConcreto implementa Builder y construye el Producto
type BuilderConcreto struct {
	producto *Producto
}

func NewBuilderConcreto() *BuilderConcreto {
	return &BuilderConcreto{
		producto: &Producto{},
	}
}

func (b *BuilderConcreto) ConstruirParteA() {
	b.producto.ParteA = "Parte A"
}

func (b *BuilderConcreto) ConstruirParteB() {
	b.producto.ParteB = "Parte B"
}

func (b *BuilderConcreto) ConstruirParteC() {
	b.producto.ParteC = "Parte C"
}

func (b *BuilderConcreto) ObtenerProducto() *Producto {
	producto := b.producto
	b.producto = &Producto{} // Reset para el siguiente uso
	return producto
}

// Director controla el algoritmo de construcción
type Director struct {
	builder Builder
}

func NewDirector(builder Builder) *Director {
	return &Director{
		builder: builder,
	}
}

func (d *Director) CambiarBuilder(builder Builder) {
	d.builder = builder
}

// Construir representa el algoritmo de construcción
func (d *Director) Construir() {
	d.builder.ConstruirParteA()
	d.builder.ConstruirParteB()
	d.builder.ConstruirParteC()
}

// Versión con método encadenado (fluent interface)
type FluentBuilder struct {
	producto *Producto
}

func NewFluentBuilder() *FluentBuilder {
	return &FluentBuilder{
		producto: &Producto{},
	}
}

func (b *FluentBuilder) ParteA(a string) *FluentBuilder {
	b.producto.ParteA = a
	return b
}

func (b *FluentBuilder) ParteB(b string) *FluentBuilder {
	b.producto.ParteB = b
	return b
}

func (b *FluentBuilder) ParteC(c string) *FluentBuilder {
	b.producto.ParteC = c
	return b
}

func (b *FluentBuilder) Build() *Producto {
	return b.producto
}

func main() {
	// Usando el patrón Builder tradicional
	builderConcreto := NewBuilderConcreto()
	director := NewDirector(builderConcreto)

	director.Construir()
	producto := builderConcreto.ObtenerProducto()

	fmt.Printf("Producto construido: %+v\n", producto)

	// Usando el builder con método encadenado
	productoFluent := NewFluentBuilder().
		ParteA("Parte A personalizada").
		ParteB("Parte B personalizada").
		ParteC("Parte C personalizada").
		Build()

	fmt.Printf("Producto fluent: %+v\n", productoFluent)
}
```

### Singleton

El patrón Singleton garantiza que una clase tenga solo una instancia y proporciona un punto de acceso global a ella. En Go, podemos implementarlo de varias maneras.

```go
package main

import (
	"fmt"
	"sync"
)

// Singleton representa la estructura que queremos como singleton
type Singleton struct {
	data string
}

var instance *Singleton
var once sync.Once

// GetInstance devuelve la única instancia del Singleton
func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{data: "Soy un singleton"}
		fmt.Println("Creando instancia singleton")
	})
	return instance
}

// Método para modificar los datos del singleton
func (s *Singleton) SetData(data string) {
	s.data = data
}

// Método para obtener los datos del singleton
func (s *Singleton) GetData() string {
	return s.data
}

func main() {
	// Obtener la instancia del singleton
	s1 := GetInstance()
	fmt.Println("Datos del singleton:", s1.GetData())

	// Modificar los datos
	s1.SetData("Datos modificados")

	// Obtener la instancia nuevamente (debe ser la misma)
	s2 := GetInstance()
	fmt.Println("Datos del singleton después de modificar:", s2.GetData())

	// Verificar que s1 y s2 son la misma instancia
	fmt.Println("¿Son la misma instancia?:", s1 == s2)

	// Demostrar concurrencia segura
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			s := GetInstance()
			fmt.Printf("Goroutine %d: %p\n", id, s)
		}(i)
	}
	wg.Wait()
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

// Objeto representa un recurso costoso de crear
type Objeto struct {
	ID int
}

// ObjectPool gestiona un conjunto de objetos reutilizables
type ObjectPool struct {
	mutex    sync.Mutex
	objetos  []*Objeto
	maxSize  int
	creados  int
	disponibles int
}

// NewObjectPool crea un nuevo pool con un tamaño máximo
func NewObjectPool(maxSize int) *ObjectPool {
	return &ObjectPool{
		objetos:  make([]*Objeto, 0, maxSize),
		maxSize:  maxSize,
		creados:  0,
		disponibles: 0,
	}
}

// Acquire obtiene un objeto del pool o crea uno nuevo si es necesario
func (p *ObjectPool) Acquire() *Objeto {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Si hay objetos disponibles, devolver uno
	if p.disponibles > 0 {
		p.disponibles--
		objeto := p.objetos[p.disponibles]
		p.objetos[p.disponibles] = nil // Evitar memory leak
		fmt.Printf("Adquirido objeto existente #%d del pool\n", objeto.ID)
		return objeto
	}

	// Si no hay objetos disponibles pero no hemos alcanzado el máximo, crear uno nuevo
	if p.creados < p.maxSize {
		p.creados++
		objeto := &Objeto{ID: p.creados}
		fmt.Printf("Creado nuevo objeto #%d\n", objeto.ID)
		return objeto
	}

	// Si llegamos aquí, el pool está lleno y todos los objetos están en uso
	fmt.Println("Pool lleno, esperando a que se libere un objeto...")
	return nil
}

// Release devuelve un objeto al pool
func (p *ObjectPool) Release(objeto *Objeto) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Verificar que el objeto no sea nil
	if objeto == nil {
		return
	}

	// Añadir el objeto de vuelta al pool
	if p.disponibles < p.maxSize {
		fmt.Printf("Devolviendo objeto #%d al pool\n", objeto.ID)
		p.objetos = append(p.objetos[:p.disponibles], objeto)
		p.disponibles++
	}
}

// GetStats devuelve estadísticas del pool
func (p *ObjectPool) GetStats() (creados, disponibles, enUso int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.creados, p.disponibles, p.creados - p.disponibles
}

func main() {
	// Crear un pool con máximo 3 objetos
	pool := NewObjectPool(3)

	// Simular uso concurrente del pool
	var wg sync.WaitGroup

	// Función para simular trabajo con un objeto
	trabajarConObjeto := func(id int) {
		defer wg.Done()

		// Adquirir un objeto
		objeto := pool.Acquire()
		if objeto == nil {
			fmt.Printf("Trabajador %d no pudo obtener un objeto\n", id)
			return
		}

		// Simular trabajo
		fmt.Printf("Trabajador %d usando objeto #%d\n", id, objeto.ID)
		time.Sleep(time.Millisecond * time.Duration(500+id*100))

		// Devolver el objeto al pool
		pool.Release(objeto)
		fmt.Printf("Trabajador %d terminó de usar objeto #%d\n", id, objeto.ID)
	}

	// Lanzar 5 trabajadores (más que el tamaño del pool)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go trabajarConObjeto(i)
	}

	// Esperar a que todos los trabajadores terminen
	wg.Wait()

	// Mostrar estadísticas finales
	creados, disponibles, enUso := pool.GetStats()
	fmt.Printf("\nEstadísticas finales - Creados: %d, Disponibles: %d, En uso: %d\n", 
		creados, disponibles, enUso)
}
```

## Patrones Estructurales

Los patrones estructurales se ocupan de cómo se componen las clases y los objetos para formar estructuras más grandes.

### Adapter

El patrón Adapter permite que interfaces incompatibles trabajen juntas, convirtiendo la interfaz de una clase en otra que el cliente espera.

```go
package main

import "fmt"

// ServicioObjetivo define la interfaz que el cliente utiliza
type ServicioObjetivo interface {
	Solicitud() string
}

// ServicioCliente es el cliente que usa ServicioObjetivo
type ServicioCliente struct{}

func (c *ServicioCliente) UsarServicio(servicio ServicioObjetivo) {
	fmt.Println(servicio.Solicitud())
}

// ServicioExistente es un servicio con una interfaz incompatible
type ServicioExistente struct{}

func (s *ServicioExistente) SolicitudEspecifica() string {
	return "Respuesta del servicio existente con interfaz incompatible"
}

// Adaptador implementa ServicioObjetivo y adapta ServicioExistente
type Adaptador struct {
	servicioExistente *ServicioExistente
}

func NewAdaptador(servicioExistente *ServicioExistente) *Adaptador {
	return &Adaptador{servicioExistente: servicioExistente}
}

func (a *Adaptador) Solicitud() string {
	// Adapta la llamada a la interfaz del servicio existente
	return fmt.Sprintf("Adaptador: %s", a.servicioExistente.SolicitudEspecifica())
}

func main() {
	cliente := &ServicioCliente{}
	servicioExistente := &ServicioExistente{}
	adaptador := NewAdaptador(servicioExistente)

	fmt.Println("Cliente: Puedo trabajar bien con el adaptador:")
	cliente.UsarServicio(adaptador)
}
```

### Bridge

El patrón Bridge separa una abstracción de su implementación para que ambas puedan variar independientemente.

```go
package main

import "fmt"

// Implementador define la interfaz para las clases de implementación
type Implementador interface {
	OperacionImpl() string
}

// ImplementadorConcreto1 es una implementación concreta
type ImplementadorConcreto1 struct{}

func (i *ImplementadorConcreto1) OperacionImpl() string {
	return "ImplementadorConcreto1: Resultado de la operación"
}

// ImplementadorConcreto2 es otra implementación concreta
type ImplementadorConcreto2 struct{}

func (i *ImplementadorConcreto2) OperacionImpl() string {
	return "ImplementadorConcreto2: Resultado de la operación"
}

// Abstraccion define la interfaz para la parte de control
type Abstraccion interface {
	Operacion() string
}

// AbstraccionBasica es una implementación básica de Abstraccion
type AbstraccionBasica struct {
	implementador Implementador
}

func NewAbstraccionBasica(implementador Implementador) *AbstraccionBasica {
	return &AbstraccionBasica{implementador: implementador}
}

func (a *AbstraccionBasica) Operacion() string {
	return fmt.Sprintf("Abstracción básica: %s", a.implementador.OperacionImpl())
}

// AbstraccionExtendida es una variante de la abstracción
type AbstraccionExtendida struct {
	implementador Implementador
}

func NewAbstraccionExtendida(implementador Implementador) *AbstraccionExtendida {
	return &AbstraccionExtendida{implementador: implementador}
}

func (a *AbstraccionExtendida) Operacion() string {
	return fmt.Sprintf("Abstracción extendida: %s", a.implementador.OperacionImpl())
}

func main() {
	implementador1 := &ImplementadorConcreto1{}
	implementador2 := &ImplementadorConcreto2{}

	abstraccion1 := NewAbstraccionBasica(implementador1)
	abstraccion2 := NewAbstraccionExtendida(implementador2)

	fmt.Println(abstraccion1.Operacion())
	fmt.Println(abstraccion2.Operacion())

	// Cambiar la implementación en tiempo de ejecución
	abstraccion1 = NewAbstraccionBasica(implementador2)
	fmt.Println(abstraccion1.Operacion())
}
```

### Composite

El patrón Composite compone objetos en estructuras de árbol para representar jerarquías parte-todo, permitiendo a los clientes tratar objetos individuales y composiciones de objetos de manera uniforme.

```go
package main

import "fmt"

// Componente define la interfaz común para objetos simples y compuestos
type Componente interface {
	Operacion() string
	Agregar(componente Componente)
	Eliminar(componente Componente)
	ObtenerHijo(indice int) Componente
}

// Hoja representa objetos finales de la composición que no tienen hijos
type Hoja struct {
	nombre string
}

func NewHoja(nombre string) *Hoja {
	return &Hoja{nombre: nombre}
}

func (h *Hoja) Operacion() string {
	return fmt.Sprintf("Hoja %s", h.nombre)
}

// Estos métodos no hacen nada en una hoja, pero deben implementarse para la interfaz
func (h *Hoja) Agregar(c Componente) {}
func (h *Hoja) Eliminar(c Componente) {}
func (h *Hoja) ObtenerHijo(indice int) Componente { return nil }

// Compuesto representa objetos complejos que pueden tener hijos
type Compuesto struct {
	nombre string
	hijos []Componente
}

func NewCompuesto(nombre string) *Compuesto {
	return &Compuesto{
		nombre: nombre,
		hijos: []Componente{},
	}
}

func (c *Compuesto) Operacion() string {
	resultado := fmt.Sprintf("Compuesto %s (", c.nombre)

	for i, hijo := range c.hijos {
		resultado += hijo.Operacion()
		if i < len(c.hijos)-1 {
			resultado += " + "
		}
	}

	return resultado + ")"
}

func (c *Compuesto) Agregar(componente Componente) {
	c.hijos = append(c.hijos, componente)
}

func (c *Compuesto) Eliminar(componente Componente) {
	for i, hijo := range c.hijos {
		if hijo == componente {
			c.hijos = append(c.hijos[:i], c.hijos[i+1:]...)
			break
		}
	}
}

func (c *Compuesto) ObtenerHijo(indice int) Componente {
	if indice < 0 || indice >= len(c.hijos) {
		return nil
	}
	return c.hijos[indice]
}

func main() {
	// Crear hojas
	hoja1 := NewHoja("A")
	hoja2 := NewHoja("B")
	hoja3 := NewHoja("C")

	// Crear compuestos
	compuesto1 := NewCompuesto("X")
	compuesto1.Agregar(hoja1)
	compuesto1.Agregar(hoja2)

	compuesto2 := NewCompuesto("Y")
	compuesto2.Agregar(hoja3)

	// Crear compuesto raíz
	compuestoRaiz := NewCompuesto("RAÍZ")
	compuestoRaiz.Agregar(compuesto1)
	compuestoRaiz.Agregar(compuesto2)

	// Mostrar la estructura
	fmt.Println(compuestoRaiz.Operacion())
}

### Decorator

El patrón Decorator añade responsabilidades adicionales a un objeto dinámicamente, proporcionando una alternativa flexible a la herencia para extender la funcionalidad.

```go
package main

import "fmt"

// Componente define la interfaz para objetos que pueden tener responsabilidades añadidas
type Componente interface {
	Operacion() string
}

// ComponenteConcreto implementa la interfaz Componente
type ComponenteConcreto struct{}

func (c *ComponenteConcreto) Operacion() string {
	return "Componente Concreto"
}

// Decorador implementa la interfaz Componente y mantiene una referencia a un objeto Componente
type Decorador struct {
	componente Componente
}

func (d *Decorador) Operacion() string {
	return d.componente.Operacion()
}

// DecoradorConcretoA extiende Decorador
type DecoradorConcretoA struct {
	Decorador
}

func NewDecoradorConcretoA(componente Componente) *DecoradorConcretoA {
	return &DecoradorConcretoA{Decorador{componente: componente}}
}

func (d *DecoradorConcretoA) Operacion() string {
	return fmt.Sprintf("DecoradorConcretoA(%s)", d.Decorador.Operacion())
}

// DecoradorConcretoB extiende Decorador
type DecoradorConcretoB struct {
	Decorador
}

func NewDecoradorConcretoB(componente Componente) *DecoradorConcretoB {
	return &DecoradorConcretoB{Decorador{componente: componente}}
}

func (d *DecoradorConcretoB) Operacion() string {
	return fmt.Sprintf("DecoradorConcretoB(%s)", d.Decorador.Operacion())
}

func main() {
	// Crear un componente simple
	componente := &ComponenteConcreto{}
	fmt.Println("Componente simple:")
	fmt.Println(componente.Operacion())

	// Decorar el componente con DecoradorConcretoA
	decoradorA := NewDecoradorConcretoA(componente)
	fmt.Println("\nComponente decorado con A:")
	fmt.Println(decoradorA.Operacion())

	// Decorar el componente con DecoradorConcretoB
	decoradorB := NewDecoradorConcretoB(componente)
	fmt.Println("\nComponente decorado con B:")
	fmt.Println(decoradorB.Operacion())

	// Decorar el componente con ambos decoradores
	decoradorBA := NewDecoradorConcretoB(NewDecoradorConcretoA(componente))
	fmt.Println("\nComponente decorado con B y A:")
	fmt.Println(decoradorBA.Operacion())
}
```

### Facade

El patrón Facade proporciona una interfaz unificada a un conjunto de interfaces en un subsistema, definiendo una interfaz de nivel superior que hace que el subsistema sea más fácil de usar.

```go
package main

import "fmt"

// Subsistema1 representa una parte compleja del sistema
type Subsistema1 struct{}

func (s *Subsistema1) Operacion1() string {
	return "Subsistema1: Listo!"
}

func (s *Subsistema1) OperacionN() string {
	return "Subsistema1: Ejecutando..."
}

// Subsistema2 representa otra parte compleja del sistema
type Subsistema2 struct{}

func (s *Subsistema2) Operacion1() string {
	return "Subsistema2: Preparado!"
}

func (s *Subsistema2) OperacionZ() string {
	return "Subsistema2: Procesando..."
}

// Fachada proporciona una interfaz simple para la lógica compleja de los subsistemas
type Fachada struct {
	subsistema1 *Subsistema1
	subsistema2 *Subsistema2
}

func NewFachada() *Fachada {
	return &Fachada{
		subsistema1: &Subsistema1{},
		subsistema2: &Subsistema2{},
	}
}

// OperacionSimple envuelve las operaciones complejas de los subsistemas
func (f *Fachada) OperacionSimple() string {
	resultado := "Fachada inicializa subsistemas:\n"
	resultado += f.subsistema1.Operacion1() + "\n"
	resultado += f.subsistema2.Operacion1() + "\n"
	resultado += "Fachada ordena a los subsistemas realizar operaciones:\n"
	resultado += f.subsistema1.OperacionN() + "\n"
	resultado += f.subsistema2.OperacionZ() + "\n"
	return resultado
}

func main() {
	fachada := NewFachada()

	// El cliente interactúa con el sistema a través de la fachada
	fmt.Println(fachada.OperacionSimple())
}
```

## Patrones de Comportamiento

Los patrones de comportamiento se ocupan de la comunicación entre objetos.

### Strategy

El patrón Strategy define una familia de algoritmos, encapsula cada uno de ellos y los hace intercambiables, permitiendo que el algoritmo varíe independientemente de los clientes que lo utilizan.

```go
package main

import "fmt"

// Estrategia define la interfaz común para todos los algoritmos soportados
type Estrategia interface {
	Ejecutar(a, b int) int
}

// EstrategiaSuma implementa la operación de suma
type EstrategiaSuma struct{}

func (s *EstrategiaSuma) Ejecutar(a, b int) int {
	return a + b
}

// EstrategiaResta implementa la operación de resta
type EstrategiaResta struct{}

func (s *EstrategiaResta) Ejecutar(a, b int) int {
	return a - b
}

// EstrategiaMultiplicacion implementa la operación de multiplicación
type EstrategiaMultiplicacion struct{}

func (s *EstrategiaMultiplicacion) Ejecutar(a, b int) int {
	return a * b
}

// Contexto utiliza una estrategia para ejecutar una operación
type Contexto struct {
	estrategia Estrategia
}

func NewContexto(estrategia Estrategia) *Contexto {
	return &Contexto{estrategia: estrategia}
}

func (c *Contexto) EstablecerEstrategia(estrategia Estrategia) {
	c.estrategia = estrategia
}

func (c *Contexto) EjecutarEstrategia(a, b int) int {
	return c.estrategia.Ejecutar(a, b)
}

func main() {
	// Crear el contexto con una estrategia inicial
	contexto := NewContexto(&EstrategiaSuma{})

	// Ejecutar la estrategia
	resultado := contexto.EjecutarEstrategia(10, 5)
	fmt.Printf("10 + 5 = %d\n", resultado)

	// Cambiar la estrategia y ejecutar de nuevo
	contexto.EstablecerEstrategia(&EstrategiaResta{})
	resultado = contexto.EjecutarEstrategia(10, 5)
	fmt.Printf("10 - 5 = %d\n", resultado)

	// Cambiar a otra estrategia
	contexto.EstablecerEstrategia(&EstrategiaMultiplicacion{})
	resultado = contexto.EjecutarEstrategia(10, 5)
	fmt.Printf("10 * 5 = %d\n", resultado)
}
```

### Observer

El patrón Observer define una dependencia uno-a-muchos entre objetos, de modo que cuando un objeto cambia de estado, todos sus dependientes son notificados y actualizados automáticamente.

```go
package main

import "fmt"

// Observador define la interfaz para todos los observadores
type Observador interface {
	Actualizar(mensaje string)
}

// Sujeto mantiene una lista de observadores y los notifica de cambios
type Sujeto struct {
	observadores []Observador
	estado       string
}

func NewSujeto() *Sujeto {
	return &Sujeto{
		observadores: make([]Observador, 0),
	}
}

func (s *Sujeto) Adjuntar(observador Observador) {
	s.observadores = append(s.observadores, observador)
}

func (s *Sujeto) Separar(observador Observador) {
	for i, obs := range s.observadores {
		if obs == observador {
			s.observadores = append(s.observadores[:i], s.observadores[i+1:]...)
			break
		}
	}
}

func (s *Sujeto) Notificar() {
	for _, observador := range s.observadores {
		observador.Actualizar(s.estado)
	}
}

func (s *Sujeto) EstablecerEstado(estado string) {
	s.estado = estado
	s.Notificar()
}

// ObservadorConcreto implementa la interfaz Observador
type ObservadorConcreto struct {
	nombre string
}

func NewObservadorConcreto(nombre string) *ObservadorConcreto {
	return &ObservadorConcreto{nombre: nombre}
}

func (o *ObservadorConcreto) Actualizar(mensaje string) {
	fmt.Printf("Observador %s recibió mensaje: %s\n", o.nombre, mensaje)
}

func main() {
	// Crear el sujeto
	sujeto := NewSujeto()

	// Crear observadores
	observador1 := NewObservadorConcreto("A")
	observador2 := NewObservadorConcreto("B")
	observador3 := NewObservadorConcreto("C")

	// Adjuntar observadores al sujeto
	sujeto.Adjuntar(observador1)
	sujeto.Adjuntar(observador2)
	sujeto.Adjuntar(observador3)

	// Cambiar el estado del sujeto
	fmt.Println("Cambiando estado a 'Primer cambio'")
	sujeto.EstablecerEstado("Primer cambio")

	// Separar un observador
	sujeto.Separar(observador2)

	// Cambiar el estado de nuevo
	fmt.Println("\nCambiando estado a 'Segundo cambio'")
	sujeto.EstablecerEstado("Segundo cambio")
}
```

### Command

El patrón Command encapsula una solicitud como un objeto, permitiendo parametrizar clientes con diferentes solicitudes, encolar o registrar solicitudes, y soportar operaciones que pueden deshacerse.

```go
package main

import "fmt"

// Comando define la interfaz para ejecutar una operación
type Comando interface {
	Ejecutar()
	Deshacer()
}

// Receptor sabe cómo realizar las operaciones asociadas con los comandos
type Receptor struct {
	estado string
}

func (r *Receptor) Accion(estado string) {
	r.estado = estado
	fmt.Printf("Receptor: Mi estado ahora es %s\n", r.estado)
}

func (r *Receptor) ObtenerEstado() string {
	return r.estado
}

// ComandoConcreto implementa Comando y define la relación entre una acción y un receptor
type ComandoConcreto struct {
	receptor     *Receptor
	estadoNuevo  string
	estadoPrevio string
}

func NewComandoConcreto(receptor *Receptor, estado string) *ComandoConcreto {
	return &ComandoConcreto{
		receptor:    receptor,
		estadoNuevo: estado,
	}
}

func (c *ComandoConcreto) Ejecutar() {
	c.estadoPrevio = c.receptor.ObtenerEstado()
	c.receptor.Accion(c.estadoNuevo)
}

func (c *ComandoConcreto) Deshacer() {
	c.receptor.Accion(c.estadoPrevio)
}

// Invocador pide al comando que ejecute la solicitud
type Invocador struct {
	comandos []Comando
	indice   int
}

func NewInvocador() *Invocador {
	return &Invocador{
		comandos: make([]Comando, 0),
		indice:   -1,
	}
}

func (i *Invocador) AlmacenarYEjecutar(comando Comando) {
	// Si hemos deshecho comandos, eliminar los comandos después del índice actual
	if i.indice < len(i.comandos)-1 {
		i.comandos = i.comandos[:i.indice+1]
	}

	// Añadir y ejecutar el nuevo comando
	i.comandos = append(i.comandos, comando)
	i.indice = len(i.comandos) - 1
	comando.Ejecutar()
}

func (i *Invocador) Deshacer() {
	if i.indice >= 0 {
		i.comandos[i.indice].Deshacer()
		i.indice--
	}
}

func (i *Invocador) Rehacer() {
	if i.indice < len(i.comandos)-1 {
		i.indice++
		i.comandos[i.indice].Ejecutar()
	}
}

func main() {
	// Crear el receptor
	receptor := &Receptor{estado: "Inicial"}

	// Crear el invocador
	invocador := NewInvocador()

	// Crear y ejecutar comandos
	comando1 := NewComandoConcreto(receptor, "Estado 1")
	invocador.AlmacenarYEjecutar(comando1)

	comando2 := NewComandoConcreto(receptor, "Estado 2")
	invocador.AlmacenarYEjecutar(comando2)

	comando3 := NewComandoConcreto(receptor, "Estado 3")
	invocador.AlmacenarYEjecutar(comando3)

	// Deshacer comandos
	fmt.Println("\nDeshaciendo...")
	invocador.Deshacer()
	fmt.Printf("Estado actual: %s\n", receptor.ObtenerEstado())

	invocador.Deshacer()
	fmt.Printf("Estado actual: %s\n", receptor.ObtenerEstado())

	// Rehacer comandos
	fmt.Println("\nRehaciendo...")
	invocador.Rehacer()
	fmt.Printf("Estado actual: %s\n", receptor.ObtenerEstado())

	invocador.Rehacer()
	fmt.Printf("Estado actual: %s\n", receptor.ObtenerEstado())
}
```

## Patrones Específicos de Go

Además de los patrones de diseño tradicionales, Go tiene algunos patrones específicos que aprovechan sus características únicas.

### Functional Options

El patrón Functional Options permite configurar objetos con múltiples opciones de manera flexible y legible.

```go
package main

import "fmt"

// Servidor representa un servidor HTTP
type Servidor struct {
	direccion string
	puerto    int
	timeout   int
	maxConns  int
	tls       bool
}

// OpcionServidor define una función que configura un Servidor
type OpcionServidor func(*Servidor)

// PorDefecto establece valores predeterminados para el servidor
func PorDefecto() *Servidor {
	return &Servidor{
		direccion: "0.0.0.0",
		puerto:    8080,
		timeout:   30,
		maxConns:  1000,
		tls:       false,
	}
}

// NuevoServidor crea un servidor con opciones personalizadas
func NuevoServidor(opciones ...OpcionServidor) *Servidor {
	servidor := PorDefecto()

	// Aplicar cada opción al servidor
	for _, opcion := range opciones {
		opcion(servidor)
	}

	return servidor
}

// ConDireccion establece la dirección del servidor
func ConDireccion(direccion string) OpcionServidor {
	return func(s *Servidor) {
		s.direccion = direccion
	}
}

// ConPuerto establece el puerto del servidor
func ConPuerto(puerto int) OpcionServidor {
	return func(s *Servidor) {
		s.puerto = puerto
	}
}

// ConTimeout establece el timeout del servidor
func ConTimeout(timeout int) OpcionServidor {
	return func(s *Servidor) {
		s.timeout = timeout
	}
}

// ConMaxConexiones establece el número máximo de conexiones
func ConMaxConexiones(maxConns int) OpcionServidor {
	return func(s *Servidor) {
		s.maxConns = maxConns
	}
}

// ConTLS habilita o deshabilita TLS
func ConTLS(tls bool) OpcionServidor {
	return func(s *Servidor) {
		s.tls = tls
	}
}

func (s *Servidor) String() string {
	return fmt.Sprintf("Servidor{direccion: %s, puerto: %d, timeout: %d, maxConns: %d, tls: %v}",
		s.direccion, s.puerto, s.timeout, s.maxConns, s.tls)
}

func main() {
	// Crear un servidor con valores predeterminados
	servidorPorDefecto := NuevoServidor()
	fmt.Println("Servidor por defecto:")
	fmt.Println(servidorPorDefecto)

	// Crear un servidor con algunas opciones personalizadas
	servidorPersonalizado := NuevoServidor(
		ConDireccion("127.0.0.1"),
		ConPuerto(443),
		ConTLS(true),
	)
	fmt.Println("\nServidor personalizado:")
	fmt.Println(servidorPersonalizado)

	// Crear un servidor con todas las opciones personalizadas
	servidorCompleto := NuevoServidor(
		ConDireccion("10.0.0.1"),
		ConPuerto(8443),
		ConTimeout(60),
		ConMaxConexiones(5000),
		ConTLS(true),
	)
	fmt.Println("\nServidor completo:")
	fmt.Println(servidorCompleto)
}
```

### Context

El patrón Context en Go se utiliza para propagar cancelación, plazos y valores a través de la cadena de llamadas de API.

```go
package main

import (
	"context"
	"fmt"
	"time"
)

// Función que simula una operación larga que puede ser cancelada
func operacionLarga(ctx context.Context) (string, error) {
	// Crear un canal para simular trabajo
	completado := make(chan string)

	// Iniciar el trabajo en una goroutine
	go func() {
		// Simular trabajo que toma tiempo
		tiempo := 2 * time.Second
		fmt.Printf("Iniciando operación larga (durará %v)...\n", tiempo)
		time.Sleep(tiempo)
		completado <- "Operación completada con éxito"
	}()

	// Esperar a que el trabajo se complete o el contexto sea cancelado
	select {
	case resultado := <-completado:
		return resultado, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	// Ejemplo 1: Contexto con cancelación manual
	fmt.Println("Ejemplo 1: Contexto con cancelación manual")
	ctx, cancelar := context.WithCancel(context.Background())

	// Cancelar después de 1 segundo
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Cancelando contexto...")
		cancelar()
	}()

	resultado, err := operacionLarga(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Resultado: %s\n", resultado)
	}

	// Ejemplo 2: Contexto con timeout
	fmt.Println("\nEjemplo 2: Contexto con timeout")
	ctxTimeout, cancelarTimeout := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelarTimeout()

	resultado, err = operacionLarga(ctxTimeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Resultado: %s\n", resultado)
	}

	// Ejemplo 3: Contexto con valores
	fmt.Println("\nEjemplo 3: Contexto con valores")
	ctxValor := context.WithValue(context.Background(), "clave", "valor")

	// Función que utiliza el valor del contexto
	func(ctx context.Context) {
		valor, ok := ctx.Value("clave").(string)
		if !ok {
			fmt.Println("Clave no encontrada o tipo incorrecto")
			return
		}
		fmt.Printf("Valor obtenido del contexto: %s\n", valor)
	}(ctxValor)
}
```

### Pipeline

El patrón Pipeline en Go se utiliza para procesar streams de datos a través de una serie de etapas de procesamiento.

```go
package main

import (
	"fmt"
	"math"
	"sync"
)

// generador crea un canal que emite los números del 1 al n
func generador(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()
	return out
}

// cuadrado recibe números de un canal y emite sus cuadrados
func cuadrado(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// esPrimo verifica si un número es primo
func esPrimo(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	limite := int(math.Sqrt(float64(n)))
	for i := 5; i <= limite; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// filtrarPrimos filtra los números primos
func filtrarPrimos(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if esPrimo(n) {
				out <- n
			}
		}
	}()
	return out
}

// merge combina múltiples canales en uno solo
func merge(cs ...<-chan int) <-chan int {
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
	// Crear un pipeline simple
	fmt.Println("Pipeline simple: generador -> cuadrado -> filtrarPrimos")
	numeros := generador(20)
	cuadrados := cuadrado(numeros)
	primos := filtrarPrimos(cuadrados)

	// Consumir los resultados
	fmt.Println("Cuadrados que son números primos:")
	for primo := range primos {
		fmt.Printf("%d es un cuadrado y es primo\n", primo)
	}

	// Crear un pipeline con paralelismo
	fmt.Println("\nPipeline con paralelismo")
	numeros = generador(50)

	// Dividir el trabajo entre múltiples etapas de cuadrado
	c1 := cuadrado(numeros)
	c2 := cuadrado(numeros)
	c3 := cuadrado(numeros)

	// Combinar los resultados
	cuadradosMerged := merge(c1, c2, c3)
	primosMerged := filtrarPrimos(cuadradosMerged)

	// Consumir los resultados
	fmt.Println("Resultados del pipeline paralelo:")
	for primo := range primosMerged {
		fmt.Printf("%d\n", primo)
	}
}
```

## Conclusión

Los patrones de diseño son herramientas valiosas para resolver problemas comunes en el desarrollo de software. En Go, algunos patrones tradicionales se implementan de manera diferente debido a las características del lenguaje, mientras que otros patrones específicos de Go aprovechan sus características únicas como goroutines, canales e interfaces implícitas.

Al aplicar patrones de diseño en Go, es importante mantener la simplicidad y claridad que caracteriza al lenguaje, evitando la sobreingeniería y utilizando los patrones solo cuando realmente aportan valor a la solución.

## Recursos Adicionales

- [Go Design Patterns](https://github.com/tmrts/go-patterns)
- [Design Patterns in Go](https://refactoring.guru/design-patterns/go)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Concurrency Patterns](https://blog.golang.org/pipelines)
- [Context Package](https://golang.org/pkg/context/)
- [Functional Options Pattern](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

---

Con esto concluimos nuestra exploración de los patrones de diseño en Go. Estos patrones te ayudarán a escribir código más mantenible, flexible y reutilizable, aprovechando al máximo las características del lenguaje Go.