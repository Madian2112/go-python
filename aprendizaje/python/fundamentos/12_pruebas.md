# Pruebas en Python

## Introducción

Las pruebas son una parte fundamental del desarrollo de software que ayudan a garantizar que el código funcione correctamente, sea mantenible y cumpla con los requisitos establecidos. Python ofrece varias herramientas y frameworks para realizar diferentes tipos de pruebas. En esta sección, exploraremos los conceptos fundamentales de las pruebas en Python, incluyendo pruebas unitarias, de integración, funcionales y más.

## Pruebas Unitarias

Las pruebas unitarias verifican que componentes individuales del código (como funciones o métodos) funcionen correctamente de forma aislada.

### Módulo unittest

El módulo `unittest` es parte de la biblioteca estándar de Python y proporciona un framework para organizar y ejecutar pruebas.

```python
import unittest

# Función que queremos probar
def suma(a, b):
    return a + b

# Clase de prueba
class TestSuma(unittest.TestCase):
    def test_suma_enteros_positivos(self):
        resultado = suma(3, 5)
        self.assertEqual(resultado, 8)
    
    def test_suma_enteros_negativos(self):
        resultado = suma(-1, -1)
        self.assertEqual(resultado, -2)
    
    def test_suma_flotantes(self):
        resultado = suma(3.5, 2.5)
        self.assertEqual(resultado, 6.0)

# Ejecutar las pruebas si este archivo es el principal
if __name__ == '__main__':
    unittest.main()
```

### Métodos de Aserción

El módulo `unittest` proporciona varios métodos de aserción para verificar diferentes condiciones:

```python
import unittest

class TestAserciones(unittest.TestCase):
    def test_igualdad(self):
        self.assertEqual(5, 5)          # Verifica que dos valores sean iguales
        self.assertNotEqual(5, 6)       # Verifica que dos valores no sean iguales
    
    def test_booleanos(self):
        self.assertTrue(True)           # Verifica que el valor sea True
        self.assertFalse(False)         # Verifica que el valor sea False
    
    def test_identidad(self):
        a = [1, 2, 3]
        b = a
        c = [1, 2, 3]
        self.assertIs(a, b)             # Verifica que dos objetos sean el mismo (identidad)
        self.assertIsNot(a, c)          # Verifica que dos objetos no sean el mismo
    
    def test_tipos(self):
        self.assertIsInstance(5, int)   # Verifica que un objeto sea instancia de una clase
        self.assertNotIsInstance(5, str) # Verifica que un objeto no sea instancia de una clase
    
    def test_contenido(self):
        lista = [1, 2, 3]
        self.assertIn(1, lista)         # Verifica que un elemento esté en un contenedor
        self.assertNotIn(4, lista)      # Verifica que un elemento no esté en un contenedor
    
    def test_excepciones(self):
        # Verifica que se lance una excepción específica
        with self.assertRaises(ZeroDivisionError):
            1 / 0
    
    def test_aproximacion(self):
        # Verifica que dos valores sean aproximadamente iguales
        self.assertAlmostEqual(1.0, 1.01, places=1)  # Igual hasta 1 decimal
        self.assertAlmostEqual(1.0, 1.001, delta=0.01)  # Diferencia menor a 0.01

if __name__ == '__main__':
    unittest.main()
```

### Configuración y Limpieza

El módulo `unittest` proporciona métodos para configurar y limpiar el entorno de prueba:

```python
import unittest
import os

class TestArchivo(unittest.TestCase):
    def setUp(self):
        """Se ejecuta antes de cada método de prueba."""
        self.archivo_temporal = 'temp.txt'
        with open(self.archivo_temporal, 'w') as f:
            f.write('Contenido de prueba')
    
    def tearDown(self):
        """Se ejecuta después de cada método de prueba."""
        if os.path.exists(self.archivo_temporal):
            os.remove(self.archivo_temporal)
    
    @classmethod
    def setUpClass(cls):
        """Se ejecuta una vez antes de todas las pruebas de la clase."""
        print('Iniciando pruebas de la clase')
        cls.directorio_temporal = 'temp_dir'
        os.makedirs(cls.directorio_temporal, exist_ok=True)
    
    @classmethod
    def tearDownClass(cls):
        """Se ejecuta una vez después de todas las pruebas de la clase."""
        print('Finalizando pruebas de la clase')
        if os.path.exists(cls.directorio_temporal):
            os.rmdir(cls.directorio_temporal)
    
    def test_leer_archivo(self):
        with open(self.archivo_temporal, 'r') as f:
            contenido = f.read()
        self.assertEqual(contenido, 'Contenido de prueba')
    
    def test_escribir_archivo(self):
        with open(self.archivo_temporal, 'w') as f:
            f.write('Nuevo contenido')
        with open(self.archivo_temporal, 'r') as f:
            contenido = f.read()
        self.assertEqual(contenido, 'Nuevo contenido')

if __name__ == '__main__':
    unittest.main()
```

