# Concurrencia Avanzada en Python

## Introducción

La concurrencia es la capacidad de un programa para ejecutar múltiples tareas que se solapan en el tiempo. En Python, existen varios modelos y herramientas para implementar concurrencia, cada uno con sus propias ventajas y casos de uso. Este módulo explora los conceptos avanzados de concurrencia en Python, incluyendo threading, multiprocessing, asyncio, y patrones de diseño concurrentes.

La concurrencia en Python es particularmente interesante debido al Global Interpreter Lock (GIL), que limita la ejecución de bytecode de Python a un solo hilo a la vez. Esto tiene implicaciones importantes para el rendimiento y el diseño de aplicaciones concurrentes.

## Modelos de Concurrencia en Python

### Threading vs Multiprocessing vs Asyncio

Python ofrece tres modelos principales de concurrencia:

1. **Threading**: Utiliza hilos para ejecutar tareas concurrentemente dentro del mismo proceso.
2. **Multiprocessing**: Utiliza procesos separados para ejecutar tareas en paralelo.
3. **Asyncio**: Utiliza corrutinas y un bucle de eventos para la concurrencia cooperativa.

Cada modelo tiene sus propias ventajas y desventajas:

```python
# Ejemplo comparativo de los tres modelos

# Threading - Bueno para tareas I/O bound
import threading
import time

def io_bound_task(name):
    print(f"Tarea {name} iniciada")
    time.sleep(1)  # Simula operación I/O
    print(f"Tarea {name} completada")

def threading_example():
    start = time.time()
    threads = []
    for i in range(10):
        thread = threading.Thread(target=io_bound_task, args=(f"thread-{i}",))
        threads.append(thread)
        thread.start()
    
    for thread in threads:
        thread.join()
    
    print(f"Threading completado en {time.time() - start:.2f} segundos")

# Multiprocessing - Bueno para tareas CPU bound
import multiprocessing

def cpu_bound_task(name):
    print(f"Tarea {name} iniciada")
    # Simula operación CPU intensiva
    result = 0
    for i in range(10**7):
        result += i
    print(f"Tarea {name} completada")

def multiprocessing_example():
    start = time.time()
    processes = []
    for i in range(10):
        process = multiprocessing.Process(target=cpu_bound_task, args=(f"process-{i}",))
        processes.append(process)
        process.start()
    
    for process in processes:
        process.join()
    
    print(f"Multiprocessing completado en {time.time() - start:.2f} segundos")

# Asyncio - Bueno para tareas I/O bound con alta concurrencia
import asyncio

async def async_io_bound_task(name):
    print(f"Tarea {name} iniciada")
    await asyncio.sleep(1)  # Simula operación I/O asíncrona
    print(f"Tarea {name} completada")

async def asyncio_example():
    start = time.time()
    tasks = [async_io_bound_task(f"async-{i}") for i in range(10)]
    await asyncio.gather(*tasks)
    print(f"Asyncio completado en {time.time() - start:.2f} segundos")

# Ejecutar ejemplos
if __name__ == "__main__":
    print("=== Threading Example ===")
    threading_example()
    
    print("\n=== Multiprocessing Example ===")
    multiprocessing_example()
    
    print("\n=== Asyncio Example ===")
    asyncio.run(asyncio_example())
```

### El Global Interpreter Lock (GIL)

El GIL es un mecanismo en el intérprete de CPython que asegura que solo un hilo ejecute código Python a la vez. Esto tiene implicaciones importantes:

- Los hilos son efectivos para tareas I/O bound, pero no para tareas CPU bound.
- Para tareas CPU bound, multiprocessing es generalmente más eficiente.
- El GIL no afecta a las operaciones que liberan el GIL, como algunas operaciones de NumPy o llamadas a C.

```python
import threading
import time
import numpy as np

# Tarea CPU bound con Python puro (afectada por el GIL)
def cpu_bound_python():
    result = 0
    for i in range(10**7):
        result += i
    return result

# Tarea CPU bound con NumPy (libera el GIL)
def cpu_bound_numpy():
    return np.sum(np.arange(10**7))

def measure_time(func, num_threads):
    start = time.time()
    threads = []
    
    for _ in range(num_threads):
        thread = threading.Thread(target=func)
        threads.append(thread)
        thread.start()
    
    for thread in threads:
        thread.join()
    
    return time.time() - start

# Comparar rendimiento
if __name__ == "__main__":
    # Ejecución secuencial
    print("Ejecución secuencial:")
    python_time = measure_time(cpu_bound_python, 1)
    numpy_time = measure_time(cpu_bound_numpy, 1)
    print(f"Python puro: {python_time:.2f} segundos")
    print(f"NumPy: {numpy_time:.2f} segundos")
    
    # Ejecución con múltiples hilos
    print("\nEjecución con 4 hilos:")
    python_time_threads = measure_time(cpu_bound_python, 4)
    numpy_time_threads = measure_time(cpu_bound_numpy, 4)
    print(f"Python puro: {python_time_threads:.2f} segundos")
    print(f"NumPy: {numpy_time_threads:.2f} segundos")
    
    # Comparación de speedup
    python_speedup = python_time / python_time_threads
    numpy_speedup = numpy_time / numpy_time_threads
    print(f"\nSpeedup con Python puro: {python_speedup:.2f}x")
    print(f"Speedup con NumPy: {numpy_speedup:.2f}x")
```

## Threading Avanzado

### Sincronización y Comunicación entre Hilos

Python proporciona varias primitivas para sincronización y comunicación entre hilos:

