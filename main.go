package main

import (
	"fmt"
	"strconv"
	"strings"

	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
		Cargar(db, &myWindow),
		Consultar(db, &myWindow),
		Mostrar(db, &myWindow),
	)
	// Configurar la ventana principal
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(854, 480))
	myWindow.ShowAndRun()

}

func Cargar(db *sql.DB, myWindow *fyne.Window) *container.TabItem {
	var precio float64
	// Crear widgets
	Entry := widget.NewEntry()

	NumEntry := widget.NewEntry()
	NumEntry.OnChanged = func(content string) {
		// Filtrar contenido para permitir solo números y puntos decimales
		NumEntry.SetText(filterNumeric(content))
	}
	button := widget.NewButton("Cargar", func() {
		insertArticulo :=
			`INSERT INTO Articulos(Nombre, Precio) VALUES (?, ?)`
		// Convertir el texto a float64
		precio, _ = strconv.ParseFloat(NumEntry.Text, 64)
		_, err := db.Exec(insertArticulo, Entry.Text, precio)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: Articulos.Nombre" {
				dialog.NewInformation("Error", "No se pueden repetir nombres.", *myWindow).Show()
			} else {
				dialog.NewError(err, *myWindow).Show()
			}

		} else {
			dialog.NewInformation("Articulo", "Articulo cargado correctamente", *myWindow).Show()
		}
		// Restablecer valores después de la inserción
		Entry.SetText("")
		NumEntry.SetText("")
		precio = 0

	})
	// Crear un contenedor para organizar los widgets
	a := container.NewVBox(
		widget.NewLabel("Nombre:"),
		Entry,
		widget.NewLabel("Precio:"),
		NumEntry,
		layout.NewSpacer(),
		button,
	)

	// Retornar la pestaña con el contenedor como contenido
	return container.NewTabItem("Cargar Articulo", a)
}

func Consultar(db *sql.DB, myWindow *fyne.Window) *container.TabItem {
	var Id int
	var Nombre string
	var Precio string
	//labels dinámicas
	LabNombre := widget.NewLabel("")
	LabPrecio := widget.NewLabel("")

	//Entry
	Entry := widget.NewEntry()
	Entry.OnChanged = func(content string) {
		// Filtrar contenido para permitir solo números
		filteredContent := filterInt(content)
		Entry.SetText(filteredContent)

		// Convertir el texto filtrado a un int
		if filteredContent != "" {
			var err error
			Id, err = strconv.Atoi(filteredContent)
			if err != nil {
				fmt.Println("Error de conversión:", err)
				Id = 0 // Valor por defecto si hay error
			}
		} else {
			Id = 0 // Valor por defecto si el texto está vacío
		}
		LabNombre.SetText("")
		LabPrecio.SetText("")
	}

	//botones
	BotConsultar := widget.NewButton("Consultar", func() {
		// Acción del botón
		row := db.QueryRow("SELECT Nombre, Precio FROM Articulos WHERE Id = ?", Id)
		err := row.Scan(&Nombre, &Precio)
		if err != nil {
			fmt.Println("Error:", err)
		}
		LabNombre.SetText(Nombre)
		LabPrecio.SetText(Precio)
	})
	BotEditar := widget.NewButton("Editar", func() {
		// Acción del botón
		fmt.Println("Modo de edición activado")

	})
	BotEliminar := widget.NewButton("Eliminar", func() {
		// Acción del botón
		fmt.Println("Botón presionado")
	})

	// Crear un contenedor para organizar los widgets
	a := container.NewVBox(
		widget.NewLabel("Id:"),
		Entry,
		widget.NewLabel("Nombre:"),
		container.NewHBox(layout.NewSpacer(), LabNombre, layout.NewSpacer()),
		widget.NewLabel("Precio:"),
		container.NewHBox(layout.NewSpacer(), LabPrecio, layout.NewSpacer()),
		layout.NewSpacer(),
		BotConsultar,
		BotEditar,
		BotEliminar,
	)
	return container.NewTabItem("Consultar por ID", a)
}

func Mostrar(db *sql.DB, myWindow *fyne.Window) *container.TabItem {
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

func filterInt(content string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' { // Permitir números
			return r
		}
		return -1 // Eliminar caracteres no válidos
	}, content)
}
