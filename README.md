# nvim [![GoDoc](https://godoc.org/github.com/minodisk/go-nvim?status.png)](https://godoc.org/github.com/minodisk/go-nvim)

Yet another wrapper for github.com/neovim/go-client.

## Installation

```
go get github.com/minodisk/go-nvim
```

## Usage

```go
import (
	"github.com/minodisk/go-nvim/buffer"
	"github.com/minodisk/go-nvim/nvim"
	"github.com/minodisk/go-nvim/window"
	cnvim "github.com/neovim/go-client/nvim"
	cplugin "github.com/neovim/go-client/nvim/plugin"
)

func main() {
	cplugin.Main(func (p *cplugin.Plugin) error {
		p.HandleCommand(&cplugin.CommandOptions{
			Name:  "Example",
		}, func (v *cnvim.Nvim) error {
			vim := nvim.New(v)
			win, _ := vim.CurrentWindow()
			buf, _ := win.Buffer()
			buf.Write([][]byte{
				[]byte("foo"),
				[]byte("bar"),
				[]byte("baz"),
			})
		})
		return nil
	}
}
```
