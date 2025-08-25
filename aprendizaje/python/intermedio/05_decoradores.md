# 05. Decoradores

## Introducción

Los decoradores son una de las características más poderosas y utilizadas en Python. Son una forma de extender o modificar el comportamiento de funciones o métodos sin alterar su código fuente.

## Temas a cubrir:

- Funciones como ciudadanos de primera clase (pasarlas como argumentos, retornarlas).
- Sintaxis de un decorador (`@my_decorator`).
- Creación de un decorador simple.
- Uso de `functools.wraps` para preservar los metadatos de la función original.
- Decoradores con argumentos.
- Casos de uso comunes: logging, caching, control de acceso, etc.

## Mejores Prácticas y Recursos

- Los decoradores deben ser claros en su propósito. Si un decorador se vuelve muy complejo, puede ser una señal de que la lógica debería estar en otro lugar.
- Usa siempre `@functools.wraps` para evitar efectos secundarios inesperados al inspeccionar funciones decoradas.
- **Guía de Referencia:** [Primer sobre Decoradores en Real Python](https://realpython.com/primer-on-python-decorators/)
- **Documentación PEP 318:** [La propuesta original de los decoradores](https://peps.python.org/pep-0318/)
