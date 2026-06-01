// Package gistid normalizes gist identifiers.
package gistid

// Sanitize removes characters that are not ASCII letters, digits, or spaces.
func Sanitize(id string) string {
	bytesID := []byte(id)
	index := 0

	for _, char := range bytesID {
		if ('a' <= char && char <= 'z') ||
			('A' <= char && char <= 'Z') ||
			('0' <= char && char <= '9') ||
			char == ' ' {
			bytesID[index] = char
			index++
		}
	}

	return string(bytesID[:index])
}
