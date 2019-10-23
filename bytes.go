package managedbigint

// SetBytes interprets buf as the bytes of a big-endian unsigned
// integer, sets dest to that value, and returns dest.
func (c *BigIntContainer) SetBytes(dest BigIntHandle, buf []byte) BigIntHandle {
	c.loadBigInt(dest, c.destination)

	destDataBefore := c.destination.Bits()

	c.destination = c.destination.SetBytes(buf)

	destDataAfter := c.destination.Bits()

	if bigIntDataMoved(destDataBefore, destDataAfter) {
		return c.Insert(c.destination)
	}

	dest.negative = false
	return dest
}

// GetBytes returns the absolute value of x as a big-endian byte slice.
func (c *BigIntContainer) GetBytes(x BigIntHandle) []byte {
	c.loadBigInt(x, c.register1)
	return c.register1.Bytes()
}
