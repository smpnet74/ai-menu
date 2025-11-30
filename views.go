package main

import (
	"fmt"
	"sort"
	"strings"
)

func (m model) renderWelcome() string {
	var b strings.Builder

	// Add top padding
	b.WriteString("\n\n")

	title := titleStyle.Render("🎉 Welcome to Scott's AI World Installation Program")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Main welcome message
	welcome := normalItemStyle.Render("Thank you for choosing Scott's AI World!")
	b.WriteString(welcome)
	b.WriteString("\n\n")

	// Key information
	info1 := summaryStyle.Render("Important Information:")
	b.WriteString(info1)
	b.WriteString("\n")

	info2 := normalItemStyle.Render("Even if you do not select any optional tools to install,")
	b.WriteString(info2)
	b.WriteString("\n")

	info3 := normalItemStyle.Render("the following will ALWAYS be installed in the ai-dev-pixi environment:")
	b.WriteString(info3)
	b.WriteString("\n\n")

	// List of guaranteed installs
	node := checkedStyle.Render("✓") + " " + normalItemStyle.Render("Node.js 22.*")
	b.WriteString(node)
	b.WriteString("\n")

	python := checkedStyle.Render("✓") + " " + normalItemStyle.Render("Python 3.12.*")
	b.WriteString(python)
	b.WriteString("\n")

	uv := checkedStyle.Render("✓") + " " + normalItemStyle.Render("uv")
	b.WriteString(uv)
	b.WriteString("\n\n")

	// Additional info
	additional := helpStyle.Render("These core dependencies will be installed once and reused for")
	b.WriteString(additional)
	b.WriteString("\n")

	additional2 := helpStyle.Render("all subsequent selections, optimizing your installation time.")
	b.WriteString(additional2)
	b.WriteString("\n\n")

	// Optional selections notice
	optional := normalItemStyle.Render("You can then select additional CLI tools, VS Code extensions,")
	b.WriteString(optional)
	b.WriteString("\n")

	optional2 := normalItemStyle.Render("and special tools to customize your development environment.")
	b.WriteString(optional2)
	b.WriteString("\n\n")

	help := helpStyle.Render("Press enter to continue • q to quit")
	b.WriteString(help)
	b.WriteString("\n")

	return b.String()
}

func (m model) renderCLITools() string {
	var b strings.Builder

	// Add top padding
	b.WriteString("\n")

	title := titleStyle.Render("🚀 Select CLI Tools to Install")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Add "Select All" option at the top
	cursor := " "
	if m.cursor == 0 {
		cursor = ">"
	}

	allSelected := len(m.selectedCLI) == len(m.cliTools)
	checked := "[ ]"
	checkStyle := uncheckedStyle
	if allSelected && len(m.cliTools) > 0 {
		checked = "[✓]"
		checkStyle = checkedStyle
	}

	itemStyle := normalItemStyle
	if m.cursor == 0 {
		itemStyle = selectedItemStyle
	}

	line := fmt.Sprintf("%s %s %s", cursor, checkStyle.Render(checked), itemStyle.Render("Select All"))
	b.WriteString(line)
	b.WriteString("\n\n")

	// Render actual tools
	for i, tool := range m.cliTools {
		cursor := " "
		if m.cursor == i+1 {
			cursor = ">"
		}

		checked := "[ ]"
		checkStyle := uncheckedStyle
		if m.selectedCLI[tool] {
			checked = "[✓]"
			checkStyle = checkedStyle
		}

		itemStyle := normalItemStyle
		if m.cursor == i+1 {
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

	// Add top padding
	b.WriteString("\n")

	title := titleStyle.Render("🔌 Select VS Code Extensions to Install")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Add "Select All" option at the top
	cursor := " "
	if m.cursor == 0 {
		cursor = ">"
	}

	allSelected := len(m.selectedVSCode) == len(m.vscodeExts)
	checked := "[ ]"
	checkStyle := uncheckedStyle
	if allSelected && len(m.vscodeExts) > 0 {
		checked = "[✓]"
		checkStyle = checkedStyle
	}

	itemStyle := normalItemStyle
	if m.cursor == 0 {
		itemStyle = selectedItemStyle
	}

	line := fmt.Sprintf("%s %s %s", cursor, checkStyle.Render(checked), itemStyle.Render("Select All"))
	b.WriteString(line)
	b.WriteString("\n\n")

	// Render actual extensions
	for i, ext := range m.vscodeExts {
		cursor := " "
		if m.cursor == i+1 {
			cursor = ">"
		}

		checked := "[ ]"
		checkStyle := uncheckedStyle
		if m.selectedVSCode[ext] {
			checked = "[✓]"
			checkStyle = checkedStyle
		}

		itemStyle := normalItemStyle
		if m.cursor == i+1 {
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

	// Add top padding
	b.WriteString("\n")

	title := titleStyle.Render("🔧 Select Special Tools to Install")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Add "Select All" option at the top
	cursor := " "
	if m.cursor == 0 {
		cursor = ">"
	}

	allSelected := len(m.selectedSpecial) == len(m.specialTools)
	checked := "[ ]"
	checkStyle := uncheckedStyle
	if allSelected && len(m.specialTools) > 0 {
		checked = "[✓]"
		checkStyle = checkedStyle
	}

	itemStyle := normalItemStyle
	if m.cursor == 0 {
		itemStyle = selectedItemStyle
	}

	line := fmt.Sprintf("%s %s %s", cursor, checkStyle.Render(checked), itemStyle.Render("Select All"))
	b.WriteString(line)
	b.WriteString("\n\n")

	// Render actual tools
	for i, tool := range m.specialTools {
		cursor := " "
		if m.cursor == i+1 {
			cursor = ">"
		}

		checked := "[ ]"
		checkStyle := uncheckedStyle
		if m.selectedSpecial[tool] {
			checked = "[✓]"
			checkStyle = checkedStyle
		}

		itemStyle := normalItemStyle
		if m.cursor == i+1 {
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

	// Add top padding
	b.WriteString("\n")

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

	// Add top padding
	b.WriteString("\n")

	title := titleStyle.Render("📦 Selected Software for Installation")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Installation path
	if len(m.selectedCLI) > 0 || len(m.selectedVSCode) > 0 || len(m.selectedSpecial) > 0 {
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
		
		// Sort tools to ensure consistent ordering
		var tools []string
		for tool := range m.selectedCLI {
			tools = append(tools, tool)
		}
		sort.Strings(tools)
		
		for _, tool := range tools {
			b.WriteString(fmt.Sprintf("  • %s\n", tool))
		}
		b.WriteString("\n")
	}

	// VS Code Extensions
	if len(m.selectedVSCode) > 0 {
		b.WriteString(summaryStyle.Render("VS Code Extensions:"))
		b.WriteString("\n")
		
		// Sort extensions to ensure consistent ordering
		var extensions []string
		for ext := range m.selectedVSCode {
			extensions = append(extensions, ext)
		}
		sort.Strings(extensions)
		
		for _, ext := range extensions {
			b.WriteString(fmt.Sprintf("  • %s\n", ext))
		}
		b.WriteString("\n")
	}

	// Special Tools
	if len(m.selectedSpecial) > 0 {
		b.WriteString(summaryStyle.Render("Special Tools:"))
		b.WriteString("\n")
		
		// Sort tools to ensure consistent ordering
		var tools []string
		for tool := range m.selectedSpecial {
			tools = append(tools, tool)
		}
		sort.Strings(tools)
		
		for _, tool := range tools {
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

	// Add top padding
	b.WriteString("\n")

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

	// Add top padding
	b.WriteString("\n")

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
