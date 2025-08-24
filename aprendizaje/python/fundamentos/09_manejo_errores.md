# Manejo de Errores en Python

## Introducción

El manejo de errores es una parte fundamental de cualquier programa robusto. Python proporciona un sistema de excepciones que permite detectar y manejar errores de forma elegante, manteniendo el control del flujo del programa incluso cuando ocurren situaciones inesperadas. En esta sección, exploraremos cómo Python maneja los errores, cómo podemos capturar y gestionar excepciones, y cómo crear nuestras propias clases de excepciones personalizadas.

## Excepciones en Python

En Python, los errores durante la ejecución se manejan mediante excepciones. Una excepción es un evento que ocurre durante la ejecución de un programa y que interrumpe el flujo normal de las instrucciones.

### Tipos Comunes de Excepciones

Python tiene muchos tipos de excepciones incorporadas. Algunos de los más comunes son:

- `SyntaxError`: Error en la sintaxis del código.
- `NameError`: Se intenta usar una variable que no ha sido definida.
- `TypeError`: Operación aplicada a un objeto de tipo inapropiado.
- `ValueError`: Operación con un valor inapropiado.
- `IndexError`: Índice fuera de rango.
- `KeyError`: Clave no encontrada en un diccionario.
- `FileNotFoundError`: Archivo no encontrado.
- `ZeroDivisionError`: División por cero.
- `ImportError`: Error al importar un módulo.
- `AttributeError`: Atributo no encontrado en un objeto.

### Jerarquía de Excepciones

Las excepciones en Python forman una jerarquía. Todas las excepciones incorporadas derivan de la clase base `BaseException`, aunque normalmente trabajaremos con excepciones que derivan de `Exception`.

```python
BaseException
 ├── SystemExit
 ├── KeyboardInterrupt
 ├── GeneratorExit
 └── Exception
      ├── StopIteration
      ├── ArithmeticError
      │    ├── FloatingPointError
      │    ├── OverflowError
      │    └── ZeroDivisionError
      ├── AssertionError
      ├── AttributeError
      ├── BufferError
      ├── EOFError
      ├── ImportError
      │    └── ModuleNotFoundError
      ├── LookupError
      │    ├── IndexError
      │    └── KeyError
      ├── MemoryError
      ├── NameError
      │    └── UnboundLocalError
      ├── OSError
      │    ├── BlockingIOError
      │    ├── ChildProcessError
      │    ├── ConnectionError
      │    │    ├── BrokenPipeError
      │    │    ├── ConnectionAbortedError
      │    │    ├── ConnectionRefusedError
      │    │    └── ConnectionResetError
      │    ├── FileExistsError
      │    ├── FileNotFoundError
      │    ├── InterruptedError
      │    ├── IsADirectoryError
      │    ├── NotADirectoryError
      │    ├── PermissionError
      │    ├── ProcessLookupError
      │    └── TimeoutError
      ├── ReferenceError
      ├── RuntimeError
      │    ├── NotImplementedError
      │    └── RecursionError
      ├── SyntaxError
      │    └── IndentationError
      │         └── TabError
      ├── SystemError
      ├── TypeError
      ├── ValueError
      │    └── UnicodeError
      │         ├── UnicodeDecodeError
      │         ├── UnicodeEncodeError
      │         └── UnicodeTranslateError
      └── Warning
           ├── DeprecationWarning
           ├── PendingDeprecationWarning
           ├── RuntimeWarning
           ├── SyntaxWarning
           ├── UserWarning
           ├── FutureWarning
           ├── ImportWarning
           ├── UnicodeWarning
           ├── BytesWarning
           └── ResourceWarning
```

## Manejo de Excepciones

### Bloque try-except

El bloque `try-except` es la estructura básica para manejar excepciones en Python:

```python
try:
    # Código que puede generar una excepción
    resultado = 10 / 0
except ZeroDivisionError:
    # Código que se ejecuta si ocurre una ZeroDivisionError
    print("¡Error! División por cero.")
```

### Capturar Múltiples Excepciones

Podemos manejar diferentes tipos de excepciones de forma específica:

```python
try:
    numero = int(input("Introduce un número: "))
    resultado = 10 / numero
except ValueError:
    print("¡Error! Debes introducir un número válido.")
except ZeroDivisionError:
    print("¡Error! No puedes dividir por cero.")
```

