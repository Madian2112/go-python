# Funciones en Python

## Introducción

Las funciones son bloques de código reutilizables que realizan una tarea específica. Permiten organizar el código, evitar la repetición y hacer que los programas sean más modulares y mantenibles. Python proporciona una sintaxis flexible y potente para definir y utilizar funciones.

## Definición de Funciones

En Python, las funciones se definen con la palabra clave `def`, seguida del nombre de la función y paréntesis que pueden contener parámetros:

```python
def nombre_funcion(parametro1, parametro2, ...):
    # Cuerpo de la función
    # Código a ejecutar
    return valor  # Opcional
```

Ejemplo simple:

```python
def saludar(nombre):
    """Esta función saluda a la persona pasada como parámetro."""
    return f"¡Hola, {nombre}!"

# Llamada a la función
mensaje = saludar("Ana")
print(mensaje)  # Imprime: ¡Hola, Ana!
```

### Docstrings

Los docstrings son cadenas de documentación que describen lo que hace una función. Se colocan justo después de la definición de la función y se delimitan con triples comillas:

```python
def area_rectangulo(base, altura):
    """Calcula el área de un rectángulo.
    
    Args:
        base (float): La base del rectángulo.
        altura (float): La altura del rectángulo.
        
    Returns:
        float: El área del rectángulo.
    """
    return base * altura
```

Los docstrings son accesibles mediante el atributo `__doc__` o la función `help()`:

```python
print(area_rectangulo.__doc__)
help(area_rectangulo)
```

## Parámetros y Argumentos

### Parámetros posicionales

Los parámetros posicionales son los más básicos y se asignan según el orden en que se pasan los argumentos:

```python
def presentar(nombre, edad):
    return f"{nombre} tiene {edad} años."

# Los argumentos se asignan por posición
presentar("Carlos", 30)  # "Carlos tiene 30 años."
```

### Parámetros con valores predeterminados

Puedes asignar valores predeterminados a los parámetros, que se utilizarán si no se proporciona un argumento:

```python
def saludar(nombre, mensaje="Hola"):
    return f"{mensaje}, {nombre}!"

saludar("Ana")  # "Hola, Ana!"
saludar("Ana", "Buenos días")  # "Buenos días, Ana!"
```

**Importante**: Los parámetros con valores predeterminados deben aparecer después de los parámetros sin valores predeterminados.

### Argumentos por nombre

Puedes pasar argumentos especificando el nombre del parámetro, lo que permite cambiar el orden:

```python
def dividir(dividendo, divisor):
    return dividendo / divisor

# Usando argumentos por nombre
resultado = dividir(divisor=2, dividendo=10)  # 5.0
```

### Número variable de argumentos posicionales (*args)

Puedes definir funciones que acepten un número variable de argumentos posicionales usando `*args`:

```python
def sumar(*numeros):
    """Suma todos los números proporcionados."""
    total = 0
    for numero in numeros:
        total += numero
    return total

sumar(1, 2)  # 3
sumar(1, 2, 3, 4, 5)  # 15
```

### Número variable de argumentos por nombre (**kwargs)

Puedes definir funciones que acepten un número variable de argumentos por nombre usando `**kwargs`:

```python
def mostrar_informacion(**datos):
    """Muestra la información proporcionada."""
    for clave, valor in datos.items():
        print(f"{clave}: {valor}")

mostrar_informacion(nombre="Ana", edad=30, ciudad="Madrid")
# Imprime:
# nombre: Ana
# edad: 30
# ciudad: Madrid
```

### Combinando diferentes tipos de parámetros

Puedes combinar todos los tipos de parámetros, pero deben seguir este orden:

1. Parámetros posicionales
2. Parámetros con valores predeterminados
3. `*args` (argumentos posicionales variables)
4. `**kwargs` (argumentos por nombre variables)

```python
def funcion_completa(pos1, pos2, default1="valor1", default2="valor2", *args, **kwargs):
    print(f"Posicionales: {pos1}, {pos2}")
    print(f"Predeterminados: {default1}, {default2}")
    print(f"Args: {args}")
    print(f"Kwargs: {kwargs}")

funcion_completa(1, 2, "personalizado", *[3, 4, 5], extra="valor", otro=42)
```

