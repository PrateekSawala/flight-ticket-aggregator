package middleware

import (
	"context"
	"testing"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
	"github.com/stretchr/testify/assert"
)

func TestMakeAcceptContentTypeFormDataValidationMiddlewareSuccess(t *testing.T) {
	m := MakeAcceptContentTypeFormDataValidationMiddleware()
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return "", nil
	}
	ctx := context.WithValue(context.Background(), domain.HttpAcceptHeader, domain.MimeTypeFormData)
	_, err := m(ep)(ctx, "")
	assert.NoError(t, err)
}
