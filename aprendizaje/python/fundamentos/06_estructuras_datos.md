# Estructuras de Datos en Python

## Introducción

Las estructuras de datos son formas de organizar y almacenar datos para que puedan ser accedidos y modificados eficientemente. Python ofrece varias estructuras de datos incorporadas, cada una con sus propias características y casos de uso. En esta sección, exploraremos las principales estructuras de datos en Python y cómo utilizarlas efectivamente.

## Listas

Las listas son una de las estructuras de datos más versátiles en Python. Son secuencias ordenadas y mutables que pueden contener elementos de diferentes tipos.

### Creación de listas

```python
# Lista vacía
lista_vacia = []

# Lista con elementos
numeros = [1, 2, 3, 4, 5]
mezcla = [1, "hola", 3.14, True, [1, 2]]

# Usando la función list()
caracteres = list("Python")
print(caracteres)  # ['P', 'y', 't', 'h', 'o', 'n']

# Usando comprensión de listas
cuadrados = [x**2 for x in range(1, 6)]
print(cuadrados)  # [1, 4, 9, 16, 25]
```

### Acceso a elementos

```python
frutas = ["manzana", "banana", "cereza", "durazno", "fresa"]

# Acceso por índice (comienza en 0)
print(frutas[0])  # manzana
print(frutas[2])  # cereza

# Índices negativos (cuentan desde el final)
print(frutas[-1])  # fresa (último elemento)
print(frutas[-2])  # durazno (penúltimo elemento)

# Slicing (rebanado)
print(frutas[1:3])    # ['banana', 'cereza'] (del índice 1 al 2)
print(frutas[:3])     # ['manzana', 'banana', 'cereza'] (desde el inicio hasta el índice 2)
print(frutas[2:])     # ['cereza', 'durazno', 'fresa'] (desde el índice 2 hasta el final)
print(frutas[-3:])    # ['cereza', 'durazno', 'fresa'] (los últimos 3 elementos)

# Slicing con paso
print(frutas[::2])    # ['manzana', 'cereza', 'fresa'] (cada 2 elementos)
print(frutas[::-1])   # ['fresa', 'durazno', 'cereza', 'banana', 'manzana'] (invertir lista)
```

### Modificación de listas

```python
frutas = ["manzana", "banana", "cereza"]

# Modificar un elemento
frutas[1] = "kiwi"
print(frutas)  # ['manzana', 'kiwi', 'cereza']

# Añadir elementos
frutas.append("naranja")  # Añade al final
print(frutas)  # ['manzana', 'kiwi', 'cereza', 'naranja']

frutas.insert(1, "pera")  # Inserta en posición específica
print(frutas)  # ['manzana', 'pera', 'kiwi', 'cereza', 'naranja']

frutas.extend(["uva", "melón"])  # Añade múltiples elementos
print(frutas)  # ['manzana', 'pera', 'kiwi', 'cereza', 'naranja', 'uva', 'melón']

# Eliminar elementos
frutas.remove("kiwi")  # Elimina por valor
print(frutas)  # ['manzana', 'pera', 'cereza', 'naranja', 'uva', 'melón']

eliminado = frutas.pop(2)  # Elimina por índice y devuelve el valor
print(eliminado)  # cereza
print(frutas)  # ['manzana', 'pera', 'naranja', 'uva', 'melón']

del frutas[0]  # Elimina por índice
print(frutas)  # ['pera', 'naranja', 'uva', 'melón']

frutas.clear()  # Elimina todos los elementos
print(frutas)  # []
```

### Operaciones comunes con listas

