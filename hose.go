package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
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
			projects: []project{},
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
		return m.ViewProjects()
	case passwordsState:
		return m.ViewPasswordPrompts()
	case buildingState:
		return m.ViewBuilding()
	}
	return tea.NewView("Error")
}
