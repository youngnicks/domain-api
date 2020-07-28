package vhostapi

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
)

// parseHandlerCaddyfile unmarshals tokens from h into a new Middleware.
func parseHandlerCaddyfile(h httpcaddyfile.Helper) ([]httpcaddyfile.ConfigValue, error) {
	if !h.Next() { // No more tokens to process.
		return nil, h.ArgErr()
	}

	var v Vhost

	// Get first token looking for a matcher token.
	// If token is matcher, ok will be true, otherwise false.
	matcherSet, ok, err := h.MatcherToken()
	if err != nil {
		return nil, err
	}
	if ok { // Token is a matcher
		// Remove matcher token from token slice. This allows us to advance to the next token.
		h.Dispenser.Delete()
	} else { // Token is not a matcher
		// First token must be a matcher for this app.
		return nil, h.ArgErr()
	}

	// Parse remaining tokens from Caddyfile using vhost method.
	err = v.UnmarshalCaddyfile(h.Dispenser)
	if err != nil {
		return nil, err
	}

	// Create new route
	m := Middleware{Vhost: v}
	return h.NewRoute(matcherSet, m), nil
}

// UnmarshalCaddyfile configures the global directive from Caddyfile.
// Syntax:
//
// vhostapi [<matcher>] [<template>] [args...] {
//     template  <text>
//     args      <text>...
// }
//
// UnmarshalCaddyfile parses the inline arguments [<template>] [args...].
func (m *Vhost) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	// vhost, if present
	if d.Next() {
		if !d.Args(&m.Template) { // Read next argument into vhost Template
			return d.ArgErr()
		}
	}
	// Everything else, if present, are args.
	m.Args = d.RemainingArgs()

	// Parse the next block.
	return m.unmarshalBlock(d)
}

// Extract vhost config from the configuration block {template, args} if defined.
func (m *Vhost) unmarshalBlock(d *caddyfile.Dispenser) error {
	for d.NextBlock(0) { // Loop over tokens starting with the first.
		switch d.Val() {
		case "template":
			if m.Template != "" {
				return d.Err("template specified twice")
			}
			if !d.args(&m.Template) { // Read argument into vhost Template
				return d.ArgErr()
			}
		case "args":
			if len(m.Args) > 0 {
				return d.Err("args specified twice")
			}
			m.Args = d.RemainingArgs() // Everything else are arguments
		}
	}

	return nil
}
