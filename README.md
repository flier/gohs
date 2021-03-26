# gohs [![travis](https://travis-ci.org/flier/gohs.svg)](https://travis-ci.org/flier/gohs) [![github](https://github.com/flier/gohs/workflows/Continuous%20integration/badge.svg)](https://github.com/flier/gohs/actions?query=workflow%3A%22Continuous+integration%22) [![go report](https://goreportcard.com/badge/github.com/flier/gohs)](https://goreportcard.com/report/github.com/flier/gohs) [![Go Reference](https://pkg.go.dev/badge/github.com/flier/gohs/hyperscan.svg)](https://pkg.go.dev/github.com/flier/gohs/hyperscan) [![Apache](https://img.shields.io/badge/license-Apache-blue.svg)](https://github.com/flier/gohs/blob/master/LICENSE-APACHE) [![MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/flier/gohs/blob/master/LICENSE-MIT)

Golang binding for Intel's HyperScan regex matching library: [hyperscan.io](https://www.hyperscan.io/)

## Build

**Note:** `gohs` will use Hyperscan v5 API by default, you can also build for Hyperscan v4 with `hyperscan_v4` tag.

```bash
go get -u -tags hyperscan_v4 github.com/flier/gohs/hyperscan
```

## License

This project is licensed under either of Apache License ([LICENSE-APACHE](LICENSE-APACHE)) or MIT license ([LICENSE-MIT](LICENSE-MIT)) at your option.

## Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in Futures by you, as defined in the Apache-2.0 license, shall be
dual licensed as above, without any additional terms or conditions.
