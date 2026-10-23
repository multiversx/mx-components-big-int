# mx-components-big-int

Small Go building blocks for working with `math/big` integers.

| Package | What it is for |
|---|---|
| [`twos-complement`](twos-complement/) | Convert `*big.Int` to and from big-endian two's complement bytes, and compute their length without building them. |
| [`managed-big-int`](managed-big-int/) | An arena that stores many big integers in one flat slice and refers to them with small pointer-free handles. |

The two packages are independent of each other.

- `twos-complement` is used in production by the MultiversX Go VM, [mx-chain-vm-go](https://github.com/multiversx/mx-chain-vm-go).
- `managed-big-int` is **experimental: it is not used anywhere and has never been tested in production.**

## Development

```
make test   # unit tests with race detector and atomic coverage profile (coverage.txt)
```

CI (`.github/workflows`) runs the tests, a coverage report with a minimum threshold, and `golangci-lint`.
