package utility

import (
	"bufio"
	"io"
	"net/http"
	"strings"
	"unicode"

	"flight-ticket-aggregator/domain"
)

func IsMethodPost(method string) error {
	if method != http.MethodPost {
		return domain.ErrInvalidRequest
	}
	return nil
}

func IsFileCSV(fileread io.Reader) error {
	bufferReader := bufio.NewReader(fileread)
	initialByte, _ := bufferReader.Peek(512)
	contentType := http.DetectContentType(initialByte)

	// Check content Type
	if !strings.Contains(contentType, domain.TypeFileContentText) {
		return domain.ErrInvalidFile
	}
	return nil
}

func IsAlphabetic(value string) bool {
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
