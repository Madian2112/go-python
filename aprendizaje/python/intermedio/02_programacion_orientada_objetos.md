# Programación Orientada a Objetos Avanzada en Python

## Introducción

Python es un lenguaje multiparadigma que soporta completamente la programación orientada a objetos (POO). Aunque los conceptos básicos de POO como clases y objetos son fundamentales, Python ofrece características avanzadas que permiten aprovechar al máximo este paradigma. En este módulo, exploraremos conceptos avanzados de POO en Python, incluyendo herencia múltiple, métodos especiales, metaclases, y patrones de diseño orientados a objetos.

## Repaso de Conceptos Básicos

### Clases y Objetos

```python
class Persona:
    def __init__(self, nombre, edad):
        self.nombre = nombre
        self.edad = edad
    
    def saludar(self):
        return f"Hola, soy {self.nombre} y tengo {self.edad} años"

# Crear una instancia
persona1 = Persona("Ana", 30)
print(persona1.saludar())  # Hola, soy Ana y tengo 30 años
```

### Herencia Simple

```python
class Empleado(Persona):
    def __init__(self, nombre, edad, salario):
        super().__init__(nombre, edad)
        self.salario = salario
    
    def info_laboral(self):
        return f"{self.nombre} gana {self.salario} euros al mes"

empleado1 = Empleado("Carlos", 35, 2500)
print(empleado1.saludar())      # Método heredado
print(empleado1.info_laboral())  # Método propio
```

## Herencia Avanzada

### Herencia Múltiple

Python permite que una clase herede de múltiples clases base, lo que puede ser muy poderoso pero también complejo.

```python
class Vehiculo:
    def __init__(self, marca, modelo):
        self.marca = marca
        self.modelo = modelo
    
    def info(self):
        return f"Vehículo: {self.marca} {self.modelo}"

class Electrico:
    def __init__(self, autonomia):
        self.autonomia = autonomia
    
    def info_electrico(self):
        return f"Autonomía: {self.autonomia} km"

class CocheElectrico(Vehiculo, Electrico):
    def __init__(self, marca, modelo, autonomia):
        Vehiculo.__init__(self, marca, modelo)
        Electrico.__init__(self, autonomia)
    
    def info_completa(self):
        return f"{self.info()} - {self.info_electrico()}"

tesla = CocheElectrico("Tesla", "Model 3", 500)
print(tesla.info_completa())  # Vehículo: Tesla Model 3 - Autonomía: 500 km
```

### Orden de Resolución de Métodos (MRO)

Cuando se utiliza herencia múltiple, Python sigue un orden específico para buscar métodos y atributos, conocido como MRO (Method Resolution Order).

```python
class A:
    def quien_soy(self):
        return "Soy A"

class B(A):
    def quien_soy(self):
        return "Soy B"

class C(A):
    def quien_soy(self):
        return "Soy C"

class D(B, C):
    pass

d = D()
print(d.quien_soy())  # Soy B
print(D.__mro__)      # Muestra el orden de resolución de métodos
```

### Mixins

Los mixins son clases diseñadas para añadir funcionalidades específicas a otras clases a través de la herencia múltiple.

```python
class SerializableMixin:
    def to_dict(self):
        return {key: value for key, value in self.__dict__.items() 
                if not key.startswith('_')}
    
    def to_json(self):
        import json
        return json.dumps(self.to_dict())

class Producto(SerializableMixin):
    def __init__(self, nombre, precio):
        self.nombre = nombre
        self.precio = precio

producto = Producto("Laptop", 1200)
print(producto.to_json())  # {"nombre": "Laptop", "precio": 1200}
```

## Métodos Especiales (Dunder Methods)

Python utiliza métodos especiales (también llamados métodos dunder o métodos mágicos) para implementar comportamientos específicos en las clases.

### Representación de Objetos

```python
class Punto:
    def __init__(self, x, y):
        self.x = x
        self.y = y
    
    def __str__(self):
        return f"Punto({self.x}, {self.y})"
    
    def __repr__(self):
        return f"Punto(x={self.x}, y={self.y})"

punto = Punto(10, 20)
print(punto)        # Punto(10, 20) - usa __str__
print(repr(punto))  # Punto(x=10, y=20) - usa __repr__
```

### Operadores Aritméticos

