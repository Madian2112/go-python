# Programación Funcional Avanzada en Python

## Introducción

La programación funcional es un paradigma que trata la computación como la evaluación de funciones matemáticas y evita el cambio de estado y los datos mutables. Python, aunque no es un lenguaje puramente funcional, ofrece muchas características que permiten adoptar este estilo de programación.

En este documento, exploraremos técnicas avanzadas de programación funcional en Python, sus beneficios y aplicaciones prácticas.

## Fundamentos de Programación Funcional en Python

### Funciones de Primera Clase y Funciones de Orden Superior

En Python, las funciones son objetos de primera clase, lo que significa que pueden ser:
- Asignadas a variables
- Pasadas como argumentos a otras funciones
- Retornadas desde otras funciones
- Almacenadas en estructuras de datos

```python
# Asignar función a variable
def add(a, b):
    return a + b

operation = add
print(operation(3, 4))  # 7

# Pasar función como argumento
def apply_operation(a, b, operation):
    return operation(a, b)

print(apply_operation(3, 4, add))  # 7

# Retornar función desde otra función
def get_multiplier(factor):
    def multiply(x):
        return x * factor
    return multiply

double = get_multiplier(2)
triple = get_multiplier(3)
print(double(5))  # 10
print(triple(5))  # 15

# Almacenar funciones en estructuras de datos
operations = {
    'add': add,
    'double': double,
    'triple': triple
}

print(operations['add'](3, 4))  # 7
```

### Funciones Puras e Inmutabilidad

Una función pura es aquella que:
1. Siempre produce el mismo resultado para los mismos argumentos
2. No tiene efectos secundarios (no modifica variables externas, no realiza I/O, etc.)

```python
# Función pura
def add_pure(a, b):
    return a + b

# Función impura (tiene efecto secundario)
total = 0
def add_impure(a, b):
    global total
    total += a + b
    return total
```

La inmutabilidad es un principio clave en la programación funcional. Python ofrece varios tipos de datos inmutables:
- Tuplas
- Strings
- Frozensets
- Números (int, float, complex)

```python
# Uso de tuplas para representar datos inmutables
point = (10, 20)

# En lugar de modificar el punto, creamos uno nuevo
def move_point(point, delta_x, delta_y):
    x, y = point
    return (x + delta_x, y + delta_y)

new_point = move_point(point, 5, 5)
print(point)      # (10, 20) - no cambia
print(new_point)  # (15, 25) - nuevo punto
```

## Técnicas Avanzadas de Programación Funcional

### Funciones Lambda

Las funciones lambda son funciones anónimas de una sola expresión.

```python
# Función lambda para sumar dos números
add = lambda a, b: a + b
print(add(3, 4))  # 7

# Uso de lambda con sorted
people = [('Alice', 25), ('Bob', 30), ('Charlie', 20)]
sorted_by_age = sorted(people, key=lambda person: person[1])
print(sorted_by_age)  # [('Charlie', 20), ('Alice', 25), ('Bob', 30)]

# Uso de lambda con filter
numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
even_numbers = list(filter(lambda x: x % 2 == 0, numbers))
print(even_numbers)  # [2, 4, 6, 8, 10]
```

### Funciones map, filter y reduce

Estas funciones son herramientas fundamentales en programación funcional:

```python
from functools import reduce

numbers = [1, 2, 3, 4, 5]

# map: aplica una función a cada elemento de un iterable
squared = list(map(lambda x: x**2, numbers))
print(squared)  # [1, 4, 9, 16, 25]

# filter: filtra elementos según un predicado
even = list(filter(lambda x: x % 2 == 0, numbers))
print(even)  # [2, 4]

# reduce: combina todos los elementos en un único valor
sum_all = reduce(lambda acc, x: acc + x, numbers, 0)
print(sum_all)  # 15

# Combinando map, filter y reduce
sum_of_squares_of_even = reduce(
    lambda acc, x: acc + x,
    map(lambda x: x**2, filter(lambda x: x % 2 == 0, numbers)),
    0
)
print(sum_of_squares_of_even)  # 20 (4 + 16)
```

