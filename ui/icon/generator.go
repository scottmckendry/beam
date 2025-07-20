package icon

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	iconVarRe = regexp.MustCompile(`icon\.([A-Z][A-Za-z0-9]*)`)
	skip      = map[string]struct{}{
		"Props": {}, "Icon": {}, "GenerateSVG": {}, "GetIconContent": {}, "Render": {},
	}
)

// ScanUsedIcons scans the given directories and files for icon usage and returns a sorted list of kebab-case icon names.
func ScanUsedIcons(searchDirs, searchFiles []string) []string {
	iconSet := make(map[string]struct{})
	for _, dir := range searchDirs {
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || (!strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".templ")) {
				return nil
			}
			processFile(path, iconSet)
			return nil
		})
	}
	for _, file := range searchFiles {
		processFile(file, iconSet)
	}
	var icons []string
	for icon := range iconSet {
		icons = append(icons, icon)
	}
	sort.Strings(icons)
	return icons
}

func processFile(path string, iconSet map[string]struct{}) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		for _, m := range iconVarRe.FindAllStringSubmatch(line, -1) {
			if _, reserved := skip[m[1]]; reserved {
				continue
			}
			iconSet[PascalToKebab(m[1])] = struct{}{}
		}
	}
}

// PascalToKebab converts a PascalCase string to kebab-case, inserting dashes before digits and uppercase letters as needed.
func PascalToKebab(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 {
			prev := rune(s[i-1])
			// Insert dash before uppercase letter (if previous is lower or digit)
			if r >= 'A' && r <= 'Z' && ((prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9')) {
				result = append(result, '-')
			}
			// Insert dash before digit if previous is a letter
			if r >= '0' && r <= '9' && ((prev >= 'a' && prev <= 'z') || (prev >= 'A' && prev <= 'Z')) {
				result = append(result, '-')
			}
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// ToPascalCase converts a kebab-case string to PascalCase.
func ToPascalCase(s string) string {
	words := strings.Split(s, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, "")
}

// ExtractSVGContent extracts the inner SVG markup from a full SVG string.
func ExtractSVGContent(svgContent string) string {
	start := strings.Index(svgContent, ">") + 1
	end := strings.LastIndex(svgContent, "</svg>")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return strings.TrimSpace(svgContent[start:end])
}

// DownloadFile fetches the content at the given URL and returns it as a byte slice.
func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("download returned status code %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// GenerateIcons downloads SVGs for the given icon names and writes Go source files for use in the icon package.
// Returns a list of missing icons (failed downloads).
func GenerateIcons(icons []string, outputDir, lucideVersion string) (missingIcons []string, err error) {
	os.MkdirAll(outputDir, os.ModePerm)

	iconDefs := []string{"package icon\n", "// This file is auto generated\n", fmt.Sprintf("// Using Lucide icons version %s\n", lucideVersion)}
	iconDataEntries := make(map[string]string)
	missingIcons = []string{}

	for i, name := range icons {
		funcName := ToPascalCase(name)
		url := fmt.Sprintf("https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/%s.svg", name)
		contentBytes, err := DownloadFile(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", name, err)
			missingIcons = append(missingIcons, name)
			iconDefs = append(iconDefs, fmt.Sprintf("// MISSING: %s (failed to download)\n", funcName))
			continue
		}
		iconDefs = append(iconDefs, fmt.Sprintf("var %s = Icon(\"%s\")\n", funcName, name))
		innerContent := ExtractSVGContent(string(contentBytes))
		iconDataEntries[name] = innerContent
		fmt.Printf("[%d/%d] %s\n", i+1, len(icons), name)
	}

	// Write icondefs.go
	outputFileDefs := filepath.Join(outputDir, "icondefs.go")
	os.WriteFile(outputFileDefs, []byte(strings.Join(iconDefs, "")), 0644)

	// Write icondata.go
	outputFileData := filepath.Join(outputDir, "icondata.go")
	var iconDataContent strings.Builder
	iconDataContent.WriteString("package icon\n\n")
	iconDataContent.WriteString("// This file is auto generated\n")
	iconDataContent.WriteString(fmt.Sprintf("// Using Lucide icons version %s\n\n", lucideVersion))
	iconDataContent.WriteString(fmt.Sprintf("const LucideVersion = %q\n\n", lucideVersion))
	iconDataContent.WriteString("var internalSvgData = map[string]string{\n")
	for name, data := range iconDataEntries {
		escapedData := strings.ReplaceAll(data, "`", "`+\"`\"+`")
		iconDataContent.WriteString(fmt.Sprintf("\t%q: `%s`,\n", name, escapedData))
	}
	iconDataContent.WriteString("}\n")
	os.WriteFile(outputFileData, []byte(iconDataContent.String()), 0644)

	return missingIcons, nil
}
