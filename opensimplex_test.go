package opensimplex

import (
	"math/big"
	"testing"
)

func TestNoise2D(t *testing.T) {

	var noizu = NewNoise(*big.NewInt(23429482934289328))
	var noise = noizu.Noise2D(0.6434, 0.5564)
	t.Log("Noise at 0.6434, 0.5564 is ", noise)
}
