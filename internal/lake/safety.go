package lake

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var safetyPatterns = []struct {
	name string
	re   *regexp.Regexp
}{
	{"local-user-path", regexp.MustCompile(`/` + `Users/[^/[:space:]]+`)},
	{"github-token", regexp.MustCompile(`gh[opsu]_[A-Za-z0-9]{20,}`)},
	{"google-access-token", regexp.MustCompile(`ya29\.[A-Za-z0-9._-]+`)},
	{"google-refresh-token", regexp.MustCompile(`1//[A-Za-z0-9._-]+`)},
	{"google-client-secret", regexp.MustCompile(`GOC` + `SPX-[A-Za-z0-9_-]+`)},
	{"private-key", regexp.MustCompile(`BEGIN (RSA |OPENSSH |EC |DSA )?PRIVATE KEY`)},
	{"youtube-channel-id", regexp.MustCompile(`\bUC[A-Za-z0-9_-]{22}\b`)},
}

func SafetyCheck(root string) error {
	var findings []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if entry.IsDir() && (rel == ".git" || rel == ".mhj/cache" || strings.HasPrefix(rel, ".git/")) {
			return filepath.SkipDir
		}
		if entry.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for line := 1; scanner.Scan(); line++ {
			for _, pattern := range safetyPatterns {
				if pattern.re.MatchString(scanner.Text()) {
					findings = append(findings, fmt.Sprintf("%s:%d:%s", rel, line, pattern.name))
				}
			}
		}
		return scanner.Err()
	})
	if err != nil {
		return err
	}
	if len(findings) > 0 {
		return errors.New("public-safety findings: " + strings.Join(findings, ", "))
	}
	return nil
}
