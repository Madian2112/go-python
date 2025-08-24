# Programación Concurrente en Python

## Introducción

La programación concurrente permite que múltiples tareas progresen simultáneamente, mejorando el rendimiento y la capacidad de respuesta de las aplicaciones. Python ofrece varias herramientas para implementar concurrencia, cada una con sus propias características y casos de uso. En esta sección, exploraremos los diferentes enfoques de concurrencia disponibles en Python.

## Hilos (Threading)

Los hilos permiten la ejecución concurrente dentro de un mismo proceso, compartiendo el mismo espacio de memoria.

### Creación de Hilos

```python
import threading
import time

def tarea(nombre, segundos):
    """Función que simula una tarea que toma tiempo."""
    print(f"Iniciando tarea {nombre}")
    time.sleep(segundos)  # Simula trabajo
    print(f"Completada tarea {nombre} después de {segundos} segundos")

# Crear hilos
hilo1 = threading.Thread(target=tarea, args=("A", 3))
hilo2 = threading.Thread(target=tarea, args=("B", 2))

# Iniciar hilos
tiempo_inicio = time.time()
hilo1.start()
hilo2.start()

# Esperar a que los hilos terminen
hilo1.join()
hilo2.join()

tiempo_total = time.time() - tiempo_inicio
print(f"Tiempo total: {tiempo_total:.2f} segundos")
# Salida aproximada:
# Iniciando tarea A
# Iniciando tarea B
# Completada tarea B después de 2 segundos
# Completada tarea A después de 3 segundos
# Tiempo total: 3.01 segundos
```

### Hilos como Clases

```python
import threading
import time

class MiHilo(threading.Thread):
    def __init__(self, nombre, segundos):
        super().__init__()
        self.nombre = nombre
        self.segundos = segundos
    
    def run(self):
        """Método que se ejecuta cuando se inicia el hilo."""
        print(f"Iniciando tarea {self.nombre}")
        time.sleep(self.segundos)  # Simula trabajo
        print(f"Completada tarea {self.nombre} después de {self.segundos} segundos")

# Crear hilos
hilo1 = MiHilo("A", 3)
hilo2 = MiHilo("B", 2)

# Iniciar hilos
tiempo_inicio = time.time()
hilo1.start()
hilo2.start()

# Esperar a que los hilos terminen
hilo1.join()
hilo2.join()

tiempo_total = time.time() - tiempo_inicio
print(f"Tiempo total: {tiempo_total:.2f} segundos")
```

### Sincronización de Hilos

Cuando múltiples hilos acceden a recursos compartidos, es necesario sincronizarlos para evitar condiciones de carrera.

#### Locks (Cerrojos)

```python
import threading
import time

# Recurso compartido
contador = 0

# Lock para proteger el recurso compartido
lock = threading.Lock()

def incrementar(n):
    """Incrementa el contador n veces."""
    global contador
    for _ in range(n):
        # Adquirir el lock
        lock.acquire()
        try:
            # Sección crítica
            valor_actual = contador
            time.sleep(0.0001)  # Simula algún procesamiento
            contador = valor_actual + 1
        finally:
            # Liberar el lock
            lock.release()

# Alternativa usando with (recomendado)
def incrementar_with(n):
    """Incrementa el contador n veces usando with."""
    global contador
    for _ in range(n):
        with lock:  # Adquiere y libera automáticamente
            valor_actual = contador
            time.sleep(0.0001)  # Simula algún procesamiento
            contador = valor_actual + 1

# Crear hilos
hilos = []
for _ in range(5):
    hilo = threading.Thread(target=incrementar_with, args=(1000,))
    hilos.append(hilo)

# Iniciar hilos
for hilo in hilos:
    hilo.start()

# Esperar a que los hilos terminen
for hilo in hilos:
    hilo.join()

print(f"Valor final del contador: {contador}")
# Salida: Valor final del contador: 5000
```

#### RLock (Cerrojos Reentrantes)

Permiten que un mismo hilo adquiera el lock múltiples veces sin bloquearse a sí mismo.

