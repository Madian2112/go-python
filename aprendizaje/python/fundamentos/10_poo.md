# Programación Orientada a Objetos en Python

## Introducción

La Programación Orientada a Objetos (POO) es un paradigma de programación que utiliza "objetos" para modelar datos y comportamientos. Python es un lenguaje multiparadigma que soporta completamente la POO, permitiendo crear programas modulares, reutilizables y fáciles de mantener. En esta sección, exploraremos los conceptos fundamentales de la POO en Python, incluyendo clases, objetos, herencia, polimorfismo, encapsulación y más.

## Clases y Objetos

### Definición de Clases

Una clase es un plano o plantilla para crear objetos. Define atributos (datos) y métodos (funciones) que serán compartidos por todos los objetos de ese tipo.

```python
class Persona:
    """Clase que representa a una persona."""
    
    def __init__(self, nombre, edad):
        """Constructor de la clase.
        
        Args:
            nombre (str): Nombre de la persona
            edad (int): Edad de la persona
        """
        self.nombre = nombre
        self.edad = edad
    
    def saludar(self):
        """Método que imprime un saludo."""
        print(f"Hola, me llamo {self.nombre} y tengo {self.edad} años.")
```

### Creación de Objetos (Instancias)

Un objeto es una instancia de una clase. Para crear un objeto, llamamos a la clase como si fuera una función.

```python
# Crear objetos de la clase Persona
persona1 = Persona("Ana", 30)
persona2 = Persona("Juan", 25)

# Usar los métodos de los objetos
persona1.saludar()  # Salida: Hola, me llamo Ana y tengo 30 años.
persona2.saludar()  # Salida: Hola, me llamo Juan y tengo 25 años.
```

### Atributos de Instancia vs. Atributos de Clase

Los atributos de instancia son únicos para cada objeto, mientras que los atributos de clase son compartidos por todas las instancias de la clase.

```python
class Persona:
    # Atributo de clase
    especie = "Homo sapiens"
    
    def __init__(self, nombre, edad):
        # Atributos de instancia
        self.nombre = nombre
        self.edad = edad

# Crear objetos
persona1 = Persona("Ana", 30)
persona2 = Persona("Juan", 25)

# Acceder a atributos de instancia
print(persona1.nombre)  # Ana
print(persona2.nombre)  # Juan

# Acceder a atributos de clase
print(persona1.especie)  # Homo sapiens
print(persona2.especie)  # Homo sapiens
print(Persona.especie)   # Homo sapiens

# Modificar un atributo de clase
Persona.especie = "Homo sapiens sapiens"
print(persona1.especie)  # Homo sapiens sapiens
print(persona2.especie)  # Homo sapiens sapiens
```

### Métodos de Instancia, Métodos de Clase y Métodos Estáticos

Python permite definir tres tipos de métodos en una clase:

```python
class Persona:
    # Atributo de clase
    contador = 0
    
    def __init__(self, nombre, edad):
        self.nombre = nombre
        self.edad = edad
        Persona.contador += 1
    
    # Método de instancia
    def saludar(self):
        return f"Hola, me llamo {self.nombre} y tengo {self.edad} años."
    
    # Método de clase
    @classmethod
    def desde_año_nacimiento(cls, nombre, año_nacimiento):
        """Crea una instancia de Persona a partir del año de nacimiento."""
        import datetime
        edad = datetime.datetime.now().year - año_nacimiento
        return cls(nombre, edad)
    
    # Método de clase para obtener el contador
    @classmethod
    def obtener_contador(cls):
        return cls.contador
    
    # Método estático
    @staticmethod
    def es_adulto(edad):
        """Verifica si una edad corresponde a un adulto."""
        return edad >= 18

# Usar método de instancia
persona1 = Persona("Ana", 30)
print(persona1.saludar())  # Hola, me llamo Ana y tengo 30 años.

# Usar método de clase
persona2 = Persona.desde_año_nacimiento("Juan", 1995)
print(persona2.saludar())  # Hola, me llamo Juan y tengo XX años. (depende del año actual)

# Usar método de clase para obtener el contador
print(Persona.obtener_contador())  # 2

# Usar método estático
print(Persona.es_adulto(20))  # True
print(Persona.es_adulto(15))  # False
```

## Herencia

La herencia permite que una clase (subclase) herede atributos y métodos de otra clase (superclase). Esto promueve la reutilización de código y la creación de jerarquías de clases.