```python
import threading
import queue
import time
import random

# Ejemplo de sincronización con Lock
class Counter:
    def __init__(self):
        self.value = 0
        self.lock = threading.Lock()
    
    def increment(self):
        with self.lock:
            current = self.value
            time.sleep(0.001)  # Simula trabajo
            self.value = current + 1

def worker(counter, num_increments):
    for _ in range(num_increments):
        counter.increment()

def demonstrate_lock():
    counter = Counter()
    threads = []
    num_threads = 10
    increments_per_thread = 100
    
    for _ in range(num_threads):
        thread = threading.Thread(target=worker, args=(counter, increments_per_thread))
        threads.append(thread)
        thread.start()
    
    for thread in threads:
        thread.join()
    
    expected = num_threads * increments_per_thread
    print(f"Valor esperado: {expected}, Valor actual: {counter.value}")

# Ejemplo de comunicación con Queue
def producer(q, num_items):
    for i in range(num_items):
        item = random.randint(1, 100)
        q.put(item)
        print(f"Producido: {item}")
        time.sleep(random.random() * 0.1)

def consumer(q, name):
    while True:
        try:
            item = q.get(timeout=1)
            print(f"Consumidor {name} consumió: {item}")
            q.task_done()
            time.sleep(random.random() * 0.2)
        except queue.Empty:
            break

def demonstrate_queue():
    q = queue.Queue()
    num_producers = 2
    num_consumers = 3
    num_items = 20
    
    # Iniciar productores
    producer_threads = []
    for i in range(num_producers):
        thread = threading.Thread(target=producer, args=(q, num_items))
        producer_threads.append(thread)
        thread.start()
    
    # Iniciar consumidores
    consumer_threads = []
    for i in range(num_consumers):
        thread = threading.Thread(target=consumer, args=(q, i))
        consumer_threads.append(thread)
        thread.daemon = True  # Hilos daemon terminan cuando el programa principal termina
        thread.start()
    
    # Esperar a que los productores terminen
    for thread in producer_threads:
        thread.join()
    
    # Esperar a que la cola se vacíe
    q.join()
    
    print("Todos los items han sido procesados")

if __name__ == "__main__":
    print("=== Demostración de Lock ===")
    demonstrate_lock()
    
    print("\n=== Demostración de Queue ===")
    demonstrate_queue()
```

### Patrones de Diseño con Threading

#### Thread Pool

```python
import concurrent.futures
import time
import random

def task(n):
    print(f"Ejecutando tarea {n}")
    sleep_time = random.uniform(0.5, 2)
    time.sleep(sleep_time)  # Simula trabajo
    return f"Resultado de tarea {n}: {sleep_time:.2f}s"

def thread_pool_example():
    tasks = range(10)
    
    # Usando ThreadPoolExecutor
    with concurrent.futures.ThreadPoolExecutor(max_workers=3) as executor:
        # Método 1: map
        print("Método map:")
        for result in executor.map(task, tasks):
            print(result)
        
        print("\nMétodo submit:")
        # Método 2: submit y as_completed
        futures = [executor.submit(task, i) for i in tasks]
        for future in concurrent.futures.as_completed(futures):
            print(future.result())

if __name__ == "__main__":
    thread_pool_example()
```

#### Worker Pool con Queue

```python
import threading
import queue
import time
import random

def worker(task_queue, result_queue, worker_id):
    while True:
        try:
            task = task_queue.get(timeout=1)
            print(f"Worker {worker_id} procesando tarea {task}")
            # Simular procesamiento
            time.sleep(random.uniform(0.5, 1.5))
            result = task * 2
            result_queue.put((task, result))
            task_queue.task_done()
        except queue.Empty:
            break

def worker_pool_example():
    num_workers = 4
    num_tasks = 20
    
    # Crear colas
    task_queue = queue.Queue()
    result_queue = queue.Queue()
    
    # Añadir tareas a la cola
    for i in range(num_tasks):
        task_queue.put(i)
    
    # Crear workers
    workers = []
    for i in range(num_workers):
        thread = threading.Thread(target=worker, args=(task_queue, result_queue, i))
        thread.daemon = True
        workers.append(thread)
        thread.start()
    
    # Esperar a que todas las tareas sean procesadas
    task_queue.join()
    
    # Recoger resultados
    results = {}
    while not result_queue.empty():
        task, result = result_queue.get()
        results[task] = result
    
    # Mostrar resultados ordenados
    for task in sorted(results.keys()):
        print(f"Tarea {task} -> Resultado {results[task]}")

if __name__ == "__main__":
    worker_pool_example()
```

## Multiprocessing Avanzado

### Comunicación entre Procesos

Python proporciona varios mecanismos para la comunicación entre procesos:

```python
import multiprocessing as mp
import time
import random

# Ejemplo con Queue
def producer_process(q):
    for i in range(10):
        item = random.randint(1, 100)
        q.put(item)
        print(f"Proceso productor generó: {item}")
        time.sleep(random.random() * 0.5)

def consumer_process(q):
    while True:
        try:
            item = q.get(timeout=2)
            print(f"Proceso consumidor recibió: {item}")
            time.sleep(random.random() * 0.5)
        except:
            print("Consumidor terminando")
            break

def queue_example():
    q = mp.Queue()
    
    # Crear procesos
    producer = mp.Process(target=producer_process, args=(q,))
    consumer = mp.Process(target=consumer_process, args=(q,))
    
    # Iniciar procesos
    producer.start()
    consumer.start()
    
    # Esperar a que el productor termine
    producer.join()
    
    # Esperar un poco más y terminar el consumidor
    time.sleep(3)
    consumer.terminate()
    consumer.join()

# Ejemplo con Pipe
def sender_process(conn):
    for i in range(10):
        msg = f"Mensaje {i}"
        conn.send(msg)
        print(f"Enviado: {msg}")
        time.sleep(random.random() * 0.5)
    conn.close()

def receiver_process(conn):
    while True:
        if conn.poll(1):  # Verificar si hay datos disponibles (timeout 1 segundo)
            msg = conn.recv()
            print(f"Recibido: {msg}")
        else:
            if not conn.closed:
                continue
            print("Receptor terminando")
            break

def pipe_example():
    # Crear un pipe
    parent_conn, child_conn = mp.Pipe()
    
    # Crear procesos
    sender = mp.Process(target=sender_process, args=(parent_conn,))
    receiver = mp.Process(target=receiver_process, args=(child_conn,))
    
    # Iniciar procesos
    sender.start()
    receiver.start()
    
    # Esperar a que terminen
    sender.join()
    receiver.join(timeout=2)
    if receiver.is_alive():
        receiver.terminate()

if __name__ == "__main__":
    print("=== Queue Example ===")
    queue_example()
    
    print("\n=== Pipe Example ===")
    pipe_example()
```

### Compartir Estado entre Procesos

