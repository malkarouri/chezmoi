package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/twpayne/chezmoi/next/internal/chezmoi"
)

type persistentStateData struct {
	ScriptOnce interface{} `json:"scriptOnce" toml:"scriptOnce" yaml:"scriptOnce"`
}

func (c *Config) newStateCmd() *cobra.Command {
	stateCmd := &cobra.Command{
		Use:   "state",
		Short: "Manipulate the state",
		// Long: mustGetLongHelp("state"), // FIXME
		Example: getExample("state"), // FIXME
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create the state if it does not already exist",
		// Long: mustGetLongHelp("state", "create"), // FIXME
		// Example: getExample("state", "create"), // FIXME
		Args: cobra.NoArgs,
		RunE: c.runStateCreateCmd,
	}
	stateCmd.AddCommand(createCmd)

	dumpCmd := &cobra.Command{
		Use:   "dump",
		Short: "Generate a dump of the state",
		// Long: mustGetLongHelp("state", "dump"), // FIXME
		// Example: getExample("state", "dump"), // FIXME
		Args: cobra.NoArgs,
		RunE: c.runStateDataCmd,
	}
	stateCmd.AddCommand(dumpCmd)

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset the state",
		// Long: mustGetLongHelp("state", "reset"), // FIXME
		// Example: getExample("state", "reset"), // FIXME
		Args: cobra.NoArgs,
		RunE: c.runStateResetCmd,
	}
	stateCmd.AddCommand(resetCmd)

	return stateCmd
}

func (c *Config) runStateCreateCmd(cmd *cobra.Command, args []string) error {
	return c.baseSystem.PersistentState().OpenOrCreate()
}

func (c *Config) runStateDataCmd(cmd *cobra.Command, args []string) error {
	scriptOnceData, err := chezmoi.ScriptOnceData(c.baseSystem.PersistentState())
	if err != nil {
		return err
	}
	return c.marshal(&persistentStateData{
		ScriptOnce: scriptOnceData,
	})
}

func (c *Config) runStateResetCmd(cmd *cobra.Command, args []string) error {
	path := c.getPersistentStateFile()
	_, err := c.baseSystem.Stat(path)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	if !c.force {
		choice, err := c.prompt(fmt.Sprintf("Remove %s", path), "yn")
		if err != nil {
			return err
		}
		if choice == 'n' {
			return nil
		}
	}
	return c.baseSystem.RemoveAll(path)
}
