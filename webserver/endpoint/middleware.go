package endpoint

import (
	"context"
	"strings"

	"flight-ticket-aggregator/domain"
	"github.com/go-kit/kit/endpoint"
)

type contextKey int

const (
	ContentType = contextKey(iota)
	AcceptHeader
	QueryValues
)

func AcceptContentTypeFormDataValidationMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			acceptHeaderValue := ctx.Value(AcceptHeader).(string)
			if !strings.Contains(acceptHeaderValue, domain.MimeTypeFormData) {
				return nil, domain.ErrInvalidContentType
			}
			return next(ctx, request)
		}
	}
}
