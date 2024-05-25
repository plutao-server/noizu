package opensimplex

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"testing"

	"github.com/anthonynsimon/bild/blur"
)

func TestNoise2D(t *testing.T) {
	var noizu = NewNoise(224234)
	// var r = rand.New(rand.NewSource(4))
	var HEIGHT = 512
	var WIDTH = 1024

	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	var reducePeaks = 0.0
	var peakSamplesFraction = 1
	for i := 0; i <= WIDTH; i++ {

		var noise = noizu.Noise2D(float64(i/peakSamplesFraction), 0.0)
		for j := HEIGHT; j > HEIGHT-int(math.Abs(noise)*float64(HEIGHT)-reducePeaks); j-- {
			img.SetRGBA(i, j, color.RGBA{R: 130, G: 235, B: 127, A: uint8(255.0)})
		}

	}

	// for i := 0; i <= WIDTH; i++ {
	// 	for j := 0; j <= HEIGHT; j++ {

	// 		var noise = noizu.Noise2D(float64(i*4)/200, float64(j*4)/150)
	// 		if 1/(1+(math.Pow(math.E, -(math.Abs((noise)))))) <= 0.56 && r.Float64() > 0.2 {
	// 			img.SetAlpha(i, j, color.Alpha{A: uint8(255.0)})
	// 		} else {
	// 			img.SetAlpha(i, j, color.Alpha{A: uint8(0.0)})
	// 		}

	// 	}
	// }

	irg := blur.Gaussian(img, 400)
	var thresholdH = 100
	for i := 0; i <= WIDTH; i++ {
		for j := 0; j <= HEIGHT; j++ {

			if irg.RGBAAt(i, j).A >= 150 {
				img.SetRGBA(i, j-thresholdH, color.RGBA{R: 130, G: 235, B: 127, A: uint8(255.0)})
			} else {
				img.SetRGBA(i, j-thresholdH, color.RGBA{R: uint8(i) % 255, G: uint8(j) % 255, B: uint8(i+j) % 255, A: uint8(0.0)})
			}

		}
	}
	for i := 0; i <= WIDTH; i++ {
		for j := HEIGHT - thresholdH; j <= HEIGHT; j++ {
			img.SetRGBA(i, j, color.RGBA{R: 130, G: 235, B: 127, A: uint8(255.0)})
		}
	}
	f, err := os.Create("./noise.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)

}