### Comprensiones de Listas, Diccionarios y Conjuntos

Las comprensiones son una forma concisa y funcional de crear colecciones:

```python
numbers = [1, 2, 3, 4, 5]

# Comprensión de lista
squared = [x**2 for x in numbers]
print(squared)  # [1, 4, 9, 16, 25]

# Comprensión de lista con condición
even_squared = [x**2 for x in numbers if x % 2 == 0]
print(even_squared)  # [4, 16]

# Comprensión de diccionario
square_dict = {x: x**2 for x in numbers}
print(square_dict)  # {1: 1, 2: 4, 3: 9, 4: 16, 5: 25}

# Comprensión de conjunto
square_set = {x**2 for x in numbers}
print(square_set)  # {1, 4, 9, 16, 25}

# Comprensión de generador
square_gen = (x**2 for x in numbers)
print(list(square_gen))  # [1, 4, 9, 16, 25]
```

### Closures y Funciones Anidadas

Los closures son funciones que capturan variables de su entorno léxico:

```python
def counter():
    count = 0
    def increment():
        nonlocal count
        count += 1
        return count
    return increment

c1 = counter()
c2 = counter()

print(c1())  # 1
print(c1())  # 2
print(c2())  # 1 (c2 tiene su propio estado)

# Uso de closures para crear funciones parciales
def multiplier(factor):
    def multiply(x):
        return x * factor
    return multiply

double = multiplier(2)
triple = multiplier(3)
print(double(5))  # 10
print(triple(5))  # 15
```

### Decoradores Avanzados

Los decoradores son una forma poderosa de modificar o extender el comportamiento de funciones:

```python
import functools
import time

# Decorador para medir tiempo de ejecución
def timer(func):
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        start_time = time.time()
        result = func(*args, **kwargs)
        end_time = time.time()
        print(f"{func.__name__} ejecutado en {end_time - start_time:.4f} segundos")
        return result
    return wrapper

# Decorador con argumentos
def repeat(n=1):
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            results = []
            for _ in range(n):
                results.append(func(*args, **kwargs))
            return results
        return wrapper
    return decorator

# Uso de decoradores
@timer
@repeat(3)
def slow_function(delay):
    time.sleep(delay)
    return delay

print(slow_function(0.1))
```

### Recursión y Recursión de Cola

La recursión es una técnica común en programación funcional:

```python
# Factorial recursivo
def factorial(n):
    if n <= 1:
        return 1
    return n * factorial(n-1)

print(factorial(5))  # 120

# Recursión de cola (no optimizada en Python)
def factorial_tail(n, acc=1):
    if n <= 1:
        return acc
    return factorial_tail(n-1, n*acc)

print(factorial_tail(5))  # 120

# Fibonacci recursivo con memoización
def fibonacci(n, memo={}):
    if n in memo:
        return memo[n]
    if n <= 1:
        return n
    memo[n] = fibonacci(n-1, memo) + fibonacci(n-2, memo)
    return memo[n]

print(fibonacci(10))  # 55
```

### Funciones Parciales

Las funciones parciales permiten fijar algunos argumentos de una función:

```python
from functools import partial

def power(base, exponent):
    return base ** exponent

# Crear función parcial para calcular cuadrados
square = partial(power, exponent=2)

# Crear función parcial para calcular cubos
cube = partial(power, exponent=3)

print(square(5))  # 25
print(cube(5))    # 125

# Otro ejemplo: convertir temperaturas
def convert_temperature(degrees, from_scale, to_scale):
    scales = {
        'C': lambda deg: deg,
        'F': lambda deg: (deg - 32) * 5/9,
        'K': lambda deg: deg - 273.15
    }
    
    # Convertir a Celsius primero
    celsius = scales[from_scale](degrees)
    
    # Luego convertir de Celsius a la escala destino
    if to_scale == 'C':
        return celsius
    elif to_scale == 'F':
        return celsius * 9/5 + 32
    elif to_scale == 'K':
        return celsius + 273.15

# Crear funciones parciales para conversiones específicas
f_to_c = partial(convert_temperature, from_scale='F', to_scale='C')
c_to_f = partial(convert_temperature, from_scale='C', to_scale='F')

print(f_to_c(32))   # 0.0
print(c_to_f(0))    # 32.0
```

