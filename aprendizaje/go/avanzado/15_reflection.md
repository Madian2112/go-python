# 15. Reflexión (Reflection) en Go

## Introducción a la Reflexión

La reflexión es la capacidad de un programa para examinar su propia estructura, particularmente a través de los tipos; es una forma de metacrogramación. En Go, esto se logra con el paquete `reflect`.

## ¿Por qué usar Reflexión?

La reflexión es una herramienta muy potente, pero también compleja y, a menudo, más lenta que el código convencional. Se debe usar con precaución.

## Temas a cubrir:

- `reflect.Type` y `reflect.Value`.
- Inspeccionar el tipo y valor de una variable.
- Modificar variables a través de punteros de reflexión (`CanSet`).
- Inspeccionar los campos de un `struct`.
- Llamar a métodos dinámicamente (`Call`).
- Advertencias y mejores prácticas.
