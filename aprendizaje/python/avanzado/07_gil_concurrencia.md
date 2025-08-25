# 07. El Global Interpreter Lock (GIL) y Estrategias de Concurrencia

## Introducción

El GIL es una de las características más discutidas y a menudo malentendidas de CPython (la implementación estándar de Python). Entenderlo es clave para diseñar aplicaciones concurrentes y paralelas de alto rendimiento.

## Temas a cubrir:

- ¿Qué es el Global Interpreter Lock (GIL)?
- ¿Por qué existe? (Gestión de memoria y simplificación de la implementación).
- Impacto del GIL en el `threading`: solo un hilo puede ejecutar bytecode de Python a la vez.
- **Cuándo usar `threading`**: Para tareas I/O-bound (operaciones de red, disco), donde los hilos pueden liberar el GIL mientras esperan.
- **Cuándo usar `multiprocessing`**: Para tareas CPU-bound (cálculos intensivos), ya que cada proceso tiene su propio intérprete de Python y su propio GIL.
- Comparativa de `threading` vs. `multiprocessing` vs. `asyncio`.

## Mejores Prácticas y Recursos

- "Usa `asyncio` para muchas conexiones de red. Usa `threading` para algunas operaciones de I/O. Usa `multiprocessing` para cálculos intensivos en CPU."
- El GIL no es un problema para la mayoría de las aplicaciones, pero es crucial conocer sus implicaciones para sistemas de alto rendimiento.
- **Explicación Detallada:** [¿Qué es el GIL? en Real Python](https://realpython.com/python-gil/)
- **Video Explicativo:** [Charla de David Beazley sobre el GIL (un clásico)](https://www.youtube.com/watch?v=Obt-vMVdM8s)