```python
numeros = [3, 1, 4, 1, 5, 9, 2, 6]

# Longitud
print(len(numeros))  # 8

# Suma y multiplicación
print([1, 2, 3] + [4, 5])  # [1, 2, 3, 4, 5] (concatenación)
print([1, 2, 3] * 2)  # [1, 2, 3, 1, 2, 3] (repetición)

# Ordenamiento
numeros.sort()  # Ordena la lista in-place
print(numeros)  # [1, 1, 2, 3, 4, 5, 6, 9]

numeros.sort(reverse=True)  # Ordena en orden descendente
print(numeros)  # [9, 6, 5, 4, 3, 2, 1, 1]

# Ordenamiento sin modificar la lista original
original = [3, 1, 4, 1, 5]
ordenada = sorted(original)
print(original)  # [3, 1, 4, 1, 5] (sin cambios)
print(ordenada)  # [1, 1, 3, 4, 5]

# Invertir
numeros.reverse()  # Invierte la lista in-place
print(numeros)  # [1, 1, 2, 3, 4, 5, 6, 9]

# Contar ocurrencias
print(numeros.count(1))  # 2 (el número 1 aparece dos veces)

# Encontrar índice
print(numeros.index(5))  # 5 (índice del valor 5)

# Verificar pertenencia
print(3 in numeros)  # True
print(7 in numeros)  # False

# Máximo, mínimo y suma
print(max(numeros))  # 9
print(min(numeros))  # 1
print(sum(numeros))  # 31
```

### Listas como pilas y colas

```python
# Pila (LIFO: Last In, First Out)
pila = []
pila.append(1)  # Añadir al final
pila.append(2)
pila.append(3)
print(pila)  # [1, 2, 3]

elemento = pila.pop()  # Sacar del final
print(elemento)  # 3
print(pila)  # [1, 2]

# Cola (FIFO: First In, First Out)
# Para colas eficientes, mejor usar collections.deque
from collections import deque
cola = deque([1, 2, 3])
cola.append(4)  # Añadir al final
print(cola)  # deque([1, 2, 3, 4])

elemento = cola.popleft()  # Sacar del inicio
print(elemento)  # 1
print(cola)  # deque([2, 3, 4])
```

### Comprensión de listas

La comprensión de listas es una forma concisa de crear listas basadas en listas existentes:

```python
# Sintaxis básica: [expresión for elemento in iterable if condición]

# Crear lista de cuadrados
cuadrados = [x**2 for x in range(1, 6)]
print(cuadrados)  # [1, 4, 9, 16, 25]

# Con condición
pares = [x for x in range(1, 11) if x % 2 == 0]
print(pares)  # [2, 4, 6, 8, 10]

# Múltiples for
combinaciones = [(x, y) for x in [1, 2, 3] for y in [3, 1, 4] if x != y]
print(combinaciones)  # [(1, 3), (1, 4), (2, 3), (2, 1), (2, 4), (3, 1), (3, 4)]

# Transformar una matriz
matriz = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
transpuesta = [[fila[i] for fila in matriz] for i in range(3)]
print(transpuesta)  # [[1, 4, 7], [2, 5, 8], [3, 6, 9]]
```

## Tuplas

Las tuplas son secuencias ordenadas e inmutables. Son similares a las listas, pero no se pueden modificar una vez creadas.

### Creación de tuplas

```python
# Tupla vacía
tupla_vacia = ()

# Tupla con elementos
numeros = (1, 2, 3, 4, 5)
mezcla = (1, "hola", 3.14, True, [1, 2])

# Tupla de un solo elemento (requiere una coma)
singleton = (42,)
print(type(singleton))  # <class 'tuple'>

# Sin paréntesis (empaquetado de tupla)
tupla = 1, 2, 3
print(type(tupla))  # <class 'tuple'>

# Usando la función tuple()
caracteres = tuple("Python")
print(caracteres)  # ('P', 'y', 't', 'h', 'o', 'n')
```

### Acceso a elementos

El acceso a elementos en tuplas funciona igual que en listas:

```python
frutas = ("manzana", "banana", "cereza", "durazno")

# Acceso por índice
print(frutas[0])  # manzana

# Índices negativos
print(frutas[-1])  # durazno

# Slicing
print(frutas[1:3])  # ('banana', 'cereza')
```

### Operaciones con tuplas

