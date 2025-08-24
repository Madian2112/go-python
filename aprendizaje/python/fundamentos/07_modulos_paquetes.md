# Módulos y Paquetes en Python

## Introducción

Los módulos y paquetes son fundamentales para organizar y reutilizar código en Python. Permiten dividir programas grandes en componentes más pequeños y manejables, facilitando el mantenimiento y la colaboración. En esta sección, aprenderemos cómo crear, importar y utilizar módulos y paquetes en Python.

## Módulos

Un módulo es simplemente un archivo Python (`.py`) que contiene definiciones de funciones, clases y variables que pueden ser utilizadas en otros programas Python.

### Creación de un módulo

Crear un módulo es tan sencillo como crear un archivo Python con el código que deseas reutilizar. Por ejemplo, podemos crear un módulo llamado `matematicas.py` con funciones matemáticas:

```python
# matematicas.py

def sumar(a, b):
    """Suma dos números y devuelve el resultado."""
    return a + b

def restar(a, b):
    """Resta b de a y devuelve el resultado."""
    return a - b

def multiplicar(a, b):
    """Multiplica dos números y devuelve el resultado."""
    return a * b

def dividir(a, b):
    """Divide a por b y devuelve el resultado.
    Lanza una excepción ZeroDivisionError si b es 0.
    """
    if b == 0:
        raise ZeroDivisionError("No se puede dividir por cero")
    return a / b

# Variable a nivel de módulo
PI = 3.14159265359

# Función que se ejecutará si el módulo se ejecuta directamente
def main():
    print(f"Suma: {sumar(10, 5)}")
    print(f"Resta: {restar(10, 5)}")
    print(f"Multiplicación: {multiplicar(10, 5)}")
    print(f"División: {dividir(10, 5)}")
    print(f"PI: {PI}")

# Bloque condicional para ejecutar código solo si el módulo se ejecuta directamente
if __name__ == "__main__":
    main()
```

### Importación de módulos

Existen varias formas de importar módulos en Python:

#### 1. Importar el módulo completo

```python
import matematicas

# Usar funciones del módulo con notación de punto
resultado = matematicas.sumar(10, 5)
print(resultado)  # 15
print(matematicas.PI)  # 3.14159265359
```

#### 2. Importar elementos específicos

```python
from matematicas import sumar, restar, PI

# Usar funciones directamente sin notación de punto
resultado = sumar(10, 5)
print(resultado)  # 15
print(PI)  # 3.14159265359
```

#### 3. Importar todo el contenido

```python
from matematicas import *

# Usar todas las funciones directamente
print(sumar(10, 5))  # 15
print(restar(10, 5))  # 5
print(PI)  # 3.14159265359
```

> **Nota**: Importar todo con `*` generalmente no se recomienda en código de producción, ya que puede causar conflictos de nombres y hacer que el código sea menos legible.

#### 4. Importar con alias

```python
import matematicas as mat

# Usar el alias para acceder a las funciones
print(mat.sumar(10, 5))  # 15
print(mat.PI)  # 3.14159265359
```

### El bloque `if __name__ == "__main__":`

Este bloque condicional es muy común en los módulos Python. Permite que el código dentro del bloque se ejecute solo cuando el módulo se ejecuta directamente como un script, pero no cuando se importa desde otro módulo.

```python
# En matematicas.py
if __name__ == "__main__":
    main()
```

Si ejecutamos `python matematicas.py`, el bloque `main()` se ejecutará.
Si importamos `matematicas` desde otro script, el bloque `main()` no se ejecutará.

### Localización de módulos

Python busca módulos en las siguientes ubicaciones (en orden):

1. El directorio actual
2. Los directorios listados en la variable de entorno `PYTHONPATH`
3. Los directorios de instalación estándar de Python

Puedes ver la lista de directorios de búsqueda en la variable `sys.path`:

```python
import sys
print(sys.path)
```

## Paquetes

Un paquete es una forma de organizar módulos relacionados en una estructura de directorios jerárquica. Es simplemente un directorio que contiene módulos Python y un archivo especial llamado `__init__.py`.

### Estructura de un paquete

Un paquete típico podría tener una estructura como esta:

```
mi_paquete/
    __init__.py
    modulo1.py
    modulo2.py
    subpaquete/
        __init__.py
        modulo3.py
        modulo4.py
```

### El archivo `__init__.py`

El archivo `__init__.py` es necesario para que Python trate un directorio como un paquete. Puede estar vacío o contener código de inicialización para el paquete.

Ejemplo de un archivo `__init__.py` simple:

```python
# mi_paquete/__init__.py

# Importar elementos específicos para que estén disponibles directamente desde el paquete
from .modulo1 import funcion1, funcion2
from .modulo2 import Clase1

# Definir qué módulos se importarán con 'from mi_paquete import *'
__all__ = ['modulo1', 'modulo2']

# Variables a nivel de paquete
__version__ = '0.1'
```

### Importación de paquetes

Hay varias formas de importar módulos desde paquetes:

#### 1. Importar un módulo específico del paquete

```python
import mi_paquete.modulo1

# Usar funciones del módulo
mi_paquete.modulo1.funcion1()
```

#### 2. Importar elementos específicos de un módulo

```python
from mi_paquete.modulo1 import funcion1, funcion2

# Usar funciones directamente
funcion1()
```

#### 3. Importar un submódulo

```python
import mi_paquete.subpaquete.modulo3

# Usar funciones del submódulo
mi_paquete.subpaquete.modulo3.funcion3()
```

#### 4. Importar con alias

```python
import mi_paquete.subpaquete.modulo3 as mod3

# Usar el alias
mod3.funcion3()
```

### Importaciones relativas

Dentro de un paquete, puedes usar importaciones relativas para referirte a módulos en el mismo paquete o en paquetes padres:

