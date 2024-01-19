package cv

import (
	"image"

	"gocv.io/x/gocv"
)

type Table struct {
	Image gocv.Mat
	Rows  [][]image.Rectangle
}

func NewTable(img gocv.Mat) Table {
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdOtsu)
	gocv.BitwiseNot(img, &img)

	size := img.Size()

	// TODO: This needs some more work, why 50 and 35?
	horizontalKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(size[0]/50, 1))
	verticalKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(1, size[1]/35))
	vhKernel := gocv.GetStructuringElement(gocv.MorphCross, image.Pt(3, 3))

	// TODO: Tweek this number or it takes longer than necessary
	iters := 8
	binaryHorizontal := gocv.NewMat()
	gocv.MorphologyExWithParams(img, &binaryHorizontal, gocv.MorphOpen, horizontalKernel, iters, gocv.BorderConstant)

	binaryVertical := gocv.NewMat()
	gocv.MorphologyExWithParams(img, &binaryVertical, gocv.MorphOpen, verticalKernel, iters, gocv.BorderConstant)

	vhLines := gocv.NewMat()
	gocv.AddWeighted(binaryVertical, 0.5, binaryHorizontal, 0.5, 0.0, &vhLines)

	gocv.BitwiseNot(vhLines, &vhLines)

	vhIters := 2
	gocv.ErodeWithParams(vhLines, &vhLines, vhKernel, image.Pt(-1, -1), vhIters, int(gocv.BorderConstant))

	gocv.Threshold(vhLines, &vhLines, 128.0, 255.0, gocv.ThresholdOtsu)

	return Table{}
}
