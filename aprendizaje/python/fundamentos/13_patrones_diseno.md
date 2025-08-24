# Patrones de Diseño en Python

## Introducción

Los patrones de diseño son soluciones probadas para problemas comunes en el diseño de software. Representan las mejores prácticas utilizadas por desarrolladores experimentados. En Python, los patrones de diseño se implementan aprovechando las características del lenguaje como tipado dinámico, funciones de primera clase, y programación orientada a objetos.

En esta sección, exploraremos cómo se implementan varios patrones de diseño en Python, adaptados a su filosofía y características únicas.

## Patrones Creacionales

Los patrones creacionales se centran en mecanismos de creación de objetos, tratando de crear objetos de manera adecuada para cada situación.

### Factory Method

El patrón Factory Method define una interfaz para crear un objeto, pero deja que las subclases decidan qué clase instanciar.

```python
from abc import ABC, abstractmethod

# Producto define la interfaz para los objetos creados por el factory
class Producto(ABC):
    @abstractmethod
    def operacion(self) -> str:
        pass

# ProductoConcreto1 implementa la interfaz Producto
class ProductoConcreto1(Producto):
    def operacion(self) -> str:
        return "Resultado de la operación del ProductoConcreto1"

# ProductoConcreto2 implementa la interfaz Producto
class ProductoConcreto2(Producto):
    def operacion(self) -> str:
        return "Resultado de la operación del ProductoConcreto2"

# CreadorProducto define la interfaz factory
class CreadorProducto(ABC):
    @abstractmethod
    def crear_producto(self) -> Producto:
        pass
    
    def alguna_operacion(self) -> str:
        # Llamar al método factory para crear un objeto Producto
        producto = self.crear_producto()
        # Usar el producto
        return f"Creador: {producto.operacion()}"

# CreadorConcreto1 implementa CreadorProducto para crear ProductoConcreto1
class CreadorConcreto1(CreadorProducto):
    def crear_producto(self) -> Producto:
        return ProductoConcreto1()

# CreadorConcreto2 implementa CreadorProducto para crear ProductoConcreto2
class CreadorConcreto2(CreadorProducto):
    def crear_producto(self) -> Producto:
        return ProductoConcreto2()

# Cliente
def cliente_code(creador: CreadorProducto) -> None:
    print(f"Cliente: No conozco la clase del creador, pero funciona.\n"
          f"{creador.alguna_operacion()}")

# Uso
if __name__ == "__main__":
    print("App: Lanzada con CreadorConcreto1.")
    cliente_code(CreadorConcreto1())
    print("\n")

    print("App: Lanzada con CreadorConcreto2.")
    cliente_code(CreadorConcreto2())
```

### Abstract Factory

El patrón Abstract Factory proporciona una interfaz para crear familias de objetos relacionados sin especificar sus clases concretas.

```python
from abc import ABC, abstractmethod

# Productos abstractos
class ProductoA(ABC):
    @abstractmethod
    def operacion_a(self) -> str:
        pass

class ProductoB(ABC):
    @abstractmethod
    def operacion_b(self) -> str:
        pass
    
    @abstractmethod
    def otra_operacion_b(self, colaborador: ProductoA) -> str:
        pass

# Productos concretos
class ProductoA1(ProductoA):
    def operacion_a(self) -> str:
        return "Producto A1"

class ProductoA2(ProductoA):
    def operacion_a(self) -> str:
        return "Producto A2"

class ProductoB1(ProductoB):
    def operacion_b(self) -> str:
        return "Producto B1"
    
    def otra_operacion_b(self, colaborador: ProductoA) -> str:
        resultado = colaborador.operacion_a()
        return f"Resultado de B1 colaborando con ({resultado})"

class ProductoB2(ProductoB):
    def operacion_b(self) -> str:
        return "Producto B2"
    
    def otra_operacion_b(self, colaborador: ProductoA) -> str:
        resultado = colaborador.operacion_a()
        return f"Resultado de B2 colaborando con ({resultado})"

# Fábrica abstracta
class FabricaAbstracta(ABC):
    @abstractmethod
    def crear_producto_a(self) -> ProductoA:
        pass
    
    @abstractmethod
    def crear_producto_b(self) -> ProductoB:
        pass

# Fábricas concretas
class FabricaConcreta1(FabricaAbstracta):
    def crear_producto_a(self) -> ProductoA:
        return ProductoA1()
    
    def crear_producto_b(self) -> ProductoB:
        return ProductoB1()

class FabricaConcreta2(FabricaAbstracta):
    def crear_producto_a(self) -> ProductoA:
        return ProductoA2()
    
    def crear_producto_b(self) -> ProductoB:
        return ProductoB2()

# Cliente
def cliente_code(fabrica: FabricaAbstracta) -> None:
    producto_a = fabrica.crear_producto_a()
    producto_b = fabrica.crear_producto_b()
    
    print(f"{producto_b.operacion_b()}")
    print(f"{producto_b.otra_operacion_b(producto_a)}")

# Uso
if __name__ == "__main__":
    print("Cliente: Probando código cliente con el primer tipo de fábrica")
    cliente_code(FabricaConcreta1())
    
    print("\nCliente: Probando el mismo código cliente con el segundo tipo de fábrica")
    cliente_code(FabricaConcreta2())
```

### Builder

El patrón Builder separa la construcción de un objeto complejo de su representación, permitiendo el mismo proceso de construcción para crear diferentes representaciones.

