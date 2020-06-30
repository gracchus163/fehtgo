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
	img := make([]*canvas.Image, 20, 20);

	img[0] = canvas.NewImageFromFile("1.jpg")
	img[0].FillMode = canvas.ImageFillContain
	img[1] = canvas.NewImageFromFile("2.jpg")
	img[1].FillMode = canvas.ImageFillContain
	img[2] = canvas.NewImageFromFile("3.jpg")
	img[2].FillMode = canvas.ImageFillContain
	img[3] = canvas.NewImageFromFile("4.jpg")
	img[3].FillMode = canvas.ImageFillContain
	img[4] = canvas.NewImageFromFile("5.jpg")
	img[4].FillMode = canvas.ImageFillContain
	img[5] = canvas.NewImageFromFile("6.jpg")
	img[5].FillMode = canvas.ImageFillContain
	img[6] = canvas.NewImageFromFile("7.jpg")
	img[6].FillMode = canvas.ImageFillContain
	img[7] = canvas.NewImageFromFile("8.jpg")
	img[7].FillMode = canvas.ImageFillContain
	//img = append(img, canvas.NewImageFromFile("1.jpg"))
	//img = append(img, canvas.NewImageFromFile("2.jpg"))
		grid := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			 img[0], img[1], img[2], img[3])
		page := 0
		w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
			if (string(ev.Name) == "Q") {app.Quit()}
			onPress(ev, w,grid, img, &page)
		})
		w.SetContent(grid)
		w.ShowAndRun()
}
func onPress(ev *fyne.KeyEvent, w fyne.Window, grid *fyne.Container, img []*canvas.Image, page *int) {
	fmt.Println("KeyDown: "+string(ev.Name))
	if (ev.Name == "Right") {
		*page += 1
		fmt.Println(*page)
		grid = fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			img[4**page], img[1+(4**page)], img[2+(4**page)], img[3+(4**page)])
		w.SetContent(grid)
	}
	if (ev.Name == "Left") {
		*page -= 1
		if (*page < 0) {*page = 0}
		fmt.Println(*page)
		grid = fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			img[4**page], img[1+(4**page)], img[2+(4**page)], img[3+(4**page)])
		w.SetContent(grid)
	}
}