## Valores de Retorno

### Retorno simple

Las funciones pueden devolver un valor usando la instrucción `return`:

```python
def cuadrado(x):
    return x * x

resultado = cuadrado(5)  # 25
```

### Retorno múltiple

Las funciones pueden devolver múltiples valores separados por comas, que se empaquetan automáticamente en una tupla:

```python
def estadisticas(numeros):
    return min(numeros), max(numeros), sum(numeros) / len(numeros)

minimo, maximo, promedio = estadisticas([1, 2, 3, 4, 5])
print(f"Mínimo: {minimo}, Máximo: {maximo}, Promedio: {promedio}")
```

### Retorno anticipado

La instrucción `return` termina inmediatamente la ejecución de la función y devuelve el valor especificado:

```python
def es_par(numero):
    if numero % 2 == 0:
        return True
    return False  # Solo se ejecuta si el número no es par
```

### Funciones sin retorno

Si una función no incluye una instrucción `return` o incluye `return` sin un valor, devuelve `None`:

```python
def saludar(nombre):
    print(f"Hola, {nombre}!")

resultado = saludar("Ana")  # La función imprime "Hola, Ana!" y devuelve None
print(resultado)  # None
```

## Alcance de Variables

### Variables locales

Las variables definidas dentro de una función solo son accesibles dentro de esa función:

```python
def funcion():
    x = 10  # Variable local
    print(x)

funcion()  # Imprime: 10
# print(x)  # Error: x no está definida fuera de la función
```

### Variables globales

Las variables definidas fuera de cualquier función son globales y pueden ser accedidas (pero no modificadas) dentro de las funciones:

```python
x = 10  # Variable global

def funcion():
    print(x)  # Accede a la variable global

funcion()  # Imprime: 10
```

### Modificación de variables globales

Para modificar una variable global dentro de una función, debes usar la palabra clave `global`:

```python
x = 10  # Variable global

def modificar():
    global x
    x = 20  # Modifica la variable global

print(x)  # 10
modificar()
print(x)  # 20
```

### Variables no locales

En funciones anidadas, puedes acceder a variables de la función externa usando la palabra clave `nonlocal`:

```python
def externa():
    x = 10
    
    def interna():
        nonlocal x
        x = 20  # Modifica la variable de la función externa
    
    print(x)  # 10
    interna()
    print(x)  # 20

externa()
```

## Funciones como Objetos de Primera Clase

En Python, las funciones son objetos de primera clase, lo que significa que pueden ser:

### Asignadas a variables

```python
def saludar(nombre):
    return f"Hola, {nombre}!"

# Asignar la función a una variable
mi_funcion = saludar

# Llamar a la función a través de la variable
print(mi_funcion("Ana"))  # Imprime: Hola, Ana!
```

### Pasadas como argumentos

```python
def saludar(nombre):
    return f"Hola, {nombre}!"

def despedir(nombre):
    return f"Adiós, {nombre}!"

def ejecutar_funcion(funcion, argumento):
    return funcion(argumento)

print(ejecutar_funcion(saludar, "Ana"))  # Imprime: Hola, Ana!
print(ejecutar_funcion(despedir, "Ana"))  # Imprime: Adiós, Ana!
```

### Devueltas por otras funciones

```python
def crear_multiplicador(factor):
    def multiplicador(numero):
        return numero * factor
    return multiplicador

duplicar = crear_multiplicador(2)
triplicar = crear_multiplicador(3)

print(duplicar(5))  # 10
print(triplicar(5))  # 15
```

## Funciones Lambda

Las funciones lambda son funciones anónimas (sin nombre) que se definen con la palabra clave `lambda`. Son útiles para funciones simples de una sola expresión:

```python
# Sintaxis: lambda argumentos: expresión

# Función lambda que suma dos números
suma = lambda x, y: x + y
print(suma(5, 3))  # 8

# Función lambda como argumento
numeros = [1, 5, 3, 9, 2, 6]
numeros_ordenados = sorted(numeros, key=lambda x: x % 3)  # Ordena por el resto de dividir por 3
print(numeros_ordenados)  # [3, 6, 9, 1, 4, 7]
```

