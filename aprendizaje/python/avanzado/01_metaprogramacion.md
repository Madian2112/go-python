# Metaprogramación en Python

## Introducción

La metaprogramación es una técnica avanzada que permite a los programas tratar el código como datos, manipulándolo y generándolo dinámicamente. Python, con su naturaleza dinámica y reflexiva, ofrece poderosas capacidades de metaprogramación que permiten crear abstracciones elegantes, reducir código repetitivo y construir frameworks flexibles.

En este módulo, exploraremos las técnicas avanzadas de metaprogramación en Python, incluyendo decoradores avanzados, metaclases, descriptores, y otras herramientas que permiten modificar el comportamiento del lenguaje y crear APIs elegantes.

## Decoradores Avanzados

### Repaso de Decoradores Básicos

Los decoradores son funciones que modifican el comportamiento de otras funciones o clases:

```python
def simple_decorator(func):
    def wrapper(*args, **kwargs):
        print("Antes de la función")
        result = func(*args, **kwargs)
        print("Después de la función")
        return result
    return wrapper

@simple_decorator
def saludar(nombre):
    print(f"Hola, {nombre}!")

saludar("Alice")  # Imprime: Antes de la función, Hola, Alice!, Después de la función
```

### Decoradores con Argumentos

Los decoradores pueden aceptar argumentos para personalizar su comportamiento:

```python
def repeat(n=1):
    def decorator(func):
        def wrapper(*args, **kwargs):
            results = []
            for _ in range(n):
                results.append(func(*args, **kwargs))
            return results
        return wrapper
    return decorator

@repeat(3)
def saludar(nombre):
    return f"Hola, {nombre}!"

print(saludar("Bob"))  # Imprime: ['Hola, Bob!', 'Hola, Bob!', 'Hola, Bob!']
```

### Decoradores de Clase

Los decoradores pueden aplicarse a clases para modificar su comportamiento:

```python
def add_greeting(cls):
    def say_hello(self):
        return f"Hola, soy {self.name}!"
    
    cls.say_hello = say_hello
    return cls

@add_greeting
class Person:
    def __init__(self, name):
        self.name = name

person = Person("Charlie")
print(person.say_hello())  # Imprime: Hola, soy Charlie!
```

### Decoradores con Estado

Los decoradores pueden mantener estado entre llamadas:

```python
def count_calls(func):
    def wrapper(*args, **kwargs):
        wrapper.calls += 1
        print(f"Llamada {wrapper.calls} a {func.__name__}")
        return func(*args, **kwargs)
    
    wrapper.calls = 0
    return wrapper

@count_calls
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

fibonacci(5)  # Imprime múltiples líneas con el conteo de llamadas
print(f"Total de llamadas: {fibonacci.calls}")
```

### Decoradores que Preservan Metadatos

Los decoradores pueden ocultar metadatos importantes como el nombre y la documentación de la función original. El módulo `functools` proporciona `wraps` para preservar estos metadatos:

```python
from functools import wraps

def log_execution(func):
    @wraps(func)
    def wrapper(*args, **kwargs):
        print(f"Ejecutando {func.__name__}")
        result = func(*args, **kwargs)
        print(f"Completado {func.__name__}")
        return result
    return wrapper

@log_execution
def suma(a, b):
    """Suma dos números."""
    return a + b

print(suma.__name__)  # Imprime: suma (en lugar de wrapper)
print(suma.__doc__)   # Imprime: Suma dos números.
```

### Decoradores Apilados

Se pueden aplicar múltiples decoradores a una función, que se ejecutan de abajo hacia arriba:

```python
def bold(func):
    @wraps(func)
    def wrapper(*args, **kwargs):
        return f"<b>{func(*args, **kwargs)}</b>"
    return wrapper

def italic(func):
    @wraps(func)
    def wrapper(*args, **kwargs):
        return f"<i>{func(*args, **kwargs)}</i>"
    return wrapper

@bold
@italic
def greeting(name):
    return f"Hola, {name}!"

print(greeting("Dave"))  # Imprime: <b><i>Hola, Dave!</i></b>
```

### Decoradores de Clase como Decoradores de Función

