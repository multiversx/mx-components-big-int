package managedbigint

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func bi(s string) *big.Int {
	x, ok := new(big.Int).SetString(s, 0)
	if !ok {
		panic("bad big int literal: " + s)
	}
	return x
}

// interesting values around word and byte boundaries, both signs
func boundaryValues() []*big.Int {
	var values []*big.Int
	for _, shift := range []uint{0, 1, 7, 8, 31, 32, 63, 64, 65, 127, 128, 129, 191, 192, 256, 300} {
		p := new(big.Int).Lsh(big.NewInt(1), shift)
		for _, delta := range []int64{-1, 0, 1} {
			v := new(big.Int).Add(p, big.NewInt(delta))
			values = append(values, v, new(big.Int).Neg(v))
		}
	}
	return values
}

func randomBigInt(r *rand.Rand, maxBits int) *big.Int {
	bits := r.Intn(maxBits + 1)
	x := new(big.Int)
	if bits > 0 {
		x.Rand(r, new(big.Int).Lsh(big.NewInt(1), uint(bits)))
	}
	if r.Intn(2) == 0 {
		x.Neg(x)
	}
	return x
}

func requireValue(t *testing.T, c *BigIntContainer, h BigIntHandle, expected *big.Int) {
	t.Helper()
	require.Equal(t, 0, expected.Cmp(c.Get(h)), "want %s, got %s", expected, c.Get(h))
}

func TestInsertGet_Boundaries(t *testing.T) {
	c := NewBigIntContainer()
	values := boundaryValues()
	handles := make([]BigIntHandle, len(values))
	for i, v := range values {
		handles[i] = c.Insert(v)
	}
	// check only after everything is inserted, to catch values overwriting each other
	for i, v := range values {
		requireValue(t, c, handles[i], v)
	}
}

func TestInsert_ZeroUsesNoStorage(t *testing.T) {
	c := NewBigIntContainer()
	h := c.Insert(big.NewInt(0))
	require.Equal(t, Zero, h)
	require.Empty(t, c.data)
	require.Equal(t, int64(0), c.Get(Zero).Int64())
}

func TestInsert_DoesNotAliasInput(t *testing.T) {
	c := NewBigIntContainer()
	x := big.NewInt(1000)
	h := c.Insert(x)
	x.SetInt64(7)
	requireValue(t, c, h, big.NewInt(1000))
}

func TestInsertUint64(t *testing.T) {
	c := NewBigIntContainer()
	for _, v := range []uint64{0, 1, 1 << 63, ^uint64(0)} {
		requireValue(t, c, c.InsertUint64(v), new(big.Int).SetUint64(v))
	}
}

func TestGet_IsCopy_GetUnsafe_IsView(t *testing.T) {
	c := NewBigIntContainer()
	h := c.Insert(big.NewInt(100))

	c.Get(h).SetInt64(1)
	requireValue(t, c, h, big.NewInt(100))

	c.GetUnsafe(h).SetInt64(1)
	requireValue(t, c, h, big.NewInt(1))
}

func TestUpdate(t *testing.T) {
	c := NewBigIntContainer()
	h := c.Insert(bi("0x10000000000000000000000000000000000000000")) // 3 words

	// shrinks in place: length must be refreshed
	h = c.Update(h, big.NewInt(-5))
	requireValue(t, c, h, big.NewInt(-5))

	// grows back within the original capacity
	big2 := bi("0x1000000000000000000000000")
	h = c.Update(h, big2)
	requireValue(t, c, h, big2)

	// nil means zero
	h = c.Update(h, nil)
	requireValue(t, c, h, big.NewInt(0))

	// from Zero
	h = c.Update(Zero, big.NewInt(42))
	requireValue(t, c, h, big.NewInt(42))
}

func TestSetInt64_SetBytes(t *testing.T) {
	c := NewBigIntContainer()
	h := c.Insert(bi("0x10000000000000000000000000000000000000000"))
	h = c.SetInt64(h, -9)
	requireValue(t, c, h, big.NewInt(-9))
	require.True(t, c.IsInt64(h))
	require.Equal(t, int64(-9), c.ToInt64(h))

	h = c.SetBytes(h, []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}) // sets sign to positive
	requireValue(t, c, h, bi("0x01000000000000000000"))
	require.False(t, c.IsInt64(h))

	h = c.SetBytes(h, []byte{0, 0, 7})
	requireValue(t, c, h, big.NewInt(7))

	h = c.SetBytes(h, nil)
	requireValue(t, c, h, big.NewInt(0))
}

