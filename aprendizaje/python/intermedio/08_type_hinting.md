# 08. Tipado Estático Gradual (Type Hinting)

## Introducción

Desde Python 3.5, es posible añadir "pistas de tipo" (type hints) al código. Esto no cambia cómo se ejecuta el código (Python sigue siendo de tipado dinámico), pero proporciona enormes beneficios para el análisis estático, la legibilidad y la robustez del código.

## Temas a cubrir:

- Sintaxis de los type hints para variables, parámetros de funciones y valores de retorno.
- Uso de tipos complejos del módulo `typing`: `List`, `Dict`, `Tuple`, `Optional`, `Union`.
- ¿Qué es el análisis estático de tipos?
- Uso de `MyPy` para verificar la consistencia de los tipos en tu código antes de ejecutarlo.
- Beneficios: autocompletado mejorado en el editor, código auto-documentado, y detección temprana de errores.

## Mejores Prácticas y Recursos

- Adopta el tipado estático gradualmente en tus proyectos. No es necesario tipar todo de una vez.
- Es especialmente útil en las interfaces públicas de tus módulos y en lógica de negocio compleja.
- Integra `MyPy` en tu flujo de CI/CD para asegurar la calidad del código.
- **Recurso Principal:** [Documentación del módulo `typing`](https://docs.python.org/es/3/library/typing.html)
- **Herramienta Esencial:** [Sitio oficial de `MyPy`](http://mypy-lang.org/)
- **Tutorial Detallado:** [Guía de Type Hinting en Real Python](https://realpython.com/python-type-checking/)
