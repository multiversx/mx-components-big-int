package twoscomplement

// CopyAlignRight yields a copy of the bytes, of specific length.
// If output is longer than input, it will pad left with zeroes.
// If output if shorter than input, input will be trimmed to the left.
// The right-most bytes always remain in place.
func CopyAlignRight(input []byte, targetLength int) []byte {
	offset := len(input) - targetLength
	resultBytes := make([]byte, targetLength)
	for i := 0; i < targetLength; i++ {
		j := offset + i
		if j < 0 {
			resultBytes[i] = 0 // pad left with 00000000
		} else {
			resultBytes[i] = input[j]
		}
	}
	return resultBytes
}
