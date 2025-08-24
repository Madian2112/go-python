# Testing Avanzado en Python

## Introducción

El testing es una parte fundamental del desarrollo de software profesional. Python cuenta con un rico ecosistema de herramientas y frameworks para realizar pruebas efectivas. En este módulo, exploraremos técnicas avanzadas de testing en Python, incluyendo pruebas unitarias, mocking, fixtures, pruebas parametrizadas y cobertura de código.

## Repaso de Conceptos Básicos

### Unittest

El módulo `unittest` es parte de la biblioteca estándar de Python y proporciona un framework para organizar y ejecutar pruebas.

```python
import unittest

def suma(a, b):
    return a + b

class TestSuma(unittest.TestCase):
    def test_suma_positivos(self):
        self.assertEqual(suma(2, 3), 5)
    
    def test_suma_negativos(self):
        self.assertEqual(suma(-1, -1), -2)
    
    def test_suma_mixto(self):
        self.assertEqual(suma(-1, 1), 0)

if __name__ == '__main__':
    unittest.main()
```

### Pytest

`pytest` es un framework de testing más moderno y flexible que simplifica la escritura de pruebas.

```python
# test_suma.py
def suma(a, b):
    return a + b

def test_suma_positivos():
    assert suma(2, 3) == 5

def test_suma_negativos():
    assert suma(-1, -1) == -2

def test_suma_mixto():
    assert suma(-1, 1) == 0
```

## Pytest Avanzado

### Fixtures

Las fixtures en pytest permiten configurar estados o recursos compartidos entre pruebas.

```python
import pytest
import tempfile
import os

@pytest.fixture
def archivo_temporal():
    # Setup
    fd, path = tempfile.mkstemp()
    with os.fdopen(fd, 'w') as f:
        f.write("contenido de prueba")
    
    # Proporcionar el recurso
    yield path
    
    # Teardown
    os.unlink(path)

def test_leer_archivo(archivo_temporal):
    with open(archivo_temporal, 'r') as f:
        contenido = f.read()
    assert contenido == "contenido de prueba"

def test_escribir_archivo(archivo_temporal):
    with open(archivo_temporal, 'a') as f:
        f.write("\nmás contenido")
    
    with open(archivo_temporal, 'r') as f:
        contenido = f.read()
    
    assert "más contenido" in contenido
```

### Fixtures con Alcance

Puedes definir el alcance de una fixture para controlar cuándo se configura y destruye.

```python
import pytest
import sqlite3

@pytest.fixture(scope="module")
def db_connection():
    # Setup - se ejecuta una vez por módulo
    conn = sqlite3.connect(':memory:')
    conn.execute("CREATE TABLE usuarios (id INTEGER PRIMARY KEY, nombre TEXT)")
    conn.execute("INSERT INTO usuarios VALUES (1, 'Alice')")
    conn.execute("INSERT INTO usuarios VALUES (2, 'Bob')")
    conn.commit()
    
    yield conn
    
    # Teardown - se ejecuta una vez al final del módulo
    conn.close()

def test_consulta_usuario(db_connection):
    cursor = db_connection.execute("SELECT nombre FROM usuarios WHERE id = 1")
    resultado = cursor.fetchone()
    assert resultado[0] == 'Alice'

def test_insertar_usuario(db_connection):
    db_connection.execute("INSERT INTO usuarios VALUES (3, 'Charlie')")
    db_connection.commit()
    
    cursor = db_connection.execute("SELECT nombre FROM usuarios WHERE id = 3")
    resultado = cursor.fetchone()
    assert resultado[0] == 'Charlie'
```

### Fixtures Parametrizadas

Las fixtures pueden ser parametrizadas para proporcionar múltiples valores.

```python
import pytest

@pytest.fixture(params=[1, 2, 3])
def valor_prueba(request):
    return request.param

def test_es_positivo(valor_prueba):
    assert valor_prueba > 0
```

### Pruebas Parametrizadas

Puedes parametrizar pruebas para ejecutarlas con diferentes conjuntos de datos.

