package op

import (
	"io"
	"os"
	"runtime"

	"bytes"

	"github.com/rs/zerolog"
)

var _ zerolog.LevelWriter = &ConsoleWriter{}

type Brush func([]byte) []byte

func NewBrush(color string) Brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text []byte) []byte {
		var buf bytes.Buffer
		buf.WriteString(pre)
		buf.WriteString(color)
		buf.WriteString("m")
		buf.Write(text)
		buf.WriteString(reset)
		return buf.Bytes()
	}
}

var colors = []Brush{
	// NewBrush("1;36"), // Trace      cyan
	NewBrush("1;34"), // Debug      blue
	NewBrush("1;32"), // Info       green
	NewBrush("1;33"), // Warn       yellow
	NewBrush("1;31"), // Error      red
	NewBrush("1;35"), // Critical   purple
	NewBrush("1;31"), // Fatal      red
}

// ConsoleWriter implements LoggerInterface and writes messages to terminal.
type ConsoleWriter struct {
	w     io.Writer
	Level zerolog.Level
}

// create ConsoleWriter returning as LoggerInterface.
func NewConsole(l zerolog.Level) zerolog.LevelWriter {
	return &ConsoleWriter{
		w:     os.Stdout,
		Level: l,
	}
}

func (cw *ConsoleWriter) Write(p []byte) (n int, err error) {
	return cw.w.Write(p)
}

func (cw *ConsoleWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	if cw.Level > l {
		return len(p), nil
	}
	if runtime.GOOS == "windows" {
		return cw.Write(p)

	}

	_, err = cw.Write(colors[l](p))
	if err != nil {
		return 0, err
	}

	return len(p), nil

}