```python
import multiprocessing as mp
import time

# Ejemplo con Value y Array
def worker_value(shared_value, lock):
    for _ in range(100):
        with lock:
            shared_value.value += 1

def value_example():
    # Crear un valor compartido
    shared_value = mp.Value('i', 0)  # 'i' para entero
    lock = mp.Lock()
    
    # Crear procesos
    processes = [mp.Process(target=worker_value, args=(shared_value, lock)) for _ in range(4)]
    
    # Iniciar procesos
    for p in processes:
        p.start()
    
    # Esperar a que terminen
    for p in processes:
        p.join()
    
    print(f"Valor final: {shared_value.value}")

def worker_array(shared_array, lock, worker_id):
    for _ in range(100):
        with lock:
            for i in range(len(shared_array)):
                shared_array[i] += 1

def array_example():
    # Crear un array compartido
    shared_array = mp.Array('i', [0, 0, 0, 0, 0])  # Array de 5 enteros
    lock = mp.Lock()
    
    # Crear procesos
    processes = [mp.Process(target=worker_array, args=(shared_array, lock, i)) for i in range(4)]
    
    # Iniciar procesos
    for p in processes:
        p.start()
    
    # Esperar a que terminen
    for p in processes:
        p.join()
    
    print(f"Array final: {list(shared_array)}")

# Ejemplo con Manager
def worker_dict(shared_dict, lock, worker_id):
    for i in range(10):
        with lock:
            if worker_id in shared_dict:
                shared_dict[worker_id].append(i)
            else:
                shared_dict[worker_id] = [i]
        time.sleep(0.1)

def manager_example():
    # Crear un manager
    with mp.Manager() as manager:
        # Crear un diccionario compartido
        shared_dict = manager.dict()
        lock = manager.Lock()
        
        # Crear procesos
        processes = [mp.Process(target=worker_dict, args=(shared_dict, lock, i)) for i in range(4)]
        
        # Iniciar procesos
        for p in processes:
            p.start()
        
        # Esperar a que terminen
        for p in processes:
            p.join()
        
        # Convertir a diccionario normal para imprimir
        result = dict(shared_dict)
        print(f"Diccionario final: {result}")

if __name__ == "__main__":
    print("=== Value Example ===")
    value_example()
    
    print("\n=== Array Example ===")
    array_example()
    
    print("\n=== Manager Example ===")
    manager_example()
```

### Pool de Procesos

```python
import multiprocessing as mp
import time
import random

def cpu_bound_task(n):
    print(f"Procesando {n}")
    start = time.time()
    # Simular tarea CPU intensiva
    count = 0
    for i in range(10**7):
        count += i
    duration = time.time() - start
    return f"Tarea {n} completada en {duration:.2f}s con resultado {count}"

def pool_example():
    # Crear un pool con el número de procesos igual al número de CPUs
    num_processes = mp.cpu_count()
    print(f"Creando pool con {num_processes} procesos")
    
    with mp.Pool(processes=num_processes) as pool:
        # Método 1: map
        print("\nUsando map:")
        results = pool.map(cpu_bound_task, range(10))
        for result in results:
            print(result)
        
        # Método 2: apply_async
        print("\nUsando apply_async:")
        results = [pool.apply_async(cpu_bound_task, args=(i,)) for i in range(10)]
        for result in results:
            print(result.get())  # .get() espera el resultado
        
        # Método 3: map_async
        print("\nUsando map_async:")
        result_obj = pool.map_async(cpu_bound_task, range(10))
        # Hacer algo mientras se procesan las tareas
        while not result_obj.ready():
            print("Esperando resultados...")
            time.sleep(1)
        
        # Obtener todos los resultados
        results = result_obj.get()
        for result in results:
            print(result)

if __name__ == "__main__":
    pool_example()
```

## Asyncio Avanzado

### Corrutinas y Tareas

```python
import asyncio
import time

async def say_after(delay, what):
    await asyncio.sleep(delay)
    print(what)
    return what

async def main():
    # Forma 1: await directo
    print(f"Inicio: {time.strftime('%X')}")
    
    await say_after(1, 'hello')
    await say_after(2, 'world')
    
    print(f"Fin forma 1: {time.strftime('%X')}")
    
    # Forma 2: crear tareas
    print(f"Inicio forma 2: {time.strftime('%X')}")
    
    task1 = asyncio.create_task(say_after(1, 'hello'))
    task2 = asyncio.create_task(say_after(2, 'world'))
    
    # Esperar a que ambas tareas terminen
    await task1
    await task2
    
    print(f"Fin forma 2: {time.strftime('%X')}")
    
    # Forma 3: gather
    print(f"Inicio forma 3: {time.strftime('%X')}")
    
    results = await asyncio.gather(
        say_after(1, 'hello'),
        say_after(2, 'world')
    )
    
    print(f"Resultados: {results}")
    print(f"Fin forma 3: {time.strftime('%X')}")

if __name__ == "__main__":
    asyncio.run(main())
```

### Manejo de Excepciones en Asyncio

```python
import asyncio
import random

async def risky_task(task_id):
    await asyncio.sleep(random.uniform(0.1, 0.5))
    # Simular error aleatorio
    if random.random() < 0.3:  # 30% de probabilidad de error
        raise Exception(f"Error en tarea {task_id}")
    return f"Resultado de tarea {task_id}"

async def handle_exceptions_individually():
    tasks = [risky_task(i) for i in range(10)]
    
    results = []
    for task_id, task_coro in enumerate(tasks):
        try:
            result = await task_coro
            results.append(result)
            print(f"Tarea {task_id} completada: {result}")
        except Exception as e:
            print(f"Tarea {task_id} falló: {e}")
    
    return results

async def handle_exceptions_with_gather():
    tasks = [risky_task(i) for i in range(10)]
    
    # Método 1: return_exceptions=True
    print("\nUsando gather con return_exceptions=True:")
    results = await asyncio.gather(*tasks, return_exceptions=True)
    
    for i, result in enumerate(results):
        if isinstance(result, Exception):
            print(f"Tarea {i} falló: {result}")
        else:
            print(f"Tarea {i} completada: {result}")
    
    # Método 2: try/except con gather
    print("\nUsando try/except con gather:")
    try:
        results = await asyncio.gather(*tasks)
        print("Todas las tareas completadas exitosamente")
    except Exception as e:
        print(f"Al menos una tarea falló: {e}")

async def handle_exceptions_with_tasks():
    # Crear tareas
    tasks = [asyncio.create_task(risky_task(i)) for i in range(10)]
    
    # Esperar a que todas las tareas terminen, incluso si algunas fallan
    done, pending = await asyncio.wait(tasks, return_when=asyncio.ALL_COMPLETED)
    
    # Procesar resultados y excepciones
    for task in done:
        try:
            result = task.result()
            print(f"Tarea completada: {result}")
        except Exception as e:
            print(f"Tarea falló: {e}")

async def main():
    print("=== Manejo individual de excepciones ===")
    await handle_exceptions_individually()
    
    print("\n=== Manejo de excepciones con gather ===")
    await handle_exceptions_with_gather()
    
    print("\n=== Manejo de excepciones con tasks ===")
    await handle_exceptions_with_tasks()

if __name__ == "__main__":
    asyncio.run(main())
```