```python
import threading

rlock = threading.RLock()

def funcion_externa():
    with rlock:  # Primera adquisición
        print("Función externa adquirió el lock")
        funcion_interna()

def funcion_interna():
    with rlock:  # Segunda adquisición (por el mismo hilo)
        print("Función interna también adquirió el lock")

# Crear y ejecutar un hilo
hilo = threading.Thread(target=funcion_externa)
hilo.start()
hilo.join()
```

#### Semáforos

Permiten que un número limitado de hilos accedan a un recurso.

```python
import threading
import time
import random

# Semáforo que permite hasta 3 hilos simultáneos
semaforo = threading.Semaphore(3)

def trabajador(id):
    print(f"Trabajador {id} esperando para acceder al recurso")
    with semaforo:
        print(f"Trabajador {id} accedió al recurso")
        time.sleep(random.uniform(0.5, 2))  # Simula trabajo
        print(f"Trabajador {id} liberó el recurso")

# Crear 10 hilos
hilos = []
for i in range(10):
    hilo = threading.Thread(target=trabajador, args=(i,))
    hilos.append(hilo)
    hilo.start()

# Esperar a que todos terminen
for hilo in hilos:
    hilo.join()
```

#### Eventos

Permiten que un hilo señalice a otros hilos que ha ocurrido algo.

```python
import threading
import time

# Crear un evento
evento = threading.Event()

def esperar_evento():
    print("Esperando evento...")
    evento.wait()  # Bloquea hasta que el evento sea establecido
    print("¡Evento recibido! Continuando...")

def establecer_evento():
    print("Trabajando...")
    time.sleep(3)  # Simula trabajo
    print("Estableciendo evento")
    evento.set()  # Establece el evento, desbloquea los hilos que esperan

# Crear hilos
hilo_espera1 = threading.Thread(target=esperar_evento)
hilo_espera2 = threading.Thread(target=esperar_evento)
hilo_establecer = threading.Thread(target=establecer_evento)

# Iniciar hilos
hilo_espera1.start()
hilo_espera2.start()
time.sleep(1)  # Asegura que los hilos de espera estén bloqueados
hilo_establecer.start()

# Esperar a que todos terminen
hilo_espera1.join()
hilo_espera2.join()
hilo_establecer.join()
```

#### Condiciones

Permiten que los hilos esperen hasta que una condición sea verdadera.

```python
import threading
import time
import random

# Recurso compartido
items = []
max_items = 10

# Condición
condicion = threading.Condition()

def productor():
    global items
    for i in range(20):
        with condicion:
            while len(items) >= max_items:
                print("Buffer lleno, productor esperando...")
                condicion.wait()  # Espera a que haya espacio
            
            item = random.randint(1, 100)
            items.append(item)
            print(f"Producido: {item}, buffer: {len(items)}")
            
            condicion.notify()  # Notifica a un consumidor
        
        time.sleep(random.uniform(0.1, 0.5))  # Simula trabajo

def consumidor():
    global items
    for _ in range(10):
        with condicion:
            while not items:  # Mientras no haya items
                print("Buffer vacío, consumidor esperando...")
                condicion.wait()  # Espera a que haya items
            
            item = items.pop(0)
            print(f"Consumido: {item}, buffer: {len(items)}")
            
            condicion.notify()  # Notifica a un productor
        
        time.sleep(random.uniform(0.2, 0.7))  # Simula trabajo

# Crear hilos
hilo_productor = threading.Thread(target=productor)
hilos_consumidores = [threading.Thread(target=consumidor) for _ in range(2)]

# Iniciar hilos
hilo_productor.start()
for hilo in hilos_consumidores:
    hilo.start()

# Esperar a que terminen
hilo_productor.join()
for hilo in hilos_consumidores:
    hilo.join()
```

#### Barreras

Permiten que múltiples hilos esperen en un punto hasta que todos lleguen.

