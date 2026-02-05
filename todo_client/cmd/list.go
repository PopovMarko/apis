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
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:          "list",
	Aliases:      []string{"l"},
	Short:        "Show all items of the Todo List ",
	SilenceUsage: true,
	Long:         ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiUrl := viper.GetString("api-url")
		return listAction(os.Stdout, apiUrl)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listAction(out io.Writer, url string) error {
	items, err := getAll(url)
	if err != nil {
		return err
	}
	printAll(out, items)
	return nil
}

func printAll(out io.Writer, items []item) error {
	w := tabwriter.NewWriter(out, 4, 2, 0, ' ', 0)
	for k, v := range items {
		done := "-"
		if v.Done {
			done = "X"
		}
		fmt.Fprintf(w, "%s\t%d\t%s\t\n", done, k+1, v.Task)
	}
	return w.Flush()
}
