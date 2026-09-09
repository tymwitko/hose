package main

import (
	"fmt"
	"os"
	"os/exec"

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

type projectModel struct {
	projects []project
	cursor   int
}

type passwordModel struct {
	storePassword textinput.Model
	keyPassword   textinput.Model
	cursor        int
}

type project struct {
	name         string
	rootPath     string
	keyStorePath string
	alias        string
}

func initialModel() model {
	return model{
		projectsState,
		projectModel{
			projects: []project{},
		},
		passwordModel{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentState {
	case projectsState:
		pr, cmd, shouldSwitch := m.projectModel.Update(msg)
		m.projectModel = pr
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

func (m projectModel) Update(msg tea.Msg) (projectModel, tea.Cmd, bool) {
	var shouldChangeState = false
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit, false

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			}

		case "enter", "space":
			pro := m.projects[m.cursor]
			s := pro.CheckGitCleanness()
			shouldChangeState = true
			storePass, keyPass := getPasswords()
			s += pro.Build(storePass, keyPass)
			if s != "" {
				fmt.Printf("error %s", s)
			} else {
				fmt.Print("Build successful!")
			}
		}
	}
	return m, nil, shouldChangeState
}

func (m passwordModel) Update(msg tea.Msg) (passwordModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < 2 {
				m.cursor++
			}

		}
	}

	m.storePassword, cmd = m.storePassword.Update(msg)
	return m, cmd
}

func (p project) CheckGitCleanness() string {
	result, err := exec.Command("git", "-C", p.rootPath, "status", "--porcelain").Output()
	if err != nil || len(result) != 0 {
		return fmt.Sprintf("Git not clean, commit or stash your changes! Error is %s, %s", err, string(result))
	}
	return ""
}

func getPasswords() (string, string) {
	return "password1", "password2" // todo: get from text inputs
}

func (p project) Build(storePass string, keyPass string) string {
	_, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd %s && ./gradlew clean assembleRelease -Pandroid.injected.signing.store.file=%s -Pandroid.injected.signing.store.password=%s -Pandroid.injected.signing.key.alias=%s -Pandroid.injected.signing.key.password=%s", p.rootPath, p.keyStorePath, storePass, p.alias, keyPass)).Output()
	if err != nil {
		return fmt.Sprintf("Build failed with message %s", err)
	}
	return ""
}

func (m model) View() tea.View {
	switch m.currentState {
	case projectsState:
		s := "Select project to build\n\n"

		for i, p := range m.projectModel.projects {
			cursor := " "
			if m.projectModel.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, p)
		}
		s += "\nPress q to quit.\n"
		return tea.NewView(s)

	case passwordsState:
		str := "Enter store password\n"
		var c *tea.Cursor
		if !m.passwordModel.storePassword.VirtualCursor() {
			c = m.passwordModel.storePassword.Cursor()
			//c.Y += lipgloss.Height(m.headerView())
		}

		// str := lipgloss.JoinVertical(lipgloss.Top, m.passwordModel.headerView(), m.textInput.View(), m.passwordModel.footerView())
		str += lipgloss.JoinVertical(lipgloss.Top, m.passwordModel.storePassword.View())
		//if m.quitting {
		//str += "\n"
		//}

		v := tea.NewView(str)
		v.Cursor = c
		return v
	}
	return tea.NewView("Error")
}