```python
import threading
import time
import random

# Barrera para 4 hilos
barrera = threading.Barrier(4)

def trabajador(id):
    print(f"Trabajador {id} iniciando fase 1")
    time.sleep(random.uniform(0.5, 2))  # Simula trabajo
    print(f"Trabajador {id} completó fase 1, esperando a los demás")
    
    # Esperar a que todos completen la fase 1
    barrera.wait()
    
    print(f"Trabajador {id} iniciando fase 2")
    time.sleep(random.uniform(0.5, 2))  # Simula trabajo
    print(f"Trabajador {id} completó fase 2")

# Crear hilos
hilos = []
for i in range(4):
    hilo = threading.Thread(target=trabajador, args=(i,))
    hilos.append(hilo)
    hilo.start()

# Esperar a que todos terminen
for hilo in hilos:
    hilo.join()
```

### Limitaciones del Threading en Python

Debido al Global Interpreter Lock (GIL) de Python, los hilos no pueden ejecutar código Python en paralelo en múltiples núcleos. El GIL asegura que solo un hilo ejecute código Python a la vez, lo que limita el rendimiento en tareas intensivas de CPU. Sin embargo, los hilos siguen siendo útiles para:

- Tareas de E/S (entrada/salida) como operaciones de red o disco
- Mantener la capacidad de respuesta de la interfaz de usuario
- Simplificar el diseño de programas concurrentes

Para tareas intensivas de CPU que requieren paralelismo real, se recomienda usar el módulo `multiprocessing`.

## Procesos (Multiprocessing)

El módulo `multiprocessing` permite crear procesos independientes, cada uno con su propio intérprete de Python y espacio de memoria, evitando las limitaciones del GIL.

### Creación de Procesos

```python
import multiprocessing as mp
import time

def tarea_intensiva(nombre, segundos):
    """Función que simula una tarea intensiva de CPU."""
    print(f"Iniciando tarea {nombre}")
    # Simula trabajo intensivo de CPU
    inicio = time.time()
    contador = 0
    while time.time() - inicio < segundos:
        contador += 1
    print(f"Completada tarea {nombre} después de {segundos} segundos")
    return contador

if __name__ == "__main__":  # Necesario en Windows
    # Crear procesos
    proceso1 = mp.Process(target=tarea_intensiva, args=("A", 3))
    proceso2 = mp.Process(target=tarea_intensiva, args=("B", 3))
    
    # Iniciar procesos
    tiempo_inicio = time.time()
    proceso1.start()
    proceso2.start()
    
    # Esperar a que los procesos terminen
    proceso1.join()
    proceso2.join()
    
    tiempo_total = time.time() - tiempo_inicio
    print(f"Tiempo total: {tiempo_total:.2f} segundos")
    # En un sistema con múltiples núcleos, el tiempo total será aproximadamente 3 segundos,
    # no 6 segundos como sería con hilos para tareas intensivas de CPU
```

### Pool de Procesos

Permite distribuir tareas entre un número fijo de procesos trabajadores.

```python
import multiprocessing as mp
import time

def procesar_item(item):
    """Procesa un solo item."""
    print(f"Procesando {item}")
    time.sleep(1)  # Simula trabajo
    return item * item

if __name__ == "__main__":
    # Crear un pool con 4 procesos trabajadores
    with mp.Pool(processes=4) as pool:
        # Método map: aplica la función a cada elemento de la lista
        items = [1, 2, 3, 4, 5, 6, 7, 8]
        resultados = pool.map(procesar_item, items)
        print(f"Resultados: {resultados}")
        
        # Método apply_async: envía tareas de forma asíncrona
        resultados_async = []
        for item in items:
            resultado = pool.apply_async(procesar_item, (item,))
            resultados_async.append(resultado)
        
        # Obtener los resultados
        resultados_finales = [res.get() for res in resultados_async]
        print(f"Resultados asíncronos: {resultados_finales}")
```

### Comunicación entre Procesos

#### Colas (Queue)

