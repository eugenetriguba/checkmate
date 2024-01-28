# Checkmate

<p>
    <a href="https://godoc.org/github.com/eugenetriguba/checkmate">
        <img src="https://godoc.org/github.com/eugenetriguba/checkmate?status.svg" alt="GoDoc">
    </a>
    <a href="https://goreportcard.com/report/github.com/eugenetriguba/checkmate">
        <img src="https://goreportcard.com/badge/github.com/eugenetriguba/checkmate" alt="Go Report Card Badge">
    </a>
    <a href="https://codecov.io/github/eugenetriguba/checkmate">
        <img src="https://codecov.io/github/eugenetriguba/checkmate/branch/main/graph/badge.svg?token=Z3ZD0JWGBr"/>
    </a>
    <img alt="Version Badge" src="https://img.shields.io/badge/version-0.2.0-blue" style="max-width:100%;">
</p>

Checkmate is a Go library designed to enhance the testing experience by providing a set of assertion functions. With Go, there is no built-in assertion library which can make writing tests feel more verbose than it needs to be, so Checkmate aims to fill that gap.

This is a simple and small library. If you need something more advanced and full-featured, [Testify](https://github.com/stretchr/testify), [Gomega](https://github.com/onsi/gomega), or [gotest.tools](https://github.com/gotestyourself/gotest.tools) are all great options.

## Installation

To install Checkmate, use the following command:

```bash
$ go get github.com/eugenetriguba/checkmate
```

## Usage

To get a feel for what Checkmate has to offer, here is a simple example with some commentary that shows all the supported assertion functions. For a reference on each function, see the [GoDoc](https://godoc.org/github.com/eugenetriguba/checkmate).

```go

import (
    "errors"
    "os"
    "testing"

    "github.com/eugenetriguba/checkmate/assert"
)

func TestAssertions(t *testing.T) {
    // Assert is the base general assertion function. It
    // takes in any boolean expression and will fail the
    // test if the expression is false.
    assert.Assert(t, true)

    // Also note, all functions in this library take in an
    // optional message argument which may have placeholders
    // that are formatted with the reaming arguments.
    assert.Assert(t, false, "expected %v to be true", false)

    // Check is like Assert, but it will not fail the test.
    assert.Check(t, true)

    assert.True(t, true)
    assert.False(t, false)
    assert.Equal(t, 1, 1)
    assert.DeepEqual(t, []int{1, 2, 3}, []int{1, 2, 3})
    assert.LenEqual(t, []int{1, 2, 3}, 3)
    assert.NotEqual(t, 1, 2)
    assert.Nil(t, nil)
    assert.NotNil(t, "not nil")
    assert.ErrorIs(t, os.ErrInvalid, os.ErrInvalid)
    assert.ErrorContains(t, errors.New("error 1"), "error 1")
}
```