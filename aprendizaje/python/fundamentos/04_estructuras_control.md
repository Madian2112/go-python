# Estructuras de Control en Python

## Introducción

Las estructuras de control son bloques de código que determinan el flujo de ejecución de un programa. Python ofrece varias estructuras de control que permiten tomar decisiones, repetir acciones y organizar el código de manera eficiente.

## Estructuras Condicionales

Las estructuras condicionales permiten ejecutar diferentes bloques de código dependiendo de si una condición es verdadera o falsa.

### if, elif, else

La estructura básica de una condición en Python es:

```python
if condicion:
    # Código a ejecutar si la condición es True
elif otra_condicion:
    # Código a ejecutar si la condición anterior es False y esta es True
else:
    # Código a ejecutar si todas las condiciones anteriores son False
```

Ejemplo:

```python
edad = 18

if edad < 18:
    print("Eres menor de edad")
elif edad == 18:
    print("Acabas de cumplir la mayoría de edad")
else:
    print("Eres mayor de edad")
```

### Operadores de comparación en condiciones

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `==` | Igual a | `a == b` |
| `!=` | Diferente de | `a != b` |
| `>` | Mayor que | `a > b` |
| `<` | Menor que | `a < b` |
| `>=` | Mayor o igual que | `a >= b` |
| `<=` | Menor o igual que | `a <= b` |

### Operadores lógicos

Permiten combinar múltiples condiciones:

| Operador | Descripción | Ejemplo |
|----------|-------------|--------|
| `and` | Verdadero si ambas condiciones son verdaderas | `a > 0 and b > 0` |
| `or` | Verdadero si al menos una condición es verdadera | `a > 0 or b > 0` |
| `not` | Invierte el valor de la condición | `not a > 0` |

Ejemplo:

```python
edad = 25
ingresos = 1500

if edad > 18 and ingresos > 1000:
    print("Cumples los requisitos para el préstamo")
else:
    print("No cumples los requisitos")
```

### Operador ternario

Python permite escribir condicionales simples en una sola línea:

```python
resultado = "Aprobado" if calificacion >= 60 else "Reprobado"
```

Este código es equivalente a:

```python
if calificacion >= 60:
    resultado = "Aprobado"
else:
    resultado = "Reprobado"
```

### Evaluación de valores como condiciones

En Python, los siguientes valores se evalúan como `False` en un contexto booleano:
- `False`
- `None`
- `0` (cero)
- `""` (cadena vacía)
- `[]` (lista vacía)
- `()` (tupla vacía)
- `{}` (diccionario vacío)
- `set()` (conjunto vacío)

Cualquier otro valor se evalúa como `True`.

Ejemplo:

```python
nombre = ""

if nombre:  # Equivalente a if nombre != ""
    print(f"Hola, {nombre}")
else:
    print("Por favor, introduce tu nombre")
```

## Estructuras de Repetición (Bucles)

Las estructuras de repetición permiten ejecutar un bloque de código múltiples veces.

### Bucle for

El bucle `for` en Python itera sobre elementos de una secuencia (lista, tupla, diccionario, conjunto o cadena):

```python
for elemento in secuencia:
    # Código a ejecutar para cada elemento
```

Ejemplos:

```python
# Iterando sobre una lista
frutas = ["manzana", "banana", "cereza"]
for fruta in frutas:
    print(fruta)

# Iterando sobre un rango de números
for i in range(5):  # 0, 1, 2, 3, 4
    print(i)

# Iterando sobre un rango con inicio y fin
for i in range(2, 8):  # 2, 3, 4, 5, 6, 7
    print(i)

# Iterando sobre un rango con paso
for i in range(0, 10, 2):  # 0, 2, 4, 6, 8
    print(i)

# Iterando sobre una cadena
for letra in "Python":
    print(letra)

# Iterando sobre un diccionario
persona = {"nombre": "Ana", "edad": 25, "ciudad": "Madrid"}
for clave in persona:
    print(f"{clave}: {persona[clave]}")

# Iterando sobre clave y valor de un diccionario
for clave, valor in persona.items():
    print(f"{clave}: {valor}")
```

