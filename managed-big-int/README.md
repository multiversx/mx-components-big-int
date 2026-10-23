# managed-big-int

> **Warning: experimental.** This package is not used anywhere and has never been tested in production. Do not rely on it for production code without extensive testing of your own.

Package `managedbigint`: a pointer-free arena for `math/big` integers.

## Problem

A program holding many `*big.Int` pays for one heap object per number (plus its word slice): more allocations and more work for the garbage collector, and the state is spread over the heap. This package stores the words of all numbers in **one** `[]big.Word` and refers to each number with a small value-type handle. The whole state is a flat slice and has nothing for the GC to chase.

This fits a runtime that manages thousands of numbers on behalf of untrusted code (for example the managed types of a VM).

## How it works

```go
c := managedbigint.NewBigIntContainer()

x := c.Insert(big.NewInt(40))          // copies the number into the arena
y := c.InsertUint64(2)
z := c.Add(managedbigint.Zero, x, y)   // result goes to a new slot
z  = c.Mul(z, z, y)                    // reuses z's slot if the result fits
fmt.Println(c.Get(z))                  // 84 (a copy, as *big.Int)
```

- `BigIntHandle` holds `start`, `length`, `capacity` (offsets into the arena) and the sign. `Zero` is a ready-made handle for 0 that uses no storage.
- Every operation loads its operands into reusable scratch `big.Int` registers by pointing them at the arena (`SetBits`) with a capacity-limited slice, so `math/big` can never grow into a neighbour's words. No `big.Int` is allocated per call.
- The result is computed into the destination's slot. If it still fits, the handle's length and sign are refreshed. If `math/big` had to reallocate, the result is appended to the arena as a new slot.

### Rules for callers

- **Always use the returned handle:** `z = c.Add(z, x, y)`. The handle can change (new slot, new length, new sign).
- **Handles are values, not pointers.** After `z` is modified in place, an old copy of its handle is stale (length and sign). Keep one handle per number.
- `dest` may alias operands (`c.Sub(z, z, y)`), and may be `Zero`.
- `Get` returns a copy. `GetUnsafe` returns a `*big.Int` that is a view into the arena: writing to it changes the stored number, and it is invalidated by later operations on that slot.
- Semantics are those of `math/big`: `Quo`/`Rem` truncate, `Div`/`Mod` are Euclidean, division by zero and `Sqrt` of a negative panic.

### Limitations

- **Memory is never reclaimed.** The arena only grows. A result that outgrows its slot abandons the old slot. There is no free or compaction.
- Not safe for concurrent use (the scratch registers are shared).
- Handles are not validated: a handle from another container, or a forged one, can panic or read other numbers.

## API overview

| Group | Methods |
|---|---|
| Storage | `Insert`, `InsertUint64`, `Set`, `Update`, `Get`, `GetUnsafe` |
| Conversion | `SetInt64`, `ToInt64`, `IsInt64`, `SetBytes`, `GetBytes`, `BitLen`, `ByteLen` |
| Arithmetic | `Add`, `Sub`, `Mul`, `Quo`, `Rem`, `Div`, `Mod`, `Abs`, `Neg`, `Sqrt`, `Cmp`, `CmpAbs` |
| Bitwise | `And`, `AndNot`, `Or`, `Xor`, `Not` |

## Tests

`container_test.go` checks every operation against plain `math/big` on boundary values (word and byte edges, both signs), covers in-place shrink/grow/sign-flip, aliasing and neighbour integrity, and runs a seeded random program on the container and on `big.Int` in lockstep, comparing all numbers after every step.
