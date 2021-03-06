package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NoArgs checks that the command is not passed any parameters
func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf("unexpected argument: %s\nSee '%s --help'", args[0], cmd.CommandPath())
}

// NoArgsCustom returns a function which checks that the command is not passed any parameters
func NoArgsCustom(errstr string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}
		return fmt.Errorf(errstr, args[0])
	}
}

// ExactArgs returns an error if the exact number of args are not passed
func ExactArgs(num int) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == num {
			return nil
		}
		return fmt.Errorf(
			"\"%s\" requires exactly %d argument(s).\nSee '%s --help'",
			cmd.CommandPath(),
			num,
			cmd.CommandPath(),
		)
	}
}