### Función enumerate()

La función `enumerate()` permite obtener el índice y el valor de cada elemento durante la iteración:

```python
frutas = ["manzana", "banana", "cereza"]
for indice, fruta in enumerate(frutas):
    print(f"{indice}: {fruta}")
```

Salida:
```
0: manzana
1: banana
2: cereza
```

### Bucle while

El bucle `while` ejecuta un bloque de código mientras una condición sea verdadera:

```python
while condicion:
    # Código a ejecutar mientras la condición sea True
```

Ejemplo:

```python
contador = 0
while contador < 5:
    print(contador)
    contador += 1
```

### Instrucciones break y continue

- `break`: Termina el bucle actual y continúa con la siguiente instrucción después del bucle.
- `continue`: Salta a la siguiente iteración del bucle, omitiendo el resto del código en la iteración actual.

Ejemplos:

```python
# Uso de break
for i in range(10):
    if i == 5:
        break  # Sale del bucle cuando i es 5
    print(i)  # Imprime: 0, 1, 2, 3, 4

# Uso de continue
for i in range(10):
    if i % 2 == 0:
        continue  # Salta a la siguiente iteración si i es par
    print(i)  # Imprime: 1, 3, 5, 7, 9
```

### Bucle else

En Python, los bucles `for` y `while` pueden tener una cláusula `else` que se ejecuta cuando el bucle termina normalmente (sin `break`):

```python
for i in range(5):
    print(i)
else:
    print("Bucle completado sin interrupciones")

# Con break
for i in range(5):
    if i == 3:
        break
    print(i)
else:
    print("Este mensaje no se imprimirá porque el bucle se interrumpió")
```

## Comprensiones de Listas

Las comprensiones de listas son una forma concisa de crear listas basadas en listas existentes:

```python
# Sintaxis básica
nueva_lista = [expresion for elemento in iterable if condicion]
```

Ejemplos:

```python
# Lista de cuadrados
cuadrados = [x**2 for x in range(10)]
# Equivalente a:
# cuadrados = []
# for x in range(10):
#     cuadrados.append(x**2)

# Lista de números pares
pares = [x for x in range(20) if x % 2 == 0]

# Convertir temperaturas de Celsius a Fahrenheit
celsius = [0, 10, 20, 30, 40]
fahrenheit = [(9/5) * c + 32 for c in celsius]

# Comprensión de lista anidada (matriz 3x3)
matriz = [[i * 3 + j + 1 for j in range(3)] for i in range(3)]
# Resultado: [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
```

## Comprensiones de Diccionarios y Conjuntos

Similar a las comprensiones de listas, Python permite crear diccionarios y conjuntos de forma concisa:

```python
# Comprensión de diccionario
cuadrados_dict = {x: x**2 for x in range(6)}
# Resultado: {0: 0, 1: 1, 2: 4, 3: 9, 4: 16, 5: 25}

# Comprensión de conjunto
cuadrados_set = {x**2 for x in range(-5, 5)}
# Resultado: {0, 1, 4, 9, 16, 25}
```

## Ejemplos Prácticos

### Calculadora simple con menú

```python
def main():
    while True:
        print("\nCalculadora Simple")
        print("1. Suma")
        print("2. Resta")
        print("3. Multiplicación")
        print("4. División")
        print("5. Salir")
        
        opcion = input("Seleccione una opción (1-5): ")
        
        if opcion == '5':
            print("¡Hasta luego!")
            break
            
        if opcion not in ['1', '2', '3', '4']:
            print("Opción no válida. Intente de nuevo.")
            continue
            
        num1 = float(input("Ingrese el primer número: "))
        num2 = float(input("Ingrese el segundo número: "))
        
        if opcion == '1':
            print(f"Resultado: {num1 + num2}")
        elif opcion == '2':
            print(f"Resultado: {num1 - num2}")
        elif opcion == '3':
            print(f"Resultado: {num1 * num2}")
        elif opcion == '4':
            if num2 == 0:
                print("Error: No se puede dividir por cero")
            else:
                print(f"Resultado: {num1 / num2}")

if __name__ == "__main__":
    main()
```