Los decoradores pueden implementarse como clases que implementan `__call__`:

```python
class Memoize:
    def __init__(self, func):
        self.func = func
        self.cache = {}
    
    def __call__(self, *args):
        if args not in self.cache:
            self.cache[args] = self.func(*args)
        return self.cache[args]

@Memoize
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

print(fibonacci(30))  # Rápido debido a la memoización
```

### Decoradores Paramétricos con Clases

Las clases también pueden usarse para crear decoradores con parámetros:

```python
class Retry:
    def __init__(self, max_attempts=3, exceptions=(Exception,)):
        self.max_attempts = max_attempts
        self.exceptions = exceptions
    
    def __call__(self, func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            attempts = 0
            while attempts < self.max_attempts:
                try:
                    return func(*args, **kwargs)
                except self.exceptions as e:
                    attempts += 1
                    if attempts == self.max_attempts:
                        raise
                    print(f"Intento {attempts} falló: {e}. Reintentando...")
            return None  # Nunca debería llegar aquí
        return wrapper

@Retry(max_attempts=5, exceptions=(ConnectionError, TimeoutError))
def request_data(url):
    # Simulación de una solicitud que puede fallar
    import random
    if random.random() < 0.7:
        raise ConnectionError("Conexión fallida")
    return "Datos recibidos"

try:
    print(request_data("https://example.com"))
except ConnectionError:
    print("Todos los intentos fallaron")
```

## Metaclases

### ¿Qué son las Metaclases?

Las metaclases son "clases de clases" - definen cómo se comportan las clases. En Python, las clases son objetos, y las metaclases son las clases que crean estos objetos.

```python
# La metaclase por defecto es type
class MiClase:
    pass

print(type(MiClase))  # Imprime: <class 'type'>

# Crear una clase dinámicamente con type
Dinamica = type('Dinamica', (object,), {'atributo': 42, 'metodo': lambda self: 'Hola'})
instancia = Dinamica()
print(instancia.atributo)  # Imprime: 42
print(instancia.metodo())  # Imprime: Hola
```

### Creando Metaclases Personalizadas

Las metaclases personalizadas permiten modificar el comportamiento de la creación de clases:

```python
class Meta(type):
    def __new__(mcs, name, bases, attrs):
        # Convertir todos los nombres de atributos a mayúsculas
        attrs_mayusculas = {
            key.upper() if not key.startswith('__') else key: value
            for key, value in attrs.items()
        }
        return super().__new__(mcs, name, bases, attrs_mayusculas)

class MiClase(metaclass=Meta):
    atributo = 42
    
    def metodo(self):
        return "Hola"

instancia = MiClase()
print(instancia.ATRIBUTO)  # Imprime: 42
print(instancia.METODO())  # Imprime: Hola
```

### Metaclases para Registro Automático

Las metaclases pueden usarse para registrar automáticamente subclases:

```python
class RegistroMeta(type):
    def __init__(cls, name, bases, attrs):
        super().__init__(name, bases, attrs)
        if not hasattr(cls, 'registro'):
            cls.registro = {}
        else:
            cls.registro[name] = cls

class Base(metaclass=RegistroMeta):
    pass

class A(Base):
    pass

class B(Base):
    pass

print(Base.registro)  # Imprime: {'A': <class '__main__.A'>, 'B': <class '__main__.B'>}
```

### Metaclases para Validación

Las metaclases pueden validar la estructura de una clase:

```python
class ValidacionMeta(type):
    def __new__(mcs, name, bases, attrs):
        # Verificar que la clase tiene un método 'validate'
        if 'validate' not in attrs:
            raise TypeError(f"La clase {name} debe implementar el método 'validate'")
        return super().__new__(mcs, name, bases, attrs)

class Validable(metaclass=ValidacionMeta):
    pass

# Esto lanzará un TypeError
# class MiClase(Validable):
#     pass

# Esto funcionará
class MiClase(Validable):
    def validate(self):
        return True
```

### Metaclases para Singletons

Las metaclases pueden implementar el patrón Singleton:

