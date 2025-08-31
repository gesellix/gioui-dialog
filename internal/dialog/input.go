package dialog

import (
	"image"
	"image/color"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// inputDialog is the internal implementation stub for a text-input dialog.
type inputDialog struct {
	Width, Height float32
	Title         string
	Label         string
	Description   string
	DefaultText   string
	Validate      func(string) error

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
func NewInputDialog(width, height float32, title, label, description, defaultText string, validate func(string) error) *inputDialog {
	if width <= 0 {
		width = 400
	}
	if height <= 0 {
		height = 200
	}
	d := &inputDialog{
		Width:       width,
		Height:      height,
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

// styledEditor creates an editor with border styling inspired by cu theme
func (d *inputDialog) styledEditor(gtx layout.Context, th *material.Theme) layout.Dimensions {
	cornerRadius := unit.Dp(4)
	inset := unit.Dp(4)

	// Set minimum size for input field
	minWidth := unit.Dp(200)
	minHeight := unit.Dp(32)

	// Apply minimum constraints
	gtx.Constraints.Min.X = max(gtx.Constraints.Min.X, gtx.Dp(minWidth))
	gtx.Constraints.Min.Y = max(gtx.Constraints.Min.Y, gtx.Dp(minHeight))

	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		// Draw the background and border
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rect := image.Rectangle{Max: gtx.Constraints.Min}
			inner := rect.Inset(gtx.Dp(inset))

			// Colors inspired by cu theme
			backgroundColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			borderColor := color.NRGBA{R: 200, G: 200, B: 200, A: 255}
			focusColor := color.NRGBA{R: 0, G: 123, B: 255, A: 255}

			rr := gtx.Dp(cornerRadius)

			// Draw focus border if focused
			if gtx.Focused(&d.textInput) {
				w := gtx.Dp(2)
				paint.FillShape(gtx.Ops, focusColor,
					clip.Stroke{
						Path:  clip.UniformRRect(rect.Inset(w), rr+w).Path(gtx.Ops),
						Width: float32(w),
					}.Op(),
				)
			}

			// Draw background
			shape := clip.UniformRRect(inner, rr)
			defer shape.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, backgroundColor)

			// Draw border
			w := gtx.Dp(1)
			paint.FillShape(gtx.Ops, borderColor,
				clip.Stroke{
					Path:  clip.UniformRRect(inner, rr).Path(gtx.Ops),
					Width: float32(w),
				}.Op(),
			)

			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),

		// Draw the text editor on top
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// Apply minimum constraints to the editor
			gtx.Constraints.Min.X = gtx.Dp(minWidth)
			gtx.Constraints.Min.Y = gtx.Dp(minHeight)

			return layout.Inset{
				Top:    8,
				Bottom: 8,
				Left:   12,
				Right:  12,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, &d.textInput, "")
				editor.TextSize = unit.Sp(14)
				return editor.Layout(gtx)
			})
		}),
	)
}

// Show runs the text-input dialog event loop and returns the entered text,
// a canceled flag, and an error if something went wrong.
func (d *inputDialog) Show() (string, bool, error) {
	w := app.Window{}
	w.Option(
		app.Title(d.Title),
		app.Size(unit.Dp(d.Width), unit.Dp(d.Height)),
	)
	go func() {
		time.Sleep(10 * time.Millisecond)
		w.Perform(system.ActionCenter)
	}()
	w.Perform(system.ActionCenter)

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
				return d.styledEditor(gtx, th)
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
