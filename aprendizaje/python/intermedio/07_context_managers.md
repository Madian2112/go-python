# 07. Manejo de Contextos (Context Managers) y la declaración `with`

## Introducción

El manejo de recursos (como archivos, conexiones de red o a bases de datos) es una tarea común que requiere una configuración inicial y una limpieza final. Los manejadores de contexto automatizan este proceso de forma elegante y segura.

## Temas a cubrir:

- La declaración `with` y su propósito.
- ¿Por qué `with` es más seguro que `try...finally`?
- Implementación de un manejador de contexto usando una clase con los métodos `__enter__` y `__exit__`.
- Creación de un manejador de contexto de forma más sencilla usando el decorador `@contextmanager` de la librería `contextlib`.

## Mejores Prácticas y Recursos

- Siempre que una librería te ofrezca un objeto que pueda ser usado con `with`, úsalo. Es el caso de `open()`, `threading.Lock()`, `requests.Session()`, etc.
- Esto garantiza que los recursos se liberen correctamente, incluso si ocurren errores.
- **Guía de Referencia:** [Manejadores de Contexto en Real Python](https://realpython.com/python-with-statement/)
- **Documentación Oficial:** [La librería `contextlib`](https://docs.python.org/es/3/library/contextlib.html)
