package main

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

func (m model) ViewPasswordPrompts() tea.View {
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

func (m passwordModel) headerView() string { return "Enter your store password\n" }

func (m passwordModel) footerView() string { return "\n(esc to quit)" }

func maskPassword(t textinput.Model) string {
	return fmt.Sprintf("> %s", strings.Repeat("*", len(t.Value())))
}
