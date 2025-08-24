# Introducción a Python

## ¿Qué es Python?

Python es un lenguaje de programación interpretado, de alto nivel y propósito general creado por Guido van Rossum y lanzado por primera vez en 1991. Python se destaca por su filosofía de diseño que enfatiza la legibilidad del código mediante el uso significativo de espacios en blanco (indentación).

## Características principales

- **Fácil de aprender y usar**: Sintaxis clara y legible, similar al pseudocódigo.
- **Interpretado**: No necesita ser compilado antes de ejecutarse.
- **Tipado dinámico**: No es necesario declarar el tipo de las variables.
- **Alto nivel**: Abstrae muchos detalles complejos de la máquina.
- **Multiplataforma**: Funciona en Windows, macOS, Linux y otros sistemas.
- **Multiparadigma**: Soporta programación orientada a objetos, imperativa y funcional.
- **Extensa biblioteca estándar**: "Baterías incluidas" para muchas tareas comunes.

## Instalación

### Windows

1. Visita [python.org](https://www.python.org/downloads/)
2. Descarga la última versión estable
3. Ejecuta el instalador y marca la opción "Add Python to PATH"
4. Verifica la instalación abriendo una terminal y escribiendo:
   ```
   python --version
   ```

### macOS

1. Instala Homebrew (si no lo tienes):
   ```
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```
2. Instala Python:
   ```
   brew install python
   ```
3. Verifica la instalación:
   ```
   python3 --version
   ```

### Linux

En la mayoría de las distribuciones Linux, Python ya viene preinstalado. Si no es así:

**Ubuntu/Debian**:
```
sudo apt update
sudo apt install python3 python3-pip
```

**Fedora**:
```
sudo dnf install python3 python3-pip
```

## Tu primer programa en Python

Es tradición comenzar con un programa "Hola Mundo". Crea un archivo llamado `hola_mundo.py` con el siguiente contenido:

```python
print("¡Hola, Mundo!")
```

Ejecuta el programa desde la terminal:

```
python hola_mundo.py
```

Deberías ver el mensaje "¡Hola, Mundo!" en la pantalla.

## El intérprete interactivo de Python

Python incluye un intérprete interactivo que te permite ejecutar código línea por línea. Para iniciarlo, simplemente escribe `python` o `python3` en la terminal:

```
$ python
Python 3.9.5 (default, May 11 2021, 08:20:37)
[GCC 10.3.0] on linux
Type "help", "copyright", "credits" or "license" for more information.
>>>
```

Ahora puedes escribir código Python directamente:

```python
>>> print("Hola desde el intérprete")
Hola desde el intérprete
>>> 2 + 3
5
>>> nombre = "Python"
>>> print(f"Me encanta {nombre}")
Me encanta Python
```

Para salir del intérprete, escribe `exit()` o presiona Ctrl+D (en Unix/macOS) o Ctrl+Z seguido de Enter (en Windows).

## Entornos de Desarrollo Integrados (IDEs)

Aunque puedes escribir código Python en cualquier editor de texto, un buen IDE puede mejorar significativamente tu productividad:

- **Visual Studio Code**: Ligero, gratuito, con excelente soporte para Python mediante extensiones.
- **PyCharm**: IDE completo específico para Python (versiones Community y Professional).
- **Jupyter Notebook**: Ideal para ciencia de datos y aprendizaje.
- **Spyder**: Diseñado para científicos, ingenieros y analistas de datos.
- **Thonny**: Excelente para principiantes, con un depurador simple.

## Zen de Python

El Zen de Python es una colección de 19 principios que influyen en el diseño del lenguaje. Puedes verlo escribiendo:

```python
import this
```

Algunos de estos principios incluyen:

- Hermoso es mejor que feo.
- Explícito es mejor que implícito.
- Simple es mejor que complejo.
- Complejo es mejor que complicado.
- La legibilidad cuenta.

## Recursos adicionales

- [Documentación oficial de Python](https://docs.python.org/es/3/)
- [Tutorial oficial de Python](https://docs.python.org/es/3/tutorial/index.html)
- [Real Python](https://realpython.com/) - Tutoriales de alta calidad
- [Python para todos](https://www.py4e.com/) - Curso gratuito de Charles Severance

---

En la siguiente sección, aprenderemos sobre variables, tipos de datos y operadores en Python.