package cmd

import (
	"strings"
	"testing"

	"github.com/pterm/pterm"
)

func TestPrintPrompt_RawOutput(t *testing.T) {
	prevRaw := pterm.RawOutput
	pterm.RawOutput = true
	t.Cleanup(func() {
		pterm.RawOutput = prevRaw
	})

	out := captureStdout(t, printPrompt)
	if out != "sim-cli> " {
		t.Fatalf("unexpected prompt: %q", out)
	}
}

func TestPrintWelcomeBanner_RawOutput(t *testing.T) {
	prevRaw := pterm.RawOutput
	pterm.RawOutput = true
	t.Cleanup(func() {
		pterm.RawOutput = prevRaw
	})

	_ = captureStdout(t, printWelcomeBanner)
}

func TestPrintPrompt_NonRawOutput(t *testing.T) {
	prevRaw := pterm.RawOutput
	pterm.RawOutput = false
	t.Cleanup(func() {
		pterm.RawOutput = prevRaw
	})

	out := captureStdout(t, printPrompt)
	if !strings.Contains(out, "sim-cli") {
		t.Fatalf("unexpected prompt: %q", out)
	}
}

func TestPrintWelcomeBanner_NonRawOutput(t *testing.T) {
	prevRaw := pterm.RawOutput
	pterm.RawOutput = false
	t.Cleanup(func() {
		pterm.RawOutput = prevRaw
	})

	_ = captureStdout(t, printWelcomeBanner)
}
