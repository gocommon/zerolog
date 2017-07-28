package op

import (
	"github.com/gocommon/rotatefile"
	"github.com/gocommon/zerolog"
)

var _ zerolog.LevelWriter = &FileLogWriter{}

// NewFileWriter NewFileWriter
func NewFileWriter(w *rotatefile.Writer, lv zerolog.Level) zerolog.LevelWriter {
	return &FileLogWriter{
		Level: lv,
		fd:    w,
	}
}

// FileLogWriter implements LevelWriter.
type FileLogWriter struct {
	Level zerolog.Level `json:"level"`
	fd    *rotatefile.Writer
}

// WriteLevel implements io.Writer.
func (w *FileLogWriter) Write(p []byte) (n int, err error) {
	return w.fd.Write(p)
}

// WriteLevel implements LevelWriter.
func (w *FileLogWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	if l < w.Level {
		return len(p), nil
	}
	return w.Write(p)

}
