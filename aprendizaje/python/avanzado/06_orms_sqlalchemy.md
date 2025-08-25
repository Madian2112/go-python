# 06. ORMs y Bases de Datos con SQLAlchemy

## Introducción

Interactuar con bases de datos SQL es una tarea fundamental. Un ORM (Object-Relational Mapper) permite trabajar con la base de datos usando objetos de Python en lugar de escribir SQL crudo, lo que aumenta la productividad y la seguridad.

## Temas a cubrir:

- ¿Qué es un ORM y qué problemas resuelve?
- **SQLAlchemy Core**: El constructor de expresiones SQL de Python. Permite escribir SQL de forma "Pythonica" y segura.
- **SQLAlchemy ORM**: El ORM completo.
  - Definición de modelos de datos usando clases.
  - Creación de una sesión y conexión a la base de datos.
  - CRUD: Crear, Leer, Actualizar y Borrar registros usando objetos.
  - Manejo de relaciones (uno a muchos, muchos a muchos).
- Migraciones de base de datos con `Alembic`.

## Mejores Prácticas y Recursos

- SQLAlchemy es el estándar de oro en el ecosistema Python. Aprenderlo te servirá en innumerables proyectos.
- Usa Alembic desde el inicio del proyecto para gestionar la evolución del esquema de tu base de datos de forma controlada.
- **Sitio Oficial:** [Documentación de SQLAlchemy](https://www.sqlalchemy.org/)
- **Herramienta de Migraciones:** [Alembic](https://alembic.sqlalchemy.org/)
- **Tutorial Detallado:** [Tutorial de SQLAlchemy en Real Python](https://realpython.com/python-sqlalchemy-database-tutorial/)
