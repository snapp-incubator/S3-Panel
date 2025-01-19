package server

import "strings"

func isServerAuthEnabled(s string) bool {
	// consider it would be disabled by default
	if strings.ToLower(s) == "true" {
		return true
	}
	return false
}