### Herencia Simple

```python
class Persona:
    def __init__(self, nombre, edad):
        self.nombre = nombre
        self.edad = edad
    
    def saludar(self):
        return f"Hola, me llamo {self.nombre} y tengo {self.edad} años."

class Estudiante(Persona):
    def __init__(self, nombre, edad, carrera):
        # Llamar al constructor de la clase padre
        super().__init__(nombre, edad)
        self.carrera = carrera
    
    def estudiar(self):
        return f"Estoy estudiando {self.carrera}."

# Crear un estudiante
estudiante = Estudiante("María", 22, "Ingeniería")

# Usar métodos heredados y propios
print(estudiante.saludar())  # Método heredado
print(estudiante.estudiar())  # Método propio
```

### Herencia Múltiple

Python permite que una clase herede de múltiples clases, lo que se conoce como herencia múltiple.

```python
class Persona:
    def __init__(self, nombre, edad):
        self.nombre = nombre
        self.edad = edad
    
    def saludar(self):
        return f"Hola, me llamo {self.nombre} y tengo {self.edad} años."

class Empleado:
    def __init__(self, salario, puesto):
        self.salario = salario
        self.puesto = puesto
    
    def trabajar(self):
        return f"Estoy trabajando como {self.puesto}."

class Profesor(Persona, Empleado):
    def __init__(self, nombre, edad, salario, puesto, materia):
        Persona.__init__(self, nombre, edad)
        Empleado.__init__(self, salario, puesto)
        self.materia = materia
    
    def enseñar(self):
        return f"Estoy enseñando {self.materia}."

# Crear un profesor
profesor = Profesor("Carlos", 45, 50000, "Profesor", "Matemáticas")

# Usar métodos de todas las clases
print(profesor.saludar())   # De Persona
print(profesor.trabajar())  # De Empleado
print(profesor.enseñar())   # De Profesor
```

### Orden de Resolución de Métodos (MRO)

Cuando se utiliza herencia múltiple, Python sigue un orden específico para buscar métodos y atributos, conocido como MRO (Method Resolution Order).

```python
class A:
    def metodo(self):
        return "Método de A"

class B(A):
    def metodo(self):
        return "Método de B"

class C(A):
    def metodo(self):
        return "Método de C"

class D(B, C):
    pass

# Crear una instancia de D
d = D()

# Ver el MRO
print(D.__mro__)  # (<class '__main__.D'>, <class '__main__.B'>, <class '__main__.C'>, <class '__main__.A'>, <class 'object'>)

# Llamar al método
print(d.metodo())  # Método de B (porque B está antes que C en el MRO)
```

## Polimorfismo

El polimorfismo permite que objetos de diferentes clases respondan al mismo método o atributo de manera diferente.

```python
class Animal:
    def sonido(self):
        pass

class Perro(Animal):
    def sonido(self):
        return "Guau"

class Gato(Animal):
    def sonido(self):
        return "Miau"

class Vaca(Animal):
    def sonido(self):
        return "Muu"

# Función que utiliza polimorfismo
def hacer_sonido(animal):
    return animal.sonido()

# Crear animales
perro = Perro()
gato = Gato()
vaca = Vaca()

# Usar polimorfismo
print(hacer_sonido(perro))  # Guau
print(hacer_sonido(gato))   # Miau
print(hacer_sonido(vaca))   # Muu

# También funciona con listas
animales = [Perro(), Gato(), Vaca()]
for animal in animales:
    print(animal.sonido())
```

## Encapsulación

La encapsulación es el principio de ocultar los detalles internos de una clase y proporcionar una interfaz controlada para interactuar con ella.

### Atributos Privados y Protegidos

En Python, la convención es usar un guion bajo (`_`) para indicar atributos protegidos y dos guiones bajos (`__`) para atributos privados.

