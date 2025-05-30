package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"runtime"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Crear la aplicación
	myApp := app.New()
	myWindow := myApp.NewWindow("Article_GO")
	Canvas := myWindow.Canvas()

	// Obtener la ruta de almacenamiento interna de la aplicación
	storagePath := getStoragePath()
	dbPath := storagePath + "/Base.db"

	// abrir base de datos
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		dialog.NewError(err, myWindow).Show()
		//return
	}
	defer db.Close()
	// Crea una tabla
	createTable :=
		`CREATE TABLE IF NOT EXISTS Articulos(
	Id INTEGER PRIMARY KEY AUTOINCREMENT, Nombre TEXT NOT NULL, Precio REAL NOT NULL);`
	_, err = db.Exec(createTable)
	if err != nil {
		dialog.NewError(err, myWindow).Show()
		//return
	}
	// crea un índice único para el campo Nombre
	createIndexNombre :=
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_nombre ON Articulos(Nombre);`
	_, err = db.Exec(createIndexNombre)
	if err != nil {
		dialog.NewError(err, myWindow).Show()
		//return
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
		dialog.NewError(err, myWindow).Show()
		//return
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
		dialog.NewError(err, myWindow).Show()
		//return
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
		dialog.NewError(err, myWindow).Show()
		//return
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
		dialog.NewError(err, myWindow).Show()
		//return
	}
	// Inicializar tabs con las pestañas iniciales
	tabs := container.NewAppTabs(
		Cargar(db, &myWindow),
		Consultar(db, &myWindow, &Canvas),
		Mostrar(db, &myWindow),
	)
	// Configurar la ventana principal
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(640, 360))
	myWindow.ShowAndRun()

}

func Cargar(db *sql.DB, myWindow *fyne.Window) *container.TabItem {
	var precio float64
	// Crear widgets
	Entry := widget.NewEntry()

	NumEntry := widget.NewEntry()
	NumEntry.OnChanged = func(content string) {
		// Filtrar contenido para permitir solo números y puntos decimales
		NumEntry.SetText(filterFloat(content))
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
		container.NewGridWithColumns(2,
			widget.NewLabel("Nombre:"), Entry,
			widget.NewLabel("Precio:"), NumEntry,
		),
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), button),
	)

	// Retornar la pestaña con el contenedor como contenido
	return container.NewTabItem("Cargar Articulo", a)
}

func Consultar(db *sql.DB, myWindow *fyne.Window, Canvas *fyne.Canvas) *container.TabItem {
	var Id int
	var Nombre string
	var Precio float64
	var FuncionConsulta func()

	//labels dinámicas
	LabNombre := widget.NewLabel("")
	LabPrecio := widget.NewLabel("")

	//Entry para el id
	Entry := widget.NewEntry()
	Entry.OnChanged = func(content string) {
		// Filtrar contenido para permitir solo números
		filteredContent := filterInt(content)
		Entry.SetText(filteredContent)

		// Convertir el texto filtrado a un int
		if filteredContent != "" {
			Id, _ = strconv.Atoi(filteredContent)
		} else {
			Id = 0 // Valor por defecto si el texto está vacío
		}
		//resetear labels y variables cuando se cambia el id
		LabNombre.SetText("")
		LabPrecio.SetText("")
		Nombre = ""
		Precio = 0
	}

	//botones
	BotConsultar := widget.NewButton("Consultar", func() {
		// Acción del botón
		row := db.QueryRow("SELECT Nombre, Precio FROM Articulos WHERE Id = ?", Id)
		err := row.Scan(&Nombre, &Precio)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				dialog.NewInformation("Error", "Articulo no encontrado.", *myWindow).Show()
			} else {
				dialog.NewError(err, *myWindow).Show()
			}
			//Limpiar variables y laves en caso de error
			Nombre = ""
			Precio = 0
			Id = 0
			Entry.SetText("")
			LabNombre.SetText("")
			LabPrecio.SetText("")
		} else {
			LabNombre.SetText(Nombre)
			LabPrecio.SetText(strconv.FormatFloat(Precio, 'f', -1, 64))
		}

	})

	FuncionConsulta = BotConsultar.OnTapped //capurar la funcion anonima del boton

	BotEditar := widget.NewButton("Editar", func() {
		var popup *widget.PopUp
		IdEditarLabel := widget.NewLabel("")
		NombreEditarEntry := widget.NewEntry()
		PrecioEditarEntry := widget.NewEntry()
		PrecioEditarEntry.OnChanged = func(content string) {
			// Filtrar contenido para permitir solo números y puntos decimales
			PrecioEditarEntry.SetText(filterFloat(content))
		}
		if Id != 0 && Nombre != "" && Precio != 0 {
			IdEditarLabel.SetText(strconv.Itoa(Id))
			NombreEditarEntry.SetText(Nombre)
			PrecioEditarEntry.SetText(strconv.FormatFloat(Precio, 'f', -1, 64))
		} else {
			dialog.NewInformation("Error", "Seleccione articulo a editar.", *myWindow).Show()
			return
		}
		content := container.NewVBox(
			container.NewGridWithColumns(2,
				widget.NewLabel("Id"), IdEditarLabel,
				widget.NewLabel("Nombre"), NombreEditarEntry,
				widget.NewLabel("Precio"), PrecioEditarEntry,
			),
			container.NewHBox(
				layout.NewSpacer(),
				widget.NewButton("Editar", func() {

					precio, err := strconv.ParseFloat(PrecioEditarEntry.Text, 64)
					if err != nil {
						if err.Error() == `strconv.ParseFloat: parsing "": invalid syntax` {
							popup.Hide()
							dialog.NewInformation("Error", "Seleccione articulo a editar.", *myWindow).Show()
						} else {
							popup.Hide()
							dialog.NewError(err, *myWindow).Show()
							fmt.Println("Error:", err)
						}
						return
					}
					updateArticulo := `UPDATE Articulos SET Nombre = ?, Precio = ? WHERE Id = ?`
					_, err1 := db.Exec(updateArticulo, NombreEditarEntry.Text, precio, IdEditarLabel.Text)
					if err1 != nil {
						//popup de error
						popup.Hide() //serrar el popup para mostrar el dialogo
						dialog.NewError(err1, *myWindow).Show()

					} else {
						FuncionConsulta()
					}
					popup.Hide()
				}),
				widget.NewButton("Cancelar", func() {
					popup.Hide()
				}),
			),
		)

		popup = widget.NewModalPopUp(content, *Canvas)
		popup.Show() // Muestra el popup

	})
	BotEliminar := widget.NewButton("Eliminar", func() {
		var popup *widget.PopUp
		IdEliminarLabel := widget.NewLabel("")
		NombreEliminarLabel := widget.NewLabel("")
		PrecioEliminarLabel := widget.NewLabel("")
		if Id != 0 && Nombre != "" && Precio != 0 {
			IdEliminarLabel.SetText(strconv.Itoa(Id))
			NombreEliminarLabel.SetText(Nombre)
			PrecioEliminarLabel.SetText(strconv.FormatFloat(Precio, 'f', -1, 64))
		} else {
			dialog.NewInformation("Error", "Seleccione articulo a eliminar.", *myWindow).Show()
			return
		}

		content := container.NewGridWithColumns(2,
			widget.NewLabel("Id:"),
			IdEliminarLabel,
			widget.NewLabel("Nombre:"),
			NombreEliminarLabel,
			widget.NewLabel("Precio:"),
			PrecioEliminarLabel,
			widget.NewButton("Eliminar", func() {
				IdEliminar, err := strconv.Atoi(IdEliminarLabel.Text)
				if err != nil {
					if err.Error() == `strconv.Atoi: parsing "": invalid syntax` {
						popup.Hide()
						dialog.NewInformation("Error", "Seleccione articulo a eliminar.", *myWindow).Show()
					} else {
						popup.Hide()
						dialog.NewError(err, *myWindow).Show()
					}
				}
				deleteArticulo := `DELETE FROM Articulos WHERE Id = ?`
				_, err1 := db.Exec(deleteArticulo, IdEliminar)
				if err1 != nil {
					//popup de error
					popup.Hide() //serrar el popup para mostrar el dialogo
					dialog.NewError(err1, *myWindow).Show()
					return
				} else {
					popup.Hide()
					Entry.SetText("")
				}
			}),
			widget.NewButton("Cancelar", func() {
				popup.Hide()
			}),
		)

		popup = widget.NewModalPopUp(content, *Canvas)
		popup.Show() // Muestra el popup
	})

	// Crear un contenedor para organizar los widgets

	content := container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Id:"), Entry,
			widget.NewLabel("Nombre:"), LabNombre,
			widget.NewLabel("Precio:"), LabPrecio,
		),
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), BotConsultar, BotEditar, BotEliminar),
	)
	return container.NewTabItem("Consultar por ID", content)

}

func Mostrar(db *sql.DB, myWindow *fyne.Window) *container.TabItem {

	gridContenedor := container.NewGridWithColumns(3)

	// Botón para añadir un nuevo elemento
	button := widget.NewButton("Actualizar", func() {

		// Limpiar todos los elementos existentes en el contenedor
		gridContenedor.Objects = []fyne.CanvasObject{}

		// Ejecuta la consulta y obtiene las filas.
		rows, err := db.Query("SELECT Id, Nombre, Precio FROM Articulos")
		if err != nil {
			fmt.Println("Error al realizar la consulta:", err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var id, nombre, precio string

			err = rows.Scan(&id, &nombre, &precio)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			Elemento1 := widget.NewLabel(id)
			Elemento2 := widget.NewLabel(nombre)
			Elemento3 := widget.NewLabel(precio)
			gridContenedor.Objects = append(gridContenedor.Objects, Elemento1, Elemento2, Elemento3) // Añadir elemento
			gridContenedor.Refresh()
		} // Refrescar el contenedor para actualizar la vista
	})
	scroleableGrid := container.NewScroll(gridContenedor)
	header := container.NewGridWithColumns(3,
		widget.NewLabel("ID"), widget.NewLabel("Nombre"), widget.NewLabel("Precio"),
	)
	cabezera := container.NewVBox(button, header)
	return container.NewTabItem("Listado Completo", (container.NewBorder(cabezera, nil, nil, nil, scroleableGrid)))
}

func filterFloat(content string) string {
	return strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' { // Permitir números y el punto decimal
			return r
		}
		return -1 // Eliminar caracteres no válidos
	}, content)
}

func filterInt(content string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' { // Permitir solo números
			return r
		}
		return -1 // Eliminar caracteres no válidos
	}, content)
}

func getStoragePath() string {
	if runtime.GOOS == "android" {
		return fyne.CurrentApp().Storage().RootURI().Path()
	}
	// Para Windows y otros sistemas de escritorio
	appData := os.Getenv("APPDATA")
	if appData != "" {
		appDir := filepath.Join(appData, "Article_GO")
		os.MkdirAll(appDir, os.ModePerm)
		return appDir
	}
	// Fallback: carpeta del ejecutable
	exePath, err := os.Executable()
	if err != nil {
		return "." // Carpeta actual
	}
	return filepath.Dir(exePath)
}
