# Specification for a Reusable GioUI Dialog Module

This specification defines the architecture, UI components, and developer API for a cross-platform GioUI-based dialog module that provides zenity-like functionality.

---

## 1. Goals and Requirements

1. Reusability  
   - Unified API for all platforms (Windows, macOS, Linux).  
   - Easy integration into existing Gio apps.

2. Dialog Types  
   - **Text Input** with display of a default value.  
   - **Single Selection** with display of current default.  
   - Cancel and Confirm in every dialog.

3. Extensibility  
   - Support for manual input also in single selection dialog.  
   - Optionally later additional dialog types (password, file selection etc.).

---

## 2. UI Components

### 2.1 Base Dialog (`BaseDialog`)
Every dialog inherits standard elements:
- **Title** (window title)  
- **Label** (heading in dialog)  
- **Description** (help text)  
- **Buttons**: OK/Confirm, Cancel/Abort

### 2.2 Text Input Dialog (`InputDialog`)
- Input field (single-line)  
- Default value initially filled  
- Validation (optional, e.g. regex, length)  
- UX: Focus on text field, Enter=Confirm, Esc=Cancel

### 2.3 Single-Select Dialog (`SelectDialog`)
Recommendation: Scrollable list with clickable items  
- List of all `Choices` with clickable selection buttons  
- "Other..." input field as last line when `AllowCustomEntry=true`  
- Pre-selected default as highlighted selection  
- UX: Scrollbar automatically shown for many entries (displays ~3 items at once), buttons as usual

---

## 3. Developer API

```go
// Options for text input dialog
type InputDialogOptions struct {
    Width, Height float32            // Dimensions of the dialog window
    Title         string             // Window title
    Label         string             // Prompt label
    Description   string             // Additional description or help text
    DefaultText   string             // Initial text shown in the input field
    Validate      func(string) error // Optional validation function; return an error on invalid input
}

// Call: shows text input dialog
func PromptInput(opts InputDialogOptions) (result string, canceled bool, err error)
```

```go
// Options for single-select dialog
type SelectDialogOptions struct {
    Width, Height    float32  // Dimensions of the dialog window
    Title            string   // Window title
    Label            string   // Prompt label
    Description      string   // Additional description or help text
    Choices          []string // Available options to select from
    DefaultSelection string   // Option pre-selected when the dialog opens
    AllowCustomEntry bool     // If true, allows the user to enter a custom value
}

// Call: shows single-select dialog
func PromptSelect(opts SelectDialogOptions) (selected string, canceled bool, err error)
```

```go
// Options for base dialog
type BaseDialogOptions struct {
    Width, Height float32 // Dimensions of the dialog window
    Title         string  // Window title
    Label         string  // Prompt label
    Description   string  // Additional description or help text
}

// Call: shows base dialog
func PromptBase(opts BaseDialogOptions) (confirmed bool, canceled bool, err error)
```

### 3.1 Optional: Builder Pattern

```go
res, canceled, err := NewInputDialog().
    Title("Name").
    Label("Enter your name").
    Default("guest").
    Validate(nonEmpty).
    Show()
```

---

## 4. Cross-Platform Strategy

- Gio abstracts render pipelines (OpenGL, Metal, DirectX).  
- Unified styles via `material` package.  
- CI tests on Windows, macOS, Linux.

---

## 5. UX Notes & Recommendations

- Clear self-explanatory single-select lists with clickable buttons.  
- "Other..." field only visible when needed.  
- Highlight selected item clearly (e.g. background color contrast).

---

## 6. Example Workflows

### 6.1 Text Query

```go
val, canceled, err := PromptInput(InputDialogOptions{
    Width:       400,
    Height:      200,
    Title:       "API Key",
    Label:       "Enter your API key",
    Description: "Will be stored on the server.",
    DefaultText: "abcd-1234",
})
```

### 6.2 Single Selection with Custom Entry

```go
chosen, canceled, err := PromptSelect(SelectDialogOptions{
    Width:            450,
    Height:           300,
    Title:            "Choose Language",
    Label:            "Select language to install",
    Choices:          []string{"Go", "Rust", "Python", "Java"},
    DefaultSelection: "Go",
    AllowCustomEntry: true,
})
```

---

## 7. Next Steps

1. Implement prototype of base components (`BaseDialog`, `InputDialog`, `SelectDialog`).  
2. Define style guide (spacing, fonts, colors).  
3. Write unit tests for API & validation.  
4. Integration tests on all platforms.