```python
from abc import ABC, abstractmethod
from typing import Any

# Producto
class Producto:
    def __init__(self) -> None:
        self.partes = []
    
    def agregar(self, parte: Any) -> None:
        self.partes.append(parte)
    
    def listar_partes(self) -> None:
        print(f"Partes del producto: {', '.join(self.partes)}")

# Builder
class Builder(ABC):
    @abstractmethod
    def reset(self) -> None:
        pass
    
    @abstractmethod
    def construir_parte_a(self) -> None:
        pass
    
    @abstractmethod
    def construir_parte_b(self) -> None:
        pass
    
    @abstractmethod
    def construir_parte_c(self) -> None:
        pass

# Builder concreto
class BuilderConcreto(Builder):
    def __init__(self) -> None:
        self.reset()
    
    def reset(self) -> None:
        self._producto = Producto()
    
    def construir_parte_a(self) -> None:
        self._producto.agregar("Parte A")
    
    def construir_parte_b(self) -> None:
        self._producto.agregar("Parte B")
    
    def construir_parte_c(self) -> None:
        self._producto.agregar("Parte C")
    
    def obtener_producto(self) -> Producto:
        producto = self._producto
        self.reset()
        return producto

# Director
class Director:
    def __init__(self) -> None:
        self._builder = None
    
    def set_builder(self, builder: Builder) -> None:
        self._builder = builder
    
    def construir_producto_minimo(self) -> None:
        self._builder.construir_parte_a()
    
    def construir_producto_completo(self) -> None:
        self._builder.construir_parte_a()
        self._builder.construir_parte_b()
        self._builder.construir_parte_c()

# Versión con método encadenado (fluent interface)
class FluentBuilder:
    def __init__(self) -> None:
        self.reset()
    
    def reset(self) -> None:
        self._producto = Producto()
        return self
    
    def parte_a(self, parte: str = "Parte A") -> 'FluentBuilder':
        self._producto.agregar(parte)
        return self
    
    def parte_b(self, parte: str = "Parte B") -> 'FluentBuilder':
        self._producto.agregar(parte)
        return self
    
    def parte_c(self, parte: str = "Parte C") -> 'FluentBuilder':
        self._producto.agregar(parte)
        return self
    
    def build(self) -> Producto:
        producto = self._producto
        self.reset()
        return producto

# Uso
if __name__ == "__main__":
    # Usando el patrón Builder tradicional
    director = Director()
    builder = BuilderConcreto()
    director.set_builder(builder)
    
    print("Producto mínimo:")
    director.construir_producto_minimo()
    builder.obtener_producto().listar_partes()
    
    print("\nProducto completo:")
    director.construir_producto_completo()
    builder.obtener_producto().listar_partes()
    
    # Usando el builder sin director
    print("\nProducto personalizado:")
    builder.construir_parte_a()
    builder.construir_parte_c()
    builder.obtener_producto().listar_partes()
    
    # Usando el builder con método encadenado
    print("\nProducto fluent:")
    producto_fluent = FluentBuilder()\
        .parte_a("Parte A personalizada")\
        .parte_b("Parte B personalizada")\
        .parte_c("Parte C personalizada")\
        .build()
    producto_fluent.listar_partes()
```

### Singleton

El patrón Singleton garantiza que una clase tenga solo una instancia y proporciona un punto de acceso global a ella. En Python, hay varias formas de implementarlo.

```python
# Implementación básica
class Singleton:
    _instancia = None
    
    def __new__(cls):
        if cls._instancia is None:
            print("Creando instancia singleton")
            cls._instancia = super(Singleton, cls).__new__(cls)
        return cls._instancia

# Implementación con decorador
def singleton(clase):
    instancias = {}
    def obtener_instancia(*args, **kwargs):
        if clase not in instancias:
            instancias[clase] = clase(*args, **kwargs)
        return instancias[clase]
    return obtener_instancia

@singleton
class SingletonDecorado:
    def __init__(self, valor=None):
        self.valor = valor

# Implementación con metaclase
class SingletonMetaclass(type):
    _instancias = {}
    def __call__(cls, *args, **kwargs):
        if cls not in cls._instancias:
            cls._instancias[cls] = super(SingletonMetaclass, cls).__call__(*args, **kwargs)
        return cls._instancias[cls]

class SingletonConMetaclase(metaclass=SingletonMetaclass):
    def __init__(self, valor=None):
        self.valor = valor

# Uso
if __name__ == "__main__":
    # Usando la implementación básica
    s1 = Singleton()
    s2 = Singleton()
    print(f"¿Son la misma instancia? {s1 is s2}")
    
    # Usando la implementación con decorador
    sd1 = SingletonDecorado(valor="A")
    sd2 = SingletonDecorado(valor="B")  # El valor no cambia porque ya existe una instancia
    print(f"Valor de sd1: {sd1.valor}")
    print(f"Valor de sd2: {sd2.valor}")
    print(f"¿Son la misma instancia? {sd1 is sd2}")
    
    # Usando la implementación con metaclase
    sm1 = SingletonConMetaclase(valor="X")
    sm2 = SingletonConMetaclase(valor="Y")  # El valor no cambia porque ya existe una instancia
    print(f"Valor de sm1: {sm1.valor}")
    print(f"Valor de sm2: {sm2.valor}")
    print(f"¿Son la misma instancia? {sm1 is sm2}")
```

### Object Pool

El patrón Object Pool mantiene un conjunto de objetos inicializados listos para usar, en lugar de asignarlos y destruirlos bajo demanda.

