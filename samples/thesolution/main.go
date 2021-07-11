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