### Composición de Funciones

La composición de funciones es una técnica para combinar múltiples funciones en una sola:

```python
# Implementación básica de composición
def compose(f, g):
    return lambda x: f(g(x))

# Funciones para componer
def add_one(x):
    return x + 1

def double(x):
    return x * 2

# Componer funciones
double_then_add_one = compose(add_one, double)
add_one_then_double = compose(double, add_one)

print(double_then_add_one(3))  # 3*2 + 1 = 7
print(add_one_then_double(3))  # (3+1)*2 = 8

# Composición de múltiples funciones
def compose_multiple(*functions):
    def compose_two(f, g):
        return lambda x: f(g(x))
    
    if not functions:
        return lambda x: x  # Función identidad
    
    return functools.reduce(compose_two, functions)

# Ejemplo de uso
def square(x):
    return x * x

pipeline = compose_multiple(add_one, double, square)
print(pipeline(3))  # square(double(add_one(3))) = square(double(4)) = square(8) = 64
```

## Patrones Funcionales Avanzados

### Mónadas en Python

Aunque Python no tiene soporte nativo para mónadas, podemos implementarlas:

```python
# Implementación de la mónada Maybe (Option)
class Maybe:
    def __init__(self, value=None):
        self.value = value
    
    @classmethod
    def just(cls, value):
        return cls(value)
    
    @classmethod
    def nothing(cls):
        return cls(None)
    
    def is_just(self):
        return self.value is not None
    
    def is_nothing(self):
        return self.value is None
    
    def map(self, func):
        if self.is_nothing():
            return self
        return Maybe.just(func(self.value))
    
    def flat_map(self, func):
        if self.is_nothing():
            return self
        return func(self.value)
    
    def get_or_else(self, default):
        if self.is_nothing():
            return default
        return self.value
    
    def __str__(self):
        if self.is_just():
            return f"Just({self.value})"
        return "Nothing"

# Ejemplo de uso
def find_user(user_id):
    users = {1: "Alice", 2: "Bob"}
    if user_id in users:
        return Maybe.just(users[user_id])
    return Maybe.nothing()

# Buscar usuario y transformar
result = find_user(1).map(lambda name: f"User: {name}")
print(result)  # Just(User: Alice)

# Buscar usuario inexistente
result = find_user(3).map(lambda name: f"User: {name}")
print(result)  # Nothing

# Encadenamiento de operaciones
result = find_user(1)\
    .map(lambda name: name.upper())\
    .map(lambda name: f"USER: {name}")
print(result)  # Just(USER: ALICE)

# Uso de flat_map para operaciones que pueden fallar
def get_department(username):
    departments = {"Alice": "Engineering", "Charlie": "Marketing"}
    if username in departments:
        return Maybe.just(departments[username])
    return Maybe.nothing()

result = find_user(1).flat_map(get_department)
print(result)  # Just(Engineering)

result = find_user(2).flat_map(get_department)
print(result)  # Nothing
```

### Implementación de Result (Either)

El patrón Result es útil para manejar operaciones que pueden fallar:

