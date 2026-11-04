package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type resultModel struct {
	title string
}

func (m model) UpdateResult(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		m.buildingModel.spinner, cmd = m.buildingModel.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) ViewResult() tea.View {
	str := fmt.Sprintf("\n\n\n   %s\n\n\n", m.resultModel.title)
	return tea.NewView(str)
}
