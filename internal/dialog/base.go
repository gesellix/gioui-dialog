package dialog

import (
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// BaseDialog contains common properties and (later) shared behavior
// used by concrete dialogs (InputDialog, SelectDialog).
// According to SPEC.md it provides Title, Label, Description and
// standard OK/Cancel handling.
type BaseDialog struct {
	Width, Height float32
	Title         string
	Label         string
	Description   string

	// internal result state
	confirmed bool
	canceled  bool

	// UI state
	okButton     widget.Clickable
	cancelButton widget.Clickable
	done         bool
}

// NewBaseDialog creates a new BaseDialog with the standard fields.
func NewBaseDialog(width, height float32, title, label, description string) *BaseDialog {
	if width <= 0 {
		width = 400
	}
	if height <= 0 {
		height = 180
	}
	return &BaseDialog{
		Width:       width,
		Height:      height,
		Title:       title,
		Label:       label,
		Description: description,
	}
}

// Show runs the base dialog event loop and returns whether the dialog was
// confirmed, canceled, and any error that occurred.
func (b *BaseDialog) Show() (confirmed bool, canceled bool, err error) {
	w := app.Window{}
	w.Option(
		app.Title(b.Title),
		app.Size(unit.Dp(b.Width), unit.Dp(b.Height)),
	)
	// TODO work around https://todo.sr.ht/~eliasnaur/gio/602 (still an issue in gio v0.8.0?)
	// this should only be required shortly after creating the window w.
	// It doesn't work with the current gio version (0.8.1-dev), which only includes a fix for os_windows.
	applyWindowOptions := sync.OnceFunc(func() {
		time.Sleep(50 * time.Millisecond)
		w.Perform(system.ActionCenter | system.ActionRaise)
	})
	w.Perform(system.ActionCenter | system.ActionRaise)

	th := material.NewTheme()
	var ops op.Ops

	for !b.done {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			applyWindowOptions()
			gtx := app.NewContext(&ops, e)
			if b.cancelButton.Clicked(gtx) {
				b.handleCancel()
				w.Perform(system.ActionClose)
			}
			if b.okButton.Clicked(gtx) {
				b.handleOK()
				w.Perform(system.ActionClose)
			}
			b.layout(gtx, th)
			e.Frame(gtx.Ops)
		case key.Event:
			// Handle Escape key for cancel
			if e.Name == key.NameEscape && e.State == key.Press {
				b.handleCancel()
				w.Perform(system.ActionClose)
			}
			// Handle Enter key for OK
			if e.Name == key.NameReturn && e.State == key.Press {
				b.handleOK()
				w.Perform(system.ActionClose)
			}
		case app.DestroyEvent:
			b.done = true
			return b.confirmed, b.canceled, e.Err
		}
	}
	//app.Main()
	return b.confirmed, b.canceled, nil
}

func (b *BaseDialog) layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceAround}.Layout(gtx,
			// Label
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(th, b.Label)
				return label.Layout(gtx)
			}),
			// Description
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if b.Description == "" {
					return layout.Dimensions{}
				}
				desc := material.Body1(th, b.Description)
				desc.Color = th.Fg
				return desc.Layout(gtx)
			}),
			// Buttons
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &b.cancelButton, "Cancel")
						return btn.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &b.okButton, "OK")
						return btn.Layout(gtx)
					}),
				)
			}),
		)
	})
}

func (b *BaseDialog) handleOK() {
	b.confirmed = true
	b.canceled = false
}

func (b *BaseDialog) handleCancel() {
	b.confirmed = false
	b.canceled = true
}
