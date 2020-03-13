package prand

import (
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// Generates a pseudo-random string of length `n`.
func StringN(n int) string {
	var rsrc = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	for i, cache, remain := n-1, rsrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rsrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// Generates a pseudo-random int, where 0 <= x < `n`.
func IntN(n int) int {
	seed := rand.NewSource(time.Now().UnixNano())
	rnew := rand.New(seed)
	return rnew.Intn(n)
}

// Generates a pseudo-random float, where 0.0 <= x < `n`.
func FloatN(n float64) float64 {
	seed := rand.NewSource(time.Now().UnixNano())
	rnew := rand.New(seed)
	return rnew.Float64() * n
}

// Generate a slice containing pseudo-random strings. Params are slice length and
// string length.
func StringSlice(sliceLen, strLen int) (out []string) {
	for i := 0; i < sliceLen; i++ {
		out = append(out, StringN(strLen))
	}

	return
}

// Same as `StringSlice`, except the slice is joined with "," into a single
// string at the end.
func StringList(sLen, strLen int) string {
	s := StringSlice(sLen, strLen)
	return strings.Join(s, ",")
}

const f50 = float64(50)

// Heads = true
func CoinFlip() bool {
	return f50 < FloatN(100)
}

// Provide desired win % for heads (heads = true). CoinFlipBias(51) would imply
// heads should win 51% of the time.
func CoinFlipBias(bias float64) bool {
	return f50-(bias-f50) < FloatN(100)
}

// Uses `StringSlice` to generate a slice of strings. Then the bias parameter
// is passed to `CoinFlipBias` for each string in the set to produce a child
// slice containing a subset of the parent. The strings in the child slice
// retain their order from the parent.
func SliceSubsets(sLen, strLen int, bias float64) (parent []string, child []string) {
	parent = StringSlice(sLen, strLen)

	for _, rs := range parent {
		if CoinFlipBias(bias) {
			child = append(child, rs)
		}
	}

	return
}

// The same as `SliceSubsets`, except the parent and child slices are each
// joined with "," into strings before being returned.
func SliceSubsetsList(sLen, strLen int, bias float64) (parent string, child string) {
	p, c := SliceSubsets(sLen, strLen, bias)
	parent = strings.Join(p, ",")
	child = strings.Join(c, ",")

	return
}
