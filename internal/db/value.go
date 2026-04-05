package db

import (
	"fmt"
	"strings"
)

// FormatValue intelligently formats a value for TUI display.
// It detects JSON, binary data, and long strings and formats them appropriately.
func FormatValue(v interface{}) string {
	if v == nil {
		return "NULL"
	}

	switch val := v.(type) {
	case []byte:
		// Check if it's valid UTF-8 (printable text)
		if isPrintable(val) {
			s := string(val)
			// Try to detect JSON
			if strings.HasPrefix(strings.TrimSpace(s), "{") || strings.HasPrefix(strings.TrimSpace(s), "[") {
				return formatJSON(s)
			}
			// Printable bytes as string, truncate if long
			if len(s) > 200 {
				return s[:197] + "..."
			}
			return s
		}
		// Binary data - show hex dump
		return hexDump(val)
	case string:
		s := val
		if strings.TrimSpace(s) == "" {
			return s
		}
		// Try to detect JSON in strings
		trimmed := strings.TrimSpace(s)
		if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
			return formatJSON(s)
		}
		// Truncate long strings
		if len(s) > 200 {
			return s[:197] + "..."
		}
		return s
	default:
		s := fmt.Sprintf("%v", v)
		if len(s) > 200 {
			return s[:197] + "..."
		}
		return s
	}
}

// isPrintable checks if all bytes in the slice are printable UTF-8.
func isPrintable(data []byte) bool {
	for _, b := range data {
		if b < 32 && b != '\t' && b != '\n' && b != '\r' {
			return false
		}
	}
	return true
}

// formatJSON attempts to pretty-print a JSON string.
func formatJSON(s string) string {
	var buf strings.Builder
	indent := 0
	inStr := false
	for _, r := range s {
		switch r {
		case '{', '[':
			buf.WriteRune(r)
			if !inStr {
				buf.WriteString("\n")
				indent++
				buf.WriteString(strings.Repeat("  ", indent))
			}
		case '}', ']':
			if !inStr {
				buf.WriteString("\n")
				indent--
				buf.WriteString(strings.Repeat("  ", indent))
			}
			buf.WriteRune(r)
		case '"':
			buf.WriteRune(r)
			inStr = !inStr
		case ',':
			buf.WriteRune(r)
			if !inStr {
				buf.WriteString("\n")
				buf.WriteString(strings.Repeat("  ", indent))
			}
		default:
			buf.WriteRune(r)
		}
	}
	result := buf.String()
	// Truncate if too long for cell display
	if len(result) > 200 {
		return result[:197] + "..."
	}
	return result
}

// hexDump formats binary data as a hex dump.
func hexDump(data []byte) string {
	var parts []string
	limit := len(data)
	if limit > 64 {
		limit = 64
	}
	for i := 0; i < limit; i++ {
		parts = append(parts, fmt.Sprintf("%02x", data[i]))
	}
	s := strings.Join(parts, " ")
	if len(data) > 64 {
		s += " ..."
	}
	return "[" + s + "]"
}
