# Manejo de Archivos en Python

## Introducción

El manejo de archivos es una parte fundamental de la programación, ya que permite almacenar y recuperar datos de forma persistente. Python proporciona funciones y métodos integrados que facilitan la lectura, escritura y manipulación de archivos. En esta sección, aprenderemos cómo trabajar con archivos en Python, desde operaciones básicas hasta técnicas más avanzadas.

## Operaciones Básicas con Archivos

### Abrir y Cerrar Archivos

Para trabajar con un archivo en Python, primero debes abrirlo usando la función `open()`. Esta función devuelve un objeto de archivo que puedes usar para leer o escribir datos.

```python
# Sintaxis básica
# open(ruta_del_archivo, modo, encoding=None)

# Abrir un archivo en modo lectura (por defecto)
archivo = open('datos.txt', 'r')

# Realizar operaciones con el archivo
# ...

# Cerrar el archivo cuando hayas terminado
archivo.close()
```

Es importante cerrar los archivos después de usarlos para liberar recursos del sistema. Sin embargo, si olvidas cerrar un archivo o si ocurre una excepción antes de llegar a la línea `close()`, el archivo podría quedar abierto. Para evitar esto, Python proporciona la declaración `with` que cierra automáticamente el archivo cuando sales del bloque:

```python
# Forma recomendada de trabajar con archivos
with open('datos.txt', 'r') as archivo:
    # Realizar operaciones con el archivo
    # ...

# El archivo se cierra automáticamente al salir del bloque with
```

### Modos de Apertura de Archivos

La función `open()` acepta varios modos que determinan cómo se abrirá el archivo:

| Modo | Descripción |
|------|-------------|
| `'r'` | Lectura (por defecto). El archivo debe existir. |
| `'w'` | Escritura. Crea el archivo si no existe o lo trunca si existe. |
| `'a'` | Anexar. Abre el archivo para añadir contenido al final. Crea el archivo si no existe. |
| `'x'` | Creación exclusiva. Falla si el archivo ya existe. |
| `'b'` | Modo binario (puede combinarse con otros modos). |
| `'t'` | Modo texto (por defecto, puede combinarse con otros modos). |
| `'+'` | Actualización (lectura y escritura, puede combinarse con otros modos). |

Ejemplos:

```python
# Abrir para lectura en modo texto
with open('datos.txt', 'r') as archivo:
    # ...

# Abrir para escritura en modo texto
with open('salida.txt', 'w') as archivo:
    # ...

# Abrir para anexar en modo texto
with open('log.txt', 'a') as archivo:
    # ...

# Abrir para lectura en modo binario
with open('imagen.jpg', 'rb') as archivo:
    # ...

# Abrir para lectura y escritura
with open('datos.txt', 'r+') as archivo:
    # ...
```

### Especificar la Codificación

Cuando trabajas con archivos de texto, es importante especificar la codificación correcta para evitar problemas con caracteres especiales:

```python
# Abrir un archivo con codificación UTF-8
with open('datos.txt', 'r', encoding='utf-8') as archivo:
    # ...

# Abrir un archivo con otra codificación
with open('datos_latin.txt', 'r', encoding='latin-1') as archivo:
    # ...
```

## Lectura de Archivos

Python ofrece varios métodos para leer datos de un archivo:

### Leer Todo el Contenido

```python
with open('datos.txt', 'r', encoding='utf-8') as archivo:
    contenido = archivo.read()  # Lee todo el archivo como una cadena
    print(contenido)
```

### Leer un Número Específico de Caracteres

```python
with open('datos.txt', 'r', encoding='utf-8') as archivo:
    primeros_10 = archivo.read(10)  # Lee los primeros 10 caracteres
    print(primeros_10)
    
    siguientes_10 = archivo.read(10)  # Lee los siguientes 10 caracteres
    print(siguientes_10)
```

### Leer Línea por Línea

