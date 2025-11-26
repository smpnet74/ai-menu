package main

// getCLITools returns the list of available CLI tools
func getCLITools() []string {
	return []string{
		"@google/gemini-cli - Google Gemini CLI",
		"@qodo/command - Qodo CLI",
		"opencode-ai - OpenCode CLI",
		"@openai/codex - OpenAI Codex CLI",
	}
}

// getVSCodeExtensions returns the list of available VS Code extensions
func getVSCodeExtensions() []string {
	return []string{
		"kilocode.kilo-code - Kilo Code",
		"augment.vscode-augment - Augment Code",
		"zencoderai.zencoder - Zencoder",
	}
}

// getSpecialTools returns the list of special tools
func getSpecialTools() []string {
	return []string{
		"helm - Kubernetes package manager",
		"jq - JSON processor",
		"yq - YAML processor",
		"bat - Better cat with syntax highlighting",
		"exa - Modern ls replacement",
		"fd - Better find alternative",
		"lazygit - Git TUI",
	}
}
