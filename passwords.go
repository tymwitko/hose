package main

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type passwordModel struct {
	project       project
	storePassword textinput.Model
	keyPassword   textinput.Model
	cursor        int
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

func (m passwordModel) headerView() string { return "Enter your store password\n" }

func (m passwordModel) footerView() string { return "\n(esc to quit)" }

func maskPassword(t textinput.Model) string {
	return fmt.Sprintf("> %s", strings.Repeat("*", len(t.Value())))
}
