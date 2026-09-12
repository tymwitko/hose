package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	project       project
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

func (m model) UpdateProjects(msg tea.Msg) (model, tea.Cmd, bool) {
	var shouldChangeState = false
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit, false

		case "up", "k":
			if m.projectModel.cursor > 0 {
				m.projectModel.cursor--
			}

		case "down", "j":
			if m.projectModel.cursor < len(m.projectModel.projects)-1 {
				m.projectModel.cursor++
			}

		case "enter", "space":
			pro := m.projectModel.projects[m.projectModel.cursor]
			s := pro.CheckGitCleanness()
			shouldChangeState = true
			if s != "" {
				fmt.Printf("error %s", s)
			} else {
				m.passwordModel.project = pro
				m.passwordModel.storePassword.Focus()
				m.passwordModel.keyPassword.SetVirtualCursor(false)
				m.passwordModel.storePassword.SetVirtualCursor(false)
				fmt.Print("Git check successful!")
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
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				if m.keyPassword.Focused() {
					m.keyPassword.Blur()
				}
				m.storePassword.Focus()
				m.cursor--
			}

		case "down":
			if m.cursor < 2 {
				if m.storePassword.Focused() {
					m.storePassword.Blur()
				}
				m.keyPassword.Focus()
				m.cursor++
			}

		case "enter":
			fmt.Printf("Starting build…")
			res := m.project.Build(m.storePassword.Value(), m.keyPassword.Value())
			fmt.Print(res)
		}
	}

	switch {
	case m.keyPassword.Focused():
		m.keyPassword, cmd = m.keyPassword.Update(msg)
	case m.storePassword.Focused():
		m.storePassword, cmd = m.storePassword.Update(msg)
	}

	return m, cmd
}

func (p project) CheckGitCleanness() string {
	result, err := exec.Command("git", "-C", p.rootPath, "status", "--porcelain").Output()
	if err != nil || len(result) != 0 {
		return fmt.Sprintf("Git not clean, commit or stash your changes! Error is %s, %s", err, string(result))
	}
	return ""
}

func (p project) Build(storePass string, keyPass string) string {
	res, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd %s && ./gradlew clean assembleRelease -Pandroid.injected.signing.store.file=%s -Pandroid.injected.signing.store.password=%s -Pandroid.injected.signing.key.alias=%s -Pandroid.injected.signing.key.password=%s", p.rootPath, p.keyStorePath, storePass, p.alias, keyPass)).Output()
	if err != nil {
		return fmt.Sprintf("Build failed with message %s %s, passwords were %s, %s", res, err, storePass, keyPass)
	}
	return ""
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
			s.WriteString(fmt.Sprintf("%s %s\n", cursor, p))
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

func (m passwordModel) headerView() string { return "Enter your store password\n" }

func (m passwordModel) footerView() string { return "\n(esc to quit)" }

func maskPassword(t textinput.Model) string {
	return fmt.Sprintf("> %s", strings.Repeat("*", len(t.Value())))
}
