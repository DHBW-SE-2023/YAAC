package yaac_mvvm

import (
	"fmt"

	yaac_backend_opencv "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
	yaac_frontend_opencv "github.com/DHBW-SE-2023/YAAC/internal/frontend/opencv"
)

func (m *MVVM) StartGoCV(img_path string) {
	_ = yaac_backend_opencv.New(m)
	gocv_window := yaac_frontend_opencv.New(m)

	var msg string
	var suc bool
	ch := make(chan int)
	go func() {
		// msg, suc = backend.StartGoCV(img_path, ch)
	}()
	for elem := range ch {
		gocv_window.UpdateProgress(float64(elem) / 100)
	}
	if suc {
		gocv_window.ShowGeneratedImage(msg)
	} else {
		fmt.Println(msg)
	}

	//fmt.Printf("Done!\n\tSuccess: %s\n\tMessage: %s\n", fmt.Sprint(suc), msg)
}
