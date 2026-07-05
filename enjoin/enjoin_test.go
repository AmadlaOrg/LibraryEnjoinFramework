package enjoin

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"enjoin-test", "version"}

	assert.NotPanics(t, func() {
		New("enjoin-test", "Enjoin Test", "1.0.0",
			[]string{"amadla.org/entity/test@^v1.0.0"},
			"Test enjoin plugin",
			func(r *io.Reader) (bool, map[string]any) {
				return true, nil
			},
			func(r *io.Reader) (bool, map[string]any) {
				return true, nil
			})
	})
}
