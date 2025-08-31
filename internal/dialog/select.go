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

// selectDialog is the internal implementation stub for a single-select dialog.
type selectDialog struct {
	Width, Height    float32
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
	list          layout.List
	done          bool
}

// NewSelectDialog initializes a selectDialog from provided parameters.
func NewSelectDialog(width, height float32, title, label, description string, choices []string, defaultSelection string, allowCustomEntry bool) *selectDialog {
	if width <= 0 {
		width = 400
	}
	if height <= 0 {
		height = 300
	}
	d := &selectDialog{
		Width:            width,
		Height:           height,
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

	// Initialize scrollable list
	d.list.Axis = layout.Vertical

	return d
}

// styledEditor creates an editor with border styling inspired by cu theme
func (d *selectDialog) styledEditor(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
			if gtx.Focused(&d.customInput) {
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
				editor := material.Editor(th, &d.customInput, "")
				editor.TextSize = unit.Sp(14)
				return editor.Layout(gtx)
			})
		}),
	)
}

// Show runs the single-selection dialog event loop and returns the selected
// item, a canceled flag, and an error if something went wrong.
func (d *selectDialog) Show() (string, bool, error) {
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
			// Choices with scrollable list
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// Set height to show approximately 3 items (button height ~40dp + spacing)
				maxHeight := unit.Dp(140)
				gtx.Constraints.Max.Y = gtx.Dp(maxHeight)

				return d.list.Layout(gtx, len(d.Choices), func(gtx layout.Context, i int) layout.Dimensions {
					return d.choiceItem(gtx, th, i)
				})
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
						return d.styledEditor(gtx, th)
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

func (d *selectDialog) choiceItem(gtx layout.Context, th *material.Theme, i int) layout.Dimensions {
	choice := d.Choices[i]

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
}

func (d *selectDialog) handleOK() {
	// Check if custom entry is provided and not empty
	if d.AllowCustomEntry {
		customText := d.customInput.Text()
		if customText != "" {
			d.selected = customText
			d.canceled = false
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
