// gioui-dialog/internal/dialog/input.go

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

// inputDialog is the internal implementation stub for a text-input dialog.
type inputDialog struct {
	Title       string
	Label       string
	Description string
	DefaultText string
	Validate    func(string) error

	// internal result state
	result   string
	canceled bool

	// UI state
	textInput    widget.Editor
	okButton     widget.Clickable
	cancelButton widget.Clickable
	done         bool
}

// NewInputDialog initializes an inputDialog from provided parameters.
func NewInputDialog(title, label, description, defaultText string, validate func(string) error) *inputDialog {
	d := &inputDialog{
		Title:       title,
		Label:       label,
		Description: description,
		DefaultText: defaultText,
		Validate:    validate,
	}
	// Initialize text input with default text
	d.textInput.SetText(defaultText)
	d.textInput.SingleLine = true
	return d
}

// Show runs the text-input dialog event loop and returns the entered text,
// a canceled flag, and an error if something went wrong.
func (d *inputDialog) Show() (string, bool, error) {
	w := app.Window{}
	w.Option(
		app.Title(d.Title),
		app.Size(unit.Dp(400), unit.Dp(200)),
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
			// Handle Enter key for OK
			if e.Name == key.NameReturn && e.State == key.Press {
				d.handleOK()
				w.Perform(system.ActionClose)
			}
		case app.DestroyEvent:
			d.done = true
			return d.result, d.canceled, e.Err
		}
	}
	//app.Main()
	return d.result, d.canceled, nil
}

func (d *inputDialog) layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
			// Text input
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, &d.textInput, "")
				editor.TextSize = unit.Sp(14)
				return editor.Layout(gtx)
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

func (d *inputDialog) handleOK() {
	text := d.textInput.Text()
	if d.Validate != nil {
		if err := d.Validate(text); err != nil {
			// TODO: Show validation error in UI
			return
		}
	}
	d.result = text
	d.canceled = false
}

func (d *inputDialog) handleCancel() {
	d.result = ""
	d.canceled = true
}
