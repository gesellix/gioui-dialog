package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gesellix/gioui-dialog/pkg/dialog"
)

// Gio-based demo UI showcasing the dialog API usage on desktop platforms.
func main() {
	// Run the Gio application in a separate goroutine and exit on close.
	go func() {
		w := new(app.Window)
		if err := eventLoop(w); err != nil {
			if err != nil {
				log.Println("window error:", err)
			}
		}
		os.Exit(0)
	}()
	app.Main()
}

func eventLoop(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops

	// Clickable buttons and result display state.
	var (
		inputBtn   widget.Clickable
		selectBtn  widget.Clickable
		baseBtn    widget.Clickable
		resultText string
	)

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Handle button clicks per frame.
			if inputBtn.Clicked(gtx) {
				go func() {
					res, canceled, err := dialog.PromptInput(dialog.InputDialogOptions{
						Title:       "Text Input",
						Label:       "Enter something",
						Description: "Please enter text.",
						DefaultText: "Default",
					})
					if err != nil {
						log.Println("Error showing input dialog:", err)
						resultText = "Error"
					} else if canceled {
						resultText = "Canceled"
					} else {
						resultText = fmt.Sprintf("Input: %q", res)
					}
				}()
			}

			if selectBtn.Clicked(gtx) {
				go func() {
					choices, canceled, err := dialog.PromptSelect(dialog.SelectDialogOptions{
						Title:             "Multiple Selection",
						Label:             "Choose options",
						Choices:           []string{"Option A", "Option B", "Option C"},
						DefaultSelections: []string{"Option B"},
						AllowCustomEntry:  true,
					})
					if err != nil {
						log.Println("Error showing select dialog:", err)
						resultText = "Error"
					} else if canceled {
						resultText = "Canceled"
					} else {
						resultText = fmt.Sprintf("Selection: %v", choices)
					}
				}()
			}

			if baseBtn.Clicked(gtx) {
				go func() {
					confirmed, canceled, err := dialog.PromptBase(dialog.BaseDialogOptions{
						Title:       "Base Dialog",
						Label:       "This is a base dialog",
						Description: "Shows only title, label, description and OK/Cancel buttons.",
					})
					if err != nil {
						log.Println("Error showing base dialog:", err)
						resultText = "Error"
					} else if canceled {
						resultText = "Base dialog canceled"
					} else if confirmed {
						resultText = "Base dialog confirmed"
					}
				}()
			}

			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceAround}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &inputBtn, "Text Dialog").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &selectBtn, "Select Dialog").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &baseBtn, "Base Dialog").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Body1(th, resultText).Layout(gtx)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
	return nil
}
