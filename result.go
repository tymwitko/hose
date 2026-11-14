package main

import (
	"fmt"
	"os/exec"
	"runtime"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type resultModel struct {
	title  string
	isSucc bool
}

func (m model) UpdateResult(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			return m, tea.Quit
		}
	default:
		return m, tea.Quit
	}
}

func (m model) ViewResult() tea.View {
	str := lipgloss.JoinVertical(
		lipgloss.Center,
		fmt.Sprintf("%s", m.resultModel.title),
	)
	content := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Top,
		lipgloss.Center,
		str,
	)
	return tea.NewView(content)
}

func openFolder(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default: // bsd
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
