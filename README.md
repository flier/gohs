# gohs [![Continuous integration](https://github.com/flier/gohs/actions/workflows/ci.yml/badge.svg?)](https://github.com/flier/gohs/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/flier/gohs?)](https://goreportcard.com/report/github.com/flier/gohs) [![codecov](https://codecov.io/gh/flier/gohs/branch/master/graph/badge.svg?token=F5CLCxpJGM)](https://codecov.io/gh/flier/gohs)  [![Apache](https://img.shields.io/badge/license-Apache-blue.svg)](https://github.com/flier/gohs/blob/master/LICENSE-APACHE) [![MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/flier/gohs/blob/master/LICENSE-MIT)

Golang binding for Intel's HyperScan regex matching library: [hyperscan.io](https://www.hyperscan.io/)

## Hyperscan [![Go Reference](https://pkg.go.dev/badge/github.com/flier/gohs/hyperscan.svg)](https://pkg.go.dev/github.com/flier/gohs/hyperscan)

Hyperscan is a software regular expression matching engine designed with high performance and flexibility in mind. It is implemented as a library that exposes a straightforward C API.

### Build

`gohs` does not enable the latest api of Hyperscan v5.4 by default, if you want to use it please pass build tags `hyperscan_v54`.

```bash
go get -u -tags hyperscan_v54 github.com/flier/gohs/hyperscan
```

`gohs` will use Hyperscan v5 API by default, you can also build for Hyperscan v4 with `hyperscan_v4` tag.

```bash
go get -u -tags hyperscan_v4 github.com/flier/gohs/hyperscan
```

## Chimera [![Go Reference](https://pkg.go.dev/badge/github.com/flier/gohs/chimera.svg)](https://pkg.go.dev/github.com/flier/gohs/chimera)

Chimera is a software regular expression matching engine that is a hybrid of Hyperscan and PCRE. The design goals of Chimera are to fully support PCRE syntax as well as to take advantage of the high performance nature of Hyperscan.

### Build

It is recommended to compile and link Chimera using static libraries.

```bash
$ mkdir build && cd build
$ cmake .. -G Ninja -DBUILD_STATIC_LIBS=on
$ ninja && ninja install
$ go get -u github.com/flier/gohs/hyperscan
```

### Note

You need to download the PCRE library source code to build Chimera, see [Chimera Requirements](https://intel.github.io/hyperscan/dev-reference/chimera.html#requirements) for more details

### Vectorscan

The Intel's Hyperscan relies on the SSE/AVX instruction set for x86 platform. If you need to use it on other platforms such as ARM, you can install [Vectorscan](https://github.com/VectorCamp/vectorscan) which is a portable fork of Hyperscan that supports the ARM NEON/ASIMD and Power VSX.

## License

This project is licensed under either of Apache License ([LICENSE-APACHE](LICENSE-APACHE)) or MIT license ([LICENSE-MIT](LICENSE-MIT)) at your option.

## Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in Futures by you, as defined in the Apache-2.0 license, shall be
dual licensed as above, without any additional terms or conditions.
