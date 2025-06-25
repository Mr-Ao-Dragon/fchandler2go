package core

import (
	"encoding/base64"
	"io"
)

func param2map(m map[string][]string) *map[string]string {
	m2 := make(map[string]string)
	for k, v := range m {
		m2[k] = ""
		for i, vv := range v {
			m2[k] += vv
			if i != len(v)-1 {
				m2[k] += ","
			}
		}
	}
	return &m2
}
func Header2map(m map[string][]string) *map[string]string {
	m2 := make(map[string]string)
	for k, v := range m {
		m2[Capitalize(k)] = ""
		for i, vv := range v {
			m2[k] += vv
			if i != len(v)-1 {
				m2[k] += ","
			}
		}
	}
	return &m2
}
func BodyReader(b io.ReadCloser, typeof string) (string, error) {
	buffer := make([]byte, 1024)
	var final []byte

	for {
		n, err := b.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			break
		}

		final = append(final, buffer[:n]...)
	}
	err := b.Close()
	if err != nil {
		return "", err
	}
	if IsBase64(typeof) {
		return base64.StdEncoding.EncodeToString(final), nil
	}

	// !!!!!!!!!!!!!!MAY CAUSE ERROR NEED TEST!!!!!!!!!!
	return string(final), nil
}
func StringPtr(s string) *string { return &s }
func IsBin(typeof string) *bool {
	q := !IsBase64(typeof)
	return &q
}
func IsBase64(typeof string) bool {
	textTypes := map[string]bool{
		"text/plain":             true,
		"application/json":       true,
		"application/ld+json":    true,
		"application/xhtml+xml":  true,
		"application/xml":        true,
		"application/atom+xml":   true,
		"application/javascript": true,
		"text/html":              true,
		"text/css":               true,
		"text/csv":               true,
		"text/javascript":        true,
		"text/xml":               true,
	}
	q := !textTypes[typeof]
	return q
}
