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

func (b *Buffer) Valid() (bool, error) {
	return b.cNvim.IsBufferValid(b.cBuffer)
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

func (b *Buffer) Focused() (bool, error) {
	cb, err := b.cNvim.CurrentBuffer()
	if err != nil {
		return false, err
	}
	return b.cBuffer == cb, nil
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

type Range struct {
	Start, End Position
}

type Position struct {
	BufferNumber int
	LineNumber   int
	Column       int
	Offset       int
}

var (
	rPosition = regexp.MustCompile(`\[(\d+), (\d+), (\d+), (\d+)\]`)
)

func NewPosition(s string) (Position, error) {
	var (
		p   Position
		err error
	)
	nums := rPosition.FindStringSubmatch(s)
	p.BufferNumber, err = strconv.Atoi(nums[1])
	if err != nil {
		return p, err
	}
	p.LineNumber, err = strconv.Atoi(nums[2])
	if err != nil {
		return p, err
	}
	p.Column, err = strconv.Atoi(nums[3])
	if err != nil {
		return p, err
	}
	p.Offset, err = strconv.Atoi(nums[4])
	if err != nil {
		return p, err
	}
	return p, nil
}

func (p Position) X() int {
	return p.Column - 1
}

func (p *Position) SetX(x int) {
	p.Column = x + 1
}

func (p Position) Y() int {
	return p.LineNumber - 1
}

func (p *Position) SetY(y int) {
	p.LineNumber = y + 1
}

func (p Position) String() string {
	return fmt.Sprintf("[%d, %d, %d, %d]", p.BufferNumber, p.LineNumber, p.Column, p.Offset)
}

func (b *Buffer) CurrentCursor() (Position, error) {
	return b.getpos(".")
}

func (b *Buffer) SetCurrentCursor(p Position) error {
	return b.setpos(".", p)
}

func (b *Buffer) SelectedRange() (Range, error) {
	r := Range{}
	var err error
	r.Start, err = b.getpos("'<")
	if err != nil {
		return r, err
	}
	r.End, err = b.getpos("'>")
	if err != nil {
		return r, err
	}
	return r, nil
}

func (b *Buffer) Mode() (string, error) {
	return b.CommandOutputf("silent echo mode()")
}

func (b *Buffer) getpos(expr string) (Position, error) {
	p, err := b.CommandOutputf(fmt.Sprintf("silent echo getpos(\"%s\")", expr))
	if err != nil {
		return Position{}, err
	}
	return NewPosition(p)
}

func (b *Buffer) setpos(expr string, p Position) error {
	return b.Commandf("call setpos(\"%s\", %s)", expr, p.String())
}

func (b *Buffer) FileType() (string, error) {
	var t string
	err := b.cNvim.BufferOption(b.cBuffer, "filetype", &t)
	return t, err
}

// SetFileType sets filetype t to buffer.
func (b *Buffer) SetFileType(t string) error {
	return b.Commandf("set filetype=%s", t)
}

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