```python
# Leer una línea
with open('datos.txt', 'r', encoding='utf-8') as archivo:
    primera_linea = archivo.readline()  # Lee la primera línea
    print(primera_linea)
    
    segunda_linea = archivo.readline()  # Lee la segunda línea
    print(segunda_linea)

# Leer todas las líneas en una lista
with open('datos.txt', 'r', encoding='utf-8') as archivo:
    lineas = archivo.readlines()  # Devuelve una lista de líneas
    print(lineas)

# Iterar sobre las líneas (forma eficiente para archivos grandes)
with open('datos.txt', 'r', encoding='utf-8') as archivo:
    for linea in archivo:  # El archivo es un iterable
        print(linea, end='')  # end='' evita doble salto de línea
```

## Escritura en Archivos

Python proporciona métodos para escribir datos en archivos:

### Escribir Texto

```python
# Escribir una cadena
with open('salida.txt', 'w', encoding='utf-8') as archivo:
    archivo.write("Hola, mundo!\n")  # \n añade un salto de línea
    archivo.write("Esta es otra línea.\n")

# Escribir varias líneas a la vez
with open('salida.txt', 'w', encoding='utf-8') as archivo:
    lineas = ["Línea 1\n", "Línea 2\n", "Línea 3\n"]
    archivo.writelines(lineas)  # writelines no añade saltos de línea automáticamente
```

### Anexar Texto

```python
# Añadir texto al final del archivo
with open('log.txt', 'a', encoding='utf-8') as archivo:
    archivo.write("Nuevo registro: " + str(datetime.now()) + "\n")
```

## Posicionamiento en Archivos

Puedes controlar la posición actual dentro de un archivo usando los métodos `seek()` y `tell()`:

```python
with open('datos.txt', 'r', encoding='utf-8') as archivo:
    # Obtener la posición actual
    posicion = archivo.tell()
    print(f"Posición inicial: {posicion}")  # Normalmente 0
    
    # Leer algunos datos
    datos = archivo.read(10)
    print(f"Datos leídos: {datos}")
    
    # Nueva posición
    posicion = archivo.tell()
    print(f"Nueva posición: {posicion}")  # Ahora es 10
    
    # Mover a una posición específica desde el inicio
    archivo.seek(0)  # Volver al inicio
    print(f"Volvemos a la posición: {archivo.tell()}")  # Ahora es 0
    
    # Leer de nuevo
    datos = archivo.read(5)
    print(f"Datos leídos: {datos}")  # Primeros 5 caracteres
    
    # Mover relativamente a la posición actual
    archivo.seek(5, 1)  # Avanzar 5 posiciones desde la posición actual
    print(f"Nueva posición: {archivo.tell()}")  # Ahora es 10
    
    # Mover relativamente al final del archivo
    archivo.seek(0, 2)  # Ir al final del archivo
    print(f"Posición final: {archivo.tell()}")  # Tamaño del archivo
```

El método `seek(offset, whence)` acepta dos parámetros:
- `offset`: El número de bytes a mover
- `whence`: Punto de referencia (opcional)
  - 0: Desde el inicio del archivo (por defecto)
  - 1: Desde la posición actual
  - 2: Desde el final del archivo

## Manejo de Archivos Binarios

Los archivos binarios contienen datos que no son texto, como imágenes, audio, video o datos estructurados. Para trabajar con ellos, debes usar el modo binario (`'b'`):

```python
# Leer un archivo binario
with open('imagen.jpg', 'rb') as archivo:
    datos = archivo.read()  # Lee todos los bytes
    print(f"Tamaño: {len(datos)} bytes")

# Escribir un archivo binario
with open('copia.jpg', 'wb') as archivo:
    archivo.write(datos)  # Escribe todos los bytes
```

### Trabajar con Datos Estructurados

Para trabajar con datos estructurados en archivos binarios, puedes usar el módulo `struct`:

