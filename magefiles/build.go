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

func Gen() error {
	return sqlGen()
}

func sqlGen() error {
	start := time.Now()
	fmt.Printf("\nrunning sql gen\n")
	err := sh.RunV("sqlc", "generate")
	elapsed := time.Since(start)

	fmt.Printf("took %s\n\n", elapsed)

	return err
}

func Debug() error {
	mg.Deps(sqlGen)
	start := time.Now()
	fmt.Printf("\nbuild-mode: debug\n")

	err := sh.RunV("go", "build", "-race", "-v", "-gcflags=all=-N", "-o", buildDebugBinPath, "./cmd/server")
	elapsed := time.Since(start)

	fmt.Printf("\npath: %s\n", buildDebugBinPath)
	fmt.Printf("took %s\n\n", elapsed)
	return err
}

func Release() error {
	fmt.Printf("\nbuild-mode: release\n")
	fmt.Printf("-------------------\n")

	mg.Deps(sqlGen, Test)
	err := sh.RunV("go", "build", "-v", "-o", buildReleaseBinPath, "./cmd/server")

	fmt.Printf("\npath: %s\n", buildReleaseBinPath)
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

func Test() error {
	err := sh.RunV("go", "test", "./...")
	return err
}
