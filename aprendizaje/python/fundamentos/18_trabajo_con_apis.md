# 18. Consumo de APIs web con `requests`

## Introducción

Una de las tareas más comunes en la programación moderna es interactuar con APIs (Interfaces de Programación de Aplicaciones) para obtener o enviar datos a través de la web.

## Temas a cubrir:

- ¿Qué es una API y cómo funciona HTTP? (GET, POST, PUT, DELETE).
- La librería `requests`, el estándar de facto en Python para hacer peticiones HTTP.
- Realizar una petición `GET` para obtener datos.
- Trabajar con la respuesta: código de estado, cabeceras y contenido (JSON).
- Enviar datos con peticiones `POST`.
- Manejo de errores y timeouts.

## Mejores Prácticas y Recursos

- Maneja siempre posibles errores de red y códigos de estado HTTP inesperados.
- No expongas claves de API (API keys) directamente en el código; usa variables de entorno.
- **Librería Esencial:** [Documentación de `requests`](https://requests.readthedocs.io/en/latest/)
- **Tutorial Completo:** [Guía sobre APIs y `requests` de Real Python](https://realpython.com/python-requests/)