```python
import threading
from typing import List, Optional, TypeVar, Generic, Callable

T = TypeVar('T')

class ObjectPool(Generic[T]):
    def __init__(self, create_func: Callable[[], T], max_size: int = 10):
        self.create_func = create_func
        self.max_size = max_size
        self.lock = threading.RLock()
        self.objects: List[T] = []
        self.in_use: List[bool] = []
    
    def acquire(self) -> Optional[T]:
        with self.lock:
            # Buscar un objeto disponible
            for i, in_use in enumerate(self.in_use):
                if not in_use:
                    self.in_use[i] = True
                    return self.objects[i]
            
            # Si no hay objetos disponibles pero no hemos alcanzado el máximo, crear uno nuevo
            if len(self.objects) < self.max_size:
                obj = self.create_func()
                self.objects.append(obj)
                self.in_use.append(True)
                return obj
            
            # Si llegamos aquí, el pool está lleno y todos los objetos están en uso
            return None
    
    def release(self, obj: T) -> None:
        with self.lock:
            try:
                idx = self.objects.index(obj)
                self.in_use[idx] = False
            except ValueError:
                # El objeto no está en el pool
                pass
    
    def get_stats(self) -> tuple:
        with self.lock:
            created = len(self.objects)
            available = self.in_use.count(False)
            in_use = self.in_use.count(True)
            return created, available, in_use

# Ejemplo de uso
class ExpensiveObject:
    def __init__(self, id_num: int = None):
        self.id = id_num
    
    def __str__(self) -> str:
        return f"ExpensiveObject(id={self.id})"

# Uso
if __name__ == "__main__":
    import time
    import random
    import concurrent.futures
    
    # Contador para asignar IDs únicos a los objetos
    counter = 0
    
    def create_expensive_object():
        global counter
        counter += 1
        print(f"Creando objeto #{counter}")
        return ExpensiveObject(counter)
    
    # Crear un pool con máximo 3 objetos
    pool = ObjectPool(create_expensive_object, 3)
    
    def worker(worker_id: int):
        # Adquirir un objeto
        obj = pool.acquire()
        if obj is None:
            print(f"Trabajador {worker_id} no pudo obtener un objeto")
            return
        
        try:
            # Simular trabajo
            print(f"Trabajador {worker_id} usando {obj}")
            time.sleep(random.uniform(0.5, 1.5))
        finally:
            # Devolver el objeto al pool
            pool.release(obj)
            print(f"Trabajador {worker_id} terminó de usar {obj}")
    
    # Simular uso concurrente del pool
    with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
        # Lanzar 5 trabajadores (más que el tamaño del pool)
        futures = [executor.submit(worker, i) for i in range(5)]
        concurrent.futures.wait(futures)
    
    # Mostrar estadísticas finales
    created, available, in_use = pool.get_stats()
    print(f"\nEstadísticas finales - Creados: {created}, Disponibles: {available}, En uso: {in_use}")
```

## Patrones Estructurales

Los patrones estructurales se ocupan de cómo se componen las clases y los objetos para formar estructuras más grandes.

### Adapter

El patrón Adapter permite que interfaces incompatibles trabajen juntas, convirtiendo la interfaz de una clase en otra que el cliente espera.

```python
from abc import ABC, abstractmethod

# Interfaz objetivo que el cliente utiliza
class ServicioObjetivo(ABC):
    @abstractmethod
    def solicitud(self) -> str:
        pass

# Cliente que usa ServicioObjetivo
class ServicioCliente:
    def usar_servicio(self, servicio: ServicioObjetivo) -> None:
        print(servicio.solicitud())

# Servicio existente con interfaz incompatible
class ServicioExistente:
    def solicitud_especifica(self) -> str:
        return "Respuesta del servicio existente con interfaz incompatible"

# Adaptador que implementa ServicioObjetivo y adapta ServicioExistente
class Adaptador(ServicioObjetivo):
    def __init__(self, servicio_existente: ServicioExistente) -> None:
        self.servicio_existente = servicio_existente
    
    def solicitud(self) -> str:
        # Adapta la llamada a la interfaz del servicio existente
        return f"Adaptador: {self.servicio_existente.solicitud_especifica()}"

# Uso
if __name__ == "__main__":
    cliente = ServicioCliente()
    servicio_existente = ServicioExistente()
    adaptador = Adaptador(servicio_existente)
    
    print("Cliente: Puedo trabajar bien con el adaptador:")
    cliente.usar_servicio(adaptador)
```

### Bridge

El patrón Bridge separa una abstracción de su implementación para que ambas puedan variar independientemente.

```python
from abc import ABC, abstractmethod

# Implementador define la interfaz para las clases de implementación
class Implementador(ABC):
    @abstractmethod
    def operacion_impl(self) -> str:
        pass

# ImplementadorConcreto1 es una implementación concreta
class ImplementadorConcreto1(Implementador):
    def operacion_impl(self) -> str:
        return "ImplementadorConcreto1: Resultado de la operación"

# ImplementadorConcreto2 es otra implementación concreta
class ImplementadorConcreto2(Implementador):
    def operacion_impl(self) -> str:
        return "ImplementadorConcreto2: Resultado de la operación"

# Abstraccion define la interfaz para la parte de control
class Abstraccion:
    def __init__(self, implementador: Implementador) -> None:
        self.implementador = implementador
    
    def operacion(self) -> str:
        return f"Abstracción básica: {self.implementador.operacion_impl()}"

# AbstraccionExtendida es una variante de la abstracción
class AbstraccionExtendida(Abstraccion):
    def operacion(self) -> str:
        return f"Abstracción extendida: {self.implementador.operacion_impl()}"

# Uso
if __name__ == "__main__":
    implementador1 = ImplementadorConcreto1()
    implementador2 = ImplementadorConcreto2()
    
    abstraccion1 = Abstraccion(implementador1)
    abstraccion2 = AbstraccionExtendida(implementador2)
    
    print(abstraccion1.operacion())
    print(abstraccion2.operacion())
    
    # Cambiar la implementación en tiempo de ejecución
    abstraccion1.implementador = implementador2
    print(abstraccion1.operacion())
```

### Composite

El patrón Composite compone objetos en estructuras de árbol para representar jerarquías parte-todo, permitiendo a los clientes tratar objetos individuales y composiciones de objetos de manera uniforme.

