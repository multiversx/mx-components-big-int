package twoscomplement

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigIntToTwosBytesOf(t *testing.T) {
	assertBigIntToTwosBytesOk(t, "0", []byte{})
	assertBigIntToTwosBytesOk(t, "1", []byte{0x01})
	assertBigIntToTwosBytesOk(t, "-1", []byte{0xFF})
	assertBigIntToTwosBytesOk(t, "-2", []byte{0xFE})
	assertBigIntToTwosBytesOk(t, "255", []byte{0x00, 0xFF})
	assertBigIntToTwosBytesOk(t, "256", []byte{0x01, 0x00})
	assertBigIntToTwosBytesOk(t, "-255", []byte{0xFF, 0x01})
	assertBigIntToTwosBytesOk(t, "-256", []byte{0xFF, 0x00})
	assertBigIntToTwosBytesOk(t, "-257", []byte{0xFE, 0xFF})
}

func assertBigIntToTwosBytesOk(t *testing.T, input string, expected []byte) {
	inputBi := big.NewInt(0)
	_ = inputBi.UnmarshalText([]byte(input))

	result := BigIntToTwosBytes(inputBi)
	assert.True(t, bytes.Equal(result, expected), "BigIntToTwosBytes returned wrong result. Want: %v. Have: %v", expected, result)
}

func TestBigIntToTwosBytesOfLength(t *testing.T) {
	assertBigIntToTwosBytesOfLengthOk(t, "0", 0, []byte{})
	assertBigIntToTwosBytesOfLengthOk(t, "0", 1, []byte{0x00})
	assertBigIntToTwosBytesOfLengthOk(t, "1", 1, []byte{0x01})
	assertBigIntToTwosBytesOfLengthOk(t, "-1", 1, []byte{0xFF})
	assertBigIntToTwosBytesOfLengthOk(t, "0", 3, []byte{0x00, 0x00, 0x00})
	assertBigIntToTwosBytesOfLengthOk(t, "1", 3, []byte{0x00, 0x00, 0x01})
	assertBigIntToTwosBytesOfLengthOk(t, "-1", 3, []byte{0xFF, 0xFF, 0xFF})
}

func assertBigIntToTwosBytesOfLengthOk(t *testing.T, input string, length int, expected []byte) {
	inputBi := big.NewInt(0)
	_ = inputBi.UnmarshalText([]byte(input))

	result := BigIntToTwosBytesOfLength(inputBi, length)
	assert.True(t, bytes.Equal(result, expected), "BigIntToTwosBytesOfLength returned wrong result")
}
