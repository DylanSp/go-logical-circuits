package bits

func IsBitSet(value uint32, bitNum int) bool {
	bitmask := uint32(1) << uint32(bitNum)
	return (value & bitmask) != 0
}
