// gioui-dialog/internal/dialog/select.go

package dialog

import (
	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// selectDialog is the internal implementation stub for a single-select dialog.
type selectDialog struct {
	Title            string
	Label            string
	Description      string
	Choices          []string
	DefaultSelection string
	AllowCustomEntry bool

	// internal result state
	selected string
	canceled bool

	// UI state
	selectedIndex int
	choiceButtons []widget.Clickable
	customInput   widget.Editor
	okButton      widget.Clickable
	cancelButton  widget.Clickable
	done          bool
}

// NewSelectDialog initializes a selectDialog from provided parameters.
func NewSelectDialog(title, label, description string, choices []string, defaultSelection string, allowCustomEntry bool) *selectDialog {
	d := &selectDialog{
		Title:            title,
		Label:            label,
		Description:      description,
		Choices:          choices,
		DefaultSelection: defaultSelection,
		AllowCustomEntry: allowCustomEntry,
		selectedIndex:    -1, // No selection initially
	}

	// Initialize clickable buttons for each choice
	d.choiceButtons = make([]widget.Clickable, len(choices))

	// Set default selection
	for i, choice := range choices {
		if choice == defaultSelection {
			d.selectedIndex = i
			break
		}
	}

	// Initialize custom input if allowed
	if allowCustomEntry {
		d.customInput.SingleLine = true
	}

	return d
}

// Show runs the single-selection dialog event loop and returns the selected
// item, a canceled flag, and an error if something went wrong.
func (d *selectDialog) Show() (string, bool, error) {
	w := app.Window{}
	w.Option(
		app.Title(d.Title),
		app.Size(unit.Dp(400), unit.Dp(300)),
	)

	th := material.NewTheme()
	var ops op.Ops

	for !d.done {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if d.cancelButton.Clicked(gtx) {
				d.handleCancel()
				w.Perform(system.ActionClose)
			}
			if d.okButton.Clicked(gtx) {
				d.handleOK()
				w.Perform(system.ActionClose)
			}
			d.layout(gtx, th)
			e.Frame(gtx.Ops)
		case key.Event:
			// Handle Escape key for cancel
			if e.Name == key.NameEscape && e.State == key.Press {
				d.handleCancel()
				w.Perform(system.ActionClose)
			}
		case app.DestroyEvent:
			d.done = true
			return d.selected, d.canceled, e.Err
		}
	}
	//app.Main()
	return d.selected, d.canceled, nil
}

func (d *selectDialog) layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceAround}.Layout(gtx,
			// Label
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(th, d.Label)
				return label.Layout(gtx)
			}),
			// Description
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if d.Description == "" {
					return layout.Dimensions{}
				}
				desc := material.Body1(th, d.Description)
				desc.Color = th.Fg
				return desc.Layout(gtx)
			}),
			// Choices with checkboxes
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					d.choiceItems(gtx, th)...,
				)
			}),
			// Custom entry if allowed
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if !d.AllowCustomEntry {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(th, "Other: ")
						return label.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, &d.customInput, "")
						editor.TextSize = unit.Sp(14)
						return editor.Layout(gtx)
					}),
				)
			}),
			// Buttons
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &d.cancelButton, "Cancel")
						return btn.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &d.okButton, "OK")
						return btn.Layout(gtx)
					}),
				)
			}),
		)
	})
}

func (d *selectDialog) choiceItems(gtx layout.Context, th *material.Theme) []layout.FlexChild {
	var items []layout.FlexChild

	for i, choice := range d.Choices {
		i, choice := i, choice // capture loop variables
		items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Check if this item was clicked
			if d.choiceButtons[i].Clicked(gtx) {
				d.selectedIndex = i
			}

			// Create button style with enhanced selection indicator
			var buttonText string
			btn := material.Button(th, &d.choiceButtons[i], "")

			if d.selectedIndex == i {
				// Selected item: use high contrast colors and add checkmark
				btn.Background = th.Palette.ContrastBg
				btn.Color = th.Palette.ContrastFg
				buttonText = "✓ " + choice
			} else {
				// Unselected item: use subtle styling
				btn.Background = th.Bg
				btn.Color = th.Fg
				buttonText = "  " + choice
			}

			btn.Text = buttonText
			return btn.Layout(gtx)
		}))
	}

	return items
}

func (d *selectDialog) handleOK() {
	// Check if custom entry is provided and not empty
	if d.AllowCustomEntry {
		customText := d.customInput.Text()
		if customText != "" {
			d.selected = customText
			d.canceled = false
			d.done = true
			return
		}
	}

	// Use selected choice if any
	if d.selectedIndex >= 0 && d.selectedIndex < len(d.Choices) {
		d.selected = d.Choices[d.selectedIndex]
	} else {
		d.selected = ""
	}
	d.canceled = false
}

func (d *selectDialog) handleCancel() {
	d.selected = ""
	d.canceled = true
}