### Saltar Pruebas

A veces es necesario saltar ciertas pruebas bajo condiciones específicas:

```python
import unittest
import sys

class TestSaltos(unittest.TestCase):
    @unittest.skip("Demostrando skip")
    def test_siempre_saltada(self):
        self.fail("No debería ejecutarse")
    
    @unittest.skipIf(sys.version_info.minor < 10, "Requiere Python 3.10 o superior")
    def test_caracteristicas_nuevas(self):
        # Código que usa características de Python 3.10+
        self.assertEqual(1, 1)
    
    @unittest.skipUnless(sys.platform.startswith('win'), "Solo para Windows")
    def test_solo_windows(self):
        # Código específico de Windows
        self.assertEqual(1, 1)
    
    def test_esperado_fallar(self):
        # Marca una prueba que sabemos que fallará
        self.assertEqual(1, 2)
    test_esperado_fallar = unittest.expectedFailure(test_esperado_fallar)

if __name__ == '__main__':
    unittest.main()
```

### Subpruebas

Las subpruebas permiten ejecutar múltiples casos de prueba dentro de un solo método de prueba:

```python
import unittest

class TestSubpruebas(unittest.TestCase):
    def test_numeros_pares(self):
        numeros = [2, 4, 6, 8, 10, 11]
        for numero in numeros:
            with self.subTest(numero=numero):
                # Si una subprueba falla, las demás siguen ejecutándose
                self.assertEqual(numero % 2, 0)
    
    def test_diccionario(self):
        persona = {
            'nombre': 'Juan',
            'edad': 30,
            'ciudad': 'Madrid'
        }
        
        atributos = {
            'nombre': str,
            'edad': int,
            'ciudad': str,
            'telefono': str  # Este no existe en persona
        }
        
        for atributo, tipo in atributos.items():
            with self.subTest(atributo=atributo):
                self.assertIn(atributo, persona)
                if atributo in persona:
                    self.assertIsInstance(persona[atributo], tipo)

if __name__ == '__main__':
    unittest.main()
```

## pytest

`pytest` es un framework de pruebas alternativo que ofrece una sintaxis más simple y características adicionales.

### Instalación

```bash
pip install pytest
```

### Pruebas Básicas

Con `pytest`, las pruebas son funciones que comienzan con `test_`:

```python
# test_ejemplo.py

# Función que queremos probar
def suma(a, b):
    return a + b

# Pruebas
def test_suma_enteros_positivos():
    assert suma(3, 5) == 8

def test_suma_enteros_negativos():
    assert suma(-1, -1) == -2

def test_suma_flotantes():
    assert suma(3.5, 2.5) == 6.0
```

Para ejecutar las pruebas:

```bash
pytest test_ejemplo.py
```

### Fixtures

Las fixtures en `pytest` son funciones que proporcionan datos o estado para las pruebas:

```python
import pytest
import os

@pytest.fixture
def archivo_temporal():
    """Crea un archivo temporal para las pruebas."""
    archivo = 'temp.txt'
    with open(archivo, 'w') as f:
        f.write('Contenido de prueba')
    
    # Proporciona el recurso a la prueba
    yield archivo
    
    # Limpieza después de la prueba
    if os.path.exists(archivo):
        os.remove(archivo)

@pytest.fixture(scope="module")
def directorio_temporal():
    """Crea un directorio temporal para todas las pruebas del módulo."""
    directorio = 'temp_dir'
    os.makedirs(directorio, exist_ok=True)
    
    yield directorio
    
    if os.path.exists(directorio):
        os.rmdir(directorio)

def test_leer_archivo(archivo_temporal):
    with open(archivo_temporal, 'r') as f:
        contenido = f.read()
    assert contenido == 'Contenido de prueba'

def test_escribir_archivo(archivo_temporal):
    with open(archivo_temporal, 'w') as f:
        f.write('Nuevo contenido')
    with open(archivo_temporal, 'r') as f:
        contenido = f.read()
    assert contenido == 'Nuevo contenido'

def test_directorio(directorio_temporal):
    assert os.path.isdir(directorio_temporal)
```

