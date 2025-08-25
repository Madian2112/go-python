# 05. Programación Asíncrona con `asyncio`

## Introducción

La programación asíncrona es un paradigma de concurrencia que permite a un solo hilo manejar múltiples operaciones de I/O (entrada/salida) de manera eficiente. En Python, `asyncio` es la librería estándar para escribir código concurrente usando la sintaxis `async`/`await`.

## Temas a cubrir:

- Concurrencia vs. Paralelismo.
- ¿Qué es I/O-bound vs. CPU-bound? ¿Cuándo usar `asyncio`?
- La sintaxis `async def` para corutinas.
- La palabra clave `await` para pausar la ejecución.
- El "event loop" (bucle de eventos) de `asyncio`.
- Ejecución de tareas concurrentes con `asyncio.gather()`.
- Ecosistema: librerías como `aiohttp` (cliente/servidor web) y `httpx` (cliente HTTP).

## Mejores Prácticas y Recursos

- `asyncio` es ideal para aplicaciones con muchas operaciones de red (APIs, websockets, bases de datos) y no para tareas que consumen mucho CPU.
- Ten cuidado de no mezclar código bloqueante síncrono con código asíncrono, ya que puede detener todo el bucle de eventos.
- **Documentación Oficial:** [Librería `asyncio`](https://docs.python.org/es/3/library/asyncio.html)
- **Guía de Referencia:** [Introducción a `async`/`await` en Real Python](https://realpython.com/async-io-python/)
