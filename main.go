package main

import (
	"fmt"
	"os"

	"github.com/abiosoft/lineprefix"
	"github.com/spf13/cobra"
)

func newDemoCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "demo",
		Short: "SHORT TEXT",
		Long:  "LONG TEXT",
	}
	c.AddCommand(newSubcommand())
	c.AddCommand(&cobra.Command{
		Use:        "hullo",
		Deprecated: "use hello instead",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("hullo wurld")
		},
	})
	return &c
}

func newSubcommand() *cobra.Command {
	c := cobra.Command{
		Use:   "hello",
		Short: "SHORT TEXT",
		Long:  "LONG TEXT",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("hello world")
		},
	}
	c.Flags().Int("num", 0, "number")
	c.Flags().Bool("flag", false, "description")
	_ = c.Flags().MarkDeprecated("flag", "please don't use it")
	return &c
}

func main() {
	c := newDemoCommand()

	c.SetArgs(os.Args[1:])

	c.SetOut(lineprefix.New(
		lineprefix.Writer(os.Stdout),
		lineprefix.Prefix("\033[34m[stdout]\033[m"),
	))
	c.SetErr(lineprefix.New(
		lineprefix.Writer(os.Stdout),
		lineprefix.Prefix("\033[31m[stderr]\033[m"),
	))

	if _, err := c.ExecuteC(); err != nil {
		fmt.Fprintf(os.Stderr, "Execute() error (%T): %s\n", err, err.Error())
		os.Exit(1)
	}
}
