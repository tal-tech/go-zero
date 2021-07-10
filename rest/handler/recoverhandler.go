package handler

import (
	"fmt"
	"github.com/tal-tech/go-zero/rest/httpx"
	"github.com/tal-tech/go-zero/rest/internal"
	"net/http"
	"runtime/debug"
)

// RecoverHandler returns a middleware that recovers if panic happens.
func RecoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if result := recover(); result != nil {
				internal.Error(r, fmt.Sprintf("%v\n%s", result, debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
				httpx.Error(w, nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