```python
import struct

# Escribir datos estructurados
with open('datos.bin', 'wb') as archivo:
    # Empaquetar un entero, un float y una cadena de 5 caracteres
    # 'i' para entero, 'f' para float, '5s' para cadena de 5 bytes
    datos = struct.pack('if5s', 42, 3.14, b'hello')
    archivo.write(datos)

# Leer datos estructurados
with open('datos.bin', 'rb') as archivo:
    datos = archivo.read()
    # Desempaquetar los datos
    entero, flotante, cadena = struct.unpack('if5s', datos)
    print(f"Entero: {entero}, Flotante: {flotante}, Cadena: {cadena.decode('ascii')}")
```

## Trabajar con Rutas de Archivos

Python proporciona el módulo `os.path` y el módulo `pathlib` para trabajar con rutas de archivos de manera compatible con diferentes sistemas operativos:

### Usando `os.path`

```python
import os.path

# Unir componentes de ruta
ruta = os.path.join('directorio', 'subdirectorio', 'archivo.txt')
print(ruta)  # 'directorio/subdirectorio/archivo.txt' en Unix/Mac, 'directorio\subdirectorio\archivo.txt' en Windows

# Obtener el directorio y el nombre de archivo
directorio = os.path.dirname(ruta)
nombre_archivo = os.path.basename(ruta)
print(f"Directorio: {directorio}, Archivo: {nombre_archivo}")

# Comprobar si una ruta existe
existe = os.path.exists(ruta)
print(f"¿Existe? {existe}")

# Comprobar si es un archivo o un directorio
es_archivo = os.path.isfile(ruta)
es_directorio = os.path.isdir(ruta)
print(f"¿Es archivo? {es_archivo}, ¿Es directorio? {es_directorio}")

# Obtener el tamaño de un archivo
tamaño = os.path.getsize(ruta)  # en bytes
print(f"Tamaño: {tamaño} bytes")

# Obtener la ruta absoluta
ruta_absoluta = os.path.abspath(ruta)
print(f"Ruta absoluta: {ruta_absoluta}")
```

### Usando `pathlib` (Python 3.4+)

```python
from pathlib import Path

# Crear un objeto Path
ruta = Path('directorio') / 'subdirectorio' / 'archivo.txt'
print(ruta)  # directorio/subdirectorio/archivo.txt (representación independiente del sistema)

# Obtener el directorio y el nombre de archivo
directorio = ruta.parent
nombre_archivo = ruta.name
print(f"Directorio: {directorio}, Archivo: {nombre_archivo}")

# Comprobar si una ruta existe
existe = ruta.exists()
print(f"¿Existe? {existe}")

# Comprobar si es un archivo o un directorio
es_archivo = ruta.is_file()
es_directorio = ruta.is_dir()
print(f"¿Es archivo? {es_archivo}, ¿Es directorio? {es_directorio}")

# Obtener el tamaño de un archivo
tamaño = ruta.stat().st_size  # en bytes
print(f"Tamaño: {tamaño} bytes")

# Obtener la ruta absoluta
ruta_absoluta = ruta.absolute()
print(f"Ruta absoluta: {ruta_absoluta}")

# Leer y escribir archivos directamente con Path
contenido = ruta.read_text(encoding='utf-8')  # Lee todo el contenido como texto
ruta.write_text("Nuevo contenido", encoding='utf-8')  # Escribe texto en el archivo
```

## Operaciones con Directorios

Python proporciona funciones para trabajar con directorios:

### Usando `os`

```python
import os

# Obtener el directorio de trabajo actual
directorio_actual = os.getcwd()
print(f"Directorio actual: {directorio_actual}")

# Cambiar el directorio de trabajo
os.chdir('/ruta/a/otro/directorio')

# Listar archivos y directorios
contenido = os.listdir('.')  # '.' representa el directorio actual
print(f"Contenido: {contenido}")

# Crear un directorio
os.mkdir('nuevo_directorio')  # Crea un solo directorio
os.makedirs('ruta/a/nuevos/directorios', exist_ok=True)  # Crea directorios anidados

# Eliminar un directorio
os.rmdir('directorio_vacio')  # Solo funciona si el directorio está vacío
import shutil
shutil.rmtree('directorio_con_contenido')  # Elimina el directorio y todo su contenido

# Renombrar un archivo o directorio
os.rename('viejo_nombre.txt', 'nuevo_nombre.txt')
```

