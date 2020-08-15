package quickrand

import (
	"math/rand"
	"time"
	"unsafe"
)

// package to make some of the rng requirements for the game easier, e.g distinct mode

// gets n random int values between 0 and max range efficiently w/ distinct option
func RandInts(max, n int, distinct bool) *[]int {
	vals := make([]int, n)
	if distinct {
		// do this so zero values in distinct mode are possible
		for i := range vals {
			vals[i] = -1
		}
	}
	if max <= 0 || n <= 0 || (distinct && n > max+1) {
		panic("Non positive argument to RandIntRange or distinct is impossibly true.")
	}
	src := rand.NewSource(time.Now().UnixNano())
	max_bitrep := bitrep(max)
	val_per_rand := 63 / max_bitrep // 63 / max_bitrep == num of possible value fills per int63 call
	var mask int64 = (1 << max_bitrep) - 1
	// can improve by recycling leftover bits that couldnt represent full value, but probably not worth it
	for bits, rem, i := src.Int63(), val_per_rand, 0; i < n; {
		if rem == 0 {
			bits = src.Int63() // refill random bits and reset rem
			rem = val_per_rand
		}
		val := bits & mask
		if val <= int64(max) && (!distinct || !InVals(int(val), &vals)) {
			vals[i] = int(val)
			i++
		}
		rem--
		bits >>= max_bitrep
	}
	return &vals
}

// exported as it is useful for some other functions in
// different packages
func InVals(val int, vals *[]int) bool {
	for i := range *vals {
		if val == (*vals)[i] {
			return true
		}
	}
	return false
}

// returns number of bits required to represent max
func bitrep(max int) int {
	size := unsafe.Sizeof(max)
	var mask int = 1 << (size*8 - 2)
	for i := int(size*8 - 1); i > 0; i-- {
		if max&mask != 0 {
			return i
		}
		mask >>= 1
	}
	return 0
}
