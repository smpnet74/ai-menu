package main

// getCLITools returns the list of available CLI tools
func getCLITools() []string {
	return []string{
		"Amp by Sourcegraph",
		"Codex by OpenAI",
		"Droid by Factory AI",
		"Gemini CLI by Google",
		"Kimi by MoonshotAI",
		"Kiro CLI by AWS",
		"OpenCode CLI",
		"OpenHands",
		"Qodo CLI",
		"Qoder by Qwen",
	}
}

// getPackageNameForCLI maps display names to package names for installation
func getPackageNameForCLI(displayName string) string {
	packageMap := map[string]string{
		"Amp by Sourcegraph":   "@sourcegraph/amp@latest",
		"Codex by OpenAI":      "@openai/codex",
		"Droid by Factory AI":  "droid",
		"Gemini CLI by Google": "@google/gemini-cli",
		"Kimi by MoonshotAI":   "kimi-cli",
		"Kiro CLI by AWS":      "kiro",
		"OpenCode CLI":         "opencode-ai",
		"OpenHands":            "openhands",
		"Qodo CLI":             "@qodo/command",
		"Qoder by Qwen":        "@qoder-ai/qodercli",
	}
	if pkg, exists := packageMap[displayName]; exists {
		return pkg
	}
	return displayName
}

// getVSCodeExtensions returns the list of available VS Code extensions
func getVSCodeExtensions() []string {
	return []string{
		"augment.vscode-augment - Augment Code",
		"kilocode.kilo-code - Kilo Code",
		"rooveterinaryinc.roo-cline - Roo Code",
		"saoudrizwan.claude-dev - Cline",
		"zencoderai.zencoder - Zencoder",
	}
}

// getSpecialTools returns the list of special tools
func getSpecialTools() []string {
	return []string{
		"helm - Kubernetes package manager",
		"gh - GitHub CLI",
		"ripgrep - Fast search tool (rg)",
		"jq - JSON processor",
		"yq - YAML processor",
		"bat - Better cat with syntax highlighting",
		"exa - Modern ls replacement (installs eza)",
		"fd - Better find alternative",
		"lazygit - Git TUI",
	}
}