```python
class Cuenta:
    def __init__(self, titular, saldo):
        self.titular = titular      # Público
        self._saldo = saldo         # Protegido (convención)
        self.__historial = []       # Privado (name mangling)
    
    def depositar(self, cantidad):
        if cantidad > 0:
            self._saldo += cantidad
            self.__registrar(f"Depósito de {cantidad}")
            return True
        return False
    
    def retirar(self, cantidad):
        if 0 < cantidad <= self._saldo:
            self._saldo -= cantidad
            self.__registrar(f"Retiro de {cantidad}")
            return True
        return False
    
    def consultar_saldo(self):
        return self._saldo
    
    def __registrar(self, operacion):
        self.__historial.append(operacion)
    
    def consultar_historial(self):
        return self.__historial.copy()

# Crear una cuenta
cuenta = Cuenta("Juan", 1000)

# Usar métodos públicos
print(cuenta.titular)           # Juan
print(cuenta.consultar_saldo()) # 1000
cuenta.depositar(500)
print(cuenta.consultar_saldo()) # 1500
cuenta.retirar(200)
print(cuenta.consultar_saldo()) # 1300
print(cuenta.consultar_historial()) # ['Depósito de 500', 'Retiro de 200']

# Acceder a atributos protegidos (posible, pero no recomendado)
print(cuenta._saldo)  # 1300

# Intentar acceder a atributos privados
try:
    print(cuenta.__historial)  # Esto generará un error
except AttributeError as e:
    print(f"Error: {e}")

# Pero en realidad, el atributo privado está disponible con un nombre diferente
print(cuenta._Cuenta__historial)  # ['Depósito de 500', 'Retiro de 200']
```

### Propiedades

Las propiedades permiten controlar el acceso a los atributos de una clase, proporcionando getters, setters y deleters.

```python
class Temperatura:
    def __init__(self):
        self._celsius = 0
    
    @property
    def celsius(self):
        """Obtener la temperatura en Celsius."""
        return self._celsius
    
    @celsius.setter
    def celsius(self, valor):
        """Establecer la temperatura en Celsius."""
        if valor < -273.15:
            raise ValueError("La temperatura no puede ser menor que el cero absoluto.")
        self._celsius = valor
    
    @property
    def fahrenheit(self):
        """Obtener la temperatura en Fahrenheit."""
        return self._celsius * 9/5 + 32
    
    @fahrenheit.setter
    def fahrenheit(self, valor):
        """Establecer la temperatura en Fahrenheit."""
        self.celsius = (valor - 32) * 5/9

# Crear un objeto temperatura
temp = Temperatura()

# Usar propiedades
temp.celsius = 25
print(f"{temp.celsius}°C = {temp.fahrenheit}°F")  # 25°C = 77.0°F

temp.fahrenheit = 68
print(f"{temp.celsius}°C = {temp.fahrenheit}°F")  # 20.0°C = 68.0°F

# Intentar establecer una temperatura inválida
try:
    temp.celsius = -300
except ValueError as e:
    print(f"Error: {e}")
```

## Métodos Especiales (Dunder Methods)

Python utiliza métodos especiales, también conocidos como "dunder methods" (double underscore methods), para definir cómo se comportan los objetos en diferentes situaciones.

```python
class Vector:
    def __init__(self, x, y):
        self.x = x
        self.y = y
    
    def __str__(self):
        """Representación en string para usuarios."""
        return f"Vector({self.x}, {self.y})"
    
    def __repr__(self):
        """Representación en string para desarrolladores."""
        return f"Vector({self.x}, {self.y})"
    
    def __add__(self, otro):
        """Suma de vectores usando el operador +."""
        return Vector(self.x + otro.x, self.y + otro.y)
    
    def __sub__(self, otro):
        """Resta de vectores usando el operador -."""
        return Vector(self.x - otro.x, self.y - otro.y)
    
    def __mul__(self, escalar):
        """Multiplicación por escalar usando el operador *."""
        return Vector(self.x * escalar, self.y * escalar)
    
    def __rmul__(self, escalar):
        """Multiplicación por escalar (orden inverso)."""
        return self.__mul__(escalar)
    
    def __eq__(self, otro):
        """Comparación de igualdad usando el operador ==."""
        if not isinstance(otro, Vector):
            return False
        return self.x == otro.x and self.y == otro.y
    
    def __len__(self):
        """Longitud del vector (redondeada)."""
        import math
        return round(math.sqrt(self.x**2 + self.y**2))
    
    def __getitem__(self, indice):
        """Acceso a componentes usando índices."""
        if indice == 0:
            return self.x
        elif indice == 1:
            return self.y
        else:
            raise IndexError("Índice fuera de rango")

# Crear vectores
v1 = Vector(3, 4)
v2 = Vector(1, 2)

# Usar métodos especiales
print(v1)          # Vector(3, 4) (usa __str__)
print(repr(v2))    # Vector(1, 2) (usa __repr__)

v3 = v1 + v2       # Vector(4, 6) (usa __add__)
print(v3)

v4 = v1 - v2       # Vector(2, 2) (usa __sub__)
print(v4)

v5 = v1 * 2        # Vector(6, 8) (usa __mul__)
print(v5)

v6 = 3 * v2        # Vector(3, 6) (usa __rmul__)
print(v6)

print(v1 == v2)    # False (usa __eq__)
print(v1 == Vector(3, 4))  # True

print(len(v1))     # 5 (usa __len__, sqrt(3^2 + 4^2) = 5)

print(v1[0])       # 3 (usa __getitem__)
print(v1[1])       # 4

try:
    print(v1[2])   # Esto generará un error
except IndexError as e:
    print(f"Error: {e}")
```

