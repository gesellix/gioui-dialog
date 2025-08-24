package dialog

import (
	internaldialog "github.com/gesellix/gioui-dialog/internal/dialog"
)

// InputDialogOptions holds the configuration for a text-input dialog.
type InputDialogOptions struct {
	Title       string             // Window title
	Label       string             // Prompt label
	Description string             // Additional description or help text
	DefaultText string             // Initial text shown in the input field
	Validate    func(string) error // Optional validation function; return an error on invalid input
}

// PromptInput displays a text-input dialog according to the provided options.
// It returns the entered text, a flag indicating whether the dialog was canceled, and any error.
func PromptInput(opts InputDialogOptions) (result string, canceled bool, err error) {
	dlg := internaldialog.NewInputDialog(opts.Title, opts.Label, opts.Description, opts.DefaultText, opts.Validate)
	return dlg.Show()
}

// SelectDialogOptions holds the configuration for a multiple-selection dialog.
type SelectDialogOptions struct {
	Title             string   // Window title
	Label             string   // Prompt label
	Description       string   // Additional description or help text
	Choices           []string // Available options to select from
	DefaultSelections []string // Options pre-selected when the dialog opens
	AllowCustomEntry  bool     // If true, allows the user to enter a custom value
}

// PromptSelect displays a multi-select dialog according to the provided options.
// It returns the selected items, a flag indicating whether the dialog was canceled, and any error.
func PromptSelect(opts SelectDialogOptions) (selected []string, canceled bool, err error) {
	dlg := internaldialog.NewSelectDialog(opts.Title, opts.Label, opts.Description, opts.Choices, opts.DefaultSelections, opts.AllowCustomEntry)
	return dlg.Show()
}

// BaseDialogOptions holds the configuration for a basic dialog with just title, label, description and OK/Cancel buttons.
type BaseDialogOptions struct {
	Title       string // Window title
	Label       string // Prompt label
	Description string // Additional description or help text
}

// PromptBase displays a base dialog according to the provided options.
// It returns whether the dialog was confirmed, a flag indicating whether it was canceled, and any error.
func PromptBase(opts BaseDialogOptions) (confirmed bool, canceled bool, err error) {
	dlg := internaldialog.NewBaseDialog(opts.Title, opts.Label, opts.Description)
	return dlg.Show()
}
