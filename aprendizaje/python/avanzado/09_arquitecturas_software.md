# 09. Arquitecturas de Software en Proyectos Python

## Introducción

Elegir la arquitectura correcta es una de las decisiones más importantes al iniciar un proyecto. Define cómo se organizarán los componentes, cómo interactuarán entre sí y cómo evolucionará el sistema a lo largo del tiempo. No existe una "arquitectura perfecta", sino arquitecturas adecuadas para diferentes contextos.

En este documento, exploraremos varias arquitecturas comunes en el ecosistema de Python, con ejemplos de su estructura de directorios y cuándo es más apropiado usar cada una.

---

## 1. Arquitectura Monolítica (Monolith)

Es el enfoque tradicional donde toda la aplicación se construye como una única unidad cohesiva y desplegable.

- **Explicación:** Todos los componentes (UI, lógica de negocio, acceso a datos) están fuertemente acoplados y se ejecutan en el mismo proceso. Frameworks como Django y Flask promueven este enfoque por defecto.
- **Cuándo usarla:** Ideal para proyectos pequeños y medianos, MVPs (Minimum Viable Products), prototipos o cuando el equipo es pequeño y el dominio del problema no es excesivamente complejo.
- **Estructura de Directorios (Ejemplo con Flask):**
  ```
  /mi_proyecto
  ├── /app                # Contenedor principal de la aplicación
  │   ├── __init__.py     # Inicializa la aplicación Flask (app factory)
  │   ├── /static         # Archivos estáticos (CSS, JS, imágenes)
  │   ├── /templates      # Plantillas HTML (Jinja2)
  │   ├── /models         # Módulos de modelos de datos (ej. SQLAlchemy)
  │   ├── /views          # Lógica de las rutas (endpoints)
  │   └── /services       # Lógica de negocio separada de las vistas
  ├── /tests              # Pruebas unitarias e de integración
  ├── config.py           # Configuraciones (claves secretas, BBDD)
  ├── requirements.txt    # Dependencias del proyecto
  └── run.py              # Punto de entrada para ejecutar la aplicación
  ```
