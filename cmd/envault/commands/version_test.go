package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionPrintsDefault(t *testing.T) {
	// Reset to known values for test.
	origVersion, origCommit, origDate := BuildVersion, BuildCommit, BuildDate
	defer func() {
		BuildVersion = origVersion
		BuildCommit = origCommit
		BuildDate = origDate
	}()

	BuildVersion = "v1.2.3"
	BuildCommit = "abc1234"
	BuildDate = "2024-01-15"

	cmd := NewVersionCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "v1.2.3") {
		t.Errorf("expected version in output, got: %q", out)
	}
	if !strings.Contains(out, "abc1234") {
		t.Errorf("expected commit in output, got: %q", out)
	}
	if !strings.Contains(out, "2024-01-15") {
		t.Errorf("expected date in output, got: %q", out)
	}
}

func TestVersionShortFlag(t *testing.T) {
	origVersion := BuildVersion
	defer func() { BuildVersion = origVersion }()

	BuildVersion = "v0.9.0"

	cmd := NewVersionCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--short"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	if out != "v0.9.0" {
		t.Errorf("expected %q, got %q", "v0.9.0", out)
	}
}

func TestVersionShortFlagExcludesCommit(t *testing.T) {
	origVersion, origCommit := BuildVersion, BuildCommit
	defer func() {
		BuildVersion = origVersion
		BuildCommit = origCommit
	}()

	BuildVersion = "v2.0.0"
	BuildCommit = "deadbeef"

	cmd := NewVersionCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"-s"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "deadbeef") {
		t.Errorf("short flag should not include commit, got: %q", out)
	}
}
