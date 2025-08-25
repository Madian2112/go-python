# Programación Funcional Avanzada en Python

## Introducción

La programación funcional es un paradigma de programación que trata la computación como la evaluación de funciones matemáticas y evita cambiar el estado y los datos mutables. Python, aunque no es un lenguaje puramente funcional como Haskell o Lisp, ofrece muchas características que permiten adoptar un estilo de programación funcional.

En este módulo, exploraremos técnicas avanzadas de programación funcional en Python, que te permitirán escribir código más conciso, mantenible y menos propenso a errores.

## Funciones de Orden Superior

Las funciones de orden superior son funciones que pueden tomar otras funciones como argumentos o devolver funciones como resultado.

### map, filter y reduce

```python
# map: aplica una función a cada elemento de un iterable
numeros = [1, 2, 3, 4, 5]
cuadrados = list(map(lambda x: x**2, numeros))  # [1, 4, 9, 16, 25]

# filter: filtra elementos según una función de predicado
pares = list(filter(lambda x: x % 2 == 0, numeros))  # [2, 4]

# reduce: reduce un iterable a un solo valor aplicando una función acumulativa
from functools import reduce
suma = reduce(lambda x, y: x + y, numeros)  # 15
```

### Funciones como Argumentos

```python
def aplicar_operacion(func, x, y):
    return func(x, y)

def suma(x, y):
    return x + y

def multiplicacion(x, y):
    return x * y

print(aplicar_operacion(suma, 5, 3))  # 8
print(aplicar_operacion(multiplicacion, 5, 3))  # 15
```

### Funciones que Devuelven Funciones

```python
def crear_multiplicador(factor):
    def multiplicador(x):
        return x * factor
    return multiplicador

duplica = crear_multiplicador(2)
triplica = crear_multiplicador(3)

print(duplica(5))  # 10
print(triplica(5))  # 15
```

## Funciones Lambda

Las funciones lambda son funciones anónimas de una sola expresión.

```python
# Función lambda básica
cuadrado = lambda x: x**2
print(cuadrado(5))  # 25

# Uso con sorted
personas = [('Alice', 25), ('Bob', 30), ('Charlie', 20)]
personas_ordenadas = sorted(personas, key=lambda persona: persona[1])
print(personas_ordenadas)  # [('Charlie', 20), ('Alice', 25), ('Bob', 30)]

# Uso con map y filter
numeros = [1, 2, 3, 4, 5]
cuadrados_pares = list(map(lambda x: x**2, filter(lambda x: x % 2 == 0, numeros)))
print(cuadrados_pares)  # [4, 16]
```

## Clausuras (Closures)

Una clausura es una función que recuerda los valores del ámbito en el que fue creada, incluso si ese ámbito ya no está disponible.

```python
def contador():
    count = 0
    def incrementar():
        nonlocal count
        count += 1
        return count
    return incrementar

mi_contador = contador()
print(mi_contador())  # 1
print(mi_contador())  # 2
print(mi_contador())  # 3

# Cada instancia de contador mantiene su propio estado
otro_contador = contador()
print(otro_contador())  # 1
print(mi_contador())  # 4
```

## Decoradores

Los decoradores son una forma poderosa de modificar o extender el comportamiento de funciones o métodos.

### Decoradores Básicos

```python
def mi_decorador(func):
    def wrapper(*args, **kwargs):
        print("Antes de llamar a la función")
        resultado = func(*args, **kwargs)
        print("Después de llamar a la función")
        return resultado
    return wrapper

@mi_decorador
def saludar(nombre):
    print(f"Hola, {nombre}!")

saludar("Alice")
# Salida:
# Antes de llamar a la función
# Hola, Alice!
# Después de llamar a la función
```

### Decoradores con Argumentos

```python
def repetir(n):
    def decorador(func):
        def wrapper(*args, **kwargs):
            for _ in range(n):
                resultado = func(*args, **kwargs)
            return resultado
        return wrapper
    return decorador

@repetir(3)
def saludar(nombre):
    print(f"Hola, {nombre}!")

saludar("Bob")
# Salida:
# Hola, Bob!
# Hola, Bob!
# Hola, Bob!
```

### Decoradores de Clase

```python
class ConContador:
    def __init__(self, func):
        self.func = func
        self.count = 0
        
    def __call__(self, *args, **kwargs):
        self.count += 1
        return self.func(*args, **kwargs)

@ConContador
def mi_funcion():
    return "Función llamada"

print(mi_funcion())  # Función llamada
print(mi_funcion())  # Función llamada
print(mi_funcion.count)  # 2
```

## Funciones Parciales

Las funciones parciales permiten fijar un número de argumentos de una función y generar una nueva función.

```python
from functools import partial

def potencia(base, exponente):
    return base ** exponente

cuadrado = partial(potencia, exponente=2)
cubo = partial(potencia, exponente=3)

print(cuadrado(5))  # 25
print(cubo(5))  # 125
```

## Módulo functools

El módulo `functools` proporciona herramientas de alto orden para trabajar con funciones y objetos invocables.

### lru_cache

```python
from functools import lru_cache

@lru_cache(maxsize=None)
def fibonacci(n):
    if n < 2:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

# Sin caché, esto sería extremadamente lento
print(fibonacci(100))  # 354224848179261915075
```

