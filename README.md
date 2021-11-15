# gosqlparser

[![Build Status](https://app.travis-ci.com/krasun/gosqlparser.svg?branch=main)](https://app.travis-ci.com/krasun/gosqlparser)
[![codecov](https://codecov.io/gh/krasun/gosqlparser/branch/main/graph/badge.svg?token=8NU6LR4FQD)](https://codecov.io/gh/krasun/gosqlparser)
[![Go Report Card](https://goreportcard.com/badge/github.com/krasun/gosqlparser)](https://goreportcard.com/report/github.com/krasun/gosqlparser)
[![GoDoc](https://godoc.org/https://godoc.org/github.com/krasun/gosqlparser?status.svg)](https://godoc.org/github.com/krasun/gosqlparser)

`gosqlparser` is a simple SQL parser.

## Installation

As simple as:

```
go get github.com/krasun/gosqlparser
```

## Supported Statements

1. SELECT: 
```
SELECT c1, c2 FROM table1 WHERE c3 == c4 AND c5 == c6
```

## Usage 

... 

## Tests 

To make sure that the code is fully tested and covered:

```
$ go test .
ok  	github.com/krasun/gosqlparser	0.470s
```

## Known Usages 

1. [krasun/gosqldb](https://github.com/krasun/gosqldb) - my experimental implementation of a simple database.

## License 

**gosqlparser** is released under [the MIT license](LICENSE).