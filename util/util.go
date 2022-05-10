package util

import (
	"math"
	"sync"

	"github.com/google/go-cmp/cmp"
)

type Epsilon struct {
	value float64
	lock  sync.Mutex
}

var DefaultEpsilon = 0.00001
var EPSILON = Epsilon{
	value: DefaultEpsilon,
}

func FloatEqual(a, b float64) bool {
	if math.Abs(a-b) < EPSILON.value {
		return true
	}

	return false
}

func SetEpsilon(num float64) {
	EPSILON.lock.Lock()
	defer EPSILON.lock.Unlock()
	EPSILON.value = num

}

var FloatComparer = cmp.Comparer(FloatEqual)

func IsNearlyEqualZero(num float64) bool {
	if math.Abs(num) < EPSILON.value {
		return true
	}

	return false
}

var Inf = math.Inf(1)
