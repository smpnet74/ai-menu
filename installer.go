package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type InstallResult struct {
	Name    string
	Success bool
	Error   error
	Message string
}

type ProgressCallback func(message string)

// EnsureCoreDependencies ensures Node 22.*, Python 3.12.*, and uv are in the pixi environment
// It only adds them if they don't already exist, preventing reinstalls
func EnsureCoreDependencies(installPath string, progress ProgressCallback) bool {
	progress("Ensuring core dependencies (Node 22.*, Python 3.12.*, and uv) are available...")

	// Append ai-dev-pixi to the provided parent path
	envDir := installPath + "/ai-dev-pixi"

	// Create directory if it doesn't exist
	if err := os.MkdirAll(envDir, 0755); err != nil {
		msg := fmt.Sprintf("✗ Failed to create directory %s: %v", envDir, err)
		progress(msg)
		return false
	}

	// Change to the environment directory
	if err := os.Chdir(envDir); err != nil {
		msg := fmt.Sprintf("✗ Failed to change to directory %s: %v", envDir, err)
		progress(msg)
		return false
	}

	progress(fmt.Sprintf("Environment directory: %s", envDir))

	// Initialize pixi project if it doesn't exist
	progress("Initializing pixi project...")
	cmd := exec.Command("pixi", "init", "--platform", "linux-64", "--platform", "linux-aarch64")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		progress(fmt.Sprintf("⚠️  Pixi init failed, project may already exist: %v", err))
	}

	// Check if nodejs already exists in the pixi.toml
	pixiContent, err := os.ReadFile("pixi.toml")
	hasNodejs := err == nil && bytes.Contains(pixiContent, []byte("nodejs"))
	hasPython := err == nil && bytes.Contains(pixiContent, []byte("python"))
	hasUv := err == nil && bytes.Contains(pixiContent, []byte("uv"))

	// Add nodejs dependency if not already present
	if !hasNodejs {
		progress("Adding nodejs 22.* to pixi environment...")
		cmd = exec.Command("pixi", "add", "nodejs=22.*")
		stdout.Reset()
		stderr.Reset()
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			msg := fmt.Sprintf("✗ Failed to add nodejs: %v", err)
			progress(msg)
			return false
		}
		progress("✓ nodejs 22.* added to pixi environment")
	} else {
		progress("✓ nodejs 22.* already in pixi environment, skipping")
	}

	// Add python dependency if not already present
	if !hasPython {
		progress("Adding python 3.12.* to pixi environment...")
		cmd = exec.Command("pixi", "add", "python=3.12.*")
		stdout.Reset()
		stderr.Reset()
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			msg := fmt.Sprintf("✗ Failed to add python: %v", err)
			progress(msg)
			return false
		}
		progress("✓ python 3.12.* added to pixi environment")
	} else {
		progress("✓ python 3.12.* already in pixi environment, skipping")
	}

	// Add uv dependency if not already present
	if !hasUv {
		progress("Adding uv to pixi environment...")
		cmd = exec.Command("pixi", "add", "uv")
		stdout.Reset()
		stderr.Reset()
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			msg := fmt.Sprintf("✗ Failed to add uv: %v", err)
			progress(msg)
			return false
		}
		progress("✓ uv added to pixi environment")
	} else {
		progress("✓ uv already in pixi environment, skipping")
	}

	progress("✓ Core dependencies are ready")
	return true
}

// getAliasName returns the desired shell alias name for a given package/tool
func getAliasName(packageName string) string {
	// Map specific package names to their desired aliases
	aliasMap := map[string]string{
		"@sourcegraph/amp@latest": "amp",
		"@openai/codex":           "codex",
		"droid":                   "droid",
		"@google/gemini-cli":      "gemini",
		"kimi-cli":                "kimi",
		"kiro":                    "kiro",
		"opencode-ai":             "opencode",
		"openhands":               "openhands",
		"@qodo/command":           "qodo",
		"@qoder-ai/qodercli":      "qoder",
	}

	// Check if we have a specific mapping
	if alias, exists := aliasMap[packageName]; exists {
		return alias
	}

	// Default: use the package name as-is (shouldn't happen with current tools)
	return packageName
}