```python
# Implementación de la mónada Result (Either)
class Result:
    def __init__(self, value=None, error=None):
        self.value = value
        self.error = error
        self.is_success = error is None
    
    @classmethod
    def success(cls, value):
        return cls(value=value)
    
    @classmethod
    def failure(cls, error):
        return cls(error=error)
    
    def map(self, func):
        if not self.is_success:
            return self
        try:
            return Result.success(func(self.value))
        except Exception as e:
            return Result.failure(str(e))
    
    def flat_map(self, func):
        if not self.is_success:
            return self
        try:
            return func(self.value)
        except Exception as e:
            return Result.failure(str(e))
    
    def get_or_else(self, default):
        if not self.is_success:
            return default
        return self.value
    
    def get_or_raise(self):
        if not self.is_success:
            raise Exception(self.error)
        return self.value
    
    def __str__(self):
        if self.is_success:
            return f"Success({self.value})"
        return f"Failure({self.error})"

# Ejemplo de uso
def divide(a, b):
    if b == 0:
        return Result.failure("Division by zero")
    return Result.success(a / b)

# División exitosa
result = divide(10, 2).map(lambda x: x * 2)
print(result)  # Success(10.0)

# División por cero
result = divide(10, 0).map(lambda x: x * 2)
print(result)  # Failure(Division by zero)

# Encadenamiento de operaciones
result = divide(10, 2)\
    .map(lambda x: x + 5)\
    .flat_map(lambda x: divide(x, 0))
print(result)  # Failure(Division by zero)
```

### Implementación de Pipe (Pipeline)

El patrón Pipe permite encadenar operaciones de manera funcional:

```python
class Pipe:
    def __init__(self, value):
        self.value = value
    
    def pipe(self, func):
        return Pipe(func(self.value))
    
    def value(self):
        return self.value

# Ejemplo de uso
result = Pipe(5)\
    .pipe(lambda x: x * 2)\
    .pipe(lambda x: x + 1)\
    .pipe(lambda x: x * x)\
    .value

print(result)  # 121
```

### Curry y Aplicación Parcial

El currying es una técnica que transforma una función de múltiples argumentos en una secuencia de funciones de un solo argumento:

```python
# Implementación básica de curry
def curry(func):
    def curried(*args, **kwargs):
        if len(args) + len(kwargs) >= func.__code__.co_argcount:
            return func(*args, **kwargs)
        return lambda *more_args, **more_kwargs: curried(*(args + more_args), **{**kwargs, **more_kwargs})
    return curried

# Ejemplo de uso
@curry
def add(a, b, c):
    return a + b + c

# Diferentes formas de llamar a la función currificada
print(add(1)(2)(3))      # 6
print(add(1, 2)(3))       # 6
print(add(1)(2, 3))       # 6
print(add(1, 2, 3))       # 6

# Uso práctico: filtrar una lista con diferentes predicados
@curry
def filter_by(predicate, iterable):
    return list(filter(predicate, iterable))

numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

is_even = lambda x: x % 2 == 0
is_odd = lambda x: x % 2 != 0
greater_than_5 = lambda x: x > 5

filter_even = filter_by(is_even)
filter_odd = filter_by(is_odd)
filter_greater_than_5 = filter_by(greater_than_5)

print(filter_even(numbers))          # [2, 4, 6, 8, 10]
print(filter_odd(numbers))           # [1, 3, 5, 7, 9]
print(filter_greater_than_5(numbers)) # [6, 7, 8, 9, 10]
```

## Aplicaciones Prácticas

### Procesamiento de Datos Funcional

```python
import csv
from functools import reduce

# Leer datos de un archivo CSV
def read_csv(file_path):
    with open(file_path, 'r') as file:
        reader = csv.DictReader(file)
        return list(reader)

# Transformar datos usando técnicas funcionales
def process_data(data):
    # Filtrar registros con precio > 100
    expensive_items = filter(lambda item: float(item['price']) > 100, data)
    
    # Transformar cada registro
    transformed = map(
        lambda item: {
            'name': item['name'].upper(),
            'price': float(item['price']),
            'category': item['category']
        },
        expensive_items
    )
    
    # Agrupar por categoría
    def group_by_category(acc, item):
        category = item['category']
        if category not in acc:
            acc[category] = []
        acc[category].append(item)
        return acc
    
    grouped = reduce(group_by_category, transformed, {})
    
    # Calcular precio promedio por categoría
    return {
        category: sum(item['price'] for item in items) / len(items)
        for category, items in grouped.items()
    }

# Ejemplo de uso
# data = read_csv('products.csv')
# result = process_data(data)
# print(result)
```

