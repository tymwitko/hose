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
)

type model struct {
	currentState  state
	projectModel  projectModel
	passwordModel passwordModel
}

func initialModel() model {
	return model{
		projectsState,
		projectModel{
			projects: []project{},
		},
		passwordModel{
			storePassword: textinput.New(),
			keyPassword:   textinput.New(),
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentState {
	case projectsState:
		pr, cmd, shouldSwitch := m.UpdateProjects(msg)
		m = pr
		if shouldSwitch {
			m.currentState = passwordsState
		}
		return m, cmd
	case passwordsState:
		pass, cmd := m.passwordModel.Update(msg)
		m.passwordModel = pass
		return m, cmd
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
	}
	return tea.NewView("Error")
}
