package main

import (
	"fmt"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

func buildInfo(ctx *cli.Context) error {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("failed to get build info")
	}
	fmt.Printf("%+v", bi)
	return nil
}