También podemos capturar múltiples excepciones en una sola cláusula:

```python
try:
    numero = int(input("Introduce un número: "))
    resultado = 10 / numero
except (ValueError, ZeroDivisionError):
    print("¡Error! Entrada inválida o división por cero.")
```

### Capturar Todas las Excepciones

Podemos capturar cualquier excepción usando `except` sin especificar un tipo, aunque generalmente no es recomendable porque puede ocultar errores inesperados:

```python
try:
    # Algún código peligroso
    resultado = 10 / 0
except:
    print("Ocurrió un error.")
```

Una mejor práctica es capturar `Exception`, que incluye todas las excepciones estándar pero no las excepciones del sistema como `KeyboardInterrupt`:

```python
try:
    # Algún código peligroso
    resultado = 10 / 0
except Exception as e:
    print(f"Ocurrió un error: {e}")
```

### Cláusula else

La cláusula `else` se ejecuta si no ocurre ninguna excepción en el bloque `try`:

```python
try:
    numero = int(input("Introduce un número: "))
    resultado = 10 / numero
except ValueError:
    print("¡Error! Debes introducir un número válido.")
except ZeroDivisionError:
    print("¡Error! No puedes dividir por cero.")
else:
    print(f"El resultado es {resultado}")
```

### Cláusula finally

La cláusula `finally` se ejecuta siempre, independientemente de si ocurre una excepción o no. Es útil para tareas de limpieza:

```python
try:
    archivo = open("datos.txt", "r")
    contenido = archivo.read()
except FileNotFoundError:
    print("El archivo no existe.")
finally:
    # Esto se ejecuta siempre, incluso si hay una excepción
    if 'archivo' in locals() and not archivo.closed:
        archivo.close()
        print("Archivo cerrado.")
```

### Combinando else y finally

Podemos combinar todas las cláusulas en un solo bloque:

```python
try:
    numero = int(input("Introduce un número: "))
    resultado = 10 / numero
except ValueError:
    print("¡Error! Debes introducir un número válido.")
except ZeroDivisionError:
    print("¡Error! No puedes dividir por cero.")
else:
    print(f"El resultado es {resultado}")
finally:
    print("Operación completada.")
```

## Lanzar Excepciones

Podemos lanzar excepciones explícitamente usando la palabra clave `raise`:

```python
def verificar_edad(edad):
    if edad < 0:
        raise ValueError("La edad no puede ser negativa")
    if edad < 18:
        raise ValueError("Debes ser mayor de edad")
    return "Edad válida"

try:
    resultado = verificar_edad(15)
    print(resultado)
except ValueError as e:
    print(f"Error: {e}")
```

### Re-lanzar Excepciones

A veces queremos capturar una excepción, hacer algo con ella, y luego re-lanzarla para que sea manejada por un nivel superior:

```python
try:
    try:
        # Código que puede generar una excepción
        resultado = 10 / 0
    except ZeroDivisionError as e:
        print("Capturada internamente, pero la re-lanzo.")
        raise  # Re-lanza la última excepción
except ZeroDivisionError:
    print("Capturada externamente.")
```

## Excepciones Personalizadas

Podemos crear nuestras propias clases de excepciones heredando de `Exception` o alguna de sus subclases:

```python
class EdadInvalidaError(Exception):
    """Excepción lanzada cuando la edad es inválida."""
    def __init__(self, edad, mensaje="Edad inválida"):
        self.edad = edad
        self.mensaje = mensaje
        super().__init__(f"{mensaje}: {edad}")

def verificar_edad(edad):
    if edad < 0:
        raise EdadInvalidaError(edad, "La edad no puede ser negativa")
    if edad < 18:
        raise EdadInvalidaError(edad, "Debes ser mayor de edad")
    return "Edad válida"

try:
    resultado = verificar_edad(15)
    print(resultado)
except EdadInvalidaError as e:
    print(f"Error: {e}")
```

## Aserciones

Las aserciones son una forma de verificar que ciertas condiciones se cumplan durante la ejecución del programa. Si la condición es falsa, se lanza una excepción `AssertionError`:

```python
def dividir(a, b):
    assert b != 0, "El divisor no puede ser cero"
    return a / b

try:
    resultado = dividir(10, 0)
except AssertionError as e:
    print(f"Error de aserción: {e}")
```

