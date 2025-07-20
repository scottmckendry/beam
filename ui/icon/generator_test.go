package icon

import (
	"strings"
	"testing"
)

func TestPascalToKebab(t *testing.T) {
	cases := map[string]string{
		"Building2":   "building-2",
		"FolderGit2":  "folder-git-2",
		"Trash2":      "trash-2",
		"GitBranch":   "git-branch",
		"CloudUpload": "cloud-upload",
		"A1B2C3":      "a-1-b-2-c-3",
		"ABCD":        "abcd",
	}
	for in, want := range cases {
		got := PascalToKebab(in)
		if got != want {
			t.Errorf("PascalToKebab(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestToPascalCase(t *testing.T) {
	cases := map[string]string{
		"building-2":   "Building2",
		"folder-git-2": "FolderGit2",
		"trash-2":      "Trash2",
		"git-branch":   "GitBranch",
		"cloud-upload": "CloudUpload",
		"a-1-b-2-c-3":  "A1B2C3",
		"abcd":         "Abcd",
	}
	for in, want := range cases {
		got := ToPascalCase(in)
		if got != want {
			t.Errorf("ToPascalCase(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestExtractSVGContent(t *testing.T) {
	svg := `<svg width="24" height="24"><path d="M1 2"/></svg>`
	want := `<path d="M1 2"/>`
	got := ExtractSVGContent(svg)
	if strings.TrimSpace(got) != want {
		t.Errorf("ExtractSVGContent() = %q, want %q", got, want)
	}
}
