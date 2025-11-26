# CLAUDE.md - AI Menu Development Documentation

This document provides comprehensive documentation about the AI Menu project, including architectural decisions, lessons learned, and the reasoning behind key implementation choices.

## Project Overview

**AI Menu** is an interactive TUI (Terminal User Interface) application built with Go and Bubble Tea that helps users install AI development tools, VS Code extensions, and utilities in an organized, user-friendly way.

### Key Features
- Beautiful TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Four-step workflow for selecting tools, extensions, and installation paths
- Isolated pixi environment for CLI tools with nodejs 22.*
- Cross-platform support (linux-64, osx-64, osx-arm64)
- Real-time installation progress with spinner
- File picker for directory selection

---

## Architecture

### File Structure

```
ai-menu/
├── main.go         # Entry point, model definition, message handling
├── data.go         # Tool/extension data sources
├── views.go        # UI rendering for each screen
├── styles.go       # Lipgloss styling definitions
├── handlers.go     # Navigation and event handling logic
├── installer.go    # Installation logic with progress callbacks
├── pixi.toml       # Pixi configuration with Go and compilers
└── README.md       # User-facing documentation
```

### State Machine

The application uses a state machine pattern with these states:

```go
const (
    cliToolsView       // Select CLI tools (4 tools)
    vscodeExtensionsView // Select VS Code extensions (3 extensions)
    specialToolsView   // Select special tools (7 tools)
    pathInputView      // File picker for directory selection
    installView        // Review summary before installation
    installingView     // Shows spinner and progress messages
    doneView          // Shows installation results
    quitView          // Exit application
)
```

### Data Flow

1. **User Selection**: Checkboxes with space bar toggle
2. **Path Selection**: File picker navigation (only if CLI tools selected)
3. **Summary Review**: Shows all selections and installation path
4. **Installation**: Async execution with progress callbacks
5. **Results Display**: Success/failure counts with details

---

## Key Technical Decisions

### 1. Pixi Environment for CLI Tools

**Decision**: Create an isolated pixi environment at `<parent-dir>/ai-dev-pixi`

**Reasoning**:
- Isolates nodejs 22.* and CLI tools from system environment
- Cross-platform package management
- Reproducible installations
- User can easily delete the entire environment without affecting system

**Implementation**: `installer.go:17-119`

```go
// Creates directory structure:
// /path/to/parent/
//   └── ai-dev-pixi/
//       ├── pixi.toml
//       ├── pixi.lock
//       └── .pixi/
//           └── envs/
//               └── default/
//                   ├── bin/
//                   └── lib/
```

### 2. Progress Callback System

**Decision**: Use callback functions instead of channels for progress reporting

**Reasoning**:
- Simpler to implement in installer functions
- Uses `program.Send()` to send messages to Bubble Tea
- Avoids complex channel management

**Implementation**: `installer.go:18`

```go
type ProgressCallback func(message string)

// Installer calls progress() for each step
progress("Installing @qodo/command...")
```

**Challenge Faced**: Progress messages were leaking to stdout instead of staying in TUI.

**Attempted Solution**: Added `tea.WithAltScreen()` to use alternate screen buffer.

**Result**: Broke the installation completely - screen became unresponsive.

**Final Solution**: Rolled back to original approach. Progress messages appear outside TUI during installation but functionality works. The completion screen properly shows results inside TUI.

**Code**: `main.go:196-200`

### 3. Directory Selection: Text Input → File Picker

**Original Implementation**: Text input field where users typed directory paths

**Problem**:
- Easy to make typos in paths
- No validation that directory exists
- Users could enter invalid paths

**Solution**: Switched to `bubbles/filepicker` component

**Benefits**:
- Visual navigation through filesystem
- Only existing directories can be selected
- Arrow key navigation
- Preview of installation path

**Implementation**: `main.go:57-62, views.go:122-152`

**Migration Notes**: Removed `textinput.Model`, added `filepicker.Model`. Changed `pathInput` field to `filepicker` field.

### 4. Default Directory: Home → Current Working Directory

**Original**: Default to `$HOME` (`~/`)

**Changed To**: Current working directory where `./ai-menu` is run

**Reasoning**:
- More intuitive - users expect things installed where they run commands
- Easier to keep projects organized
- No need to cd to home and back

**Implementation**: `main.go:51-55`

```go
currentDir, err := os.Getwd()
if err != nil {
    currentDir = "."
}
```

### 5. Reduced Tool Selection