### Otros Métodos Especiales Comunes

```python
class MiLista:
    def __init__(self, elementos=None):
        self.elementos = elementos or []
    
    def __str__(self):
        return str(self.elementos)
    
    def __len__(self):
        return len(self.elementos)
    
    def __getitem__(self, indice):
        return self.elementos[indice]
    
    def __setitem__(self, indice, valor):
        self.elementos[indice] = valor
    
    def __delitem__(self, indice):
        del self.elementos[indice]
    
    def __iter__(self):
        return iter(self.elementos)
    
    def __contains__(self, item):
        return item in self.elementos
    
    def __bool__(self):
        return bool(self.elementos)

# Crear una lista personalizada
mi_lista = MiLista([1, 2, 3, 4, 5])

# Usar métodos especiales
print(mi_lista)        # [1, 2, 3, 4, 5] (usa __str__)
print(len(mi_lista))   # 5 (usa __len__)
print(mi_lista[2])     # 3 (usa __getitem__)

mi_lista[1] = 10       # (usa __setitem__)
print(mi_lista)        # [1, 10, 3, 4, 5]

del mi_lista[0]        # (usa __delitem__)
print(mi_lista)        # [10, 3, 4, 5]

for item in mi_lista:  # (usa __iter__)
    print(item, end=" ")  # 10 3 4 5
print()

print(3 in mi_lista)   # True (usa __contains__)
print(6 in mi_lista)   # False

if mi_lista:           # (usa __bool__)
    print("La lista no está vacía")
```

## Clases Abstractas

Las clases abstractas definen una interfaz que las subclases deben implementar. En Python, se utilizan a través del módulo `abc` (Abstract Base Classes).

```python
from abc import ABC, abstractmethod

class FiguraGeometrica(ABC):
    @abstractmethod
    def area(self):
        """Calcular el área de la figura."""
        pass
    
    @abstractmethod
    def perimetro(self):
        """Calcular el perímetro de la figura."""
        pass
    
    def descripcion(self):
        """Método no abstracto que pueden usar las subclases."""
        return f"Soy una figura geométrica con área {self.area()} y perímetro {self.perimetro()}."

class Rectangulo(FiguraGeometrica):
    def __init__(self, base, altura):
        self.base = base
        self.altura = altura
    
    def area(self):
        return self.base * self.altura
    
    def perimetro(self):
        return 2 * (self.base + self.altura)

class Circulo(FiguraGeometrica):
    def __init__(self, radio):
        self.radio = radio
    
    def area(self):
        import math
        return math.pi * self.radio ** 2
    
    def perimetro(self):
        import math
        return 2 * math.pi * self.radio

# Intentar crear una instancia de la clase abstracta
try:
    figura = FiguraGeometrica()  # Esto generará un error
except TypeError as e:
    print(f"Error: {e}")

# Crear instancias de las subclases
rectangulo = Rectangulo(5, 3)
circulo = Circulo(4)

# Usar métodos
print(f"Área del rectángulo: {rectangulo.area()}")        # 15
print(f"Perímetro del rectángulo: {rectangulo.perimetro()}")  # 16
print(rectangulo.descripcion())

print(f"Área del círculo: {circulo.area():.2f}")          # ~50.27
print(f"Perímetro del círculo: {circulo.perimetro():.2f}")    # ~25.13
print(circulo.descripcion())
```

## Mixins

Los mixins son clases que proporcionan funcionalidades adicionales a otras clases a través de la herencia múltiple.

