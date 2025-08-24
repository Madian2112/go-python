# Variables y Tipos de Datos en Python

## Variables

En Python, una variable es un nombre que hace referencia a un valor almacenado en la memoria. A diferencia de otros lenguajes, no es necesario declarar el tipo de variable explícitamente.

### Asignación de variables

Para crear una variable, simplemente asigna un valor usando el operador `=`:

```python
nombre = "Ana"
edad = 25
altura = 1.75
es_estudiante = True
```

### Reglas para nombrar variables

- Deben comenzar con una letra (a-z, A-Z) o un guion bajo (_)
- El resto del nombre puede contener letras, números y guiones bajos
- Son sensibles a mayúsculas y minúsculas (edad, Edad y EDAD son variables diferentes)
- No pueden ser palabras reservadas (como `if`, `for`, `while`, etc.)

### Convenciones de nomenclatura (PEP 8)

- Usa nombres en minúsculas separados por guiones bajos para variables: `mi_variable`
- Usa nombres descriptivos: `nombre_usuario` es mejor que `nu`
- Las constantes suelen escribirse en mayúsculas: `PI = 3.14159`

## Tipos de datos básicos

Python tiene varios tipos de datos incorporados:

### Números

#### Enteros (int)

Números enteros sin parte decimal:

```python
edad = 25
temperatura_negativa = -10
numero_grande = 1_000_000  # Los guiones bajos mejoran la legibilidad
```

#### Punto flotante (float)

Números con parte decimal:

```python
altura = 1.75
pi = 3.14159
notacion_cientifica = 1.23e-4  # 1.23 × 10^-4 = 0.000123
```

#### Complejos (complex)

Números con parte real e imaginaria:

```python
z = 2 + 3j
```

### Booleanos (bool)

Representan valores de verdad, solo pueden ser `True` o `False`:

```python
es_mayor_de_edad = True
tiene_descuento = False
```

### Cadenas de texto (str)

Secuencias de caracteres, delimitadas por comillas simples o dobles:

```python
nombre = "María"
apellido = 'González'
```

#### Cadenas multilínea

Usa triple comillas para cadenas que abarcan varias líneas:

```python
descripcion = """Este es un texto
que ocupa varias
líneas."""
```

#### Operaciones con cadenas

```python
# Concatenación
nombre_completo = nombre + " " + apellido  # "María González"

# Repetición
linea = "-" * 10  # "----------"

# Indexación (acceso a caracteres individuales)
primera_letra = nombre[0]  # "M"

# Slicing (rebanado)
primeras_tres = nombre[0:3]  # "Mar"

# Longitud
longitud = len(nombre)  # 5
```

### None

Representa la ausencia de valor o un valor nulo:

```python
resultado = None
```

## Conversión entre tipos

Python permite convertir entre diferentes tipos de datos:

```python
# De string a número
edad_str = "25"
edad_int = int(edad_str)  # 25

# De número a string
temperatura = 36.6
temperatura_str = str(temperatura)  # "36.6"

# De entero a float
numero = 10
numero_float = float(numero)  # 10.0

# De float a entero (trunca la parte decimal)
precio = 23.99
precio_int = int(precio)  # 23

# A booleano
# Nota: 0, cadenas vacías, None, y colecciones vacías se convierten a False
# Todo lo demás se convierte a True
print(bool(0))  # False
print(bool(1))  # True
print(bool(""))  # False
print(bool("Hola"))  # True
```

## Verificación de tipos

Puedes verificar el tipo de una variable con la función `type()`:

```python
x = 10
y = "Hola"
z = 3.14

print(type(x))  # <class 'int'>
print(type(y))  # <class 'str'>
print(type(z))  # <class 'float'>
```

## Strings formateados (f-strings)

Desde Python 3.6, puedes usar f-strings para incluir variables dentro de cadenas de texto:

```python
nombre = "Carlos"
edad = 30

# Formato tradicional
presentacion1 = "Me llamo {} y tengo {} años".format(nombre, edad)

# f-string (más legible y conciso)
presentacion2 = f"Me llamo {nombre} y tengo {edad} años"

# Puedes incluir expresiones
presentacion3 = f"Me llamo {nombre} y el año que viene tendré {edad + 1} años"
```

## Constantes

Python no tiene un mecanismo incorporado para declarar constantes, pero por convención, se usan nombres en mayúsculas para indicar que una variable no debe modificarse:

```python
PI = 3.14159
GRAVEDAD = 9.8
DIAS_SEMANA = 7
```

## Ejemplo práctico

Veamos un ejemplo que utiliza diferentes tipos de datos:

```python
# Programa para calcular el área de un círculo
import math

# Entrada del usuario (siempre se obtiene como string)
radio_str = input("Ingresa el radio del círculo: ")

# Conversión a float
radio = float(radio_str)

# Cálculo del área
area = math.pi * radio ** 2

# Mostrar resultado
print(f"Un círculo con radio {radio} tiene un área de {area:.2f} unidades cuadradas")
```

## Buenas prácticas

- Usa nombres descriptivos para tus variables
- Sigue las convenciones de PEP 8 para la nomenclatura
- Evita usar variables globales cuando sea posible
- Inicializa tus variables antes de usarlas
- Usa comentarios para explicar el propósito de variables complejas
- Utiliza f-strings para formatear cadenas (en Python 3.6+)

## Recursos adicionales

- [Documentación oficial sobre tipos de datos](https://docs.python.org/es/3/library/stdtypes.html)
- [PEP 8 - Guía de estilo para Python](https://www.python.org/dev/peps/pep-0008/)
- [Real Python: Variables en Python](https://realpython.com/python-variables/)

---

En la siguiente sección, aprenderemos sobre operadores y expresiones en Python.