// getCommandName returns the actual command name in the pixi environment for a given package/tool
func getCommandName(packageName string) string {
	// Map package names to their actual command names
	commandMap := map[string]string{
		"droid":                   "droid",
		"kimi-cli":                "kimi",
		"kiro":                    "kiro-cli",
		"openhands":               "openhands",
		"@qoder-ai/qodercli":      "qodercli",
		"@sourcegraph/amp@latest": "amp",
	}

	// Check if we have a specific mapping
	if command, exists := commandMap[packageName]; exists {
		return command
	}

	// Default: use the package name as the command name
	return packageName
}

// InstallCLITools installs the selected CLI tools in a new pixi environment
func InstallCLITools(tools []string, installPath string, progress ProgressCallback) []InstallResult {
	results := make([]InstallResult, 0, len(tools))

	if len(tools) == 0 {
		return results
	}

	// Extract tool names from display strings
	toolNames := make([]string, 0, len(tools))
	for _, tool := range tools {
		parts := strings.Split(tool, " - ")
		toolName := strings.TrimSpace(parts[0])
		toolNames = append(toolNames, toolName)
	}

	progress("Installing CLI tools...")

	// Append ai-dev-pixi to the provided parent path
	envDir := installPath + "/ai-dev-pixi"

	// Change to the environment directory
	if err := os.Chdir(envDir); err != nil {
		msg := fmt.Sprintf("✗ Failed to change to directory %s: %v", envDir, err)
		progress(msg)
		return results
	}

	progress(fmt.Sprintf("Using pixi environment: %s", envDir))

	// Install each CLI tool using pixi run
	var stdout, stderr bytes.Buffer
	for _, toolName := range toolNames {
		progress(fmt.Sprintf("Installing %s...", toolName))

		var cmd *exec.Cmd
		var err error

		// Handle special CLI tools installed via curl scripts or custom installers
		if toolName == "droid" {
			cmd = exec.Command("bash", "-c", "curl -fsSL https://app.factory.ai/cli | sh")
			stdout.Reset()
			stderr.Reset()
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
		} else if toolName == "kiro" {
			cmd = exec.Command("bash", "-c", "curl -fsSL https://cli.kiro.dev/install | bash")
			stdout.Reset()
			stderr.Reset()
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
		} else if toolName == "kimi-cli" {
			cmd = exec.Command("pixi", "run", "uv", "tool", "install", "--python", "3.13", "kimi-cli")
			stdout.Reset()
			stderr.Reset()
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
		} else if toolName == "openhands" {
			cmd = exec.Command("pixi", "run", "uv", "tool", "install", "openhands")
			stdout.Reset()
			stderr.Reset()
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
		} else {
			// Install npm packages via pixi
			cmd = exec.Command("pixi", "run", "npm", "install", "-g", toolName)
			stdout.Reset()
			stderr.Reset()
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
		}

		var msg string
		if err == nil {
			msg = fmt.Sprintf("✓ %s installed successfully", toolName)
		} else {
			msg = fmt.Sprintf("✗ Failed to install %s: %v", toolName, err)
		}
		progress(msg)

		results = append(results, InstallResult{
			Name:    toolName,
			Success: err == nil,
			Error:   err,
			Message: msg,
		})
	}

	progress(fmt.Sprintf("📦 CLI tools installed in pixi environment at: %s", envDir))

	// Add aliases to ~/.zshrc for easy access
	homeDir, err := os.UserHomeDir()
	if err == nil {
		zshrcPath := homeDir + "/.zshrc"
		progress("Adding aliases to ~/.zshrc...")

		// Read existing .zshrc content
		existingContent, err := os.ReadFile(zshrcPath)
		if err != nil && !os.IsNotExist(err) {
			progress(fmt.Sprintf("⚠️  Could not read ~/.zshrc: %v", err))
		}

		// Open .zshrc for appending
		f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			progress(fmt.Sprintf("⚠️  Could not open ~/.zshrc for writing: %v", err))
		} else {
			defer f.Close()

			// Add a comment header if this is the first time
			markerComment := "# AI Menu CLI Tool Aliases"
			if !bytes.Contains(existingContent, []byte(markerComment)) {
				f.WriteString("\n" + markerComment + "\n")
			}

			// Add alias for each successfully installed tool
			aliasesAdded := 0
			for _, result := range results {
				if result.Success {
					// Map npm package names to desired alias/command names
					aliasName := getAliasName(result.Name)
					commandName := getCommandName(result.Name)

					// Create alias with the correct command name
					aliasLine := fmt.Sprintf("alias %s='pixi run --manifest-path %s %s'\n",
						aliasName, envDir, commandName)

					// Check if alias already exists
					if !bytes.Contains(existingContent, []byte(aliasLine)) {
						if _, err := f.WriteString(aliasLine); err == nil {
							aliasesAdded++
						}
					}
				}
			}

			if aliasesAdded > 0 {
				progress(fmt.Sprintf("✓ Added %d alias(es) to ~/.zshrc", aliasesAdded))
				progress("Run 'source ~/.zshrc' or restart your shell to use the aliases")
			} else {
				progress("Aliases already exist in ~/.zshrc")
			}
		}
	}

	progress(fmt.Sprintf("To use the tools, run: cd %s && pixi shell", envDir))
	progress("Or use the aliases added to ~/.zshrc (restart shell or run: source ~/.zshrc)")

	return results
}

