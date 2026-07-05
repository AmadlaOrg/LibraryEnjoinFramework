package command

import (
	"encoding/json"
	"fmt"

	"github.com/AmadlaOrg/LibraryEnjoinFramework/internal/command/enjoin"
	"github.com/spf13/cobra"
)

// PluginMeta holds metadata for the enjoin plugin.
type PluginMeta struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Supports    []string `json:"supports"`
	Description string   `json:"description"`
}

// New wires up the info, apply, and validate subcommands.
func New(
	cmd *cobra.Command,
	meta PluginMeta,
	runApply enjoin.RunApply,
	runValidate enjoin.RunValidate) {

	enjoinService := enjoin.New(runApply, runValidate)
	cmd.AddCommand(enjoinService.ApplyCmd())
	cmd.AddCommand(enjoinService.ValidateCmd())

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Display plugin metadata",
		Run: func(c *cobra.Command, args []string) {
			data, _ := json.Marshal(meta)
			fmt.Fprintln(c.OutOrStdout(), string(data))
		},
	}
	cmd.AddCommand(infoCmd)
}
