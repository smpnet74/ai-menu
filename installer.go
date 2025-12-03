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
		"@augmentcode/auggie":     "auggie",
		"@openai/codex":           "codex",
		"droid":                   "droid",
		"forgecode@latest":        "forge",
		"@google/gemini-cli":      "gemini",
		"goose":                   "goose",
		"kimi-cli":                "kimi",
		"kiro":                    "kiro",
		"opencode-ai":             "opencode",
		"openhands":               "openhands",
		"plandex":                 "plandex",
		"@qodo/command":           "qodo",
		"@qoder-ai/qodercli":      "qoder",
		"claude-flow@alpha":       "claude-flow",
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
		"@augmentcode/auggie":     "auggie",
		"droid":                   "droid",
		"forgecode@latest":        "forge",
		"goose":                   "goose",
		"kimi-cli":                "kimi",
		"kiro":                    "kiro-cli",
		"openhands":               "openhands",
		"plandex":                 "plandex",
		"@qodo/command":           "qodo",
		"@qoder-ai/qodercli":      "qodercli",
		"@sourcegraph/amp@latest": "amp",
		"@google/gemini-cli":      "gemini",
		"@openai/codex":           "codex",
		"opencode-ai":             "opencode",
		"claude-flow@alpha":       "claude-flow",
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
		} else if toolName == "goose" {
			cmd = exec.Command("bash", "-c", "curl -fsSL https://github.com/block/goose/releases/download/stable/download_cli.sh | CONFIGURE=false bash")
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
		} else if toolName == "plandex" {
			cmd = exec.Command("bash", "-c", "curl -sL https://plandex.ai/install.sh | bash")
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

			// Add npx and npm aliases if any CLI tools were successfully installed
			hasSuccess := false
			for _, result := range results {
				if result.Success {
					hasSuccess = true
					break
				}
			}
			if hasSuccess {
				npxAliasLine := fmt.Sprintf("alias npx='pixi run --manifest-path %s npx'\n", envDir)
				if !bytes.Contains(existingContent, []byte(npxAliasLine)) {
					if _, err := f.WriteString(npxAliasLine); err == nil {
						progress("✓ Added npx alias to ~/.zshrc")
					}
				}
				npmAliasLine := fmt.Sprintf("alias npm='pixi run --manifest-path %s npm'\n", envDir)
				if !bytes.Contains(existingContent, []byte(npmAliasLine)) {
					if _, err := f.WriteString(npmAliasLine); err == nil {
						progress("✓ Added npm alias to ~/.zshrc")
					}
				}
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
func InstallSpecialTools(tools []string, installPath string, progress ProgressCallback) []InstallResult {
	results := make([]InstallResult, 0, len(tools))

	// Append ai-dev-pixi to the provided parent path for pixi environment
	envDir := installPath + "/ai-dev-pixi"

	// Store current directory to restore later
	originalDir, err := os.Getwd()
	if err != nil {
		progress(fmt.Sprintf("⚠️  Could not get current directory: %v", err))
	}

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
			// Install GitHub CLI using official install script
			cmd = exec.Command("bash", "-c", "curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && echo \"deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main\" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null && sudo apt update && sudo apt install -y gh")
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
		case "modal":
			// Install modal via uv pip in the pixi environment
			// First change to the pixi environment directory
			if err := os.Chdir(envDir); err != nil {
				progress(fmt.Sprintf("✗ Failed to change to directory %s: %v", envDir, err))
				continue
			}
			cmd = exec.Command("pixi", "run", "uv", "pip", "install", "modal")
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

		result := InstallResult{
			Name:    toolName,
			Success: err == nil,
			Error:   err,
			Message: msg,
		}

		results = append(results, result)

		// Add alias for bat to .zshrc if installation was successful
		if toolName == "bat" && result.Success {
			addBatAlias(progress)
		}

		// Add alias for modal to .zshrc if installation was successful
		if toolName == "modal" && result.Success {
			addModalAlias(envDir, progress)
		}
	}

	// Restore original directory
	if originalDir != "" {
		os.Chdir(originalDir)
	}

	return results
}

// addBatAlias adds an alias for bat pointing to batcat in .zshrc
func addBatAlias(progress ProgressCallback) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		progress(fmt.Sprintf("⚠️  Could not get home directory: %v", err))
		return
	}

	zshrcPath := homeDir + "/.zshrc"

	// Read existing .zshrc content
	existingContent, err := os.ReadFile(zshrcPath)
	if err != nil && !os.IsNotExist(err) {
		progress(fmt.Sprintf("⚠️  Could not read ~/.zshrc: %v", err))
		return
	}

	// Check if alias already exists
	batAliasLine := "alias bat='batcat'\n"
	if bytes.Contains(existingContent, []byte(batAliasLine)) {
		progress("✓ bat alias already exists in ~/.zshrc")
		return
	}

	// Open .zshrc for appending
	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		progress(fmt.Sprintf("⚠️  Could not open ~/.zshrc for writing: %v", err))
		return
	}
	defer f.Close()

	// Add a comment header if this is the first time
	specialToolsMarker := "# AI Menu Special Tools Aliases"
	if !bytes.Contains(existingContent, []byte(specialToolsMarker)) {
		f.WriteString("\n" + specialToolsMarker + "\n")
	}

	// Add the bat alias
	if _, err := f.WriteString(batAliasLine); err != nil {
		progress(fmt.Sprintf("⚠️  Could not write to ~/.zshrc: %v", err))
		return
	}

	progress("✓ Added bat alias to ~/.zshrc")
	progress("Run 'source ~/.zshrc' or restart your shell to use the alias")
}