### singledispatch

```python
from functools import singledispatch

@singledispatch
def procesar(obj):
    raise NotImplementedError(f"No se puede procesar objeto de tipo {type(obj)}")

@procesar.register
def _(obj: int):
    return f"Procesando entero: {obj}"

@procesar.register
def _(obj: str):
    return f"Procesando cadena: {obj}"

@procesar.register(list)
def _(obj):
    return f"Procesando lista con {len(obj)} elementos"

print(procesar(10))  # Procesando entero: 10
print(procesar("hola"))  # Procesando cadena: hola
print(procesar([1, 2, 3]))  # Procesando lista con 3 elementos
```

## Comprensiones y Expresiones Generadoras

Las comprensiones y expresiones generadoras proporcionan una sintaxis concisa para crear listas, diccionarios, conjuntos y generadores.

### Comprensiones Avanzadas

```python
# Comprensión de lista con múltiples condiciones
numeros = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
pares_mayores_que_5 = [x for x in numeros if x % 2 == 0 if x > 5]
print(pares_mayores_que_5)  # [6, 8, 10]

# Comprensión de lista con if-else
paridad = ["par" if x % 2 == 0 else "impar" for x in numeros]
print(paridad)  # ['impar', 'par', 'impar', 'par', 'impar', 'par', 'impar', 'par', 'impar', 'par']

# Comprensión de diccionario
cuadrados_dict = {x: x**2 for x in numeros}
print(cuadrados_dict)  # {1: 1, 2: 4, 3: 9, 4: 16, 5: 25, 6: 36, 7: 49, 8: 64, 9: 81, 10: 100}

# Comprensión de conjunto
cuadrados_set = {x**2 for x in numeros}
print(cuadrados_set)  # {64, 1, 4, 36, 100, 9, 16, 49, 81, 25}
```

### Expresiones Generadoras

```python
# Expresión generadora
generador = (x**2 for x in range(1, 6))
print(next(generador))  # 1
print(next(generador))  # 4

# Uso en funciones que aceptan iterables
suma_cuadrados = sum(x**2 for x in range(1, 6))
print(suma_cuadrados)  # 55
```

## Inmutabilidad y Programación Funcional Pura

La programación funcional pura evita los efectos secundarios y favorece la inmutabilidad.

```python
# Función impura (con efecto secundario)
def agregar_elemento_impuro(lista, elemento):
    lista.append(elemento)  # Modifica la lista original
    return lista

# Función pura (sin efectos secundarios)
def agregar_elemento_puro(lista, elemento):
    return lista + [elemento]  # Crea una nueva lista

lista_original = [1, 2, 3]

# Usando la función impura
lista_modificada = agregar_elemento_impuro(lista_original, 4)
print(lista_original)  # [1, 2, 3, 4] - La lista original fue modificada
print(lista_modificada)  # [1, 2, 3, 4] - Ambas variables apuntan a la misma lista

# Reiniciamos la lista original
lista_original = [1, 2, 3]

# Usando la función pura
lista_nueva = agregar_elemento_puro(lista_original, 4)
print(lista_original)  # [1, 2, 3] - La lista original no fue modificada
print(lista_nueva)  # [1, 2, 3, 4] - Se creó una nueva lista
```

## Módulo itertools

El módulo `itertools` proporciona funciones para trabajar con iteradores de manera eficiente.

```python
import itertools

# Combinaciones
print(list(itertools.combinations('ABC', 2)))  # [('A', 'B'), ('A', 'C'), ('B', 'C')]

# Permutaciones
print(list(itertools.permutations('ABC', 2)))  # [('A', 'B'), ('A', 'C'), ('B', 'A'), ('B', 'C'), ('C', 'A'), ('C', 'B')]

# Producto cartesiano
print(list(itertools.product('AB', '12')))  # [('A', '1'), ('A', '2'), ('B', '1'), ('B', '2')]

# Ciclo infinito
ciclo = itertools.cycle([1, 2, 3])
print([next(ciclo) for _ in range(6)])  # [1, 2, 3, 1, 2, 3]

# Agrupación
animales = ['perro', 'pato', 'gato', 'pez', 'gallina']
for letra, grupo in itertools.groupby(sorted(animales), key=lambda x: x[0]):
    print(letra, list(grupo))
# g ['gallina', 'gato']
# p ['pato', 'perro', 'pez']
```

## Ejercicios Prácticos

1. Implementa una función `compose` que tome múltiples funciones y devuelva una nueva función que sea la composición de todas ellas.

2. Crea un decorador que mida el tiempo de ejecución de una función.

3. Implementa un sistema de validación de datos usando funciones parciales.

4. Crea un pipeline de procesamiento de datos usando programación funcional.

## Conclusión

La programación funcional en Python ofrece herramientas poderosas para escribir código más limpio, modular y mantenible. Aunque Python no es un lenguaje puramente funcional, incorporar estos principios y técnicas puede mejorar significativamente la calidad de tu código.

Recuerda que no es necesario adoptar un enfoque puramente funcional; a menudo, el mejor enfoque es combinar paradigmas según las necesidades específicas de tu proyecto. La programación funcional es especialmente útil para operaciones de transformación de datos, procesamiento en paralelo y creación de código con menos efectos secundarios y, por lo tanto, más fácil de probar y depurar.