//go:build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/magefile/mage/sh"
)

var (
	PROJECT_DIR, _ = os.Getwd()
	PROJECT_BIN    = fmt.Sprintf("%s/bin", PROJECT_DIR)
	SWAGGER        = fmt.Sprintf("%s/swagger", PROJECT_BIN)
)

func runInDir(dir string, env map[string]string, cmd []string) error {
	os.Chdir(dir)
	defer os.Chdir(PROJECT_DIR)
	return sh.RunWith(env, cmd[0], cmd[1:]...)
}

func createBin() error {
	_, err := os.Stat(PROJECT_BIN)
	if !os.IsNotExist(err) {
		return err
	}
	return os.Mkdir(PROJECT_BIN, 0750)
}

func goGetTool(target, source string) error {
	_, err := os.Stat(target)
	if !os.IsNotExist(err) {
		return err
	}
	err = createBin()
	if err != nil {
		return err
	}
	dir, err := ioutil.TempDir("/tmp", "*-isoboard-get-tool")
	defer sh.Rm(dir)
	end := map[string]string{"GOBIN": PROJECT_BIN}
	cmd := []string{"go", "install", source}
	err = runInDir(dir, end, cmd)
	if err != nil {
		return err
	}
	return nil
}

func swagger() error {
	return goGetTool(
		SWAGGER,
		"github.com/go-swagger/go-swagger/cmd/swagger@latest",
	)
}

// Installing dependencies
func Prepare() error {
	err := swagger()
	if err != nil {
		return err
	}
	return sh.Run("go", "mod", "download")
}

func Generate() error {
	return sh.Run(
		SWAGGER,
		"generate", "server",
		"--spec", "openapi.yaml",
		"--target", "backend/generated",
	)
	return nil
}

// Linting any errors
func Lint() error {
	return sh.Run("go", "vet", "./...")
}

// Building single binary
func Build() error {
	return sh.Run("go", "build", "main.go")
}

func Run() error {
	return sh.Run("go", "run", "main.go")
}

// Test package
func Test() error {
	return sh.Run("go", "test", "./...")
}
