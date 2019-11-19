package twoscomplement

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToSignedBigInt(t *testing.T) {
	assertBytesToSignedBigIntOk(t, []byte{}, "0")
	assertBytesToSignedBigIntOk(t, []byte{0x00}, "0")
	assertBytesToSignedBigIntOk(t, []byte{0x01}, "1")
	assertBytesToSignedBigIntOk(t, []byte{0xFF}, "-1")
	assertBytesToSignedBigIntOk(t, []byte{0x00, 0xFF}, "255")
	assertBytesToSignedBigIntOk(t, []byte{0x01, 0x00}, "256")
	assertBytesToSignedBigIntOk(t, []byte{0xFF, 0xFF}, "-1")
	assertBytesToSignedBigIntOk(t, []byte{0xFF, 0xFE}, "-2")
	assertBytesToSignedBigIntOk(t, []byte{0xFF, 0x00}, "-256")
}

func assertBytesToSignedBigIntOk(t *testing.T, input []byte, expected string) {
	conv := BytesToSignedBigInt(input)
	expectedBi := big.NewInt(0)
	_ = expectedBi.UnmarshalText([]byte(expected))
	assert.True(t, conv.Cmp(expectedBi) == 0, "BytesToSignedBigInt returned wrong result")
}