```python
# Desempaquetado de tupla
tupla = (1, 2, 3)
a, b, c = tupla
print(a, b, c)  # 1 2 3

# Intercambio de variables
x, y = 10, 20
x, y = y, x  # Usando tuplas implícitas
print(x, y)  # 20 10

# Concatenación
tupla1 = (1, 2, 3)
tupla2 = (4, 5, 6)
tupla3 = tupla1 + tupla2
print(tupla3)  # (1, 2, 3, 4, 5, 6)

# Repetición
tupla = (1, 2) * 3
print(tupla)  # (1, 2, 1, 2, 1, 2)

# Longitud, máximo, mínimo
print(len(tupla1))  # 3
print(max(tupla1))  # 3
print(min(tupla1))  # 1

# Verificar pertenencia
print(2 in tupla1)  # True

# Contar ocurrencias
tupla = (1, 2, 2, 3, 2, 4)
print(tupla.count(2))  # 3

# Encontrar índice
print(tupla.index(3))  # 3
```

### Ventajas de las tuplas sobre las listas

1. **Inmutabilidad**: Las tuplas no pueden ser modificadas, lo que las hace más seguras para datos que no deben cambiar.
2. **Rendimiento**: Las tuplas son ligeramente más rápidas que las listas para operaciones de acceso.
3. **Uso como claves de diccionario**: A diferencia de las listas, las tuplas pueden ser usadas como claves en diccionarios (siempre que contengan solo elementos inmutables).
4. **Claridad de intención**: Usar tuplas indica que los datos no deben ser modificados.

## Conjuntos (Sets)

Los conjuntos son colecciones no ordenadas de elementos únicos. Son útiles para eliminar duplicados y realizar operaciones matemáticas de conjuntos.

### Creación de conjuntos

```python
# Conjunto vacío
conjunto_vacio = set()

# Conjunto con elementos
numeros = {1, 2, 3, 4, 5}
frutas = {"manzana", "banana", "cereza"}

# Eliminar duplicados
numeros_con_duplicados = [1, 2, 2, 3, 4, 4, 5]
numeros_unicos = set(numeros_con_duplicados)
print(numeros_unicos)  # {1, 2, 3, 4, 5}

# Comprensión de conjuntos
cuadrados = {x**2 for x in range(1, 6)}
print(cuadrados)  # {1, 4, 9, 16, 25}
```

### Operaciones con conjuntos

```python
# Añadir y eliminar elementos
frutas = {"manzana", "banana", "cereza"}

frutas.add("naranja")  # Añadir un elemento
print(frutas)  # {'manzana', 'banana', 'cereza', 'naranja'}

frutas.update(["uva", "melón"])  # Añadir múltiples elementos
print(frutas)  # {'manzana', 'banana', 'cereza', 'naranja', 'uva', 'melón'}

frutas.remove("banana")  # Eliminar (lanza error si no existe)
print(frutas)  # {'manzana', 'cereza', 'naranja', 'uva', 'melón'}

frutas.discard("kiwi")  # Eliminar si existe (no lanza error)
print(frutas)  # {'manzana', 'cereza', 'naranja', 'uva', 'melón'}

eliminado = frutas.pop()  # Elimina y devuelve un elemento arbitrario
print(eliminado)  # (alguna fruta)
print(frutas)  # (conjunto sin la fruta eliminada)

frutas.clear()  # Elimina todos los elementos
print(frutas)  # set()
```

### Operaciones matemáticas de conjuntos

```python
A = {1, 2, 3, 4, 5}
B = {4, 5, 6, 7, 8}

# Unión
print(A | B)  # {1, 2, 3, 4, 5, 6, 7, 8}
print(A.union(B))  # {1, 2, 3, 4, 5, 6, 7, 8}

# Intersección
print(A & B)  # {4, 5}
print(A.intersection(B))  # {4, 5}

# Diferencia
print(A - B)  # {1, 2, 3}
print(A.difference(B))  # {1, 2, 3}

# Diferencia simétrica
print(A ^ B)  # {1, 2, 3, 6, 7, 8}
print(A.symmetric_difference(B))  # {1, 2, 3, 6, 7, 8}

# Subconjunto y superconjunto
C = {1, 2}
print(C.issubset(A))  # True (C es subconjunto de A)
print(A.issuperset(C))  # True (A es superconjunto de C)
```

### Conjuntos inmutables (frozenset)

Python también ofrece conjuntos inmutables mediante la clase `frozenset`:

