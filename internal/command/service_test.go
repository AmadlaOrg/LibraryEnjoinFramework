package command

import (
	"bytes"
	"io"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	meta := PluginMeta{
		Name:        "enjoin-test",
		Version:     "1.0.0",
		Supports:    []string{"amadla.org/entity/test@^v1.0.0"},
		Description: "Test plugin",
	}
	runApply := func(r *io.Reader) (bool, map[string]any) {
		return true, nil
	}
	runValidate := func(r *io.Reader) (bool, map[string]any) {
		return true, nil
	}

	New(cmd, meta, runApply, runValidate)

	var names []string
	for _, sub := range cmd.Commands() {
		names = append(names, sub.Use)
	}
	assert.Contains(t, names, "apply")
	assert.Contains(t, names, "validate")
	assert.Contains(t, names, "info")
}

func TestInfoSubcommand(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	meta := PluginMeta{
		Name:        "enjoin-test",
		Version:     "1.0.0",
		Supports:    []string{"amadla.org/entity/test@^v1.0.0"},
		Description: "Test plugin",
	}
	runApply := func(r *io.Reader) (bool, map[string]any) {
		return true, nil
	}
	runValidate := func(r *io.Reader) (bool, map[string]any) {
		return true, nil
	}

	New(cmd, meta, runApply, runValidate)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"info"})
	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"name":"enjoin-test"`)
	assert.Contains(t, buf.String(), `"version":"1.0.0"`)
	assert.Contains(t, buf.String(), `"supports"`)
}
