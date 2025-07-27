package acs

import (
	"encoding/base64"
	"testing"
)

func TestNewAttachment(t *testing.T) {
	data := []byte("hello world")
	att := NewAttachment("file.txt", "text/plain", data)
	if att.Name != "file.txt" {
		t.Errorf("Name = %q, want file.txt", att.Name)
	}
	if att.ContentType != "text/plain" {
		t.Errorf("ContentType = %q, want text/plain", att.ContentType)
	}
	if att.ContentInBase64 != base64.StdEncoding.EncodeToString(data) {
		t.Errorf("ContentInBase64 = %q, want %q", att.ContentInBase64, base64.StdEncoding.EncodeToString(data))
	}
}

func TestGenerateContentHash(t *testing.T) {
	input := []byte("abc")
	hash := generateContentHash(input)
	if len(hash) == 0 {
		t.Error("hash is empty")
	}
}
