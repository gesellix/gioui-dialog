# GioUI Dialog

A cross-platform dialog library for [Gio](https://gioui.org/) applications, providing zenity-like functionality with native Go implementation.

## Features

- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Native Gio Integration**: Built specifically for Gio applications
- **Multiple Dialog Types**: Text input, single-select, and base dialogs
- **Scrollable Lists**: Selection dialogs support scrolling for large option lists
- **Styled Input Fields**: Professional-looking input fields with borders and focus highlighting
- **Keyboard Shortcuts**: Support for Enter (confirm) and Escape (cancel)
- **Validation Support**: Optional input validation for text dialogs
- **Custom Entries**: Allow custom input in selection dialogs

## Installation

```bash
go get github.com/gesellix/gioui-dialog
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/gesellix/gioui-dialog/pkg/dialog"
)

func main() {
    // Text input dialog
    result, canceled, err := dialog.PromptInput(dialog.InputDialogOptions{
        Title:       "User Input",
        Label:       "Enter your name:",
        Description: "Please provide your full name",
        DefaultText: "John Doe",
    })

    if err != nil {
        log.Fatal(err)
    }

    if !canceled {
        fmt.Printf("Hello, %s!\n", result)
    }
}
```

## Dialog Types

### Text Input Dialog

Prompts the user for text input with optional validation.

```go
result, canceled, err := dialog.PromptInput(dialog.InputDialogOptions{
    Title:       "API Key",
    Label:       "Enter your API key",
    Description: "This will be stored securely",
    DefaultText: "abcd-1234",
    Validate: func(s string) error {
        if len(s) < 4 {
            return fmt.Errorf("API key must be at least 4 characters")
        }
        return nil
    },
})
```

### Single-Select Dialog

Allows users to select one option from a list, with optional custom entry.

```go
selected, canceled, err := dialog.PromptSelect(dialog.SelectDialogOptions{
    Title:            "Choose Language",
    Label:            "Select programming language",
    Description:      "Choose your favorite language to install",
    Choices:          []string{"Go", "Rust", "Python", "Java", "JavaScript"},
    DefaultSelection: "Go",
    AllowCustomEntry: true,
})
```

### Base Dialog

Simple confirmation dialog with OK/Cancel buttons.

```go
confirmed, canceled, err := dialog.PromptBase(dialog.BaseDialogOptions{
    Title:       "Confirmation",
    Label:       "Are you sure?",
    Description: "This action cannot be undone.",
})
```

## API Reference

### InputDialogOptions

| Field | Type | Description |
|-------|------|-------------|
| `Title` | `string` | Window title |
| `Label` | `string` | Main prompt text |
| `Description` | `string` | Additional help text (optional) |
| `DefaultText` | `string` | Pre-filled text in input field |
| `Validate` | `func(string) error` | Input validation function (optional) |

### SelectDialogOptions

| Field | Type | Description |
|-------|------|-------------|
| `Title` | `string` | Window title |
| `Label` | `string` | Main prompt text |
| `Description` | `string` | Additional help text (optional) |
| `Choices` | `[]string` | Available options to select from |
| `DefaultSelection` | `string` | Pre-selected option |
| `AllowCustomEntry` | `bool` | Allow user to enter custom values |

### BaseDialogOptions

| Field | Type | Description |
|-------|------|-------------|
| `Title` | `string` | Window title |
| `Label` | `string` | Main prompt text |
| `Description` | `string` | Additional help text (optional) |

## Demo Application

Run the included demo to see all dialog types in action:

```bash
go run ./cmd/gioui-dialog
```

The demo application provides buttons to test each dialog type and displays the results.

## Keyboard Shortcuts

- **Enter**: Confirm/OK (in text dialogs, also works when input field has focus)
- **Escape**: Cancel/Close dialog

## Development

### Building

```bash
go build ./cmd/gioui-dialog
```

### Project Structure

```
.
├── cmd/gioui-dialog/           # Demo application
│   └── main.go
├── pkg/dialog/                 # Public API
│   └── dialog.go
├── internal/dialog/            # Internal implementations
│   ├── base.go                 # Base dialog
│   ├── input.go               # Text input dialog
│   └── select.go              # Single-select dialog
├── SPEC.md                    # Technical specification
├── README.md                  # This file
├── LICENSE                    # MIT License
└── go.mod                     # Go module definition
```

## Requirements

- Go 1.24 or later
- Gio v0.8.0 or later

## Platform Support

- **Windows**: Full support
- **macOS**: Full support  
- **Linux**: Full support (requires X11 or Wayland)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Gio](https://gioui.org/) - immediate mode GUI for Go
- Inspired by [zenity](https://github.com/ncruces/zenity) and similar dialog libraries
- Thanks to the Gio community for their excellent documentation and examples