```python
import multiprocessing as mp
import random
import time

def productor(cola):
    """Produce items y los coloca en la cola."""
    for i in range(10):
        item = random.randint(1, 100)
        cola.put(item)
        print(f"Producido: {item}")
        time.sleep(random.uniform(0.1, 0.5))
    
    # Señal de finalización
    cola.put(None)
    print("Productor terminado")

def consumidor(cola):
    """Consume items de la cola."""
    while True:
        item = cola.get()
        if item is None:  # Señal de finalización
            break
        print(f"Consumido: {item}")
        time.sleep(random.uniform(0.2, 0.7))
    
    print("Consumidor terminado")

if __name__ == "__main__":
    # Crear una cola
    cola = mp.Queue()
    
    # Crear procesos
    proc_productor = mp.Process(target=productor, args=(cola,))
    proc_consumidor = mp.Process(target=consumidor, args=(cola,))
    
    # Iniciar procesos
    proc_productor.start()
    proc_consumidor.start()
    
    # Esperar a que terminen
    proc_productor.join()
    proc_consumidor.join()
```

#### Pipes

```python
import multiprocessing as mp

def emisor(conexion, mensajes):
    """Envía mensajes a través de la conexión."""
    for msg in mensajes:
        conexion.send(msg)
        print(f"Enviado: {msg}")
    
    conexion.close()

def receptor(conexion):
    """Recibe mensajes de la conexión."""
    while True:
        try:
            msg = conexion.recv()
            print(f"Recibido: {msg}")
        except EOFError:
            break  # La conexión se cerró
    
    print("Receptor terminado")

if __name__ == "__main__":
    # Crear un pipe
    conexion_receptor, conexion_emisor = mp.Pipe()
    
    # Mensajes a enviar
    mensajes = ["Hola", 42, {"clave": "valor"}, [1, 2, 3]]
    
    # Crear procesos
    proc_emisor = mp.Process(target=emisor, args=(conexion_emisor, mensajes))
    proc_receptor = mp.Process(target=receptor, args=(conexion_receptor,))
    
    # Iniciar procesos
    proc_emisor.start()
    proc_receptor.start()
    
    # Esperar a que terminen
    proc_emisor.join()
    proc_receptor.join()
```

#### Memoria Compartida

```python
import multiprocessing as mp
import time

def incrementar(valor, lock):
    """Incrementa el valor compartido."""
    for _ in range(100):
        with lock:
            valor.value += 1
        time.sleep(0.01)

def imprimir_valor(valor):
    """Imprime el valor compartido periódicamente."""
    for _ in range(10):
        print(f"Valor actual: {valor.value}")
        time.sleep(0.5)

if __name__ == "__main__":
    # Crear un valor compartido
    valor_compartido = mp.Value('i', 0)  # 'i' indica entero
    
    # Crear un lock para sincronización
    lock = mp.Lock()
    
    # Crear procesos
    procesos_incrementar = [
        mp.Process(target=incrementar, args=(valor_compartido, lock))
        for _ in range(4)
    ]
    
    proceso_imprimir = mp.Process(target=imprimir_valor, args=(valor_compartido,))
    
    # Iniciar procesos
    for p in procesos_incrementar:
        p.start()
    
    proceso_imprimir.start()
    
    # Esperar a que terminen
    for p in procesos_incrementar:
        p.join()
    
    proceso_imprimir.join()
    
    print(f"Valor final: {valor_compartido.value}")
```

#### Array Compartido

```python
import multiprocessing as mp
import random

def llenar_seccion(array, inicio, fin):
    """Llena una sección del array con números aleatorios."""
    for i in range(inicio, fin):
        array[i] = random.randint(1, 100)

def calcular_suma(array, resultado):
    """Calcula la suma de todos los elementos del array."""
    resultado.value = sum(array)

if __name__ == "__main__":
    # Tamaño del array
    n = 10000000
    
    # Crear un array compartido
    array_compartido = mp.Array('i', n)  # 'i' indica entero
    
    # Dividir el trabajo entre 4 procesos
    tamaño_seccion = n // 4
    procesos = []
    
    for i in range(4):
        inicio = i * tamaño_seccion
        fin = inicio + tamaño_seccion if i < 3 else n
        p = mp.Process(target=llenar_seccion, args=(array_compartido, inicio, fin))
        procesos.append(p)
        p.start()
    
    # Esperar a que terminen
    for p in procesos:
        p.join()
    
    # Calcular la suma en otro proceso
    resultado = mp.Value('i', 0)
    p_suma = mp.Process(target=calcular_suma, args=(array_compartido, resultado))
    p_suma.start()
    p_suma.join()
    
    print(f"Suma de los elementos: {resultado.value}")
```

