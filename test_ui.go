package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fmt"
)

func main() {
	app := app.New()

	w := app.NewWindow("Hello")
	image1 := canvas.NewImageFromFile("1.jpg")
	image2 := canvas.NewImageFromFile("2.jpg")
	image3 := canvas.NewImageFromFile("3.jpg")
	image4 := canvas.NewImageFromFile("4.jpg")
	image1.FillMode = canvas.ImageFillContain
	image2.FillMode = canvas.ImageFillContain
	image3.FillMode = canvas.ImageFillContain
	image4.FillMode = canvas.ImageFillContain
	/*w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		image,
		widget.NewButton("Quit", func() {app.Quit()}),
		))*/
		grid := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			image1,image2, image3, image4)
		w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		fmt.Println("KeyDown: "+string(ev.Name))
		})
		w.SetContent(grid)
		w.ShowAndRun()
}
