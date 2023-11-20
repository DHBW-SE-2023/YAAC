package yaac_mvvm

import (
	"fmt"

	yaac_backend "github.com/DHBW-SE-2023/YAAC/internal/backend"
	yaac_frontend "github.com/DHBW-SE-2023/YAAC/internal/frontend"
)

func (m *MVVM) StartGoCV(img_path string) {
	backend := yaac_backend.New(m)
	frontend := yaac_frontend.New(m)

	var msg string
	var suc bool
	ch := make(chan int)
	go func() {
		msg, suc = backend.StartGoCV(img_path, ch)
	}()
	for elem := range ch {
		frontend.UpdateProgress(float64(elem) / 100)
	}
	if suc {
		frontend.ShowGeneratedImage(msg)
	} else {
		fmt.Println(msg)
	}

	//fmt.Printf("Done!\n\tSuccess: %s\n\tMessage: %s\n", fmt.Sprint(suc), msg)
}
