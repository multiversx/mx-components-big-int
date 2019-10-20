package managedbigint

import "math/big"

type binaryOperation func(destination, x, y *big.Int) *big.Int

type unaryOperation func(destination, x *big.Int) *big.Int

func (c *BigIntContainer) performBinaryOperation(op binaryOperation, dest, x, y BigIntHandle) BigIntHandle {
	c.loadBigInt(x, c.register1)
	c.loadBigInt(y, c.register2)
	c.loadBigInt(dest, c.destination)

	destNegBefore := c.destination.Sign() < 0
	destDataBefore := c.destination.Bits()

	c.destination = op(c.destination, c.register1, c.register2)

	destNegAfter := c.destination.Sign() < 0
	destDataAfter := c.destination.Bits()

	if bigIntDataMoved(destDataBefore, destDataAfter) {
		return c.Insert(c.destination)
	}

	// maybe dest changed sign
	if destNegBefore != destNegAfter {
		dest.negative = !dest.negative
	}

	return dest
}

// true if big.Int got reallocated
func bigIntDataMoved(before, after []big.Word) bool {
	if len(before) == 0 {
		return len(after) != 0
	}
	if len(after) == 0 {
		return true
	}

	return &before[0] != &after[0]
}

func (c *BigIntContainer) Add(dest, x, y BigIntHandle) BigIntHandle {
	return c.performBinaryOperation((*big.Int).Add, dest, x, y)
}

func (c *BigIntContainer) Sub(dest, x, y BigIntHandle) BigIntHandle {
	return c.performBinaryOperation((*big.Int).Sub, dest, x, y)
}

func (c *BigIntContainer) Mul(dest, x, y BigIntHandle) BigIntHandle {
	return c.performBinaryOperation((*big.Int).Mul, dest, x, y)
}
