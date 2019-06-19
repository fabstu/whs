package httpmanual

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var addr = "localhost:8081"

var CMD = &cobra.Command{
	Use:   "http-manual",
	Short: "starts a http server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		fmt.Println("Accepting on", addr)
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("failed to accept connection:", err)
				continue
			}
			go handleConnection(conn)
		}
	},
}

func init() {

}