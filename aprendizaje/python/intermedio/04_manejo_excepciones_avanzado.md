# Manejo Avanzado de Excepciones en Python

## Introducción

El manejo de excepciones es una parte fundamental de la programación en Python. Permite que los programas respondan de manera elegante a situaciones inesperadas, mejorando la robustez y la experiencia del usuario. En este módulo, exploraremos técnicas avanzadas de manejo de excepciones, incluyendo la creación de excepciones personalizadas, patrones de manejo de excepciones, y mejores prácticas para escribir código más resiliente y mantenible.

## Repaso de Conceptos Básicos

### Estructura try-except-else-finally

```python
try:
    # Código que puede generar una excepción
    resultado = 10 / 0
except ZeroDivisionError as e:
    # Se ejecuta si ocurre una ZeroDivisionError
    print(f"Error: {e}")
except (TypeError, ValueError) as e:
    # Se ejecuta si ocurre TypeError o ValueError
    print(f"Error de tipo o valor: {e}")
else:
    # Se ejecuta si no ocurre ninguna excepción en el bloque try
    print(f"Resultado: {resultado}")
finally:
    # Se ejecuta siempre, haya o no excepción
    print("Limpieza final")
```

## Excepciones Personalizadas

### Creación de Excepciones Personalizadas

Las excepciones personalizadas permiten definir errores específicos de tu aplicación:

```python
class MiErrorPersonalizado(Exception):
    """Clase base para excepciones en este módulo."""
    pass

class ValorInvalidoError(MiErrorPersonalizado):
    """Se lanza cuando un valor no cumple con ciertos criterios."""
    def __init__(self, valor, mensaje="Valor no válido"):
        self.valor = valor
        self.mensaje = mensaje
        super().__init__(f"{mensaje}: {valor}")

class RecursoNoEncontradoError(MiErrorPersonalizado):
    """Se lanza cuando no se encuentra un recurso solicitado."""
    def __init__(self, recurso, mensaje="Recurso no encontrado"):
        self.recurso = recurso
        self.mensaje = mensaje
        super().__init__(f"{mensaje}: {recurso}")

# Uso
def procesar_valor(valor):
    if valor < 0:
        raise ValorInvalidoError(valor, "El valor no puede ser negativo")
    # Procesamiento normal...
    return valor * 2

try:
    resultado = procesar_valor(-5)
except ValorInvalidoError as e:
    print(f"Error: {e}")
    print(f"Valor problemático: {e.valor}")
```

### Jerarquía de Excepciones

Es recomendable crear una jerarquía de excepciones para tu aplicación:

```python
class AppError(Exception):
    """Clase base para todas las excepciones de la aplicación."""
    pass

class DatabaseError(AppError):
    """Errores relacionados con la base de datos."""
    pass

class ConnectionError(DatabaseError):
    """Error de conexión a la base de datos."""
    pass

class QueryError(DatabaseError):
    """Error en la ejecución de una consulta."""
    pass

class ValidationError(AppError):
    """Errores de validación de datos."""
    pass

class AuthenticationError(AppError):
    """Errores de autenticación."""
    pass

# Uso
try:
    # Código que puede generar diferentes tipos de errores
    pass
except ConnectionError as e:
    print(f"Error de conexión: {e}")
except QueryError as e:
    print(f"Error en la consulta: {e}")
except DatabaseError as e:
    print(f"Error de base de datos: {e}")
except ValidationError as e:
    print(f"Error de validación: {e}")
except AppError as e:
    print(f"Error de la aplicación: {e}")
```

### Excepciones con Contexto

Las excepciones pueden contener información contextual adicional:

