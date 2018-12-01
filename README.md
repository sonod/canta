# CANTA - ConsAl eveNT Accepter

## Description
canta reads, displays or executes the payload of consul event.

## Usage

```
Usage of canta:
  canta [OPTIONS]
OPTIONS
  --version  Print version information and quit.
  --run  running command in event payload
```


#### with consul
- Display the payload of consul event.

*Watch the consul event.*
```
consul watch -type event -name "hello" canta
```

*Fire the event.*
```
consul event -name "hello" 'hoge'
```

- Executes the payload of consul event.

*Watch the consul event.*
```
consul watch -type event -name "hello" canta --run
```

*Fire the event.*
```
consul event -name "hello" 'ls'
```

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