```python
class SerializableMixin:
    def a_json(self):
        import json
        return json.dumps(self.__dict__)
    
    def desde_json(self, json_str):
        import json
        datos = json.loads(json_str)
        for clave, valor in datos.items():
            setattr(self, clave, valor)
        return self

class LoggableMixin:
    def log(self, mensaje):
        print(f"[LOG] {mensaje}")
    
    def log_metodo(self, nombre_metodo):
        def decorador(metodo):
            def wrapper(*args, **kwargs):
                self.log(f"Llamando al método {nombre_metodo}")
                resultado = metodo(*args, **kwargs)
                self.log(f"Método {nombre_metodo} completado")
                return resultado
            return wrapper
        return decorador

class Producto(SerializableMixin, LoggableMixin):
    def __init__(self, nombre, precio):
        self.nombre = nombre
        self.precio = precio
    
    @property
    def precio_con_iva(self):
        return self.precio * 1.21

# Crear un producto
producto = Producto("Laptop", 1000)

# Usar funcionalidades de SerializableMixin
json_str = producto.a_json()
print(json_str)  # {"nombre": "Laptop", "precio": 1000}

producto_nuevo = Producto("", 0).desde_json(json_str)
print(f"{producto_nuevo.nombre}: {producto_nuevo.precio}€")  # Laptop: 1000€

# Usar funcionalidades de LoggableMixin
producto.log(f"Precio con IVA: {producto.precio_con_iva}€")  # [LOG] Precio con IVA: 1210.0€

# Decorar un método con log_metodo
@producto.log_metodo("aplicar_descuento")
def aplicar_descuento(self, porcentaje):
    descuento = self.precio * (porcentaje / 100)
    self.precio -= descuento
    return self.precio

# Añadir el método al producto
Producto.aplicar_descuento = aplicar_descuento

# Usar el método decorado
producto.aplicar_descuento(10)  # [LOG] Llamando al método aplicar_descuento
                               # [LOG] Método aplicar_descuento completado
print(f"Precio con descuento: {producto.precio}€")  # Precio con descuento: 900.0€
```

## Descriptores

Los descriptores son objetos que definen cómo se comportan los atributos de una clase cuando se accede a ellos, se modifican o se eliminan.

```python
class Validador:
    def __init__(self, nombre, validacion, mensaje_error):
        self.nombre = nombre
        self.validacion = validacion
        self.mensaje_error = mensaje_error
    
    def __set_name__(self, propietario, nombre):
        self.nombre_privado = f"_{nombre}"
    
    def __get__(self, instancia, propietario):
        if instancia is None:
            return self
        return getattr(instancia, self.nombre_privado, None)
    
    def __set__(self, instancia, valor):
        if not self.validacion(valor):
            raise ValueError(f"{self.nombre}: {self.mensaje_error}")
        setattr(instancia, self.nombre_privado, valor)

class Persona:
    nombre = Validador(
        "Nombre",
        lambda x: isinstance(x, str) and 2 <= len(x) <= 30,
        "debe ser una cadena de entre 2 y 30 caracteres"
    )
    
    edad = Validador(
        "Edad",
        lambda x: isinstance(x, int) and 0 <= x <= 120,
        "debe ser un entero entre 0 y 120"
    )
    
    email = Validador(
        "Email",
        lambda x: isinstance(x, str) and "@" in x,
        "debe ser una cadena con formato de email válido"
    )
    
    def __init__(self, nombre, edad, email):
        self.nombre = nombre
        self.edad = edad
        self.email = email
    
    def __str__(self):
        return f"{self.nombre}, {self.edad} años, {self.email}"

# Crear una persona válida
persona = Persona("Juan", 30, "juan@example.com")
print(persona)  # Juan, 30 años, juan@example.com

# Intentar crear una persona con datos inválidos
try:
    persona_invalida = Persona("J", 30, "juan@example.com")  # Nombre muy corto
except ValueError as e:
    print(f"Error: {e}")

try:
    persona.edad = 150  # Edad fuera de rango
except ValueError as e:
    print(f"Error: {e}")

try:
    persona.email = "juanexample.com"  # Email sin @
except ValueError as e:
    print(f"Error: {e}")
```

## Metaclases

Las metaclases son clases que definen el comportamiento de otras clases. Son un concepto avanzado en Python.

