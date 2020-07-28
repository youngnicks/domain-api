package vhostapi

import (
	"net/http"

	"go.uber.org/zap"
)

// API runs vhost api.
type API interface {
	Run() error
}

type runnerFunc func() error

func (r runnerFunc) Run() error { return r() }

func (m *Vhost) run(w http.ResponseWriter, r *http.Request) error {
	apiInfo := zap.Any("custom domain", append([]string{m.Template}, m.Args...))
	log := m.log.With(apiInfo)

	log.Info()

	return nil
}
