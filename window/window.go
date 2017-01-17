package window

import (
	"fmt"

	"github.com/minodisk/go-nvim/buffer"
	cnvim "github.com/neovim/go-client/nvim"
)

type Window struct {
	cNvim   *cnvim.Nvim
	cWindow cnvim.Window
	width   int
	height  int
}

func New(v *cnvim.Nvim, w cnvim.Window) *Window {
	return &Window{v, w, 0, 0}
}

func (w *Window) Close() error {
	c, err := w.cNvim.CurrentWindow()
	if err != nil {
		return err
	}
	defer w.cNvim.SetCurrentWindow(c)

	if err := w.Focus(); err != nil {
		return err
	}
	return w.cNvim.Command("quit")
}

func (w *Window) SetDefaultWidth(width int) {
	w.width = width
}

func (w *Window) SetDefaultHeight(height int) {
	w.height = height
}

func (w *Window) ResizeToDefaultWidth() error {
	if w.width <= 0 {
		return nil
	}
	return w.cNvim.SetWindowWidth(w.cWindow, w.width)
}

func (w *Window) ResizeToDefaultHeight() error {
	if w.height <= 0 {
		return nil
	}
	return w.cNvim.SetWindowHeight(w.cWindow, w.height)
}

func (w *Window) Buffer() (*buffer.Buffer, error) {
	b, err := w.cNvim.WindowBuffer(w.cWindow)
	if err != nil {
		return nil, err
	}
	return buffer.New(w.cNvim, b), nil
}

func (w *Window) Focus() error {
	return w.cNvim.SetCurrentWindow(w.cWindow)
}

func (w *Window) Open(name string) error {
	c, err := w.cNvim.CurrentWindow()
	if err != nil {
		return err
	}
	defer w.cNvim.SetCurrentWindow(c)

	if err := w.Focus(); err != nil {
		return err
	}
	// n := Escape(name)
	// return w.cNvim.WriteOut(fmt.Sprintf("edit '%s'", n))
	return w.cNvim.Command(fmt.Sprintf("edit `='%s'`", name))
}

// func (w *Window) Option(name string) (interface{}, error) {
// 	var v interface{}
// 	if err := w.cNvim.WindowOption(w.cWindow, name, &v); err != nil {
// 		return nil, err
// 	}
// 	return v, nil
// }
//
// func (w *Window) SetOption(name string, value interface{}) error {
// 	return w.cNvim.SetWindowOption(w.cWindow, name, value)
// }

func (w *Window) Option() (Option, error) {
	var o Option
	for k, p := range o.MapPointer() {
		if err := w.cNvim.WindowOption(w.cWindow, k, p); err != nil {
			return o, err
		}
	}
	return o, nil
}

func (w *Window) SetOption(o Option) error {
	for n, v := range o.MapValue() {
		if err := w.cNvim.SetWindowOption(w.cWindow, n, v); err != nil {
			return err
		}
	}
	return nil
}
