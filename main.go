package main

import (
	"fmt"
	_ "io/ioutil"
	_ "os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	osdialog "github.com/sqweek/dialog"
)

var (
	listItems []string
)

func main() {
	app := app.New()

	w := app.NewWindow("File Uploader")
	w.Resize(fyne.NewSize(640, 480))

	vbox := widget.NewVBox(
			widget.NewLabel("This files will be uploaded!"),
		)
	fileListContainer := widget.NewScrollContainer(vbox)


	chooseButton := widget.NewButton("Choose Folder...", func() {
		dir, err := osdialog.Directory().Title("Select Folder to upload the csv files").Browse()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		pattern := dir + "/*.go"
		files, err := filepath.Glob(pattern)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		vbox.Children = make([]fyne.CanvasObject, 0)
		for _, file := range files {
			vbox.Append(
				widget.NewLabel(file),
			)
		}
		vbox.Refresh()
	})

	uploadButton := widget.NewButton("Upload!", func() {
		if len(listItems) == 0 {
			dialog.ShowError(fmt.Errorf("Please choose a folder!"), w)
			return
		}
	})

	box := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(chooseButton, uploadButton, nil, nil),
		chooseButton, fileListContainer, uploadButton,
	)
	w.SetContent(box)

	w.ShowAndRun()
}
