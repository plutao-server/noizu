package opensimplex

import (
	"math"
)

const STRETCH float64 = -0.21_132_486_540_518_713
const SQUISH float64 = 0.366_025_403_784_439
const ROOT2OVER2 float64 = 0.7_071_067_811_865_476

const PRIME_X int64 = 0x5205402B9270C86F
const PRIME_Y int64 = 0x598CD327003817B5
const HASH_MULTIPLIER int64 = 0x53A3F72DEEC546F5
const RSQUARED float64 = 2.0 / 3.0

const N_GRADS_EXPONENT int32 = 7
const N_GRADS int32 = 1 << N_GRADS_EXPONENT

const GRADS_NORMALIZER float64 = 0.05_481_866_495_625_118

var GRADIENT_TABLE [16]int8 = [16]int8{
	5, 2, 2, 5,
	-5, 2, -2, 5,
	5, -2, 2, -5,
	-5, -2, -2, -5,
}
var GRADIENTS_TABLE [N_GRADS * 2]float64 = [N_GRADS * 2]float64{}

type PermsArray = [256]int64

type Noizu struct {
	seed       int64
	permsArray PermsArray
}

func NewNoise(seed int64) *Noizu {
	var grads2 = []float64{
		0.38268343236509, 0.923879532511287,
		0.923879532511287, 0.38268343236509,
		0.923879532511287, -0.38268343236509,
		0.38268343236509, -0.923879532511287,
		-0.38268343236509, -0.923879532511287,
		-0.923879532511287, -0.38268343236509,
		-0.923879532511287, 0.38268343236509,
		-0.38268343236509, 0.923879532511287,
		//-------------------------------------//
		0.130526192220052, 0.99144486137381,
		0.608761429008721, 0.793353340291235,
		0.793353340291235, 0.608761429008721,
		0.99144486137381, 0.130526192220051,
		0.99144486137381, -0.130526192220051,
		0.793353340291235, -0.60876142900872,
		0.608761429008721, -0.793353340291235,
		0.130526192220052, -0.99144486137381,
		-0.130526192220052, -0.99144486137381,
		-0.608761429008721, -0.793353340291235,
		-0.793353340291235, -0.608761429008721,
		-0.99144486137381, -0.130526192220052,
		-0.99144486137381, 0.130526192220051,
		-0.793353340291235, 0.608761429008721,
		-0.608761429008721, 0.793353340291235,
		-0.130526192220052, 0.99144486137381,
	}
	for i := range grads2 {
		grads2[i] = float64(grads2[i] / GRADS_NORMALIZER)
	}
	var j = 0
	for i := range GRADIENTS_TABLE {
		if j == len(grads2) {
			j = 0
		}
		GRADIENTS_TABLE[i] = grads2[j]
		j++
	}

	return &Noizu{
		seed:       seed,
		permsArray: generatePermutationsArray(seed),
	}
}