```python
class Vector:
    def __init__(self, x, y):
        self.x = x
        self.y = y
    
    def __add__(self, otro):
        return Vector(self.x + otro.x, self.y + otro.y)
    
    def __sub__(self, otro):
        return Vector(self.x - otro.x, self.y - otro.y)
    
    def __mul__(self, escalar):
        return Vector(self.x * escalar, self.y * escalar)
    
    def __str__(self):
        return f"Vector({self.x}, {self.y})"

v1 = Vector(2, 3)
v2 = Vector(5, 1)
print(v1 + v2)    # Vector(7, 4)
print(v1 - v2)    # Vector(-3, 2)
print(v1 * 2)     # Vector(4, 6)
```

### Comparación

```python
class Persona:
    def __init__(self, nombre, edad):
        self.nombre = nombre
        self.edad = edad
    
    def __eq__(self, otro):
        if not isinstance(otro, Persona):
            return NotImplemented
        return self.nombre == otro.nombre and self.edad == otro.edad
    
    def __lt__(self, otro):
        if not isinstance(otro, Persona):
            return NotImplemented
        return self.edad < otro.edad

p1 = Persona("Ana", 30)
p2 = Persona("Ana", 30)
p3 = Persona("Carlos", 25)

print(p1 == p2)  # True - mismos atributos
print(p1 == p3)  # False - diferentes atributos
print(p3 < p1)   # True - Carlos es más joven que Ana
```

### Contenedores

```python
class Carrito:
    def __init__(self):
        self.productos = []
    
    def agregar(self, producto):
        self.productos.append(producto)
    
    def __len__(self):
        return len(self.productos)
    
    def __getitem__(self, indice):
        return self.productos[indice]
    
    def __iter__(self):
        return iter(self.productos)

carrito = Carrito()
carrito.agregar("Laptop")
carrito.agregar("Mouse")
carrito.agregar("Teclado")

print(len(carrito))     # 3
print(carrito[1])       # Mouse

for producto in carrito:
    print(producto)     # Imprime cada producto
```

### Context Managers

```python
class Archivo:
    def __init__(self, nombre, modo):
        self.nombre = nombre
        self.modo = modo
    
    def __enter__(self):
        self.archivo = open(self.nombre, self.modo)
        return self.archivo
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.archivo.close()

# Usar la clase como context manager
with Archivo("ejemplo.txt", "w") as f:
    f.write("Hola, mundo!")
# El archivo se cierra automáticamente al salir del bloque with
```

## Propiedades y Descriptores

### Propiedades

Las propiedades permiten controlar el acceso a los atributos de una clase.

```python
class Temperatura:
    def __init__(self):
        self._celsius = 0
    
    @property
    def celsius(self):
        return self._celsius
    
    @celsius.setter
    def celsius(self, valor):
        if valor < -273.15:
            raise ValueError("La temperatura no puede ser menor que -273.15°C")
        self._celsius = valor
    
    @property
    def fahrenheit(self):
        return self._celsius * 9/5 + 32
    
    @fahrenheit.setter
    def fahrenheit(self, valor):
        self.celsius = (valor - 32) * 5/9

temp = Temperatura()
temp.celsius = 25
print(temp.celsius)     # 25
print(temp.fahrenheit)  # 77.0

temp.fahrenheit = 68
print(temp.celsius)     # 20.0
```

### Descriptores

Los descriptores son objetos que definen cómo se accede a los atributos de otras clases.

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
            raise ValueError(f"{self.nombre} no puede ser menor que {self.min_valor}")
        if self.max_valor is not None and valor > self.max_valor:
            raise ValueError(f"{self.nombre} no puede ser mayor que {self.max_valor}")
        instancia.__dict__[self.nombre] = valor

class Persona:
    edad = Validador(min_valor=0, max_valor=120)
    altura = Validador(min_valor=0, max_valor=250)
    
    def __init__(self, nombre, edad, altura):
        self.nombre = nombre
        self.edad = edad
        self.altura = altura

try:
    p = Persona("Ana", 30, 165)
    print(f"{p.nombre}: {p.edad} años, {p.altura} cm")
    
    p.edad = 130  # Esto lanzará una excepción
except ValueError as e:
    print(f"Error: {e}")