Limitaciones de las funciones lambda:
- Solo pueden contener una expresión
- No pueden contener instrucciones (como `if`, `for`, etc.)
- No pueden tener docstrings

## Funciones Integradas

Python proporciona muchas funciones integradas útiles:

| Función | Descripción | Ejemplo |
|---------|-------------|--------|
| `print()` | Imprime objetos en la consola | `print("Hola")` |
| `len()` | Devuelve la longitud de un objeto | `len([1, 2, 3])` |
| `type()` | Devuelve el tipo de un objeto | `type(5)` |
| `int()`, `float()`, `str()` | Conversión de tipos | `int("5")`, `float(5)`, `str(5)` |
| `min()`, `max()` | Valor mínimo/máximo | `min(1, 2, 3)`, `max([1, 2, 3])` |
| `sum()` | Suma los elementos de un iterable | `sum([1, 2, 3])` |
| `sorted()` | Devuelve una lista ordenada | `sorted([3, 1, 2])` |
| `enumerate()` | Devuelve pares (índice, valor) | `list(enumerate(['a', 'b', 'c']))` |
| `zip()` | Combina iterables | `list(zip([1, 2], ['a', 'b']))` |
| `map()` | Aplica una función a cada elemento | `list(map(lambda x: x*2, [1, 2, 3]))` |
| `filter()` | Filtra elementos según una función | `list(filter(lambda x: x>0, [-1, 0, 1]))` |

## Decoradores

Los decoradores son funciones que modifican el comportamiento de otras funciones. Se aplican usando el símbolo `@` seguido del nombre del decorador encima de la definición de la función:

```python
def mi_decorador(funcion):
    def wrapper(*args, **kwargs):
        print("Antes de llamar a la función")
        resultado = funcion(*args, **kwargs)
        print("Después de llamar a la función")
        return resultado
    return wrapper

@mi_decorador
def saludar(nombre):
    print(f"Hola, {nombre}!")

saludar("Ana")
# Imprime:
# Antes de llamar a la función
# Hola, Ana!
# Después de llamar a la función
```

Los decoradores son útiles para:
- Registrar llamadas a funciones (logging)
- Medir el tiempo de ejecución
- Verificar precondiciones
- Cachear resultados
- Gestionar permisos o autenticación

## Generadores

Los generadores son funciones especiales que devuelven un iterador. Se definen como funciones normales pero usan `yield` en lugar de `return`:

```python
def contar_hasta(n):
    i = 1
    while i <= n:
        yield i
        i += 1

# Usar el generador
for numero in contar_hasta(5):
    print(numero)  # Imprime: 1, 2, 3, 4, 5
```

Ventajas de los generadores:
- Eficiencia de memoria (generan valores bajo demanda)
- Pueden representar secuencias infinitas
- Simplifican el código para secuencias complejas

## Recursión

La recursión es cuando una función se llama a sí misma. Es útil para problemas que pueden dividirse en subproblemas más pequeños del mismo tipo:

```python
def factorial(n):
    """Calcula el factorial de n de forma recursiva."""
    if n == 0 or n == 1:
        return 1
    else:
        return n * factorial(n - 1)

print(factorial(5))  # 120 (5 * 4 * 3 * 2 * 1)
```

Consideraciones sobre la recursión:
- Debe tener un caso base para evitar recursión infinita
- Puede consumir mucha memoria para valores grandes (límite de recursión)
- A veces es menos eficiente que soluciones iterativas
- Python tiene un límite de recursión (normalmente 1000)

## Ejemplos Prácticos

### Calculadora con funciones

```python
def suma(a, b):
    return a + b

def resta(a, b):
    return a - b

def multiplicacion(a, b):
    return a * b

def division(a, b):
    if b == 0:
        return "Error: División por cero"
    return a / b

def calculadora():
    operaciones = {
        '1': ('Suma', suma),
        '2': ('Resta', resta),
        '3': ('Multiplicación', multiplicacion),
        '4': ('División', division),
    }
    
    while True:
        print("\nCalculadora Simple")
        for key, (nombre, _) in operaciones.items():
            print(f"{key}. {nombre}")
        print("5. Salir")
        
        opcion = input("Seleccione una opción (1-5): ")
        
        if opcion == '5':
            print("¡Hasta luego!")
            break
            
        if opcion not in operaciones:
            print("Opción no válida. Intente de nuevo.")
            continue
            
        try:
            num1 = float(input("Ingrese el primer número: "))
            num2 = float(input("Ingrese el segundo número: "))
            
            nombre_operacion, funcion = operaciones[opcion]
            resultado = funcion(num1, num2)
            print(f"Resultado de {nombre_operacion}: {resultado}")
        except ValueError:
            print("Error: Ingrese números válidos")

if __name__ == "__main__":
    calculadora()
```

