# Operadores y Expresiones en Python

## Introducción

Los operadores son símbolos especiales que realizan operaciones sobre variables y valores. Las expresiones son combinaciones de valores, variables y operadores que se evalúan para producir un resultado.

## Operadores Aritméticos

Realizan operaciones matemáticas básicas:

| Operador | Descripción | Ejemplo | Resultado |
|----------|-------------|---------|----------|
| `+` | Suma | `5 + 3` | `8` |
| `-` | Resta | `5 - 3` | `2` |
| `*` | Multiplicación | `5 * 3` | `15` |
| `/` | División | `5 / 3` | `1.6666...` |
| `//` | División entera | `5 // 3` | `1` |
| `%` | Módulo (resto) | `5 % 3` | `2` |
| `**` | Potencia | `5 ** 3` | `125` |

```python
# Ejemplos de operadores aritméticos
a = 10
b = 3

suma = a + b          # 13
resta = a - b         # 7
multiplicacion = a * b  # 30
division = a / b      # 3.3333...
div_entera = a // b   # 3
modulo = a % b        # 1
potencia = a ** b     # 1000
```

### Precedencia de operadores aritméticos

Los operadores se evalúan en el siguiente orden:

1. Paréntesis `()`
2. Potencia `**`
3. Multiplicación `*`, división `/`, división entera `//`, módulo `%`
4. Suma `+`, resta `-`

```python
resultado = 2 + 3 * 4    # 14 (no 20, porque * tiene mayor precedencia que +)
resultado = (2 + 3) * 4  # 20 (los paréntesis cambian la precedencia)
```

## Operadores de Asignación

Asignan valores a variables:

| Operador | Descripción | Ejemplo | Equivalente a |
|----------|-------------|---------|---------------|
| `=` | Asignación simple | `x = 5` | `x = 5` |
| `+=` | Suma y asignación | `x += 3` | `x = x + 3` |
| `-=` | Resta y asignación | `x -= 3` | `x = x - 3` |
| `*=` | Multiplicación y asignación | `x *= 3` | `x = x * 3` |
| `/=` | División y asignación | `x /= 3` | `x = x / 3` |
| `//=` | División entera y asignación | `x //= 3` | `x = x // 3` |
| `%=` | Módulo y asignación | `x %= 3` | `x = x % 3` |
| `**=` | Potencia y asignación | `x **= 3` | `x = x ** 3` |

```python
# Ejemplos de operadores de asignación
x = 10

x += 5   # x ahora es 15
x -= 3   # x ahora es 12
x *= 2   # x ahora es 24
x /= 4   # x ahora es 6.0 (nota: ahora x es float)
x //= 2  # x ahora es 3.0 (sigue siendo float)
x %= 2   # x ahora es 1.0
x **= 3  # x ahora es 1.0
```

## Operadores de Comparación

Comparan valores y devuelven un resultado booleano (True o False):

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `==` | Igual a | `a == b` |
| `!=` | Diferente de | `a != b` |
| `>` | Mayor que | `a > b` |
| `<` | Menor que | `a < b` |
| `>=` | Mayor o igual que | `a >= b` |
| `<=` | Menor o igual que | `a <= b` |

```python
# Ejemplos de operadores de comparación
a = 10
b = 5

print(a == b)  # False
print(a != b)  # True
print(a > b)   # True
print(a < b)   # False
print(a >= b)  # True
print(a <= b)  # False

# También funcionan con strings (comparación lexicográfica)
print("apple" < "banana")  # True
print("apple" == "Apple")  # False (sensible a mayúsculas/minúsculas)
```

## Operadores Lógicos

Combina expresiones booleanas:

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `and` | Y lógico (True si ambos operandos son True) | `a and b` |
| `or` | O lógico (True si al menos un operando es True) | `a or b` |
| `not` | Negación (invierte el valor booleano) | `not a` |

```python
# Ejemplos de operadores lógicos
x = 5
y = 10

print(x > 0 and y > 0)    # True (ambas condiciones son verdaderas)
print(x > 10 or y > 5)    # True (la segunda condición es verdadera)
print(not x > 10)         # True (x > 10 es falso, not lo invierte)

# Cortocircuito en operadores lógicos
# and: si el primer operando es False, el segundo no se evalúa
# or: si el primer operando es True, el segundo no se evalúa
```

### Tabla de verdad para operadores lógicos

| a | b | a and b | a or b | not a |
|---|---|---------|--------|-------|
| True | True | True | True | False |
| True | False | False | True | False |
| False | True | False | True | True |
| False | False | False | False | True |

## Operadores de Identidad

Comprueban si dos variables hacen referencia al mismo objeto en memoria:

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `is` | True si ambas variables apuntan al mismo objeto | `a is b` |
| `is not` | True si las variables apuntan a objetos diferentes | `a is not b` |

```python
# Ejemplos de operadores de identidad
a = [1, 2, 3]
b = [1, 2, 3]
c = a

print(a is b)      # False (mismos valores pero objetos diferentes)
print(a is c)      # True (ambos apuntan al mismo objeto)
print(a is not b)  # True

# Casos especiales con None
x = None
print(x is None)   # True (forma correcta de comprobar si es None)
```

