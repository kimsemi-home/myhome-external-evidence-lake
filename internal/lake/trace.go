package lake

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type traceManifest struct {
	SchemaVersion string       `json:"schema_version"`
	Groups        []traceGroup `json:"groups"`
}

type traceGroup struct {
	ID                 string   `json:"id"`
	DocumentSources    []string `json:"document_sources"`
	GeneratedDocuments []string `json:"generated_documents"`
	Code               []string `json:"code"`
	Tests              []string `json:"tests"`
	ChangePolicy       string   `json:"change_policy"`
}

func loadTrace(root string) (traceManifest, error) {
	var manifest traceManifest
	if err := decodeStrict(filepath.Join(root, "traceability.json"), &manifest); err != nil {
		return manifest, err
	}
	return manifest, verifyTracePaths(root, manifest)
}

func VerifyTrace(root string) error {
	_, err := loadTrace(root)
	return err
}

func verifyTracePaths(root string, manifest traceManifest) error {
	if manifest.SchemaVersion != "repo.traceability/v1" || len(manifest.Groups) == 0 {
		return errors.New("traceability manifest header invalid")
	}
	for _, group := range manifest.Groups {
		if group.ID == "" || group.ChangePolicy != "bidirectional" || len(group.DocumentSources) == 0 || len(group.GeneratedDocuments) == 0 || len(group.Code) == 0 || len(group.Tests) == 0 {
			return fmt.Errorf("traceability group %q is incomplete", group.ID)
		}
		paths := append(append(append([]string{}, group.DocumentSources...), group.GeneratedDocuments...), append(group.Code, group.Tests...)...)
		for _, rel := range paths {
			if filepath.IsAbs(rel) || filepath.Clean(rel) != rel || strings.HasPrefix(rel, "..") {
				return fmt.Errorf("unsafe trace path %q", rel)
			}
			if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
				return fmt.Errorf("missing trace path %q", rel)
			}
		}
	}
	return nil
}

func TraceStaged(root string) error {
	files, err := gitFiles(root, "diff", "--cached", "--name-only", "--diff-filter=ACMR")
	if err != nil {
		return err
	}
	return checkTraceChanges(root, files)
}

func TraceRange(root, base string) error {
	if base == "" {
		return errors.New("base ref is required")
	}
	files, err := gitFiles(root, "diff", "--name-only", "--diff-filter=ACMR", base+"...HEAD")
	if err != nil {
		return err
	}
	return checkTraceChanges(root, files)
}

func checkTraceChanges(root string, files []string) error {
	manifest, err := loadTrace(root)
	if err != nil {
		return err
	}
	for _, group := range manifest.Groups {
		source := matches(files, group.DocumentSources)
		generated := matches(files, group.GeneratedDocuments)
		technical := matches(files, append(append([]string{}, group.Code...), group.Tests...))
		if (source || generated || technical) && (source != generated || source != technical) {
			return fmt.Errorf("one-sided document/code change in %s: source=%t generated=%t technical=%t", group.ID, source, generated, technical)
		}
	}
	return nil
}

func gitFiles(root string, args ...string) ([]string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var files []string
	for _, line := range strings.Split(string(out), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			files = append(files, filepath.ToSlash(line))
		}
	}
	sort.Strings(files)
	return files, nil
}

func matches(files, paths []string) bool {
	for _, file := range files {
		for _, path := range paths {
			if file == path || strings.HasPrefix(file, strings.TrimSuffix(path, "/")+"/") {
				return true
			}
		}
	}
	return false
}
