package main

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type addingModel struct {
	projectName  textinput.Model
	rootPath     textinput.Model
	keyStorePath textinput.Model
	alias        textinput.Model
	cursor       int
}

func (m model) UpdateProject(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.currentState = projectsState

		case "up":
			m.addingModel.focusOnProperField()
			m.addingModel.cursor--

		case "down":
			m.addingModel.focusOnProperField()
			m.addingModel.cursor++

		case "enter":
			m.submit()
			return m, loadProjects()
		}

		switch {
		case m.addingModel.projectName.Focused():
			m.addingModel.projectName, cmd = m.addingModel.projectName.Update(msg)
		case m.addingModel.rootPath.Focused():
			m.addingModel.rootPath, cmd = m.addingModel.rootPath.Update(msg)
		case m.addingModel.keyStorePath.Focused():
			m.addingModel.keyStorePath, cmd = m.addingModel.keyStorePath.Update(msg)
		case m.addingModel.alias.Focused():
			m.addingModel.alias, cmd = m.addingModel.alias.Update(msg)
		}

	}
	return m, cmd
}

func (m model) ViewAddingProjectPrompt() tea.View {
	var c *tea.Cursor
	if !m.addingModel.projectName.VirtualCursor() && m.addingModel.projectName.Focused() {
		c = m.addingModel.projectName.Cursor()
		c.Y += lipgloss.Height(m.addingModel.headerView())
	}

	if !m.addingModel.rootPath.VirtualCursor() && m.addingModel.rootPath.Focused() {
		c = m.addingModel.rootPath.Cursor()
		c.Y += lipgloss.Height(m.addingModel.headerView()) + 4
	}

	if !m.addingModel.keyStorePath.VirtualCursor() && m.addingModel.keyStorePath.Focused() {
		c = m.addingModel.keyStorePath.Cursor()
		c.Y += lipgloss.Height(m.addingModel.headerView()) + 8
	}

	if !m.addingModel.alias.VirtualCursor() && m.addingModel.alias.Focused() {
		c = m.addingModel.alias.Cursor()
		c.Y += lipgloss.Height(m.addingModel.headerView()) + 12
	}

	str := lipgloss.JoinVertical(lipgloss.Top, m.addingModel.headerView(), "\nName:\n", m.addingModel.projectName.View(), "\nProject root path:\n", m.addingModel.rootPath.View(), "\nKeystore file path:\n", m.addingModel.keyStorePath.View(), "\nKey alias:\n", m.addingModel.alias.View(), m.addingModel.footerView())

	v := tea.NewView(str)
	v.Cursor = c
	return v

}

func (m addingModel) headerView() string { return "Add new project\n" }

func (m addingModel) footerView() string { return "\n(esc to cancel)" }

func (m *addingModel) focusOnProperField() {
	switch m.cursor {
	case 0:
		m.projectName.Focus()
		m.rootPath.Blur()
	case 1:
		m.projectName.Blur()
		m.rootPath.Focus()
		m.keyStorePath.Blur()
	case 2:
		m.rootPath.Blur()
		m.keyStorePath.Focus()
		m.alias.Blur()
	case 3:
		m.keyStorePath.Blur()
		m.alias.Focus()

	}
}

func (m addingModel) toProject() (project, error) {
	if m.projectName.Value() == "" || m.rootPath.Value() == "" || m.keyStorePath.Value() == "" || m.alias.Value() == "" {
		return project{}, errors.New("invalid inputs")
	}
	return project{
		name:         m.projectName.Value(),
		rootPath:     m.rootPath.Value(),
		keyStorePath: m.keyStorePath.Value(),
		alias:        m.alias.View(),
	}, nil
}

func loadProjects() tea.Cmd {
	return func() tea.Msg {
		projects, _ := getAllEntries()
		return projects
	}
}

func (m model) submit() {
	pro, err := m.addingModel.toProject()
	if err != nil {
		fmt.Print(err.Error())
	} else {
		insertProject(pro)
		m.currentState = projectsState
	}
}
