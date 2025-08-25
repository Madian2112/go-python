# 06. Generadores y la palabra clave `yield`

## Introducción

Los generadores ofrecen una forma de crear iteradores de una manera sencilla y eficiente en memoria. Son fundamentales para trabajar con secuencias de datos muy grandes o infinitas.

## Temas a cubrir:

- ¿Qué es un iterador y un iterable?
- La palabra clave `yield` vs. `return`.
- Cómo funciona una función generadora.
- Expresiones generadoras (similares a list comprehensions pero con paréntesis).
- Ventajas de los generadores: eficiencia de memoria y evaluación perezosa (lazy evaluation).
- Pipelines de datos con generadores.

## Mejores Prácticas y Recursos

- Utiliza generadores siempre que trabajes con conjuntos de datos que no necesites cargar completamente en memoria.
- Las expresiones generadoras son preferibles para casos simples por su sintaxis concisa.
- **Referencia Principal:** [Documentación oficial sobre Generadores](https://docs.python.org/es/3/howto/functional.html#generators)
- **Guía Práctica:** [Introducción a Generadores en Real Python](https://realpython.com/introduction-to-python-generators/)
