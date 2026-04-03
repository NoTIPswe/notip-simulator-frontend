package cmd

import "github.com/pterm/pterm"

type spinner interface {
	Success(text string)
	Fail(text string)
	Warning(text string)
}

type noopSpinner struct{}

func (noopSpinner) Success(string) {
	// no-op: spinner is disabled in raw output mode
}

func (noopSpinner) Fail(string) {
	// no-op: spinner is disabled in raw output mode
}

func (noopSpinner) Warning(string) {
	// no-op: spinner is disabled in raw output mode
}

type ptermSpinner struct {
	inner *pterm.SpinnerPrinter
}

func (s ptermSpinner) Success(text string) {
	if s.inner == nil {
		return
	}
	s.inner.Success(text)
}

func (s ptermSpinner) Fail(text string) {
	if s.inner == nil {
		return
	}
	s.inner.Fail(text)
}

func (s ptermSpinner) Warning(text string) {
	if s.inner == nil {
		return
	}
	s.inner.Warning(text)
}

// startSpinner avoids spawning PTerm's spinner goroutine in RawOutput mode.
func startSpinner(text string) spinner {
	if pterm.RawOutput {
		return noopSpinner{}
	}
	sp, _ := pterm.DefaultSpinner.Start(text)
	if sp == nil {
		return noopSpinner{}
	}
	return ptermSpinner{inner: sp}
}
