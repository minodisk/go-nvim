package nvim

import (
	"fmt"
	"strings"

	"github.com/minodisk/go-nvim/window"
	cnvim "github.com/neovim/go-client/nvim"
)

type WindowDirection string
type WindowPosition string

const (
	WindowHorizontal  WindowDirection = "horizontal"
	WindowVertical    WindowDirection = "vertical"
	WindowTopLeft     WindowPosition  = "topleft"
	WindowBottomRight WindowPosition  = "botright"
)

type Nvim struct {
	cNvim *cnvim.Nvim
}

func New(v *cnvim.Nvim) *Nvim {
	return &Nvim{v}
}

func (v *Nvim) VarBool(name string) (bool, error) {
	var b bool
	if err := v.cNvim.Var(name, &b); err != nil {
		return false, err
	}
	return b, nil
}

func (v *Nvim) VarInt(name string) (int, error) {
	var i int
	if err := v.cNvim.Var(name, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (v *Nvim) VarString(name string) (string, error) {
	var s string
	if err := v.cNvim.Var(name, &s); err != nil {
		return "", err
	}
	return s, nil
}

func (v *Nvim) CreateWindowLeft(width int, name string) (*window.Window, error) {
	return v.CreateWindow(WindowVertical, WindowTopLeft, width, name)
}

func (v *Nvim) CreateWindowRight(height int, name string) (*window.Window, error) {
	return v.CreateWindow(WindowVertical, WindowBottomRight, height, name)
}

func (v *Nvim) CreateWindow(d WindowDirection, p WindowPosition, size int, name string) (*window.Window, error) {
	c, err := v.cNvim.CurrentWindow()
	if err != nil {
		return nil, err
	}
	defer v.cNvim.SetCurrentWindow(c)

	var s string
	if size > 0 {
		s = fmt.Sprintf("%d", size)
	}
	if err := v.cNvim.Command(fmt.Sprintf("%s %s %ssplit %s", d, p, s, name)); err != nil {
		return nil, err
	}
	w, err := v.CurrentWindow()
	if err != nil {
		return nil, err
	}
	switch d {
	case WindowVertical:
		w.SetDefaultWidth(size)
	case WindowHorizontal:
		w.SetDefaultHeight(size)
	}
	return w, nil
}

func (v *Nvim) CurrentWindow() (*window.Window, error) {
	w, err := v.cNvim.CurrentWindow()
	if err != nil {
		return nil, err
	}
	return window.New(v.cNvim, w), nil
}

func (v *Nvim) CurrentBufferName() (string, error) {
	buf, err := v.cNvim.CurrentBuffer()
	if err != nil {
		return "", err
	}
	return v.cNvim.BufferName(buf)
}

// func (v *Nvim) Buffer(name string) (*buffer.Buffer, error) {
// 	bufs, err := v.cNvim.Buffers()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, buf := range bufs {
// 		b := buffer.New(v.cNvim, buf)
// 		n, err := b.Name()
// 		if err != nil {
// 			continue
// 		}
// 		if n == name {
// 			return b, nil
// 		}
// 	}
// 	return nil, errors.New("not found")
// }

func (v *Nvim) Windows() ([]*window.Window, error) {
	ws, err := v.cNvim.Windows()
	if err != nil {
		return nil, err
	}
	windows := make([]*window.Window, len(ws))
	for i, w := range ws {
		windows[i] = window.New(v.cNvim, w)
	}
	return windows, nil
}

func (v *Nvim) InputString(prompt string) (string, error) {
	return v.Input(prompt, "")
}

func (v *Nvim) InputStrings(prompt string) ([]string, error) {
	out, err := v.Input(fmt.Sprintf("%s, separated by commas", prompt), "")
	if err != nil {
		return nil, err
	}
	ss := strings.Split(out, ",")
	for i, s := range ss {
		ss[i] = strings.TrimSpace(s)
	}
	return ss, err
}

func (v *Nvim) InputBool(prompt string) (bool, error) {
	out, err := v.Input(fmt.Sprintf("%s [y/n]", prompt), "")
	if err != nil {
		return false, err
	}
	return strings.ToLower(out) == "y", nil
}

func (v *Nvim) Input(prompt, defaultText string) (string, error) {
	out, err := v.cNvim.CommandOutput(fmt.Sprintf(`echo input("%s: ", "%s")`, prompt, defaultText))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (v *Nvim) Printf(format string, args ...interface{}) error {
	return v.cNvim.WriteOut(fmt.Sprintf(format, args...))
}

func (v *Nvim) PrintError(err error) error {
	if err == nil {
		return nil
	}
	return v.cNvim.WriteErr(err.Error() + "\n")
}
