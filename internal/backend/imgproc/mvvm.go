package yaac_backend_imgproc

import (
	"github.com/otiai10/gosseract"
	"gocv.io/x/gocv"
)

type mvvm interface{}
type BackendImgproc struct {
	MVVM            mvvm
	tesseractClient *gosseract.Client
}

func NewBackend(mvvm mvvm) *BackendImgproc {
	tesseractClient := gosseract.NewClient()
	return &BackendImgproc{
		MVVM:            mvvm,
		tesseractClient: tesseractClient,
	}
}

func (mvvm *BackendImgproc) ValidateTable(imgBuf []byte) (Table, error) {
	img, err := gocv.IMDecode(imgBuf, gocv.IMReadAnyColor)
	if err != nil {
		return Table{}, err
	}

	img = FindTable(img)
	table, err := ReviewTable(img, mvvm.tesseractClient)
	if err != nil {
		return Table{}, err
	}

	return table, nil
}
