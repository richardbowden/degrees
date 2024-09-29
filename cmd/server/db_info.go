package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func dbInfo(ctx *cli.Context) error {
	fmt.Println("db info command")
	return nil
}
