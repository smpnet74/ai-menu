package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
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
	filepicker      filepicker.Model
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
	// Default parent directory is the current working directory
	// (ai-dev-pixi will be created inside this)
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	// Create filepicker
	fp := filepicker.New()
	fp.AllowedTypes = []string{} // Empty means directories only
	fp.CurrentDirectory = currentDir
	fp.ShowHidden = false
	fp.DirAllowed = true
	fp.FileAllowed = false

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = selectedItemStyle

	return model{
		state:           cliToolsView,
		cliTools:        getCLITools(),
		selectedCLI:     make(map[string]bool),
		vscodeExts:      getVSCodeExtensions(),
		selectedVSCode:  make(map[string]bool),
		specialTools:    getSpecialTools(),
		selectedSpecial: make(map[string]bool),
		cursor:          0,
		filepicker:      fp,
		installPath:     currentDir,
		spinner:         s,
		installMessages: []string{},
		installResults:  []InstallResult{},
	}
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
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
			case "q":
				// Go back to special tools view
				m.state = specialToolsView
				return m, nil
			}
		}

		m.filepicker, cmd = m.filepicker.Update(msg)

		// Check if a directory was selected
		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.installPath = path
			m.state = installView
			return m, nil
		}

		// Check if a directory was disabled (means it was selected as current dir)
		if didDisable, path := m.filepicker.DidSelectDisabledFile(msg); didDisable {
			m.installPath = path
			m.state = installView
			return m, nil
		}

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