func (noizu *Noizu) Noise2D(x, y float64) float64 {
	var xx = x * ROOT2OVER2
	var yy = y * (ROOT2OVER2 * (1 + 2*SQUISH))

	return noise2DBase(noizu.seed, xx+yy, yy-xx) //change 47 to smth else its a normalizing scalar
}
func noise2DBase(seed int64, xs, ys float64) float64 { //xs,ys = x skewed, y squished
	var xsb = int64(math.Floor(xs))
	var ysb = int64(math.Floor(ys))
	var xi = float64(xs - float64(xsb))
	var yi = float64(ys - float64(ysb))

	var xsbp = xsb * PRIME_X
	var ysbp = ysb * PRIME_Y

	var t = (xi + yi) * STRETCH
	var dx0 = xi + t
	var dy0 = yi + t

	var a0 = RSQUARED - dx0*dx0 - dy0*dy0
	var value = (a0 * a0) * (a0 * a0) * grad(seed, xsbp, ysbp, dx0, dy0)

	// Second vertex.
	var a1 = float64(2*(1+2*STRETCH)*(1/STRETCH+2))*t + (float64(-2*(1+2*STRETCH)*(1+2*STRETCH)) + a0)
	var dx1 = dx0 - float64(1+2*STRETCH)
	var dy1 = dy0 - float64(1+2*STRETCH)
	value += (a1 * a1) * (a1 * a1) * grad(seed, xsbp+PRIME_X, ysbp+PRIME_Y, dx1, dy1)

	var xmyi = xi - yi
	if t < STRETCH {
		if xi+xmyi > 1 {
			var dx2 = dx0 - float64(3*STRETCH+2)
			var dy2 = dy0 - float64(3*STRETCH+1)
			var a2 = RSQUARED - dx2*dx2 - dy2*dy2
			if a2 > 0 {
				value += (a2 * a2) * (a2 * a2) * grad(seed, xsbp+(2*(PRIME_X-math.MaxInt64)+math.MaxInt64), ysbp+PRIME_Y, dx2, dy2)
			}
		} else {
			var dx2 = dx0 - float64(STRETCH)
			var dy2 = dy0 - float64(STRETCH+1)
			var a2 = RSQUARED - dx2*dx2 - dy2*dy2
			if a2 > 0 {
				value += (a2 * a2) * (a2 * a2) * grad(seed, xsbp, ysbp+PRIME_Y, dx2, dy2)
			}
		}

		if yi-xmyi > 1 {
			var dx3 = dx0 - float64(3*STRETCH+1)
			var dy3 = dy0 - float64(3*STRETCH+2)
			var a3 = RSQUARED - dx3*dx3 - dy3*dy3
			if a3 > 0 {
				value += (a3 * a3) * (a3 * a3) * grad(seed, xsbp+PRIME_X, ysbp+(2*(PRIME_Y-math.MaxInt64)+math.MaxInt64), dx3, dy3)
			}
		} else {
			var dx3 = dx0 - float64(STRETCH+1)
			var dy3 = dy0 - float64(STRETCH)
			var a3 = RSQUARED - dx3*dx3 - dy3*dy3
			if a3 > 0 {
				value += (a3 * a3) * (a3 * a3) * grad(seed, xsbp+PRIME_X, ysbp, dx3, dy3)
			}
		}
	} else {
		if xi+xmyi < 0 {
			var dx2 = dx0 + float64(1+STRETCH)
			var dy2 = dy0 + float64(STRETCH)
			var a2 = RSQUARED - dx2*dx2 - dy2*dy2
			if a2 > 0 {
				value += (a2 * a2) * (a2 * a2) * grad(seed, xsbp-PRIME_X, ysbp, dx2, dy2)
			}
		} else {
			var dx2 = dx0 - float64(STRETCH+1)
			var dy2 = dy0 - float64(STRETCH)
			var a2 = RSQUARED - dx2*dx2 - dy2*dy2
			if a2 > 0 {
				value += (a2 * a2) * (a2 * a2) * grad(seed, xsbp+PRIME_X, ysbp, dx2, dy2)
			}
		}

		if yi < xmyi {
			var dx2 = dx0 + float64(STRETCH)
			var dy2 = dy0 + float64(STRETCH+1)
			var a2 = RSQUARED - dx2*dx2 - dy2*dy2
			if a2 > 0 {
				value += (a2 * a2) * (a2 * a2) * grad(seed, xsbp, ysbp-PRIME_Y, dx2, dy2)
			}
		} else {
			var dx2 = dx0 - float64(STRETCH)
			var dy2 = dy0 + float64(STRETCH+1)
			var a2 = RSQUARED - dx2*dx2 - dy2*dy2
			if a2 > 0 {
				value += (a2 * a2) * (a2 * a2) * grad(seed, xsbp, ysbp+PRIME_Y, dx2, dy2)
			}
		}
	}
	return value
}
func grad(seed, xsvp, ysvp int64, dx, dy float64) float64 {
	var hash = seed ^ xsvp ^ ysvp
	hash *= HASH_MULTIPLIER
	hash ^= hash >> (64 - N_GRADS_EXPONENT + 1)
	var gi = int32(hash) & ((N_GRADS - 1) << 1)
	return GRADIENTS_TABLE[gi]*dx + GRADIENTS_TABLE[gi|1]*dy
}
func evaluateInsideTriangle(insX, insY, originX, originY, gridX, gridY float64, permsArray PermsArray) float64 {
	var ixpy = insX + insY
	var xsv_ext = 1.0
	var ysv_ext = 1.0
	var dx_ext = 1.0
	var dy_ext = 1.0
	if ixpy <= 1.0 {
		var zins = 1.0 - ixpy
		if zins > insX || zins > insY {
			if insX > insY {
				xsv_ext = gridX + 1.0
				ysv_ext = gridY - 1.0
				dx_ext = originX - 1.0
				dy_ext = originY + 1.0
			} else {
				xsv_ext = gridX - 1.0
				ysv_ext = gridY + 1.0
				dx_ext = originX + 1.0
				dy_ext = originY - 1.0
			}
		} else {
			xsv_ext = gridX + 1.0
			ysv_ext = gridY + 1.0
			dy_ext = originX - 1 - 2*SQUISH
			dx_ext = originY - 1 - 2*SQUISH
		}
	} else {
		var zins = 2 - ixpy
		if zins < insX || zins < insY {
			if insX > insY {
				xsv_ext = gridX + 2.0
				ysv_ext = gridY + 0.0
				dx_ext = originX - 2 - 2*SQUISH
				dy_ext = originY + 0 - 2*SQUISH
			} else {
				xsv_ext = gridX + 0.0
				ysv_ext = gridY + 2.0
				dx_ext = originX + 0 - 2*SQUISH
				dy_ext = originY - 2 - 2*SQUISH
			}
		} else {
			dx_ext = originX
			dy_ext = originY
			xsv_ext = gridX
			ysv_ext = gridY
		}
		gridX += 1
		gridY += 1
		originX = originX - 1 - 2*SQUISH
		originY = originY - 1 - 2*SQUISH

	}
	var attn = 2 - originX*originX - originY*originY
	var value = 0.0
	if attn > 0 {
		value += math.Pow(attn, 4.0) * extrapolate(gridX, gridY, originX, originY, permsArray)
	}
	var attn_ext = 2 - dx_ext*dx_ext - dy_ext*dy_ext
	if attn_ext > 0 {
		value += math.Pow(attn, 4.0) * extrapolate(xsv_ext, ysv_ext, dx_ext, dy_ext, permsArray)
	}

	return value
}
func contribute(deltaX, deltaY, originX, originY, gridX, gridY float64, permsArray PermsArray) float64 {
	// var dxpy = deltaX + deltaY
	var shiftedX = originX - deltaX - SQUISH
	var shiftedY = originY - deltaY - SQUISH
	var attenuation = 2.0 - getAttenuationFactor(shiftedX, shiftedY)

	if attenuation > 0.0 {
		return math.Pow(attenuation, 4.0) * extrapolate(gridX+deltaX, gridY+deltaY, shiftedX, shiftedY, permsArray) // multiply with extrapolation
	} else {
		return 0.0
	}

}
func extrapolate(gridX, gridY, deltaX, deltaY float64, permsArray PermsArray) float64 {
	var index = getGradientTableIndex(gridX, gridY, permsArray)
	var pointX = float64(GRADIENT_TABLE[index])
	var pointY = float64(GRADIENT_TABLE[index+1])
	return pointX*deltaX + pointY*deltaY
}
func getAttenuationFactor(x, y float64) float64 {
	return (x * x) + (y * y)
}
func getGradientTableIndex(gridX, gridY float64, permsArray PermsArray) int64 {
	var index = (int32(permsArray[int32(gridX)&0xFF]) + int32(gridY)) & 0xFF
	return (permsArray[index] & 0x0E)
}
func generatePermutationsArray(seed int64) PermsArray {
	var permsArray = PermsArray{}
	var source = PermsArray{}
	for i := range source {
		source[i] = int64(i)
	}
	seed = seed*6_364_136_223_846_793_005 + 1_442_695_040_888_963_407
	seed = seed*6_364_136_223_846_793_005 + 1_442_695_040_888_963_407
	seed = seed*6_364_136_223_846_793_005 + 1_442_695_040_888_963_407
	for i := int32(len(permsArray) - 1); i > -1; i-- {
		seed = seed*6_364_136_223_846_793_005 + 1_442_695_040_888_963_407
		var r = int32((seed + 32) % int64(i+1))

		if r < 0 {
			r += i + 1
		}
		permsArray[i] = source[r]
		source[r] = source[i]

	}

	return permsArray

}