- **Recursos:**
  - [The Monolith First Approach](https://www.mongodb.com/developer/products/mongodb/monolith-first-approach-pattern/)
  - [Tutorial de la mega-aplicación de Flask por Miguel Grinberg](https://blog.miguelgrinberg.com/post/the-flask-mega-tutorial-part-i-hello-world) (Un excelente ejemplo de un monolito bien estructurado)

---

## 2. Arquitectura de Microservicios

La aplicación se descompone en un conjunto de servicios pequeños, independientes y débilmente acoplados. Cada servicio es responsable de una capacidad de negocio específica.

- **Explicación:** Cada microservicio tiene su propia base de datos, se desarrolla y despliega de forma independiente. Se comunican entre sí a través de la red, comúnmente mediante APIs REST o mensajería asíncrona.
- **Cuándo usarla:** Para aplicaciones grandes y complejas, cuando se requiere alta escalabilidad, en equipos de desarrollo grandes que pueden trabajar en paralelo, o cuando se quiere usar diferentes tecnologías para diferentes servicios.
- **Estructura de Directorios (Ejemplo Monorepo):**
  ```
  /mi_plataforma
  ├── /services                   # Directorio para todos los microservicios
  │   ├── /auth-service           # Servicio de autenticación
  │   │   ├── /app
  │   │   ├── Dockerfile
  │   │   └── requirements.txt
  │   ├── /product-service        # Servicio de productos
  │   │   ├── /app
  │   │   ├── Dockerfile
  │   │   └── requirements.txt
  │   └── /order-service          # Servicio de pedidos
  │       ├── /app
  │       ├── Dockerfile
  │       └── requirements.txt
  ├── /shared-libs                # Librerías compartidas entre servicios (opcional)
  └── docker-compose.yml          # Orquestación de los servicios para desarrollo local
  ```
- **Recursos:**
  - [Microservices.io por Chris Richardson](https://microservices.io/) (Referencia fundamental)
  - [Patrones de Diseño para Microservicios](https://martinfowler.com/articles/microservices.html) (Artículo de Martin Fowler)

---

## 3. Arquitectura Limpia (Clean Architecture)

Es un patrón de diseño de software que separa el código en capas concéntricas, con un fuerte énfasis en la independencia del framework, la UI y la base de datos.

- **Explicación:** La regla principal es que las dependencias solo pueden apuntar hacia adentro. La capa más interna (`Entities`) contiene la lógica de negocio más pura y no sabe nada de las capas externas.
  - **Capas (de adentro hacia afuera):** Entities -> Use Cases -> Interface Adapters (Controllers, Presenters) -> Frameworks & Drivers (Web, DB).
- **Cuándo usarla:** En proyectos complejos de larga duración donde la mantenibilidad y la facilidad para realizar pruebas son cruciales. Cuando se quiere evitar estar atado a un framework o tecnología específica.
- **Estructura de Directorios (Ejemplo conceptual):**
  ```
  /mi_proyecto
  ├── /src
  │   ├── /domain               # Capa de "Entities" y "Use Cases"
  │   │   ├── /entities         # Modelos de dominio puros (clases Python)
  │   │   └── /use_cases        # Lógica de negocio, orquesta las entidades
  │   ├── /adapters             # Capa de "Interface Adapters"
  │   │   ├── /http             # Controladores de API, Presenters
  │   │   ├── /repositories     # Implementaciones de repositorios (ej. para BBDD)
  │   │   └── /event_bus        # Adaptadores para sistemas de mensajería
  │   └── /infrastructure       # Capa "Frameworks & Drivers"
  │       ├── /database         # Configuración de BBDD, migraciones (SQLAlchemy)
  │       ├── /web              # Framework web (Flask, FastAPI), enrutamiento
  │       └── /dependencies.py  # Inyección de dependencias
  ├── /tests
  └── ...
  ```
- **Recursos:**
  - [Post original de "The Clean Architecture" por Robert C. Martin (Uncle Bob)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
  - [Artículo práctico sobre Arquitectura Limpia en Python](https://medium.com/@hugo.nobre.leite/clean-architecture-in-python-a-practical-approach-to-organize-your-code-123df9434423)

---

## 4. Arquitectura Hexagonal (Puertos y Adaptadores)

Similar en filosofía a la Arquitectura Limpia, busca aislar la lógica de negocio de los detalles externos (infraestructura, UI).

- **Explicación:** El "hexágono" representa el núcleo de la aplicación (dominio). La comunicación con el exterior se realiza a través de "puertos" (que son interfaces definidas por el dominio) y "adaptadores" (que son las implementaciones concretas de esos puertos).
  - **Adaptadores "Driving"**: Inician la interacción (ej. un controlador de API).
  - **Adaptadores "Driven"**: Son invocados por el núcleo (ej. un repositorio de base de datos).
- **Cuándo usarla:** Mismos casos que la Arquitectura Limpia. A menudo, la elección entre ambas es una cuestión de preferencia terminológica y de estructuración.
- **Estructura de Directorios:**
  ```
  /mi_proyecto
  ├── /src
  │   ├── /domain                 # El núcleo o "hexágono"
  │   │   ├── /models             # Modelos de dominio
  │   │   ├── /services           # Lógica de negocio
  │   │   └── /ports              # Interfaces (ej. "ProductsRepositoryPort")
  │   └── /adapters               # Implementaciones de los puertos
  │       ├── /http               # Adaptador "driving" (API, controladores)
  │       ├── /db                 # Adaptador "driven" (implementación del repositorio)
  │       └── /cli                # Otro adaptador "driving" (comandos de consola)
  ├── /tests
  └── ...
  ```
- **Recursos:**
  - [Artículo original de Alistair Cockburn sobre Arquitectura Hexagonal](https://alistair.cockburn.us/hexagonal-architecture/)
  - [Implementación de Arquitectura Hexagonal en Python (video)](https://www.youtube.com/watch?v=n4Gz1gH5I94)
