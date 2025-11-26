package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	cliToolsView sessionState = iota
	vscodeExtensionsView
	specialToolsView
	pathInputView
	installView
	installingView
	doneView
	quitView
)

type model struct {
	state           sessionState
	cliTools        []string
	selectedCLI     map[string]bool
	vscodeExts      []string
	selectedVSCode  map[string]bool
	specialTools    []string
	selectedSpecial map[string]bool
	cursor          int
	pathInput       textinput.Model
	installPath     string
	spinner         spinner.Model
	installing      bool
	installMessages []string
	installResults  []InstallResult
	err             error
}

// Installation messages
type installMsgStart struct{}
type installMsg struct{ message string }
type installCompleteMsg struct{ results []InstallResult }

func initialModel() model {
	// Default installation directory
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	// Create text input for path
	ti := textinput.New()
	ti.Placeholder = "/path/to/parent/directory"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50
	ti.SetValue(currentDir)

	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = spinnerStyle

	return model{
		state:           cliToolsView,
		cliTools:        getCLITools(),
		selectedCLI:     make(map[string]bool),
		vscodeExts:      getVSCodeExtensions(),
		selectedVSCode:  make(map[string]bool),
		specialTools:    getSpecialTools(),
		selectedSpecial: make(map[string]bool),
		cursor:          0,
		pathInput:       ti,
		installPath:     currentDir,
		spinner:         s,
		installMessages: []string{},
		installResults:  []InstallResult{},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle installation messages
	switch msg := msg.(type) {
	case installMsgStart:
		m.installing = true
		m.state = installingView
		return m, tea.Batch(m.spinner.Tick, m.performInstallation())

	case installMsg:
		m.installMessages = append(m.installMessages, msg.message)
		return m, nil

	case installCompleteMsg:
		m.installing = false
		m.installResults = msg.results
		m.state = doneView
		return m, nil

	case spinner.TickMsg:
		if m.installing {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	// Handle path input separately
	if m.state == pathInputView {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				m.state = quitView
				return m, tea.Quit
			case "esc":
				// Go back to special tools view
				m.state = specialToolsView
				return m, nil
			case "enter":
				// Validate and accept the path
				path := m.pathInput.Value()
				if path != "" {
					m.installPath = path
					m.state = installView
				}
				return m, nil
			}
		}

		m.pathInput, cmd = m.pathInput.Update(msg)
		return m, cmd
	}

	// Handle done view
	if m.state == doneView {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "q", " ":
				return m, tea.Quit
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.state = quitView
			return m, tea.Quit

		case "enter":
			return m.handleEnter()

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			m.cursor = m.handleDown()

		case " ":
			m.toggleSelection()
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case cliToolsView:
		return m.renderCLITools()
	case vscodeExtensionsView:
		return m.renderVSCodeExtensions()
	case specialToolsView:
		return m.renderSpecialTools()
	case pathInputView:
		return m.renderPathInput()
	case installView:
		return m.renderInstallSummary()
	case installingView:
		return m.renderInstalling()
	case doneView:
		return m.renderDone()
	case quitView:
		return "Goodbye!\n"
	}
	return ""
}

var program *tea.Program

func main() {
	program = tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
