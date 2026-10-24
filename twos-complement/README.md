# twos-complement

> Used in production by the MultiversX Go VM: [mx-chain-vm-go](https://github.com/multiversx/mx-chain-vm-go).

Package `twoscomplement`: big-endian two's complement encoding of `*big.Int`.

`math/big` only offers sign-magnitude bytes (`Bytes()` / `SetBytes()` ignore the sign). This package adds the signed encoding, with the *minimal* number of bytes.

## Encoding

- `0` is the empty byte string.
- Positive values are the magnitude, with a leading `0x00` if the top bit would otherwise be 1.
- Negative values are the bitwise negation of `|x|-1`, with a leading `0xFF` if the top bit would otherwise be 0.

| Value | Bytes |
|---|---|
| `0` | `` (empty) |
| `127` | `7f` |
| `128` | `00 80` |
| `-1` | `ff` |
| `-128` | `80` |
| `-129` | `ff 7f` |

## API

| Function | Description |
|---|---|
| `ToBytes(bi)` | Minimal two's complement bytes. Does not alter the input. |
| `ToBytesOfLength(bi, n)` | Exactly `n` bytes, sign-extended; error if the value does not fit. |
| `FromBytes(b)` | Decode two's complement bytes to a new `*big.Int`. Empty input is 0. |
| `SetBytes(z, b)` | Same as `FromBytes`, but stores into `z`. |
| `ByteLength(bi)` | `len(ToBytes(bi))`, computed from the bit length, no allocation. |
| `UnsignedByteLength(bi)` | `len(bi.Bytes())` (sign disregarded), no allocation. |
| `CopyAlignRight(b, n)` | Copy of `b` as `n` bytes, left-padded with zeroes or left-trimmed. |

`ByteLength` and `UnsignedByteLength` treat `nil` as 0 bytes.

```go
b := twoscomplement.ToBytes(big.NewInt(-129)) // ff 7f
x := twoscomplement.FromBytes(b)              // -129
n := twoscomplement.ByteLength(x)             // 2
```
