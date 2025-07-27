package utils

import (
	"encoding/base64"
	"path"
	"strings"
)

// DecodeBase64Image decodes a base64 string, stripping any data URL prefix if present.
func DecodeBase64Image(s string) ([]byte, error) {
	if idx := strings.Index(s, ","); idx != -1 {
		s = s[idx+1:]
	}
	return base64.StdEncoding.DecodeString(s)
}

// GetImageExtension returns the file extension for an image based on MIME type or filename.
func GetImageExtension(mimes []string, names []string, defaultExt string) string {
	ext := defaultExt
	if len(mimes) > 0 && mimes[0] != "" {
		switch mimes[0] {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		}
	} else if len(names) > 0 && names[0] != "" {
		ext = path.Ext(names[0])
		if ext == "" {
			ext = defaultExt
		}
	}
	return ext
}