```python
class ContextualError(Exception):
    """Excepción que proporciona contexto adicional sobre el error."""
    def __init__(self, mensaje, contexto=None):
        self.mensaje = mensaje
        self.contexto = contexto or {}
        super().__init__(mensaje)
    
    def __str__(self):
        if not self.contexto:
            return self.mensaje
        
        contexto_str = ", ".join(f"{k}={v}" for k, v in self.contexto.items())
        return f"{self.mensaje} [{contexto_str}]"

# Uso
def procesar_usuario(usuario_id):
    try:
        # Código que puede fallar
        if not isinstance(usuario_id, int):
            raise TypeError("usuario_id debe ser un entero")
        
        if usuario_id <= 0:
            raise ValueError("usuario_id debe ser positivo")
        
        # Más código...
    except (TypeError, ValueError) as e:
        # Capturar y enriquecer la excepción
        contexto = {
            "usuario_id": usuario_id,
            "operacion": "procesar_usuario",
            "timestamp": datetime.now().isoformat()
        }
        raise ContextualError(str(e), contexto) from e

try:
    procesar_usuario("abc")
except ContextualError as e:
    print(f"Error: {e}")
    print(f"Contexto: {e.contexto}")
    print(f"Causa original: {e.__cause__}")
```

## Patrones Avanzados de Manejo de Excepciones

### Encadenamiento de Excepciones (Exception Chaining)

Python permite encadenar excepciones para mantener el contexto de error:

```python
try:
    try:
        # Código que puede fallar
        1 / 0
    except ZeroDivisionError as e:
        # Capturar y lanzar una nueva excepción con contexto
        raise ValueError("No se pudo completar el cálculo") from e
except ValueError as e:
    print(f"Error: {e}")
    print(f"Causado por: {e.__cause__}")
```

### Supresión de Contexto

En algunos casos, puedes querer suprimir el encadenamiento de excepciones:

```python
try:
    # Código que puede fallar
    1 / 0
except ZeroDivisionError:
    # Lanzar una nueva excepción sin encadenar
    raise ValueError("No se pudo completar el cálculo") from None
```

### Manejo de Excepciones en Contextos

Los manejadores de contexto (`with`) son útiles para garantizar la limpieza de recursos:

```python
class ManejoArchivo:
    def __init__(self, nombre_archivo, modo):
        self.nombre_archivo = nombre_archivo
        self.modo = modo
        self.archivo = None
    
    def __enter__(self):
        try:
            self.archivo = open(self.nombre_archivo, self.modo)
            return self.archivo
        except IOError as e:
            raise ContextualError(f"No se pudo abrir {self.nombre_archivo}", 
                                 {"modo": self.modo}) from e
    
    def __exit__(self, tipo_exc, valor_exc, traceback_exc):
        if self.archivo:
            self.archivo.close()
        
        # Manejar excepciones específicas
        if tipo_exc is IOError:
            print(f"Error de IO al cerrar el archivo: {valor_exc}")
            return True  # Suprimir la excepción
        
        return False  # Propagar otras excepciones

# Uso
try:
    with ManejoArchivo("archivo.txt", "r") as f:
        contenido = f.read()
        # Procesar contenido...
except ContextualError as e:
    print(f"Error: {e}")
```

### Manejo de Múltiples Excepciones

```python
def validar_formulario(datos):
    errores = []
    
    try:
        if not datos.get("nombre"):
            errores.append("El nombre es obligatorio")
        
        edad = datos.get("edad")
        if edad is None:
            errores.append("La edad es obligatoria")
        elif not isinstance(edad, int):
            errores.append("La edad debe ser un número entero")
        elif edad < 0 or edad > 120:
            errores.append("La edad debe estar entre 0 y 120")
        
        email = datos.get("email")
        if not email:
            errores.append("El email es obligatorio")
        elif "@" not in email:
            errores.append("El email no es válido")
        
        if errores:
            raise ValidationError("Errores de validación", errores=errores)
        
        return True
    except Exception as e:
        if not isinstance(e, ValidationError):
            # Capturar excepciones inesperadas
            errores.append(f"Error inesperado: {str(e)}")
            raise ValidationError("Errores de validación", errores=errores) from e
        raise

class ValidationError(Exception):
    def __init__(self, mensaje, errores=None):
        self.mensaje = mensaje
        self.errores = errores or []
        super().__init__(mensaje)
    
    def __str__(self):
        if not self.errores:
            return self.mensaje
        
        errores_str = "\n - ".join(self.errores)
        return f"{self.mensaje}:\n - {errores_str}"

# Uso
try:
    validar_formulario({"nombre": "", "edad": "treinta", "email": "correo-invalido"})
except ValidationError as e:
    print(f"Error de validación: {e}")
```