### Parametrización

La parametrización permite ejecutar una prueba con diferentes conjuntos de datos:

```python
import pytest

def es_palindromo(texto):
    texto = texto.lower().replace(' ', '')
    return texto == texto[::-1]

@pytest.mark.parametrize("texto,esperado", [
    ("radar", True),
    ("Ana", True),
    ("Anita lava la tina", True),
    ("Python", False),
    ("", True),
    ("a", True),
])
def test_es_palindromo(texto, esperado):
    assert es_palindromo(texto) == esperado

# También se puede parametrizar con múltiples decoradores
@pytest.mark.parametrize("a", [1, 2])
@pytest.mark.parametrize("b", [3, 4])
def test_combinaciones(a, b):
    # Se ejecutará con todas las combinaciones: (1,3), (1,4), (2,3), (2,4)
    print(f"Probando con a={a}, b={b}")
    assert a < b
```

### Marcadores

Los marcadores permiten categorizar las pruebas y ejecutarlas selectivamente:

```python
import pytest
import sys

@pytest.mark.slow
def test_proceso_lento():
    # Una prueba que toma mucho tiempo
    import time
    time.sleep(1)
    assert True

@pytest.mark.skipif(sys.version_info < (3, 10), reason="Requiere Python 3.10+")
def test_caracteristicas_nuevas():
    # Código que usa características de Python 3.10+
    assert True

@pytest.mark.xfail
def test_esperado_fallar():
    assert False

@pytest.mark.parametrize("n,esperado", [(1, 2), (2, 4), (3, 6)])
@pytest.mark.matematicas
def test_doble(n, esperado):
    assert n * 2 == esperado
```

Para ejecutar pruebas con marcadores específicos:

```bash
pytest -m slow  # Ejecuta pruebas marcadas como "slow"
pytest -m "not slow"  # Ejecuta pruebas que no están marcadas como "slow"
pytest -m "matematicas and not slow"  # Combinación de marcadores
```

### Captura de Salida

`pytest` puede capturar y verificar la salida estándar y de error:

```python
def test_print(capsys):
    print("Hola, mundo!")
    salida_capturada = capsys.readouterr()
    assert "Hola, mundo!" in salida_capturada.out

def test_error(capsys):
    import sys
    sys.stderr.write("Error de prueba")
    salida_capturada = capsys.readouterr()
    assert "Error de prueba" in salida_capturada.err
```

### Temporales

`pytest` proporciona fixtures para trabajar con archivos y directorios temporales:

```python
def test_archivo_temporal(tmp_path):
    # tmp_path es un objeto Path que apunta a un directorio temporal único
    archivo = tmp_path / "test.txt"
    archivo.write_text("Contenido de prueba")
    assert archivo.read_text() == "Contenido de prueba"

def test_directorio_temporal(tmp_path_factory):
    # tmp_path_factory permite crear múltiples directorios temporales
    directorio1 = tmp_path_factory.mktemp("dir1")
    directorio2 = tmp_path_factory.mktemp("dir2")
    
    archivo1 = directorio1 / "test.txt"
    archivo1.write_text("Contenido 1")
    
    archivo2 = directorio2 / "test.txt"
    archivo2.write_text("Contenido 2")
    
    assert archivo1.read_text() != archivo2.read_text()
```

## Mocking

El mocking permite simular componentes externos o comportamientos complejos durante las pruebas.

### unittest.mock