Las aserciones son útiles durante el desarrollo y las pruebas, pero no deben usarse para manejar errores en producción, ya que pueden desactivarse con la opción `-O` del intérprete de Python.

## Context Managers (with)

Los context managers, utilizados con la declaración `with`, son una forma elegante de manejar recursos que necesitan ser liberados o cerrados después de su uso, incluso si ocurre una excepción:

```python
# Sin context manager
try:
    archivo = open("datos.txt", "r")
    contenido = archivo.read()
finally:
    archivo.close()

# Con context manager
with open("datos.txt", "r") as archivo:
    contenido = archivo.read()
# El archivo se cierra automáticamente al salir del bloque with
```

### Creando Context Managers Personalizados

Podemos crear nuestros propios context managers implementando los métodos `__enter__` y `__exit__`:

```python
class MiContextManager:
    def __init__(self, nombre):
        self.nombre = nombre
        
    def __enter__(self):
        print(f"Entrando en el contexto: {self.nombre}")
        return self  # El objeto devuelto se asigna a la variable después de 'as'
        
    def __exit__(self, tipo_exc, valor_exc, traceback_exc):
        print(f"Saliendo del contexto: {self.nombre}")
        if tipo_exc:
            print(f"Ocurrió una excepción: {valor_exc}")
            # Devolver True para suprimir la excepción
            return True

# Uso del context manager
with MiContextManager("Ejemplo") as cm:
    print("Dentro del bloque with")
    raise ValueError("¡Error de prueba!")

print("Después del bloque with")  # Esto se ejecuta porque la excepción fue suprimida
```

También podemos crear context managers más simples usando el decorador `@contextmanager` del módulo `contextlib`:

```python
from contextlib import contextmanager

@contextmanager
def mi_context_manager(nombre):
    print(f"Entrando en el contexto: {nombre}")
    try:
        yield nombre  # El valor cedido se asigna a la variable después de 'as'
        print("Saliendo normalmente del contexto")
    except Exception as e:
        print(f"Saliendo del contexto con excepción: {e}")
        # No re-lanzar la excepción (suprimirla)

# Uso del context manager
with mi_context_manager("Ejemplo") as nombre:
    print(f"Dentro del bloque with con {nombre}")
    raise ValueError("¡Error de prueba!")

print("Después del bloque with")  # Esto se ejecuta porque la excepción fue suprimida
```

## Ejemplo Práctico: Validador de Datos

Vamos a crear un validador de datos que utilice excepciones personalizadas para manejar diferentes tipos de errores de validación:

```python
class ValidacionError(Exception):
    """Clase base para excepciones de validación."""
    pass

class LongitudInvalidaError(ValidacionError):
    """Excepción lanzada cuando la longitud es inválida."""
    def __init__(self, valor, min_longitud, max_longitud):
        self.valor = valor
        self.min_longitud = min_longitud
        self.max_longitud = max_longitud
        mensaje = f"La longitud de '{valor}' debe estar entre {min_longitud} y {max_longitud} caracteres"
        super().__init__(mensaje)

class FormatoInvalidoError(ValidacionError):
    """Excepción lanzada cuando el formato es inválido."""
    def __init__(self, valor, patron):
        self.valor = valor
        self.patron = patron
        mensaje = f"'{valor}' no cumple con el formato requerido: {patron}"
        super().__init__(mensaje)

class ValorNoPermitidoError(ValidacionError):
    """Excepción lanzada cuando el valor no está permitido."""
    def __init__(self, valor, valores_permitidos):
        self.valor = valor
        self.valores_permitidos = valores_permitidos
        mensaje = f"'{valor}' no está en la lista de valores permitidos: {valores_permitidos}"
        super().__init__(mensaje)

class Validador:
    """Clase para validar datos según diferentes reglas."""
    
    @staticmethod
    def validar_longitud(valor, min_longitud, max_longitud):
        """Valida que la longitud del valor esté dentro del rango especificado."""
        if not min_longitud <= len(str(valor)) <= max_longitud:
            raise LongitudInvalidaError(valor, min_longitud, max_longitud)
        return valor
    
    @staticmethod
    def validar_formato(valor, patron):
        """Valida que el valor cumpla con el patrón especificado."""
        import re
        if not re.match(patron, str(valor)):
            raise FormatoInvalidoError(valor, patron)
        return valor
    
    @staticmethod
    def validar_permitido(valor, valores_permitidos):
        """Valida que el valor esté en la lista de valores permitidos."""
        if valor not in valores_permitidos:
            raise ValorNoPermitidoError(valor, valores_permitidos)
        return valor
    
    @classmethod
    def validar_email(cls, email):
        """Valida que el email tenga un formato correcto."""
        patron = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        cls.validar_longitud(email, 6, 100)
        cls.validar_formato(email, patron)
        return email
    
    @classmethod
    def validar_telefono(cls, telefono):
        """Valida que el teléfono tenga un formato correcto."""
        patron = r'^\+?[0-9]{8,15}$'
        cls.validar_formato(telefono, patron)
        return telefono
    
    @classmethod
    def validar_codigo_pais(cls, codigo):
        """Valida que el código de país sea válido."""
        codigos_validos = ['ES', 'US', 'UK', 'FR', 'DE', 'IT', 'PT', 'MX', 'AR', 'CL']
        cls.validar_longitud(codigo, 2, 2)
        cls.validar_permitido(codigo.upper(), codigos_validos)
        return codigo.upper()

# Ejemplo de uso
def registrar_usuario(nombre, email, telefono, pais):
    try:
        # Validar los datos
        Validador.validar_longitud(nombre, 2, 50)
        Validador.validar_email(email)
        Validador.validar_telefono(telefono)
        Validador.validar_codigo_pais(pais)
        
        # Si todas las validaciones pasan, registrar el usuario
        print(f"Usuario registrado: {nombre}, {email}, {telefono}, {pais}")
        return True
    except LongitudInvalidaError as e:
        print(f"Error de longitud: {e}")
    except FormatoInvalidoError as e:
        print(f"Error de formato: {e}")
    except ValorNoPermitidoError as e:
        print(f"Error de valor no permitido: {e}")
    except ValidacionError as e:
        print(f"Error de validación: {e}")
    except Exception as e:
        print(f"Error inesperado: {e}")
    
    return False

# Pruebas
registrar_usuario("Ana", "ana@example.com", "+34123456789", "ES")  # Válido
registrar_usuario("B", "ana@example.com", "+34123456789", "ES")    # Nombre muy corto
registrar_usuario("Ana", "ana@example", "+34123456789", "ES")      # Email inválido
registrar_usuario("Ana", "ana@example.com", "123", "ES")           # Teléfono inválido
registrar_usuario("Ana", "ana@example.com", "+34123456789", "XX")  # País inválido
```

## Buenas Prácticas

1. **Ser específico con las excepciones**: Captura solo las excepciones que puedes manejar adecuadamente. Evita usar `except:` sin especificar un tipo.

2. **Usar excepciones para situaciones excepcionales**: No uses excepciones para controlar el flujo normal del programa.

3. **Crear jerarquías de excepciones**: Si creas excepciones personalizadas, organízalas en una jerarquía lógica.

4. **Proporcionar mensajes útiles**: Incluye información detallada en los mensajes de error para facilitar la depuración.

5. **Limpiar recursos**: Usa `finally` o context managers (`with`) para asegurarte de que los recursos se liberen correctamente.

6. **Documentar excepciones**: Documenta qué excepciones puede lanzar una función y en qué circunstancias.

7. **No suprimir excepciones silenciosamente**: Si capturas una excepción, haz algo útil con ella. No uses bloques `except` vacíos.

8. **Usar logging en lugar de print**: En aplicaciones reales, usa el módulo `logging` para registrar errores en lugar de `print`.

9. **Considerar el rendimiento**: El manejo de excepciones tiene un costo de rendimiento. No lo uses para casos que ocurren frecuentemente en el flujo normal.

10. **Probar el manejo de excepciones**: Escribe pruebas que verifiquen que tu código maneja correctamente las excepciones.

## Recursos Adicionales

- [Documentación oficial de Python sobre excepciones](https://docs.python.org/es/3/tutorial/errors.html)
- [PEP 343 – The "with" Statement](https://peps.python.org/pep-0343/)
- [Python Exception Handling Techniques](https://realpython.com/python-exceptions/)
- [Python's contextlib module](https://docs.python.org/es/3/library/contextlib.html)
- [Logging en Python](https://docs.python.org/es/3/howto/logging.html)

---

En la siguiente sección, exploraremos la programación orientada a objetos en Python, incluyendo clases, herencia, polimorfismo y otros conceptos avanzados.