```

## Metaclases

Las metaclases son clases de clases, es decir, clases cuyas instancias son clases.

```python
class Meta(type):
    def __new__(mcs, nombre, bases, attrs):
        # Convertir todos los métodos a mayúsculas
        attrs_mayusculas = {}
        for key, value in attrs.items():
            if key.startswith('__') and key.endswith('__'):
                attrs_mayusculas[key] = value
            else:
                attrs_mayusculas[key.upper()] = value
        
        return super().__new__(mcs, nombre, bases, attrs_mayusculas)

class MiClase(metaclass=Meta):
    def metodo(self):
        return "Este método estará en mayúsculas"

obj = MiClase()
print(obj.METODO())  # Este método estará en mayúsculas
```

### Registro Automático de Clases

```python
class RegistroMeta(type):
    clases = {}
    
    def __new__(mcs, nombre, bases, attrs):
        clase = super().__new__(mcs, nombre, bases, attrs)
        if nombre != 'ModeloBase':
            mcs.clases[nombre] = clase
        return clase

class ModeloBase(metaclass=RegistroMeta):
    pass

class Usuario(ModeloBase):
    pass

class Producto(ModeloBase):
    pass

print(RegistroMeta.clases)  # {'Usuario': <class '__main__.Usuario'>, 'Producto': <class '__main__.Producto'>}
```

## Patrones de Diseño Orientados a Objetos

### Singleton

El patrón Singleton garantiza que una clase tenga una única instancia y proporciona un punto de acceso global a ella.

```python
class Singleton(type):
    _instances = {}
    
    def __call__(cls, *args, **kwargs):
        if cls not in cls._instances:
            cls._instances[cls] = super().__call__(*args, **kwargs)
        return cls._instances[cls]

class ConfiguracionApp(metaclass=Singleton):
    def __init__(self):
        self.tema = "claro"
        self.idioma = "es"

# Ambas variables apuntan a la misma instancia
config1 = ConfiguracionApp()
config2 = ConfiguracionApp()

config1.tema = "oscuro"
print(config2.tema)  # oscuro - los cambios en config1 afectan a config2
```

### Factory Method

El patrón Factory Method define una interfaz para crear objetos, pero permite a las subclases decidir qué clase instanciar.

```python
from abc import ABC, abstractmethod

class Animal(ABC):
    @abstractmethod
    def hablar(self):
        pass

class Perro(Animal):
    def hablar(self):
        return "Guau!"

class Gato(Animal):
    def hablar(self):
        return "Miau!"

class FabricaAnimales:
    def crear_animal(self, tipo):
        if tipo == "perro":
            return Perro()
        elif tipo == "gato":
            return Gato()
        else:
            raise ValueError(f"Tipo de animal desconocido: {tipo}")

fabrica = FabricaAnimales()
animales = [fabrica.crear_animal("perro"), fabrica.crear_animal("gato")]

for animal in animales:
    print(animal.hablar())
```

### Observer

El patrón Observer define una dependencia uno a muchos entre objetos, de modo que cuando un objeto cambia de estado, todos sus dependientes son notificados y actualizados automáticamente.

```python
class Sujeto:
    def __init__(self):
        self._observadores = []
    
    def registrar(self, observador):
        if observador not in self._observadores:
            self._observadores.append(observador)
    
    def desregistrar(self, observador):
        if observador in self._observadores:
            self._observadores.remove(observador)
    
    def notificar(self, *args, **kwargs):
        for observador in self._observadores:
            observador.actualizar(self, *args, **kwargs)

class Observador(ABC):
    @abstractmethod
    def actualizar(self, sujeto, *args, **kwargs):
        pass

class EstacionMeteorologica(Sujeto):
    def __init__(self):
        super().__init__()
        self._temperatura = 0
    
    @property
    def temperatura(self):
        return self._temperatura
    
    @temperatura.setter
    def temperatura(self, valor):
        self._temperatura = valor
        self.notificar()

class PantallaTemperatura(Observador):
    def actualizar(self, sujeto, *args, **kwargs):
        print(f"La temperatura actual es: {sujeto.temperatura}°C")

class AlertaTemperatura(Observador):
    def actualizar(self, sujeto, *args, **kwargs):
        if sujeto.temperatura > 30:
            print("¡Alerta! Temperatura alta")

