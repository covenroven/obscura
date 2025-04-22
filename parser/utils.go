package parser

import "strings"

// cleanText clean up invalid characters and replace them with dash
func cleanText(txt string) string {
	txt = strings.ReplaceAll(txt, "/", "-")
	txt = strings.ReplaceAll(txt, "\\", "-")

	return strings.TrimSpace(txt)
}
