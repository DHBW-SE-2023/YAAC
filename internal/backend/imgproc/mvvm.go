package yaac_backend_imgproc

import (
	"github.com/otiai10/gosseract/v2"
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

// Take in an image as a byte array in a valid format
// and parse it into a table which has information about the students and
// the validity of the signature.
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