// InstallVSCodeExtensions installs the selected VS Code extensions
func InstallVSCodeExtensions(extensions []string, progress ProgressCallback) []InstallResult {
	results := make([]InstallResult, 0, len(extensions))

	// Check if code CLI is available
	_, err := exec.LookPath("code")
	if err != nil {
		progress("⚠️  VS Code CLI not found. Extensions cannot be installed automatically.")
		progress("Please install VS Code and ensure 'code' command is in your PATH.")
		return results
	}

	for _, ext := range extensions {
		// Extract extension ID from the display string
		parts := strings.Split(ext, " - ")
		extID := strings.TrimSpace(parts[0])

		progress(fmt.Sprintf("Installing %s...", extID))

		cmd := exec.Command("code", "--install-extension", extID)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		var msg string
		if err == nil {
			msg = fmt.Sprintf("✓ %s installed successfully", extID)
		} else {
			msg = fmt.Sprintf("✗ Failed to install %s: %v", extID, err)
		}
		progress(msg)

		results = append(results, InstallResult{
			Name:    extID,
			Success: err == nil,
			Error:   err,
			Message: msg,
		})
	}

	return results
}

// InstallSpecialTools installs the selected special tools
func InstallSpecialTools(tools []string, progress ProgressCallback) []InstallResult {
	results := make([]InstallResult, 0, len(tools))

	for _, tool := range tools {
		// Extract tool name from the display string
		parts := strings.Split(tool, " - ")
		toolName := strings.TrimSpace(parts[0])

		progress(fmt.Sprintf("Installing %s...", toolName))

		var cmd *exec.Cmd
		switch toolName {
		case "helm":
			cmd = exec.Command("bash", "-c", "curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash")
		case "gh":
			// Install GitHub CLI via apt-get
			cmd = exec.Command("sudo", "apt-get", "install", "-y", "gh")
		case "ripgrep":
			// Install ripgrep via apt-get
			cmd = exec.Command("sudo", "apt-get", "install", "-y", "ripgrep")
		case "jq", "yq", "bat":
			// Install via apt-get
			cmd = exec.Command("sudo", "apt-get", "install", "-y", toolName)
		case "fd":
			// fd is packaged as fd-find in Ubuntu
			cmd = exec.Command("sudo", "apt-get", "install", "-y", "fd-find")
		case "exa":
			// exa has been replaced by eza in Ubuntu 24.04
			cmd = exec.Command("sudo", "apt-get", "install", "-y", "eza")
		case "lazygit":
			cmd = exec.Command("bash", "-c", "LAZYGIT_VERSION=$(curl -s \"https://api.github.com/repos/jesseduffield/lazygit/releases/latest\" | grep -Po '\"tag_name\": \"v\\K[^\"]*') && curl -Lo lazygit.tar.gz \"https://github.com/jesseduffield/lazygit/releases/latest/download/lazygit_${LAZYGIT_VERSION}_Linux_x86_64.tar.gz\" && tar xf lazygit.tar.gz lazygit && sudo install lazygit /usr/local/bin && rm lazygit lazygit.tar.gz")
		default:
			progress(fmt.Sprintf("⚠️  Unknown tool: %s", toolName))
			continue
		}

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		var msg string
		if err == nil {
			msg = fmt.Sprintf("✓ %s installed successfully", toolName)
		} else {
			msg = fmt.Sprintf("✗ Failed to install %s: %v", toolName, err)
		}
		progress(msg)

		results = append(results, InstallResult{
			Name:    toolName,
			Success: err == nil,
			Error:   err,
			Message: msg,
		})
	}

	return results
}