```python
import unittest
from unittest.mock import Mock, MagicMock, patch
import requests

# Función que queremos probar
def obtener_datos(url):
    response = requests.get(url)
    if response.status_code == 200:
        return response.json()
    else:
        return None

class TestObtenerDatos(unittest.TestCase):
    def test_obtener_datos_exitoso(self):
        # Crear un mock para requests.get
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"nombre": "Juan", "edad": 30}
        
        # Reemplazar requests.get con nuestro mock
        with patch('requests.get', return_value=mock_response):
            resultado = obtener_datos("https://api.ejemplo.com/datos")
        
        # Verificar el resultado
        self.assertEqual(resultado, {"nombre": "Juan", "edad": 30})
    
    def test_obtener_datos_fallido(self):
        # Crear un mock para requests.get
        mock_response = Mock()
        mock_response.status_code = 404
        
        # Reemplazar requests.get con nuestro mock
        with patch('requests.get', return_value=mock_response):
            resultado = obtener_datos("https://api.ejemplo.com/datos")
        
        # Verificar el resultado
        self.assertIsNone(resultado)

if __name__ == '__main__':
    unittest.main()
```

### Verificación de Llamadas

```python
import unittest
from unittest.mock import Mock, call

class TestVerificacionLlamadas(unittest.TestCase):
    def test_verificar_llamadas(self):
        # Crear un mock
        mock_funcion = Mock()
        
        # Llamar al mock con diferentes argumentos
        mock_funcion(1, 2)
        mock_funcion("a", "b")
        mock_funcion(x=1, y=2)
        
        # Verificar que el mock fue llamado
        self.assertTrue(mock_funcion.called)
        
        # Verificar el número de llamadas
        self.assertEqual(mock_funcion.call_count, 3)
        
        # Verificar los argumentos de las llamadas
        mock_funcion.assert_any_call(1, 2)
        mock_funcion.assert_any_call("a", "b")
        mock_funcion.assert_any_call(x=1, y=2)
        
        # Verificar todas las llamadas de una vez
        expected_calls = [call(1, 2), call("a", "b"), call(x=1, y=2)]
        mock_funcion.assert_has_calls(expected_calls, any_order=True)

if __name__ == '__main__':
    unittest.main()
```

### Mocking de Clases y Métodos

```python
import unittest
from unittest.mock import patch, MagicMock

class BaseDatos:
    def conectar(self):
        # En realidad, esto se conectaría a una base de datos
        pass
    
    def ejecutar_consulta(self, consulta):
        # En realidad, esto ejecutaría una consulta SQL
        pass
    
    def cerrar(self):
        # En realidad, esto cerraría la conexión
        pass

class Servicio:
    def __init__(self):
        self.db = BaseDatos()
    
    def obtener_usuarios(self):
        self.db.conectar()
        resultado = self.db.ejecutar_consulta("SELECT * FROM usuarios")
        self.db.cerrar()
        return resultado

class TestServicio(unittest.TestCase):
    @patch('__main__.BaseDatos')
    def test_obtener_usuarios(self, MockBaseDatos):
        # Configurar el mock
        instancia_mock = MockBaseDatos.return_value
        instancia_mock.ejecutar_consulta.return_value = [
            {"id": 1, "nombre": "Juan"},
            {"id": 2, "nombre": "Ana"}
        ]
        
        # Crear el servicio (que usará el mock)
        servicio = Servicio()
        resultado = servicio.obtener_usuarios()
        
        # Verificar el resultado
        self.assertEqual(len(resultado), 2)
        self.assertEqual(resultado[0]["nombre"], "Juan")
        
        # Verificar que los métodos fueron llamados correctamente
        instancia_mock.conectar.assert_called_once()
        instancia_mock.ejecutar_consulta.assert_called_once_with("SELECT * FROM usuarios")
        instancia_mock.cerrar.assert_called_once()

if __name__ == '__main__':
    unittest.main()
```

### Mocking con pytest

`pytest` proporciona soporte para mocking a través del paquete `pytest-mock`:

```python
import requests

def obtener_datos(url):
    response = requests.get(url)
    if response.status_code == 200:
        return response.json()
    else:
        return None

def test_obtener_datos_exitoso(mocker):
    # Crear un mock para requests.get
    mock_get = mocker.patch('requests.get')
    mock_response = mocker.Mock()
    mock_response.status_code = 200
    mock_response.json.return_value = {"nombre": "Juan", "edad": 30}
    mock_get.return_value = mock_response
    
    # Llamar a la función
    resultado = obtener_datos("https://api.ejemplo.com/datos")
    
    # Verificar el resultado
    assert resultado == {"nombre": "Juan", "edad": 30}
    mock_get.assert_called_once_with("https://api.ejemplo.com/datos")

def test_obtener_datos_fallido(mocker):
    # Crear un mock para requests.get
    mock_get = mocker.patch('requests.get')
    mock_response = mocker.Mock()
    mock_response.status_code = 404
    mock_get.return_value = mock_response
    
    # Llamar a la función
    resultado = obtener_datos("https://api.ejemplo.com/datos")
    
    # Verificar el resultado
    assert resultado is None
```

