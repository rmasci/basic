// This is the continuation of 2.json/main.go adding menus to the file.
package main

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/pretty"
)

type Config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.TextGrid
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
	Action        string
}

var config Config

func main() {
	// Steps to creating a fyne app:
	// 1. create a fyne app
	config.Action = "Pretty"
	a := app.New()

	// 2. create a window for the app
	win := a.NewWindow("JSON Pretty")

	// 3. get the user interface
	edit, preview := config.makeUI()
	// create menu -- createMenuItems is in menu.go
	win.SetMainMenu(config.CreateMenuItems(win))
	// 4. set the content of the window
	win.SetContent(container.NewHSplit(edit, preview))

	// 5. show window and run app, size is super vga!!
	win.Resize(fyne.Size{Width: 800, Height: 600}) //anybody remember when 800x600 was a large monotor?
	win.CenterOnScreen()                           // Put it right in the middle.
	win.ShowAndRun()                               // display the app
}

func (config *Config) makeUI() (*widget.Entry, *widget.TextGrid) {
	edit := widget.NewMultiLineEntry()          // changed from richtext to just text
	edit.Wrapping = fyne.TextWrapBreak          // line breaks.
	preview := widget.NewTextGridFromString("") // Initially nothing is displayed on the right.
	config.EditWidget = edit                    // Assign editwidget the edit pane
	config.PreviewWidget = preview              // assign previewwidget the preview pane
	// edit.OnChanged equals a function. That function is run every time the window is changed rendering the json.
	// in this example we create an 'inline' function. In 1.basic it also equaled a function called. 'preview.ParseMarkdown'
	// This is more manual, but shows you could really implement any type of parser.
	edit.OnChanged = func(content string) {
		config.onChange(content)
	}
	return edit, preview // Goes back to main line 28
}

func (config *Config) onChange(content string) {
	switch config.Action {
	case "Pretty":
		j := pretty.Pretty([]byte(content))
		config.PreviewWidget.SetText(string(j))
	case "Ugly":
		j := pretty.Ugly([]byte(content))
		config.EditWidget.SetText(string(j))
		config.Action = "Pretty"
	case "EditPretty":
		j := pretty.Pretty([]byte(content))
		config.PreviewWidget.SetText(string(j))
		config.EditWidget.SetText(string(j))
		config.Action = "Pretty"
	}
}

func serialize(content string) string {
	j := pretty.Ugly([]byte(content))
	return url.PathEscape(string(j))
}

func (config *Config) saveEscaped(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if write == nil {
				// user clicked cancel, return empty
				return
			}
			// Kinda juggling here between string and bytes...
			contentBytes := pretty.UglyInPlace([]byte(config.PreviewWidget.Text()))
			content := url.QueryEscape(string(contentBytes))
			write.Write([]byte(content))
			// When the user saves the file, that is stored in write.URI(). Need to copy this
			// to the config so that it can be used outside this function.
			config.CurrentFile = write.URI()
			// Defer puts off when the function ends to close the file. We could do it later but
			// putting it after the write makes sure we don't forget it.
			defer write.Close()
		}, win)
		// So now we can display the saveDialog box
		saveDialog.Show()
	}
}
