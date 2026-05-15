package redact_test

import (
	"testing"

	"github.com/envault/envault/internal/redact"
)

func TestMaskStarsReplacesAllValues(t *testing.T) {
	r := redact.New(redact.ModeStars)
	env := map[string]string{
		"SECRET": "supersecret",
		"TOKEN":  "abc123",
	}
	masked := r.Mask(env)
	for k, v := range masked {
		if v != "********" {
			t.Errorf("key %s: expected '********', got %q", k, v)
		}
	}
}

func TestMaskPartialRevealsPrefix(t *testing.T) {
	r := redact.New(redact.ModePartial)
	got := r.MaskValue("mysecretvalue")
	if got[:2] != "my" {
		t.Errorf("expected prefix 'my', got %q", got[:2])
	}
	for _, ch := range got[2:] {
		if ch != '*' {
			t.Errorf("expected '*' after prefix, got %q", string(ch))
		}
	}
}

func TestMaskPartialShortValue(t *testing.T) {
	r := redact.New(redact.ModePartial)
	got := r.MaskValue("ab")
	if got != "**" {
		t.Errorf("expected '**', got %q", got)
	}
}

func TestMaskHashShowsLength(t *testing.T) {
	r := redact.New(redact.ModeHash)
	got := r.MaskValue("hello")
	if got != "<len:5>" {
		t.Errorf("expected '<len:5>', got %q", got)
	}
}

func TestMaskEmptyValuePassthrough(t *testing.T) {
	for _, mode := range []redact.Mode{redact.ModeStars, redact.ModePartial, redact.ModeHash} {
		r := redact.New(mode)
		if got := r.MaskValue(""); got != "" {
			t.Errorf("mode %d: expected empty string, got %q", mode, got)
		}
	}
}

func TestSkipKeysAreNotMasked(t *testing.T) {
	r := redact.New(redact.ModeStars, "APP_ENV", "LOG_LEVEL")
	env := map[string]string{
		"APP_ENV":   "production",
		"LOG_LEVEL": "info",
		"DB_PASS":   "s3cr3t",
	}
	masked := r.Mask(env)
	if masked["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be masked, got %q", masked["APP_ENV"])
	}
	if masked["LOG_LEVEL"] != "info" {
		t.Errorf("LOG_LEVEL should not be masked, got %q", masked["LOG_LEVEL"])
	}
	if masked["DB_PASS"] != "********" {
		t.Errorf("DB_PASS should be masked, got %q", masked["DB_PASS"])
	}
}

func TestSkipKeysAreCaseInsensitive(t *testing.T) {
	r := redact.New(redact.ModeStars, "app_env")
	env := map[string]string{"APP_ENV": "staging"}
	masked := r.Mask(env)
	if masked["APP_ENV"] != "staging" {
		t.Errorf("expected unmasked value, got %q", masked["APP_ENV"])
	}
}
