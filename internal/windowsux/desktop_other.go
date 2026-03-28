//go:build !windows

package windowsux

import (
	"github.com/pasqualotodiogenes/sleepoff/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

type Integration struct{}

func NewIntegration(_ func()) *Integration { return &Integration{} }

func (i *Integration) AttachProgram(_ *tea.Program) {}

func (i *Integration) UpdateDesktop(_ model.DesktopState) {}

func (i *Integration) Close() {}
