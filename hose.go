package main

import (
	"bytes"
	"fmt"
	"os"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func main() {

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Print("Oopsie")
		os.Exit(1)
	}
}

type state int

const (
	projectsState state = iota
	passwordsState
	buildingState
)

type model struct {
	currentState  state
	projectModel  projectModel
	passwordModel passwordModel
	buildingModel buildingModel
}

func initialModel() model {
	return model{
		projectsState,
		projectModel{
			projects: []project{
				project{
					"Recents",
					"/home/tymon/Kodzenie/Kotlin/Recents",
					"/home/tymon/Kodzenie/Kotlin/Recents/recents_keystore.jks",
					"recents",
				},
				project{
					"ToReplace",
					"/home/tymon/Kodzenie/Kotlin/ToReplace",
					"/home/tymon/Kodzenie/Kotlin/Recents/recents_keystore.jks",
					"recents",
				},
			},
		},
		passwordModel{
			storePassword: textinput.New(),
			keyPassword:   textinput.New(),
		},
		initializeSpinner(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currentState {
	case projectsState:
		m, cmd = m.UpdateProjects(msg)
		return m, cmd
	case passwordsState:
		m, cmd = m.UpdatePassword(msg)
		return m, cmd
	case buildingState:
		m, cmd = m.UpdateBuilding(msg)
	}
	return m, nil
}

func (m model) View() tea.View {
	switch m.currentState {
	case projectsState:
		var s bytes.Buffer
		s.WriteString("Select project to build\n\n")

		for i, p := range m.projectModel.projects {
			cursor := " "
			if m.projectModel.cursor == i {
				cursor = ">"
			}
			fmt.Fprintf(&s, "%s %s\n", cursor, p)
		}
		s.WriteString("\nPress q to quit.\n")
		return tea.NewView(s.String())

	case passwordsState:
		var c *tea.Cursor
		if !m.passwordModel.storePassword.VirtualCursor() && m.passwordModel.storePassword.Focused() {
			c = m.passwordModel.storePassword.Cursor()
			c.Y += lipgloss.Height(m.passwordModel.headerView())
		}

		if !m.passwordModel.keyPassword.VirtualCursor() && m.passwordModel.keyPassword.Focused() {
			c = m.passwordModel.keyPassword.Cursor()
			c.Y += lipgloss.Height(m.passwordModel.headerView()) + 4
		}

		str := lipgloss.JoinVertical(lipgloss.Top, m.passwordModel.headerView(), maskPassword(m.passwordModel.storePassword), "\nEnter your key password\n", maskPassword(m.passwordModel.keyPassword), m.passwordModel.footerView())

		v := tea.NewView(str)
		v.Cursor = c
		return v

	case buildingState:
		str := fmt.Sprintf("\n\n\n   %s Building…\n\n\n", m.buildingModel.spinner.View())
		return tea.NewView(str)
	}
	return tea.NewView("Error")
}
