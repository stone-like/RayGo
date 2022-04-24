package util

import (
	"math"

	"github.com/google/go-cmp/cmp"
)

var EPSILON = 0.00001

func FloatEqual(a, b float64) bool {
	if math.Abs(a-b) < EPSILON {
		return true
	}

	return false
}

var FloatComparer = cmp.Comparer(FloatEqual)