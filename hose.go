package main

import (
	"fmt"
	"os"
	"os/exec"

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

type model struct {
	currentState state
	projects     []project
	cursor       int
}

type project struct {
	name         string
	rootPath     string
	keyStorePath string
	alias        string
}

func initialModel() model {
	return model{
		projects: []project{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

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
			storePass, keyPass := getPasswords()
			s += pro.Build(storePass, keyPass)
			if s != "" {
				fmt.Printf("error %s", s)
			} else {
				fmt.Print("Build successful!")
			}
		}
	}
	return m, nil
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
	s := "Select project to build\n\n"

	for i, p := range m.projects {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, p)
	}
	s += "\nPress q to quit.\n"

	return tea.NewView(s)
}