## Pruebas de Integración

Las pruebas de integración verifican que diferentes componentes del sistema funcionen correctamente juntos.

```python
import unittest
import tempfile
import os
import json

# Componentes que queremos probar juntos
class Almacenamiento:
    def __init__(self, ruta_archivo):
        self.ruta_archivo = ruta_archivo
    
    def guardar_datos(self, datos):
        with open(self.ruta_archivo, 'w') as f:
            json.dump(datos, f)
    
    def cargar_datos(self):
        if not os.path.exists(self.ruta_archivo):
            return {}
        with open(self.ruta_archivo, 'r') as f:
            return json.load(f)

class GestorUsuarios:
    def __init__(self, almacenamiento):
        self.almacenamiento = almacenamiento
    
    def agregar_usuario(self, id, nombre, email):
        datos = self.almacenamiento.cargar_datos()
        if 'usuarios' not in datos:
            datos['usuarios'] = {}
        
        datos['usuarios'][id] = {
            'nombre': nombre,
            'email': email
        }
        
        self.almacenamiento.guardar_datos(datos)
    
    def obtener_usuario(self, id):
        datos = self.almacenamiento.cargar_datos()
        if 'usuarios' not in datos or id not in datos['usuarios']:
            return None
        return datos['usuarios'][id]
    
    def eliminar_usuario(self, id):
        datos = self.almacenamiento.cargar_datos()
        if 'usuarios' in datos and id in datos['usuarios']:
            del datos['usuarios'][id]
            self.almacenamiento.guardar_datos(datos)
            return True
        return False

# Prueba de integración
class TestIntegracionUsuarios(unittest.TestCase):
    def setUp(self):
        # Crear un archivo temporal para las pruebas
        self.archivo_temp = tempfile.NamedTemporaryFile(delete=False).name
        self.almacenamiento = Almacenamiento(self.archivo_temp)
        self.gestor = GestorUsuarios(self.almacenamiento)
    
    def tearDown(self):
        # Eliminar el archivo temporal
        if os.path.exists(self.archivo_temp):
            os.remove(self.archivo_temp)
    
    def test_flujo_completo(self):
        # Agregar usuarios
        self.gestor.agregar_usuario("1", "Juan", "juan@example.com")
        self.gestor.agregar_usuario("2", "Ana", "ana@example.com")
        
        # Verificar que se pueden obtener
        usuario1 = self.gestor.obtener_usuario("1")
        self.assertEqual(usuario1["nombre"], "Juan")
        self.assertEqual(usuario1["email"], "juan@example.com")
        
        usuario2 = self.gestor.obtener_usuario("2")
        self.assertEqual(usuario2["nombre"], "Ana")
        
        # Eliminar un usuario
        self.assertTrue(self.gestor.eliminar_usuario("1"))
        
        # Verificar que ya no existe
        self.assertIsNone(self.gestor.obtener_usuario("1"))
        
        # Verificar que el otro usuario sigue existiendo
        usuario2 = self.gestor.obtener_usuario("2")
        self.assertEqual(usuario2["nombre"], "Ana")
        
        # Verificar que los datos persisten en el archivo
        with open(self.archivo_temp, 'r') as f:
            datos = json.load(f)
        
        self.assertIn("usuarios", datos)
        self.assertIn("2", datos["usuarios"])
        self.assertNotIn("1", datos["usuarios"])

if __name__ == '__main__':
    unittest.main()
```

## Pruebas Funcionales

Las pruebas funcionales verifican que el sistema cumpla con los requisitos funcionales desde la perspectiva del usuario.

### Selenium para Pruebas Web

