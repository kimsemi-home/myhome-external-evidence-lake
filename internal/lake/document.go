package lake

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type document struct {
	SchemaVersion string    `json:"schema_version"`
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Output        string    `json:"output"`
	Summary       string    `json:"summary"`
	Sections      []section `json:"sections"`
}

type section struct {
	Heading      string   `json:"heading"`
	Paragraphs   []string `json:"paragraphs,omitempty"`
	Bullets      []string `json:"bullets,omitempty"`
	CodeLanguage string   `json:"code_language,omitempty"`
	Code         string   `json:"code,omitempty"`
}

func GenerateDocuments(root string) error {
	doc, err := loadDocument(root)
	if err != nil {
		return err
	}
	output := filepath.Join(root, filepath.FromSlash(doc.Output))
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return err
	}
	return os.WriteFile(output, renderDocument(doc), 0o644)
}

func CheckDocuments(root string) error {
	doc, err := loadDocument(root)
	if err != nil {
		return err
	}
	body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(doc.Output)))
	if err != nil {
		return err
	}
	if !bytes.Equal(body, renderDocument(doc)) {
		return fmt.Errorf("generated document is stale: %s", doc.Output)
	}
	return nil
}

func loadDocument(root string) (document, error) {
	var doc document
	path := filepath.Join(root, "docs-src", "shorts-trend-evidence.json")
	if err := decodeStrict(path, &doc); err != nil {
		return document{}, err
	}
	if doc.SchemaVersion != "repo.document/v1" || doc.ID == "" || doc.Title == "" || doc.Summary == "" || len(doc.Sections) == 0 {
		return document{}, errors.New("document source is incomplete")
	}
	if filepath.IsAbs(doc.Output) || filepath.Clean(doc.Output) != doc.Output || strings.HasPrefix(doc.Output, "..") || !strings.HasSuffix(doc.Output, ".md") {
		return document{}, errors.New("document output is unsafe")
	}
	return doc, nil
}

func renderDocument(doc document) []byte {
	var out bytes.Buffer
	fmt.Fprintf(&out, "# %s\n\n%s\n", doc.Title, doc.Summary)
	for _, section := range doc.Sections {
		fmt.Fprintf(&out, "\n## %s\n", section.Heading)
		for _, paragraph := range section.Paragraphs {
			fmt.Fprintf(&out, "\n%s\n", paragraph)
		}
		if len(section.Bullets) > 0 {
			out.WriteByte('\n')
			for _, item := range section.Bullets {
				fmt.Fprintf(&out, "- %s\n", item)
			}
		}
		if section.Code != "" {
			fmt.Fprintf(&out, "\n```%s\n%s\n```\n", section.CodeLanguage, strings.TrimSuffix(section.Code, "\n"))
		}
	}
	return out.Bytes()
}