```python
conjunto_normal = {1, 2, 3}
conjunto_inmutable = frozenset([1, 2, 3])

# No se puede modificar
# conjunto_inmutable.add(4)  # Esto lanzaría un error

# Pero se puede usar en operaciones que no modifican
union = conjunto_inmutable.union({4, 5})
print(union)  # frozenset({1, 2, 3, 4, 5})

# Los frozenset pueden ser usados como claves de diccionario
diccionario = {conjunto_inmutable: "valor"}
print(diccionario[conjunto_inmutable])  # valor
```

## Diccionarios

Los diccionarios son colecciones no ordenadas de pares clave-valor. Son extremadamente útiles para almacenar y recuperar datos asociados.

### Creación de diccionarios

```python
# Diccionario vacío
diccionario_vacio = {}

# Diccionario con elementos
persona = {"nombre": "Ana", "edad": 30, "profesion": "ingeniera"}

# Usando dict()
colores = dict(rojo="#FF0000", verde="#00FF00", azul="#0000FF")
print(colores)  # {'rojo': '#FF0000', 'verde': '#00FF00', 'azul': '#0000FF'}

# A partir de una lista de tuplas
items = [("a", 1), ("b", 2), ("c", 3)]
diccionario = dict(items)
print(diccionario)  # {'a': 1, 'b': 2, 'c': 3}

# Comprensión de diccionarios
cuadrados = {x: x**2 for x in range(1, 6)}
print(cuadrados)  # {1: 1, 2: 4, 3: 9, 4: 16, 5: 25}
```

### Acceso y modificación

```python
persona = {"nombre": "Ana", "edad": 30, "profesion": "ingeniera"}

# Acceso por clave
print(persona["nombre"])  # Ana

# Acceso seguro con get (devuelve None o un valor por defecto si la clave no existe)
print(persona.get("altura"))  # None
print(persona.get("altura", 165))  # 165

# Modificar valores
persona["edad"] = 31
print(persona)  # {'nombre': 'Ana', 'edad': 31, 'profesion': 'ingeniera'}

# Añadir nuevos pares clave-valor
persona["ciudad"] = "Madrid"
print(persona)  # {'nombre': 'Ana', 'edad': 31, 'profesion': 'ingeniera', 'ciudad': 'Madrid'}

# Actualizar múltiples valores
persona.update({"edad": 32, "estado_civil": "casada"})
print(persona)  # {'nombre': 'Ana', 'edad': 32, 'profesion': 'ingeniera', 'ciudad': 'Madrid', 'estado_civil': 'casada'}

# Eliminar pares clave-valor
del persona["estado_civil"]
print(persona)  # {'nombre': 'Ana', 'edad': 32, 'profesion': 'ingeniera', 'ciudad': 'Madrid'}

eliminado = persona.pop("ciudad")
print(eliminado)  # Madrid
print(persona)  # {'nombre': 'Ana', 'edad': 32, 'profesion': 'ingeniera'}

# Eliminar y devolver el último par insertado (en Python 3.7+ el orden de inserción se preserva)
ultimo = persona.popitem()
print(ultimo)  # ('profesion', 'ingeniera')
print(persona)  # {'nombre': 'Ana', 'edad': 32}

# Limpiar el diccionario
persona.clear()
print(persona)  # {}
```

### Operaciones comunes con diccionarios

```python
persona = {"nombre": "Ana", "edad": 30, "profesion": "ingeniera"}

# Verificar si una clave existe
print("nombre" in persona)  # True
print("altura" in persona)  # False

# Obtener todas las claves
claves = persona.keys()
print(claves)  # dict_keys(['nombre', 'edad', 'profesion'])

# Obtener todos los valores
valores = persona.values()
print(valores)  # dict_values(['Ana', 30, 'ingeniera'])

# Obtener todos los pares clave-valor
items = persona.items()
print(items)  # dict_items([('nombre', 'Ana'), ('edad', 30), ('profesion', 'ingeniera')])

# Iterar sobre un diccionario
for clave in persona:
    print(clave, persona[clave])

# Iterar sobre pares clave-valor
for clave, valor in persona.items():
    print(clave, valor)

# Crear una copia
copia = persona.copy()
print(copia)  # {'nombre': 'Ana', 'edad': 30, 'profesion': 'ingeniera'}

# Diccionario con valor por defecto
from collections import defaultdict

# defaultdict con valor por defecto 0
contador = defaultdict(int)
palabras = ["manzana", "banana", "manzana", "cereza", "banana", "manzana"]
for palabra in palabras:
    contador[palabra] += 1
print(dict(contador))  # {'manzana': 3, 'banana': 2, 'cereza': 1}

# defaultdict con valor por defecto lista vacía
grupos = defaultdict(list)
personas = [("Grupo A", "Ana"), ("Grupo B", "Carlos"), ("Grupo A", "David")]
for grupo, persona in personas:
    grupos[grupo].append(persona)
print(dict(grupos))  # {'Grupo A': ['Ana', 'David'], 'Grupo B': ['Carlos']}
```