### Usando `pathlib`

```python
from pathlib import Path

# Obtener el directorio de trabajo actual
directorio_actual = Path.cwd()
print(f"Directorio actual: {directorio_actual}")

# Listar archivos y directorios
ruta = Path('.')
for elemento in ruta.iterdir():
    print(elemento)

# Filtrar por tipo o patrón
archivos = [p for p in ruta.iterdir() if p.is_file()]
print(f"Archivos: {archivos}")

python_files = list(ruta.glob('*.py'))  # Archivos .py en el directorio actual
print(f"Archivos Python: {python_files}")

all_python_files = list(ruta.rglob('*.py'))  # Archivos .py en el directorio actual y subdirectorios
print(f"Todos los archivos Python: {all_python_files}")

# Crear un directorio
nuevo_dir = Path('nuevo_directorio')
nuevo_dir.mkdir(exist_ok=True)  # exist_ok=True evita errores si ya existe

# Crear directorios anidados
nuevos_dirs = Path('ruta/a/nuevos/directorios')
nuevos_dirs.mkdir(parents=True, exist_ok=True)  # parents=True crea directorios padres

# Eliminar un archivo
archivo = Path('archivo_a_eliminar.txt')
if archivo.exists():
    archivo.unlink()

# Renombrar un archivo
archivo = Path('viejo_nombre.txt')
archivo.rename('nuevo_nombre.txt')
```

## Formatos de Archivo Comunes

Python proporciona módulos para trabajar con formatos de archivo comunes:

### CSV (Valores Separados por Comas)

```python
import csv

# Escribir un archivo CSV
with open('datos.csv', 'w', newline='', encoding='utf-8') as archivo:
    escritor = csv.writer(archivo)
    escritor.writerow(['Nombre', 'Edad', 'Ciudad'])  # Encabezados
    escritor.writerow(['Ana', 25, 'Madrid'])
    escritor.writerow(['Carlos', 30, 'Barcelona'])
    escritor.writerow(['Elena', 28, 'Valencia'])

# Leer un archivo CSV
with open('datos.csv', 'r', newline='', encoding='utf-8') as archivo:
    lector = csv.reader(archivo)
    encabezados = next(lector)  # Leer encabezados
    print(f"Encabezados: {encabezados}")
    
    for fila in lector:
        print(f"Nombre: {fila[0]}, Edad: {fila[1]}, Ciudad: {fila[2]}")

# Usar DictReader y DictWriter para trabajar con diccionarios
with open('datos.csv', 'r', newline='', encoding='utf-8') as archivo:
    lector = csv.DictReader(archivo)
    for fila in lector:
        print(f"Nombre: {fila['Nombre']}, Edad: {fila['Edad']}, Ciudad: {fila['Ciudad']}")

with open('nuevos_datos.csv', 'w', newline='', encoding='utf-8') as archivo:
    encabezados = ['Nombre', 'Edad', 'Ciudad']
    escritor = csv.DictWriter(archivo, fieldnames=encabezados)
    escritor.writeheader()  # Escribir encabezados
    escritor.writerow({'Nombre': 'Ana', 'Edad': 25, 'Ciudad': 'Madrid'})
    escritor.writerow({'Nombre': 'Carlos', 'Edad': 30, 'Ciudad': 'Barcelona'})
```

### JSON (JavaScript Object Notation)

```python
import json

# Datos Python
datos = {
    'nombre': 'Ana',
    'edad': 25,
    'ciudad': 'Madrid',
    'intereses': ['programación', 'música', 'viajes'],
    'activo': True,
    'altura': 1.65
}

# Escribir JSON
with open('datos.json', 'w', encoding='utf-8') as archivo:
    json.dump(datos, archivo, ensure_ascii=False, indent=4)
    # ensure_ascii=False permite caracteres no ASCII
    # indent=4 formatea el JSON con 4 espacios de indentación

# Convertir a cadena JSON
json_str = json.dumps(datos, ensure_ascii=False, indent=4)
print(json_str)

# Leer JSON
with open('datos.json', 'r', encoding='utf-8') as archivo:
    datos_leidos = json.load(archivo)
    print(datos_leidos['nombre'])  # Ana
    print(datos_leidos['intereses'])  # ['programación', 'música', 'viajes']

# Convertir cadena JSON a objeto Python
json_str = '{"nombre": "Carlos", "edad": 30}'
datos_desde_str = json.loads(json_str)
print(datos_desde_str['nombre'])  # Carlos
```

