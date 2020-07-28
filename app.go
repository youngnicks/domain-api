package vhostapi

import (
	"fmt"
	"sync/atomic"

	caddy "github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"
)

// Interface guards
var (
	_ caddy.App         = (*App)(nil)
	_ caddy.Module      = (*App)(nil)
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.Validator   = (*App)(nil)
)

var lifeCycle int32

func init() {
	caddy.RegisterModule(App{})
}

// App is top level module that runs Custom Domain API.
type App struct {
	Vhost string `json:vhost,omitempty"`

	vhosts map[string]Vhost
	log    *zap.Logger
}

// Provision implements caddy.Provisioner
func (a *App) Provision(ctx caddy.Context) error {
	a.log = ctx.Logger(a)

	return nil
}

// Validate implements caddy.Validator
func (a App) Validate() error {
	if a.RollDomain == "" {
		return fmt.Errorf("roll_domain is required")
	}
	return nil
}

// Start starts the app.
func (a App) Start() error {
	count := atomic.AddInt32(&lifeCycle, 1)
	if count > 1 {
		// not the first startup, maybe a reload
		return nil
	}
	return nil
}

// Stop stops the app.
func (a *App) Stop() error {
	count := atomic.AddInt32(&lifeCycle, -1)
	if count > 0 {
		// not shutdown, maybe a prior config reload.
		return nil
	}
	return nil
}

// CaddyModule implements caddy.ModuleInfo
func (a App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "custom_domain_api",
		New: func() caddy.Module { return new(App) },
	}
}