```python
from abc import ABC, abstractmethod
from typing import List

# Componente define la interfaz común para objetos simples y compuestos
class Componente(ABC):
    @abstractmethod
    def operacion(self) -> str:
        pass
    
    def agregar(self, componente: 'Componente') -> None:
        pass
    
    def eliminar(self, componente: 'Componente') -> None:
        pass
    
    def es_compuesto(self) -> bool:
        return False

# Hoja representa objetos finales de la composición que no tienen hijos
class Hoja(Componente):
    def __init__(self, nombre: str) -> None:
        self.nombre = nombre
    
    def operacion(self) -> str:
        return f"Hoja {self.nombre}"

# Compuesto representa objetos complejos que pueden tener hijos
class Compuesto(Componente):
    def __init__(self, nombre: str) -> None:
        self.nombre = nombre
        self.hijos: List[Componente] = []
    
    def operacion(self) -> str:
        resultados = []
        for hijo in self.hijos:
            resultados.append(hijo.operacion())
        return f"Compuesto {self.nombre} ({' + '.join(resultados)})"
    
    def agregar(self, componente: Componente) -> None:
        self.hijos.append(componente)
    
    def eliminar(self, componente: Componente) -> None:
        self.hijos.remove(componente)
    
    def es_compuesto(self) -> bool:
        return True

# Uso
if __name__ == "__main__":
    # Crear hojas
    hoja1 = Hoja("A")
    hoja2 = Hoja("B")
    hoja3 = Hoja("C")
    
    # Crear compuestos
    compuesto1 = Compuesto("X")
    compuesto1.agregar(hoja1)
    compuesto1.agregar(hoja2)
    
    compuesto2 = Compuesto("Y")
    compuesto2.agregar(hoja3)
    
    # Crear compuesto raíz
    compuesto_raiz = Compuesto("RAÍZ")
    compuesto_raiz.agregar(compuesto1)
    compuesto_raiz.agregar(compuesto2)
    
    # Mostrar la estructura
    print(compuesto_raiz.operacion())
    
    # Modificar la estructura
    compuesto1.eliminar(hoja1)
    compuesto_raiz.agregar(hoja1)
    
    # Mostrar la estructura modificada
    print(compuesto_raiz.operacion())
```

### Decorator

El patrón Decorator añade responsabilidades adicionales a un objeto dinámicamente, proporcionando una alternativa flexible a la herencia para extender la funcionalidad.

```python
from abc import ABC, abstractmethod

# Componente define la interfaz para objetos que pueden tener responsabilidades añadidas
class Componente(ABC):
    @abstractmethod
    def operacion(self) -> str:
        pass

# ComponenteConcreto implementa la interfaz Componente
class ComponenteConcreto(Componente):
    def operacion(self) -> str:
        return "Componente Concreto"

# Decorador implementa la interfaz Componente y mantiene una referencia a un objeto Componente
class Decorador(Componente):
    def __init__(self, componente: Componente) -> None:
        self._componente = componente
    
    @abstractmethod
    def operacion(self) -> str:
        pass

# DecoradorConcretoA extiende Decorador
class DecoradorConcretoA(Decorador):
    def operacion(self) -> str:
        return f"DecoradorConcretoA({self._componente.operacion()})"

# DecoradorConcretoB extiende Decorador
class DecoradorConcretoB(Decorador):
    def operacion(self) -> str:
        return f"DecoradorConcretoB({self._componente.operacion()})"

# Uso
if __name__ == "__main__":
    # Crear un componente simple
    componente = ComponenteConcreto()
    print("Componente simple:")
    print(componente.operacion())
    
    # Decorar el componente con DecoradorConcretoA
    decorador_a = DecoradorConcretoA(componente)
    print("\nComponente decorado con A:")
    print(decorador_a.operacion())
    
    # Decorar el componente con DecoradorConcretoB
    decorador_b = DecoradorConcretoB(componente)
    print("\nComponente decorado con B:")
    print(decorador_b.operacion())
    
    # Decorar el componente con ambos decoradores
    decorador_ba = DecoradorConcretoB(DecoradorConcretoA(componente))
    print("\nComponente decorado con B y A:")
    print(decorador_ba.operacion())
```

### Facade

El patrón Facade proporciona una interfaz unificada a un conjunto de interfaces en un subsistema, definiendo una interfaz de nivel superior que hace que el subsistema sea más fácil de usar.

```python
# Subsistema1 representa una parte compleja del sistema
class Subsistema1:
    def operacion1(self) -> str:
        return "Subsistema1: Listo!"
    
    def operacion_n(self) -> str:
        return "Subsistema1: Ejecutando..."

# Subsistema2 representa otra parte compleja del sistema
class Subsistema2:
    def operacion1(self) -> str:
        return "Subsistema2: Preparado!"
    
    def operacion_z(self) -> str:
        return "Subsistema2: Procesando..."

# Fachada proporciona una interfaz simple para la lógica compleja de los subsistemas
class Fachada:
    def __init__(self) -> None:
        self.subsistema1 = Subsistema1()
        self.subsistema2 = Subsistema2()
    
    def operacion_simple(self) -> str:
        resultado = ["Fachada inicializa subsistemas:"]
        resultado.append(self.subsistema1.operacion1())
        resultado.append(self.subsistema2.operacion1())
        resultado.append("Fachada ordena a los subsistemas realizar operaciones:")
        resultado.append(self.subsistema1.operacion_n())
        resultado.append(self.subsistema2.operacion_z())
        return "\n".join(resultado)

# Uso
if __name__ == "__main__":
    fachada = Fachada()
    
    # El cliente interactúa con el sistema a través de la fachada
    print(fachada.operacion_simple())
```

## Patrones de Comportamiento

Los patrones de comportamiento se ocupan de la comunicación entre objetos.

### Strategy

El patrón Strategy define una familia de algoritmos, encapsula cada uno de ellos y los hace intercambiables, permitiendo que el algoritmo varíe independientemente de los clientes que lo utilizan.

