package yaac_frontend_opencv

import (
	"image/color"
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

var gv GlobalVars

type GlobalVars struct {
	App                  fyne.App
	Window               fyne.Window
	ImagePath            string
	ProgBar              *widget.ProgressBar
	InputImageContainer  *fyne.Container
	OutputImageContainer *fyne.Container
}

func (f *WindowOpenCV) Open() {
	gv = GlobalVars{}
	gv.App = *yaac_shared.GetApp()

	// setuping window
	gv.Window = gv.App.NewWindow("OpenCV Demo")

	// set icon
	gv.Window.SetIcon(yaac_shared.ResourceIconPng)

	// handle main window
	gv.Window.SetContent(makeWindow(f))
	gv.Window.Resize(fyne.NewSize(800, 600))
	gv.Window.Show()

	gv.App.Run()
}

func (f *WindowOpenCV) UpdateProgress(value float64) {
	gv.ProgBar.SetValue(value)
}

func makeWindow(f *WindowOpenCV) *fyne.Container {
	header := widget.NewLabel("Please select an Input image:")

	inputImage := canvas.NewLinearGradient(color.Transparent, color.Black, 0)
	inputImageScroll := container.NewScroll(inputImage)
	gv.InputImageContainer = container.NewAdaptiveGrid(1, inputImageScroll)
	gv.OutputImageContainer = container.NewAdaptiveGrid(1, inputImageScroll)

	openFile := widget.NewButton("File Open With Filter (.jpg or .png)", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, gv.Window)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			gv.ImagePath = reader.URI().Path()
			showImage(reader, gv.InputImageContainer)
		}, gv.Window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})

	startOpenCV := widget.NewButton("Run OpenCV", func() {
		f.MVVM.StartGoCV(gv.ImagePath)
	})
	gv.ProgBar = widget.NewProgressBar()

	/*
		box := container.NewVBox(
			header,
			openFile,
			inputImageContainer,
			startOpenCV,
			opencvDemoWindow.ProgBar,
		)
	*/
	return container.NewAdaptiveGrid(1, container.NewScroll(container.NewAdaptiveGrid(
		1,
		container.NewVBox(
			header,
			openFile,
		),
		gv.InputImageContainer,
		container.NewVBox(
			startOpenCV,
			gv.ProgBar,
		),
		gv.OutputImageContainer,
	)))
}

func (f *WindowOpenCV) ShowGeneratedImage(out_Path string) {
	// Load the image resource directly from the file path
	res, err := fyne.LoadResourceFromPath(out_Path)
	if err != nil {
		log.Println("Error loading generated image:", err)
		return
	}

	// Create an image from the resource
	img := canvas.NewImageFromResource(res)
	if img == nil {
		log.Println("Error creating image from resource")
		return
	}

	img.FillMode = canvas.ImageFillContain

	imgScroll := container.NewScroll(img)
	gv.OutputImageContainer.Objects = []fyne.CanvasObject{imgScroll}

	// Refresh the content to display the new image
	gv.Window.Content().Refresh()
	gv.Window.RequestFocus()
	gv.Window.Show()
}

func loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := io.ReadAll(f)
	if err != nil {
		fyne.LogError("Error at loading file", err)
		return nil
	}
	res := fyne.NewStaticResource(f.URI().Name(), data)
	img := canvas.NewImageFromResource(res)
	if img == nil {
		fyne.LogError("Error at creating file object", err)
		return nil
	}

	return img
}

func showImage(f fyne.URIReadCloser, imgContainer *fyne.Container) {
	if f == nil {
		log.Println("Cancelled")
		return
	}
	defer f.Close()
	img := loadImage(f)
	if img == nil {
		log.Println("Error at loading image")
		return
	}

	img.FillMode = canvas.ImageFillContain

	// Create a container with dynamic sizing
	//containerWithDynamicSizing := fyne.NewContainer(img)

	// Set the content of the main container to the new container with dynamic sizing
	//imgContainer.Objects = []fyne.CanvasObject{containerWithDynamicSizing}

	//inputImage := canvas.NewImageFromFile(img.File)
	//inputImageScroll := container.NewScroll(inputImage)
	imgScroll := container.NewScroll(img)
	//imgScroll.Resize(img.Size())
	imgContainer.Objects = []fyne.CanvasObject{imgScroll}

	//imgContainer.Resize(img.Size())
	//fmt.Println(img.Size())

	// Actualize and show window
	gv.Window.Content().Refresh()
	gv.Window.RequestFocus()
	gv.Window.Show()
}
