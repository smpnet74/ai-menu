package main

import (
	"fmt"
	"strings"
)

func (m model) renderCLITools() string {
	var b strings.Builder

	title := titleStyle.Render("🚀 Select CLI Tools to Install")
	b.WriteString(title)
	b.WriteString("\n\n")

	for i, tool := range m.cliTools {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := "[ ]"
		checkStyle := uncheckedStyle
		if m.selectedCLI[tool] {
			checked = "[✓]"
			checkStyle = checkedStyle
		}

		itemStyle := normalItemStyle
		if m.cursor == i {
			itemStyle = selectedItemStyle
		}

		line := fmt.Sprintf("%s %s %s", cursor, checkStyle.Render(checked), itemStyle.Render(tool))
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	help := helpStyle.Render("↑/k up • ↓/j down • space toggle • enter next • q quit")
	b.WriteString(help)
	b.WriteString("\n")

	return b.String()
}

func (m model) renderVSCodeExtensions() string {
	var b strings.Builder

	title := titleStyle.Render("🔌 Select VS Code Extensions to Install")
	b.WriteString(title)
	b.WriteString("\n\n")

	for i, ext := range m.vscodeExts {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := "[ ]"
		checkStyle := uncheckedStyle
		if m.selectedVSCode[ext] {
			checked = "[✓]"
			checkStyle = checkedStyle
		}

		itemStyle := normalItemStyle
		if m.cursor == i {
			itemStyle = selectedItemStyle
		}

		line := fmt.Sprintf("%s %s %s", cursor, checkStyle.Render(checked), itemStyle.Render(ext))
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	help := helpStyle.Render("↑/k up • ↓/j down • space toggle • enter next • q quit")
	b.WriteString(help)
	b.WriteString("\n")

	return b.String()
}

func (m model) renderSpecialTools() string {
	var b strings.Builder

	title := titleStyle.Render("🛠️  Select Special Tools to Install")
	b.WriteString(title)
	b.WriteString("\n\n")

	for i, tool := range m.specialTools {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := "[ ]"
		checkStyle := uncheckedStyle
		if m.selectedSpecial[tool] {
			checked = "[✓]"
			checkStyle = checkedStyle
		}

		itemStyle := normalItemStyle
		if m.cursor == i {
			itemStyle = selectedItemStyle
		}

		line := fmt.Sprintf("%s %s %s", cursor, checkStyle.Render(checked), itemStyle.Render(tool))
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	help := helpStyle.Render("↑/k up • ↓/j down • space toggle • enter review • q quit")
	b.WriteString(help)
	b.WriteString("\n")

	return b.String()
}

func (m model) renderPathInput() string {
	var b strings.Builder

	title := titleStyle.Render("📁 Enter Installation Directory")
	b.WriteString(title)
	b.WriteString("\n\n")

	b.WriteString(normalItemStyle.Render("Enter the parent directory path (ai-dev-pixi will be created inside):"))
	b.WriteString("\n\n")

	// Show the text input
	b.WriteString(m.pathInput.View())
	b.WriteString("\n\n")

	// Show the full path that will be created
	currentPath := m.pathInput.Value()
	if currentPath == "" {
		currentPath = m.installPath
	}
	fullPath := currentPath + "/ai-dev-pixi"
	pathPreview := helpStyle.Render(fmt.Sprintf("Installation path: %s", fullPath))
	b.WriteString(pathPreview)
	b.WriteString("\n\n")

	infoText := helpStyle.Render("A pixi environment with nodejs 22.* will be created at this location.\nSupports: linux-64, linux-aarch64")
	b.WriteString(infoText)
	b.WriteString("\n\n")

	help := helpStyle.Render("enter confirm • esc back")
	b.WriteString(help)
	b.WriteString("\n")

	return b.String()
}

func (m model) renderInstallSummary() string {
	var b strings.Builder

	title := titleStyle.Render("📦 Installation Summary")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Installation path
	if len(m.selectedCLI) > 0 {
		b.WriteString(summaryStyle.Render("Installation Path:"))
		b.WriteString("\n")
		fullPath := m.installPath + "/ai-dev-pixi"
		b.WriteString(fmt.Sprintf("  📁 %s\n", fullPath))
		b.WriteString("\n")
	}

	// CLI Tools
	if len(m.selectedCLI) > 0 {
		b.WriteString(summaryStyle.Render("CLI Tools:"))
		b.WriteString("\n")
		for tool := range m.selectedCLI {
			b.WriteString(fmt.Sprintf("  • %s\n", tool))
		}
		b.WriteString("\n")
	}

	// VS Code Extensions
	if len(m.selectedVSCode) > 0 {
		b.WriteString(summaryStyle.Render("VS Code Extensions:"))
		b.WriteString("\n")
		for ext := range m.selectedVSCode {
			b.WriteString(fmt.Sprintf("  • %s\n", ext))
		}
		b.WriteString("\n")
	}

	// Special Tools
	if len(m.selectedSpecial) > 0 {
		b.WriteString(summaryStyle.Render("Special Tools:"))
		b.WriteString("\n")
		for tool := range m.selectedSpecial {
			b.WriteString(fmt.Sprintf("  • %s\n", tool))
		}
		b.WriteString("\n")
	}

	if len(m.selectedCLI) == 0 && len(m.selectedVSCode) == 0 && len(m.selectedSpecial) == 0 {
		b.WriteString(helpStyle.Render("No items selected for installation."))
		b.WriteString("\n\n")
	}

	help := helpStyle.Render("enter to start installation • q quit without installing")
	b.WriteString(help)
	b.WriteString("\n")

	return b.String()
}

func (m model) renderInstalling() string {
	var b strings.Builder

	title := titleStyle.Render("⏳ Installing...")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Show spinner
	b.WriteString(m.spinner.View())
	b.WriteString(" Installing selected tools...\n\n")

	// Show recent messages (last 10)
	startIdx := 0
	if len(m.installMessages) > 10 {
		startIdx = len(m.installMessages) - 10
	}
	for i := startIdx; i < len(m.installMessages); i++ {
		b.WriteString(normalItemStyle.Render(m.installMessages[i]))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	help := helpStyle.Render("Please wait... Installation in progress")
	b.WriteString(help)
	b.WriteString("\n")

	return b.String()
}

func (m model) renderDone() string {
	var b strings.Builder

	title := titleStyle.Render("✅ Installation Complete!")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Count successes and failures
	successCount := 0
	failCount := 0
	for _, result := range m.installResults {
		if result.Success {
			successCount++
		} else {
			failCount++
		}
	}

	// Summary
	if successCount > 0 {
		b.WriteString(summaryStyle.Render(fmt.Sprintf("✓ %d tools installed successfully", successCount)))
		b.WriteString("\n")
	}
	if failCount > 0 {
		b.WriteString(helpStyle.Render(fmt.Sprintf("✗ %d tools failed to install", failCount)))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Detailed results
	cliToolInstalled := false
	for _, result := range m.installResults {
		if result.Success {
			b.WriteString(checkedStyle.Render(fmt.Sprintf("✓ %s", result.Name)))
			// Check if this is a CLI tool (npm package)
			if len(m.selectedCLI) > 0 {
				for cliTool := range m.selectedCLI {
					if result.Name == cliTool || result.Name == strings.Split(cliTool, " - ")[0] {
						cliToolInstalled = true
						break
					}
				}
			}
		} else {
			b.WriteString(uncheckedStyle.Render(fmt.Sprintf("✗ %s: %v", result.Name, result.Error)))
		}
		b.WriteString("\n")
	}

	// Show reminder to source .zshrc only if CLI tools were successfully installed
	if cliToolInstalled {
		b.WriteString("\n")
		b.WriteString(summaryStyle.Render("⚠️  Important: To use the CLI tool aliases, run:"))
		b.WriteString("\n")
		b.WriteString(selectedItemStyle.Render("  source ~/.zshrc"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	help := helpStyle.Render("Press enter or q to exit")
	b.WriteString(help)
	b.WriteString("\n")

	return b.String()
}