### Diccionarios anidados

```python
# Diccionario con diccionarios anidados
equipo = {
    "gerente": {
        "nombre": "Ana",
        "edad": 45,
        "departamento": "Ventas"
    },
    "desarrollador": {
        "nombre": "Carlos",
        "edad": 30,
        "departamento": "IT"
    }
}

# Acceso a valores anidados
print(equipo["gerente"]["nombre"])  # Ana

# Modificar valores anidados
equipo["desarrollador"]["edad"] = 31

# Añadir un nuevo empleado
equipo["diseñador"] = {"nombre": "Elena", "edad": 28, "departamento": "Diseño"}

# Iterar sobre un diccionario anidado
for puesto, datos in equipo.items():
    print(f"{puesto.capitalize()}: {datos['nombre']}, {datos['edad']} años, {datos['departamento']}")
```

## Otras Estructuras de Datos

### Colas de doble extremo (deque)

```python
from collections import deque

# Crear una deque
cola = deque(["a", "b", "c"])

# Añadir elementos
cola.append("d")        # Añade a la derecha
cola.appendleft("z")    # Añade a la izquierda
print(cola)  # deque(['z', 'a', 'b', 'c', 'd'])

# Eliminar elementos
cola.pop()              # Elimina de la derecha
cola.popleft()          # Elimina de la izquierda
print(cola)  # deque(['a', 'b', 'c'])

# Rotar elementos
cola.rotate(1)          # Rota 1 posición a la derecha
print(cola)  # deque(['c', 'a', 'b'])
cola.rotate(-2)         # Rota 2 posiciones a la izquierda
print(cola)  # deque(['b', 'c', 'a'])
```

### Contadores

```python
from collections import Counter

# Contar elementos
palabras = ["manzana", "banana", "manzana", "cereza", "banana", "manzana"]
contador = Counter(palabras)
print(contador)  # Counter({'manzana': 3, 'banana': 2, 'cereza': 1})

# Elementos más comunes
print(contador.most_common(2))  # [('manzana', 3), ('banana', 2)]

# Operaciones
c1 = Counter(a=3, b=1)
c2 = Counter(a=1, b=2)
print(c1 + c2)  # Counter({'a': 4, 'b': 3})
print(c1 - c2)  # Counter({'a': 2})
```

### OrderedDict

```python
from collections import OrderedDict

# En Python 3.7+, los diccionarios normales mantienen el orden de inserción,
# pero OrderedDict sigue siendo útil para algunas operaciones específicas

ord_dict = OrderedDict()
ord_dict["a"] = 1
ord_dict["b"] = 2
ord_dict["c"] = 3

# Mover un elemento al final
ord_dict.move_to_end("a")
print(list(ord_dict.items()))  # [('b', 2), ('c', 3), ('a', 1)]

# Eliminar y devolver el primer elemento insertado
primer_elemento = ord_dict.popitem(last=False)
print(primer_elemento)  # ('b', 2)
print(ord_dict)  # OrderedDict([('c', 3), ('a', 1)])
```

### Colas de prioridad (heapq)

