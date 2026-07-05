package enjoin

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Enjoin defines the interface for the enjoin subcommands.
type Enjoin interface {
	ApplyCmd() *cobra.Command
	ValidateCmd() *cobra.Command
}

type enjoinImpl struct {
	applyCmd    *cobra.Command
	validateCmd *cobra.Command
}

var (
	osOpen = os.Open
	osExit = os.Exit
)

// RunApply is the callback type for apply.
// Returns: success bool, details of changes made.
type RunApply func(*io.Reader) (bool, map[string]any)

// RunValidate is the callback type for dry-run validation.
// Returns: valid bool, details of what would change.
type RunValidate func(*io.Reader) (bool, map[string]any)

// attachFlags adds entity input flags to a command.
func attachFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(
		"entity",
		"e",
		"",
		"Specify the entity file path (reads from stdin if omitted)",
	)
	cmd.Flags().StringP("output", "o", "table", "Output format: table, json, or yaml")
}

// inputEntity reads the entity from the --entity flag or stdin.
func inputEntity(cmd *cobra.Command) (*io.Reader, error) {
	entityPath, err := cmd.Flags().GetString("entity")
	if err != nil {
		return nil, err
	}

	var entity io.Reader

	if entityPath == "" || entityPath == "-" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to stat stdin: %v", err)
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			entity = os.Stdin
		} else {
			return nil, fmt.Errorf("no entity file path or stdin input provided")
		}
	} else {
		entityFile, err := osOpen(entityPath)
		if err != nil {
			if os.IsNotExist(err) {
				cmd.PrintErrln("Entity file does not exist")
			} else {
				cmd.PrintErrf("Failed to open entity file: %v\n", err)
			}
			return nil, err
		}
		entity = entityFile
	}

	return &entity, nil
}

// apply runs the apply callback and sets exit code per UNIX protocol.
func (s *enjoinImpl) apply(runApply RunApply) {
	input, err := inputEntity(s.applyCmd)
	if err != nil {
		s.applyCmd.PrintErrln(err)
		osExit(2)
		return
	}

	success, details := runApply(input)

	result := map[string]any{
		"status":  "failed",
		"details": details,
	}
	if success {
		result["status"] = "applied"
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(result)

	if !success {
		osExit(1)
	}
}

// validate runs the validate callback (dry-run).
func (s *enjoinImpl) validate(runValidate RunValidate) {
	input, err := inputEntity(s.validateCmd)
	if err != nil {
		s.validateCmd.PrintErrln(err)
		osExit(2)
		return
	}

	valid, details := runValidate(input)

	result := map[string]any{
		"status":  "invalid",
		"details": details,
	}
	if valid {
		result["status"] = "valid"
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(result)

	if !valid {
		osExit(1)
	}
}

// ApplyCmd returns the apply cobra command.
func (s *enjoinImpl) ApplyCmd() *cobra.Command {
	return s.applyCmd
}

// ValidateCmd returns the validate cobra command.
func (s *enjoinImpl) ValidateCmd() *cobra.Command {
	return s.validateCmd
}
