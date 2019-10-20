package managedbigint

import "math/big"

// BigIntContainer is a structure holding the data for big.Int numbers without pointers.
type BigIntContainer struct {
	data []big.Word

	register1   *big.Int
	register2   *big.Int
	destination *big.Int
}

// NewBigIntContainer constructs a new BigIntContainer.
func NewBigIntContainer() *BigIntContainer {
	return &BigIntContainer{
		data:        nil,
		register1:   big.NewInt(0),
		register2:   big.NewInt(0),
		destination: big.NewInt(0),
	}
}

// BigIntHandle acts like a pointer to a big.Int value in a BigIntContainer.
type BigIntHandle struct {
	start    int
	length   int
	capacity int
	negative bool
}

// Zero is a reference to Zero value. Zero doesn't need storage.
var Zero = BigIntHandle{
	start:    0,
	length:   0,
	capacity: 0,
	negative: false,
}

// Insert adds a copy of a big number into the BigIntContainer.
func (c *BigIntContainer) Insert(bi *big.Int) BigIntHandle {
	if bi.Sign() == 0 {
		return Zero
	}
	words := bi.Bits()
	start := len(c.data)
	c.data = append(c.data, words[:cap(words)]...) // copy full capacity, to allow later extension
	return BigIntHandle{
		start:    start,
		length:   len(words),
		capacity: cap(words),
		negative: bi.Sign() < 0,
	}
}

// InsertUint64 adds a uint64 number into the BigIntContainer.
func (c *BigIntContainer) InsertUint64(x uint64) BigIntHandle {
	bi := big.NewInt(0).SetUint64(x)
	return c.Insert(bi)
}

func (c *BigIntContainer) loadBigInt(handler BigIntHandle, target *big.Int) {
	// setting the capacity is very important
	// the math/big library will sometimes try to extend the slice, but not beyond its capacity
	// if we do not specify slice capacity, the default capacity might extend over other number data
	// potentially causing hard to detect bugs
	target.SetBits(c.data[handler.start : handler.start+handler.length : handler.start+handler.capacity])
	if handler.negative {
		target.Neg(target)
	}
}

// Get yields a copy of a BigIntContainer number, as big.Int.
func (c *BigIntContainer) Get(handler BigIntHandle) *big.Int {
	result := big.NewInt(0)
	c.loadBigInt(handler, result)
	return big.NewInt(0).Set(result) // clone, to prevent accidental changing of underlying structure
}

// GetUnsafe casts a BigIntContainer number to big.Int.
// Changing the resulting big.Int will also change the underlying data.
func (c *BigIntContainer) GetUnsafe(handler BigIntHandle) *big.Int {
	result := big.NewInt(0)
	c.loadBigInt(handler, result)
	return result
}
