/*
Copyright © 2026 The Pragmatic Programmers LLC
Copyright apply to this codebase.
Check license for detailes.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <task name>",
	Aliases:      []string{"a"},
	Short:        "Add task to list",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	Long:         ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiUrl := viper.GetString("api-url")
		return addAction(os.Stdout, apiUrl, args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addAction(w io.Writer, url string, args []string) error {
	task := strings.Join(args, " ")
	if err := addItem(url, task); err != nil {
		return err
	}
	return printAdd(w, task)
}

func printAdd(w io.Writer, task string) error {
	_, err := fmt.Fprintf(w, "Task: %s, added to the list\n", task)
	return err
}
