package enjoin

import (
	"github.com/AmadlaOrg/LibraryEnjoinFramework/internal/command"
	internalenjoin "github.com/AmadlaOrg/LibraryEnjoinFramework/internal/command/enjoin"
	"github.com/AmadlaOrg/LibraryFramework/cli"
	"github.com/spf13/cobra"
)

// RunApply is the callback type that enjoin plugin implementors must satisfy for apply.
type RunApply = internalenjoin.RunApply

// RunValidate is the callback type for dry-run validation (no changes made).
type RunValidate = internalenjoin.RunValidate

// New sets up the enjoin CLI application with UNIX plugin protocol support.
func New(
	name, title, version string,
	supports []string,
	description string,
	runApply RunApply,
	runValidate RunValidate) {

	meta := command.PluginMeta{
		Name:        name,
		Version:     version,
		Supports:    supports,
		Description: description,
	}

	cli.New(name, title, version, func(c *cobra.Command) {
		command.New(c, meta, runApply, runValidate)
	})
}
