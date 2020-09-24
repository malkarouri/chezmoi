package cmd

// FIXME update REFERENCE.md
// FIXME add per-shell Long and Example

import (
	"io"
	"strings"

	"github.com/spf13/cobra"
)

func (c *Config) newCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion code",
		// Long:      mustGetLongHelp("completion"),
		// Example:   getExample("completion"),
	}

	makeRunE := func(genCompletionFunc func(io.Writer) error) func(*cobra.Command, []string) error {
		return func(cmd *cobra.Command, args []string) error {
			sb := &strings.Builder{}
			if err := genCompletionFunc(sb); err != nil {
				return err
			}
			return c.writeOutputString(sb.String())
		}
	}

	bashCmd := &cobra.Command{
		Use:   "bash",
		Args:  cobra.NoArgs,
		Short: "Generate bash completion code",
		RunE:  makeRunE(rootCmd.GenBashCompletion),
		Annotations: map[string]string{
			doesNotRequireValidConfig: "true",
		},
	}
	completionCmd.AddCommand(bashCmd)

	fishCmd := &cobra.Command{
		Use:   "fish",
		Args:  cobra.NoArgs,
		Short: "Generate fish completion code",
		RunE: makeRunE(func(w io.Writer) error {
			return rootCmd.GenFishCompletion(w, true)
		}),
		Annotations: map[string]string{
			doesNotRequireValidConfig: "true",
		},
	}
	completionCmd.AddCommand(fishCmd)

	powerShellCmd := &cobra.Command{
		Use:   "powershell",
		Args:  cobra.NoArgs,
		Short: "Generate PowerShell completion code",
		RunE:  makeRunE(rootCmd.GenPowerShellCompletion),
		Annotations: map[string]string{
			doesNotRequireValidConfig: "true",
		},
	}
	completionCmd.AddCommand(powerShellCmd)

	zshCmd := &cobra.Command{
		Use:   "zsh",
		Args:  cobra.NoArgs,
		Short: "Generate zsh completion code",
		RunE:  makeRunE(rootCmd.GenZshCompletion),
		Annotations: map[string]string{
			doesNotRequireValidConfig: "true",
		},
	}
	completionCmd.AddCommand(zshCmd)

	return completionCmd
}
