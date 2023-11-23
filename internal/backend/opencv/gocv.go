package yaac_backend_opencv

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

// Performs a warp-perspective transform on the image under image_path
//
//	 Input: image_path string -> Path to the image in the FS
//	 Output: (string, bool)
//
//		string -> Depends on bool
//			true -> Path to the transormed image in the FS
//			false -> Error Message
//		bool -> If the trasnformation was successfull
func (b *BackendOpenCV) StartGoCV(image_path string, prog chan int) (string, bool) {
	if prog != nil {
		defer close(prog)
		prog <- 0
	}

	image := gocv.IMRead(image_path, gocv.IMReadColor)
	msg := ""
	suc := false
	if image.Empty() {
		msg = "Empty image"
		return msg, suc
	} else {
		//fmt.Print("Loaded image!\n\tSize: ")
		//fmt.Println(image.Size())
		debug := false
		err := transform_image(image, debug, prog)
		if err != nil {
			//fmt.Println(err)
			return err.Error(), suc
		}
		out_path, saved := save_image_state("list_output", image, debug)
		if prog != nil {
			prog <- 100
		}
		if saved {
			return out_path, saved
		} else {
			msg = "ERROR: Failed to save image"
			//fmt.Println(msg)
			return msg, suc
		}
	}
}

func save_image_state(state_name string, img gocv.Mat, debug bool) (string, bool) {
	if debug {
		fmt.Printf("%s:\n\tSize: ", state_name)
		fmt.Print(img.Size())
		fmt.Print(" Saved: ")
	}

	path := fmt.Sprintf("./assets/%s.jpg", state_name)
	suc := gocv.IMWrite(path, img)
	if debug {
		fmt.Println(suc)
	}

	return path, suc
}

