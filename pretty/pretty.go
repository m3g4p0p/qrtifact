package pretty

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/alecthomas/chroma/quick"
)

type Console struct {
	w io.Writer
}

func NewConsole(w io.Writer) *Console {
	return &Console{w: w}
}

func (c *Console) Pretty(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	defer fmt.Fprintln(c.w)
	return quick.Highlight(
		c.w,
		string(data),
		"json",
		"terminal256",
		"fruity",
	)
}

func (c *Console) Write(b []byte) (int, error) {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return len(b), err
	}
	return len(b), c.Pretty(v)
}
