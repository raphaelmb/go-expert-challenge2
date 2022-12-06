package main

import (
	"context"
	"strings"
)

func SanitizeInput(ctx context.Context, s string) string {
	if strings.Contains(s, "-") {
		if ctx.Value("trim") == "true" {
			return strings.Replace(s, "-", "", 1)
		}
		return s
	} else {
		if ctx.Value("trim") == "false" {
			return s[:5] + "-" + s[5:]
		}
		return s
	}
}
