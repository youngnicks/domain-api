package vhostapi

import (
	"encoding/json"
	"net/http"

	caddy "github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

var (
	_ caddy.Module                = (*Middleware)(nil)
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterDirective("vhapi", parseHandlerCaddyfile)
}

// Middleware implements an HTTP handler that creates a vhost.
type Middleware struct {
	Vhost
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.vhostapi",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	if err := m.Vhost.provision(ctx, m); err != nil {
		return err
	}

	// only non-routes gets added to the App
	if m.Vhost.isRoute() {
		return nil
	}

	// load or bootstrap App
	appI, err := ctx.App(App{}.CaddyModule().String())
	if err != nil {
		return err
	}
	app := appI.(*App)
	app.addVhost(m.Vhost)
	return nil
}

// Validate implements caddy.Validator
func (m Middleware) Validate() error {
	return m.Vhost.validate()
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	var resp struct {
		Status string `json:"status,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	err := m.run(w, r)

	if err == nil {
		resp.Status = "success"
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err.Error()
	}

	return json.NewEncoder(w).Encode(resp)
}