### Validación Funcional

```python
# Implementación de validadores componibles
def validate(value, *validators):
    errors = [error for validator in validators
              if (error := validator(value)) is not None]
    return errors if errors else None

# Validadores específicos
def min_length(min_len):
    return lambda value: f"Debe tener al menos {min_len} caracteres" \
                        if len(value) < min_len else None

def max_length(max_len):
    return lambda value: f"Debe tener como máximo {max_len} caracteres" \
                        if len(value) > max_len else None

def contains_digit():
    return lambda value: "Debe contener al menos un dígito" \
                        if not any(c.isdigit() for c in value) else None

def contains_uppercase():
    return lambda value: "Debe contener al menos una letra mayúscula" \
                        if not any(c.isupper() for c in value) else None

# Ejemplo de uso
def validate_password(password):
    return validate(
        password,
        min_length(8),
        max_length(20),
        contains_digit(),
        contains_uppercase()
    )

# Probar validación
passwords = ["abc", "abcdefgh", "abcdefgh1", "Abcdefgh1"]
for password in passwords:
    errors = validate_password(password)
    if errors:
        print(f"'{password}' no es válido: {', '.join(errors)}")
    else:
        print(f"'{password}' es válido")
```

### Manejo de Configuración Inmutable

```python
from types import MappingProxyType
import copy

class Config:
    def __init__(self, config_dict=None):
        self._config = MappingProxyType(copy.deepcopy(config_dict or {}))
    
    def get(self, key, default=None):
        return self._config.get(key, default)
    
    def with_update(self, **kwargs):
        new_config = dict(self._config)
        new_config.update(kwargs)
        return Config(new_config)
    
    def __getitem__(self, key):
        return self._config[key]
    
    def __contains__(self, key):
        return key in self._config
    
    def __str__(self):
        return str(dict(self._config))

# Ejemplo de uso
default_config = Config({
    'host': 'localhost',
    'port': 8080,
    'debug': False,
    'database': {
        'host': 'localhost',
        'port': 5432,
        'name': 'mydb'
    }
})

# Crear nueva configuración basada en la anterior
dev_config = default_config.with_update(
    debug=True,
    database={'host': 'localhost', 'port': 5432, 'name': 'devdb'}
)

print(default_config)  # No cambia
print(dev_config)      # Configuración actualizada
```

## Rendimiento y Optimización

### Memoización

La memoización es una técnica de optimización que almacena los resultados de llamadas a funciones costosas:

```python
import functools
import time

# Decorador de memoización
def memoize(func):
    cache = {}
    
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        # Convertir kwargs a tupla de tuplas para hacerlo hashable
        kwargs_tuple = tuple(sorted(kwargs.items()))
        key = (args, kwargs_tuple)
        
        if key not in cache:
            cache[key] = func(*args, **kwargs)
        return cache[key]
    
    return wrapper

# Ejemplo: Fibonacci con y sin memoización
def fibonacci_slow(n):
    if n <= 1:
        return n
    return fibonacci_slow(n-1) + fibonacci_slow(n-2)

@memoize
def fibonacci_fast(n):
    if n <= 1:
        return n
    return fibonacci_fast(n-1) + fibonacci_fast(n-2)

# Comparar rendimiento
def measure_time(func, *args, **kwargs):
    start = time.time()
    result = func(*args, **kwargs)
    end = time.time()
    print(f"{func.__name__} tomó {end - start:.6f} segundos")
    return result

# Calcular fibonacci(30) con ambas implementaciones
n = 30
print(f"Calculando fibonacci({n})...")
measure_time(fibonacci_fast, n)
measure_time(fibonacci_slow, n)
```

### Evaluación Perezosa (Lazy Evaluation)