estacion = EstacionMeteorologica()
pantalla = PantallaTemperatura()
alerta = AlertaTemperatura()

estacion.registrar(pantalla)
estacion.registrar(alerta)

estacion.temperatura = 25  # La temperatura actual es: 25°C
estacion.temperatura = 32  # La temperatura actual es: 32°C ¡Alerta! Temperatura alta
```

## Composición vs Herencia

La composición es a menudo preferible a la herencia para crear relaciones entre clases.

```python
# Enfoque de herencia
class Animal:
    def __init__(self, nombre):
        self.nombre = nombre
    
    def comer(self):
        return f"{self.nombre} está comiendo"

class Volador:
    def volar(self):
        return f"{self.nombre} está volando"

class Ave(Animal, Volador):
    def __init__(self, nombre):
        super().__init__(nombre)

# Enfoque de composición
class Comportamiento:
    def __init__(self, nombre):
        self.nombre = nombre

class ComportamientoComer(Comportamiento):
    def ejecutar(self):
        return f"{self.nombre} está comiendo"

class ComportamientoVolar(Comportamiento):
    def ejecutar(self):
        return f"{self.nombre} está volando"

class Animal:
    def __init__(self, nombre):
        self.nombre = nombre
        self.comportamientos = []
    
    def agregar_comportamiento(self, comportamiento):
        self.comportamientos.append(comportamiento)
    
    def ejecutar_comportamientos(self):
        return [comportamiento.ejecutar() for comportamiento in self.comportamientos]

# Crear un ave con composición
ave = Animal("Águila")
ave.agregar_comportamiento(ComportamientoComer(ave.nombre))
ave.agregar_comportamiento(ComportamientoVolar(ave.nombre))

print(ave.ejecutar_comportamientos())  # ['Águila está comiendo', 'Águila está volando']
```

## Mejores Prácticas

1. **Prefiere composición sobre herencia**: La composición es más flexible y menos propensa a problemas.

2. **Sigue el principio SOLID**:
   - **S**ingle Responsibility: Una clase debe tener una sola razón para cambiar.
   - **O**pen/Closed: Las clases deben estar abiertas para extensión pero cerradas para modificación.
   - **L**iskov Substitution: Las subclases deben poder sustituir a sus clases base.
   - **I**nterface Segregation: Muchas interfaces específicas son mejores que una interfaz general.
   - **D**ependency Inversion: Depende de abstracciones, no de implementaciones concretas.

3. **Usa propiedades en lugar de getters y setters**: Las propiedades proporcionan una sintaxis más limpia.

4. **Implementa métodos especiales cuando sea apropiado**: Mejoran la integración con el lenguaje.

5. **Documenta tus clases y métodos**: Usa docstrings para explicar el propósito y uso de tus clases.

6. **Evita la herencia múltiple profunda**: Puede llevar a problemas difíciles de depurar.

7. **Usa mixins para comportamientos reutilizables**: Los mixins son ideales para añadir funcionalidades específicas.

8. **Considera el uso de dataclasses para datos simples**: A partir de Python 3.7, las dataclasses simplifican la creación de clases centradas en datos.

## Ejercicios Prácticos

1. Implementa una jerarquía de clases para un sistema de gestión de empleados, utilizando herencia y polimorfismo.

2. Crea una clase que implemente varios métodos especiales para que se comporte como un contenedor personalizado.

3. Diseña un sistema utilizando el patrón Observer para notificar a múltiples componentes cuando cambia el estado de un objeto.

4. Implementa una metaclase que valide automáticamente los tipos de atributos de las clases que la utilizan.

5. Refactoriza un diseño basado en herencia para utilizar composición, siguiendo el principio "prefiere composición sobre herencia".

## Conclusión

La programación orientada a objetos en Python ofrece herramientas poderosas para crear código modular, reutilizable y mantenible. Dominar conceptos avanzados como herencia múltiple, métodos especiales, propiedades, descriptores y metaclases te permitirá aprovechar al máximo las capacidades de Python para resolver problemas complejos de manera elegante.

Recuerda que, aunque Python proporciona muchas características avanzadas de POO, la simplicidad sigue siendo un valor fundamental. Usa estas herramientas cuando aporten claridad y valor a tu código, no solo porque estén disponibles.