package system

import (
	"bufio"
	"io"
	"net/http"
	"strings"
	"unicode"

	"flight-ticket-aggregator/domain"
)

func IsAlphabetic(value string) bool {
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsFileCSV(fileread io.Reader) error {
	bufferReader := bufio.NewReader(fileread)
	initialByte, _ := bufferReader.Peek(512)
	contentType := http.DetectContentType(initialByte)
	if !strings.Contains(contentType, domain.TypeFileContentText) {
		return domain.ErrInvalidFile
	}
	return nil
}
