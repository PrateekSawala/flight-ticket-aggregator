package serveroption

import (
	"context"
	"net/http"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
)

func ExtractAcceptHeaderIntoContext(ctx context.Context, r *http.Request) context.Context {
	if acceptHeaderValue := r.Header.Get(domain.HttpContentTypeHeaderContentType); acceptHeaderValue != "" {
		return context.WithValue(ctx, domain.HttpAcceptHeader, acceptHeaderValue)
	}
	return context.WithValue(ctx, domain.HttpAcceptHeader, domain.MimeTypeJsonFormData)
}
