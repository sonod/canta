# CANTA - ConsAl eveNT Accepter

## Description
canta reads the payload of the consul event and performs deletion of Nginx's cache.

## Usage
consul watch -type event -name "hello" ./canta

## Install

To install, use `go get`:

```bash
$ go get -d github.com/noda/canta
```

## Contribution

1. Fork ([https://github.com/noda/canta/fork](https://github.com/noda/canta/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[noda](https://github.com/noda)