```python
import heapq

# Crear una cola de prioridad (heap)
numeros = [3, 1, 4, 1, 5, 9, 2, 6]
heapq.heapify(numeros)  # Convierte la lista en un heap
print(numeros)  # [1, 1, 2, 3, 5, 9, 4, 6] (representación interna del heap)

# Extraer el elemento más pequeño
print(heapq.heappop(numeros))  # 1
print(numeros)  # [1, 3, 2, 6, 5, 9, 4]

# Añadir un elemento
heapq.heappush(numeros, 0)
print(numeros)  # [0, 1, 2, 6, 3, 9, 4, 5]

# Los n elementos más pequeños
print(heapq.nsmallest(3, [3, 1, 4, 1, 5, 9, 2, 6]))  # [1, 1, 2]

# Los n elementos más grandes
print(heapq.nlargest(3, [3, 1, 4, 1, 5, 9, 2, 6]))  # [9, 6, 5]
```

### Pilas (usando listas)

```python
# Implementación de una pila usando una lista
class Pila:
    def __init__(self):
        self.items = []
    
    def esta_vacia(self):
        return len(self.items) == 0
    
    def apilar(self, item):
        self.items.append(item)
    
    def desapilar(self):
        if not self.esta_vacia():
            return self.items.pop()
        return None
    
    def ver_tope(self):
        if not self.esta_vacia():
            return self.items[-1]
        return None
    
    def tamaño(self):
        return len(self.items)

# Uso de la pila
pila = Pila()
pila.apilar(1)
pila.apilar(2)
pila.apilar(3)

print(pila.ver_tope())  # 3
print(pila.desapilar())  # 3
print(pila.tamaño())    # 2
```

## Ejemplos Prácticos

### Contador de palabras

```python
from collections import Counter
import re

def contar_palabras(texto):
    # Convertir a minúsculas y eliminar caracteres no alfanuméricos
    texto = texto.lower()
    palabras = re.findall(r'\b\w+\b', texto)
    
    # Contar palabras
    contador = Counter(palabras)
    return contador

# Ejemplo de uso
texto = """Python es un lenguaje de programación interpretado cuya filosofía hace 
        hincapié en la legibilidad de su código. Es un lenguaje interpretado, 
        dinámico y multiplataforma."""

resultado = contar_palabras(texto)
print("Palabras más comunes:")
for palabra, frecuencia in resultado.most_common(5):
    print(f"{palabra}: {frecuencia}")
```

### Agenda de contactos

```python
class Agenda:
    def __init__(self):
        self.contactos = {}
    
    def agregar_contacto(self, nombre, telefono, email="", direccion=""):
        self.contactos[nombre] = {
            "telefono": telefono,
            "email": email,
            "direccion": direccion
        }
        print(f"Contacto {nombre} agregado correctamente.")
    
    def eliminar_contacto(self, nombre):
        if nombre in self.contactos:
            del self.contactos[nombre]
            print(f"Contacto {nombre} eliminado correctamente.")
        else:
            print(f"No se encontró el contacto {nombre}.")
    
    def buscar_contacto(self, nombre):
        if nombre in self.contactos:
            print(f"Información de {nombre}:")
            for clave, valor in self.contactos[nombre].items():
                if valor:  # Solo mostrar campos con valor
                    print(f"  {clave.capitalize()}: {valor}")
        else:
            print(f"No se encontró el contacto {nombre}.")
    
    def listar_contactos(self):
        if not self.contactos:
            print("La agenda está vacía.")
            return
        
        print("Lista de contactos:")
        for i, (nombre, datos) in enumerate(sorted(self.contactos.items()), 1):
            print(f"{i}. {nombre}: {datos['telefono']}")
    
    def actualizar_contacto(self, nombre, **kwargs):
        if nombre in self.contactos:
            for clave, valor in kwargs.items():
                if clave in self.contactos[nombre]:
                    self.contactos[nombre][clave] = valor
            print(f"Contacto {nombre} actualizado correctamente.")
        else:
            print(f"No se encontró el contacto {nombre}.")

# Ejemplo de uso
agenda = Agenda()

# Agregar contactos
agenda.agregar_contacto("Ana García", "555-1234", "ana@example.com", "Calle Principal 123")
agenda.agregar_contacto("Carlos López", "555-5678", "carlos@example.com")
agenda.agregar_contacto("Elena Martínez", "555-9012")

# Listar contactos
agenda.listar_contactos()

# Buscar contacto
agenda.buscar_contacto("Ana García")

# Actualizar contacto
agenda.actualizar_contacto("Elena Martínez", telefono="555-3456", email="elena@example.com")
agenda.buscar_contacto("Elena Martínez")

# Eliminar contacto
agenda.eliminar_contacto("Carlos López")
agenda.listar_contactos()
```