**Original Plans**:
- CLI Tools: 10 tools
- VS Code Extensions: 10 extensions
- Special Tools: 14 tools

**Final Selection**:
- **CLI Tools (4)**:
  - @google/gemini-cli
  - @qodo/command
  - opencode-ai
  - @openai/codex

- **VS Code Extensions (3)**:
  - kilocode.kilo-code
  - augment.vscode-augment
  - zencoderai.zencoder

- **Special Tools (7)**:
  - helm
  - jq, yq
  - bat, exa, fd
  - lazygit

**Removed Tools**: docker, kubectl, terraform, ripgrep, fzf, gh, ansible

**Reasoning**:
- Focus on AI-specific tools only
- Removed infrastructure/DevOps tools (docker, kubectl, terraform)
- Removed tools commonly already installed (ripgrep, fzf, gh)
- Simplified maintenance and testing

**Code Changes**:
- `data.go:4-40` - Updated lists
- `installer.go:177-190` - Removed install cases
- `README.md` - Updated documentation

---

## Installation System

### How It Works

1. **CLI Tools Installation** (`installer.go:17-119`)
   ```
   1. Create parent directory if needed
   2. Run: pixi init --platform linux-64 --platform osx-64 --platform osx-arm64
   3. Run: pixi add nodejs=22.*
   4. For each tool: pixi run npm install -g <tool>
   ```

2. **VS Code Extensions** (`installer.go:121-163`)
   ```
   1. Check if 'code' CLI exists
   2. For each extension: code --install-extension <extension-id>
   ```

3. **Special Tools** (`installer.go:165-225`)
   ```
   Each tool has custom installation:
   - helm: curl script installation
   - jq, yq, bat, fd: apt-get install
   - exa: apt-get install
   - lazygit: GitHub release download
   ```

### Progress Reporting

**Architecture**:
```
Installer Function
  └─> calls progress("message")
       └─> program.Send(installMsg{message: msg})
            └─> Update() receives installMsg
                 └─> appends to m.installMessages
                      └─> View() renders last 10 messages
```

**Limitation**: Messages appear in stdout AND TUI. This is acceptable because:
- Installation still completes successfully
- Results screen shows proper summary
- Alternative approaches (alternate screen) broke functionality

---

## Message Types

### Custom Messages

```go
type installMsgStart struct{}
// Triggers installation, transitions to installingView

type installMsg struct{ message string }
// Progress message from installer

type installCompleteMsg struct{ results []InstallResult }
// Installation finished with results

type InstallResult struct {
    Name    string
    Success bool
    Error   error
    Message string
}
```

### Flow

```
installView (press enter)
  └─> Send installMsgStart
       └─> Update() receives it
            └─> state = installingView
            └─> spinner starts
            └─> performInstallation() runs in goroutine
                 └─> sends installMsg for each step
                 └─> sends installCompleteMsg when done
                      └─> state = doneView
                      └─> spinner stops
```

---

## Styling with Lipgloss

All styles defined in `styles.go`:

```go
titleStyle       // Bold, bordered, purple
selectedItemStyle // Bold, purple (cursor on item)
normalItemStyle   // White text
checkedStyle      // Green, bold (selected checkboxes)
uncheckedStyle    // Gray (unselected checkboxes)
helpStyle         // Gray, bottom help text
summaryStyle      // Orange, bold (section headers)
```

**Usage Pattern**:
```go
b.WriteString(titleStyle.Render("🚀 Select CLI Tools"))
b.WriteString(checkedStyle.Render("[✓]"))
```

---

## Navigation

### Keyboard Controls

**Selection Views** (CLI Tools, VS Code Extensions, Special Tools):
- `↑/k` - Move cursor up
- `↓/j` - Move cursor down
- `Space` - Toggle selection
- `Enter` - Next screen
- `q/Ctrl+C` - Quit

**File Picker View**:
- `↑/↓` - Navigate directories
- `Enter` - Select directory
- `q` - Go back
- `Ctrl+C` - Quit

**Installation Views**:
- Installing: No input (wait)
- Done: `Enter/q/Space` - Exit

### State Transitions

```
cliToolsView
  └─> [Enter] → vscodeExtensionsView
       └─> [Enter] → specialToolsView
            └─> [Enter] → pathInputView (if CLI tools selected)
                 │          OR installView (if no CLI tools)
                 └─> [Select Dir] → installView
                      └─> [Enter] → installingView
                           └─> [Auto] → doneView
                                └─> [Enter] → quit
```

