package icon

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/a-h/templ"
)

var (
	iconContents = make(map[string]string)
	iconMutex    sync.RWMutex
)

type Props struct {
	Size        int
	Color       string
	Fill        string
	Stroke      string
	StrokeWidth string // Stroke Width of Icon, Usage: "2.5"
	Class       string
}

func Icon(name string) func(...Props) templ.Component {
	return func(props ...Props) templ.Component {
		var p Props
		if len(props) > 0 {
			p = props[0]
		}
		cacheKey := fmt.Sprintf("%s|s:%d|c:%s|f:%s|sk:%s|sw:%s|cl:%s",
			name, p.Size, p.Color, p.Fill, p.Stroke, p.StrokeWidth, p.Class)
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			iconMutex.RLock()
			svg, cached := iconContents[cacheKey]
			iconMutex.RUnlock()
			if cached {
				_, err = w.Write([]byte(svg))
				return err
			}
			generatedSvg, err := generateSVG(name, p)
			if err != nil {
				return fmt.Errorf("failed to generate svg for icon '%s' with props %+v: %w", name, p, err)
			}
			iconMutex.Lock()
			iconContents[cacheKey] = generatedSvg
			iconMutex.Unlock()
			_, err = w.Write([]byte(generatedSvg))
			return err
		})
	}
}

func generateSVG(name string, props Props) (string, error) {
	content, err := getIconContent(name)
	if err != nil {
		return "", err
	}
	size := props.Size
	if size <= 0 {
		size = 24
	}
	fill := props.Fill
	if fill == "" {
		fill = "none"
	}
	stroke := props.Stroke
	if stroke == "" {
		stroke = props.Color
	}
	if stroke == "" {
		stroke = "currentColor"
	}
	strokeWidth := props.StrokeWidth
	if strokeWidth == "" {
		strokeWidth = "2"
	}
	return fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\" viewBox=\"0 0 24 24\" fill=\"%s\" stroke=\"%s\" stroke-width=\"%s\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"%s\" data-lucide=\"icon\">%s</svg>",
		size, size, fill, stroke, strokeWidth, props.Class, content), nil
}

func getIconContent(name string) (string, error) {
	content, exists := internalSvgData[name]
	if !exists {
		return "", fmt.Errorf("icon '%s' not found in internalSvgData map", name)
	}
	return content, nil
}
