package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case welcomeView:
		m.state = cliToolsView
		m.cursor = 0
	case cliToolsView:
		m.state = vscodeExtensionsView
		m.cursor = 0
	case vscodeExtensionsView:
		m.state = specialToolsView
		m.cursor = 0
	case specialToolsView:
		m.state = cliEnhancersView
		m.cursor = 0
	case cliEnhancersView:
		// ALWAYS go to path input because core dependencies (Node 22 and Python 3.12)
		// are guaranteed to be installed regardless of selections
		m.state = pathInputView
		m.pathInput.Focus()
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
		// +1 for "Select All" option at the top
		maxLen = len(m.cliTools) + 1
	case vscodeExtensionsView:
		// +1 for "Select All" option at the top
		maxLen = len(m.vscodeExts) + 1
	case specialToolsView:
		// +1 for "Select All" option at the top
		maxLen = len(m.specialTools) + 1
	case cliEnhancersView:
		// +1 for "Select All" option at the top
		maxLen = len(m.cliEnhancers) + 1
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
		if m.cursor == 0 {
			// Toggle Select All
			if len(m.selectedCLI) == len(m.cliTools) {
				// All selected, deselect all
				m.selectedCLI = make(map[string]bool)
			} else {
				// Not all selected, select all
				for _, tool := range m.cliTools {
					m.selectedCLI[tool] = true
				}
			}
		} else if m.cursor <= len(m.cliTools) {
			tool := m.cliTools[m.cursor-1]
			m.selectedCLI[tool] = !m.selectedCLI[tool]
			if !m.selectedCLI[tool] {
				delete(m.selectedCLI, tool)
			}
		}
	case vscodeExtensionsView:
		if m.cursor == 0 {
			// Toggle Select All
			if len(m.selectedVSCode) == len(m.vscodeExts) {
				// All selected, deselect all
				m.selectedVSCode = make(map[string]bool)
			} else {
				// Not all selected, select all
				for _, ext := range m.vscodeExts {
					m.selectedVSCode[ext] = true
				}
			}
		} else if m.cursor <= len(m.vscodeExts) {
			ext := m.vscodeExts[m.cursor-1]
			m.selectedVSCode[ext] = !m.selectedVSCode[ext]
			if !m.selectedVSCode[ext] {
				delete(m.selectedVSCode, ext)
			}
		}
	case specialToolsView:
		if m.cursor == 0 {
			// Toggle Select All
			if len(m.selectedSpecial) == len(m.specialTools) {
				// All selected, deselect all
				m.selectedSpecial = make(map[string]bool)
			} else {
				// Not all selected, select all
				for _, tool := range m.specialTools {
					m.selectedSpecial[tool] = true
				}
			}
		} else if m.cursor <= len(m.specialTools) {
			tool := m.specialTools[m.cursor-1]
			m.selectedSpecial[tool] = !m.selectedSpecial[tool]
			if !m.selectedSpecial[tool] {
				delete(m.selectedSpecial, tool)
			}
		}
	case cliEnhancersView:
		if m.cursor == 0 {
			// Toggle Select All
			if len(m.selectedCLIEnhancers) == len(m.cliEnhancers) {
				// All selected, deselect all
				m.selectedCLIEnhancers = make(map[string]bool)
			} else {
				// Not all selected, select all
				for _, enhancer := range m.cliEnhancers {
					m.selectedCLIEnhancers[enhancer] = true
				}
			}
		} else if m.cursor <= len(m.cliEnhancers) {
			enhancer := m.cliEnhancers[m.cursor-1]
			m.selectedCLIEnhancers[enhancer] = !m.selectedCLIEnhancers[enhancer]
			if !m.selectedCLIEnhancers[enhancer] {
				delete(m.selectedCLIEnhancers, enhancer)
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

		// Collect all results
		allResults := []InstallResult{}

		// ALWAYS ensure core dependencies first (Node 22.* and Python 3.12.*)
		progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		if !EnsureCoreDependencies(m.installPath, progress) {
			progress("✗ Failed to ensure core dependencies")
			return installCompleteMsg{results: allResults}
		}
		progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		progress("")

		// Convert selected maps to slices
		cliTools := make([]string, 0, len(m.selectedCLI))
		for tool := range m.selectedCLI {
			// Convert display name to package name for installation
			packageName := getPackageNameForCLI(tool)
			cliTools = append(cliTools, packageName)
		}

		vscodeExts := make([]string, 0, len(m.selectedVSCode))
		for ext := range m.selectedVSCode {
			vscodeExts = append(vscodeExts, ext)
		}

		specialTools := make([]string, 0, len(m.selectedSpecial))
		for tool := range m.selectedSpecial {
			specialTools = append(specialTools, tool)
		}

		cliEnhancers := make([]string, 0, len(m.selectedCLIEnhancers))
		for enhancer := range m.selectedCLIEnhancers {
			cliEnhancers = append(cliEnhancers, enhancer)
		}

		// Perform installations
		if len(cliTools) > 0 {
			progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			results := InstallCLITools(cliTools, m.installPath, progress)
			allResults = append(allResults, results...)
			progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			progress("")
		}

		if len(vscodeExts) > 0 {
			progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			results := InstallVSCodeExtensions(vscodeExts, progress)
			allResults = append(allResults, results...)
			progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			progress("")
		}

		if len(specialTools) > 0 {
			progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			results := InstallSpecialTools(specialTools, progress)
			allResults = append(allResults, results...)
			progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			progress("")
		}

		if len(cliEnhancers) > 0 {
			progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			results := InstallCLIEnhancers(cliEnhancers, m.installPath, progress)
			allResults = append(allResults, results...)
			progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			progress("")
		}

		return installCompleteMsg{results: allResults}
	}
}
