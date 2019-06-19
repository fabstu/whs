package main

import (
	"aduu.dev/lib/completion"
	"aduu.dev/tools/aduu/helper"
	"fmt"
	"aduu.dev/tools/whs/inp"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "whs",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		helper.CheckErr(cmd.Help())
	},
}

func init() {
	RootCmd.AddCommand(completion.NewCompletionCMD())
	RootCmd.AddCommand(inp.CMD)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