```python
import unittest
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

class TestPaginaWeb(unittest.TestCase):
    def setUp(self):
        # Inicializar el navegador
        self.driver = webdriver.Chrome()  # Requiere ChromeDriver
        self.driver.implicitly_wait(10)  # Espera implícita
    
    def tearDown(self):
        # Cerrar el navegador
        self.driver.quit()
    
    def test_busqueda_google(self):
        # Abrir Google
        self.driver.get("https://www.google.com")
        
        # Aceptar cookies si aparece el diálogo
        try:
            boton_aceptar = WebDriverWait(self.driver, 5).until(
                EC.element_to_be_clickable((By.ID, "L2AGLb"))
            )
            boton_aceptar.click()
        except:
            pass  # No apareció el diálogo de cookies
        
        # Encontrar el campo de búsqueda y escribir
        campo_busqueda = self.driver.find_element(By.NAME, "q")
        campo_busqueda.send_keys("Python testing")
        campo_busqueda.send_keys(Keys.RETURN)
        
        # Esperar a que aparezcan los resultados
        WebDriverWait(self.driver, 10).until(
            EC.presence_of_element_located((By.ID, "search"))
        )
        
        # Verificar que hay resultados
        resultados = self.driver.find_elements(By.CSS_SELECTOR, "#search .g")
        self.assertGreater(len(resultados), 0)
        
        # Verificar que "Python" aparece en el título de la página
        self.assertIn("Python", self.driver.title)

if __name__ == "__main__":
    unittest.main()
```

## Cobertura de Código

La cobertura de código mide qué partes del código se ejecutan durante las pruebas.

### coverage.py

```bash
pip install coverage
```

#### Uso Básico

```bash
# Ejecutar pruebas con cobertura
coverage run -m unittest discover

# Generar informe
coverage report

# Generar informe HTML
coverage html
```

#### Integración con pytest

```bash
pip install pytest-cov
```

```bash
# Ejecutar pytest con cobertura
pytest --cov=mi_paquete tests/

# Generar informe HTML
pytest --cov=mi_paquete --cov-report=html tests/
```

## Pruebas de Rendimiento

Las pruebas de rendimiento verifican que el sistema cumpla con los requisitos de rendimiento.

### timeit

```python
import timeit

# Función que queremos medir
def ordenar_lista(n):
    return sorted([i for i in range(n, 0, -1)])

# Medir el tiempo de ejecución
tiempo = timeit.timeit("ordenar_lista(1000)", globals=globals(), number=100)
print(f"Tiempo promedio para ordenar_lista(1000): {tiempo/100:.6f} segundos")

# Comparar diferentes implementaciones
def ordenar_lista_alternativa(n):
    lista = list(range(n, 0, -1))
    lista.sort()
    return lista

tiempo1 = timeit.timeit("ordenar_lista(1000)", globals=globals(), number=100)
tiempo2 = timeit.timeit("ordenar_lista_alternativa(1000)", globals=globals(), number=100)

print(f"Tiempo ordenar_lista: {tiempo1/100:.6f} segundos")
print(f"Tiempo ordenar_lista_alternativa: {tiempo2/100:.6f} segundos")
print(f"Diferencia: {abs(tiempo1-tiempo2)/100:.6f} segundos")
```

### cProfile

```python
import cProfile
import pstats
import io

def funcion_compleja():
    resultado = 0
    for i in range(1000000):
        resultado += i
    return resultado

# Perfilar la función
pr = cProfile.Profile()
pr.enable()

funcion_compleja()

pr.disable()

# Imprimir estadísticas
s = io.StringIO()
ps = pstats.Stats(pr, stream=s).sort_stats('cumulative')
ps.print_stats(10)  # Mostrar las 10 funciones que más tiempo consumen
print(s.getvalue())
```

## Pruebas de Carga

Las pruebas de carga verifican cómo se comporta el sistema bajo una carga específica.

### locust

```bash
pip install locust
```

```python
# locustfile.py
from locust import HttpUser, task, between

class WebsiteUser(HttpUser):
    wait_time = between(1, 5)  # Tiempo de espera entre tareas
    
    @task
    def index_page(self):
        self.client.get("/")
    
    @task(3)  # Esta tarea se ejecuta 3 veces más que las otras
    def view_products(self):
        self.client.get("/products")
    
    @task
    def view_about(self):
        self.client.get("/about")
```

Para ejecutar locust:

```bash
locust -f locustfile.py --host=http://example.com
```

## Pruebas de Seguridad

Las pruebas de seguridad verifican que el sistema sea seguro contra amenazas conocidas.

### bandit

```bash
pip install bandit
```

```bash
# Analizar un archivo
bandit mi_script.py

# Analizar un directorio recursivamente
bandit -r mi_proyecto/

# Generar informe HTML
bandit -r mi_proyecto/ -f html -o informe_seguridad.html
```

