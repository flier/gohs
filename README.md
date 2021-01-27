# gohs [![travis](https://travis-ci.org/flier/gohs.svg)](https://travis-ci.org/flier/gohs) [![build](https://github.com/flier/gohs/workflows/Continuous%20integration/badge.svg)](https://github.com/flier/gohs/actions?query=workflow%3A%22Continuous+integration%22) [![Go Reference](https://pkg.go.dev/badge/github.com/flier/gohs/hyperscan.svg)](https://pkg.go.dev/github.com/flier/gohs/hyperscan) [![Apache](https://img.shields.io/badge/license-Apache-blue.svg)](https://github.com/flier/gohs/blob/master/LICENSE-APACHE) [![MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/flier/gohs/blob/master/LICENSE-MIT)

GoLang Binding of Intel's HyperScan regex matching library: https://www.hyperscan.io/

[API Reference](https://godoc.org/github.com/flier/gohs/hyperscan)

# Build

**Note:** `gohs` will use Hyperscan v5 API by default, you can also build for Hyperscan v4 with `hyperscan_v4` tag.

```bash
$ go get -u -tags hyperscan_v4 github.com/flier/gohs/hyperscan
```

# License

This project is licensed under either of

 * Apache License, Version 2.0, ([LICENSE-APACHE](LICENSE-APACHE) or
   http://www.apache.org/licenses/LICENSE-2.0)
 * MIT license ([LICENSE-MIT](LICENSE-MIT) or
   http://opensource.org/licenses/MIT)

at your option.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in Futures by you, as defined in the Apache-2.0 license, shall be
dual licensed as above, without any additional terms or conditions.
