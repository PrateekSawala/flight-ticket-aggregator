package middleware

import (
	"context"
	"strings"

	"flight-ticket-aggregator/domain"
	"github.com/go-kit/kit/endpoint"
)

func MakeAcceptContentTypeFormDataValidationMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			acceptHeaderValue := ctx.Value(domain.HttpAcceptHeader).(string)
			if !strings.Contains(acceptHeaderValue, domain.MimeTypeFormData) {
				return nil, domain.ErrInvalidContentType
			}
			return next(ctx, request)
		}
	}
}