func TestGetBytes_BitLen_ByteLen(t *testing.T) {
	c := NewBigIntContainer()
	for _, v := range boundaryValues() {
		h := c.Insert(v)
		abs := new(big.Int).Abs(v)
		require.Equal(t, abs.Bytes(), c.GetBytes(h), v.String())
		require.Equal(t, abs.BitLen(), c.BitLen(h), v.String())
		require.Equal(t, len(abs.Bytes()), c.ByteLen(h), v.String())
	}
	require.Equal(t, 0, c.ByteLen(Zero))
	require.Equal(t, 0, c.BitLen(Zero))
	require.Empty(t, c.GetBytes(Zero))
}

// The destination handle has spare capacity, so results stay in place
// and the handle length and sign have to be refreshed.
func TestInPlace_LengthAndSignRefreshed(t *testing.T) {
	c := NewBigIntContainer()
	wide := bi("0x100000000000000000000000000000000000000000000000000000000000000") // 4 words
	neighbour := c.Insert(big.NewInt(12345))

	h := c.Insert(wide)
	small := c.Insert(big.NewInt(5))

	// shrink
	h = c.Set(h, small)
	requireValue(t, c, h, big.NewInt(5))
	require.Equal(t, 1, h.length)

	// grow within capacity
	h = c.Set(h, c.Insert(wide))
	requireValue(t, c, h, wide)

	// sign flip
	h = c.Neg(h, h)
	requireValue(t, c, h, new(big.Int).Neg(wide))

	// to zero, and back
	h = c.Sub(h, h, h)
	requireValue(t, c, h, big.NewInt(0))
	h = c.Add(h, small, small)
	requireValue(t, c, h, big.NewInt(10))

	requireValue(t, c, neighbour, big.NewInt(12345))
	requireValue(t, c, small, big.NewInt(5))
}

func TestAliasing(t *testing.T) {
	c := NewBigIntContainer()
	for _, v := range boundaryValues() {
		x := c.Insert(v)
		x = c.Add(x, x, x)
		requireValue(t, c, x, new(big.Int).Add(v, v))
		x = c.Mul(x, x, x)
		requireValue(t, c, x, new(big.Int).Mul(new(big.Int).Add(v, v), new(big.Int).Add(v, v)))
		x = c.Sub(x, x, x)
		requireValue(t, c, x, big.NewInt(0))
	}
}

func TestSetCopiesData(t *testing.T) {
	c := NewBigIntContainer()
	x := c.Insert(big.NewInt(77))
	y := c.Set(Zero, x)
	y = c.Add(y, y, c.Insert(big.NewInt(1)))
	requireValue(t, c, x, big.NewInt(77))
	requireValue(t, c, y, big.NewInt(78))
}

func TestCmp(t *testing.T) {
	c := NewBigIntContainer()
	values := boundaryValues()
	handles := make([]BigIntHandle, len(values))
	for i, v := range values {
		handles[i] = c.Insert(v)
	}
	for i, a := range values {
		for j, b := range values {
			require.Equal(t, a.Cmp(b), c.Cmp(handles[i], handles[j]))
			require.Equal(t, a.CmpAbs(b), c.CmpAbs(handles[i], handles[j]))
		}
	}
}

func TestUnaryOperations(t *testing.T) {
	c := NewBigIntContainer()
	for _, v := range boundaryValues() {
		x := c.Insert(v)
		requireValue(t, c, c.Abs(Zero, x), new(big.Int).Abs(v))
		requireValue(t, c, c.Neg(Zero, x), new(big.Int).Neg(v))
		requireValue(t, c, c.Not(Zero, x), new(big.Int).Not(v))
		if v.Sign() >= 0 {
			requireValue(t, c, c.Sqrt(Zero, x), new(big.Int).Sqrt(v))
		}
		requireValue(t, c, x, v) // operand untouched
	}
}

