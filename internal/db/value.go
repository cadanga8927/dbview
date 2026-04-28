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
var indentSpaces = strings.Repeat("  ", 32)

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
				buf.WriteString(indentSpaces[:indent*2])
			}
		case '}', ']':
			if !inStr {
				buf.WriteString("\n")
				indent--
				buf.WriteString(indentSpaces[:indent*2])
			}
			buf.WriteRune(r)
		case '"':
			buf.WriteRune(r)
			inStr = !inStr
		case ',':
			buf.WriteRune(r)
			if !inStr {
				buf.WriteString("\n")
				buf.WriteString(indentSpaces[:indent*2])
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

const hexTable = "0123456789abcdef"

func hexDump(data []byte) string {
	limit := min(len(data), 64)
	buf := make([]byte, 0, limit*3-1)
	for i := range limit {
		if i > 0 {
			buf = append(buf, ' ')
		}
		buf = append(buf, hexTable[data[i]>>4], hexTable[data[i]&0x0f])
	}
	s := string(buf)
	if len(data) > 64 {
		s += " ..."
	}
	return "[" + s + "]"
}