```python
from abc import ABC, abstractmethod
from typing import List

# Estrategia define la interfaz común para todos los algoritmos soportados
class Estrategia(ABC):
    @abstractmethod
    def ejecutar(self, a: int, b: int) -> int:
        pass

# EstrategiaSuma implementa la operación de suma
class EstrategiaSuma(Estrategia):
    def ejecutar(self, a: int, b: int) -> int:
        return a + b

# EstrategiaResta implementa la operación de resta
class EstrategiaResta(Estrategia):
    def ejecutar(self, a: int, b: int) -> int:
        return a - b

# EstrategiaMultiplicacion implementa la operación de multiplicación
class EstrategiaMultiplicacion(Estrategia):
    def ejecutar(self, a: int, b: int) -> int:
        return a * b

# Contexto utiliza una estrategia para ejecutar una operación
class Contexto:
    def __init__(self, estrategia: Estrategia = None) -> None:
        self._estrategia = estrategia
    
    @property
    def estrategia(self) -> Estrategia:
        return self._estrategia
    
    @estrategia.setter
    def estrategia(self, estrategia: Estrategia) -> None:
        self._estrategia = estrategia
    
    def ejecutar_estrategia(self, a: int, b: int) -> int:
        if self._estrategia is None:
            raise Exception("Estrategia no establecida")
        return self._estrategia.ejecutar(a, b)

# Uso
if __name__ == "__main__":
    # Crear el contexto con una estrategia inicial
    contexto = Contexto(EstrategiaSuma())
    
    # Ejecutar la estrategia
    resultado = contexto.ejecutar_estrategia(10, 5)
    print(f"10 + 5 = {resultado}")
    
    # Cambiar la estrategia y ejecutar de nuevo
    contexto.estrategia = EstrategiaResta()
    resultado = contexto.ejecutar_estrategia(10, 5)
    print(f"10 - 5 = {resultado}")
    
    # Cambiar a otra estrategia
    contexto.estrategia = EstrategiaMultiplicacion()
    resultado = contexto.ejecutar_estrategia(10, 5)
    print(f"10 * 5 = {resultado}")
```

### Observer

El patrón Observer define una dependencia uno-a-muchos entre objetos, de modo que cuando un objeto cambia de estado, todos sus dependientes son notificados y actualizados automáticamente.

```python
from abc import ABC, abstractmethod
from typing import List

# Observador define la interfaz para todos los observadores
class Observador(ABC):
    @abstractmethod
    def actualizar(self, mensaje: str) -> None:
        pass

# Sujeto mantiene una lista de observadores y los notifica de cambios
class Sujeto:
    def __init__(self) -> None:
        self._observadores: List[Observador] = []
        self._estado: str = None
    
    def adjuntar(self, observador: Observador) -> None:
        if observador not in self._observadores:
            self._observadores.append(observador)
    
    def separar(self, observador: Observador) -> None:
        self._observadores.remove(observador)
    
    def notificar(self) -> None:
        for observador in self._observadores:
            observador.actualizar(self._estado)
    
    @property
    def estado(self) -> str:
        return self._estado
    
    @estado.setter
    def estado(self, estado: str) -> None:
        self._estado = estado
        self.notificar()

# ObservadorConcreto implementa la interfaz Observador
class ObservadorConcreto(Observador):
    def __init__(self, nombre: str) -> None:
        self.nombre = nombre
    
    def actualizar(self, mensaje: str) -> None:
        print(f"Observador {self.nombre} recibió mensaje: {mensaje}")

# Uso
if __name__ == "__main__":
    # Crear el sujeto
    sujeto = Sujeto()
    
    # Crear observadores
    observador1 = ObservadorConcreto("A")
    observador2 = ObservadorConcreto("B")
    observador3 = ObservadorConcreto("C")
    
    # Adjuntar observadores al sujeto
    sujeto.adjuntar(observador1)
    sujeto.adjuntar(observador2)
    sujeto.adjuntar(observador3)
    
    # Cambiar el estado del sujeto
    print("Cambiando estado a 'Primer cambio'")
    sujeto.estado = "Primer cambio"
    
    # Separar un observador
    sujeto.separar(observador2)
    
    # Cambiar el estado de nuevo
    print("\nCambiando estado a 'Segundo cambio'")
    sujeto.estado = "Segundo cambio"
```

### Command

El patrón Command encapsula una solicitud como un objeto, permitiendo parametrizar clientes con diferentes solicitudes, encolar o registrar solicitudes, y soportar operaciones que pueden deshacerse.

```python
from abc import ABC, abstractmethod
from typing import List

# Comando define la interfaz para ejecutar una operación
class Comando(ABC):
    @abstractmethod
    def ejecutar(self) -> None:
        pass
    
    @abstractmethod
    def deshacer(self) -> None:
        pass

# Receptor sabe cómo realizar las operaciones asociadas con los comandos
class Receptor:
    def __init__(self) -> None:
        self._estado = "Inicial"
    
    def accion(self, estado: str) -> None:
        self._estado = estado
        print(f"Receptor: Mi estado ahora es {self._estado}")
    
    @property
    def estado(self) -> str:
        return self._estado

# ComandoConcreto implementa Comando y define la relación entre una acción y un receptor
class ComandoConcreto(Comando):
    def __init__(self, receptor: Receptor, estado: str) -> None:
        self._receptor = receptor
        self._estado_nuevo = estado
        self._estado_previo = None
    
    def ejecutar(self) -> None:
        self._estado_previo = self._receptor.estado
        self._receptor.accion(self._estado_nuevo)
    
    def deshacer(self) -> None:
        self._receptor.accion(self._estado_previo)

# Invocador pide al comando que ejecute la solicitud
class Invocador:
    def __init__(self) -> None:
        self._comandos: List[Comando] = []
        self._indice = -1
    
    def almacenar_y_ejecutar(self, comando: Comando) -> None:
        # Si hemos deshecho comandos, eliminar los comandos después del índice actual
        if self._indice < len(self._comandos) - 1:
            self._comandos = self._comandos[:self._indice + 1]
        
        # Añadir y ejecutar el nuevo comando
        self._comandos.append(comando)
        self._indice = len(self._comandos) - 1
        comando.ejecutar()
    
    def deshacer(self) -> None:
        if self._indice >= 0:
            self._comandos[self._indice].deshacer()
            self._indice -= 1
    
    def rehacer(self) -> None:
        if self._indice < len(self._comandos) - 1:
            self._indice += 1
            self._comandos[self._indice].ejecutar()

# Uso
if __name__ == "__main__":
    # Crear el receptor
    receptor = Receptor()
    
    # Crear el invocador
    invocador = Invocador()
    
    # Crear y ejecutar comandos
    comando1 = ComandoConcreto(receptor, "Estado 1")
    invocador.almacenar_y_ejecutar(comando1)
    
    comando2 = ComandoConcreto(receptor, "Estado 2")
    invocador.almacenar_y_ejecutar(comando2)
    
    comando3 = ComandoConcreto(receptor, "Estado 3")
    invocador.almacenar_y_ejecutar(comando3)
    
    # Deshacer comandos
    print("\nDeshaciendo...")
    invocador.deshacer()
    print(f"Estado actual: {receptor.estado}")
    
    invocador.deshacer()
    print(f"Estado actual: {receptor.estado}")
    
    # Rehacer comandos
    print("\nRehaciendo...")
    invocador.rehacer()
    print(f"Estado actual: {receptor.estado}")
    
    invocador.rehacer()
    print(f"Estado actual: {receptor.estado}")
```

