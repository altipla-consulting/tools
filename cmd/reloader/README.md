
# reloader

Ejecuta los tests o un binario de Go cada vez que cambia algún fichero del código fuente. Se queda esperando a los cambios y los ejecuta automáticamente al guardar para facilitar el desarrollo en local.


## Tests

Ejecutar los tests de uno o varios paquetes a cada cambio:
```shell
reloader test ./pkg/foo ./pkg/bar
```

Ejecutar los tests en modo *verbose* que muestra toda la salida:
```shell
reloader test -v ./pkg/foo
```

Ejecutar un test concreto por nombre:
```shell
reloader test -v ./pkg/foo -r TestNameHere$
```

Ejecutar todos los tests que comienzan por un prefijo:
```shell
reloader test -v ./pkg/foo -r TestGet
```


## Aplicaciones

Ejecutarlas y reiniciar al hacer cambios en el paquete:
```shell
reloader run ./cmd/myapp
```

Monitorizar carpetas adicionales en busca de cambios:
```shell
reloader run ./cmd/myapp -w ./pkg
```

Reiniciar la aplicación cuando cambie el código o cambien ficheros de configuración:
```shell
reloader run ./pkg/foo ./pkg/bar -e .pbtext -e .yml
```

Reiniciar la aplicación cada pocos segundos cuando se cierra inesperadamente Eespecialmente útil para los servidores que deben abrirse tras un pánico para seguir experimentando:

```shell
reloader run ./pkg/foo -r
```