### Pickle (Serialización de Objetos Python)

```python
import pickle

# Datos Python (puede ser cualquier objeto serializable)
class Persona:
    def __init__(self, nombre, edad):
        self.nombre = nombre
        self.edad = edad
    
    def __str__(self):
        return f"{self.nombre}, {self.edad} años"

personas = [
    Persona("Ana", 25),
    Persona("Carlos", 30),
    Persona("Elena", 28)
]

# Serializar (guardar) objetos
with open('personas.pickle', 'wb') as archivo:
    pickle.dump(personas, archivo)

# Deserializar (cargar) objetos
with open('personas.pickle', 'rb') as archivo:
    personas_cargadas = pickle.load(archivo)
    for persona in personas_cargadas:
        print(persona)  # Usa el método __str__ de la clase Persona
```

> **Nota de seguridad**: Nunca deserialices datos pickle de fuentes no confiables, ya que puede ejecutar código arbitrario.

## Manejo de Excepciones en Operaciones de Archivo

Las operaciones de archivo pueden generar varias excepciones. Es importante manejarlas adecuadamente:

```python
try:
    with open('archivo_inexistente.txt', 'r') as archivo:
        contenido = archivo.read()
except FileNotFoundError:
    print("El archivo no existe.")
except PermissionError:
    print("No tienes permiso para acceder al archivo.")
except IsADirectoryError:
    print("La ruta especificada es un directorio, no un archivo.")
except UnicodeDecodeError:
    print("No se puede decodificar el archivo con la codificación especificada.")
except Exception as e:
    print(f"Ocurrió un error inesperado: {e}")
```

## Ejemplo Práctico: Analizador de Logs

Vamos a crear un programa que analice un archivo de log y genere un informe:

```python
import re
from collections import Counter
from datetime import datetime

def analizar_log(ruta_log):
    # Patrones para extraer información
    patron_fecha = r'\[(\d{2}/\w+/\d{4}:\d{2}:\d{2}:\d{2})\]'
    patron_ip = r'(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})'
    patron_metodo = r'"(GET|POST|PUT|DELETE)'
    patron_url = r'\s(/[^\s]*)'
    patron_codigo = r'\s(\d{3})\s'
    
    # Contadores
    ips = Counter()
    metodos = Counter()
    urls = Counter()
    codigos = Counter()
    horas = Counter()
    
    # Leer el archivo de log
    try:
        with open(ruta_log, 'r', encoding='utf-8') as archivo:
            for linea in archivo:
                # Extraer información
                match_ip = re.search(patron_ip, linea)
                if match_ip:
                    ips[match_ip.group(1)] += 1
                
                match_fecha = re.search(patron_fecha, linea)
                if match_fecha:
                    fecha_str = match_fecha.group(1)
                    try:
                        fecha = datetime.strptime(fecha_str, '%d/%b/%Y:%H:%M:%S')
                        horas[fecha.hour] += 1
                    except ValueError:
                        pass
                
                match_metodo = re.search(patron_metodo, linea)
                if match_metodo:
                    metodos[match_metodo.group(1)] += 1
                
                match_url = re.search(patron_url, linea)
                if match_url:
                    urls[match_url.group(1)] += 1
                
                match_codigo = re.search(patron_codigo, linea)
                if match_codigo:
                    codigos[match_codigo.group(1)] += 1
        
        # Generar informe
        print("=== INFORME DE ANÁLISIS DE LOG ===")
        print(f"\nTotal de líneas procesadas: {sum(ips.values())}")
        
        print("\nIPs más frecuentes:")
        for ip, count in ips.most_common(5):
            print(f"  {ip}: {count}")
        
        print("\nMétodos HTTP:")
        for metodo, count in metodos.most_common():
            print(f"  {metodo}: {count}")
        
        print("\nCódigos de respuesta:")
        for codigo, count in codigos.most_common():
            print(f"  {codigo}: {count}")
        
        print("\nURLs más solicitadas:")
        for url, count in urls.most_common(5):
            print(f"  {url}: {count}")
        
        print("\nDistribución por hora:")
        for hora in sorted(horas.keys()):
            print(f"  {hora:02d}:00 - {hora+1:02d}:00: {horas[hora]}")
        
        # Guardar informe en un archivo
        with open('informe_log.txt', 'w', encoding='utf-8') as informe:
            informe.write("=== INFORME DE ANÁLISIS DE LOG ===\n")
            informe.write(f"\nTotal de líneas procesadas: {sum(ips.values())}\n")
            
            informe.write("\nIPs más frecuentes:\n")
            for ip, count in ips.most_common(5):
                informe.write(f"  {ip}: {count}\n")
            
            # ... (escribir el resto del informe)
        
        print("\nInforme guardado en 'informe_log.txt'")
    
    except FileNotFoundError:
        print(f"Error: El archivo '{ruta_log}' no existe.")
    except Exception as e:
        print(f"Error inesperado: {e}")

# Uso
if __name__ == "__main__":
    ruta_log = input("Introduce la ruta del archivo de log: ")
    analizar_log(ruta_log)
```

