package window

import (
	"fmt"

	"github.com/minodisk/go-nvim/buffer"
	cnvim "github.com/neovim/go-client/nvim"
)

type Window struct {
	cNvim   *cnvim.Nvim
	cWindow cnvim.Window
}

func New(v *cnvim.Nvim, w cnvim.Window) *Window {
	return &Window{v, w}
}

func (w *Window) Valid() (bool, error) {
	return w.cNvim.IsWindowValid(w.cWindow)
}

func (w *Window) Close() error {
	ok, err := w.Valid()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

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

func (w *Window) SetWidth(width int) error {
	if width <= 0 {
		return nil
	}
	return w.cNvim.SetWindowWidth(w.cWindow, width)
}

func (w *Window) SetHeight(height int) error {
	if height <= 0 {
		return nil
	}
	return w.cNvim.SetWindowHeight(w.cWindow, height)
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
	return w.cNvim.Command(fmt.Sprintf("edit `='%s'`", name))
}

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
