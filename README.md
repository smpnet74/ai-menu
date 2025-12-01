# AI Menu - Interactive CLI Tool Installer

A beautiful, interactive CLI tool for selecting and installing AI development tools, VS Code extensions, and special utilities.

## Features

- **Five-step workflow:**
  1. Select CLI tools (Amp, Codex, Droid, Gemini CLI, Kimi, Kiro, OpenCode, OpenHands, Plandex, Qodo, Qoder)
  2. Select VS Code extensions (Kilo Code, Zencoder, Augment Code)
  3. Select special tools (helm, jq, bat, exa, lazygit, etc.)
  4. Select CLI tool enhancers (Claude Flow, etc.)
  5. Configure installation path (if CLI tools or enhancers selected)

- **Configurable installation path** - Choose where CLI tools are installed (default: current directory/ai-dev-pixi)
- **Isolated pixi environment** - CLI tools are installed in a dedicated pixi environment with nodejs 22.*
- **Cross-platform support** - Pixi environment configured for linux-64, osx-64, and osx-arm64

- **Interactive TUI** built with Bubble Tea framework
- **Beautiful styling** using Lipgloss
- **Keyboard navigation** - Vi-style (hjkl) or arrow keys
- **Multi-select** - Use spacebar to toggle selections
- **Installation summary** before proceeding

## Prerequisites

- Pixi package manager
- Go 1.21+ (managed by pixi)

## Installation

Clone the repository and navigate to the directory:

```bash
git clone https://github.com/smpnet74/ai-menu.git
cd ai-menu
```

## Usage

### Run with pixi

```bash
pixi run run
```

### Build the binary

```bash
pixi run build
```

Then run directly:

```bash
./ai-menu
```

### Run without pixi shell

You can run the program directly using pixi without entering a shell:

```bash
pixi run -- ./ai-menu
```

## Keyboard Controls

- **↑/k** - Move cursor up
- **↓/j** - Move cursor down
- **Space** - Toggle selection
- **Enter** - Move to next workflow
- **Esc** - Go back to previous screen
- **q / Ctrl+C** - Quit

## Workflows

### 1. CLI Tools Selection
Select from popular AI CLI tools (installed in isolated pixi environment):
- **@sourcegraph/amp@latest** - Amp by Sourcegraph
- **@openai/codex** - OpenAI Codex CLI
- **droid** - Droid by Factory AI
- **@google/gemini-cli** - Google Gemini CLI
- **kimi-cli** - Kimi by MoonshotAI
- **kiro** - Kiro CLI by AWS
- **opencode-ai** - OpenCode CLI
- **openhands** - OpenHands
- **plandex** - Plandex
- **@qodo/command** - Qodo CLI
- **@qoder-ai/qodercli** - Qoder by Qwen

These tools will be installed in a configurable directory (default: current directory + `/ai-dev-pixi`) with nodejs 22.* in a pixi environment that supports linux-64, osx-64, and osx-arm64 platforms.

### 2. Installation Path Configuration
If CLI tools are selected, you'll be prompted to select the parent directory:
- A file picker lets you browse and navigate directories
- Default starting directory: Current working directory (where you ran `ai-menu`)
- Use arrow keys to navigate, Enter to select a directory
- The `ai-dev-pixi` directory will be created inside your chosen parent directory
- For example, if you select `/workspaces/devcontainer`, the installation will be at `/workspaces/devcontainer/ai-dev-pixi`

### 3. VS Code Extensions
Select from AI-powered VS Code extensions:
- **Kilo Code** - AI coding assistant
- **Zencoder** - AI pair programmer
- **Augment Code** - AI code completion

### 4. Special Tools
Select from development utilities:
- Container tools (helm)
- CLI utilities (jq, yq, bat, exa, fd)
- Git tools (lazygit)

### 5. CLI Tool Enhancers
Select from additional tools that enhance CLI functionality:
- **Claude Flow by ruvnet** - Advanced CLI workflow tool

## Installation Notes

### CLI Tools
- **CLI tools are installed in an isolated pixi environment** at a configurable location (default: current directory + `/ai-dev-pixi`)
- You specify a parent directory, and `ai-dev-pixi` is created inside it
- The pixi environment includes nodejs 22.* and is cross-platform (linux-64, linux-aarch64)
- To use the CLI tools after installation, run: `cd <parent-dir>/ai-dev-pixi && pixi shell`
- All npm packages are installed globally within the pixi environment
- Shell aliases for the installed CLI tools, npx, and npm are automatically added to ~/.zshrc for convenient access

### VS Code Extensions
- VS Code extensions require the `code` CLI to be available in your PATH
- Extensions are installed globally in VS Code

### Special Tools
- All tools can be installed automatically
- Some tools may require `sudo` permissions

## Tagging and Pushing Releases

To create and push a new release of the ai-menu project:

1. **Update version in code** (if applicable):
   ```bash
   # Update any version constants in the code if needed
   ```

2. **Commit your changes**:
   ```bash
   git add .
   git commit -m "Release vX.Y.Z: Brief description of changes"
   ```

3. **Create a git tag**:
   ```bash
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   ```

4. **Push the tag to GitHub**:
   ```bash
   git push origin vX.Y.Z
   ```

5. **Create a GitHub release** (optional):
   - Go to the repository on GitHub
   - Navigate to "Releases"
   - Click "Create a new release"
   - Select the tag you just pushed
   - Add release notes describing the changes
   - Publish the release

## Project Structure

```
ai-menu/
├── main.go         # Entry point and main model
├── data.go         # Tool/extension data sources
├── views.go        # UI rendering logic
├── styles.go       # Lipgloss styling
├── handlers.go     # Event handlers and navigation
├── installer.go    # Installation logic
├── pixi.toml       # Pixi configuration
└── README.md       # This file
```

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling library
- [Pixi](https://pixi.sh) - Package manager

## License

MIT
# Updated
