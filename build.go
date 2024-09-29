//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	Default        = Debug
	debugBinName   = "p402-debug"
	releaseBinName = "p402-release"
)

var (
	buildDebugDir       string
	buildDebugBinPath   string
	buildReleaseDir     string
	buildReleaseBinPath string
)

var cwd string

func init() {
	cwd, _ = os.Getwd()

	buildBaseDir := filepath.Join(cwd, "build")

	buildDebugDir = filepath.Join(buildBaseDir, "debug/")
	buildDebugBinPath = filepath.Join(buildDebugDir, debugBinName)
	buildReleaseDir = filepath.Join(buildBaseDir, "release/")
	buildReleaseBinPath = filepath.Join(buildReleaseDir, releaseBinName)
}

func Debug() error {
	fmt.Printf("\nbuild-mode: debug\n")

	start := time.Now()

	err := sh.RunV("go", "build", "-race", "-v", "-gcflags=all=-N", "-o", buildDebugBinPath, "./cmd/server")
	elapsed := time.Since(start)

	fmt.Printf("\npath: %s\n", buildDebugBinPath)
	fmt.Printf("took %s\n\n", elapsed)
	return err
}

func Release() error {
	err := sh.RunV("go", "build", "-v", "-o", buildReleaseBinPath, "./cmd/server")
	fmt.Printf("build-mode: release\n")
	fmt.Printf("path: %s\n", buildReleaseBinPath)
	return err
}

func Clean() error {
	fmt.Println("clean command")
	return nil
}

func Run() error {
	mg.SerialDeps(Debug)
	// err := Debug()

	// if err != nil {
	// 	return err
	// }

	return sh.RunV(buildDebugBinPath, "--human-logs", "server", "run")
}
