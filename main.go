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
	tab1 := container.NewTabItem("Cargar Articulo", widget.NewLabel("Contenido de la Pestaña 1"))
	tab2 := container.NewTabItem("Consultar por ID", widget.NewLabel("Contenido de la Pestaña 2"))
	tab3 := container.NewTabItem("Listado Completo", widget.NewLabel("Contenido de la Pestaña 3"))

	// Crear las pestañas y añadirlas al contenedor
	tabs := container.NewAppTabs(tab1, tab2, tab3)
	// Obtener la pestaña seleccionada
	tabs.OnSelected = func(selectedTab *container.TabItem) {
		for index, tab := range tabs.Items {
			if tab == selectedTab {
				fmt.Printf("Índice de la pestaña seleccionada: %d\n", index)
				break
			}
		}
	}

	// Configurar la ventana principal
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.ShowAndRun()
}
