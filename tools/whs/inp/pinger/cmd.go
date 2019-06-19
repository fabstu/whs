package pinger

import (
	"github.com/spf13/cobra"

	"aduu.dev/tools/aduu/helper"

	"aduu.dev/lib/completion"
)

var CMD = &cobra.Command{
	Use:   "pinger",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		helper.CheckErr(cmd.Help())
	},
	//Args: cobra.ExactArgs(1),
}

var (
	count         *int
	sourceAddress *string
)

func init() {
	CMD.AddCommand(completion.NewCompletionCMD())
	CMD.AddCommand(pingCMD, tracerouteCMD)

	count = CMD.Flags().IntP("count", "c", -1, "sets the amount of pings to send")
	sourceAddress = CMD.Flags().String("source", "::", "the address to bind to")
}
