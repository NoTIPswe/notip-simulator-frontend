package cmd

import (
	"reflect"
	"testing"

	"github.com/pterm/pterm"
)

func TestStartSpinner_RawOutputReturnsNoop(t *testing.T) {
	prevRaw := pterm.RawOutput
	pterm.RawOutput = true
	t.Cleanup(func() {
		pterm.RawOutput = prevRaw
	})

	sp := startSpinner("test")
	if _, ok := sp.(noopSpinner); !ok {
		t.Fatalf("expected noopSpinner in raw output mode, got %T", sp)
	}

	// No-op methods should be safe to call.
	sp.Success("ok")
	sp.Fail("err")
	sp.Warning("warn")
}

func TestPtermSpinner_NilInnerIsSafe(t *testing.T) {
	sp := ptermSpinner{}

	sp.Success("ok")
	sp.Fail("err")
	sp.Warning("warn")
}

func TestNoopSpinner_DirectMethods(t *testing.T) {
	var sp noopSpinner

	sp.Success("ok")
	sp.Fail("err")
	sp.Warning("warn")
}

func TestStartSpinner_NonRawReturnsPtermSpinner(t *testing.T) {
	prevRaw := pterm.RawOutput
	pterm.RawOutput = false
	t.Cleanup(func() {
		pterm.RawOutput = prevRaw
	})

	sp := startSpinner("test")

	if reflect.TypeOf(sp) != reflect.TypeOf(ptermSpinner{}) {
		t.Fatalf("expected ptermSpinner in non-raw mode, got %T", sp)
	}

	sp.Success("ok")
}

func TestPtermSpinner_WithInnerIsSafe(t *testing.T) {
	prevRaw := pterm.RawOutput
	pterm.RawOutput = false
	t.Cleanup(func() {
		pterm.RawOutput = prevRaw
	})

	raw := startSpinner("test")
	sp, ok := raw.(ptermSpinner)
	if !ok {
		t.Skipf("expected ptermSpinner in non-raw mode, got %T", raw)
	}

	sp.Success("ok")
	sp.Fail("err")
	sp.Warning("warn")
}
