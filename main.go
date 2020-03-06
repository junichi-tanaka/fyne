package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/pkg/errors"
	osdialog "github.com/sqweek/dialog"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob"
)

var (
	gcsURL    string
	uploadFiles []string
)

func uploadFile(ctx context.Context, urlstr, file string) error {
	bucket, err := blob.OpenBucket(ctx, urlstr)
	if err != nil {
		return errors.Wrap(err, "failed to open blob")
	}
	defer bucket.Close()

	key := filepath.Base(file)
	w, err := bucket.NewWriter(ctx, key, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create the writer")
	}

	f, err := os.Open(file)
	if err != nil {
		return errors.Wrap(err, "failed to open src file")
	}
	defer f.Close()

	written, err := io.Copy(w, f)
	if err != nil {
		return errors.Wrap(err, "failed to copy file")
	}
	fmt.Println(written)

	err = w.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close the writer")
	}

	return nil
}

func main() {
	ctx := context.Background()

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
		uploadFiles = uploadFiles[:0]
		for _, file := range files {
			vbox.Append(
				widget.NewLabel(file),
			)
			uploadFiles = append(uploadFiles, file)
		}
		vbox.Refresh()
	})

	gcsPathEntry := widget.NewEntry()
	uploadButton := widget.NewButton("Upload!", func() {
		fmt.Println("upload!")
		if len(uploadFiles) == 0 {
			dialog.ShowError(fmt.Errorf("Please choose a folder!"), w)
			return
		}

		ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second*300)
		gcsURL = gcsPathEntry.Text
		for _, file := range uploadFiles {
			fmt.Println(file)
			err := uploadFile(ctxWithTimeout, gcsURL, file)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
		}
	})

	gcsPathFormItem := widget.NewFormItem("gcs url", gcsPathEntry)
	vbox2 := widget.NewVBox(
		widget.NewForm(
			gcsPathFormItem,
		),
		uploadButton,
	)
	gcsPathEntry.SetText("")
	gcsPathEntry.SetPlaceHolder("gs://...")

	box := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(chooseButton, vbox2, nil, nil),
		chooseButton, fileListContainer, vbox2,
	)
	w.SetContent(box)

	w.ShowAndRun()
}