### Análisis de frecuencia de caracteres

```python
def analizar_frecuencia(texto):
    # Eliminar espacios y convertir a minúsculas
    texto = texto.lower().replace(" ", "")
    
    # Contar frecuencia de cada carácter
    frecuencia = {}
    for caracter in texto:
        if caracter in frecuencia:
            frecuencia[caracter] += 1
        else:
            frecuencia[caracter] = 1
    
    # Calcular porcentajes
    total = len(texto)
    porcentajes = {caracter: (count / total) * 100 for caracter, count in frecuencia.items()}
    
    # Ordenar por frecuencia descendente
    ordenados = sorted(porcentajes.items(), key=lambda x: x[1], reverse=True)
    
    return ordenados

# Ejemplo de uso
texto = "Python es un lenguaje de programación de alto nivel"
resultado = analizar_frecuencia(texto)

print("Análisis de frecuencia de caracteres:")
print("Carácter | Frecuencia (%)")
print("-" * 25)
for caracter, porcentaje in resultado:
    print(f"   {caracter}    | {porcentaje:.2f}%")
```

## Buenas Prácticas

1. **Elegir la estructura adecuada**:
   - Usa listas cuando necesites una colección ordenada y mutable.
   - Usa tuplas para datos inmutables o cuando quieras asegurar que los datos no cambien.
   - Usa conjuntos para eliminar duplicados o realizar operaciones de conjuntos.
   - Usa diccionarios cuando necesites asociar valores con claves.

2. **Comprensiones**:
   - Usa comprensiones de listas, diccionarios y conjuntos para código más conciso y legible.
   - Pero evita comprensiones demasiado complejas que dificulten la legibilidad.

3. **Inmutabilidad**:
   - Prefiere estructuras inmutables (tuplas, frozenset) cuando los datos no deban cambiar.
   - Esto evita efectos secundarios no deseados y permite usar estos objetos como claves de diccionarios.

4. **Eficiencia**:
   - Para operaciones frecuentes de inserción/eliminación al principio de una secuencia, usa `collections.deque` en lugar de listas.
   - Para contar elementos, usa `collections.Counter` en lugar de diccionarios manuales.
   - Para colas de prioridad, usa el módulo `heapq`.

5. **Verificación de pertenencia**:
   - Para verificar si un elemento está en una colección, usa el operador `in`.
   - Esta operación es más eficiente en conjuntos y diccionarios (O(1)) que en listas y tuplas (O(n)).

6. **Copias**:
   - Ten cuidado con las copias superficiales vs. profundas.
   - Usa `.copy()` para copias superficiales y `copy.deepcopy()` para copias profundas.

7. **Acceso seguro a diccionarios**:
   - Usa `get()` o `setdefault()` para acceder a valores de diccionarios de forma segura.
   - Considera usar `defaultdict` para valores por defecto.

8. **Nombres descriptivos**:
   - Usa nombres descriptivos para tus variables que indiquen qué tipo de estructura contienen.

## Recursos Adicionales

- [Documentación oficial de Python - Estructuras de datos](https://docs.python.org/es/3/tutorial/datastructures.html)
- [Python Collections Module](https://docs.python.org/3/library/collections.html)
- [Real Python - Listas y tuplas](https://realpython.com/python-lists-tuples/)
- [Real Python - Diccionarios](https://realpython.com/python-dicts/)
- [Real Python - Sets](https://realpython.com/python-sets/)
- [Python Cookbook - Recetas para estructuras de datos](https://github.com/dabeaz/python-cookbook)

---

En la siguiente sección, exploraremos la programación orientada a objetos en Python, que nos permitirá crear estructuras de datos personalizadas y organizar nuestro código de manera más modular y reutilizable.