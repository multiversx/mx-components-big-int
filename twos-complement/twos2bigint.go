package twoscomplement

import "math/big"

// BytesToSignedBigInt convert a byte array to a number
// interprets input as a 2's complement representation if the first bit (most significant) is 1
// big endian
func BytesToSignedBigInt(twosBytes []byte) *big.Int {
	if len(twosBytes) == 0 {
		return big.NewInt(0)
	}

	testBit := twosBytes[0] >> 7
	result := new(big.Int)
	if testBit == 0 {
		// positive number, no further processing required
		result.SetBytes(twosBytes)
	} else {
		// convert to negative number
		notBytes := make([]byte, len(twosBytes))
		for i, b := range twosBytes {
			notBytes[i] = ^b // negate every bit
		}
		result.SetBytes(notBytes)
		result.Neg(result)
		result.Sub(result, bigOne) // -1
	}

	return result
}
