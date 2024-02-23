package yaac_backend_imgproc

import "gocv.io/x/gocv"

type mvvm interface {
}

type BackendImgproc struct {
	MVVM mvvm
}

func NewBackend(mvvm mvvm) *BackendImgproc {
	return &BackendImgproc{
		MVVM: mvvm,
	}
}

func (mvvm *BackendImgproc) ValidateTable(imgBuf []byte) (Table, error) {
	img, err := gocv.IMDecode(imgBuf, gocv.IMReadAnyColor)
	if err != nil {
		return Table{}, err
	}

	img = FindTable(img)
	table, err := ReviewTable(img)
	if err != nil {
		return Table{}, err
	}

	return table, nil
}
