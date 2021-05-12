package server

import (
	"net/http"
	"strings"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/logging"
)

func AcceptContentTypeValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logging.Log("AcceptContentTypeValidationMiddleware")
		acceptContentType := r.Header.Get(domain.ContentType)
		if !strings.Contains(acceptContentType, domain.MimeFormData) {
			log.Debugf("Invalid content type %s", acceptContentType)
			http.Error(w, domain.ErrInvalidContentType.Error(), http.StatusUnauthorized)
			return
		}
		// on to the next handler
		next.ServeHTTP(w, r)
	})
}
