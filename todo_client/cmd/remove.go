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

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:          "remove <item id>",
	Short:        "Delete item by id",
	Aliases:      []string{"d"},
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	Long:         ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiUrl := viper.GetString("api-url")
		return removeAction(os.Stdout, apiUrl, args)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func removeAction(w io.Writer, url string, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("%w, Arg must be a number", ErrNotNumber)
	}
	if err := deleteItem(url, id); err != nil {
		return err
	}
	return printDelete(w, id)
}

func printDelete(w io.Writer, id int) error {
	_, err := fmt.Fprintf(w, "Item No %d removed", id)
	return err
}
