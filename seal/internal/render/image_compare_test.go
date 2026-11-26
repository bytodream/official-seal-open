package render

import (
	"image"
	"image/color"
	"testing"
)

func TestCompareImagesIdentical(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			img.Set(x, y, color.RGBA{200, 10, 10, 255})
		}
	}
	d := CompareImages(img, img, 0)
	if d.MismatchedPixels != 0 {
		t.Fatalf("%+v", d)
	}
}
