package display

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// New creates a Display for the given command and data, writing to stdout.
func New(cmd *cobra.Command, data map[string]any) Display {
	display := &displayImpl{
		cmd:    cmd,
		data:   data,
		writer: os.Stdout,
	}

	display.attachFlags()

	return display
}

// NewWithWriter creates a Display writing to a custom writer.
func NewWithWriter(cmd *cobra.Command, data map[string]any, writer io.Writer) Display {
	display := &displayImpl{
		cmd:    cmd,
		data:   data,
		writer: writer,
	}

	display.attachFlags()

	return display
}