func TestSqrt_NegativePanics(t *testing.T) {
	c := NewBigIntContainer()
	require.Panics(t, func() { c.Sqrt(Zero, c.Insert(big.NewInt(-4))) })
}

func TestDivisionByZeroPanics(t *testing.T) {
	c := NewBigIntContainer()
	x := c.Insert(big.NewInt(4))
	for name, op := range map[string]func(d, x, y BigIntHandle) BigIntHandle{
		"Quo": c.Quo, "Rem": c.Rem, "Div": c.Div, "Mod": c.Mod,
	} {
		require.Panics(t, func() { op(Zero, x, Zero) }, name)
	}
}

type binaryCase struct {
	name     string
	op       func(c *BigIntContainer, d, x, y BigIntHandle) BigIntHandle
	ref      func(z, x, y *big.Int) *big.Int
	nonZeroY bool
}

var binaryCases = []binaryCase{
	{"Add", (*BigIntContainer).Add, (*big.Int).Add, false},
	{"Sub", (*BigIntContainer).Sub, (*big.Int).Sub, false},
	{"Mul", (*BigIntContainer).Mul, (*big.Int).Mul, false},
	{"Quo", (*BigIntContainer).Quo, (*big.Int).Quo, true},
	{"Rem", (*BigIntContainer).Rem, (*big.Int).Rem, true},
	{"Div", (*BigIntContainer).Div, (*big.Int).Div, true},
	{"Mod", (*BigIntContainer).Mod, (*big.Int).Mod, true},
	{"And", (*BigIntContainer).And, (*big.Int).And, false},
	{"AndNot", (*BigIntContainer).AndNot, (*big.Int).AndNot, false},
	{"Or", (*BigIntContainer).Or, (*big.Int).Or, false},
	{"Xor", (*BigIntContainer).Xor, (*big.Int).Xor, false},
}

func TestBinaryOperations_Boundaries(t *testing.T) {
	values := boundaryValues()
	for _, tc := range binaryCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewBigIntContainer()
			for _, a := range values {
				for _, b := range values {
					if tc.nonZeroY && b.Sign() == 0 {
						continue
					}
					x, y := c.Insert(a), c.Insert(b)
					expected := tc.ref(new(big.Int), a, b)

					requireValue(t, c, tc.op(c, Zero, x, y), expected)

					// reused destination, with leftover size from a previous result
					d := c.Insert(bi("0x100000000000000000000000000000000000000000000000000"))
					requireValue(t, c, tc.op(c, d, x, y), expected)

					requireValue(t, c, x, a)
					requireValue(t, c, y, b)
				}
			}
		})
	}
}

// Random program run on the container and on plain big.Int in lockstep.
// A pool of numbers is kept; every step writes into a random pool slot (often aliasing the
// operands), then the whole pool is compared, which catches any operation corrupting neighbours.
func TestRandomOperations_MatchBigInt(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	c := NewBigIntContainer()

	const poolSize = 8
	handles := make([]BigIntHandle, poolSize)
	refs := make([]*big.Int, poolSize)
	for i := range handles {
		refs[i] = randomBigInt(r, 200)
		handles[i] = c.Insert(refs[i])
	}

	for step := 0; step < 20000; step++ {
		d, a, b := r.Intn(poolSize), r.Intn(poolSize), r.Intn(poolSize)
		tc := binaryCases[r.Intn(len(binaryCases))]

		switch {
		case r.Intn(8) == 0:
			// refill from a fresh random value, of any size
			refs[d] = randomBigInt(r, 300)
			handles[d] = c.Update(handles[d], refs[d])
		case tc.nonZeroY && refs[b].Sign() == 0:
			continue
		default:
			expected := tc.ref(new(big.Int), refs[a], refs[b])
			if expected.BitLen() > 600 { // keep numbers bounded
				expected.Rsh(expected, uint(expected.BitLen()-300))
				refs[d] = expected
				handles[d] = c.Update(handles[d], expected)
				break
			}
			handles[d] = tc.op(c, handles[d], handles[a], handles[b])
			refs[d] = expected
		}

		for i := range handles {
			if refs[i].Cmp(c.Get(handles[i])) != 0 {
				t.Fatalf("step %d (%s): slot %d want %s, got %s", step, tc.name, i, refs[i], c.Get(handles[i]))
			}
		}
	}
}