```python
class Meta(type):
    def __new__(mcs, nombre, bases, attrs):
        # Añadir un prefijo a todos los métodos que no sean especiales
        for clave, valor in list(attrs.items()):
            if callable(valor) and not clave.startswith("__"):
                attrs[f"metodo_{clave}"] = valor
                del attrs[clave]
        
        # Añadir un método a todas las clases
        attrs["hablar"] = lambda self: f"{nombre} dice: Hola"
        
        # Crear la clase
        return super().__new__(mcs, nombre, bases, attrs)

class MiClase(metaclass=Meta):
    def saludar(self):
        return "¡Hola!"
    
    def despedir(self):
        return "¡Adiós!"

# Crear una instancia
obj = MiClase()

# Los métodos originales ya no existen
try:
    print(obj.saludar())
except AttributeError as e:
    print(f"Error: {e}")

# Pero existen con el prefijo
print(obj.metodo_saludar())  # ¡Hola!
print(obj.metodo_despedir())  # ¡Adiós!

# Y también existe el método añadido por la metaclase
print(obj.hablar())  # MiClase dice: Hola
```

## Ejemplo Práctico: Sistema de Gestión de Biblioteca

Vamos a crear un sistema de gestión de biblioteca que utilice los conceptos de POO que hemos aprendido.