---

## Lessons Learned

### 1. Alternate Screen Mode Issues

**What Happened**: Added `tea.WithAltScreen()` to prevent stdout leakage.

**Result**: Installation broke completely - TUI became unresponsive, wouldn't exit.

**Lesson**: Alternate screen mode requires careful management of stdin/stdout. When installer functions write to stdout, they need special handling.

**Solution**: Removed alternate screen mode, accepted minor stdout leakage during installation.

**Code Reference**: Reverted in `main.go:196`

### 2. TextInput vs FilePicker Trade-offs

**TextInput Pros**:
- Simple to implement
- Users can type any path
- Fast for power users

**TextInput Cons**:
- Typos are easy
- No validation
- Confusing for new users

**FilePicker Pros**:
- Visual confirmation
- Only valid directories
- Harder to make mistakes

**FilePicker Cons**:
- More complex to integrate
- Takes more screen space
- Requires understanding filepicker API

**Decision**: FilePicker wins for better UX

### 3. Progress Callback Pattern

**Why Not Channels?**
```go
// Could have done this:
progressChan := make(chan string)
go func() {
    for msg := range progressChan {
        program.Send(installMsg{message: msg})
    }
}()
```

**Problems**:
- Need to manage goroutine lifecycle
- Need to close channels properly
- More complexity

**Callback Approach**:
```go
progress := func(msg string) {
    program.Send(installMsg{message: msg})
}
```

**Benefits**:
- Simpler to understand
- No channel management
- Direct message sending

### 4. Global Program Variable

**Implementation**: `var program *tea.Program` in `main.go:193`

**Reasoning**: Needed to call `program.Send()` from installer goroutine.

**Trade-off**: Global state isn't ideal, but:
- Alternative is complex channel passing
- Only one program instance exists
- Simplifies callback implementation

---

## Building for Multiple Platforms

### Current Build

```bash
pixi run build  # Builds for current platform
```

### Cross-Platform Build Commands

```bash
# Linux x86_64
GOOS=linux GOARCH=amd64 go build -o ai-menu-linux-amd64 .

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o ai-menu-darwin-amd64 .

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o ai-menu-darwin-arm64 .
```

### Future Enhancement

Could add to `pixi.toml`:
```toml
[tasks]
build-all = { depends-on = ["build-linux", "build-macos-intel", "build-macos-arm"] }
build-linux = "GOOS=linux GOARCH=amd64 go build -o ai-menu-linux-amd64 ."
build-macos-intel = "GOOS=darwin GOARCH=amd64 go build -o ai-menu-darwin-amd64 ."
build-macos-arm = "GOOS=darwin GOARCH=arm64 go build -o ai-menu-darwin-arm64 ."
```

---

## Testing Considerations

### Manual Testing Checklist

1. **Selection Flow**:
   - [ ] Can select/deselect tools with space
   - [ ] Cursor moves with arrow keys and j/k
   - [ ] Enter progresses to next screen

2. **File Picker**:
   - [ ] Shows current directory contents
   - [ ] Can navigate up/down directory tree
   - [ ] Shows preview path with /ai-dev-pixi
   - [ ] Selecting directory progresses to summary

3. **Installation**:
   - [ ] Spinner appears during installation
   - [ ] Progress messages update (even if outside TUI)
   - [ ] Done screen shows correct success/failure counts
   - [ ] Can exit after completion

4. **Edge Cases**:
   - [ ] No tools selected - should still show summary
   - [ ] Only VS Code extensions - skips path selection
   - [ ] Press q at various screens - should quit or go back
   - [ ] Invalid directory permissions - shows error

### Known Issues

1. **Progress Messages Outside TUI**:
   - During installation, messages appear in terminal outside TUI box
   - This is cosmetic; functionality works correctly
   - Attempted fix (alternate screen) broke functionality

2. **Filepicker Performance**:
   - May be slow in directories with thousands of files
   - Could add file count limit in future

---

## Future Enhancements

### Potential Improvements

1. **Better Progress Display**:
   - Investigate proper alternate screen buffer usage
   - Research how other TUI apps handle child process output
   - Consider running installers in isolated pty

2. **Configuration File**:
   - Save user preferences (default directory, favorite tools)
   - Remember last selections
   - Custom tool additions

3. **Installation Validation**:
   - Verify tools installed successfully
   - Check tool versions
   - Post-install health check

