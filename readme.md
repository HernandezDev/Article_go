# Article_GO

Este proyecto es una aplicación desarrollada en Go utilizando el framework Fyne. La aplicación permite gestionar artículos en una base de datos SQLite y está disponible tanto para escritorio como para dispositivos Android.

## Estructura del Proyecto

```
go.mod
go.sum
main.go
```

## Dependencias

Las dependencias del proyecto están especificadas en el archivo `go.mod`:

- `fyne.io/fyne/v2` v2.5.5
- `github.com/mattn/go-sqlite3` v1.14.24
- Otras dependencias indirectas

## Instalación

### Opción 1: Descargar la aplicación ya compilada

Puedes descargar tanto la versión de escritorio como la de Android desde la página de [releases](https://github.com/HernandezDev/Article_go/releases/tag/v0.2).

### Opción 2: Compilar desde el código fuente

1. Clona el repositorio:
    ```sh
    git clone <URL_DEL_REPOSITORIO>
    ```
2. Navega al directorio del proyecto:
    ```sh
    cd fyne
    ```
3. Instala las dependencias:
    ```sh
    go mod tidy
    ```

## Uso

### En Escritorio

Para ejecutar la aplicación en un entorno de escritorio, utiliza el siguiente comando:
```sh
go run main.go
```

Si deseas generar un archivo ejecutable:
```sh
go build -o Article_GO.exe main.go
```
Esto creará el archivo `Article_GO.exe` en el directorio actual.

### En Android

Para compilar y ejecutar la aplicación en Android, asegúrate de tener configurado el entorno de desarrollo para Android con Go y Fyne.  
Luego, utiliza el siguiente comando para compilar:
```sh
fyne package -os android -appID com.example.article_go
```

## Funcionalidades

### Cargar Artículo

Permite cargar un nuevo artículo en la base de datos especificando su nombre y precio.

### Consultar por ID

Permite consultar un artículo por su ID, mostrando su nombre y precio. También permite editar o eliminar el artículo consultado.

### Listado Completo

Muestra un listado completo de todos los artículos en la base de datos.


