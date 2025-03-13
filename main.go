package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Crear la aplicación
	myApp := app.New()
	myWindow := myApp.NewWindow("Ventana con Pestañas")

	// Crear los contenidos para las pestañas
	tab1 := container.NewTabItem("Pestaña 1", widget.NewLabel("Contenido de la Pestaña 1"))
	tab2 := container.NewTabItem("Pestaña 2", widget.NewLabel("Contenido de la Pestaña 2"))
	tab3 := container.NewTabItem("Pestaña 3", widget.NewLabel("Contenido de la Pestaña 3"))

	// Crear las pestañas y añadirlas al contenedor
	tabs := container.NewAppTabs(tab1, tab2, tab3)

	// Configurar la ventana principal
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
