package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Crear la aplicación
	myApp := app.New()
	myWindow := myApp.NewWindow("Ventana con Pestañas")

	// Crear los contenidos para las pestañas
	tabs := container.NewAppTabs(createTab1(),
		createTab2(),
		createTab3(),
	)

	// Configurar la ventana principal
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.ShowAndRun()
}

func createTab1() *container.TabItem {
	// Crear widgets
	entry := widget.NewEntry()
	button := widget.NewButton("Botón de acción", func() {
		// Acción del botón
		fmt.Println("Botón presionado")
	})

	// Crear un contenedor para organizar los widgets
	a := container.NewVBox(
		widget.NewLabel("Nombre:"),
		entry,
		button,
	)

	// Retornar la pestaña con el contenedor como contenido
	return container.NewTabItem("Cargar Articulo", a)
}

func createTab2() *container.TabItem {
	a := widget.NewLabel("Contenido de la Pestaña 2")
	return container.NewTabItem("Consultar por ID", a)
}

func createTab3() *container.TabItem {
	a := widget.NewLabel("Contenido de la Pestaña 3")
	return container.NewTabItem("Listado Completo", a)
}