```python
class SingletonMeta(type):
    _instances = {}
    
    def __call__(cls, *args, **kwargs):
        if cls not in cls._instances:
            cls._instances[cls] = super().__call__(*args, **kwargs)
        return cls._instances[cls]

class Singleton(metaclass=SingletonMeta):
    def __init__(self, value):
        self.value = value

# Ambas variables referencian la misma instancia
s1 = Singleton(1)
s2 = Singleton(2)
print(s1 is s2)  # Imprime: True
print(s1.value)  # Imprime: 1 (no 2, porque s1 fue creado primero)
```

### Metaclases para ORM

Las metaclases son útiles para implementar ORMs (Object-Relational Mappers):

```python
class ModelMeta(type):
    def __new__(mcs, name, bases, attrs):
        # No aplicar a la clase base Model
        if name == 'Model':
            return super().__new__(mcs, name, bases, attrs)
        
        # Extraer campos
        fields = {}
        for key, value in list(attrs.items()):
            if isinstance(value, Field):
                fields[key] = value
                value.name = key
        
        attrs['_fields'] = fields
        attrs['_table'] = attrs.get('_table', name.lower())
        
        return super().__new__(mcs, name, bases, attrs)

class Field:
    def __init__(self, field_type, primary_key=False):
        self.field_type = field_type
        self.primary_key = primary_key
        self.name = None

class Model(metaclass=ModelMeta):
    def __init__(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self, key, value)
    
    @classmethod
    def create_table(cls):
        fields_sql = []
        for name, field in cls._fields.items():
            field_type = field.field_type
            primary_key = "PRIMARY KEY" if field.primary_key else ""
            fields_sql.append(f"{name} {field_type} {primary_key}".strip())
        
        sql = f"CREATE TABLE {cls._table} ({', '.join(fields_sql)});"
        return sql

class User(Model):
    _table = "users"
    id = Field("INTEGER", primary_key=True)
    name = Field("TEXT")
    email = Field("TEXT")

print(User.create_table())  # Imprime: CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, email TEXT);
```

## Descriptores

### ¿Qué son los Descriptores?

Los descriptores son objetos que implementan los métodos `__get__`, `__set__` o `__delete__`, permitiendo personalizar el acceso a atributos.

```python
class Descriptor:
    def __get__(self, instance, owner):
        print(f"Obteniendo desde {instance} de {owner}")
        return 42
    
    def __set__(self, instance, value):
        print(f"Estableciendo {value} en {instance}")
    
    def __delete__(self, instance):
        print(f"Eliminando de {instance}")

class MiClase:
    atributo = Descriptor()

obj = MiClase()
print(obj.atributo)  # Imprime: Obteniendo desde <__main__.MiClase object at ...> de <class '__main__.MiClase'>, 42
obj.atributo = 100   # Imprime: Estableciendo 100 en <__main__.MiClase object at ...>
del obj.atributo     # Imprime: Eliminando de <__main__.MiClase object at ...>
```

### Descriptores de Datos vs. No-Datos

- **Descriptores de datos**: Implementan `__set__` y/o `__delete__`
- **Descriptores de no-datos**: Solo implementan `__get__`

Los descriptores de datos tienen prioridad sobre los atributos de instancia, mientras que los descriptores de no-datos no.

```python
class DescriptorDeDatos:
    def __get__(self, instance, owner):
        return instance._valor
    
    def __set__(self, instance, value):
        instance._valor = value

class DescriptorDeNoDatos:
    def __get__(self, instance, owner):
        return 42

class Ejemplo:
    datos = DescriptorDeDatos()
    no_datos = DescriptorDeNoDatos()
    
    def __init__(self):
        self._valor = 0

obj = Ejemplo()
obj.datos = 100
print(obj.datos)  # Imprime: 100

obj.no_datos = 200  # Esto crea un atributo de instancia 'no_datos'
print(obj.no_datos)  # Imprime: 200 (el atributo de instancia tiene prioridad)
del obj.no_datos     # Elimina el atributo de instancia
print(obj.no_datos)  # Imprime: 42 (ahora usa el descriptor)
```

### Descriptores para Validación

Los descriptores son útiles para validar atributos:

```python
class Validado:
    def __init__(self, validacion=None, error_msg=None):
        self.validacion = validacion
        self.error_msg = error_msg
    
    def __set_name__(self, owner, name):
        self.name = name
        self.private_name = f"_{name}"
    
    def __get__(self, instance, owner):
        if instance is None:
            return self
        return getattr(instance, self.private_name, None)
    
    def __set__(self, instance, value):
        if self.validacion and not self.validacion(value):
            raise ValueError(self.error_msg or f"Valor inválido para {self.name}: {value}")
        setattr(instance, self.private_name, value)

class Persona:
    nombre = Validado(lambda x: isinstance(x, str) and len(x) > 0, "El nombre debe ser una cadena no vacía")
    edad = Validado(lambda x: isinstance(x, int) and 0 <= x <= 120, "La edad debe estar entre 0 y 120")
    
    def __init__(self, nombre, edad):
        self.nombre = nombre
        self.edad = edad

try:
    p = Persona("", -5)
except ValueError as e:
    print(f"Error: {e}")  # Imprime: Error: El nombre debe ser una cadena no vacía

p = Persona("Alice", 30)
print(p.nombre, p.edad)  # Imprime: Alice 30
```

### Descriptores para Propiedades Calculadas

Los descriptores pueden implementar propiedades calculadas:

```python
class Calculado:
    def __init__(self, func):
        self.func = func
    
    def __get__(self, instance, owner):
        if instance is None:
            return self
        return self.func(instance)

class Circulo:
    def __init__(self, radio):
        self.radio = radio
    
    @Calculado
    def area(self):
        import math
        return math.pi * self.radio ** 2
    
    @Calculado
    def perimetro(self):
        import math
        return 2 * math.pi * self.radio

c = Circulo(5)
print(f"Área: {c.area:.2f}")        # Imprime: Área: 78.54
print(f"Perímetro: {c.perimetro:.2f}")  # Imprime: Perímetro: 31.42
```

### Descriptores para Lazy Loading

Los descriptores pueden implementar carga perezosa de atributos:

```python
class LazyProperty:
    def __init__(self, func):
        self.func = func
        self.name = func.__name__
    
    def __get__(self, instance, owner):
        if instance is None:
            return self
        
        value = self.func(instance)
        setattr(instance, self.name, value)  # Reemplaza el descriptor con el valor calculado
        return value

class Datos:
    def __init__(self, filename):
        self.filename = filename
    
    @LazyProperty
    def contenido(self):
        print(f"Cargando contenido de {self.filename}...")
        # Simulación de carga de archivo
        return f"Contenido de {self.filename}"

d = Datos("archivo.txt")
print("Objeto creado, contenido no cargado aún")
print(d.contenido)  # Imprime: Cargando contenido de archivo.txt, Contenido de archivo.txt
print(d.contenido)  # Imprime: Contenido de archivo.txt (sin mensaje de carga)
```

## Introspección y Reflexión

### Introspección de Objetos

Python proporciona varias funciones para examinar objetos en tiempo de ejecución:

```python
class Ejemplo:
    """Una clase de ejemplo."""
    
    atributo = 42
    
    def __init__(self, valor):
        self.valor = valor
    
    def metodo(self, x):
        """Un método de ejemplo."""
        return x + self.valor

obj = Ejemplo(10)

# Examinar atributos y métodos
print(dir(obj))  # Lista todos los atributos y métodos

# Verificar tipos
print(isinstance(obj, Ejemplo))  # True
print(issubclass(Ejemplo, object))  # True

# Obtener documentación
print(Ejemplo.__doc__)  # Una clase de ejemplo.
print(Ejemplo.metodo.__doc__)  # Un método de ejemplo.

# Obtener tipo
print(type(obj))  # <class '__main__.Ejemplo'>

# Obtener atributos
print(getattr(obj, 'valor'))  # 10
print(hasattr(obj, 'atributo'))  # True
setattr(obj, 'nuevo_atributo', 100)
print(obj.nuevo_atributo)  # 100
delattr(obj, 'nuevo_atributo')
print(hasattr(obj, 'nuevo_atributo'))  # False
```

### Módulo `inspect`

El módulo `inspect` proporciona funciones más avanzadas para introspección:

```python
import inspect

# Obtener firma de función
def func(a, b=1, *args, **kwargs):
    pass

signature = inspect.signature(func)
print(signature)  # (a, b=1, *args, **kwargs)

for name, param in signature.parameters.items():
    print(f"{name}: {param.kind}, default={param.default}")

# Obtener código fuente
print(inspect.getsource(Ejemplo))

# Obtener línea de definición
print(inspect.getsourcelines(Ejemplo)[1])  # Número de línea

# Obtener miembros
for name, member in inspect.getmembers(Ejemplo):
    if not name.startswith('__'):
        print(f"{name}: {member}")
```

### Modificación Dinámica de Clases

Python permite modificar clases en tiempo de ejecución:

```python
class Dinamica:
    pass

# Añadir atributos y métodos dinámicamente
Dinamica.atributo = 42
Dinamica.metodo = lambda self, x: x * 2

# Añadir métodos con funciones definidas externamente
def nuevo_metodo(self, x, y):
    return x + y

Dinamica.suma = nuevo_metodo

# Usar la clase modificada
obj = Dinamica()
print(obj.atributo)  # 42
print(obj.metodo(5))  # 10
print(obj.suma(3, 4))  # 7
```

## Técnicas Avanzadas

### Módulo `abc` para Clases Abstractas

El módulo `abc` permite definir clases y métodos abstractos:

```python
from abc import ABC, abstractmethod

class Animal(ABC):
    @abstractmethod
    def hacer_sonido(self):
        pass
    
    @abstractmethod
    def moverse(self):
        pass

# No se puede instanciar directamente
# animal = Animal()  # TypeError

class Perro(Animal):
    def hacer_sonido(self):
        return "Guau"
    
    def moverse(self):
        return "Corriendo"

perro = Perro()
print(perro.hacer_sonido())  # Guau
print(perro.moverse())  # Corriendo
```

### Módulo `contextlib` para Manejadores de Contexto

El módulo `contextlib` simplifica la creación de manejadores de contexto:

```python
from contextlib import contextmanager

@contextmanager
def tempdir():
    import tempfile
    import shutil
    import os
    
    dir_path = tempfile.mkdtemp()
    try:
        print(f"Creado directorio temporal: {dir_path}")
        yield dir_path
    finally:
        print(f"Eliminando directorio temporal: {dir_path}")
        shutil.rmtree(dir_path)

with tempdir() as path:
    # Crear un archivo en el directorio temporal
    with open(os.path.join(path, "test.txt"), "w") as f:
        f.write("Hola, mundo!")
    print(f"Archivo creado en {path}")
# El directorio temporal se elimina automáticamente al salir del bloque with
```

### Módulo `functools` para Programación Funcional

El módulo `functools` proporciona herramientas para programación funcional:

```python
from functools import partial, lru_cache, reduce

# partial: crea una nueva función con algunos argumentos fijos
def multiplicar(x, y):
    return x * y

duplicar = partial(multiplicar, 2)
print(duplicar(5))  # 10

# lru_cache: memoriza los resultados de una función
@lru_cache(maxsize=None)
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

print(fibonacci(100))  # Rápido debido a la memorización

# reduce: aplica una función acumulativa a una secuencia
numeros = [1, 2, 3, 4, 5]
producto = reduce(lambda x, y: x * y, numeros)
print(producto)  # 120
```

### Módulo `operator` para Operaciones Eficientes

El módulo `operator` proporciona funciones para operaciones comunes:

```python
import operator

# Operaciones aritméticas
print(operator.add(5, 3))  # 8
print(operator.mul(5, 3))  # 15

# Operaciones de comparación
print(operator.lt(5, 3))  # False (5 < 3)
print(operator.eq(5, 5))  # True (5 == 5)

# Operaciones de atributos e ítems
class Ejemplo:
    def __init__(self):
        self.valor = 42

obj = Ejemplo()
print(operator.attrgetter('valor')(obj))  # 42

diccionario = {'a': 1, 'b': 2}
print(operator.itemgetter('a')(diccionario))  # 1

# Uso con funciones de orden superior
from functools import reduce
numeros = [1, 2, 3, 4, 5]
print(reduce(operator.add, numeros))  # 15
```

### Módulo `itertools` para Iteración Eficiente

