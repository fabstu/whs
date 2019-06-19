package inp

import (
	"github.com/spf13/cobra"

	"aduu.dev/tools/whs/inp/dns"
	"aduu.dev/tools/whs/inp/httpmanual"
	"aduu.dev/tools/whs/inp/pinger"

	"aduu.dev/tools/aduu/helper"
)

var CMD = &cobra.Command{
	Use:   "inp",
	Short: "Collects programs programmed in Internet-Protokolle",
	Run: func(cmd *cobra.Command, args []string) {
		helper.CheckErr(cmd.Help())
	},
}

func init() {
	CMD.AddCommand(pinger.CMD, httpmanual.CMD, dns.CMD)
}