## Ejemplo Práctico: Pruebas para una API REST

Vamos a crear pruebas para una API REST simple utilizando Flask y pytest.

### API de Tareas

```python
# app.py
from flask import Flask, request, jsonify

app = Flask(__name__)

# Base de datos en memoria
tareas = {}
id_contador = 1

@app.route('/tareas', methods=['GET'])
def obtener_tareas():
    return jsonify(list(tareas.values()))

@app.route('/tareas/<int:id>', methods=['GET'])
def obtener_tarea(id):
    tarea = tareas.get(id)
    if tarea:
        return jsonify(tarea)
    return jsonify({"error": "Tarea no encontrada"}), 404

@app.route('/tareas', methods=['POST'])
def crear_tarea():
    global id_contador
    datos = request.json
    
    if not datos or 'titulo' not in datos:
        return jsonify({"error": "Se requiere un título"}), 400
    
    tarea = {
        "id": id_contador,
        "titulo": datos["titulo"],
        "descripcion": datos.get("descripcion", ""),
        "completada": False
    }
    
    tareas[id_contador] = tarea
    id_contador += 1
    
    return jsonify(tarea), 201

@app.route('/tareas/<int:id>', methods=['PUT'])
def actualizar_tarea(id):
    if id not in tareas:
        return jsonify({"error": "Tarea no encontrada"}), 404
    
    datos = request.json
    if not datos:
        return jsonify({"error": "No hay datos para actualizar"}), 400
    
    tarea = tareas[id]
    
    if "titulo" in datos:
        tarea["titulo"] = datos["titulo"]
    
    if "descripcion" in datos:
        tarea["descripcion"] = datos["descripcion"]
    
    if "completada" in datos:
        tarea["completada"] = bool(datos["completada"])
    
    return jsonify(tarea)

@app.route('/tareas/<int:id>', methods=['DELETE'])
def eliminar_tarea(id):
    if id not in tareas:
        return jsonify({"error": "Tarea no encontrada"}), 404
    
    tarea = tareas.pop(id)
    return jsonify(tarea)

if __name__ == '__main__':
    app.run(debug=True)
```

### Pruebas para la API

```python
# test_app.py
import pytest
import json
from app import app, tareas, id_contador

@pytest.fixture
def cliente():
    app.config['TESTING'] = True
    with app.test_client() as cliente:
        yield cliente

@pytest.fixture
def limpiar_datos():
    # Limpiar datos antes de cada prueba
    tareas.clear()
    global id_contador
    id_contador = 1
    yield

def test_obtener_tareas_vacio(cliente, limpiar_datos):
    respuesta = cliente.get('/tareas')
    assert respuesta.status_code == 200
    assert json.loads(respuesta.data) == []

def test_crear_tarea(cliente, limpiar_datos):
    # Crear una tarea
    datos = {"titulo": "Comprar leche", "descripcion": "2 litros"}
    respuesta = cliente.post('/tareas', json=datos)
    
    # Verificar respuesta
    assert respuesta.status_code == 201
    tarea = json.loads(respuesta.data)
    assert tarea["id"] == 1
    assert tarea["titulo"] == "Comprar leche"
    assert tarea["descripcion"] == "2 litros"
    assert tarea["completada"] == False
    
    # Verificar que la tarea se guardó
    respuesta = cliente.get('/tareas')
    tareas_guardadas = json.loads(respuesta.data)
    assert len(tareas_guardadas) == 1
    assert tareas_guardadas[0]["id"] == 1

def test_crear_tarea_sin_titulo(cliente, limpiar_datos):
    # Intentar crear una tarea sin título
    datos = {"descripcion": "Sin título"}
    respuesta = cliente.post('/tareas', json=datos)
    
    # Verificar error
    assert respuesta.status_code == 400
    assert "error" in json.loads(respuesta.data)

def test_obtener_tarea_existente(cliente, limpiar_datos):
    # Crear una tarea
    datos = {"titulo": "Comprar leche"}
    cliente.post('/tareas', json=datos)
    
    # Obtener la tarea
    respuesta = cliente.get('/tareas/1')
    
    # Verificar respuesta
    assert respuesta.status_code == 200
    tarea = json.loads(respuesta.data)
    assert tarea["id"] == 1
    assert tarea["titulo"] == "Comprar leche"

def test_obtener_tarea_inexistente(cliente, limpiar_datos):
    # Intentar obtener una tarea que no existe
    respuesta = cliente.get('/tareas/999')
    
    # Verificar error
    assert respuesta.status_code == 404
    assert "error" in json.loads(respuesta.data)

def test_actualizar_tarea(cliente, limpiar_datos):
    # Crear una tarea
    datos = {"titulo": "Comprar leche"}
    cliente.post('/tareas', json=datos)
    
    # Actualizar la tarea
    datos_actualizados = {"titulo": "Comprar leche desnatada", "completada": True}
    respuesta = cliente.put('/tareas/1', json=datos_actualizados)
    
    # Verificar respuesta
    assert respuesta.status_code == 200
    tarea = json.loads(respuesta.data)
    assert tarea["titulo"] == "Comprar leche desnatada"
    assert tarea["completada"] == True
    
    # Verificar que la tarea se actualizó
    respuesta = cliente.get('/tareas/1')
    tarea_guardada = json.loads(respuesta.data)
    assert tarea_guardada["titulo"] == "Comprar leche desnatada"
    assert tarea_guardada["completada"] == True

def test_eliminar_tarea(cliente, limpiar_datos):
    # Crear una tarea
    datos = {"titulo": "Comprar leche"}
    cliente.post('/tareas', json=datos)
    
    # Eliminar la tarea
    respuesta = cliente.delete('/tareas/1')
    
    # Verificar respuesta
    assert respuesta.status_code == 200
    
    # Verificar que la tarea se eliminó
    respuesta = cliente.get('/tareas/1')
    assert respuesta.status_code == 404
    
    # Verificar que no hay tareas
    respuesta = cliente.get('/tareas')
    tareas_guardadas = json.loads(respuesta.data)
    assert len(tareas_guardadas) == 0
```

