package yaac_frontend

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
	resource "github.com/DHBW-SE-2023/yaac-go-prototype/pkg/resource_manager"
)

type OpencvDemoWindow struct {
	Window               fyne.Window
	ImagePath            string
	ProgBar              *widget.ProgressBar
	InputImageContainer  *fyne.Container
	OutputImageContainer *fyne.Container
}

var opencvDemoWindow OpencvDemoWindow

func (f *Frontend) OpenOpencvDemoWindow() {
	opencvDemoWindow = OpencvDemoWindow{}

	// setuping window
	opencvDemoWindow.Window = App.NewWindow("OpenCV Demo")

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	opencvDemoWindow.Window.SetIcon(r)

	// handle main window
	opencvDemoWindow.Window.SetContent(makeOpencvDemoWindow(f))
	opencvDemoWindow.Window.Resize(fyne.NewSize(800, 600))
	opencvDemoWindow.Window.Show()

	App.Run()
}

func (f *Frontend) UpdateProgress(value float64) {
	opencvDemoWindow.ProgBar.SetValue(value)
}

func makeOpencvDemoWindow(f *Frontend) *fyne.Container {
	header := widget.NewLabel("Please select an Input image:")

	inputImage := canvas.NewLinearGradient(color.Transparent, color.Black, 0)
	inputImageScroll := container.NewScroll(inputImage)
	opencvDemoWindow.InputImageContainer = container.NewAdaptiveGrid(1, inputImageScroll)
	opencvDemoWindow.OutputImageContainer = container.NewAdaptiveGrid(1, inputImageScroll)

	openFile := widget.NewButton("File Open With Filter (.jpg or .png)", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, opencvDemoWindow.Window)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			opencvDemoWindow.ImagePath = reader.URI().Path()
			showImage(reader, opencvDemoWindow.InputImageContainer)
		}, opencvDemoWindow.Window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})

	startOpenCV := widget.NewButton("Run OpenCV", func() {
		f.MVVM.StartGoCV(opencvDemoWindow.ImagePath)
	})
	opencvDemoWindow.ProgBar = widget.NewProgressBar()

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
		opencvDemoWindow.InputImageContainer,
		container.NewVBox(
			startOpenCV,
			opencvDemoWindow.ProgBar,
		),
		opencvDemoWindow.OutputImageContainer,
	)))
}

func (f *Frontend) ShowGeneratedImage(out_Path string) {
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
	opencvDemoWindow.OutputImageContainer.Objects = []fyne.CanvasObject{imgScroll}

	// Refresh the content to display the new image
	opencvDemoWindow.Window.Content().Refresh()
	opencvDemoWindow.Window.RequestFocus()
	opencvDemoWindow.Window.Show()
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
	opencvDemoWindow.Window.Content().Refresh()
	opencvDemoWindow.Window.RequestFocus()
	opencvDemoWindow.Window.Show()
}
