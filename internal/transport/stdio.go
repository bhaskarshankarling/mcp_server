// Package transport provides communication transports for MCP
package transport

import (
	"bufio"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// StdioTransport implements MCP over stdio
type StdioTransport struct {
	reader *bufio.Reader
	writer io.Writer
	logger *logrus.Logger
}

// NewStdioTransport creates a new stdio transport
func NewStdioTransport() *StdioTransport {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &StdioTransport{
		reader: bufio.NewReader(os.Stdin),
		writer: os.Stdout,
		logger: logger,
	}
}

// ReadMessage reads a message from stdin
func (t *StdioTransport) ReadMessage() ([]byte, error) {
	line, err := t.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// Remove the newline character
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}

	t.logger.Debugf("Received: %s", string(line))
	return line, nil
}

// WriteMessage writes a message to stdout
func (t *StdioTransport) WriteMessage(data []byte) error {
	_, err := t.writer.Write(append(data, '\n'))
	if err != nil {
		return err
	}

	t.logger.Debugf("Sent: %s", string(data))
	return nil
}

// SetLogLevel sets the logging level
func (t *StdioTransport) SetLogLevel(level logrus.Level) {
	t.logger.SetLevel(level)
}
