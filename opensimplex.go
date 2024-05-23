package opensimplex

import (
	"fmt"
	"math"
	"math/big"
)

const STRETCH float64 = -0.211_324_865_405_187
const SQUISH float64 = 0.366_025_403_784_439

var GRADIENT_TABLE [8][2]float64 = [8][2]float64{
	{5.0, 2.0},
	{2.0, 5.0},
	{-5.0, 2.0},
	{-2.0, 5.0},
	{5.0, -2.0},
	{2.0, -5.0},
	{-5.0, -2.0},
	{-2.0, -5.0},
}

type PermsArray = [2048]int64

type Noizu struct {
	seed       big.Int
	permsArray PermsArray
}

func NewNoise(seed big.Int) *Noizu {
	return &Noizu{
		seed:       seed,
		permsArray: generatePermutationsArray(seed),
	}
}

func (noizu *Noizu) Noise2D(x, y float64) float64 {

	var xpy = x + y
	var stretchX = x + float64(STRETCH*xpy)
	var stretchY = y + float64(STRETCH*xpy)

	var gridX = math.Floor(stretchX)
	var gridY = math.Floor(stretchY)
	var gxpy = gridX + gridY

	var squishedX = gridX + float64(SQUISH+gxpy)
	var squishedY = gridY + float64(SQUISH+gxpy)
	var insX = stretchX - gridX
	var insY = stretchY - gridY
	var originX = x - squishedX
	var originY = y - squishedY
	println(contribute(1.0, 0.0, originX, originY, gridX, gridY, noizu.permsArray))
	var value = contribute(1.0, 0.0, originX, originY, gridX, gridY, noizu.permsArray) + contribute(0.0, 1.0, originX, originY, gridX, gridY, noizu.permsArray) + evaluateInsideTriangle(insX, insY, originX, originY, gridX, gridY, noizu.permsArray)

	return value / 47.0 //change 47 to smth else its a normalizing scalar
}
func evaluateInsideTriangle(insX, insY, originX, originY, gridX, gridY float64, permsArray PermsArray) float64 {
	var ixpy = insX + insY
	var factorPointX = 1.0
	var factorPointY = 1.0
	if ixpy <= 1.0 {
		factorPointX = 0.0
		factorPointY = 0.0
	}
	var zins = 1.0 + factorPointX - ixpy

	var pointX = 1.0 - factorPointX
	var pointY = 1.0 - factorPointY
	if zins > insX || zins > insY {
		if insX > insY {
			pointX = 1.0 + factorPointX
			pointY = -1.0 + factorPointY
		} else {
			pointX = -1.0 + factorPointX
			pointY = 1.0 + factorPointY
		}
	}
	return (contribute(0.0+factorPointX, 0.0+factorPointY, originX, originY, gridX, gridY, permsArray) + contribute(pointX, pointY, originX, originY, gridX, gridY, permsArray))
}
func contribute(deltaX, deltaY, originX, originY, gridX, gridY float64, permsArray PermsArray) float64 {
	var dxpy = deltaX + deltaY
	var shiftedX = originX - deltaX - SQUISH*dxpy
	var shiftedY = originY - deltaY - SQUISH*dxpy
	var attenuation = 2.0 - getAttenuationFactor(shiftedX, shiftedY)
	fmt.Println(attenuation)
	if attenuation > 0.0 {
		return math.Pow(attenuation, 4.0) * extrapolate(gridX+deltaX, gridY+deltaY, shiftedX, shiftedY, permsArray) // multiply with extrapolation
	} else {
		return 0.0
	}

}
func extrapolate(gridX, gridY, deltaX, deltaY float64, permsArray PermsArray) float64 {
	var point = GRADIENT_TABLE[getGradientTableIndex(gridX, gridY, permsArray)]
	var pointX = point[0]
	var pointY = point[1]
	return pointX*deltaX + pointY*deltaY
}
func getAttenuationFactor(x, y float64) float64 {
	return (x * x) + (y * y)
}
func getGradientTableIndex(gridX, gridY float64, permsArray PermsArray) int64 {
	var index = (permsArray[int64(gridX)&0xFF] + int64(gridY)) & 0xFF
	return (permsArray[index] & 0x0E) >> 1
}
func generatePermutationsArray(seed big.Int) PermsArray {
	var permsArray = PermsArray{}
	var source = PermsArray{}
	for i := range source {
		source[i] = int64(i)
	}
	seed = *seed.Add(seed.Mul(&seed, big.NewInt(6_364_136_223_846_793_005)), big.NewInt(1_442_695_040_888_963_407))
	for i := range source {
		i := int64(i)
		var r = seed.Mod(seed.Add(&seed, big.NewInt(31)), big.NewInt(i+1)).Int64()
		if r < 0 {
			r += i + 1
		}
		permsArray[i] = source[r]
		source[r] = source[i]

	}
	return permsArray

}
