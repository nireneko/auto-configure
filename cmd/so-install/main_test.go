package main

import (
	"bytes"
	"testing"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/presentation/tui"
	"github.com/so-install/pkg/mocks"
)

func TestRun_PrivilegeError(t *testing.T) {
	oldGetuid := osGetuid
	osGetuid = func() int { return 1000 }
	defer func() { osGetuid = oldGetuid }()

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	
	exitCode := Run([]string{"1x-so-install"}, out, errOut)
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}

func TestRun_OSDetectError(t *testing.T) {
	oldGetuid := osGetuid
	osGetuid = func() int { return 0 }
	defer func() { osGetuid = oldGetuid }()

	oldNewDetector := newDetector
	newDetector = func() domain.OSDetector {
		return &mocks.MockOSDetector{ReturnErr: errors.New("detect fail")}
	}
	defer func() { newDetector = oldNewDetector }()

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	exitCode := Run([]string{"1x-so-install"}, out, errOut)
	if exitCode != 1 {
		t.Errorf("expected exit code 1 on OS error, got %d", exitCode)
	}
}

func TestRun_TUIError(t *testing.T) {
	oldGetuid := osGetuid
	osGetuid = func() int { return 0 }
	defer func() { osGetuid = oldGetuid }()

	oldNewDetector := newDetector
	newDetector = func() domain.OSDetector {
		return &mocks.MockOSDetector{ReturnID: "debian", ReturnVersionID: "12"}
	}
	defer func() { newDetector = oldNewDetector }()

	oldRunProgram := runProgram
	runProgram = func(p *tea.Program) (tea.Model, error) {
		return nil, errors.New("tui fail")
	}
	defer func() { runProgram = oldRunProgram }()

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	exitCode := Run([]string{"1x-so-install"}, out, errOut)
	if exitCode != 1 {
		t.Errorf("expected exit code 1 on TUI error, got %d", exitCode)
	}
}

func TestRun_Success(t *testing.T) {
	oldGetuid := osGetuid
	osGetuid = func() int { return 0 }
	defer func() { osGetuid = oldGetuid }()

	oldGetenv := osGetenv
	osGetenv = func(key string) string {
		if key == "SUDO_USER" { return "testuser" }
		return ""
	}
	defer func() { osGetenv = oldGetenv }()

	oldNewDetector := newDetector
	newDetector = func() domain.OSDetector {
	        return &mocks.MockOSDetector{
	                ReturnID: "debian",
	                ReturnVersionID: "12",
	        }
	}
	defer func() { newDetector = oldNewDetector }()

	oldRunProgramLocal := runProgram
	runProgram = func(p *tea.Program) (tea.Model, error) {
	        return tui.NewModel(nil, nil), nil
	}
	defer func() { runProgram = oldRunProgramLocal }()

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	
	exitCode := Run([]string{"1x-so-install"}, out, errOut)
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d. Error: %s", exitCode, errOut.String())
	}
}
