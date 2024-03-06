package pages

import (
	"log"

	imgproc "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
	shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type FrontendMain struct {
	MVVM shared.MVVM
}

var myMVVM shared.MVVM = nil

func New(mvvm shared.MVVM) *FrontendMain {
	myMVVM = mvvm
	return &FrontendMain{
		MVVM: mvvm,
	}
}

func (*FrontendMain) ReceiveNewTable(table imgproc.Table) {
	log.Fatal("Not yet implemented")
}
