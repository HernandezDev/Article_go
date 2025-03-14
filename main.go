package main

import (
	"fmt"
	"strings"

	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

type Articulo struct {
	Id     int
	Nombre string
	Precio float64
}

func main() {
	// Crear la aplicación
	myApp := app.New()
	myWindow := myApp.NewWindow("Ventana con Pestañas")

	// abrir base de datos
	db, err := sql.Open("sqlite3", "./Base.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	// Crea una tabla
	createTable :=
		`CREATE TABLE IF NOT EXISTS Articulos(
	Id INTEGER PRIMARY KEY AUTOINCREMENT, Nombre TEXT NOT NULL, Precio REAL NOT NULL);`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}
	// crea un índice único para el campo Nombre
	createIndexNombre :=
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_nombre ON Articulos(Nombre);`
	_, err = db.Exec(createIndexNombre)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Trigger para evitar que se inserten registros con el campo Nombre vacío o solo con espacios
	createTriggerInsertNombre :=
		`CREATE TRIGGER IF NOT EXISTS Insert_Nombre
		BEFORE INSERT ON Articulos
		FOR EACH ROW
			BEGIN
				SELECT CASE 
				WHEN TRIM(NEW.Nombre) = '' THEN
				RAISE(ABORT, 'Nombre no puede estar vacío ni contener solo espacios')
			END;
		END;`
	_, err = db.Exec(createTriggerInsertNombre)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Trigger para evitar que se actualicen registros con el campo Nombre vacío o solo con espacios
	createTriggerUpdateNombre :=
		`CREATE TRIGGER IF NOT EXISTS Update_Nombre
		BEFORE UPDATE ON Articulos		
		FOR EACH ROW
			BEGIN
	 			SELECT CASE 
				WHEN TRIM(NEW.Nombre) = '' THEN
				RAISE(ABORT, 'Nombre no puede estar vacío ni contener solo espacios')
			END;
		END;`
	_, err = db.Exec(createTriggerUpdateNombre)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Trigger para evitar que se inserten registros con el campo Precio menor o igual a cero
	createTriggerInsertPrecio :=
		`CREATE TRIGGER IF NOT EXISTS Insert_Precio
		BEFORE INSERT ON Articulos
		FOR EACH ROW
			BEGIN
				SELECT CASE
				WHEN NEW.Precio <= 0 THEN
				RAISE(ABORT, 'Precio debe ser mayor a cero')
			END;
		END;`
	_, err = db.Exec(createTriggerInsertPrecio)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Trigger para evitar que se actualicen registros con el campo Precio menor o igual a cero
	createTriggerUpdatePrecio :=
		`CREATE TRIGGER IF NOT EXISTS Update_Precio
		BEFORE UPDATE ON Articulos
		FOR EACH ROW
			BEGIN
				SELECT CASE
				WHEN NEW.Precio <= 0 THEN
				RAISE(ABORT, 'Precio debe ser mayor a cero')
			END;
		END;`
	_, err = db.Exec(createTriggerUpdatePrecio)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Inicializar tabs con las pestañas iniciales
	tabs := container.NewAppTabs(
		Cargar(db),
		Consultar(db),
		Mostrar(db),
	)
	// Configurar la ventana principal
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(854, 480))
	myWindow.ShowAndRun()

}

func Cargar(db *sql.DB) *container.TabItem {

	// Crear widgets
	entry := widget.NewEntry()
	button := widget.NewButton("Cargar", func() {
		// Acción del botón
		fmt.Println("Botón presionado")
	})
	NumEntry := widget.NewEntry()
	NumEntry.OnChanged = func(content string) {
		// Filtrar contenido para permitir solo números y puntos decimales
		NumEntry.SetText(filterNumeric(content))
	}

	// Crear un contenedor para organizar los widgets
	a := container.NewVBox(
		widget.NewLabel("Nombre:"),
		entry,
		widget.NewLabel("Precio:"),
		NumEntry,
		layout.NewSpacer(),
		button,
	)

	// Retornar la pestaña con el contenedor como contenido
	return container.NewTabItem("Cargar Articulo", a)
}

func Consultar(db *sql.DB) *container.TabItem {
	button1 := widget.NewButton("Consultar", func() {
		// Acción del botón
		fmt.Println("Botón presionado")
	})
	button2 := widget.NewButton("Editar", func() {
		// Acción del botón
		fmt.Println("Modo de edición activado")

	})
	button3 := widget.NewButton("Eliminar", func() {
		// Acción del botón
		fmt.Println("Botón presionado")
	})

	// Crear un contenedor para organizar los widgets
	a := container.NewVBox(
		layout.NewSpacer(),
		button1,
		button2,
		button3,
	)
	return container.NewTabItem("Consultar por ID", a)
}

func Mostrar(db *sql.DB) *container.TabItem {
	a := widget.NewLabel("Contenido de la Pestaña 3")
	return container.NewTabItem("Listado Completo", a)
}

func filterNumeric(content string) string {
	return strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' { // Permitir números y el punto decimal
			return r
		}
		return -1 // Eliminar caracteres no válidos
	}, content)
}
