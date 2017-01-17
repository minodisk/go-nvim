package buffer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	cnvim "github.com/neovim/go-client/nvim"
)

type Buffer struct {
	cNvim   *cnvim.Nvim
	cBuffer cnvim.Buffer
}

func New(v *cnvim.Nvim, b cnvim.Buffer) *Buffer {
	return &Buffer{v, b}
}

func (b *Buffer) Name() (string, error) {
	return b.cNvim.BufferName(b.cBuffer)
}

func (b *Buffer) SetName(name string) error {
	return b.cNvim.SetBufferName(b.cBuffer, name)
}

func (b *Buffer) Focus() error {
	return b.cNvim.SetCurrentBuffer(b.cBuffer)
}

func (b *Buffer) Commandf(format string, args ...interface{}) error {
	c, err := b.cNvim.CurrentBuffer()
	if err != nil {
		return err
	}
	if err := b.Focus(); err != nil {
		return err
	}
	defer b.cNvim.SetCurrentBuffer(c)

	return b.cNvim.Command(fmt.Sprintf(format, args...))
}

func (b *Buffer) CommandOutputf(format string, args ...interface{}) (string, error) {
	c, err := b.cNvim.CurrentBuffer()
	if err != nil {
		return "", err
	}
	if err := b.Focus(); err != nil {
		return "", err
	}
	defer b.cNvim.SetCurrentBuffer(c)

	s, err := b.cNvim.CommandOutput(fmt.Sprintf(format, args...))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

type Position struct {
	bufnum int
	X      int
	Y      int
	off    int
}

var (
	rPosition = regexp.MustCompile(`\[(\d+), (\d+), (\d+), (\d+)\]`)
)

func NewPosition(s string) (Position, error) {
	nums := rPosition.FindStringSubmatch(s)
	bufnum, err := strconv.Atoi(nums[1])
	if err != nil {
		return Position{}, err
	}
	lnum, err := strconv.Atoi(nums[2])
	if err != nil {
		return Position{}, err
	}
	col, err := strconv.Atoi(nums[3])
	if err != nil {
		return Position{}, err
	}
	off, err := strconv.Atoi(nums[4])
	if err != nil {
		return Position{}, err
	}
	return Position{bufnum, lnum - 1, col - 1, off}, nil
}

func (p Position) String() string {
	return fmt.Sprintf("[%d, %d, %d, %d]", p.bufnum, p.X+1, p.Y+1, p.off)
}

func (b *Buffer) CurrentCursor() (Position, error) {
	c, err := b.cNvim.CurrentBuffer()
	if err != nil {
		return Position{}, err
	}
	if err := b.cNvim.SetCurrentBuffer(b.cBuffer); err != nil {
		return Position{}, err
	}
	defer b.cNvim.SetCurrentBuffer(c)

	p, err := b.CommandOutputf("silent echo getpos('.')")
	if err != nil {
		return Position{}, err
	}
	return NewPosition(p)
}

func (b *Buffer) SetCurrentCursor(p Position) error {
	c, err := b.cNvim.CurrentBuffer()
	if err != nil {
		return err
	}
	if err := b.cNvim.SetCurrentBuffer(b.cBuffer); err != nil {
		return err
	}
	defer b.cNvim.SetCurrentBuffer(c)

	return b.Commandf("call setpos('.', %s)", p.String())
}

func (b *Buffer) FileType() (string, error) {
	var t string
	err := b.cNvim.BufferOption(b.cBuffer, "filetype", &t)
	return t, err
}

// SetFileType sets filetype t to buffer.
//
// This method is WORKAROUND.
// Using github.com/nvim/go-client nvim.SetBufferOption to set filetype finder,
// FileType event of autocmd doesn't fire.
// Though using nvim.Command("set filetype=foo"), FileType event fired.
func (b *Buffer) SetFileType(t string) error {
	return b.Commandf("set filetype=%s", t)
}

// func (b *Buffer) Option(name string) (interface{}, error) {
// 	var res interface{}
// 	if err := b.cNvim.BufferOption(b.cBuffer, name, &res); err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }

// func (b *Buffer) SetOption(name string, value interface{}) error {
// 	return b.cNvim.SetBufferOption(b.cBuffer, name, value)
// }

func (b *Buffer) Option() (Option, error) {
	var o Option
	for k, p := range o.MapPointer() {
		if err := b.cNvim.BufferOption(b.cBuffer, k, p); err != nil {
			return o, err
		}
	}
	return o, nil
}

func (b *Buffer) SetOption(o Option) error {
	for n, v := range o.MapValue() {
		if err := b.cNvim.SetBufferOption(b.cBuffer, n, v); err != nil {
			return err
		}
	}
	return nil
}

func (b *Buffer) Write(lines [][]byte) error {
	m, err := b.Modifiable()
	if err != nil {
		return err
	}
	if err := b.SetModifiable(true); err != nil {
		return err
	}
	defer b.SetModifiable(m)

	p, err := b.CurrentCursor()
	if err != nil {
		return err
	}
	defer b.SetCurrentCursor(p)

	l, err := b.LineCount()
	if err != nil {
		return err
	}
	return b.cNvim.SetBufferLines(b.cBuffer, 0, l, true, lines)
}

func (b *Buffer) LineCount() (int, error) {
	return b.cNvim.BufferLineCount(b.cBuffer)
}

func (b *Buffer) Clear() error {
	l, err := b.LineCount()
	if err != nil {
		return err
	}
	return b.cNvim.SetBufferLines(b.cBuffer, 0, l-1, true, [][]byte{})
}

func (b *Buffer) Modifiable() (bool, error) {
	var m bool
	if err := b.cNvim.BufferOption(b.cBuffer, "modifiable", &m); err != nil {
		return false, err
	}
	return m, nil
}

func (b *Buffer) SetModifiable(m bool) error {
	return b.cNvim.SetBufferOption(b.cBuffer, "modifiable", m)
}