```python
import pytest

def es_palindromo(texto):
    texto = texto.lower().replace(" ", "")
    return texto == texto[::-1]

@pytest.mark.parametrize("texto,esperado", [
    ("radar", True),
    ("Ana", True),
    ("A man a plan a canal Panama", True),
    ("hello", False),
    ("", True),
])
def test_es_palindromo(texto, esperado):
    assert es_palindromo(texto) == esperado
```

### Marcadores

Los marcadores permiten categorizar pruebas y ejecutar subconjuntos específicos.

```python
import pytest
import sys

@pytest.mark.slow
def test_proceso_lento():
    # Una prueba que toma mucho tiempo
    import time
    time.sleep(2)
    assert True

@pytest.mark.skipif(sys.version_info < (3, 9), reason="Requiere Python 3.9 o superior")
def test_caracteristica_nueva():
    # Prueba para una característica disponible solo en Python 3.9+
    assert True

@pytest.mark.xfail
def test_caracteristica_no_implementada():
    # Esta prueba se espera que falle
    assert False
```

Para ejecutar solo las pruebas marcadas como "slow":

```bash
pytest -m slow
```

## Mocking

El mocking es una técnica que permite reemplazar partes del sistema bajo prueba con objetos simulados.

### unittest.mock

```python
from unittest.mock import Mock, patch
import requests

def obtener_datos(url):
    response = requests.get(url)
    if response.status_code == 200:
        return response.json()
    else:
        return None

def test_obtener_datos():
    # Crear un mock para requests.get
    mock_response = Mock()
    mock_response.status_code = 200
    mock_response.json.return_value = {"data": "valor"}
    
    with patch('requests.get', return_value=mock_response):
        resultado = obtener_datos("https://api.example.com/data")
    
    assert resultado == {"data": "valor"}

def test_obtener_datos_error():
    # Simular un error en la API
    mock_response = Mock()
    mock_response.status_code = 404
    
    with patch('requests.get', return_value=mock_response):
        resultado = obtener_datos("https://api.example.com/data")
    
    assert resultado is None
```

### Mocking con pytest-mock

`pytest-mock` proporciona una fixture `mocker` que simplifica el uso de mocks en pytest.

```python
def obtener_datos(url):
    response = requests.get(url)
    if response.status_code == 200:
        return response.json()
    else:
        return None

def test_obtener_datos(mocker):
    # Crear un mock para requests.get
    mock_get = mocker.patch('requests.get')
    mock_response = mocker.Mock()
    mock_response.status_code = 200
    mock_response.json.return_value = {"data": "valor"}
    mock_get.return_value = mock_response
    
    resultado = obtener_datos("https://api.example.com/data")
    
    assert resultado == {"data": "valor"}
    mock_get.assert_called_once_with("https://api.example.com/data")
```

### Spy

Un spy permite observar las llamadas a un objeto real sin reemplazarlo completamente.

```python
def test_con_spy(mocker):
    # Crear un spy para una función real
    spy = mocker.spy(requests, 'get')
    
    # La función real se ejecuta, pero podemos verificar cómo se llamó
    try:
        requests.get("https://example.com")
    except:
        pass  # Ignorar errores de conexión en la prueba
    
    assert spy.call_count > 0
    spy.assert_called_with("https://example.com")
```

### Mock de Clases y Métodos

```python
class BaseDatos:
    def __init__(self, url):
        self.url = url
    
    def conectar(self):
        # En una aplicación real, esto se conectaría a una base de datos
        pass
    
    def obtener_usuario(self, id_usuario):
        # En una aplicación real, esto consultaría la base de datos
        pass

class ServicioUsuario:
    def __init__(self, db):
        self.db = db
    
    def obtener_nombre_usuario(self, id_usuario):
        usuario = self.db.obtener_usuario(id_usuario)
        if usuario:
            return usuario.get('nombre')
        return None

def test_servicio_usuario(mocker):
    # Crear un mock para BaseDatos
    mock_db = mocker.Mock(spec=BaseDatos)
    mock_db.obtener_usuario.return_value = {"id": 1, "nombre": "Alice"}
    
    # Usar el mock en lugar de una base de datos real
    servicio = ServicioUsuario(mock_db)
    nombre = servicio.obtener_nombre_usuario(1)
    
    assert nombre == "Alice"
    mock_db.obtener_usuario.assert_called_once_with(1)
```