### Excepciones en Funciones Asíncronas

```python
import asyncio

async def operacion_asincrona():
    try:
        await asyncio.sleep(1)
        raise ValueError("Error en operación asíncrona")
    except ValueError as e:
        print(f"Capturado dentro de la función asíncrona: {e}")
        raise  # Re-lanzar para propagar

async def main():
    try:
        await operacion_asincrona()
    except ValueError as e:
        print(f"Capturado en main: {e}")

# Ejecutar
asyncio.run(main())
```

## Técnicas Avanzadas

### Decoradores para Manejo de Excepciones

Los decoradores pueden simplificar el manejo de excepciones repetitivo:

```python
import functools
import logging

def manejar_excepciones(func=None, *, reintento=0, logger=None):
    if func is None:
        return lambda f: manejar_excepciones(f, reintento=reintento, logger=logger)
    
    logger = logger or logging.getLogger(func.__module__)
    
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        intentos = 0
        while True:
            try:
                return func(*args, **kwargs)
            except Exception as e:
                intentos += 1
                logger.error(f"Error en {func.__name__}: {e}")
                
                if intentos <= reintento:
                    logger.info(f"Reintentando {func.__name__} ({intentos}/{reintento})")
                    continue
                
                raise
    
    return wrapper

# Uso
@manejar_excepciones(reintento=3)
def operacion_con_reintentos():
    import random
    if random.random() < 0.8:  # 80% de probabilidad de error
        raise ConnectionError("Falló la conexión")
    return "Éxito"

try:
    resultado = operacion_con_reintentos()
    print(f"Resultado: {resultado}")
except ConnectionError as e:
    print(f"Error final después de reintentos: {e}")
```

### Context Manager para Transacciones

```python
class Transaccion:
    def __init__(self, db):
        self.db = db
        self.tx = None
    
    def __enter__(self):
        self.tx = self.db.begin()
        return self.tx
    
    def __exit__(self, tipo_exc, valor_exc, traceback_exc):
        if tipo_exc is None:
            # Sin excepciones, confirmar la transacción
            self.tx.commit()
            return True
        else:
            # Con excepciones, revertir la transacción
            self.tx.rollback()
            # Propagar la excepción
            return False

# Uso
class BaseDatos:
    def begin(self):
        print("Iniciando transacción")
        return self
    
    def commit(self):
        print("Confirmando transacción")
    
    def rollback(self):
        print("Revirtiendo transacción")
    
    def ejecutar(self, query):
        if "ERROR" in query:
            raise ValueError(f"Error en la consulta: {query}")
        print(f"Ejecutando: {query}")

db = BaseDatos()

# Transacción exitosa
try:
    with Transaccion(db) as tx:
        tx.ejecutar("INSERT INTO usuarios VALUES (1, 'Alice')")
        tx.ejecutar("UPDATE contadores SET valor = valor + 1")
    print("Transacción completada con éxito")
except ValueError as e:
    print(f"Error: {e}")

# Transacción con error
try:
    with Transaccion(db) as tx:
        tx.ejecutar("INSERT INTO usuarios VALUES (2, 'Bob')")
        tx.ejecutar("ERROR: consulta inválida")
    print("Este mensaje no se imprimirá si hay error")
except ValueError as e:
    print(f"Error: {e}")
```

### Captura Selectiva con Filtrado

```python
def es_error_temporal(e):
    """Determina si un error es temporal y se puede reintentar."""
    if isinstance(e, (ConnectionError, TimeoutError)):
        return True
    if isinstance(e, OSError) and e.errno in (10054, 10053):  # Códigos de error de red
        return True
    return False

def operacion_con_filtro():
    try:
        # Código que puede fallar
        raise ConnectionError("Conexión rechazada")
    except Exception as e:
        if es_error_temporal(e):
            print("Error temporal detectado, reintentando...")
            # Lógica de reintento
        else:
            print("Error permanente, propagando...")
            raise

operacion_con_filtro()
```

## Mejores Prácticas

### 1. Ser Específico con las Excepciones