El módulo `itertools` proporciona funciones para trabajar con iteradores:

```python
import itertools

# count: contador infinito
for i in itertools.islice(itertools.count(10, 2), 5):
    print(i, end=' ')  # 10 12 14 16 18
print()

# cycle: ciclo infinito
for i, c in enumerate(itertools.islice(itertools.cycle('ABC'), 7)):
    print(c, end=' ')  # A B C A B C A
print()

# combinations: todas las combinaciones posibles
for combo in itertools.combinations('ABC', 2):
    print(''.join(combo), end=' ')  # AB AC BC
print()

# permutations: todas las permutaciones posibles
for perm in itertools.permutations('ABC', 2):
    print(''.join(perm), end=' ')  # AB AC BA BC CA CB
print()

# product: producto cartesiano
for prod in itertools.product('AB', '12'):
    print(''.join(prod), end=' ')  # A1 A2 B1 B2
print()

# groupby: agrupa elementos consecutivos
for key, group in itertools.groupby('AAABBBCCAABB'):
    print(key, list(group), end=' | ')  # A ['A', 'A', 'A'] | B ['B', 'B', 'B'] | C ['C', 'C'] | A ['A', 'A'] | B ['B', 'B'] |
print()
```

## Patrones de Metaprogramación

### Patrón Singleton con Metaclases

```python
class SingletonMeta(type):
    _instances = {}
    
    def __call__(cls, *args, **kwargs):
        if cls not in cls._instances:
            cls._instances[cls] = super().__call__(*args, **kwargs)
        return cls._instances[cls]

class Config(metaclass=SingletonMeta):
    def __init__(self):
        self.database = None
        self.api_key = None
    
    def setup(self, database, api_key):
        self.database = database
        self.api_key = api_key

# Ambas variables referencian la misma instancia
config1 = Config()
config1.setup("mysql://localhost", "abc123")

config2 = Config()
print(config2.database)  # mysql://localhost
print(config2.api_key)   # abc123
print(config1 is config2)  # True
```

### Patrón Factory con Registro Automático

```python
class PluginRegistry(type):
    plugins = {}
    
    def __new__(mcs, name, bases, attrs):
        cls = super().__new__(mcs, name, bases, attrs)
        if bases:  # No registrar la clase base
            mcs.plugins[name] = cls
        return cls

class Plugin(metaclass=PluginRegistry):
    def process(self, data):
        raise NotImplementedError

class TextPlugin(Plugin):
    def process(self, data):
        return f"Procesando texto: {data}"

class ImagePlugin(Plugin):
    def process(self, data):
        return f"Procesando imagen: {data}"

# Crear un plugin dinámicamente
def create_plugin(name, process_func):
    return type(name, (Plugin,), {'process': process_func})

VideoPlugin = create_plugin('VideoPlugin', lambda self, data: f"Procesando video: {data}")

# Usar el registro de plugins
print(PluginRegistry.plugins)  # {'TextPlugin': <class '__main__.TextPlugin'>, 'ImagePlugin': <class '__main__.ImagePlugin'>, 'VideoPlugin': <class '__main__.VideoPlugin'>}

def process_data(data_type, data):
    plugin_class = PluginRegistry.plugins.get(f"{data_type.capitalize()}Plugin")
    if not plugin_class:
        raise ValueError(f"No hay plugin para el tipo {data_type}")
    plugin = plugin_class()
    return plugin.process(data)

print(process_data("text", "Hola, mundo"))  # Procesando texto: Hola, mundo
print(process_data("image", "foto.jpg"))   # Procesando imagen: foto.jpg
print(process_data("video", "video.mp4"))  # Procesando video: video.mp4
```

### Patrón Builder con Fluent Interface