## Testing de Excepciones

### Verificar que se Lance una Excepción

```python
import pytest

def dividir(a, b):
    if b == 0:
        raise ValueError("No se puede dividir por cero")
    return a / b

def test_division_por_cero():
    with pytest.raises(ValueError) as excinfo:
        dividir(10, 0)
    
    assert "No se puede dividir por cero" in str(excinfo.value)

def test_division_normal():
    assert dividir(10, 2) == 5.0
```

### Verificar Mensajes de Excepción

```python
def validar_edad(edad):
    if not isinstance(edad, int):
        raise TypeError("La edad debe ser un número entero")
    if edad < 0:
        raise ValueError("La edad no puede ser negativa")
    if edad > 120:
        raise ValueError("La edad es demasiado alta")
    return True

def test_validar_edad_tipo_incorrecto():
    with pytest.raises(TypeError, match="La edad debe ser un número entero"):
        validar_edad("treinta")

def test_validar_edad_negativa():
    with pytest.raises(ValueError, match="La edad no puede ser negativa"):
        validar_edad(-5)

def test_validar_edad_demasiado_alta():
    with pytest.raises(ValueError, match="La edad es demasiado alta"):
        validar_edad(150)

def test_validar_edad_valida():
    assert validar_edad(30) is True
```

## Testing de Aplicaciones Asíncronas

### Testing con asyncio

```python
import asyncio
import pytest

async def obtener_datos_async(id):
    # Simular una operación asíncrona
    await asyncio.sleep(0.1)
    return {"id": id, "nombre": f"Usuario {id}"}

@pytest.mark.asyncio
async def test_obtener_datos_async():
    resultado = await obtener_datos_async(1)
    assert resultado["id"] == 1
    assert resultado["nombre"] == "Usuario 1"
```

Para ejecutar pruebas asíncronas, necesitas instalar el plugin `pytest-asyncio`:

```bash
pip install pytest-asyncio
```

## Testing de Aplicaciones Web

### Testing de Flask

```python
import pytest
from flask import Flask, jsonify

# Aplicación Flask de ejemplo
app = Flask(__name__)

@app.route('/api/usuarios/<int:id>')
def obtener_usuario(id):
    usuarios = {1: {"nombre": "Alice"}, 2: {"nombre": "Bob"}}
    if id in usuarios:
        return jsonify(usuarios[id])
    return jsonify({"error": "Usuario no encontrado"}), 404

# Pruebas
@pytest.fixture
def cliente():
    app.config['TESTING'] = True
    with app.test_client() as cliente:
        yield cliente

def test_obtener_usuario_existente(cliente):
    response = cliente.get('/api/usuarios/1')
    assert response.status_code == 200
    json_data = response.get_json()
    assert json_data["nombre"] == "Alice"

def test_obtener_usuario_inexistente(cliente):
    response = cliente.get('/api/usuarios/999')
    assert response.status_code == 404
    json_data = response.get_json()
    assert "error" in json_data
```

### Testing de Django

```python
import pytest
from django.test import Client
from django.urls import reverse
from myapp.models import Usuario

@pytest.fixture
def cliente():
    return Client()

@pytest.fixture
def usuario_ejemplo():
    return Usuario.objects.create(nombre="Alice", email="alice@example.com")

def test_vista_detalle_usuario(cliente, usuario_ejemplo):
    url = reverse('detalle_usuario', args=[usuario_ejemplo.id])
    response = cliente.get(url)
    
    assert response.status_code == 200
    assert usuario_ejemplo.nombre in response.content.decode()

def test_vista_crear_usuario(cliente):
    url = reverse('crear_usuario')
    datos = {"nombre": "Charlie", "email": "charlie@example.com"}
    
    response = cliente.post(url, datos)
    
    assert response.status_code == 302  # Redirección después de crear
    assert Usuario.objects.filter(email="charlie@example.com").exists()
```

