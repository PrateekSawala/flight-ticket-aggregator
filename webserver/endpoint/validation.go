package endpoint

import (
	"net/http"
	"strings"

	"flight-ticket-aggregator/domain"
)

func IsMethodPost(method string) error {
	if method != http.MethodPost {
		return domain.ErrInvalidRequestMethod
	}
	return nil
}

func IsMethodGet(method string) error {
	if method != http.MethodGet {
		return domain.ErrInvalidRequestMethod
	}
	return nil
}

func IsContentTypeFormData(contentType string) error {
	if !strings.Contains(contentType, domain.MimeTypeFormData) {
		return domain.ErrInvalidFile
	}
	return nil
}