```python
from abc import ABC, abstractmethod
from datetime import datetime, timedelta

# Clase base abstracta para elementos de la biblioteca
class ElementoBiblioteca(ABC):
    def __init__(self, codigo, titulo, año):
        self.codigo = codigo
        self.titulo = titulo
        self.año = año
        self.disponible = True
    
    @abstractmethod
    def tipo(self):
        pass
    
    def __str__(self):
        estado = "disponible" if self.disponible else "no disponible"
        return f"{self.tipo()} - {self.titulo} ({self.año}) - {estado}"

# Clases concretas para diferentes tipos de elementos
class Libro(ElementoBiblioteca):
    def __init__(self, codigo, titulo, año, autor, genero, paginas):
        super().__init__(codigo, titulo, año)
        self.autor = autor
        self.genero = genero
        self.paginas = paginas
    
    def tipo(self):
        return "Libro"

class Revista(ElementoBiblioteca):
    def __init__(self, codigo, titulo, año, editorial, periodicidad, numero):
        super().__init__(codigo, titulo, año)
        self.editorial = editorial
        self.periodicidad = periodicidad
        self.numero = numero
    
    def tipo(self):
        return "Revista"

class DVD(ElementoBiblioteca):
    def __init__(self, codigo, titulo, año, director, duracion, genero):
        super().__init__(codigo, titulo, año)
        self.director = director
        self.duracion = duracion
        self.genero = genero
    
    def tipo(self):
        return "DVD"

# Clase para usuarios de la biblioteca
class Usuario:
    def __init__(self, id, nombre, email):
        self.id = id
        self.nombre = nombre
        self.email = email
        self.prestamos = []
    
    def __str__(self):
        return f"Usuario: {self.nombre} ({self.id})"

# Clase para préstamos
class Prestamo:
    def __init__(self, elemento, usuario, dias=14):
        self.elemento = elemento
        self.usuario = usuario
        self.fecha_prestamo = datetime.now()
        self.fecha_devolucion_prevista = self.fecha_prestamo + timedelta(days=dias)
        self.fecha_devolucion_real = None
    
    @property
    def activo(self):
        return self.fecha_devolucion_real is None
    
    @property
    def dias_restantes(self):
        if not self.activo:
            return 0
        dias = (self.fecha_devolucion_prevista - datetime.now()).days
        return max(0, dias)
    
    @property
    def retrasado(self):
        if not self.activo:
            return False
        return datetime.now() > self.fecha_devolucion_prevista
    
    def devolver(self):
        self.fecha_devolucion_real = datetime.now()
        self.elemento.disponible = True
    
    def __str__(self):
        estado = "activo" if self.activo else "devuelto"
        return f"Préstamo de {self.elemento.titulo} a {self.usuario.nombre} - {estado}"

# Clase principal para la biblioteca
class Biblioteca:
    def __init__(self, nombre):
        self.nombre = nombre
        self.catalogo = {}
        self.usuarios = {}
        self.prestamos = []
    
    def agregar_elemento(self, elemento):
        self.catalogo[elemento.codigo] = elemento
    
    def agregar_usuario(self, usuario):
        self.usuarios[usuario.id] = usuario
    
    def buscar_por_titulo(self, titulo):
        return [e for e in self.catalogo.values() if titulo.lower() in e.titulo.lower()]
    
    def buscar_disponibles(self):
        return [e for e in self.catalogo.values() if e.disponible]
    
    def prestar(self, codigo_elemento, id_usuario, dias=14):
        if codigo_elemento not in self.catalogo:
            raise ValueError(f"Elemento con código {codigo_elemento} no encontrado")
        
        if id_usuario not in self.usuarios:
            raise ValueError(f"Usuario con ID {id_usuario} no encontrado")
        
        elemento = self.catalogo[codigo_elemento]
        usuario = self.usuarios[id_usuario]
        
        if not elemento.disponible:
            raise ValueError(f"El elemento {elemento.titulo} no está disponible")
        
        # Verificar si el usuario tiene préstamos activos retrasados
        prestamos_activos = [p for p in usuario.prestamos if p.activo]
        if any(p.retrasado for p in prestamos_activos):
            raise ValueError(f"El usuario {usuario.nombre} tiene préstamos retrasados")
        
        # Verificar si el usuario ya tiene 3 préstamos activos
        if len(prestamos_activos) >= 3:
            raise ValueError(f"El usuario {usuario.nombre} ya tiene el máximo de préstamos permitidos")
        
        # Crear el préstamo
        prestamo = Prestamo(elemento, usuario, dias)
        self.prestamos.append(prestamo)
        usuario.prestamos.append(prestamo)
        elemento.disponible = False
        
        return prestamo
    
    def devolver(self, codigo_elemento, id_usuario):
        # Buscar el préstamo activo para este elemento y usuario
        for prestamo in self.prestamos:
            if (prestamo.elemento.codigo == codigo_elemento and 
                prestamo.usuario.id == id_usuario and 
                prestamo.activo):
                prestamo.devolver()
                return prestamo
        
        raise ValueError(f"No se encontró un préstamo activo para el elemento {codigo_elemento} y el usuario {id_usuario}")
    
    def listar_prestamos_activos(self):
        return [p for p in self.prestamos if p.activo]
    
    def listar_prestamos_retrasados(self):
        return [p for p in self.prestamos if p.activo and p.retrasado]
    
    def __str__(self):
        return f"Biblioteca {self.nombre} - {len(self.catalogo)} elementos, {len(self.usuarios)} usuarios"

# Demostración del sistema
def demo_biblioteca():
    # Crear la biblioteca
    biblioteca = Biblioteca("Biblioteca Municipal")
    
    # Agregar elementos
    libro1 = Libro("L001", "Cien años de soledad", 1967, "Gabriel García Márquez", "Realismo mágico", 432)
    libro2 = Libro("L002", "El código Da Vinci", 2003, "Dan Brown", "Thriller", 589)
    libro3 = Libro("L003", "Harry Potter y la piedra filosofal", 1997, "J.K. Rowling", "Fantasía", 223)
    
    revista1 = Revista("R001", "National Geographic", 2022, "National Geographic Society", "Mensual", 256)
    revista2 = Revista("R002", "Scientific American", 2021, "Springer Nature", "Mensual", 124)
    
    dvd1 = DVD("D001", "El Padrino", 1972, "Francis Ford Coppola", 175, "Drama")
    dvd2 = DVD("D002", "Matrix", 1999, "Hermanas Wachowski", 136, "Ciencia ficción")
    
    biblioteca.agregar_elemento(libro1)
    biblioteca.agregar_elemento(libro2)
    biblioteca.agregar_elemento(libro3)
    biblioteca.agregar_elemento(revista1)
    biblioteca.agregar_elemento(revista2)
    biblioteca.agregar_elemento(dvd1)
    biblioteca.agregar_elemento(dvd2)
    
    # Agregar usuarios
    usuario1 = Usuario("U001", "Ana García", "ana@example.com")
    usuario2 = Usuario("U002", "Juan Pérez", "juan@example.com")
    usuario3 = Usuario("U003", "María López", "maria@example.com")
    
    biblioteca.agregar_usuario(usuario1)
    biblioteca.agregar_usuario(usuario2)
    biblioteca.agregar_usuario(usuario3)
    
    print(biblioteca)
    print("\nCatálogo:")
    for elemento in biblioteca.catalogo.values():
        print(f"- {elemento}")
    
    print("\nUsuarios:")
    for usuario in biblioteca.usuarios.values():
        print(f"- {usuario}")
    
    # Realizar préstamos
    print("\nRealizando préstamos:")
    try:
        prestamo1 = biblioteca.prestar("L001", "U001")
        print(f"Préstamo realizado: {prestamo1}")
        
        prestamo2 = biblioteca.prestar("D001", "U001")
        print(f"Préstamo realizado: {prestamo2}")
        
        prestamo3 = biblioteca.prestar("L002", "U002")
        print(f"Préstamo realizado: {prestamo3}")
    except ValueError as e:
        print(f"Error: {e}")
    
    # Listar préstamos activos
    print("\nPréstamos activos:")
    for prestamo in biblioteca.listar_prestamos_activos():
        print(f"- {prestamo} (Días restantes: {prestamo.dias_restantes})")
    
    # Devolver un préstamo
    print("\nDevolviendo préstamo:")
    try:
        prestamo_devuelto = biblioteca.devolver("L001", "U001")
        print(f"Préstamo devuelto: {prestamo_devuelto}")
    except ValueError as e:
        print(f"Error: {e}")
    
    # Listar préstamos activos después de la devolución
    print("\nPréstamos activos después de la devolución:")
    for prestamo in biblioteca.listar_prestamos_activos():
        print(f"- {prestamo} (Días restantes: {prestamo.dias_restantes})")
    
    # Buscar elementos por título
    print("\nBúsqueda por título 'harry':")
    resultados = biblioteca.buscar_por_titulo("harry")
    for elemento in resultados:
        print(f"- {elemento}")
    
    # Intentar prestar un elemento no disponible
    print("\nIntentando prestar un elemento no disponible:")
    try:
        prestamo_error = biblioteca.prestar("D001", "U002")  # D001 ya está prestado a U001
        print(f"Préstamo realizado: {prestamo_error}")
    except ValueError as e:
        print(f"Error: {e}")
    
    # Intentar prestar demasiados elementos a un usuario
    print("\nIntentando prestar demasiados elementos a un usuario:")
    try:
        prestamo4 = biblioteca.prestar("R001", "U001")
        print(f"Préstamo realizado: {prestamo4}")
        
        prestamo5 = biblioteca.prestar("R002", "U001")
        print(f"Préstamo realizado: {prestamo5}")
    except ValueError as e:
        print(f"Error: {e}")

# Ejecutar la demostración
demo_biblioteca()
```

