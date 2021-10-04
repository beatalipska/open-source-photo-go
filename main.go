package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/frameloss/prettyfyne"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type fyneapp struct {
	fyne        fyne.App
	parentW     fyne.Window
	inputPath   string
	outputPath  string
	labelInput  *widget.Label
	labelOutput *widget.Label
	logo        *canvas.Image
	progress    *widget.ProgressBar
}

func main() {
	f := app.New()
	w := f.NewWindow("")

	fyneapp := &fyneapp{
		fyne:        f,
		parentW:     w,
		labelInput:  widget.NewLabel("no input selected"),
		labelOutput: widget.NewLabel(""),
		logo:        canvas.NewImageFromFile("logo.png"),
		progress:    widget.NewProgressBar(),
	}

	myTheme := prettyfyne.PrettyTheme{
		BackgroundColor:     color.RGBA{R: 223, G: 224, B: 233, A: 255},
		ButtonColor:         color.RGBA{R: 123, G: 136, B: 167, A: 255},
		DisabledButtonColor: color.RGBA{R: 15, G: 15, B: 17, A: 255},
		HyperlinkColor:      color.RGBA{R: 170, G: 100, B: 20, A: 64},
		TextColor:           color.RGBA{R: 44, G: 50, B: 66, A: 255},
		DisabledTextColor:   color.RGBA{R: 155, G: 155, B: 155, A: 127},
		IconColor:           color.RGBA{R: 150, G: 80, B: 0, A: 255},
		DisabledIconColor:   color.RGBA{R: 155, G: 155, B: 155, A: 127},
		PlaceHolderColor:    color.RGBA{R: 150, G: 80, B: 0, A: 255},
		PrimaryColor:        color.RGBA{R: 110, G: 40, B: 0, A: 127},
		HoverColor:          color.RGBA{R: 54, G: 72, B: 109, A: 127},
		FocusColor:          color.RGBA{R: 99, G: 99, B: 99, A: 255},
		ScrollBarColor:      color.RGBA{R: 35, G: 35, B: 35, A: 8},
		ShadowColor:         color.RGBA{R: 0, G: 0, B: 0, A: 64},
		TextSize:            13,
		TextFont:            theme.DarkTheme().TextFont(),
		TextBoldFont:        theme.DarkTheme().TextBoldFont(),
		TextItalicFont:      theme.DarkTheme().TextItalicFont(),
		TextBoldItalicFont:  theme.DarkTheme().TextBoldItalicFont(),
		TextMonospaceFont:   theme.DarkTheme().TextMonospaceFont(),
		Padding:             4,
		IconInlineSize:      24,
		ScrollBarSize:       10,
		ScrollBarSmallSize:  4,
	}

	fyne.CurrentApp().Settings().SetTheme(myTheme.ToFyneTheme())

	label1 := widget.NewLabelWithStyle("PhotoDB", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	selectB := widget.NewButton("Select input", fyneapp.selectInput)
	inputLabel := widget.NewLabelWithStyle("Input Folder: ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	inputLabel2 := fyneapp.labelInput
	outputLabel := widget.NewLabelWithStyle("Output Folder: ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	outputLabel2 := fyneapp.labelOutput
	copy := widget.NewButton("Copy", fyneapp.copy)
	orLabel := widget.NewLabelWithStyle("or", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	move := widget.NewButton("Move", fyneapp.move)
	copywrightLabel := widget.NewLabel("Created with heart at Paddle.com by Deanna Green, Mario Arranz, Jeppe Fensholt, Michał Dobrzyński and Beata Lipska")
	fyneapp.progress.SetValue(float64(0.00))

	fyneapp.logo = &canvas.Image{
		File:         "logo.png",
		Resource:     nil,
		Image:        nil,
		Translucency: 0,
		FillMode:     canvas.ImageFillOriginal,
		ScaleMode:    0,
	}

	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewVBoxLayout(),
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), label1, layout.NewSpacer()),
			fyneapp.logo,
			layout.NewSpacer(),
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), selectB, layout.NewSpacer()),
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), inputLabel, inputLabel2, layout.NewSpacer()),
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), outputLabel, outputLabel2, layout.NewSpacer()),
			layout.NewSpacer(),
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), copy, orLabel, move, layout.NewSpacer()),
			fyneapp.progress,
			layout.NewSpacer(),
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), copywrightLabel, layout.NewSpacer()),
		),
	)

	w.Resize(fyne.Size{Height: 500, Width: 500})
	w.ShowAndRun()
}

func (a *fyneapp) showInfo() {
	dialog.ShowInformation("", "test info", a.parentW)
}

func (a *fyneapp) selectInput() {
	dialog.ShowFileOpen(a.selectInputCallback, a.parentW)
}

func (a *fyneapp) selectInputCallback(rc fyne.URIReadCloser, err error) {
	if err != nil {
		dialog.ShowError(err, a.parentW)
		return
	}
	if rc == nil {
		return
	}
	// keep only the folder, not the file
	uri := rc.URI().String()
	inputPath := uri

	// slashDelimiter := "/"
	// inputList := strings.Split(uri, slashDelimiter)
	// inputPath := strings.Join(inputList[:len(inputList)-1], "/")

	fileDelimiter := "file://"
	if strings.HasPrefix(inputPath, fileDelimiter) {
		pathList := strings.Split(uri, fileDelimiter)
		inputPath = strings.Join(pathList[1:], "")
	}

	a.inputPath = inputPath
	a.labelInput.SetText(a.inputPath)
}

func (a *fyneapp) copy() {
	ok := isValidPath(a.inputPath, a.outputPath, &a.parentW)
	if !ok {
		return
	}
	err := organise(a.inputPath, a.outputPath, "copy", &a.parentW, a.progress)
	if err != nil {
		dialog.ShowError(err, a.parentW)
		return
	}
}

func (a *fyneapp) move() {
	ok := isValidPath(a.inputPath, a.outputPath, &a.parentW)
	if !ok {
		return
	}
	err := organise(a.inputPath, a.outputPath, "move", &a.parentW, a.progress)
	if err != nil {
		dialog.ShowError(err, a.parentW)
		return
	}
}

func setBarValue(v int) float64 {
	numberOfFiles := 10
	p := float64(v) / float64(numberOfFiles)
	return p
}

func isValidPath(input, output string, parentW *fyne.Window) bool {
	if input == "" {
		err := fmt.Errorf("%s", "no input selected")
		dialog.ShowError(err, *parentW)
		return false
	}
	return true
}

func (a *fyneapp) quit() {
	a.fyne.Quit()
}