### Patrones Avanzados con Asyncio

#### Semáforos y Limitación de Concurrencia

```python
import asyncio
import aiohttp
import time

async def fetch_url(url, session, semaphore):
    async with semaphore:
        print(f"Fetching {url}")
        try:
            async with session.get(url) as response:
                return await response.text()
        except Exception as e:
            print(f"Error fetching {url}: {e}")
            return None

async def semaphore_example():
    # Lista de URLs a consultar
    urls = [
        "https://www.python.org",
        "https://www.google.com",
        "https://www.github.com",
        "https://www.stackoverflow.com",
        "https://www.wikipedia.org",
        "https://www.reddit.com",
        "https://www.twitter.com",
        "https://www.facebook.com",
        "https://www.instagram.com",
        "https://www.linkedin.com"
    ]
    
    # Crear un semáforo para limitar a 3 conexiones simultáneas
    semaphore = asyncio.Semaphore(3)
    
    async with aiohttp.ClientSession() as session:
        tasks = [fetch_url(url, session, semaphore) for url in urls]
        results = await asyncio.gather(*tasks)
    
    # Contar cuántas consultas fueron exitosas
    successful = sum(1 for r in results if r is not None)
    print(f"Consultas exitosas: {successful}/{len(urls)}")

if __name__ == "__main__":
    start = time.time()
    asyncio.run(semaphore_example())
    print(f"Tiempo total: {time.time() - start:.2f} segundos")
```

#### Productor-Consumidor Asíncrono

```python
import asyncio
import random

async def producer(queue, id):
    for i in range(5):
        item = f"Producer {id} - Item {i}"
        # Simular tiempo de producción variable
        await asyncio.sleep(random.uniform(0.1, 0.5))
        await queue.put(item)
        print(f"Produced: {item}")
    
    # Señal de finalización para este productor
    await queue.put(None)

async def consumer(queue, id):
    while True:
        # Esperar por un item
        item = await queue.get()
        
        # Verificar señal de finalización
        if item is None:
            print(f"Consumer {id} shutting down")
            queue.task_done()
            break
        
        # Simular tiempo de procesamiento variable
        await asyncio.sleep(random.uniform(0.2, 0.6))
        print(f"Consumer {id} consumed: {item}")
        queue.task_done()

async def producer_consumer_example():
    # Crear una cola asíncrona
    queue = asyncio.Queue(maxsize=10)  # Limitar tamaño de cola
    
    # Crear productores y consumidores
    producers = [producer(queue, i) for i in range(3)]
    consumers = [consumer(queue, i) for i in range(2)]
    
    # Iniciar productores
    producer_tasks = [asyncio.create_task(p) for p in producers]
    
    # Iniciar consumidores
    consumer_tasks = [asyncio.create_task(c) for c in consumers]
    
    # Esperar a que los productores terminen
    await asyncio.gather(*producer_tasks)
    
    # Enviar señales de finalización a los consumidores (uno por consumidor)
    for _ in range(len(consumers)):
        await queue.put(None)
    
    # Esperar a que los consumidores terminen
    await asyncio.gather(*consumer_tasks)
    
    print("Todos los productores y consumidores han terminado")

if __name__ == "__main__":
    asyncio.run(producer_consumer_example())
```

#### Timeouts y Cancelación

```python
import asyncio
import random

async def long_running_task(task_id):
    try:
        delay = random.uniform(1, 5)
        print(f"Tarea {task_id} iniciada, durará {delay:.2f}s")
        await asyncio.sleep(delay)
        print(f"Tarea {task_id} completada")
        return f"Resultado de tarea {task_id}"
    except asyncio.CancelledError:
        print(f"Tarea {task_id} cancelada")
        raise  # Re-lanzar la excepción para que sea manejada correctamente

async def timeout_example():
    # Crear tareas
    tasks = [long_running_task(i) for i in range(5)]
    
    # Ejecutar tareas con timeout
    try:
        results = await asyncio.wait_for(asyncio.gather(*tasks), timeout=3)
        print(f"Todas las tareas completadas: {results}")
    except asyncio.TimeoutError:
        print("Timeout alcanzado, algunas tareas no completaron a tiempo")

async def cancellation_example():
    # Crear tareas
    tasks = [asyncio.create_task(long_running_task(i)) for i in range(5)]
    
    # Esperar un poco y luego cancelar algunas tareas
    await asyncio.sleep(2)
    
    # Cancelar tareas con ID par
    for i, task in enumerate(tasks):
        if i % 2 == 0 and not task.done():
            print(f"Cancelando tarea {i}")
            task.cancel()
    
    # Esperar a que todas las tareas terminen (completadas o canceladas)
    done, pending = await asyncio.wait(tasks, return_when=asyncio.ALL_COMPLETED)
    
    # Verificar resultados
    for i, task in enumerate(tasks):
        try:
            if task.done() and not task.cancelled():
                result = task.result()
                print(f"Tarea {i} completada con resultado: {result}")
            elif task.cancelled():
                print(f"Tarea {i} fue cancelada")
            else:
                print(f"Tarea {i} está en estado desconocido")
        except Exception as e:
            print(f"Tarea {i} lanzó excepción: {e}")

async def shield_example():
    # shield protege una corrutina de ser cancelada
    task = long_running_task(999)
    shielded_task = asyncio.shield(task)
    
    # Intentar cancelar después de un tiempo
    await asyncio.sleep(1)
    print("Intentando cancelar la tarea protegida")
    
    # Esto no cancelará la tarea subyacente debido al shield
    try:
        shielded_task.cancel()
        await shielded_task
    except asyncio.CancelledError:
        print("La tarea protegida fue marcada como cancelada")
        # Pero la tarea original sigue ejecutándose
        try:
            result = await task
            print(f"La tarea original completó con: {result}")
        except Exception as e:
            print(f"La tarea original falló: {e}")

async def main():
    print("=== Timeout Example ===")
    await timeout_example()
    
    print("\n=== Cancellation Example ===")
    await cancellation_example()
    
    print("\n=== Shield Example ===")
    await shield_example()

if __name__ == "__main__":
    asyncio.run(main())
```

