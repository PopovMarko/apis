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
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:          "complete <item No>",
	Short:        "Set the item as completed",
	Aliases:      []string{"c"},
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	Long:         ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiUrl := viper.GetString("api-url")
		return completeAction(os.Stdout, apiUrl, args)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}

func completeAction(w io.Writer, url string, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("%w Argument must be a number.", ErrNotNumber)
	}
	if err := completeItem(url, id); err != nil {
		return err
	}
	return printComplete(w, id)
}

func printComplete(w io.Writer, id int) error {
	_, err := fmt.Fprintf(w, "Item No %d set as completed", id)
	return err
}
