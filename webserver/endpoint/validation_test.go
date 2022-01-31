package endpoint

import (
	"net/http"
	"testing"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"

	"github.com/stretchr/testify/assert"
)

func TestIsMethodPostValidationSuccessful(t *testing.T) {
	err := IsMethodPost(http.MethodPost)
	assert.NoError(t, err)
}

func TestIsMethodGetValidationSuccessful(t *testing.T) {
	err := IsMethodGet(http.MethodGet)
	assert.NoError(t, err)
}

func TestIsContentTypeFormDataValidationSuccessful(t *testing.T) {
	err := IsContentTypeFormData(domain.MimeTypeFormData)
	assert.NoError(t, err)
}
