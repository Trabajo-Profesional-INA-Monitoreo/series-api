# Inputs API

## Ejecución
La aplicación se puede ejecutar de forma local o contanerizada. 

Se provee un  `Makefile` para simplificar la ejecución de los comandos. Se puede ejecutar como `make <target>`

El `Makefile` tiene las acciones:
* `docker-image`: Buildea la imagen de la aplicación para Docker
* `test`: Ejecuta los tests de la aplicación
* `build`: Realiza el build de los servicios localmente. Es necesario Go a partir de 1.21 por lo menos.
* `run`: Ejecuta la aplicación

### Configuración
Se deben de configurar las siguientes propiedades:

```yaml
server:
  port: 1234 # puerto de la aplicacion

datasource:
  connection: "host={URL} user={USUARIO} password={PASSWORD} dbname={NOMBRE DE DB} port={PUERTO DE DB}"
```

## Documentación

Se puede ver la especificación de la API en `[host]:[puerto]/swagger/index.html` o en la carpeta `/docs`

Para actualizarla ejecutar `swag init` y seguir las especificaciones de la [documentación](https://github.com/swaggo/swag#api-operation)