```python
class QueryBuilder:
    def __init__(self):
        self.table = None
        self.fields = []
        self.conditions = []
        self.order_by = None
        self.limit = None
    
    def select(self, *fields):
        self.fields = fields if fields else ['*']
        return self
    
    def from_table(self, table):
        self.table = table
        return self
    
    def where(self, condition):
        self.conditions.append(condition)
        return self
    
    def order(self, field, direction='ASC'):
        self.order_by = (field, direction)
        return self
    
    def limit_to(self, limit):
        self.limit = limit
        return self
    
    def build(self):
        if not self.table:
            raise ValueError("Debe especificar una tabla")
        
        query = f"SELECT {', '.join(self.fields)} FROM {self.table}"
        
        if self.conditions:
            query += f" WHERE {' AND '.join(self.conditions)}"
        
        if self.order_by:
            field, direction = self.order_by
            query += f" ORDER BY {field} {direction}"
        
        if self.limit is not None:
            query += f" LIMIT {self.limit}"
        
        return query + ";"

# Uso del builder
query = QueryBuilder()\
    .select("id", "name", "email")\
    .from_table("users")\
    .where("age > 18")\
    .where("status = 'active'")\
    .order("name")\
    .limit_to(10)\
    .build()

print(query)  # SELECT id, name, email FROM users WHERE age > 18 AND status = 'active' ORDER BY name ASC LIMIT 10;
```

### Patrón Proxy con Descriptores

```python
class LazyLoader:
    def __init__(self, cls, *args, **kwargs):
        self.cls = cls
        self.args = args
        self.kwargs = kwargs
        self._instance = None
    
    def __get__(self, instance, owner):
        if instance is None:
            return self
        
        if self._instance is None:
            print(f"Cargando instancia de {self.cls.__name__}...")
            self._instance = self.cls(*self.args, **self.kwargs)
        
        return self._instance

class ExpensiveResource:
    def __init__(self, name):
        self.name = name
        print(f"Inicializando recurso costoso: {name}")
    
    def use(self):
        return f"Usando {self.name}"

class App:
    # El recurso no se carga hasta que se accede a él
    resource = LazyLoader(ExpensiveResource, "Database Connection")
    
    def __init__(self):
        print("Inicializando aplicación")
    
    def run(self):
        print("Ejecutando aplicación")
        print(self.resource.use())  # Aquí se carga el recurso

app = App()
print("Aplicación creada, recurso no cargado aún")
app.run()
```

## Ejercicios Prácticos

1. **Implementa un decorador `@validate` que valide los argumentos de una función según un esquema**:
   - Debe permitir especificar tipos y restricciones para cada argumento.
   - Debe lanzar excepciones descriptivas cuando la validación falla.

2. **Crea una metaclase `EnumMeta` para implementar enumeraciones**:
   - Debe convertir atributos de clase en instancias de una clase Enum.
   - Debe prevenir la creación de instancias duplicadas.
   - Debe proporcionar métodos para iterar sobre todos los valores.

3. **Implementa un descriptor `@property_with_history` que mantenga un historial de valores**:
   - Debe funcionar como `@property` pero guardar todos los valores asignados.
   - Debe proporcionar un método para acceder al historial de valores.

4. **Crea un framework de inyección de dependencias usando metaprogramación**:
   - Debe permitir registrar servicios y sus dependencias.
   - Debe resolver automáticamente las dependencias al crear instancias.
   - Debe manejar dependencias circulares.

5. **Implementa un ORM simple usando metaclases y descriptores**:
   - Debe permitir definir modelos con campos tipados.
   - Debe generar automáticamente SQL para crear tablas.
   - Debe proporcionar métodos para guardar, cargar y consultar objetos.

## Conclusión

La metaprogramación en Python es una herramienta poderosa que permite crear abstracciones elegantes, reducir código repetitivo y construir frameworks flexibles. A través de decoradores, metaclases, descriptores y otras técnicas, puedes modificar el comportamiento del lenguaje y crear APIs que sean más expresivas y fáciles de usar.

Sin embargo, con gran poder viene gran responsabilidad. La metaprogramación puede hacer que el código sea más difícil de entender y depurar si se usa en exceso o de manera inapropiada. Es importante encontrar un equilibrio entre la elegancia y la claridad, y documentar bien cualquier uso de metaprogramación para que otros desarrolladores (incluido tu futuro yo) puedan entender fácilmente cómo funciona el código.

Cuando se usa sabiamente, la metaprogramación puede elevar tu código Python a nuevos niveles de expresividad y reutilización, permitiéndote escribir código más limpio, más mantenible y más poderoso.