## Integración de Modelos de Concurrencia

### Combinando Threading y Asyncio

```python
import asyncio
import threading
import time

# Función para ejecutar en un hilo separado
def blocking_io():
    print(f"Hilo: iniciando tarea IO-bound")
    time.sleep(2)  # Simular operación IO bloqueante
    print(f"Hilo: tarea IO-bound completada")
    return "Resultado de IO"

# Función para ejecutar código asíncrono en un hilo separado
def run_async_in_thread(coro):
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        return loop.run_until_complete(coro)
    finally:
        loop.close()

# Corrutina asíncrona
async def async_task():
    print("Corrutina: iniciando tarea asíncrona")
    await asyncio.sleep(1)  # Operación asíncrona
    print("Corrutina: tarea asíncrona completada")
    return "Resultado asíncrono"

# Ejecutar tarea bloqueante en un hilo desde asyncio
async def run_blocking_in_executor():
    print("Ejecutando tarea bloqueante en executor")
    loop = asyncio.get_running_loop()
    # Ejecutar la función bloqueante en un ThreadPoolExecutor
    result = await loop.run_in_executor(None, blocking_io)
    print(f"Resultado del executor: {result}")

# Ejecutar tarea asíncrona en un hilo separado
def run_async_task_in_thread():
    print("Ejecutando tarea asíncrona en hilo separado")
    result = run_async_in_thread(async_task())
    print(f"Resultado del hilo: {result}")

async def main():
    # 1. Ejecutar tarea bloqueante en executor desde asyncio
    await run_blocking_in_executor()
    
    # 2. Ejecutar tarea asíncrona en hilo separado
    thread = threading.Thread(target=run_async_task_in_thread)
    thread.start()
    thread.join()
    
    # 3. Ejecutar múltiples tareas bloqueantes en paralelo
    print("\nEjecutando múltiples tareas bloqueantes en paralelo")
    loop = asyncio.get_running_loop()
    tasks = [loop.run_in_executor(None, blocking_io) for _ in range(3)]
    results = await asyncio.gather(*tasks)
    print(f"Resultados: {results}")

if __name__ == "__main__":
    asyncio.run(main())
```

### Combinando Multiprocessing y Asyncio

```python
import asyncio
import multiprocessing as mp
import time

# Función CPU-bound para ejecutar en un proceso separado
def cpu_bound_task(n):
    print(f"Proceso: iniciando tarea CPU-bound {n}")
    # Simular tarea CPU intensiva
    result = 0
    for i in range(10**7):
        result += i
    print(f"Proceso: tarea CPU-bound {n} completada")
    return result

# Ejecutar tarea CPU-bound en un proceso desde asyncio
async def run_in_process(fn, *args):
    loop = asyncio.get_running_loop()
    # Usar un futuro para recibir el resultado
    future = loop.create_future()
    
    # Función que se ejecutará cuando el proceso termine
    def done_callback(result):
        loop.call_soon_threadsafe(future.set_result, result)
    
    # Crear y configurar el proceso
    process = mp.Process(target=process_wrapper, args=(fn, done_callback, *args))
    process.start()
    
    # Esperar el resultado
    return await future

# Wrapper para ejecutar la función y enviar el resultado a través del callback
def process_wrapper(fn, callback, *args):
    result = fn(*args)
    callback(result)

async def main():
    # Ejecutar una tarea CPU-bound en un proceso
    print("Ejecutando tarea CPU-bound en un proceso")
    result = await run_in_process(cpu_bound_task, 1)
    print(f"Resultado: {result}")
    
    # Ejecutar múltiples tareas CPU-bound en procesos paralelos
    print("\nEjecutando múltiples tareas CPU-bound en procesos paralelos")
    tasks = [run_in_process(cpu_bound_task, i) for i in range(4)]
    results = await asyncio.gather(*tasks)
    print(f"Resultados: {results}")

if __name__ == "__main__":
    # Necesario para multiprocessing en Windows
    mp.set_start_method('spawn', force=True)
    asyncio.run(main())
```

## Patrones de Diseño Concurrentes

### Patrón Pipeline

```python
import asyncio
import random

# Etapas del pipeline
async def producer(queue, num_items):
    for i in range(num_items):
        item = random.randint(1, 100)
        await queue.put(item)
        print(f"Producido: {item}")
        await asyncio.sleep(random.uniform(0.1, 0.3))
    
    # Señal de finalización
    await queue.put(None)

async def stage1(input_queue, output_queue):
    while True:
        item = await input_queue.get()
        
        # Verificar señal de finalización
        if item is None:
            print("Etapa 1 finalizando")
            await output_queue.put(None)  # Propagar señal
            break
        
        # Procesar item
        result = item * 2
        print(f"Etapa 1: {item} -> {result}")
        await asyncio.sleep(random.uniform(0.1, 0.2))  # Simular procesamiento
        await output_queue.put(result)

async def stage2(input_queue, output_queue):
    while True:
        item = await input_queue.get()
        
        # Verificar señal de finalización
        if item is None:
            print("Etapa 2 finalizando")
            await output_queue.put(None)  # Propagar señal
            break
        
        # Procesar item
        result = item + 10
        print(f"Etapa 2: {item} -> {result}")
        await asyncio.sleep(random.uniform(0.1, 0.2))  # Simular procesamiento
        await output_queue.put(result)

async def stage3(input_queue, output_queue):
    while True:
        item = await input_queue.get()
        
        # Verificar señal de finalización
        if item is None:
            print("Etapa 3 finalizando")
            await output_queue.put(None)  # Propagar señal
            break
        
        # Procesar item
        result = item ** 2
        print(f"Etapa 3: {item} -> {result}")
        await asyncio.sleep(random.uniform(0.1, 0.2))  # Simular procesamiento
        await output_queue.put(result)

async def consumer(queue):
    results = []
    while True:
        item = await queue.get()
        
        # Verificar señal de finalización
        if item is None:
            print("Consumidor finalizando")
            break
        
        # Consumir item
        print(f"Consumido: {item}")
        results.append(item)
    
    return results

async def pipeline_example():
    # Crear colas para conectar las etapas
    queue1 = asyncio.Queue()  # Productor -> Etapa 1
    queue2 = asyncio.Queue()  # Etapa 1 -> Etapa 2
    queue3 = asyncio.Queue()  # Etapa 2 -> Etapa 3
    queue4 = asyncio.Queue()  # Etapa 3 -> Consumidor
    
    # Iniciar todas las etapas
    producer_task = asyncio.create_task(producer(queue1, 10))
    stage1_task = asyncio.create_task(stage1(queue1, queue2))
    stage2_task = asyncio.create_task(stage2(queue2, queue3))
    stage3_task = asyncio.create_task(stage3(queue3, queue4))
    consumer_task = asyncio.create_task(consumer(queue4))
    
    # Esperar a que todas las etapas terminen
    await producer_task
    await stage1_task
    await stage2_task
    await stage3_task
    results = await consumer_task
    
    print(f"Resultados finales: {results}")
    return results

if __name__ == "__main__":
    asyncio.run(pipeline_example())
```

