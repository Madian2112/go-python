# 16. Usando CGO para interoperar con C

## Introducción a CGO

CGO es una herramienta que permite a los programas de Go llamar a código C. No es magia, es una herramienta que genera el "pegamento" necesario para que Go y C puedan comunicarse.

## ¿Cuándo usar CGO?

Usar CGO tiene un costo en el rendimiento de la compilación y en la llamada entre lenguajes. Generalmente, se usa cuando necesitas:
- Utilizar una librería de C existente que no tiene un equivalente en Go.
- Escribir código de muy bajo nivel o interactuar directamente con el sistema operativo.

## Temas a cubrir:

- Habilitar CGO.
- `import "C"`.
- Tipos de datos de Go vs. C.
- Llamar a funciones de C desde Go.
- Directivas de compilador y enlazador (`#cgo CFLAGS`, `#cgo LDFLAGS`).
- Pasar punteros entre Go y C.
- Reglas y gestión de memoria.
- Exportar funciones de Go para ser llamadas desde C.