// addModalAlias adds an alias for modal pointing to the modal pip package in .zshrc
func addModalAlias(envDir string, progress ProgressCallback) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		progress(fmt.Sprintf("⚠️  Could not get home directory: %v", err))
		return
	}

	zshrcPath := homeDir + "/.zshrc"

	// Read existing .zshrc content
	existingContent, err := os.ReadFile(zshrcPath)
	if err != nil && !os.IsNotExist(err) {
		progress(fmt.Sprintf("⚠️  Could not read ~/.zshrc: %v", err))
		return
	}

	// Create alias using pixi environment
	modalAliasLine := fmt.Sprintf("alias modal='pixi run --manifest-path %s/pixi.toml python -m modal'\n", envDir)
	if bytes.Contains(existingContent, []byte(modalAliasLine)) {
		progress("✓ modal alias already exists in ~/.zshrc")
		return
	}

	// Open .zshrc for appending
	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		progress(fmt.Sprintf("⚠️  Could not open ~/.zshrc for writing: %v", err))
		return
	}
	defer f.Close()

	// Add a comment header if this is the first time
	specialToolsMarker := "# AI Menu Special Tools Aliases"
	if !bytes.Contains(existingContent, []byte(specialToolsMarker)) {
		f.WriteString("\n" + specialToolsMarker + "\n")
	}

	// Add the modal alias
	if _, err := f.WriteString(modalAliasLine); err != nil {
		progress(fmt.Sprintf("⚠️  Could not write to ~/.zshrc: %v", err))
		return
	}

	progress("✓ Added modal alias to ~/.zshrc")
	progress("Run 'source ~/.zshrc' or restart your shell to use the alias")
}

// InstallCLIEnhancers installs the selected CLI tool enhancers in the pixi environment
func InstallCLIEnhancers(enhancers []string, installPath string, progress ProgressCallback) []InstallResult {
	results := make([]InstallResult, 0, len(enhancers))

	if len(enhancers) == 0 {
		return results
	}

	progress("Installing CLI tool enhancers...")

	// Append ai-dev-pixi to the provided parent path
	envDir := installPath + "/ai-dev-pixi"

	// Change to the environment directory
	if err := os.Chdir(envDir); err != nil {
		msg := fmt.Sprintf("✗ Failed to change to directory %s: %v", envDir, err)
		progress(msg)
		return results
	}

	progress(fmt.Sprintf("Using pixi environment: %s", envDir))

	// Install each CLI enhancer using pixi run
	var stdout, stderr bytes.Buffer
	for _, enhancer := range enhancers {
		// Convert display name to package name
		packageName := getPackageNameForCLIEnhancer(enhancer)
		progress(fmt.Sprintf("Installing %s...", enhancer))

		// Install npm packages via pixi
		cmd := exec.Command("pixi", "run", "npm", "install", "-g", packageName)
		stdout.Reset()
		stderr.Reset()
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		var msg string
		if err == nil {
			msg = fmt.Sprintf("✓ %s installed successfully", enhancer)
		} else {
			msg = fmt.Sprintf("✗ Failed to install %s: %v", enhancer, err)
		}
		progress(msg)

		results = append(results, InstallResult{
			Name:    packageName,
			Success: err == nil,
			Error:   err,
			Message: msg,
		})
	}

	progress(fmt.Sprintf("📦 CLI enhancers installed in pixi environment at: %s", envDir))

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

			// Add alias for each successfully installed enhancer
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

	progress(fmt.Sprintf("To use the enhancers, run: cd %s && pixi shell", envDir))
	progress("Or use the aliases added to ~/.zshrc (restart shell or run: source ~/.zshrc)")

	return results
}