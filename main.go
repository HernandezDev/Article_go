package main

import (
	"fmt"

	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Crear la aplicación
	myApp := app.New()
	myWindow := myApp.NewWindow("Ventana con Pestañas")
	// Crear los contenidos para las pestañas

	// Inicializar tabs con las pestañas iniciales
	tabs := container.NewAppTabs(
		createTab1(),
		createTab2(),
		createTab3(),
	)
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
	Id INTEGER PRIMARY KEY AUTOINCREMENT, Nombre TEXT, Precio REAL);`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Configurar la ventana principal
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(854, 480))
	myWindow.ShowAndRun()

}

func createTab1() *container.TabItem {
	// Crear widgets
	entry := widget.NewEntry()
	button := widget.NewButton("Cargar", func() {
		// Acción del botón
		fmt.Println("Botón presionado")
	})
	NumEntry := widget.NewEntry()

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

func createTab2() *container.TabItem {
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

func createTab3() *container.TabItem {
	a := widget.NewLabel("Contenido de la Pestaña 3")
	return container.NewTabItem("Listado Completo", a)
}