### Administrador (Manager)

Permite compartir objetos más complejos entre procesos.

```python
import multiprocessing as mp
import time

def trabajador(diccionario, lista, lock):
    """Modifica estructuras de datos compartidas."""
    for i in range(5):
        with lock:
            # Modificar el diccionario
            clave = f"proceso-{mp.current_process().name}-{i}"
            diccionario[clave] = i
            
            # Modificar la lista
            lista.append(i)
        
        time.sleep(0.1)

if __name__ == "__main__":
    # Crear un manager
    with mp.Manager() as manager:
        # Crear estructuras de datos compartidas
        diccionario_compartido = manager.dict()
        lista_compartida = manager.list()
        lock = manager.Lock()
        
        # Crear procesos
        procesos = [
            mp.Process(target=trabajador, args=(diccionario_compartido, lista_compartida, lock))
            for _ in range(3)
        ]
        
        # Iniciar procesos
        for p in procesos:
            p.start()
        
        # Esperar a que terminen
        for p in procesos:
            p.join()
        
        # Mostrar resultados
        print(f"Diccionario compartido: {dict(diccionario_compartido)}")
        print(f"Lista compartida: {list(lista_compartida)}")
```

## Programación Asíncrona (asyncio)

El módulo `asyncio` proporciona una infraestructura para escribir código concurrente utilizando la sintaxis `async`/`await`. Es especialmente útil para aplicaciones con muchas operaciones de E/S.

### Conceptos Básicos

```python
import asyncio

async def saludar(nombre, delay):
    """Función asíncrona que saluda después de un delay."""
    await asyncio.sleep(delay)  # Pausa sin bloquear el bucle de eventos
    print(f"Hola, {nombre}!")
    return f"{nombre} saludado después de {delay} segundos"

async def main():
    # Ejecutar tareas concurrentemente
    resultados = await asyncio.gather(
        saludar("Alice", 1),
        saludar("Bob", 2),
        saludar("Charlie", 3)
    )
    
    print(f"Resultados: {resultados}")

# Python 3.7+
asyncio.run(main())
```

### Tareas (Tasks)

```python
import asyncio
import time

async def tarea_larga():
    print("Iniciando tarea larga...")
    await asyncio.sleep(5)  # Simula una operación larga
    print("Tarea larga completada")
    return "Resultado de tarea larga"

async def tarea_corta():
    print("Iniciando tarea corta...")
    await asyncio.sleep(1)  # Simula una operación corta
    print("Tarea corta completada")
    return "Resultado de tarea corta"

async def main():
    # Crear tareas
    tarea1 = asyncio.create_task(tarea_larga())
    tarea2 = asyncio.create_task(tarea_corta())
    
    # Esperar a que ambas tareas terminen
    await tarea1
    await tarea2
    
    # O esperar a que todas terminen con gather
    # resultados = await asyncio.gather(tarea1, tarea2)
    # print(f"Resultados: {resultados}")

asyncio.run(main())
```

### Espera con Timeout

```python
import asyncio

async def operacion_lenta():
    print("Iniciando operación lenta...")
    await asyncio.sleep(10)  # Simula una operación muy lenta
    print("Operación lenta completada")
    return "Resultado de operación lenta"

async def main():
    try:
        # Esperar la operación con un timeout de 3 segundos
        resultado = await asyncio.wait_for(operacion_lenta(), timeout=3)
        print(f"Resultado: {resultado}")
    except asyncio.TimeoutError:
        print("La operación tardó demasiado y se canceló")

asyncio.run(main())
```

### Ejecutar en Segundo Plano

