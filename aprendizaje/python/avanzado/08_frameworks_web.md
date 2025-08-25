# 08. Panorama de Frameworks Web

## Introducción

Python tiene un ecosistema de desarrollo web extremadamente rico y maduro. Elegir el framework adecuado depende de las necesidades del proyecto, la escala y la experiencia del equipo.

## Temas a cubrir:

### 1. Django: El framework "con baterías incluidas"
- **Filosofía**: Full-stack, promueve el desarrollo rápido y el diseño pragmático.
- **Componentes clave**: ORM, sistema de administración, autenticación, sistema de plantillas.
- **Ideal para**: Proyectos grandes y complejos, portales de noticias, e-commerce, CMS.

### 2. Flask: El micro-framework
- **Filosofía**: Minimalista, flexible y extensible. Proporciona lo esencial y deja el resto a elección del desarrollador.
- **Componentes clave**: Enrutamiento y motor de plantillas Jinja2.
- **Ideal para**: APIs, microservicios, aplicaciones pequeñas y medianas, prototipos.

### 3. FastAPI: El framework moderno de alto rendimiento
- **Filosofía**: Basado en estándares modernos (OpenAPI, JSON Schema), `asyncio` y type hints.
- **Componentes clave**: Validación de datos con Pydantic, documentación automática de API, soporte asíncrono.
- **Ideal para**: APIs de alto rendimiento, aplicaciones que se benefician de `asyncio` y tipado estático.

## Mejores Prácticas y Recursos

- **Django**: Sigue sus convenciones, te ahorrará mucho trabajo. Es excelente para proyectos donde necesitas una solución completa y robusta rápidamente. [Sitio oficial de Django](https://www.djangoproject.com/)
- **Flask**: Empieza pequeño y añade extensiones conforme las necesites. Ideal para aprender los fundamentos del desarrollo web. [Sitio oficial de Flask](https://flask.palletsprojects.com/)
- **FastAPI**: Úsalo si tu prioridad es el rendimiento de I/O y quieres aprovechar las ventajas del tipado estático y `asyncio`. [Sitio oficial de FastAPI](https://fastapi.tiangolo.com/)
