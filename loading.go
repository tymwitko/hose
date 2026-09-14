package main

import (
	"fmt"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type buildingModel struct {
	spinner spinner.Model
}

func initializeSpinner() buildingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return buildingModel{spinner: s}
}

func (m model) UpdateBuilding(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		m.buildingModel.spinner, cmd = m.buildingModel.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) ViewBuilding() tea.View {
	str := fmt.Sprintf("\n\n\n   %s Building…\n\n\n", m.buildingModel.spinner.View())
	return tea.NewView(str)
}