```python
import asyncio

async def tarea_segundo_plano():
    """Tarea que se ejecuta en segundo plano."""
    for i in range(10):
        print(f"Tarea en segundo plano: paso {i}")
        await asyncio.sleep(1)

async def main():
    # Iniciar tarea en segundo plano
    tarea = asyncio.create_task(tarea_segundo_plano())
    
    # Hacer otras cosas mientras la tarea se ejecuta
    print("Haciendo trabajo principal...")
    await asyncio.sleep(3)
    print("Trabajo principal completado")
    
    # Cancelar la tarea en segundo plano
    tarea.cancel()
    try:
        await tarea
    except asyncio.CancelledError:
        print("Tarea en segundo plano cancelada")

asyncio.run(main())
```

### Comunicación entre Corrutinas

#### Colas Asíncronas

```python
import asyncio
import random

async def productor(cola):
    """Produce items y los coloca en la cola."""
    for i in range(10):
        item = random.randint(1, 100)
        await cola.put(item)
        print(f"Producido: {item}")
        await asyncio.sleep(random.uniform(0.1, 0.5))
    
    # Señal de finalización
    await cola.put(None)
    print("Productor terminado")

async def consumidor(cola):
    """Consume items de la cola."""
    while True:
        item = await cola.get()
        if item is None:  # Señal de finalización
            break
        print(f"Consumido: {item}")
        await asyncio.sleep(random.uniform(0.2, 0.7))
        cola.task_done()
    
    print("Consumidor terminado")

async def main():
    # Crear una cola
    cola = asyncio.Queue()
    
    # Crear tareas
    productor_task = asyncio.create_task(productor(cola))
    consumidor_task = asyncio.create_task(consumidor(cola))
    
    # Esperar a que ambas tareas terminen
    await asyncio.gather(productor_task, consumidor_task)

asyncio.run(main())
```

### Combinando asyncio con Threads

Para operaciones de E/S bloqueantes que no tienen versiones asíncronas.

```python
import asyncio
import time
import concurrent.futures

def operacion_bloqueante(segundos):
    """Función bloqueante que no puede ser async directamente."""
    print(f"Iniciando operación bloqueante de {segundos} segundos")
    time.sleep(segundos)  # Operación bloqueante
    print(f"Operación bloqueante de {segundos} segundos completada")
    return f"Resultado después de {segundos} segundos"

async def main():
    print("Iniciando programa")
    
    # Crear un executor de hilos
    with concurrent.futures.ThreadPoolExecutor() as executor:
        # Ejecutar operaciones bloqueantes en hilos
        loop = asyncio.get_running_loop()
        resultados = await asyncio.gather(
            loop.run_in_executor(executor, operacion_bloqueante, 3),
            loop.run_in_executor(executor, operacion_bloqueante, 2),
            loop.run_in_executor(executor, operacion_bloqueante, 1)
        )
    
    print(f"Resultados: {resultados}")

asyncio.run(main())
```

## Concurrencia con concurrent.futures

El módulo `concurrent.futures` proporciona una interfaz de alto nivel para ejecutar tareas de forma asíncrona utilizando hilos o procesos.

### ThreadPoolExecutor

```python
import concurrent.futures
import time
import requests

def descargar_url(url):
    """Descarga el contenido de una URL."""
    print(f"Descargando {url}")
    response = requests.get(url)
    return f"{url}: {len(response.content)} bytes"

# Lista de URLs para descargar
urls = [
    "https://www.python.org",
    "https://www.google.com",
    "https://www.github.com",
    "https://www.stackoverflow.com",
    "https://www.wikipedia.org"
]

# Usar ThreadPoolExecutor para descargas concurrentes
tiempo_inicio = time.time()

with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
    # Método map
    resultados = list(executor.map(descargar_url, urls))
    
    # Alternativa con submit
    # futuros = [executor.submit(descargar_url, url) for url in urls]
    # resultados = [futuro.result() for futuro in concurrent.futures.as_completed(futuros)]

tiempo_total = time.time() - tiempo_inicio

for resultado in resultados:
    print(resultado)

print(f"Tiempo total: {tiempo_total:.2f} segundos")
```

### ProcessPoolExecutor

