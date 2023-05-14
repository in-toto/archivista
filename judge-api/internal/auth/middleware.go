package auth

import (
	"net/http"

	"gitlab.com/testifysec/judge-platform/judge-api/viewer"
)

func Middleware(authProvider AuthProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookieHeader := r.Header.Get("Cookie")
			v, err := authProvider.ValidateAndGetViewer(r.Context(), cookieHeader)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}

			ctx := viewer.NewContext(r.Context(), v)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
