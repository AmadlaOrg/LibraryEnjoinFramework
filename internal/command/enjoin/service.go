package enjoin

import "github.com/spf13/cobra"

// New creates an Enjoin instance with apply and validate subcommands.
func New(runApply RunApply, runValidate RunValidate) Enjoin {
	e := &enjoinImpl{
		applyCmd: &cobra.Command{
			Use:   "apply",
			Short: "Apply system state configuration from an entity",
		},
		validateCmd: &cobra.Command{
			Use:   "validate",
			Short: "Validate entity configuration without making changes",
		},
	}

	attachFlags(e.applyCmd)
	attachFlags(e.validateCmd)

	e.applyCmd.Run = func(cmd *cobra.Command, args []string) {
		e.apply(runApply)
	}

	e.validateCmd.Run = func(cmd *cobra.Command, args []string) {
		e.validate(runValidate)
	}

	return e
}
