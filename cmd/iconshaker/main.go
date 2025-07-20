package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/scottmckendry/beam/ui/icon"
)

var (
	iconVarRe   = regexp.MustCompile(`icon\.([A-Z][A-Za-z0-9]*)`) // icon.SomeIcon
	searchDirs  = []string{"ui", "handlers", "acs"}
	searchFiles = []string{"main.go"}
	skip        = map[string]struct{}{
		"Props": {}, "Icon": {}, "GenerateSVG": {}, "GetIconContent": {}, "Render": {},
	}
)

// main is the entry point for the icon generator CLI.
func main() {
	// Minimal CLI: scan for used icons and generate them
	icons := icon.ScanUsedIcons([]string{"ui", "handlers", "acs"}, []string{"main.go"})
	if len(icons) == 0 {
		fmt.Fprintln(os.Stderr, "No icons found to generate.")
		os.Exit(1)
	}
	missing, err := icon.GenerateIcons(icons, "ui/icon", "0.525.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating icons: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nRequested %d icons, generated %d icons in ui/icon\n", len(icons), len(icons)-len(missing))
	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "\nWARNING: The following icons could not be generated (missing or download failed):\n")
		for _, name := range missing {
			fmt.Fprintf(os.Stderr, "  - %s\n", name)
		}
	}
}

func readIconListFile(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open icon list file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	var icons []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		icon := strings.TrimSpace(scanner.Text())
		if icon != "" {
			icons = append(icons, icon)
		}
	}
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
			iconSet[icon.PascalToKebab(m[1])] = struct{}{}
		}
	}
}

// End of file
