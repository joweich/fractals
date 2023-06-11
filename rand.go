package main

func RandUint64(rng *uint64) uint64 {
	*rng = *rng*0x3243f6a8885a308d + 1
	r := *rng
	r ^= r >> 32
	r *= 1111111111111111111
	r ^= r >> 32
	return r
}

func RandFloat64(rng *uint64) float64 {
	return float64(RandUint64(rng)/2) / (1 << 63)
}