## Buenas Prácticas

1. **Nombres Descriptivos**: Usa nombres claros y descriptivos para clases, métodos y atributos.

2. **Encapsulación**: Usa convenciones de nombres para indicar la visibilidad de atributos y métodos (p. ej., `_atributo` para protegido, `__atributo` para privado).

3. **Propiedades**: Usa propiedades en lugar de getters y setters explícitos para mantener una interfaz limpia.

4. **Herencia**: Usa herencia cuando exista una relación "es un". Para relaciones "tiene un", usa composición.

5. **Clases Abstractas**: Usa clases abstractas para definir interfaces comunes que deben implementar las subclases.

6. **Métodos Especiales**: Implementa métodos especiales para hacer que tus clases se comporten como tipos integrados de Python.

7. **Documentación**: Documenta tus clases y métodos con docstrings para facilitar su uso y mantenimiento.

8. **Principio de Responsabilidad Única**: Cada clase debe tener una única responsabilidad y razón para cambiar.

9. **Principio Abierto/Cerrado**: Las clases deben estar abiertas para extensión pero cerradas para modificación.

10. **Principio de Sustitución de Liskov**: Las subclases deben poder sustituir a sus clases base sin alterar el comportamiento esperado.

11. **Principio de Segregación de Interfaces**: Es mejor tener muchas interfaces específicas que una interfaz general.

12. **Principio de Inversión de Dependencias**: Depende de abstracciones, no de implementaciones concretas.

## Recursos Adicionales

- [Documentación oficial de Python sobre clases](https://docs.python.org/es/3/tutorial/classes.html)
- [Python's Instance, Class, and Static Methods Demystified](https://realpython.com/instance-class-and-static-methods-demystified/)
- [Inheritance and Composition: A Python OOP Guide](https://realpython.com/inheritance-composition-python/)
- [Python's super() considered super!](https://rhettinger.wordpress.com/2011/05/26/super-considered-super/)
- [Python 3 Patterns, Recipes and Idioms](https://python-3-patterns-idioms-test.readthedocs.io/en/latest/)
- [Design Patterns in Python](https://refactoring.guru/design-patterns/python)
- [Python Descriptors: An Introduction](https://realpython.com/python-descriptors/)
- [A Guide to Python's Magic Methods](https://rszalski.github.io/magicmethods/)

---

En la siguiente sección, exploraremos la programación concurrente en Python, incluyendo hilos, procesos, asyncio y más.