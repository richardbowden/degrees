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

	extended := ctx.Bool("extended")

	// Extract VCS information
	var (
		vcsRevision = "unknown"
		vcsTime     = "unknown"
		vcsModified = "false"
		goVersion   = bi.GoVersion
		path        = bi.Path
		version     = bi.Main.Version
	)

	for _, setting := range bi.Settings {
		switch setting.Key {
		case "vcs.revision":
			vcsRevision = setting.Value
		case "vcs.time":
			vcsTime = setting.Value
		case "vcs.modified":
			vcsModified = setting.Value
		}
	}

	// Basic info
	fmt.Println("Build Information:")
	fmt.Printf("  Version:    %s\n", version)
	fmt.Printf("  Commit:     %s\n", vcsRevision)
	fmt.Printf("  Built:      %s\n", vcsTime)
	fmt.Printf("  Modified:   %s\n", vcsModified)
	fmt.Printf("  Go Version: %s\n", goVersion)

	if extended {
		fmt.Printf("\nModule Information:\n")
		fmt.Printf("  Path:       %s\n", path)

		if len(bi.Deps) > 0 {
			fmt.Printf("\nDependencies (%d):\n", len(bi.Deps))
			for _, dep := range bi.Deps {
				fmt.Printf("  %s %s\n", dep.Path, dep.Version)
			}
		}

		fmt.Printf("\nBuild Settings:\n")
		for _, setting := range bi.Settings {
			fmt.Printf("  %s = %s\n", setting.Key, setting.Value)
		}
	}

	return nil
}
