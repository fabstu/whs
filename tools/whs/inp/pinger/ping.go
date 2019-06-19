package pinger

import (
	"aduu.dev/tools/whs/inp/pinger/domain"
	"aduu.dev/tools/whs/inp/pinger/domain/xicmp"

	"fmt"

	"github.com/spf13/cobra"
)

var pingCMD = &cobra.Command{
	Use:   "ping",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {

		//localAddresses()

		/*
			settings := domain.PingSettings{
				TargetAddress: strings.TrimSpace(args[0]),
				Count:         *count,
				SourceAddress: *sourceAddress,
			}
		*/
		for i, settings := range domain.DefaultSettings() {
			fmt.Println("\nsettings i:", i)
			fmt.Printf("%#v\n", settings)
			var implementations = map[string]domain.Pinger{
				"xicmp": xicmp.Make(settings),
			}

			for name, impl := range implementations {
				fmt.Println("---", name, "---")
				res, err := impl.Ping()
				if err != nil {
					fmt.Println("ping error:", err)
					continue
				}

				fmt.Printf("peer=%s time=%s\n", res.Peer.String(), res.TimeTaken.String())
			}
		}
	},
}

func init() {

}