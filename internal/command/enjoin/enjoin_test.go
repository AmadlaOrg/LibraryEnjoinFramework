package enjoin

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestApply_Success(t *testing.T) {
	originalExit := osExit
	defer func() { osExit = originalExit }()
	var exitCode int
	osExit = func(code int) { exitCode = code }

	runApply := func(r *io.Reader) (bool, map[string]any) {
		return true, map[string]any{"changed": "user created"}
	}

	s := &enjoinImpl{applyCmd: &cobra.Command{Use: "apply"}}
	attachFlags(s.applyCmd)

	tmpFile, err := os.CreateTemp(t.TempDir(), "entity")
	assert.NoError(t, err)
	_, _ = tmpFile.WriteString("test entity content")
	tmpFile.Close()

	originalOpen := osOpen
	defer func() { osOpen = originalOpen }()
	osOpen = func(name string) (*os.File, error) {
		return os.Open(tmpFile.Name())
	}

	s.applyCmd.Flags().Set("entity", tmpFile.Name())
	s.apply(runApply)

	assert.Equal(t, 0, exitCode)
}

func TestApply_Fail(t *testing.T) {
	originalExit := osExit
	defer func() { osExit = originalExit }()
	var exitCode int
	osExit = func(code int) { exitCode = code }

	runApply := func(r *io.Reader) (bool, map[string]any) {
		return false, map[string]any{"error": "permission denied"}
	}

	s := &enjoinImpl{applyCmd: &cobra.Command{Use: "apply"}}
	attachFlags(s.applyCmd)

	tmpFile, err := os.CreateTemp(t.TempDir(), "entity")
	assert.NoError(t, err)
	tmpFile.Close()

	originalOpen := osOpen
	defer func() { osOpen = originalOpen }()
	osOpen = func(name string) (*os.File, error) {
		return os.Open(tmpFile.Name())
	}

	s.applyCmd.Flags().Set("entity", tmpFile.Name())
	s.apply(runApply)

	assert.Equal(t, 1, exitCode)
}

func TestValidate_Success(t *testing.T) {
	originalExit := osExit
	defer func() { osExit = originalExit }()
	var exitCode int
	osExit = func(code int) { exitCode = code }

	runValidate := func(r *io.Reader) (bool, map[string]any) {
		return true, map[string]any{"would_change": "nothing"}
	}

	s := &enjoinImpl{validateCmd: &cobra.Command{Use: "validate"}}
	attachFlags(s.validateCmd)

	tmpFile, err := os.CreateTemp(t.TempDir(), "entity")
	assert.NoError(t, err)
	_, _ = tmpFile.WriteString("test entity content")
	tmpFile.Close()

	originalOpen := osOpen
	defer func() { osOpen = originalOpen }()
	osOpen = func(name string) (*os.File, error) {
		return os.Open(tmpFile.Name())
	}

	s.validateCmd.Flags().Set("entity", tmpFile.Name())
	s.validate(runValidate)

	assert.Equal(t, 0, exitCode)
}

func TestValidate_Fail(t *testing.T) {
	originalExit := osExit
	defer func() { osExit = originalExit }()
	var exitCode int
	osExit = func(code int) { exitCode = code }

	runValidate := func(r *io.Reader) (bool, map[string]any) {
		return false, map[string]any{"error": "invalid config"}
	}

	s := &enjoinImpl{validateCmd: &cobra.Command{Use: "validate"}}
	attachFlags(s.validateCmd)

	tmpFile, err := os.CreateTemp(t.TempDir(), "entity")
	assert.NoError(t, err)
	tmpFile.Close()

	originalOpen := osOpen
	defer func() { osOpen = originalOpen }()
	osOpen = func(name string) (*os.File, error) {
		return os.Open(tmpFile.Name())
	}

	s.validateCmd.Flags().Set("entity", tmpFile.Name())
	s.validate(runValidate)

	assert.Equal(t, 1, exitCode)
}

func TestInputEntity_FileNotFound(t *testing.T) {
	originalOpen := osOpen
	defer func() { osOpen = originalOpen }()
	osOpen = func(name string) (*os.File, error) {
		return nil, os.ErrNotExist
	}

	cmd := &cobra.Command{Use: "apply"}
	attachFlags(cmd)
	cmd.Flags().Set("entity", "/nonexistent/path")

	reader, err := inputEntity(cmd)
	assert.Error(t, err)
	assert.Nil(t, reader)
}

func TestInputEntity_FileSuccess(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "entity")
	assert.NoError(t, err)
	_, _ = tmpFile.WriteString("test content")
	tmpFile.Close()

	cmd := &cobra.Command{Use: "apply"}
	attachFlags(cmd)
	cmd.Flags().Set("entity", tmpFile.Name())

	reader, err := inputEntity(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, reader)

	data, err := io.ReadAll(*reader)
	assert.NoError(t, err)
	assert.Equal(t, "test content", string(data))

	if closer, ok := (*reader).(io.Closer); ok {
		closer.Close()
	}
}

func TestInputEntity_EmptyPath_NoStdin(t *testing.T) {
	cmd := &cobra.Command{Use: "apply"}
	attachFlags(cmd)

	reader, err := inputEntity(cmd)
	if err != nil {
		assert.Nil(t, reader)
	}
}

func TestInputEntity_Dash_NoStdin(t *testing.T) {
	cmd := &cobra.Command{Use: "apply"}
	attachFlags(cmd)
	cmd.Flags().Set("entity", "-")

	reader, err := inputEntity(cmd)
	if err != nil {
		assert.Nil(t, reader)
	}
}

func TestAttachFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "apply"}
	attachFlags(cmd)

	flag := cmd.Flags().Lookup("entity")
	assert.NotNil(t, flag)
	assert.Equal(t, "e", flag.Shorthand)

	outputFlag := cmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "o", outputFlag.Shorthand)
	assert.Equal(t, "table", outputFlag.DefValue)
}

func TestInputEntity_ReaderNotClosed(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "entity")
	assert.NoError(t, err)
	content := "entity data for reading"
	_, _ = tmpFile.WriteString(content)
	tmpFile.Close()

	cmd := &cobra.Command{Use: "apply"}
	attachFlags(cmd)
	cmd.Flags().Set("entity", tmpFile.Name())

	reader, err := inputEntity(cmd)
	assert.NoError(t, err)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, *reader)
	assert.NoError(t, err)
	assert.Equal(t, content, buf.String())

	if closer, ok := (*reader).(io.Closer); ok {
		closer.Close()
	}
}

func TestApplyCmd(t *testing.T) {
	cmd := &cobra.Command{Use: "apply"}
	s := &enjoinImpl{applyCmd: cmd}
	assert.Equal(t, cmd, s.ApplyCmd())
}

func TestValidateCmd(t *testing.T) {
	cmd := &cobra.Command{Use: "validate"}
	s := &enjoinImpl{validateCmd: cmd}
	assert.Equal(t, cmd, s.ValidateCmd())
}