## Patrones Específicos de Python

Además de los patrones de diseño tradicionales, Python tiene algunos patrones específicos que aprovechan sus características únicas.

### Descriptor

El patrón Descriptor permite a una clase delegar el acceso a sus atributos a otra clase.

```python
class Validador:
    def __init__(self, min_valor=None, max_valor=None):
        self.min_valor = min_valor
        self.max_valor = max_valor
        self.nombre = None
    
    def __set_name__(self, propietario, nombre):
        self.nombre = nombre
    
    def __get__(self, instancia, propietario):
        if instancia is None:
            return self
        return instancia.__dict__.get(self.nombre, None)
    
    def __set__(self, instancia, valor):
        if self.min_valor is not None and valor < self.min_valor:
            raise ValueError(f"{self.nombre} debe ser >= {self.min_valor}")
        if self.max_valor is not None and valor > self.max_valor:
            raise ValueError(f"{self.nombre} debe ser <= {self.max_valor}")
        instancia.__dict__[self.nombre] = valor

class Persona:
    edad = Validador(0, 120)
    altura = Validador(0, 250)
    peso = Validador(0, 300)
    
    def __init__(self, nombre, edad, altura, peso):
        self.nombre = nombre
        self.edad = edad
        self.altura = altura
        self.peso = peso
    
    def __str__(self):
        return f"Persona(nombre={self.nombre}, edad={self.edad}, altura={self.altura}, peso={self.peso})"

# Uso
if __name__ == "__main__":
    # Crear una persona con valores válidos
    persona = Persona("Juan", 30, 175, 70)
    print(persona)
    
    # Intentar establecer un valor inválido
    try:
        persona.edad = 150  # Esto lanzará una excepción
    except ValueError as e:
        print(f"Error: {e}")
    
    # Modificar un valor válido
    persona.edad = 35
    print(f"Nueva edad: {persona.edad}")
```

### Context Manager

El patrón Context Manager permite encapsular operaciones comunes que se realizan antes y después de un bloque de código.

```python
class ManejoArchivo:
    def __init__(self, nombre_archivo, modo):
        self.nombre_archivo = nombre_archivo
        self.modo = modo
        self.archivo = None
    
    def __enter__(self):
        self.archivo = open(self.nombre_archivo, self.modo)
        return self.archivo
    
    def __exit__(self, tipo_excepcion, valor_excepcion, traceback):
        self.archivo.close()
        # Retornar True para suprimir la excepción, False para propagarla
        return False

# Implementación con decorador contextlib
from contextlib import contextmanager

@contextmanager
def manejo_archivo(nombre_archivo, modo):
    try:
        archivo = open(nombre_archivo, modo)
        yield archivo
    finally:
        archivo.close()

# Uso
if __name__ == "__main__":
    # Usando la clase
    with ManejoArchivo("ejemplo.txt", "w") as f:
        f.write("Hola, mundo!")
    
    # Usando el decorador
    with manejo_archivo("ejemplo2.txt", "w") as f:
        f.write("Hola de nuevo!")
    
    # Verificar que los archivos se cerraron correctamente
    print("Archivos creados y cerrados correctamente.")
```

### Mixin

El patrón Mixin permite añadir funcionalidades a clases sin usar herencia múltiple de manera compleja.

```python
class SerializableMixin:
    def to_dict(self):
        return {key: value for key, value in self.__dict__.items() if not key.startswith('_')}
    
    def to_json(self):
        import json
        return json.dumps(self.to_dict())

class LoggableMixin:
    def log(self, mensaje):
        print(f"[LOG] {self.__class__.__name__}: {mensaje}")

class Producto(SerializableMixin, LoggableMixin):
    def __init__(self, id, nombre, precio):
        self.id = id
        self.nombre = nombre
        self.precio = precio
    
    def aplicar_descuento(self, porcentaje):
        descuento = self.precio * (porcentaje / 100)
        self.precio -= descuento
        self.log(f"Aplicado descuento de {porcentaje}%. Nuevo precio: {self.precio}")

# Uso
if __name__ == "__main__":
    producto = Producto(1, "Laptop", 1000)
    
    # Usar funcionalidad de LoggableMixin
    producto.log("Producto creado")
    
    # Aplicar descuento (método propio que usa el mixin)
    producto.aplicar_descuento(10)
    
    # Usar funcionalidad de SerializableMixin
    print(producto.to_dict())
    print(producto.to_json())
```

### Metaclase

Las metaclases permiten controlar la creación de clases.