### Generador de contraseñas

```python
import random
import string

def generar_contrasena(longitud=12, incluir_mayusculas=True, incluir_numeros=True, incluir_especiales=True):
    """Genera una contraseña aleatoria con los criterios especificados.
    
    Args:
        longitud (int): Longitud de la contraseña.
        incluir_mayusculas (bool): Incluir letras mayúsculas.
        incluir_numeros (bool): Incluir números.
        incluir_especiales (bool): Incluir caracteres especiales.
        
    Returns:
        str: Contraseña generada.
    """
    # Definir los conjuntos de caracteres
    minusculas = string.ascii_lowercase
    mayusculas = string.ascii_uppercase if incluir_mayusculas else ""
    numeros = string.digits if incluir_numeros else ""
    especiales = string.punctuation if incluir_especiales else ""
    
    # Combinar todos los caracteres disponibles
    todos_caracteres = minusculas + mayusculas + numeros + especiales
    
    if not todos_caracteres:
        return "Error: No hay caracteres disponibles para generar la contraseña"
    
    # Asegurar que al menos un carácter de cada tipo esté presente si se solicita
    contrasena = []
    if incluir_mayusculas and mayusculas:
        contrasena.append(random.choice(mayusculas))
    if incluir_numeros and numeros:
        contrasena.append(random.choice(numeros))
    if incluir_especiales and especiales:
        contrasena.append(random.choice(especiales))
    
    # Completar con caracteres aleatorios hasta alcanzar la longitud deseada
    while len(contrasena) < longitud:
        contrasena.append(random.choice(todos_caracteres))
    
    # Mezclar los caracteres para evitar patrones predecibles
    random.shuffle(contrasena)
    
    # Convertir la lista de caracteres a una cadena
    return ''.join(contrasena)

def main():
    print("Generador de Contraseñas\n")
    
    try:
        longitud = int(input("Longitud de la contraseña (mínimo 4): "))
        if longitud < 4:
            print("La longitud mínima es 4. Se usará 4 como longitud.")
            longitud = 4
    except ValueError:
        print("Valor no válido. Se usará la longitud predeterminada (12).")
        longitud = 12
    
    incluir_mayusculas = input("¿Incluir mayúsculas? (s/n): ").lower() == 's'
    incluir_numeros = input("¿Incluir números? (s/n): ").lower() == 's'
    incluir_especiales = input("¿Incluir caracteres especiales? (s/n): ").lower() == 's'
    
    contrasena = generar_contrasena(
        longitud, 
        incluir_mayusculas, 
        incluir_numeros, 
        incluir_especiales
    )
    
    print(f"\nContraseña generada: {contrasena}")

if __name__ == "__main__":
    main()
```

### Función para calcular estadísticas