### Patrón Fan-Out, Fan-In

```python
import asyncio
import random
import time

# Generador de tareas
async def generate_tasks(num_tasks):
    tasks = []
    for i in range(num_tasks):
        tasks.append(i)
    return tasks

# Worker que procesa tareas
async def worker(name, queue, results):
    while True:
        # Obtener tarea de la cola
        task = await queue.get()
        
        # Verificar señal de finalización
        if task is None:
            print(f"Worker {name} finalizando")
            break
        
        # Procesar tarea
        print(f"Worker {name} procesando tarea {task}")
        # Simular tiempo de procesamiento variable
        await asyncio.sleep(random.uniform(0.5, 2.0))
        result = task * 2  # Operación simple para demostración
        
        # Almacenar resultado
        await results.put((task, result))
        
        # Marcar tarea como completada
        queue.task_done()

# Recolector de resultados
async def result_collector(results_queue, num_tasks):
    collected_results = {}
    tasks_processed = 0
    
    while tasks_processed < num_tasks:
        # Obtener resultado
        task_id, result = await results_queue.get()
        
        # Almacenar resultado
        collected_results[task_id] = result
        print(f"Resultado recolectado para tarea {task_id}: {result}")
        
        # Incrementar contador
        tasks_processed += 1
        
        # Marcar como procesado
        results_queue.task_done()
    
    return collected_results

async def fan_out_fan_in_example():
    # Parámetros
    num_workers = 4
    num_tasks = 20
    
    # Crear colas
    task_queue = asyncio.Queue()
    results_queue = asyncio.Queue()
    
    # Generar tareas
    tasks = await generate_tasks(num_tasks)
    for task in tasks:
        await task_queue.put(task)
    
    # Crear workers (fan-out)
    workers = []
    for i in range(num_workers):
        worker_name = f"worker-{i}"
        workers.append(asyncio.create_task(
            worker(worker_name, task_queue, results_queue)
        ))
    
    # Crear recolector de resultados (fan-in)
    collector = asyncio.create_task(
        result_collector(results_queue, num_tasks)
    )
    
    # Esperar a que todas las tareas sean procesadas
    await task_queue.join()
    
    # Enviar señales de finalización a los workers
    for _ in range(num_workers):
        await task_queue.put(None)
    
    # Esperar a que los workers terminen
    await asyncio.gather(*workers)
    
    # Esperar a que el recolector termine
    results = await collector
    
    # Mostrar resultados ordenados
    print("\nResultados finales:")
    for task_id in sorted(results.keys()):
        print(f"Tarea {task_id} -> {results[task_id]}")
    
    return results

if __name__ == "__main__":
    start = time.time()
    results = asyncio.run(fan_out_fan_in_example())
    elapsed = time.time() - start
    print(f"\nTiempo total: {elapsed:.2f} segundos")
```

### Patrón Actor

