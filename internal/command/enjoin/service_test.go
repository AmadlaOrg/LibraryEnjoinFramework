package enjoin

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	runApply := func(r *io.Reader) (bool, map[string]any) {
		return true, nil
	}
	runValidate := func(r *io.Reader) (bool, map[string]any) {
		return true, nil
	}

	svc := New(runApply, runValidate)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.ApplyCmd())
	assert.Equal(t, "apply", svc.ApplyCmd().Use)
	assert.NotNil(t, svc.ValidateCmd())
	assert.Equal(t, "validate", svc.ValidateCmd().Use)
}
