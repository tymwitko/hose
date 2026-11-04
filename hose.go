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
	addingState
	resultState
)

type model struct {
	currentState  state
	projectModel  projectModel
	passwordModel passwordModel
	buildingModel buildingModel
	addingModel   addingModel
	resultModel   resultModel
}

type buildFinishedMsg struct{}

func initialModel() model {
	projects, _ := getAllEntries()
	return model{
		projectsState,
		projectModel{
			projects: projects,
		},
		passwordModel{
			storePassword: textinput.New(),
			keyPassword:   textinput.New(),
		},
		initializeSpinner(),
		addingModel{
			projectName:  textinput.New(),
			rootPath:     textinput.New(),
			keyStorePath: textinput.New(),
			alias:        textinput.New(),
		},
		resultModel{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case buildFinishedMsg:
		m.currentState = resultState
		return m, cmd
	}

	switch m.currentState {
	case projectsState:
		m, cmd = m.UpdateProjects(msg)
		return m, cmd
	case passwordsState:
		m, cmd = m.UpdatePassword(msg)
		return m, cmd
	case buildingState:
		m, cmd = m.UpdateBuilding(msg)
		return m, cmd
	case addingState:
		m, cmd = m.UpdateProject(msg)
		return m, cmd
	case resultState:
		m, cmd = m.UpdateResult(msg)
		return m, cmd
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
	case addingState:
		return m.ViewAddingProjectPrompt()
	case resultState:
		return m.ViewResult()
	}
	return tea.NewView("Error")
}