```python
import asyncio
import random
from enum import Enum, auto

# Tipos de mensajes
class MessageType(Enum):
    COMPUTE = auto()
    GET_STATE = auto()
    SHUTDOWN = auto()

# Mensaje
class Message:
     def __init__(self, type, data=None, sender=None, reply_to=None):
        self.type = type
        self.data = data
        self.sender = sender
        self.reply_to = reply_to

## Ejercicios Prácticos

1. **Implementar un sistema de caché concurrente**:
   - Crear una caché que permita múltiples lecturas concurrentes pero bloquee durante escrituras.
   - Implementar expiración de entradas.
   - Añadir estadísticas de hit/miss.

2. **Crear un servidor de chat con salas**:
   - Implementar un servidor que maneje múltiples clientes concurrentemente.
   - Permitir que los clientes se unan a diferentes salas de chat.
   - Implementar broadcast de mensajes dentro de cada sala.

3. **Desarrollar un crawler web concurrente**:
   - Implementar un crawler que visite URLs concurrentemente.
   - Limitar el número máximo de tareas concurrentes.
   - Implementar detección de ciclos para evitar visitar la misma URL múltiples veces.
   - Añadir timeouts para evitar bloqueos en sitios lentos.

4. **Implementar un sistema de procesamiento de datos en pipeline**:
   - Crear varias etapas de procesamiento (lectura, transformación, filtrado, agregación, escritura).
   - Conectar las etapas mediante colas.
   - Implementar manejo de errores a lo largo del pipeline.
   - Añadir capacidad para cancelar todo el pipeline.

5. **Crear un pool de conexiones a base de datos**:
   - Implementar un pool que gestione un número limitado de conexiones.
   - Permitir que múltiples hilos o corrutinas soliciten y liberen conexiones.
   - Implementar timeouts para solicitudes de conexión.
   - Añadir health checks para las conexiones.

## Conclusiones

La concurrencia avanzada en Python ofrece herramientas poderosas para resolver problemas complejos de manera eficiente. Cada modelo de concurrencia (threading, multiprocessing, asyncio) tiene sus propias ventajas y casos de uso ideales.

Al elegir un modelo de concurrencia, es importante considerar:

1. **Naturaleza de la tarea**: I/O bound vs CPU bound.
2. **Requisitos de estado compartido**: Necesidad de compartir datos entre tareas.
3. **Escala de concurrencia**: Número de tareas concurrentes necesarias.
4. **Complejidad de coordinación**: Patrones de comunicación y sincronización requeridos.

Los patrones de diseño concurrentes como Pipeline, Fan-out/Fan-in, y Actor proporcionan soluciones estructuradas para problemas comunes de concurrencia.

Finalmente, siempre es importante seguir las mejores prácticas para evitar problemas como race conditions, deadlocks, y fugas de recursos, que son especialmente críticos en aplicaciones concurrentes.

## Referencias

1. Python Documentation: Threading - https://docs.python.org/3/library/threading.html
2. Python Documentation: Multiprocessing - https://docs.python.org/3/library/multiprocessing.html
3. Python Documentation: Asyncio - https://docs.python.org/3/library/asyncio.html
4. Luciano Ramalho. (2015). Fluent Python. O'Reilly Media.
5. Caleb Hattingh. (2020). Using Asyncio in Python. O'Reilly Media.
6. David Beazley. (2019). Python Concurrency with asyncio. Manning Publications.
7. Brett Slatkin. (2019). Effective Python: 90 Specific Ways to Write Better Python. Addison-Wesley Professional. type, data=None, sender=None, reply_to=None):
        self.type = type
        self.data = data
        self.sender = sender
        self.reply_to = reply_to

# Actor base
class Actor:
    def __init__(self, name):
        self.name = name
        self.mailbox = asyncio.Queue()
        self.running = False
    
    async def start(self):
        self.running = True
        await self.run()
    
    async def stop(self):
        self.running = False
        # Enviar mensaje de apagado
        await self.send(Message(MessageType.SHUTDOWN))
    
    async def send(self, message):
        await self.mailbox.put(message)
    
    async def run(self):
        while self.running:
            # Obtener mensaje
            message = await self.mailbox.get()
            
            # Procesar mensaje
            if message.type == MessageType.SHUTDOWN:
                print(f"Actor {self.name} apagándose")
                self.running = False
                break
            
            # Procesar otros tipos de mensajes
            await self.process_message(message)
            
            # Marcar mensaje como procesado
            self.mailbox.task_done()
    
    async def process_message(self, message):
        # Implementado por subclases
        pass

# Actor de cálculo
class CalculatorActor(Actor):
    def __init__(self, name):
        super().__init__(name)
        self.state = 0
    
    async def process_message(self, message):
        if message.type == MessageType.COMPUTE:
            # Simular cálculo
            operation, value = message.data
            await asyncio.sleep(random.uniform(0.1, 0.5))  # Simular trabajo
            
            if operation == "add":
                self.state += value
            elif operation == "multiply":
                self.state *= value
            elif operation == "subtract":
                self.state -= value
            elif operation == "divide":
                if value != 0:
                    self.state /= value
            
            print(f"Actor {self.name} realizó {operation} {value}, nuevo estado: {self.state}")
            
            # Responder si se solicita
            if message.reply_to:
                await message.reply_to.send(
                    Message(MessageType.COMPUTE, self.state, sender=self)
                )
        
        elif message.type == MessageType.GET_STATE:
            print(f"Actor {self.name} enviando estado: {self.state}")
            # Responder con el estado actual
            if message.reply_to:
                await message.reply_to.send(
                    Message(MessageType.GET_STATE, self.state, sender=self)
                )

# Actor supervisor
class SupervisorActor(Actor):
    def __init__(self, name):
        super().__init__(name)
        self.workers = {}
        self.results = {}
    
    async def create_worker(self, worker_name):
        worker = CalculatorActor(worker_name)
        self.workers[worker_name] = worker
        asyncio.create_task(worker.start())
        print(f"Supervisor {self.name} creó worker {worker_name}")
        return worker
    
    async def process_message(self, message):
        if message.type == MessageType.COMPUTE:
            # Recibir resultado de un worker
            worker_name = message.sender.name
            self.results[worker_name] = message.data
            print(f"Supervisor {self.name} recibió resultado de {worker_name}: {message.data}")
        
        elif message.type == MessageType.GET_STATE:
            # Solicitar estado de todos los workers
            for worker_name, worker in self.workers.items():
                await worker.send(Message(
                    MessageType.GET_STATE,
                    reply_to=self
                ))

async def actor_example():
    # Crear supervisor
    supervisor = SupervisorActor("main-supervisor")
    supervisor_task = asyncio.create_task(supervisor.start())
    
    # Crear workers
    workers = []
    for i in range(3):
        worker = await supervisor.create_worker(f"calculator-{i}")
        workers.append(worker)
    
    # Enviar mensajes a los workers
    operations = [("add", 5), ("multiply", 2), ("subtract", 3), ("divide", 2)]
    
    for worker in workers:
        for operation, value in operations:
            await worker.send(Message(
                MessageType.COMPUTE,
                (operation, value),
                reply_to=supervisor
            ))
    
    # Esperar un poco para que se procesen los mensajes
    await asyncio.sleep(3)
    
    # Solicitar estado de todos los workers
    await supervisor.send(Message(MessageType.GET_STATE))
    
    # Esperar un poco más
    await asyncio.sleep(1)
    
    # Detener todos los actores
    for worker in workers:
        await worker.stop()
    await supervisor.stop()
    
    # Esperar a que el supervisor termine
    await supervisor_task
    
    print("\nResultados finales:")
    for worker_name, result in supervisor.results.items():
        print(f"{worker_name}: {result}")

if __name__ == "__main__":
    asyncio.run(actor_example())
```

## Mejores Prácticas

### Elegir el Modelo Correcto

- **Threading**: Ideal para tareas I/O bound con estado compartido.
- **Multiprocessing**: Ideal para tareas CPU bound sin necesidad de compartir mucho estado.
- **Asyncio**: Ideal para tareas I/O bound con alta concurrencia.

### Evitar Race Conditions