```python
import concurrent.futures
import time
import math

def es_primo(n):
    """Verifica si un número es primo."""
    if n < 2:
        return False
    for i in range(2, int(math.sqrt(n)) + 1):
        if n % i == 0:
            return False
    return True

def contar_primos(inicio, fin):
    """Cuenta los números primos en un rango."""
    print(f"Procesando rango {inicio}-{fin}")
    contador = sum(1 for n in range(inicio, fin) if es_primo(n))
    return contador

# Rango grande para buscar primos
inicio = 1
fin = 10000000

# Dividir el trabajo en chunks
num_procesos = 8
tamaño_chunk = (fin - inicio) // num_procesos
rangos = [
    (inicio + i * tamaño_chunk, inicio + (i + 1) * tamaño_chunk)
    for i in range(num_procesos)
]

# Usar ProcessPoolExecutor para procesamiento paralelo
tiempo_inicio = time.time()

with concurrent.futures.ProcessPoolExecutor(max_workers=num_procesos) as executor:
    # Enviar cada rango a un proceso diferente
    futuros = [executor.submit(contar_primos, rango[0], rango[1]) for rango in rangos]
    
    # Recoger resultados a medida que se completan
    total_primos = 0
    for futuro in concurrent.futures.as_completed(futuros):
        total_primos += futuro.result()

tiempo_total = time.time() - tiempo_inicio

print(f"Total de números primos encontrados: {total_primos}")
print(f"Tiempo total: {tiempo_total:.2f} segundos")
```

## Ejemplo Práctico: Servidor Web Concurrente

Vamos a crear un servidor web concurrente utilizando diferentes enfoques.

### Versión con Threading

```python
import http.server
import socketserver
import threading
import time

class MiHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        """Maneja solicitudes GET."""
        print(f"[{threading.current_thread().name}] Solicitud recibida: {self.path}")
        
        # Simular procesamiento
        if self.path == "/rapido":
            time.sleep(0.1)  # Respuesta rápida
        elif self.path == "/lento":
            time.sleep(3)    # Respuesta lenta
        else:
            time.sleep(0.5)  # Respuesta normal
        
        # Enviar respuesta
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()
        self.wfile.write(f"<html><body><h1>Hola desde {self.path}</h1><p>Tiempo: {time.time()}</p></body></html>".encode())

# Crear servidor con múltiples hilos
class ThreadedHTTPServer(socketserver.ThreadingMixIn, http.server.HTTPServer):
    """Servidor HTTP que maneja cada solicitud en un hilo separado."""
    pass

def ejecutar_servidor_threading():
    # Configurar y ejecutar el servidor
    puerto = 8000
    servidor = ThreadedHTTPServer(("localhost", puerto), MiHandler)
    print(f"Servidor iniciado en http://localhost:{puerto}")
    
    try:
        servidor.serve_forever()
    except KeyboardInterrupt:
        print("Servidor detenido")
    finally:
        servidor.server_close()

if __name__ == "__main__":
    ejecutar_servidor_threading()
```

### Versión con asyncio

```python
import asyncio
import time
from aiohttp import web

async def handle_rapido(request):
    """Manejador para solicitudes rápidas."""
    print(f"Solicitud rápida recibida: {request.path}")
    await asyncio.sleep(0.1)  # Simular procesamiento rápido
    return web.Response(text=f"<html><body><h1>Respuesta rápida</h1><p>Tiempo: {time.time()}</p></body></html>", content_type="text/html")

async def handle_lento(request):
    """Manejador para solicitudes lentas."""
    print(f"Solicitud lenta recibida: {request.path}")
    await asyncio.sleep(3)  # Simular procesamiento lento
    return web.Response(text=f"<html><body><h1>Respuesta lenta</h1><p>Tiempo: {time.time()}</p></body></html>", content_type="text/html")

async def handle_normal(request):
    """Manejador para solicitudes normales."""
    print(f"Solicitud normal recibida: {request.path}")
    await asyncio.sleep(0.5)  # Simular procesamiento normal
    return web.Response(text=f"<html><body><h1>Respuesta normal</h1><p>Tiempo: {time.time()}</p></body></html>", content_type="text/html")

async def iniciar_servidor_asyncio():
    # Configurar la aplicación
    app = web.Application()
    app.add_routes([
        web.get('/rapido', handle_rapido),
        web.get('/lento', handle_lento),
        web.get('/{name}', handle_normal),
    ])
    
    # Iniciar el servidor
    runner = web.AppRunner(app)
    await runner.setup()
    site = web.TCPSite(runner, 'localhost', 8000)
    print("Servidor iniciado en http://localhost:8000")
    await site.start()
    
    # Mantener el servidor en ejecución
    try:
        while True:
            await asyncio.sleep(3600)  # Esperar indefinidamente
    except asyncio.CancelledError:
        print("Servidor detenido")
    finally:
        await runner.cleanup()

if __name__ == "__main__":
    asyncio.run(iniciar_servidor_asyncio())
```

