package middlewares

import (
	"context"
	"net/http"

	"github.com/pavlegich/gophkeeper/internal/infra/hash"
	"github.com/pavlegich/gophkeeper/internal/server/utils"
)

// WithAuth checks and validates authorization token.
func WithAuth(token *hash.Token) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI == "/api/user/register" || r.RequestURI == "/api/user/login" ||
				r.RequestURI == "/" {
				h.ServeHTTP(w, r)
				return
			}
			cookie, err := r.Cookie("auth")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			id, err := token.Validate(cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), utils.ContextIDKey, id)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
