package handler

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/dgrijalva/jwt-go"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"github.com/tal-tech/go-zero/rest/token"
)

const (
	jwtAudience    = "aud"
	jwtExpire      = "exp"
	jwtId          = "jti"
	jwtIssueAt     = "iat"
	jwtIssuer      = "iss"
	jwtNotBefore   = "nbf"
	jwtSubject     = "sub"
	noDetailReason = "no detail reason"
)

var (
	errInvalidToken = errors.New("invalid auth token")
	errNoClaims     = errors.New("no auth params")
)

type (
	// A AuthorizeOptions is authorize options.
	AuthorizeOptions struct {
		PrevSecret string
		Callback   UnauthorizedCallback
	}

	// UnauthorizedCallback defines the method of unauthorized callback.
	UnauthorizedCallback func(w http.ResponseWriter, r *http.Request, err error)
	// AuthorizeOption defines the method to customize an AuthorizeOptions.
	AuthorizeOption func(opts *AuthorizeOptions)
)

// Authorize returns an authorize middleware.
func Authorize(secret string, opts ...AuthorizeOption) func(http.Handler) http.Handler {
	var authOpts AuthorizeOptions
	for _, opt := range opts {
		opt(&authOpts)
	}

	parser := token.NewTokenParser()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tok, err := parser.ParseToken(r, secret, authOpts.PrevSecret)
			if err != nil {
				unauthorized(w, r, err, authOpts.Callback)
				return
			}

			if !tok.Valid {
				unauthorized(w, r, errInvalidToken, authOpts.Callback)
				return
			}

			claims, ok := tok.Claims.(jwt.MapClaims)
			if !ok {
				unauthorized(w, r, errNoClaims, authOpts.Callback)
				return
			}

			ctx := r.Context()
			for k, v := range claims {
				switch k {
				case jwtAudience, jwtExpire, jwtId, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
					// ignore the standard claims
				default:
					ctx = context.WithValue(ctx, k, v)
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WithPrevSecret returns an AuthorizeOption with setting previous secret.
func WithPrevSecret(secret string) AuthorizeOption {
	return func(opts *AuthorizeOptions) {
		opts.PrevSecret = secret
	}
}

// WithUnauthorizedCallback returns an AuthorizeOption with setting unauthorized callback.
func WithUnauthorizedCallback(callback UnauthorizedCallback) AuthorizeOption {
	return func(opts *AuthorizeOptions) {
		opts.Callback = callback
	}
}

func detailAuthLog(r *http.Request, reason string) {
	// discard dump error, only for debug purpose
	details, _ := httputil.DumpRequest(r, true)
	logx.Errorf("authorize failed: %s\n=> %+v", reason, string(details))
}

func unauthorized(w http.ResponseWriter, r *http.Request, err error, callback UnauthorizedCallback) {
	writer := newGuardedResponseWriter(w)

	if err != nil {
		detailAuthLog(r, err.Error())
	} else {
		detailAuthLog(r, noDetailReason)
	}
	if callback != nil {
		callback(writer, r, err)
	}

	httpx.Error(w, err)
	writer.WriteHeader(http.StatusUnauthorized)
}

type guardedResponseWriter struct {
	writer      http.ResponseWriter
	wroteHeader bool
}

func newGuardedResponseWriter(w http.ResponseWriter) *guardedResponseWriter {
	return &guardedResponseWriter{
		writer: w,
	}
}

func (grw *guardedResponseWriter) Flush() {
	if flusher, ok := grw.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (grw *guardedResponseWriter) Header() http.Header {
	return grw.writer.Header()
}

// Hijack implements the http.Hijacker interface.
// This expands the Response to fulfill http.Hijacker if the underlying http.ResponseWriter supports it.
func (grw *guardedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := grw.writer.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

func (grw *guardedResponseWriter) Write(body []byte) (int, error) {
	return grw.writer.Write(body)
}

func (grw *guardedResponseWriter) WriteHeader(statusCode int) {
	if grw.wroteHeader {
		return
	}

	grw.wroteHeader = true
	grw.writer.WriteHeader(statusCode)
}
