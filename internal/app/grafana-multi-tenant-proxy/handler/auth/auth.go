package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/config"
)

type key int

const (
	// OrgIDKey Key used to pass tenant id though the middleware context
	OrgIDKey key = iota
)

// INTERFACE to handle different type of authentication
type Authenticator interface {
	Authenticate(r *http.Request, targetServer *config.TargetServer) (bool, string)
	OnAuthenticationError(w http.ResponseWriter)
}

type AuthenticationMiddleware struct {
	handler http.HandlerFunc
	config  *config.Config
	logger  *zap.Logger
}

<<<<<<< HEAD:internal/app/grafana-multi-tenant-proxy/handler/auth/auth.go
func NewAuthenticationMiddleware(config *config.Config, logger *zap.Logger, handler http.HandlerFunc) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		handler: handler,
		config:  config,
=======
func NewAuthenticationMiddleware(config config.Config, logger *zap.Logger, handler http.HandlerFunc) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		handler: handler,
		config:  &config,
>>>>>>> 2eb33b2 (Improve config management):internal/app/grafana-multi-tenant-proxy/auth/auth.go
		logger:  logger,
	}
}

// ////////////////////////////////////////////////////////////////////////////////////
// Authenticate can be used as a middleware chain to authenticate every request before proxying the request
func (am AuthenticationMiddleware) Authenticate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticator, err := newAuthenticator(r, am.config, am.logger)
		if err != nil {
			am.logger.Error("Error while authenticating request", zap.String("url", r.URL.String()), zap.Error(err))
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised\n"))
			return
		}

		targetServer := am.config.Proxy.FindTargetServer(r.Host)
<<<<<<< HEAD
		if targetServer == nil {
			am.logger.Error("Target server not configured",
				zap.String("host", r.Host),
				zap.String("url", r.URL.String()),
				zap.Error(err),
			)
			w.WriteHeader(404)
			w.Write([]byte("Not found\n"))
			return
=======
		if targetServer != nil {
			am.logger.Error("Target server not configured")
			w.WriteHeader(404)
			w.Write([]byte("Not found\n"))
>>>>>>> e5eff05 (support-multiple-hosts-from-one-config)
		}

		am.logger.Debug(fmt.Sprintf("Authentication mode: %T", authenticator))
		ok, orgID := authenticator.Authenticate(r, targetServer)
		if !ok {
			authenticator.OnAuthenticationError(w)
			return
		}
		ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
		am.handler(w, r.WithContext(ctx))
	}
}

<<<<<<< HEAD:internal/app/grafana-multi-tenant-proxy/handler/auth/auth.go
func (am AuthenticationMiddleware) ApplyConfig(config *config.Config) {
	*am.config = *config
=======
func (am AuthenticationMiddleware) ApplyConfig(config config.Config) {
	*am.config = config
>>>>>>> 2eb33b2 (Improve config management):internal/app/grafana-multi-tenant-proxy/auth/auth.go
}

// newAuthenticator returns the authentication mode used by the request and its credentials
func newAuthenticator(r *http.Request, config *config.Config, logger *zap.Logger) (Authenticator, error) {
	// OAuth token is favorite authentication mode
	token := r.Header.Get("X-Id-Token")
	if token != "" {
		return OAuthAuthenticator{
			token:  token,
			config: config,
			logger: logger,
		}, nil
	}
	// If no oauth token, we are looking for basicAuth
	user, pwd, ok := r.BasicAuth()
	if ok {
		return BasicAuthenticator{
			user:   user,
			pwd:    pwd,
			config: config,
			logger: logger,
		}, nil
	}
	return nil, errors.New("unsupported authentication")
}