```python
class MetaSingleton(type):
    _instancias = {}
    
    def __call__(cls, *args, **kwargs):
        if cls not in cls._instancias:
            cls._instancias[cls] = super(MetaSingleton, cls).__call__(*args, **kwargs)
        return cls._instancias[cls]

class RegistroClases(type):
    clases = {}
    
    def __new__(mcs, nombre, bases, atributos):
        clase = super().__new__(mcs, nombre, bases, atributos)
        mcs.clases[nombre] = clase
        return clase

class ValidacionAtributos(type):
    def __new__(mcs, nombre, bases, atributos):
        # Verificar que todos los métodos tienen docstrings
        for nombre_attr, attr in atributos.items():
            if callable(attr) and not nombre_attr.startswith('__') and not attr.__doc__:
                raise TypeError(f"El método {nombre_attr} debe tener un docstring")
        
        return super().__new__(mcs, nombre, bases, atributos)

# Uso de MetaSingleton
class BaseDatos(metaclass=MetaSingleton):
    def __init__(self):
        self.conexion = "Conexión a la base de datos establecida"

# Uso de RegistroClases
class Animal(metaclass=RegistroClases):
    pass

class Perro(Animal):
    pass

class Gato(Animal):
    pass

# Uso de ValidacionAtributos
class ClaseConDocstrings(metaclass=ValidacionAtributos):
    def metodo_con_docstring(self):
        """Este método tiene un docstring."""
        pass

# Esto lanzará un error porque el método no tiene docstring
try:
    class ClaseSinDocstrings(metaclass=ValidacionAtributos):
        def metodo_sin_docstring(self):
            pass
except TypeError as e:
    print(f"Error esperado: {e}")

# Uso
if __name__ == "__main__":
    # Probar MetaSingleton
    db1 = BaseDatos()
    db2 = BaseDatos()
    print(f"¿Son la misma instancia? {db1 is db2}")
    
    # Probar RegistroClases
    print("Clases registradas:")
    for nombre, clase in RegistroClases.clases.items():
        print(f"  - {nombre}")
    
    # Probar ValidacionAtributos
    obj = ClaseConDocstrings()
    obj.metodo_con_docstring()
```

## Ejemplo Práctico: Sistema de Notificaciones

A continuación, se presenta un ejemplo práctico que combina varios patrones de diseño para crear un sistema de notificaciones flexible y extensible.

