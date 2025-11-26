package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case cliToolsView:
		m.state = vscodeExtensionsView
		m.cursor = 0
	case vscodeExtensionsView:
		m.state = specialToolsView
		m.cursor = 0
	case specialToolsView:
		// Go to path input if CLI tools are selected
		if len(m.selectedCLI) > 0 {
			m.state = pathInputView
		} else {
			m.state = installView
		}
		m.cursor = 0
	case pathInputView:
		m.state = installView
		m.cursor = 0
	case installView:
		// Trigger installation
		return m, func() tea.Msg { return installMsgStart{} }
	}
	return m, nil
}

func (m model) handleDown() int {
	var maxLen int
	switch m.state {
	case cliToolsView:
		maxLen = len(m.cliTools)
	case vscodeExtensionsView:
		maxLen = len(m.vscodeExts)
	case specialToolsView:
		maxLen = len(m.specialTools)
	default:
		return m.cursor
	}

	if m.cursor < maxLen-1 {
		return m.cursor + 1
	}
	return m.cursor
}

func (m *model) toggleSelection() {
	switch m.state {
	case cliToolsView:
		if m.cursor < len(m.cliTools) {
			tool := m.cliTools[m.cursor]
			m.selectedCLI[tool] = !m.selectedCLI[tool]
			if !m.selectedCLI[tool] {
				delete(m.selectedCLI, tool)
			}
		}
	case vscodeExtensionsView:
		if m.cursor < len(m.vscodeExts) {
			ext := m.vscodeExts[m.cursor]
			m.selectedVSCode[ext] = !m.selectedVSCode[ext]
			if !m.selectedVSCode[ext] {
				delete(m.selectedVSCode, ext)
			}
		}
	case specialToolsView:
		if m.cursor < len(m.specialTools) {
			tool := m.specialTools[m.cursor]
			m.selectedSpecial[tool] = !m.selectedSpecial[tool]
			if !m.selectedSpecial[tool] {
				delete(m.selectedSpecial, tool)
			}
		}
	}
}

func (m model) performInstallation() tea.Cmd {
	return func() tea.Msg {
		// Create a progress callback that sends messages via the program
		progress := func(msg string) {
			if program != nil {
				program.Send(installMsg{message: msg})
			}
		}

		// Convert selected maps to slices
		cliTools := make([]string, 0, len(m.selectedCLI))
		for tool := range m.selectedCLI {
			cliTools = append(cliTools, tool)
		}

		vscodeExts := make([]string, 0, len(m.selectedVSCode))
		for ext := range m.selectedVSCode {
			vscodeExts = append(vscodeExts, ext)
		}

		specialTools := make([]string, 0, len(m.selectedSpecial))
		for tool := range m.selectedSpecial {
			specialTools = append(specialTools, tool)
		}

		// Collect all results
		allResults := []InstallResult{}

		// Perform installations
		if len(cliTools) > 0 {
			results := InstallCLITools(cliTools, m.installPath, progress)
			allResults = append(allResults, results...)
		}

		if len(vscodeExts) > 0 {
			results := InstallVSCodeExtensions(vscodeExts, progress)
			allResults = append(allResults, results...)
		}

		if len(specialTools) > 0 {
			results := InstallSpecialTools(specialTools, progress)
			allResults = append(allResults, results...)
		}

		return installCompleteMsg{results: allResults}
	}
}
