// Package redact provides utilities for masking sensitive values
// in .env content before displaying or logging them.
package redact

import (
	"strings"
)

// Mode controls how values are redacted.
type Mode int

const (
	// ModeStars replaces the value with a fixed number of asterisks.
	ModeStars Mode = iota
	// ModePartial reveals the first two characters and masks the rest.
	ModePartial
	// ModeHash shows only the length of the value.
	ModeHash
)

// Redactor masks env var values according to a chosen mode.
type Redactor struct {
	mode Mode
	skip map[string]bool
}

// New creates a Redactor with the given mode.
// Keys listed in skip are left unmasked.
func New(mode Mode, skip ...string) *Redactor {
	s := make(map[string]bool, len(skip))
	for _, k := range skip {
		s[strings.ToUpper(k)] = true
	}
	return &Redactor{mode: mode, skip: s}
}

// Mask returns a redacted copy of the provided key=value map.
func (r *Redactor) Mask(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.skip[strings.ToUpper(k)] {
			out[k] = v
			continue
		}
		out[k] = r.maskValue(v)
	}
	return out
}

// MaskValue masks a single value according to the configured mode.
func (r *Redactor) MaskValue(v string) string {
	return r.maskValue(v)
}

func (r *Redactor) maskValue(v string) string {
	if v == "" {
		return ""
	}
	switch r.mode {
	case ModePartial:
		if len(v) <= 2 {
			return strings.Repeat("*", len(v))
		}
		return v[:2] + strings.Repeat("*", len(v)-2)
	case ModeHash:
		return fmt.Sprintf("<len:%d>", len(v))
	default: // ModeStars
		return "********"
	}
}