## Buenas Prácticas

1. **Usar `with` para manejar archivos**: Siempre usa la declaración `with` para asegurarte de que los archivos se cierren correctamente, incluso si ocurren excepciones.

2. **Especificar la codificación**: Siempre especifica la codificación al abrir archivos de texto para evitar problemas con caracteres especiales.

3. **Manejar excepciones**: Anticipa y maneja las excepciones que pueden ocurrir durante las operaciones de archivo.

4. **Usar rutas relativas con cuidado**: Las rutas relativas dependen del directorio de trabajo actual. Si no estás seguro, usa rutas absolutas o `pathlib` para construir rutas de manera segura.

5. **Limitar el tamaño de lectura**: Para archivos grandes, lee por partes o línea por línea en lugar de cargar todo el archivo en memoria.

6. **Usar los módulos adecuados**: Utiliza módulos especializados como `csv`, `json` o `pickle` para formatos específicos en lugar de implementar tu propia lógica de análisis.

7. **Verificar permisos y existencia**: Verifica si tienes permisos para acceder a un archivo y si existe antes de intentar operaciones que puedan fallar.

8. **Usar nombres de archivo seguros**: Evita caracteres especiales o espacios en nombres de archivo para garantizar la compatibilidad entre sistemas.

9. **Hacer copias de seguridad**: Antes de modificar archivos importantes, haz una copia de seguridad.

10. **Usar `pathlib` para operaciones modernas**: Para Python 3.4+, considera usar `pathlib` en lugar de `os.path` para un manejo más orientado a objetos y legible.

## Recursos Adicionales

- [Documentación oficial de Python sobre archivos](https://docs.python.org/es/3/tutorial/inputoutput.html#reading-and-writing-files)
- [Documentación del módulo `os.path`](https://docs.python.org/es/3/library/os.path.html)
- [Documentación del módulo `pathlib`](https://docs.python.org/es/3/library/pathlib.html)
- [Documentación del módulo `csv`](https://docs.python.org/es/3/library/csv.html)
- [Documentación del módulo `json`](https://docs.python.org/es/3/library/json.html)
- [Documentación del módulo `pickle`](https://docs.python.org/es/3/library/pickle.html)

---

En la siguiente sección, exploraremos la programación orientada a objetos en Python, que nos permitirá crear estructuras de datos personalizadas y organizar nuestro código de manera más modular y reutilizable.