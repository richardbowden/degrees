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
	debugBinName   = "degrees-debug"
	releaseBinName = "degrees-release"
)

var (
	buildDebugDir       string
	buildDebugBinPath   string
	buildReleaseDir     string
	buildReleaseBinPath string
)

var cwd string
var hasBuildThisRun bool = false

func init() {
	cwd, _ = os.Getwd()

	buildBaseDir := filepath.Join(cwd, "build")

	buildDebugDir = filepath.Join(buildBaseDir, "debug/")
	buildDebugBinPath = filepath.Join(buildDebugDir, debugBinName)
	buildReleaseDir = filepath.Join(buildBaseDir, "release/")
	buildReleaseBinPath = filepath.Join(buildReleaseDir, releaseBinName)
}

func Gen() error {
	err := protoGen()
	if err != nil {
		return err
	}

	err = sqlGen()
	if err != nil {
		return err
	}
	return goGen()
}

func protoGen() error {
	start := time.Now()
	fmt.Printf("\nrunning proto gen\n")
	// Change to proto directory and run buf generate
	err := sh.RunV("sh", "-c", "cd proto && go run github.com/bufbuild/buf/cmd/buf@latest generate")
	elapsed := time.Since(start)

	fmt.Printf("took %s\n\n", elapsed)

	return err
}

func sqlGen() error {
	start := time.Now()
	fmt.Printf("\nrunning sql gen\n")
	err := sh.RunV("go", "tool", "sqlc", "generate")
	elapsed := time.Since(start)

	fmt.Printf("took %s\n\n", elapsed)

	return err
}

func goGen() error {
	start := time.Now()
	fmt.Println("running go generate")
	err := sh.RunV("go", "generate", "./...")
	elapsed := time.Since(start)

	fmt.Printf("took %s\n\n", elapsed)

	return err
}

func Migrate() error {
	mg.Deps(Debug)

	return sh.RunV(buildDebugBinPath, "db", "migration", "up")

}

func Debug() error {
	mg.Deps(sqlGen, goGen)
	start := time.Now()
	fmt.Printf("\nbuild-mode: debug\n")

	if !hasBuildThisRun{
		err := sh.RunV("go", "build", "-race", "-v", "-gcflags=all=-N", "-o", buildDebugBinPath, "./cmd/server")
		elapsed := time.Since(start)
		hasBuildThisRun = true
		fmt.Printf("\npath: %s\n", buildDebugBinPath)
		fmt.Printf("took %s\n\n", elapsed)
		return err
	}

	fmt.Printf("\npath(already build this run): %s\n", buildDebugBinPath)
	fmt.Printf("")
	return nil
}

func Release() error {
	fmt.Printf("\nbuild-mode: release\n")
	fmt.Printf("-------------------\n")

	mg.Deps(sqlGen, goGen, Test)
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

// Proto namespace for proto-related commands
type Proto mg.Namespace

// Lint runs buf lint on proto files
func (Proto) Lint() error {
	fmt.Println("linting proto files")
	return sh.RunV("sh", "-c", "cd proto && go run github.com/bufbuild/buf/cmd/buf@latest lint")
}

// Format formats proto files
func (Proto) Format() error {
	fmt.Println("formatting proto files")
	return sh.RunV("sh", "-c", "cd proto && go run github.com/bufbuild/buf/cmd/buf@latest format -w")
}

// Breaking checks for breaking changes against main branch
func (Proto) Breaking() error {
	fmt.Println("checking for breaking changes")
	return sh.RunV("sh", "-c", "cd proto && go run github.com/bufbuild/buf/cmd/buf@latest breaking --against '.git#branch=main'")
}

// GRPC namespace for gRPC testing commands
type GRPC mg.Namespace

// List lists all gRPC services (requires server running on :9090)
func (GRPC) List() error {
	fmt.Println("listing gRPC services")
	return sh.RunV("go", "run", "github.com/fullstorydev/grpcurl/cmd/grpcurl@latest", "-plaintext", "localhost:9090", "list")
}

// UI opens interactive gRPC UI (requires server running on :9090)
func (GRPC) UI() error {
	fmt.Println("opening gRPC UI at http://localhost:8081")
	return sh.RunV("go", "run", "github.com/fullstorydev/grpcui/cmd/grpcui@latest", "-plaintext", "-port", "8081", "localhost:9090")
}
