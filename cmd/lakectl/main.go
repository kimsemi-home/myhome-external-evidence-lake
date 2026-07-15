package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kimsemi-home/myhome-external-evidence-lake/internal/lake"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	root, err := repositoryRoot()
	if err != nil {
		return err
	}
	if len(args) == 1 && args[0] == "verify" {
		return report("public evidence lake verified", lake.Verify(root))
	}
	if len(args) == 2 && args[0] == "docs" && args[1] == "generate" {
		return report("documents generated", lake.GenerateDocuments(root))
	}
	if len(args) == 2 && args[0] == "docs" && args[1] == "check" {
		return report("documents verified", lake.CheckDocuments(root))
	}
	if len(args) == 2 && args[0] == "trace" && args[1] == "verify" {
		return report("traceability verified", lake.VerifyTrace(root))
	}
	if len(args) == 2 && args[0] == "trace" && args[1] == "staged" {
		return report("staged traceability verified", lake.TraceStaged(root))
	}
	if len(args) == 3 && args[0] == "trace" && args[1] == "range" {
		return report("range traceability verified", lake.TraceRange(root, args[2]))
	}
	if len(args) == 2 && args[0] == "safety" && args[1] == "check" {
		return report("public-safety scan passed", lake.SafetyCheck(root))
	}
	return errors.New("usage: lakectl <verify|docs generate|docs check|trace verify|trace staged|trace range BASE|safety check>")
}

func repositoryRoot() (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			return root, nil
		}
		parent := filepath.Dir(root)
		if parent == root {
			return "", errors.New("repository root not found")
		}
		root = parent
	}
}

func report(message string, err error) error {
	if err != nil {
		return err
	}
	fmt.Println(message)
	return nil
}
