package main

// getCLITools returns the list of available CLI tools
func getCLITools() []string {
	return []string{
		"Amp by Sourcegraph",
		"Auggie by Augment Code",
		"Codex by OpenAI",
		"Droid by Factory AI",
		"Forgecode",
		"Gemini CLI by Google",
		"Goose",
		"Kimi by MoonshotAI",
		"Kiro CLI by AWS",
		"OpenCode CLI",
		"OpenHands",
		"Plandex",
		"Qodo CLI",
		"Qoder by Qwen",
	}
}

// getPackageNameForCLI maps display names to package names for installation
func getPackageNameForCLI(displayName string) string {
	packageMap := map[string]string{
		"Amp by Sourcegraph":      "@sourcegraph/amp@latest",
		"Auggie by Augment Code":  "@augmentcode/auggie",
		"Codex by OpenAI":         "@openai/codex",
		"Droid by Factory AI":     "droid",
		"Forgecode":               "forgecode@latest",
		"Gemini CLI by Google":    "@google/gemini-cli",
		"Goose":                   "goose",
		"Kimi by MoonshotAI":      "kimi-cli",
		"Kiro CLI by AWS":         "kiro",
		"OpenCode CLI":            "opencode-ai",
		"OpenHands":               "openhands",
		"Plandex":                 "plandex",
		"Qodo CLI":                "@qodo/command",
		"Qoder by Qwen":           "@qoder-ai/qodercli",
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

// getCLIEnhancers returns the list of CLI tool enhancers
func getCLIEnhancers() []string {
	return []string{
		"Claude Flow by ruvnet - Claude CLI enhancer",
	}
}

// getPackageNameForCLIEnhancer maps display names to package names for installation
func getPackageNameForCLIEnhancer(displayName string) string {
	packageMap := map[string]string{
		"Claude Flow by ruvnet - Claude CLI enhancer": "claude-flow@alpha",
	}
	if pkg, exists := packageMap[displayName]; exists {
		return pkg
	}
	return displayName
}
