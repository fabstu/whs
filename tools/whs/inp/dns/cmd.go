package dns

import (
	"github.com/spf13/cobra"

	"aduu.dev/tools/aduu/helper"
)

var CMD = &cobra.Command{
	Use:   "dns",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		helper.CheckErr(cmd.Help())
	},
}

func init() {
	CMD.AddCommand(resolveCMD)
}
