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

func (m model) UpdatePassword(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "up":
			if m.passwordModel.cursor > 0 {
				if m.passwordModel.keyPassword.Focused() {
					m.passwordModel.keyPassword.Blur()
				}
				m.passwordModel.storePassword.Focus()
				m.passwordModel.cursor--
			}

		case "down":
			if m.passwordModel.cursor < 2 {
				if m.passwordModel.storePassword.Focused() {
					m.passwordModel.storePassword.Blur()
				}
				m.passwordModel.keyPassword.Focus()
				m.passwordModel.cursor++
			}

		case "enter":
			fmt.Printf("Starting build…")
			m.currentState = buildingState
			return m, tea.Batch(
				m.buildingModel.spinner.Tick,
				m.passwordModel.project.Build(m.passwordModel.storePassword.Value(), m.passwordModel.keyPassword.Value()),
			)
		}

		switch {
		case m.passwordModel.keyPassword.Focused():
			m.passwordModel.keyPassword, cmd = m.passwordModel.keyPassword.Update(msg)
		case m.passwordModel.storePassword.Focused():
			m.passwordModel.storePassword, cmd = m.passwordModel.storePassword.Update(msg)
		}

	}
	return m, cmd
}

func (m passwordModel) headerView() string { return "Enter your store password\n" }

func (m passwordModel) footerView() string { return "\n(esc to quit)" }

func maskPassword(t textinput.Model) string {
	return fmt.Sprintf("> %s", strings.Repeat("*", len(t.Value())))
}