```python
from abc import ABC, abstractmethod
from typing import List, Dict, Any

# Observer Pattern
class Observador(ABC):
    @abstractmethod
    def actualizar(self, evento: str, datos: Dict[str, Any]) -> None:
        pass

class Sujeto:
    def __init__(self) -> None:
        self._observadores: Dict[str, List[Observador]] = {}
    
    def registrar(self, evento: str, observador: Observador) -> None:
        if evento not in self._observadores:
            self._observadores[evento] = []
        if observador not in self._observadores[evento]:
            self._observadores[evento].append(observador)
    
    def eliminar(self, evento: str, observador: Observador) -> None:
        if evento in self._observadores and observador in self._observadores[evento]:
            self._observadores[evento].remove(observador)
    
    def notificar(self, evento: str, datos: Dict[str, Any]) -> None:
        if evento in self._observadores:
            for observador in self._observadores[evento]:
                observador.actualizar(evento, datos)

# Strategy Pattern
class EstrategiaNotificacion(ABC):
    @abstractmethod
    def enviar(self, destinatario: str, mensaje: str) -> None:
        pass

class NotificacionEmail(EstrategiaNotificacion):
    def enviar(self, destinatario: str, mensaje: str) -> None:
        print(f"Enviando email a {destinatario}: {mensaje}")

class NotificacionSMS(EstrategiaNotificacion):
    def enviar(self, destinatario: str, mensaje: str) -> None:
        print(f"Enviando SMS a {destinatario}: {mensaje}")

class NotificacionPush(EstrategiaNotificacion):
    def enviar(self, destinatario: str, mensaje: str) -> None:
        print(f"Enviando notificación push a {destinatario}: {mensaje}")

# Factory Method Pattern
class CreadorNotificacion(ABC):
    @abstractmethod
    def crear_notificador(self) -> EstrategiaNotificacion:
        pass

class CreadorNotificacionEmail(CreadorNotificacion):
    def crear_notificador(self) -> EstrategiaNotificacion:
        return NotificacionEmail()

class CreadorNotificacionSMS(CreadorNotificacion):
    def crear_notificador(self) -> EstrategiaNotificacion:
        return NotificacionSMS()

class CreadorNotificacionPush(CreadorNotificacion):
    def crear_notificador(self) -> EstrategiaNotificacion:
        return NotificacionPush()

# Decorator Pattern
class DecoradorNotificacion(EstrategiaNotificacion):
    def __init__(self, notificador: EstrategiaNotificacion) -> None:
        self._notificador = notificador
    
    def enviar(self, destinatario: str, mensaje: str) -> None:
        self._notificador.enviar(destinatario, mensaje)

class LoggingDecorador(DecoradorNotificacion):
    def enviar(self, destinatario: str, mensaje: str) -> None:
        print(f"[LOG] Enviando notificación a {destinatario}")
        super().enviar(destinatario, mensaje)
        print(f"[LOG] Notificación enviada a {destinatario}")

class RetryDecorador(DecoradorNotificacion):
    def __init__(self, notificador: EstrategiaNotificacion, max_intentos: int = 3) -> None:
        super().__init__(notificador)
        self.max_intentos = max_intentos
    
    def enviar(self, destinatario: str, mensaje: str) -> None:
        intentos = 0
        while intentos < self.max_intentos:
            try:
                super().enviar(destinatario, mensaje)
                break
            except Exception as e:
                intentos += 1
                print(f"[RETRY] Intento {intentos} fallido: {e}")
                if intentos == self.max_intentos:
                    print(f"[RETRY] Máximo de intentos alcanzado. No se pudo enviar la notificación.")

# Singleton Pattern (usando metaclase)
class Singleton(type):
    _instancias = {}
    
    def __call__(cls, *args, **kwargs):
        if cls not in cls._instancias:
            cls._instancias[cls] = super(Singleton, cls).__call__(*args, **kwargs)
        return cls._instancias[cls]

# Sistema de Notificaciones (Facade Pattern)
class SistemaNotificaciones(metaclass=Singleton):
    def __init__(self) -> None:
        self.sujeto = Sujeto()
        self.creadores: Dict[str, CreadorNotificacion] = {
            "email": CreadorNotificacionEmail(),
            "sms": CreadorNotificacionSMS(),
            "push": CreadorNotificacionPush()
        }
    
    def registrar_observador(self, evento: str, observador: Observador) -> None:
        self.sujeto.registrar(evento, observador)
    
    def eliminar_observador(self, evento: str, observador: Observador) -> None:
        self.sujeto.eliminar(evento, observador)
    
    def notificar_evento(self, evento: str, datos: Dict[str, Any]) -> None:
        self.sujeto.notificar(evento, datos)
    
    def crear_notificador(self, tipo: str, con_logging: bool = False, con_retry: bool = False) -> EstrategiaNotificacion:
        if tipo not in self.creadores:
            raise ValueError(f"Tipo de notificación no soportado: {tipo}")
        
        notificador = self.creadores[tipo].crear_notificador()
        
        if con_logging:
            notificador = LoggingDecorador(notificador)
        
        if con_retry:
            notificador = RetryDecorador(notificador)
        
        return notificador
    
    def enviar_notificacion(self, tipo: str, destinatario: str, mensaje: str, con_logging: bool = False, con_retry: bool = False) -> None:
        notificador = self.crear_notificador(tipo, con_logging, con_retry)
        notificador.enviar(destinatario, mensaje)

# Observadores concretos
class ObservadorUsuario(Observador):
    def __init__(self, nombre: str, email: str, telefono: str, token_dispositivo: str) -> None:
        self.nombre = nombre
        self.email = email
        self.telefono = telefono
        self.token_dispositivo = token_dispositivo
        self.sistema = SistemaNotificaciones()
    
    def actualizar(self, evento: str, datos: Dict[str, Any]) -> None:
        if evento == "nuevo_mensaje":
            self.sistema.enviar_notificacion("email", self.email, f"Nuevo mensaje: {datos['contenido']}", con_logging=True)
        elif evento == "alerta_seguridad":
            self.sistema.enviar_notificacion("sms", self.telefono, f"Alerta de seguridad: {datos['contenido']}", con_logging=True, con_retry=True)
            self.sistema.enviar_notificacion("push", self.token_dispositivo, f"Alerta de seguridad: {datos['contenido']}")

class ObservadorAdmin(Observador):
    def __init__(self, nombre: str, email: str) -> None:
        self.nombre = nombre
        self.email = email
        self.sistema = SistemaNotificaciones()
    
    def actualizar(self, evento: str, datos: Dict[str, Any]) -> None:
        if evento == "error_sistema":
            self.sistema.enviar_notificacion("email", self.email, f"Error en el sistema: {datos['contenido']}", con_logging=True, con_retry=True)

# Uso del sistema
if __name__ == "__main__":
    # Obtener la instancia del sistema de notificaciones
    sistema = SistemaNotificaciones()
    
    # Crear observadores
    usuario1 = ObservadorUsuario("Juan", "juan@example.com", "+1234567890", "token123")
    usuario2 = ObservadorUsuario("María", "maria@example.com", "+0987654321", "token456")
    admin = ObservadorAdmin("Admin", "admin@example.com")
    
    # Registrar observadores para diferentes eventos
    sistema.registrar_observador("nuevo_mensaje", usuario1)
    sistema.registrar_observador("nuevo_mensaje", usuario2)
    sistema.registrar_observador("alerta_seguridad", usuario1)
    sistema.registrar_observador("error_sistema", admin)
    
    # Generar eventos
    print("\n=== Generando evento: nuevo_mensaje ===")
    sistema.notificar_evento("nuevo_mensaje", {"contenido": "Hola a todos!"})  
    
    print("\n=== Generando evento: alerta_seguridad ===")
    sistema.notificar_evento("alerta_seguridad", {"contenido": "Intento de acceso no autorizado"})  
    
    print("\n=== Generando evento: error_sistema ===")
    sistema.notificar_evento("error_sistema", {"contenido": "Error en la base de datos"})  
    
    # Eliminar un observador y generar otro evento
    sistema.eliminar_observador("nuevo_mensaje", usuario2)
    
    print("\n=== Generando evento después de eliminar un observador: nuevo_mensaje ===")
    sistema.notificar_evento("nuevo_mensaje", {"contenido": "Mensaje importante!"})  
```

## Conclusión

Los patrones de diseño son herramientas valiosas para resolver problemas comunes en el desarrollo de software. En Python, algunos patrones tradicionales se implementan de manera diferente debido a las características del lenguaje, mientras que otros patrones específicos de Python aprovechan sus características únicas como decoradores, generadores, comprensiones de listas y tipado dinámico.

Al aplicar patrones de diseño en Python, es importante mantener la simplicidad y claridad que caracteriza al lenguaje, siguiendo el principio "Pythonic" y evitando la sobreingeniería. Recuerda que los patrones son guías, no reglas estrictas, y deben adaptarse a las necesidades específicas de cada proyecto.

## Recursos Adicionales

- [Python Design Patterns: For Sleek And Fashionable Code](https://www.toptal.com/python/python-design-patterns)
- [Python Patterns](https://python-patterns.guide/)
- [Design Patterns in Python](https://refactoring.guru/design-patterns/python)
- [Python 3 Patterns, Recipes and Idioms](https://python-3-patterns-idioms-test.readthedocs.io/en/latest/)
- [Real Python: Python Design Patterns](https://realpython.com/tutorials/design-patterns/)