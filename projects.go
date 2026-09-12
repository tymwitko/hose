package main

import (
	"fmt"
	"os/exec"

	tea "charm.land/bubbletea/v2"
)

type projectModel struct {
	projects []project
	cursor   int
}

type project struct {
	name         string
	rootPath     string
	keyStorePath string
	alias        string
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
