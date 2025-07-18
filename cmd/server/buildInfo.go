package main

import (
	"fmt"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

func buildInfo(ctx *cli.Context) error {

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("not ok")
	}
	fmt.Printf("%+v", bi)
	return nil
}