## Buenas Prácticas

1. **Pruebas Independientes**: Cada prueba debe ser independiente de las demás. No debe depender del estado dejado por otras pruebas.

2. **Pruebas Claras**: Las pruebas deben ser fáciles de entender. Usa nombres descriptivos y comentarios cuando sea necesario.

3. **Pruebas Rápidas**: Las pruebas deben ejecutarse rápidamente. Si una prueba es lenta, considera marcarla como tal y ejecutarla solo cuando sea necesario.

4. **Pruebas Completas**: Las pruebas deben cubrir tanto los casos normales como los casos límite y de error.

5. **Pruebas Deterministas**: Las pruebas deben producir el mismo resultado cada vez que se ejecuten. Evita dependencias de factores externos como la hora actual o valores aleatorios.

6. **Pruebas Automatizadas**: Las pruebas deben poder ejecutarse automáticamente, sin intervención manual.

7. **Pruebas Mantenibles**: Las pruebas deben ser fáciles de mantener. Usa fixtures y helpers para evitar la duplicación de código.

8. **Pruebas como Documentación**: Las pruebas pueden servir como documentación del comportamiento esperado del código.

9. **Pruebas Primero (TDD)**: Considera escribir las pruebas antes de escribir el código (Test-Driven Development).

10. **Integración Continua**: Ejecuta las pruebas automáticamente en cada cambio del código.

## Recursos Adicionales

- [Documentación oficial de unittest](https://docs.python.org/es/3/library/unittest.html)
- [Documentación oficial de pytest](https://docs.pytest.org/)
- [Documentación oficial de unittest.mock](https://docs.python.org/es/3/library/unittest.mock.html)
- [Documentación oficial de coverage.py](https://coverage.readthedocs.io/)
- [Python Testing with pytest (libro)](https://pragprog.com/titles/bopytest/python-testing-with-pytest/)
- [Test-Driven Development with Python (libro)](https://www.obeythetestinggoat.com/)
- [Effective Python Testing with pytest (artículo)](https://realpython.com/pytest-python-testing/)
- [Python Mocking: A Guide to Better Unit Tests (artículo)](https://realpython.com/python-mock-library/)

---

Con esto concluimos nuestra exploración de las pruebas en Python. Las pruebas son una parte esencial del desarrollo de software profesional y Python proporciona herramientas excelentes para implementarlas. Recuerda que las pruebas no solo ayudan a encontrar errores, sino que también mejoran la calidad del código, facilitan el mantenimiento y proporcionan confianza al realizar cambios.