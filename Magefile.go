//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Installing dependencies
func Prepare() error {
	return sh.Run("go", "mod", "download")
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