func transform_image(img gocv.Mat, debug bool, prog chan int) error {
	// Save initial State
	if debug {
		save_image_state("state00_initial", img, debug)
	}
	tmp := img.Clone()
	if prog != nil {
		prog <- 1
	}

	// Make grey
	gocv.CvtColor(tmp, &tmp, gocv.ColorBGRToGray)
	image_size := tmp.Size()
	if debug {
		save_image_state("state01_grey", tmp, debug)
	}
	if prog != nil {
		prog <- 10
	}

	// Blur and theshhold
	gocv.GaussianBlur(tmp, &tmp, image.Point{3, 3}, 2, 0, gocv.BorderDefault)
	//gocv.AdaptiveThreshold(tmp, &tmp, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, 11, 2)
	gocv.Threshold(tmp, &tmp, 128, 255, gocv.ThresholdOtsu) // Better according to @Leander
	gocv.FastNlMeansDenoisingWithParams(tmp, &tmp, 11, 31, 9)
	if debug {
		save_image_state("state02_blur_thresh", tmp, debug)
	}
	if prog != nil {
		prog <- 20
	}

	/*
	   # We need two threshold values, minVal and maxVal. Any edges with intensity gradient more than maxVal
	   # are sure to be edges and those below minVal are sure to be non-edges, so discarded.
	   #  Those who lie between these two thresholds are classified edges or non-edges based on their connectivity.
	   # If they are connected to "sure-edge" pixels, they are considered to be part of edges.
	   #  Otherwise, they are also discarded
	*/
	edges := gocv.NewMat()
	gocv.Canny(tmp, &edges, 50, 150) // FIXME - apertureSize=7 UNKNOWN
	hierarchy := gocv.NewMat()
	contours := gocv.FindContoursWithParams(edges, &hierarchy, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	simplified_contours := gocv.NewPointsVector()

	for i := 0; i < contours.Size(); i++ {
		cnt := contours.At(i)
		hull := gocv.NewMat()
		gocv.ConvexHull(cnt, &hull, false, true)
		poly := gocv.ApproxPolyDP(gocv.NewPointVectorFromMat(hull), (0.001 * gocv.ArcLength(gocv.NewPointVectorFromMat(hull), true)), true)

		if poly.Size() != 4 {
			continue
		}

		simplified_contours.Append(poly)
	}
	if debug {
		save_image_state("state03_edges", edges, debug)
	}
	if prog != nil {
		prog <- 30
	}

	// Find bigges contours
	min_area := float64(image_size[0] * image_size[1])
	max_area := 0.0
	approx_contour := gocv.NewPointVector()

	for n := 0; n < simplified_contours.Size(); n++ {
		cnt := simplified_contours.At(n)
		area := gocv.ContourArea(cnt)

		if area > float64(min_area/10) {
			peri := gocv.ArcLength(cnt, true)
			approx := gocv.ApproxPolyDP(cnt, 0.02*peri, true)
			if area > max_area && approx.Size() == 4 {
				max_area = area
				approx_contour = approx
			}
		}
	}
	gocv.DrawContours(&img, simplified_contours, -1, color.RGBA{255, 0, 0, 0}, 1)
	if debug {
		save_image_state("state04_contours", img, debug)
	}
	if prog != nil {
		prog <- 40
	}

	// Four point transform
	// ## Find the exact (x,y) coordinates of the biggest contour and crop it out
	if approx_contour.IsNil() || approx_contour.Size() != 4 {
		fmt.Println("Conditions for four point transform not fulfilled")
		return errors.New("ERROR: Did not finish processing")
	}
	// obtain a consistent order of the points and unpack them individually
	/*
			initialzie a list of coordinates that will be ordered
		    such that the first entry in the list is the top-left,
		    the second entry is the top-right, the third is the
		    bottom-right, and the fourth is the bottom-left
	*/
	pts := approx_contour.ToPoints()
	rec_arr := []image.Point{
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
	}

	// Get Rectangle Points
	arg_min_sum := math.MaxInt
	arg_min_sum_ind := -1
	arg_max_sum := math.MinInt
	arg_max_sum_ind := -1
	arg_min_dif := math.MaxInt
	arg_min_dif_ind := -1
	arg_max_dif := math.MinInt
	arg_max_dif_ind := -1

	for i := 0; i < approx_contour.Size(); i++ {
		pnt := approx_contour.At(i)

		// the top-left point will have the smallest sum, whereas
		// the bottom-right point will have the largest sum
		sum := pnt.X + pnt.Y
		if sum < arg_min_sum {
			arg_min_sum = sum
			arg_min_sum_ind = i
		}
		if sum > arg_max_sum {
			arg_max_sum = sum
			arg_max_sum_ind = i
		}

		// now, compute the difference between the points, the
		// top-right point will have the smallest difference,
		// whereas the bottom-left will have the largest difference
		dif := pnt.Y - pnt.X
		if dif < arg_min_dif {
			arg_min_dif = dif
			arg_min_dif_ind = i
		}
		if dif > arg_max_dif {
			arg_max_dif = dif
			arg_max_dif_ind = i
		}
	}
	if prog != nil {
		prog <- 50
	}
	/*
		fmt.Printf(
			"pts:\n\tSize: %d\nargs:\n\tmin_sum: %d\n\tmax_sum: %d\n\tmin_dif: %d\n\tmax_dif: %d\n\t\n",
			len(pts),
			arg_min_sum_ind,
			arg_max_sum_ind,
			arg_min_dif_ind,
			arg_max_dif_ind,
		)
	*/
	rec_arr[0] = pts[arg_min_sum_ind]
	rec_arr[2] = pts[arg_max_sum_ind]
	rec_arr[1] = pts[arg_min_dif_ind]
	rec_arr[3] = pts[arg_max_dif_ind]

	// return the ordered coordinates
	rect := gocv.NewPointVectorFromPoints(rec_arr)

	//(tl, tr, br, bl) := rect @LEANDER
	tl := rec_arr[0]
	tr := rec_arr[1]
	br := rec_arr[2]
	bl := rec_arr[3]

	// compute the width of the new image, which will be the
	// maximum distance between bottom-right and bottom-left
	// x-coordiates or the top-right and top-left x-coordinates
	widthA := math.Sqrt(math.Pow(float64(br.X-bl.X), 2) + math.Pow(float64(br.Y-bl.Y), 2))
	widthB := math.Sqrt(math.Pow(float64(tr.X-tl.X), 2) + math.Pow(float64(tr.Y-tl.Y), 2))
	maxWidth := max(int(widthA), int(widthB))

	// compute the height of the new image, which will be the
	// maximum distance between the top-right and bottom-right
	// y-coordinates or the top-left and bottom-left y-coordinates
	heightA := math.Sqrt(math.Pow(float64(tr.X-br.X), 2) + math.Pow(float64(tr.Y-br.Y), 2))
	heightB := math.Sqrt(math.Pow(float64(tl.X-bl.X), 2) + math.Pow(float64(tl.Y-bl.Y), 2))
	maxHeight := max(int(heightA), int(heightB))

	/*
		now that we have the dimensions of the new image, construct
		the set of destination points to obtain a "birds eye view",
		(i.e. top-down view) of the image, again specifying points
		in the top-left, top-right, bottom-right, and bottom-left
		order
	*/
	dst := gocv.NewPointVectorFromPoints([]image.Point{
		{0, 0},
		{maxWidth - 1, 0},
		{maxWidth - 1, maxHeight - 1},
		{0, maxHeight - 1},
	})

	if prog != nil {
		prog <- 60
	}

	// compute the perspective transform matrix and then apply it
	M := gocv.GetPerspectiveTransform(rect, dst)
	gocv.WarpPerspective(img, &img, M, image.Point{maxWidth, maxHeight})

	if debug {
		save_image_state("state05_perspective_transform", img, debug)
	}
	if prog != nil {
		prog <- 70
	}

	// Create our shapening kernel, it must equal to one eventually
	kernel_sharpening := gocv.NewMatWithSize(3, 3, gocv.MatTypeCV32S) // @LEANDER
	kernel_sharpening.SetIntAt(0, 0, 0)
	kernel_sharpening.SetIntAt(0, 1, -1)
	kernel_sharpening.SetIntAt(0, 2, 0)
	kernel_sharpening.SetIntAt(1, 0, -1)
	kernel_sharpening.SetIntAt(1, 1, 5)
	kernel_sharpening.SetIntAt(1, 2, -1)
	kernel_sharpening.SetIntAt(2, 0, 0)
	kernel_sharpening.SetIntAt(2, 1, -1)
	kernel_sharpening.SetIntAt(2, 2, 0)
	if prog != nil {
		prog <- 80
	}
	// applying the sharpening kernel to the input image & displaying it.
	gocv.Filter2D(img, &img, -1, kernel_sharpening, image.Point{-1, -1}, 0, gocv.BorderDefault)
	if prog != nil {
		prog <- 90
	}
	if debug {
		save_image_state("state06_final", img, debug)
	}

	return nil
}
