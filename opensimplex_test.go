package opensimplex

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestNoise2D(t *testing.T) {
	var noizu = NewNoise(2342342)
	var HEIGHT = 1024
	var WIDTH = 1024
	img := image.NewAlpha(image.Rect(0, 0, WIDTH, HEIGHT))
	for i := 0; i <= WIDTH; i++ {
		for j := 0; j <= HEIGHT; j++ {

			var noise = noizu.Noise2D(float64(i)/float64(WIDTH), float64(j)/float64(HEIGHT))
			img.SetAlpha(i, j, color.Alpha{A: uint8(noise * 127.0)})

		}
	}

	f, err := os.Create("./noise.png")
	if err != nil {
		panic(err)
	}

	png.Encode(f, img)

}