```python
def calcular_estadisticas(numeros):
    """Calcula varias estadísticas para una lista de números.
    
    Args:
        numeros (list): Lista de números.
        
    Returns:
        dict: Diccionario con las estadísticas calculadas.
    """
    if not numeros:
        return {"error": "La lista está vacía"}
    
    # Ordenar la lista para facilitar algunos cálculos
    ordenados = sorted(numeros)
    n = len(ordenados)
    
    # Calcular estadísticas básicas
    minimo = ordenados[0]
    maximo = ordenados[-1]
    suma = sum(ordenados)
    media = suma / n
    
    # Calcular la mediana
    if n % 2 == 0:  # Si hay un número par de elementos
        mediana = (ordenados[n//2 - 1] + ordenados[n//2]) / 2
    else:  # Si hay un número impar de elementos
        mediana = ordenados[n//2]
    
    # Calcular la moda (valor más frecuente)
    frecuencias = {}
    for num in ordenados:
        frecuencias[num] = frecuencias.get(num, 0) + 1
    
    max_frecuencia = max(frecuencias.values())
    moda = [num for num, freq in frecuencias.items() if freq == max_frecuencia]
    
    # Calcular la desviación estándar
    suma_cuadrados_diff = sum((x - media) ** 2 for x in numeros)
    desviacion_estandar = (suma_cuadrados_diff / n) ** 0.5
    
    # Calcular rango y rango intercuartil
    rango = maximo - minimo
    q1_pos = n // 4
    q3_pos = 3 * n // 4
    q1 = ordenados[q1_pos] if n % 4 != 0 else (ordenados[q1_pos - 1] + ordenados[q1_pos]) / 2
    q3 = ordenados[q3_pos] if n % 4 != 0 else (ordenados[q3_pos - 1] + ordenados[q3_pos]) / 2
    rango_intercuartil = q3 - q1
    
    return {
        "cantidad": n,
        "minimo": minimo,
        "maximo": maximo,
        "suma": suma,
        "media": media,
        "mediana": mediana,
        "moda": moda[0] if len(moda) == 1 else moda,
        "desviacion_estandar": desviacion_estandar,
        "rango": rango,
        "q1": q1,
        "q3": q3,
        "rango_intercuartil": rango_intercuartil
    }

def main():
    print("Calculadora de Estadísticas\n")
    
    entrada = input("Ingrese una lista de números separados por espacios: ")
    
    try:
        numeros = [float(x) for x in entrada.split()]
        if not numeros:
            print("No se ingresaron números.")
            return
        
        estadisticas = calcular_estadisticas(numeros)
        
        print("\nEstadísticas:")
        for clave, valor in estadisticas.items():
            print(f"{clave.capitalize()}: {valor}")
    
    except ValueError:
        print("Error: Ingrese solo números válidos separados por espacios.")

if __name__ == "__main__":
    main()
```

## Buenas Prácticas

1. **Nombres descriptivos**: Usa nombres de funciones que describan claramente lo que hacen (verbos o frases verbales).

2. **Funciones pequeñas y específicas**: Cada función debe hacer una sola cosa y hacerla bien. Divide funciones grandes en funciones más pequeñas y específicas.

3. **Docstrings**: Documenta tus funciones con docstrings que expliquen qué hace la función, sus parámetros y valores de retorno.

4. **Valores predeterminados inmutables**: Usa solo objetos inmutables (como números, cadenas, tuplas) como valores predeterminados para los parámetros.

5. **Evita efectos secundarios**: Las funciones deben evitar modificar variables globales o parámetros mutables, a menos que ese sea su propósito explícito.

6. **Manejo de errores**: Incluye validación de parámetros y manejo de casos especiales.

7. **Principio DRY (Don't Repeat Yourself)**: Si encuentras código duplicado, considera extraerlo a una función.

8. **Funciones puras**: Cuando sea posible, escribe funciones puras (que siempre producen el mismo resultado para los mismos argumentos y no tienen efectos secundarios).

9. **Limita el número de parámetros**: Intenta mantener el número de parámetros bajo (idealmente menos de 5). Si necesitas más, considera usar un objeto o diccionario.

10. **Usa anotaciones de tipo**: Para proyectos más grandes, considera usar anotaciones de tipo para mejorar la legibilidad y permitir verificaciones estáticas.

```python
def saludar(nombre: str, edad: int = 30) -> str:
    return f"Hola {nombre}, tienes {edad} años."
```

## Recursos Adicionales

- [Documentación oficial de Python - Definiendo funciones](https://docs.python.org/3/tutorial/controlflow.html#defining-functions)
- [Real Python - Defining Your Own Python Function](https://realpython.com/defining-your-own-python-function/)
- [Python Function Arguments: A Deep Dive](https://realpython.com/python-kwargs-and-args/)
- [Decorators in Python](https://realpython.com/primer-on-python-decorators/)
- [Generator Functions in Python](https://realpython.com/introduction-to-python-generators/)
- [Python Type Checking](https://realpython.com/python-type-checking/)

---

En la siguiente sección, aprenderemos sobre estructuras de datos en Python, que nos permitirán organizar y manipular datos de manera eficiente.