```python
# Malo: captura demasiado amplia
try:
    # Código
except Exception as e:
    # Manejo genérico

# Bueno: captura específica
try:
    # Código
except ValueError as e:
    # Manejo específico para ValueError
except IOError as e:
    # Manejo específico para IOError
except Exception as e:
    # Manejo para otras excepciones inesperadas
    logging.error(f"Error inesperado: {e}", exc_info=True)
```

### 2. Proporcionar Información Contextual

```python
try:
    procesar_archivo("datos.csv")
except Exception as e:
    logging.error(f"Error al procesar datos.csv: {e}", exc_info=True)
    raise RuntimeError(f"Error al procesar el archivo datos.csv") from e
```

### 3. Usar finally para Limpieza

```python
recurso = None
try:
    recurso = abrir_recurso()
    # Usar el recurso
finally:
    if recurso:
        recurso.cerrar()
```

### 4. Evitar Pasar Silenciosamente las Excepciones

```python
# Malo: silenciar excepciones sin acción
try:
    operacion_importante()
except Exception:
    pass  # Silenciar el error

# Bueno: al menos registrar el error
try:
    operacion_importante()
except Exception as e:
    logging.error(f"Error en operación importante: {e}", exc_info=True)
    # Decidir si se propaga o no
```

### 5. Crear Jerarquías de Excepciones Significativas

```python
class ErrorAplicacion(Exception):
    """Base para todas las excepciones de la aplicación."""

class ErrorEntrada(ErrorAplicacion):
    """Error en los datos de entrada."""

class ErrorProcesamiento(ErrorAplicacion):
    """Error durante el procesamiento."""

class ErrorSalida(ErrorAplicacion):
    """Error al generar la salida."""
```

### 6. Documentar el Comportamiento de Excepciones

```python
def procesar_archivo(ruta):
    """Procesa un archivo de datos.
    
    Args:
        ruta: Ruta al archivo a procesar.
        
    Returns:
        dict: Datos procesados.
        
    Raises:
        FileNotFoundError: Si el archivo no existe.
        PermissionError: Si no hay permisos para leer el archivo.
        ValueError: Si el formato del archivo es inválido.
    """
    # Implementación...
```

### 7. Usar Context Managers para Recursos

```python
# En lugar de:
f = open("archivo.txt", "r")
try:
    datos = f.read()
finally:
    f.close()

# Usar:
with open("archivo.txt", "r") as f:
    datos = f.read()
```

### 8. Considerar el Nivel de Detalle Apropiado

```python
# Para usuarios finales
try:
    # Operación que puede fallar
except Exception as e:
    print("Lo sentimos, ocurrió un error. Por favor, inténtelo más tarde.")
    logging.error(f"Error detallado: {e}", exc_info=True)

# Para desarrolladores/depuración
try:
    # Operación que puede fallar
except Exception as e:
    print(f"Error: {e}")
    print(f"Tipo: {type(e).__name__}")
    import traceback
    traceback.print_exc()
```

## Ejercicios Prácticos

1. Implementa una jerarquía de excepciones personalizadas para una aplicación de gestión de inventario.

2. Crea un decorador que maneje diferentes tipos de excepciones de manera específica, con opciones para reintentar, registrar o propagar.

3. Diseña un context manager que implemente un sistema de transacciones para operaciones que deben ser atómicas.

4. Implementa un sistema de validación que recopile múltiples errores antes de lanzar una excepción con todos ellos.

5. Refactoriza un código existente para mejorar su manejo de excepciones siguiendo las mejores prácticas.

## Conclusión

El manejo avanzado de excepciones en Python permite crear aplicaciones más robustas y mantenibles. Al utilizar excepciones personalizadas, patrones de manejo de excepciones y seguir las mejores prácticas, puedes crear código que responda elegantemente a situaciones inesperadas.

Recuerda que el objetivo del manejo de excepciones no es solo evitar que el programa falle, sino proporcionar información útil sobre lo que salió mal y, cuando sea posible, recuperarse de manera adecuada.

Las excepciones no deben usarse para controlar el flujo normal del programa, sino para manejar situaciones excepcionales. Un buen diseño de manejo de excepciones mejora la legibilidad, mantenibilidad y robustez de tu código.