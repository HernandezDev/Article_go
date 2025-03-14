package main

import (
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
	a := widget.NewLabel("Contenido de la Pestaña 1")
	return container.NewTabItem("Cargar Articulo", a)
}

func createTab2() *container.TabItem {
	a := widget.NewLabel("Contenido de la Pestaña 2")
	return container.NewTabItem("Consultar por ID", a)
}

func createTab3() *container.TabItem {
	a := widget.NewLabel("Contenido de la Pestaña 3")
	return container.NewTabItem("Mostrar Lista", a)
}