### Verificación de número primo

```python
def es_primo(numero):
    if numero <= 1:
        return False
    if numero <= 3:
        return True
    if numero % 2 == 0 or numero % 3 == 0:
        return False
    
    i = 5
    while i * i <= numero:
        if numero % i == 0 or numero % (i + 2) == 0:
            return False
        i += 6
    
    return True

def main():
    while True:
        try:
            numero = int(input("Ingrese un número (0 para salir): "))
            
            if numero == 0:
                print("¡Hasta luego!")
                break
                
            if es_primo(numero):
                print(f"{numero} es un número primo")
            else:
                print(f"{numero} no es un número primo")
                
        except ValueError:
            print("Por favor, ingrese un número entero válido")

if __name__ == "__main__":
    main()
```

### Generador de patrones

```python
def imprimir_patron(filas):
    # Patrón de triángulo
    print("\nPatrón 1:")
    for i in range(1, filas + 1):
        print("*" * i)
    
    # Patrón de pirámide
    print("\nPatrón 2:")
    for i in range(1, filas + 1):
        espacios = filas - i
        estrellas = 2 * i - 1
        print(" " * espacios + "*" * estrellas)
    
    # Patrón de diamante
    print("\nPatrón 3:")
    for i in range(1, filas + 1):
        espacios = filas - i
        estrellas = 2 * i - 1
        print(" " * espacios + "*" * estrellas)
    
    for i in range(filas - 1, 0, -1):
        espacios = filas - i
        estrellas = 2 * i - 1
        print(" " * espacios + "*" * estrellas)

def main():
    try:
        filas = int(input("Ingrese el número de filas para los patrones: "))
        if filas <= 0:
            print("Por favor, ingrese un número positivo")
        else:
            imprimir_patron(filas)
    except ValueError:
        print("Por favor, ingrese un número entero válido")

if __name__ == "__main__":
    main()
```

## Buenas Prácticas

1. **Indentación consistente**: Python utiliza la indentación para definir bloques de código. Usa 4 espacios por nivel de indentación (recomendación de PEP 8).

2. **Evita bucles anidados profundos**: Los bucles anidados pueden hacer que el código sea difícil de leer y mantener. Intenta limitar la anidación a 2-3 niveles.

3. **Usa comprensiones de listas con moderación**: Son concisas y eficientes, pero pueden dificultar la lectura si son demasiado complejas.

4. **Prefiere `for` sobre `while` cuando sea posible**: Los bucles `for` son generalmente más seguros y menos propensos a errores que los bucles `while`.

5. **Evita `break` y `continue` excesivos**: Aunque son útiles, su uso excesivo puede hacer que el flujo del programa sea difícil de seguir.

6. **Usa nombres descriptivos para variables de control**: En lugar de `i`, `j`, `k`, considera nombres como `indice`, `fila`, `columna`.

7. **Verifica condiciones de borde**: Asegúrate de que tus bucles y condiciones manejen correctamente los casos límite.

8. **Evita condiciones complejas**: Si una condición es demasiado compleja, considera dividirla en partes más pequeñas o usar variables intermedias con nombres descriptivos.

## Recursos Adicionales

- [Documentación oficial de Python - Control Flow](https://docs.python.org/3/tutorial/controlflow.html)
- [Real Python - Conditional Statements in Python](https://realpython.com/python-conditional-statements/)
- [Real Python - Python for Loops](https://realpython.com/python-for-loop/)
- [Real Python - Python while Loops](https://realpython.com/python-while-loop/)
- [Python List Comprehensions: Explained Visually](https://treyhunner.com/2015/12/python-list-comprehensions-now-in-color/)

---

En la siguiente sección, aprenderemos sobre funciones en Python, que nos permitirán organizar y reutilizar nuestro código de manera más eficiente.