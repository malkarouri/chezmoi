package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var onepasswordCmd = &cobra.Command{
	Use:     "onepassword [args...]",
	Short:   "Execute the 1Password CLI (op)",
	PreRunE: config.ensureNoError,
	RunE:    config.runOnepasswordCmd,
}

type onepasswordCmdConfig struct {
	Command string
}

var onepasswordCache = make(map[string]interface{})

func init() {
	config.Onepassword.Command = "op"
	config.addTemplateFunc("onepassword", config.onepasswordFunc)

	secretCmd.AddCommand(onepasswordCmd)
}

func (c *Config) runOnepasswordCmd(cmd *cobra.Command, args []string) error {
	return c.run("", c.Onepassword.Command, args...)
}

func (c *Config) onepasswordFunc(item string) interface{} {
	if data, ok := onepasswordCache[item]; ok {
		return data
	}
	name := c.Onepassword.Command
	args := []string{"get", "item", item}
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	output, err := c.mutator.IdempotentCmdOutput(cmd)
	if err != nil {
		panic(fmt.Errorf("onepassword: %s %s: %w\n%s", name, strings.Join(args, " "), err, output))
	}
	var data interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		panic(fmt.Errorf("onepassword: %s %s: %w\n%s", name, strings.Join(args, " "), err, output))
	}
	onepasswordCache[item] = data
	return data
}