## Operadores de Pertenencia

Comprueban si un valor está presente en una secuencia:

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `in` | True si el valor está en la secuencia | `a in b` |
| `not in` | True si el valor no está en la secuencia | `a not in b` |

```python
# Ejemplos de operadores de pertenencia
frutas = ["manzana", "banana", "cereza"]

print("manzana" in frutas)       # True
print("pera" in frutas)          # False
print("pera" not in frutas)      # True

# También funcionan con strings
print("a" in "manzana")          # True
print("z" not in "manzana")      # False

# Y con diccionarios (comprueba las claves, no los valores)
diccionario = {"a": 1, "b": 2}
print("a" in diccionario)        # True
print(1 in diccionario)          # False
```

## Operadores Bit a Bit

Realizan operaciones a nivel de bits:

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `&` | AND bit a bit | `a & b` |
| `\|` | OR bit a bit | `a \| b` |
| `^` | XOR bit a bit | `a ^ b` |
| `~` | NOT bit a bit (complemento a uno) | `~a` |
| `<<` | Desplazamiento a la izquierda | `a << n` |
| `>>` | Desplazamiento a la derecha | `a >> n` |

```python
# Ejemplos de operadores bit a bit
a = 60  # 0011 1100 en binario
b = 13  # 0000 1101 en binario

print(a & b)   # 12 (0000 1100)
print(a | b)   # 61 (0011 1101)
print(a ^ b)   # 49 (0011 0001)
print(~a)      # -61 (complemento a uno)
print(a << 2)  # 240 (1111 0000)
print(a >> 2)  # 15 (0000 1111)
```

## Expresiones

Las expresiones son combinaciones de valores, variables y operadores que se evalúan para producir un resultado.

### Tipos de expresiones

1. **Expresiones aritméticas**: Realizan cálculos matemáticos
   ```python
   resultado = (a + b) * c / d
   ```

2. **Expresiones relacionales**: Comparan valores
   ```python
   es_mayor = edad >= 18
   ```

3. **Expresiones lógicas**: Combinan condiciones
   ```python
   puede_votar = edad >= 18 and es_ciudadano
   ```

4. **Expresiones de asignación**: Asignan valores a variables
   ```python
   x = y = z = 0  # Asignación múltiple
   a, b = 10, 20  # Desempaquetado
   ```

### Evaluación de expresiones

Python evalúa las expresiones siguiendo reglas de precedencia. Puedes usar paréntesis para cambiar el orden de evaluación:

```python
# Sin paréntesis
resultado1 = 2 + 3 * 4  # 14

# Con paréntesis
resultado2 = (2 + 3) * 4  # 20
```

## Expresiones condicionales (operador ternario)

Python permite escribir expresiones condicionales en una sola línea:

```python
# Sintaxis: valor_si_verdadero if condicion else valor_si_falso

# Ejemplo tradicional con if-else
if edad >= 18:
    estado = "adulto"
else:
    estado = "menor"

# Equivalente con expresión condicional
estado = "adulto" if edad >= 18 else "menor"
```

## Ejemplos prácticos

### Calculadora simple

```python
# Calculadora básica
num1 = float(input("Ingresa el primer número: "))
num2 = float(input("Ingresa el segundo número: "))
operacion = input("Ingresa la operación (+, -, *, /): ")

if operacion == "+":
    resultado = num1 + num2
elif operacion == "-":
    resultado = num1 - num2
elif operacion == "*":
    resultado = num1 * num2
elif operacion == "/":
    resultado = num1 / num2 if num2 != 0 else "Error: División por cero"
else:
    resultado = "Operación no válida"

print(f"Resultado: {resultado}")
```

### Verificación de año bisiesto

```python
# Un año es bisiesto si es divisible por 4,
# excepto los divisibles por 100 que no son divisibles por 400
año = int(input("Ingresa un año: "))

es_bisiesto = (año % 4 == 0 and año % 100 != 0) or (año % 400 == 0)

print(f"{año} {'es' if es_bisiesto else 'no es'} un año bisiesto")
```

## Buenas prácticas

- Usa paréntesis para clarificar el orden de evaluación, incluso cuando no son estrictamente necesarios
- Evita expresiones demasiado complejas; divide en partes más simples si es necesario
- Usa variables con nombres descriptivos para almacenar resultados intermedios
- Aprovecha las expresiones condicionales para código más conciso, pero sin sacrificar legibilidad
- Ten cuidado con las divisiones por cero
- Recuerda que `==` compara valores, mientras que `is` compara identidad (referencias)

## Recursos adicionales

- [Documentación oficial sobre operadores](https://docs.python.org/es/3/reference/expressions.html)
- [Real Python: Operadores y expresiones](https://realpython.com/python-operators-expressions/)
- [Python Tutor](http://pythontutor.com/) - Visualiza la ejecución de expresiones paso a paso

---

En la siguiente sección, aprenderemos sobre estructuras de control en Python (condicionales y bucles).