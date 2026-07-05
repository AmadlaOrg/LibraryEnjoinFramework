package display

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func newTestCmd() (*cobra.Command, *bytes.Buffer) {
	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	return cmd, buf
}

func TestDisplay_Table(t *testing.T) {
	cmd, _ := newTestCmd()
	var buf bytes.Buffer
	data := map[string]any{
		"system": map[string]string{
			"os": "linux",
		},
	}
	svc := NewWithWriter(cmd, data, &buf)

	err := svc.Display()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "CATEGORY")
}

func TestDisplay_JSON(t *testing.T) {
	cmd, _ := newTestCmd()
	var buf bytes.Buffer
	data := map[string]any{
		"status": "pass",
	}
	svc := NewWithWriter(cmd, data, &buf)
	cmd.Flags().Set("output", "json")

	err := svc.Display()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"status": "pass"`)
}

func TestDisplay_YAML(t *testing.T) {
	cmd, _ := newTestCmd()
	var buf bytes.Buffer
	data := map[string]any{
		"status": "pass",
	}
	svc := NewWithWriter(cmd, data, &buf)
	cmd.Flags().Set("output", "yaml")

	err := svc.Display()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "status: pass")
}

func TestDisplay_UnknownFormat(t *testing.T) {
	cmd, _ := newTestCmd()
	var buf bytes.Buffer
	data := map[string]any{}
	svc := NewWithWriter(cmd, data, &buf)
	cmd.Flags().Set("output", "xml")

	err := svc.Display()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown output format")
}

func TestDisplay_JSONMarshalError(t *testing.T) {
	originalMarshal := jsonMarshalIndent
	defer func() { jsonMarshalIndent = originalMarshal }()

	jsonMarshalIndent = func(v any, prefix, indent string) ([]byte, error) {
		return nil, errors.New("marshal failed")
	}

	cmd, _ := newTestCmd()
	var buf bytes.Buffer
	data := map[string]any{"key": "value"}
	svc := NewWithWriter(cmd, data, &buf)
	cmd.Flags().Set("output", "json")

	err := svc.Display()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error encoding JSON")
}

func TestDisplay_YAMLMarshalError(t *testing.T) {
	originalMarshal := yamlMarshal
	defer func() { yamlMarshal = originalMarshal }()

	yamlMarshal = func(v any) ([]byte, error) {
		return nil, errors.New("marshal failed")
	}

	cmd, _ := newTestCmd()
	var buf bytes.Buffer
	data := map[string]any{"key": "value"}
	svc := NewWithWriter(cmd, data, &buf)
	cmd.Flags().Set("output", "yaml")

	err := svc.Display()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error encoding YAML")
}

func TestSetTableHeaders(t *testing.T) {
	cmd, _ := newTestCmd()
	data := map[string]any{}
	svc := New(cmd, data)

	result := svc.SetTableHeaders([]string{"A", "B"})
	assert.Equal(t, svc, result)
}

func TestGetTableHeaders_Default(t *testing.T) {
	s := &displayImpl{}
	headers := s.getTableHeaders()
	assert.Equal(t, []string{"Category", "Key", "Value"}, headers)
}

func TestGetTableHeaders_Custom(t *testing.T) {
	s := &displayImpl{tableHeaders: []string{"X", "Y"}}
	headers := s.getTableHeaders()
	assert.Equal(t, []string{"X", "Y"}, headers)
}
