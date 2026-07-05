package display

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	data := map[string]any{"key": "value"}

	svc := New(cmd, data)

	assert.NotNil(t, svc)
}

func TestNewWithWriter(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	data := map[string]any{"key": "value"}
	var buf bytes.Buffer

	svc := NewWithWriter(cmd, data, &buf)

	assert.NotNil(t, svc)
}
