package opensimplex

import (
	"math/big"
	"os"
	"strconv"
	"testing"
)

func TestNoise2D(t *testing.T) {
	var noiseList string = ""
	var noizu = NewNoise(*big.NewInt(2423423555))
	for i := 0; i <= 100; i++ {
		for j := 0; j <= 50; j++ {

			var noise = noizu.Noise2D(float64(i)/100.0, float64(j)/50.0)
			noiseList += strconv.FormatFloat(noise, 'f', -1, 64) + "\n"

		}
	}

	f, err := os.Create("./noises.txt")
	if err != nil {
		panic(err)
	}

	f.WriteString(noiseList)
	f.Sync()
}
