package serveroption

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
	"github.com/stretchr/testify/assert"
)

func TestExtractAcceptHeaderIntoContextSuccess(t *testing.T) {
	req := httptest.NewRequest("", "http://localhost:3002", nil)
	req.Header.Set("Accept", "application/json")

	ctx := context.Background()
	textCtx := ExtractAcceptHeaderIntoContext(ctx, req)

	actual := textCtx.Value(domain.HttpAcceptHeader).(string)
	expected := "application/json"

	assert.Equal(t, expected, actual)
}