## Buenas Prácticas

1. **Elegir el Enfoque Adecuado**:
   - Usa `threading` para tareas de E/S concurrentes cuando la simplicidad es importante
   - Usa `multiprocessing` para tareas intensivas de CPU que requieren paralelismo real
   - Usa `asyncio` para aplicaciones con muchas operaciones de E/S, especialmente servidores web
   - Usa `concurrent.futures` para una interfaz más simple y unificada

2. **Evitar Condiciones de Carrera**:
   - Usa mecanismos de sincronización (locks, semáforos, etc.) para proteger recursos compartidos
   - Minimiza el código en secciones críticas para reducir contención
   - Considera usar estructuras de datos inmutables o thread-safe

3. **Manejo de Recursos**:
   - Libera recursos (cerrar archivos, conexiones, etc.) incluso si ocurren excepciones
   - Usa `with` para gestionar recursos automáticamente
   - Cierra explícitamente pools de hilos y procesos cuando ya no se necesiten

4. **Evitar Deadlocks**:
   - Adquiere múltiples locks siempre en el mismo orden
   - Usa timeouts al adquirir locks para evitar bloqueos indefinidos
   - Considera usar `RLock` si un mismo hilo necesita adquirir un lock múltiples veces

5. **Diseño para Concurrencia**:
   - Divide el trabajo en unidades independientes
   - Minimiza la comunicación y dependencias entre hilos/procesos/tareas
   - Usa patrones como productor-consumidor, mapeo-reducción, etc.

6. **Depuración**:
   - Usa nombres descriptivos para hilos y procesos
   - Implementa logging detallado para rastrear la ejecución
   - Considera usar herramientas como `threading.settrace()` o depuradores concurrentes

7. **Rendimiento**:
   - Mide y compara diferentes enfoques para tu caso específico
   - Ajusta el número de hilos/procesos según la carga y recursos disponibles
   - Considera el overhead de creación y comunicación

8. **Cancelación y Timeouts**:
   - Implementa mecanismos para cancelar tareas largas
   - Usa timeouts para operaciones que podrían bloquearse indefinidamente
   - Maneja correctamente las excepciones relacionadas con cancelación

## Recursos Adicionales

- [Documentación oficial de threading](https://docs.python.org/es/3/library/threading.html)
- [Documentación oficial de multiprocessing](https://docs.python.org/es/3/library/multiprocessing.html)
- [Documentación oficial de asyncio](https://docs.python.org/es/3/library/asyncio.html)
- [Documentación oficial de concurrent.futures](https://docs.python.org/es/3/library/concurrent.futures.html)
- [Python Concurrency From the Ground Up (PyCon 2015)](https://www.youtube.com/watch?v=MCs5OvhV9S4)
- [Thinking About Concurrency (Raymond Hettinger)](https://www.youtube.com/watch?v=Bv25Dwe84g0)
- [Async IO in Python: A Complete Walkthrough](https://realpython.com/async-io-python/)
- [Speed Up Your Python Program With Concurrency](https://realpython.com/python-concurrency/)
- [Python Parallel Programming Cookbook](https://www.packtpub.com/product/python-parallel-programming-cookbook/9781785289583)

---

En la siguiente sección, exploraremos las pruebas y testing en Python, incluyendo unittest, pytest, mocking y más.