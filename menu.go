package main

import (
	"fmt"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/tidwall/pretty"
)

// Create the new menu. Returns a pointer to fyne.MainMenu.
//
// Usage: win.SetMainMenu(createMenuItems())
func (config *Config) CreateMenuItems(win fyne.Window) *fyne.MainMenu {
	// File Menu
	openMenuItem := fyne.NewMenuItem("Open", config.openFile(win))
	saveMenuItem := fyne.NewMenuItem("Save", config.saveFile(win))
	// Disable save until we get a save as or open.
	saveMenuItem.Disabled = true
	config.SaveMenuItem = saveMenuItem

	saveAsMenuItem := fyne.NewMenuItem("Save As", config.saveAs(win))
	saveEscapeMenuItem := fyne.NewMenuItem("Save HTML Escaped Line", config.saveEscaped(win))
	exitMenuItem := fyne.NewMenuItem("Exit", func() {
		os.Exit(1)
	})
	// Order Matters.
	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem, saveEscapeMenuItem, exitMenuItem)

	// Tools menu
	prettyMenuItem := fyne.NewMenuItem("Pretty", func() {
		config.Action = "EditPretty"
		config.onChange(config.PreviewWidget.Text())
	})
	uglyMenuItem := fyne.NewMenuItem("One Line", func() {
		config.Action = "Ugly"
		config.onChange(config.PreviewWidget.Text())
	})
	toolMenu := fyne.NewMenu("Tools", prettyMenuItem, uglyMenuItem)

	copyPrettyItem := fyne.NewMenuItem("Copy Pretty", func() {
		win.Clipboard().SetContent(config.PreviewWidget.Text())
		d := dialog.NewInformation("Copy Pretty", "Copied Right side (Pretty) to clipboard.", win)
		d.Show()
	})
	copyOneLineItem := fyne.NewMenuItem("Copy One line", func() {
		content := config.PreviewWidget.Text()
		pp := pretty.Ugly([]byte(content))
		win.Clipboard().SetContent(string(pp))
		dialog.NewInformation("Copy One Line", "Copied Left side (One Line) to clipboard.", win).Show()
	})
	copyEscapeItem := fyne.NewMenuItem("Copy HTML Escape", func() {
		content := serialize(config.PreviewWidget.Text())
		win.Clipboard().SetContent(content)
		d := dialog.NewInformation("HTML Escaped String", "HTML Escaped string has been copied to your clipboard", win)
		d.Show()

	})
	editMenu := fyne.NewMenu("Edit", copyPrettyItem, copyOneLineItem, copyEscapeItem)
	return fyne.NewMainMenu(fileMenu, editMenu, toolMenu)
}

func (config *Config) saveAs(win fyne.Window) func() {
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
			// User selects a file, then write the contents of preview widget to the file. We don't
			// want the edit widget, we want to write the json that is formatted.
			write.Write([]byte(config.PreviewWidget.Text()))
			// When the user saves the file, that is stored in write.URI(). Need to copy this
			// to the config so that it can be used outside this function.
			config.CurrentFile = write.URI()
			// Defer puts off when the function ends to close the file. We could do it later but
			// putting it after the write makes sure we don't forget it.
			defer write.Close()
			// Now the title of the window can contain the file name. win.Title, and
			// write.URI().Name() return strings.
			newTitle := fmt.Sprintf("%s - %s", win.Title(), write.URI().Name())
			win.SetTitle(newTitle)
			// Now that it's saved, the save menu can now be enabled.
			config.SaveMenuItem.Disabled = false
			// Save the fyne.URI to config.
			config.CurrentFile = write.URI()
		}, win)
		// So now we can display the saveDialog box
		saveDialog.Show()
	}
}
func (config *Config) saveFile(win fyne.Window) func() {
	return func() {
		if config.CurrentFile == nil {
			return
		}
		write, err := storage.Writer(config.CurrentFile)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		write.Write([]byte(config.PreviewWidget.Text()))
		defer write.Close()
	}
}

func (config *Config) openFile(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if read == nil {
				dialog.NewInformation("Open File", "Can't open file", win)
				return
			}
			data, err := io.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			config.EditWidget.SetText(string(data))
			config.SaveMenuItem.Disabled = false
			newTitle := fmt.Sprintf("%s - %s", win.Title(), read.URI().Name())
			win.SetTitle(newTitle)
			config.CurrentFile = read.URI()
		}, win)
		openDialog.Show()
	}
}
