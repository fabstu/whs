package pinger

import (
	"aduu.dev/tools/whs/inp/pinger/domain"
	"aduu.dev/tools/whs/inp/pinger/domain/xicmp"

	"fmt"

	"github.com/spf13/cobra"
)

var tracerouteCMD = &cobra.Command{
	Use:   "traceroute",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		settings := domain.DefaultSettings()[1]

		fmt.Printf("settings: %#v\n", settings)

		fmt.Println("---", "xicmp", "---")
		impl := xicmp.Make(settings)

		fmt.Println(fmt.Sprintf("-%d", 3))

		return impl.Traceroute()
	},
}

func init() {

}
