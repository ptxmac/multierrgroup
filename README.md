# multierrgroup

[![Go Reference](https://pkg.go.dev/badge/go.ptx.dk/multierrgroup.svg)](https://pkg.go.dev/go.ptx.dk/multierrgroup)
[![Github Workflow](https://github.com/ptxmac/multierrgroup/actions/workflows/go.yml/badge.svg)](https://github.com/ptxmac/multierrgroup/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/ptxmac/multierrgroup/branch/master/graph/badge.svg)](https://codecov.io/gh/ptxmac/multierrgroup)

A simple combination of [golang.org/x/sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) and
[go.uber.org/multierr](https://pkg.go.dev/go.uber.org/multierr).

`multierrgroup` have the exact same behaviour as `errgroup`,
_except_ that it captures _all_ errors returned from the goroutines instead of just one.

## Usage

`multierrgroup` is a drop-in replacement for `errgroup` and have the same methods and behaviour.

When any goroutine returns an error, it will colleceted using `multierr.Append` and `Wait` will return the multierr.
Thus error can be split into their original errors, or worked on using `errors.As/Is`.

See the [multierr documentation](https://pkg.go.dev/go.uber.org/multierr) for more information.

## Sample

```go
package main

import (
	"errors"
	"fmt"
	"os"

	"go.ptx.dk/multierrgroup"
	"go.uber.org/multierr"
)

func main() {
	g := multierrgroup.Group{}

	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("Error: %d", i)
		g.Go(func() error {
			return errors.New(str)
		})
		if i == 5 {
			g.Go(func() error {
				return os.ErrNotExist
			})
		}
	}
	err := g.Wait()
	fmt.Println("Got", len(multierr.Errors(err)), "errors")
	fmt.Println(err.Error())
	fmt.Println("Was one of the errors a ErrNotExists?", errors.Is(err, os.ErrNotExist))
}
```

Running the sample shows that errors from all goroutines are kept:

```
$ go run thesolution/main.go
Got 11 errors
Error: 9; Error: 5; file does not exist; Error: 6; Error: 7; Error: 8; Error: 1; Error: 0; Error: 2; Error: 3; Error: 4
Was one of the errors a ErrNotExists? true
```