## Cobertura de Código

La cobertura de código mide qué partes de tu código se ejecutan durante las pruebas.

### Usando pytest-cov

```bash
pip install pytest-cov

# Ejecutar pruebas con informe de cobertura
pytest --cov=mimodulo

# Generar informe HTML detallado
pytest --cov=mimodulo --cov-report=html
```

### Ignorar Código en la Cobertura

```python
def funcion_principal():
    # Código normal
    resultado = procesar_datos()
    
    if __debug__:  # pragma: no cover
        # Código de depuración que no necesita cobertura
        print(f"Resultado: {resultado}")
    
    return resultado
```

## Testing de Rendimiento

### Profiling

```python
import cProfile
import pstats
import io

def test_rendimiento():
    # Crear un profiler
    pr = cProfile.Profile()
    pr.enable()
    
    # Ejecutar el código que queremos medir
    resultado = funcion_intensiva()
    
    # Desactivar el profiler
    pr.disable()
    
    # Analizar resultados
    s = io.StringIO()
    ps = pstats.Stats(pr, stream=s).sort_stats('cumulative')
    ps.print_stats(10)  # Mostrar las 10 funciones que más tiempo consumen
    
    # Verificar que el rendimiento es aceptable
    assert resultado == valor_esperado
    
    # Opcional: imprimir resultados del profiling
    print(s.getvalue())
```

### Benchmarking con pytest-benchmark

```python
def fibonacci_recursivo(n):
    if n <= 1:
        return n
    return fibonacci_recursivo(n-1) + fibonacci_recursivo(n-2)

def fibonacci_iterativo(n):
    a, b = 0, 1
    for _ in range(n):
        a, b = b, a + b
    return a

def test_benchmark_fibonacci(benchmark):
    # Medir el rendimiento de la función
    resultado = benchmark(fibonacci_iterativo, 20)
    
    # Verificar que el resultado es correcto
    assert resultado == 6765
```

## Mejores Prácticas

1. **Escribe tests primero (TDD)**: Considera escribir las pruebas antes de implementar la funcionalidad.

2. **Mantén las pruebas simples y enfocadas**: Cada prueba debe verificar una sola cosa.

3. **Usa fixtures para configuración común**: Evita duplicar código de configuración.

4. **Separa pruebas unitarias e integración**: Las pruebas unitarias deben ser rápidas y no depender de servicios externos.

5. **Usa mocks para dependencias externas**: Evita depender de servicios externos en pruebas unitarias.

6. **Ejecuta pruebas regularmente**: Integra las pruebas en tu flujo de trabajo de desarrollo.

7. **Mantén alta cobertura de código**: Apunta a una cobertura de al menos 80%.

8. **Usa nombres descriptivos para las pruebas**: El nombre debe indicar qué se está probando y qué se espera.

9. **Documenta casos de prueba complejos**: Explica qué estás probando y por qué.

10. **Evita pruebas frágiles**: Las pruebas no deben romperse con cambios menores en la implementación.

## Ejercicios Prácticos

1. Implementa un conjunto de pruebas para una API REST usando pytest y mocks.

2. Crea fixtures parametrizadas para probar una función con múltiples conjuntos de datos.

3. Escribe pruebas para una aplicación asíncrona usando pytest-asyncio.

4. Implementa pruebas de integración para una aplicación que interactúa con una base de datos.

5. Usa pytest-benchmark para comparar el rendimiento de diferentes implementaciones de un algoritmo.

## Conclusión

El testing es una parte esencial del desarrollo de software profesional. Python proporciona un rico ecosistema de herramientas para escribir y ejecutar pruebas efectivas. Dominar técnicas avanzadas como fixtures, mocking, parametrización y cobertura de código te permitirá escribir código más confiable y mantenible.

Recuerda que el objetivo del testing no es solo encontrar errores, sino también documentar el comportamiento esperado del código y facilitar futuros cambios. Un buen conjunto de pruebas te da la confianza para refactorizar y mejorar tu código sin miedo a romper la funcionalidad existente.