
# impsort

Ordena los imports por bloques según su procedencia.


## Uso

```shell
impsort folder1 folder2 -p code.altipla.consulting
impsort . -p code.altipla.consulting
```

Con esos argumentos separará 3 grupos:

1. Librerías de stdlib (`fmt`, `net/http`, ...)
2. Librerías de terceros (`github.com/sirupsen/logrus`, `libs.altipla.consulting`, ...)
3. Imports locales de `code.altipla.consulting` que se configura con el argumento.