4. **Uninstall Feature**:
   - Remove tools from pixi environment
   - Uninstall VS Code extensions
   - Clean up ai-dev-pixi directory

5. **Update Feature**:
   - Check for newer tool versions
   - Update installed tools
   - Update ai-menu itself

---

## Dependencies

### Direct Dependencies

```go
github.com/charmbracelet/bubbletea    // TUI framework
github.com/charmbracelet/bubbles      // TUI components
  ├── spinner                          // Loading animation
  └── filepicker                       // Directory browser
github.com/charmbracelet/lipgloss      // Styling
```

### External Tools Required

- **pixi**: Package manager for creating environments
- **nodejs**: Installed by pixi (22.*)
- **npm**: Comes with nodejs
- **code** (optional): For VS Code extension installation
- **sudo** (optional): For special tools that need system install

---

## Troubleshooting Guide

### Installation Fails

**Symptom**: Tools show as failed in results

**Possible Causes**:
1. No internet connection
2. npm registry unreachable
3. Insufficient permissions
4. Pixi not installed

**Debug**:
```bash
# Check pixi
pixi --version

# Check nodejs in environment
cd ~/ai-dev-pixi
pixi shell
node --version
npm --version
```

### VS Code Extensions Don't Install

**Symptom**: Warning about 'code' CLI not found

**Solution**:
1. Install VS Code
2. Open VS Code
3. Cmd+Shift+P → "Shell Command: Install 'code' command in PATH"
4. Restart terminal
5. Run `code --version` to verify

### Directory Not Found

**Symptom**: Filepicker shows empty or wrong directory

**Solution**:
- The filepicker starts in current working directory
- Make sure you're running `./ai-menu` from desired location
- Or navigate using filepicker to desired location

---

## Code Conventions

### Naming

- **Views**: `render*()` functions in `views.go`
- **Handlers**: `handle*()` functions in `handlers.go`
- **Data**: `get*()` functions in `data.go`
- **Installers**: `Install*()` functions in `installer.go`

### Error Handling

- Most errors are logged with `progress()` callback
- Critical errors (can't create directory) return early
- Non-critical errors (tool install fails) continue to next tool

### Styling

- All lipgloss styles defined in `styles.go`
- Use existing styles rather than inline styling
- Keep colors consistent across views

---

## Version History

### Development Journey

1. **Initial Setup**
   - Created basic TUI structure
   - Added tool selection views
   - Implemented installation outside TUI

2. **Progress Display**
   - Added spinner and progress messages
   - Moved installation into TUI
   - Fought with stdout leakage issue

3. **Alternate Screen Attempt**
   - Added `tea.WithAltScreen()`
   - Installation broke
   - Rolled back changes

4. **Directory Selection Evolution**
   - Started with textinput
   - Realized validation issues
   - Migrated to filepicker

5. **Tool Curation**
   - Reduced from 34 total tools to 14
   - Focused on AI-specific tools only
   - Removed redundant utilities

6. **Polish**
   - Updated documentation
   - Added CLAUDE.md
   - Cleaned up comments

---

## Contributing Guidelines

### Adding New Tools

1. **Add to data.go**:
   ```go
   func getCLITools() []string {
       return []string{
           // ...existing tools
           "new-tool - Description",
       }
   }
   ```

2. **Add installer case** in `installer.go`:
   ```go
   case "new-tool":
       cmd = exec.Command("npm", "install", "-g", "new-tool")
   ```

3. **Update README.md** with new tool info

4. **Test installation** thoroughly

### Code Style

- Run `gofmt` before committing
- Keep functions under 50 lines when possible
- Add comments for non-obvious logic
- Use descriptive variable names

### Testing Changes

```bash
# Build
pixi run build

# Test
./ai-menu

# Test different scenarios:
# - Select each tool type
# - Navigate filepicker
# - Complete installation
# - Verify results screen
```

---

## Credits

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) by Charm
- [Bubbles](https://github.com/charmbracelet/bubbles) by Charm
- [Lipgloss](https://github.com/charmbracelet/lipgloss) by Charm
- [Pixi](https://pixi.sh) by prefix.dev

Created for the AI developer community to simplify setup of AI coding tools.

---

## Contact & Support

For issues, questions, or contributions:
- Create an issue in the repository
- Check README.md for basic usage
- Refer to this CLAUDE.md for development details

---

*This documentation was created to help future developers (including Claude AI instances) understand the project's architecture, decisions, and lessons learned during development.*