```python
# En mi_paquete/subpaquete/modulo3.py

# Importación relativa desde el mismo paquete
from . import modulo4

# Importación relativa desde un módulo en el mismo paquete
from .modulo4 import funcion4

# Importación relativa desde el paquete padre
from .. import modulo1

# Importación relativa desde un módulo en el paquete padre
from ..modulo1 import funcion1
```

## Ejemplo Práctico: Creación de un Paquete

Vamos a crear un paquete llamado `utilidades` con varios módulos:

### Estructura del paquete

```
utilidades/
    __init__.py
    matematicas.py
    texto.py
    archivos.py
```

### Contenido de los archivos

```python
# utilidades/__init__.py

__version__ = '0.1'
__author__ = 'Tu Nombre'

# Importaciones para facilitar el acceso
from .matematicas import sumar, restar, multiplicar, dividir
from .texto import invertir, contar_palabras

# Definir qué módulos se importarán con 'from utilidades import *'
__all__ = ['matematicas', 'texto', 'archivos']
```

```python
# utilidades/matematicas.py

def sumar(a, b):
    return a + b

def restar(a, b):
    return a - b

def multiplicar(a, b):
    return a * b

def dividir(a, b):
    if b == 0:
        raise ZeroDivisionError("No se puede dividir por cero")
    return a / b
```

```python
# utilidades/texto.py

def invertir(texto):
    """Invierte una cadena de texto."""
    return texto[::-1]

def contar_palabras(texto):
    """Cuenta el número de palabras en un texto."""
    return len(texto.split())

def es_palindromo(texto):
    """Verifica si un texto es un palíndromo (se lee igual al derecho y al revés)."""
    # Eliminar espacios y convertir a minúsculas
    texto = texto.lower().replace(" ", "")
    return texto == texto[::-1]
```

```python
# utilidades/archivos.py

import os

def leer_archivo(ruta):
    """Lee el contenido de un archivo y lo devuelve como una cadena."""
    with open(ruta, 'r', encoding='utf-8') as archivo:
        return archivo.read()

def escribir_archivo(ruta, contenido):
    """Escribe contenido en un archivo."""
    with open(ruta, 'w', encoding='utf-8') as archivo:
        archivo.write(contenido)

def listar_archivos(directorio):
    """Lista todos los archivos en un directorio."""
    return [f for f in os.listdir(directorio) if os.path.isfile(os.path.join(directorio, f))]
```

### Uso del paquete

```python
# Importar el paquete
import utilidades

# Usar funciones importadas en __init__.py
print(utilidades.sumar(10, 5))  # 15
print(utilidades.invertir("Hola"))  # aloH

# Importar módulos específicos
from utilidades import archivos

# Listar archivos en el directorio actual
archivos_actuales = archivos.listar_archivos('.')
print(archivos_actuales)

# Importar funciones específicas
from utilidades.texto import es_palindromo

# Verificar palíndromos
print(es_palindromo("anita lava la tina"))  # True
print(es_palindromo("python"))  # False

# Verificar la versión del paquete
print(utilidades.__version__)  # 0.1
```

## Distribución de Paquetes

Para compartir tu paquete con otros, puedes distribuirlo de varias formas:

### 1. Instalación en modo desarrollo

Para instalar tu paquete en modo desarrollo (los cambios en el código se reflejan inmediatamente):

```bash
pip install -e /ruta/a/tu/paquete
```

### 2. Creación de un paquete distribuible

Para crear un paquete distribuible, necesitas un archivo `setup.py`:

```python
# setup.py
from setuptools import setup, find_packages

setup(
    name="utilidades",
    version="0.1",
    packages=find_packages(),
    author="Tu Nombre",
    author_email="tu@email.com",
    description="Un paquete de utilidades para Python",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/tuusuario/utilidades",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)
```

Luego, puedes crear un archivo distribuible:

```bash
python setup.py sdist bdist_wheel
```

Esto creará archivos en los directorios `dist/` que puedes compartir o subir a PyPI.

### 3. Publicación en PyPI

Para publicar tu paquete en el Python Package Index (PyPI):

```bash
pip install twine
twine upload dist/*
```

## Buenas Prácticas

1. **Organización clara**: Organiza tu código en módulos y paquetes lógicos basados en funcionalidad.

2. **Nombres descriptivos**: Usa nombres claros y descriptivos para tus módulos y paquetes.

3. **Documentación**: Incluye docstrings en tus módulos, clases y funciones para explicar su propósito y uso.

4. **Importaciones específicas**: Prefiere importar elementos específicos (`from modulo import funcion`) en lugar de importar todo (`from modulo import *`).

5. **Evita la circularidad**: Evita las importaciones circulares entre módulos, ya que pueden causar problemas difíciles de depurar.

6. **Usa `__all__`**: Define la variable `__all__` en tus módulos y paquetes para controlar qué se importa con `from modulo import *`.

7. **Pruebas**: Incluye pruebas unitarias para tus módulos y paquetes para asegurar su correcto funcionamiento.

8. **Versionado**: Usa versionado semántico para tus paquetes (MAJOR.MINOR.PATCH).

## Recursos Adicionales

- [Documentación oficial de Python sobre módulos](https://docs.python.org/es/3/tutorial/modules.html)
- [Documentación oficial de Python sobre paquetes](https://docs.python.org/es/3/tutorial/modules.html#packages)
- [Python Packaging User Guide](https://packaging.python.org/)
- [Setuptools Documentation](https://setuptools.readthedocs.io/)
- [PyPI - Python Package Index](https://pypi.org/)

---

En la siguiente sección, exploraremos el manejo de archivos y directorios en Python, que nos permitirá leer y escribir datos en archivos, navegar por el sistema de archivos y trabajar con diferentes formatos de archivo.