package endpoint

import (
	"context"
	"testing"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
	"github.com/stretchr/testify/assert"
)

func TestMakeAcceptContentTypeFormDataValidationMiddlewareSuccess(t *testing.T) {
	m := AcceptContentTypeFormDataValidationMiddleware()
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return "", nil
	}
	ctx := context.WithValue(context.Background(), AcceptHeader, domain.MimeTypeFormData)
	_, err := m(ep)(ctx, "")
	assert.NoError(t, err)
}