```python
import threading
import time

# Mal: Race condition
class UnsafeCounter:
    def __init__(self):
        self.count = 0
    
    def increment(self):
        current = self.count
        time.sleep(0.001)  # Simular trabajo
        self.count = current + 1

# Bien: Usando Lock
class SafeCounter:
    def __init__(self):
        self.count = 0
        self.lock = threading.Lock()
    
    def increment(self):
        with self.lock:
            current = self.count
            time.sleep(0.001)  # Simular trabajo
            self.count = current + 1

def test_counter(counter_class, num_threads, increments_per_thread):
    counter = counter_class()
    threads = []
    
    for _ in range(num_threads):
        thread = threading.Thread(
            target=lambda: [counter.increment() for _ in range(increments_per_thread)]
        )
        threads.append(thread)
    
    # Iniciar hilos
    for thread in threads:
        thread.start()
    
    # Esperar a que terminen
    for thread in threads:
        thread.join()
    
    expected = num_threads * increments_per_thread
    actual = counter.count
    print(f"{counter_class.__name__}: Esperado {expected}, Actual {actual}, {'Correcto' if expected == actual else 'Incorrecto'}")

if __name__ == "__main__":
    test_counter(UnsafeCounter, 10, 100)
    test_counter(SafeCounter, 10, 100)
```

### Evitar Deadlocks

```python
import threading
import time

# Mal: Potencial deadlock
class DeadlockRisk:
    def __init__(self):
        self.lock_a = threading.Lock()
        self.lock_b = threading.Lock()
    
    def method_a(self):
        with self.lock_a:
            print("Método A adquirió lock_a")
            time.sleep(0.1)  # Simular trabajo
            with self.lock_b:
                print("Método A adquirió lock_b")
                # Hacer algo con ambos recursos
    
    def method_b(self):
        with self.lock_b:
            print("Método B adquirió lock_b")
            time.sleep(0.1)  # Simular trabajo
            with self.lock_a:
                print("Método B adquirió lock_a")
                # Hacer algo con ambos recursos

# Bien: Evitar deadlock con orden consistente
class DeadlockSafe:
    def __init__(self):
        self.lock_a = threading.Lock()
        self.lock_b = threading.Lock()
    
    def method_a(self):
        with self.lock_a:
            print("Método A adquirió lock_a")
            time.sleep(0.1)  # Simular trabajo
            with self.lock_b:
                print("Método A adquirió lock_b")
                # Hacer algo con ambos recursos
    
    def method_b(self):
        with self.lock_a:  # Mismo orden que method_a
            print("Método B adquirió lock_a")
            time.sleep(0.1)  # Simular trabajo
            with self.lock_b:
                print("Método B adquirió lock_b")
                # Hacer algo con ambos recursos

def test_deadlock(cls):
    instance = cls()
    
    thread_a = threading.Thread(target=instance.method_a)
    thread_b = threading.Thread(target=instance.method_b)
    
    # Iniciar hilos
    thread_a.start()
    thread_b.start()
    
    # Esperar con timeout
    thread_a.join(timeout=2)
    thread_b.join(timeout=2)
    
    # Verificar si los hilos siguen vivos (posible deadlock)
    if thread_a.is_alive() or thread_b.is_alive():
        print(f"{cls.__name__}: Posible deadlock detectado")
    else:
        print(f"{cls.__name__}: No se detectó deadlock")

if __name__ == "__main__":
    print("Probando clase con riesgo de deadlock:")
    test_deadlock(DeadlockRisk)
    
    print("\nProbando clase segura contra deadlock:")
    test_deadlock(DeadlockSafe)
```

### Manejo de Recursos

```python
import threading
import asyncio

# Patrón de manejo de recursos con contexto
class ThreadSafeResource:
    def __init__(self, name):
        self.name = name
        self.lock = threading.Lock()
        self.in_use = False
    
    def __enter__(self):
        self.lock.acquire()
        self.in_use = True
        print(f"Recurso {self.name} adquirido")
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.in_use = False
        self.lock.release()
        print(f"Recurso {self.name} liberado")
        # Manejar excepciones si es necesario
        if exc_type is not None:
            print(f"Excepción manejada: {exc_type.__name__}: {exc_val}")
            # Devolver True para suprimir la excepción o False para propagarla
            return False

# Patrón de manejo de recursos asíncronos
class AsyncResource:
    def __init__(self, name):
        self.name = name
        self.lock = asyncio.Lock()
        self.in_use = False
    
    async def __aenter__(self):
        await self.lock.acquire()
        self.in_use = True
        print(f"Recurso asíncrono {self.name} adquirido")
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        self.in_use = False
        self.lock.release()
        print(f"Recurso asíncrono {self.name} liberado")
        # Manejar excepciones si es necesario
        if exc_type is not None:
            print(f"Excepción manejada: {exc_type.__name__}: {exc_val}")
            # Devolver True para suprimir la excepción o False para propagarla
            return False

# Ejemplo de uso
def use_thread_safe_resource():
    resource = ThreadSafeResource("database")
    
    def worker(worker_id):
        try:
            with resource:
                print(f"Worker {worker_id} usando el recurso")
                time.sleep(0.5)  # Simular trabajo
                if worker_id % 3 == 0:
                    # Simular error ocasional
                    raise ValueError(f"Error en worker {worker_id}")
        except ValueError as e:
            print(f"Error capturado fuera del contexto: {e}")
    
    threads = []
    for i in range(5):
        thread = threading.Thread(target=worker, args=(i,))
        threads.append(thread)
        thread.start()
    
    for thread in threads:
        thread.join()

async def use_async_resource():
    resource = AsyncResource("api")
    
    async def worker(worker_id):
        try:
            async with resource:
                print(f"Worker asíncrono {worker_id} usando el recurso")
                await asyncio.sleep(0.5)  # Simular trabajo asíncrono
                if worker_id % 3 == 0:
                    # Simular error ocasional
                    raise ValueError(f"Error en worker asíncrono {worker_id}")
        except ValueError as e:
            print(f"Error capturado fuera del contexto asíncrono: {e}")
    
    tasks = [worker(i) for i in range(5)]
    await asyncio.gather(*tasks)

if __name__ == "__main__":
    print("=== Ejemplo de recurso thread-safe ===")
    use_thread_safe_resource()
    
    print("\n=== Ejemplo de recurso asíncrono ===")
    asyncio.run(use_async_resource())
    def __init__(self,