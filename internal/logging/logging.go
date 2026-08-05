package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Setup redirects the standard logger to write to both stdout and a
// daily log file, and returns that same combined writer so other
// loggers (like chi's request logger, which does not use the standard
// log package) can be pointed at it explicitly.
func Setup(logDir string) (io.Writer, *os.File, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, nil, err
	}

	filename := time.Now().Format("2006-01-02") + ".log"
	path := filepath.Join(logDir, filename)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	writer := io.MultiWriter(os.Stdout, file)

	log.SetOutput(writer)
	log.SetFlags(log.Ldate | log.Ltime)

	return writer, file, nil
}
