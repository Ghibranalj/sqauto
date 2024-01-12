package sqauto

import "strings"

// Assumptions:
// all fields are pascal case since all exported fields are pascal case
func toSnakeCase(s string) string {

	bu := make([]byte, 0, len(s)*2)
	prev := s[0]
	for i := 1; i < len(s); i++ {
		// if prev is lower and current is upper, insert _
		if prev >= 'a' && prev <= 'z' && s[i] >= 'A' && s[i] <= 'Z' {
			bu = append(bu, prev, '_')
		} else {
			bu = append(bu, prev)
		}
		prev = s[i]
	}
	bu = append(bu, prev)
	str := string(bu)
	return strings.ToLower(str)
}
