/*
This file
*/

package vhostapi

import (
	"fmt"
	"os"

	caddy "github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"
)

// Vhost is the module configuration
type Vhost struct {
	// The template file to use.
	Template string `json:"template,omitempty"`

	// The template args.
	Args map[string]interface{} `json:"args,omitempty"`

	log *zap.Logger
}

// Provision implements caddy.Provisioner
func (m *Vhost) provision(ctx caddy.Context, cm caddy.Module) error {
	m.log = ctx.Logger(cm)

	return nil
}

// Validate implements caddy.Validator.
func (m Vhost) validate() error {
	// Check for missing template variable
	if m.Template == "" {
		return fmt.Errorf("template is required")
	}

	// Validate template file exists
	if err := isValidFile(m.Template); err != nil {
		return err
	}

	return nil
}

func isValidFile(file string) error {
	s, err := os.Stat(file)

	// Check if file exists
	if err != nil {
		return err
	}

	// Check if file is a directory
	if s.IsDir() {
		return fmt.Errorf("not a file '%s'", file)
	}

	return nil
}