La evaluación perezosa permite calcular valores solo cuando son necesarios:

```python
class LazySequence:
    def __init__(self, generator_func):
        self.generator_func = generator_func
        self._cached = []
    
    def __getitem__(self, index):
        while index >= len(self._cached):
            try:
                self._cached.append(next(self.generator_func))
            except StopIteration:
                raise IndexError("Sequence index out of range")
        return self._cached[index]
    
    def __iter__(self):
        for item in self._cached:
            yield item
        
        while True:
            try:
                next_item = next(self.generator_func)
                self._cached.append(next_item)
                yield next_item
            except StopIteration:
                break

# Ejemplo: secuencia infinita de Fibonacci
def fibonacci_generator():
    a, b = 0, 1
    while True:
        yield a
        a, b = b, a + b

# Crear secuencia perezosa
fib_sequence = LazySequence(fibonacci_generator())

# Acceder a elementos específicos
print(fib_sequence[0])   # 0
print(fib_sequence[1])   # 1
print(fib_sequence[10])  # 55

# Iterar sobre los primeros n elementos
for i, fib in enumerate(fib_sequence):
    if i >= 15:
        break
    print(f"Fibonacci({i}) = {fib}")
```

## Ejercicios Prácticos

1. **Implementar una biblioteca de utilidades funcionales**:
   - Crear implementaciones de Map, Filter, Reduce que soporten operaciones en paralelo.
   - Implementar funciones para composición y currying.
   - Añadir documentación y pruebas.

2. **Refactorizar código imperativo a estilo funcional**:
   - Tomar un programa existente con muchos efectos secundarios.
   - Refactorizarlo para usar principios funcionales.
   - Comparar legibilidad, mantenibilidad y rendimiento.

3. **Implementar un procesador de datos en pipeline**:
   - Crear un pipeline para procesar datos de un archivo.
   - Usar composición de funciones para transformar los datos.
   - Implementar manejo de errores funcional.

4. **Crear una mini-biblioteca de estructuras de datos inmutables**:
   - Implementar versiones inmutables de listas, mapas, etc.
   - Proporcionar operaciones funcionales para estas estructuras.
   - Comparar con las estructuras de datos mutables de Python.

5. **Desarrollar un mini-framework para manejo de efectos**:
   - Implementar mónadas para IO, estado, etc.
   - Crear un sistema para componer efectos.
   - Demostrar su uso en una aplicación real.

## Conclusiones

La programación funcional en Python ofrece herramientas poderosas para escribir código más claro, modular y fácil de probar. Aunque Python no es un lenguaje puramente funcional, muchos principios funcionales pueden aplicarse con éxito.

Las técnicas funcionales son especialmente útiles para:

1. **Transformación de datos**: Operaciones en colecciones.
2. **Manejo de errores**: Patrones como Maybe y Result.
3. **Composición de comportamientos**: Funciones de orden superior y composición.

Sin embargo, es importante encontrar un equilibrio entre los principios funcionales y el estilo idiomático de Python. La clave está en usar técnicas funcionales cuando mejoren la claridad y mantenibilidad del código, sin sacrificar la legibilidad y simplicidad que caracterizan a Python.

## Referencias

1. Luciano Ramalho. (2015). Fluent Python. O'Reilly Media.
2. Harry Percival & Bob Gregory. (2020). Architecture Patterns with Python. O'Reilly Media.
3. Brett Slatkin. (2019). Effective Python: 90 Specific Ways to Write Better Python. Addison-Wesley Professional.
4. David Mertz. (2015). Functional Programming in Python. O'Reilly Media.
5. Scott Wlaschin. (2018). Domain Modeling Made Functional: Tackle Software Complexity with Domain-Driven Design and F#.
6. Eric Elliott. (2016). Composing Software: An Exploration of Functional Programming and Object Composition in JavaScript.
7. Python Documentation: functools - https://docs.python.org/3/library/functools.html
8. Python Documentation: itertools - https://docs.python.org/3/library/itertools.html