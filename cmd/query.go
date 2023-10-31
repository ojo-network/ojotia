package cmd

import (
	"github.com/ojo-network/ojotia/tia"
	"github.com/spf13/cobra"
)

const (
	flagFormat = "format"
)

func getQueryCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "query [auth-token] [celestia-rpc-addr] [commitment] [height]",
		Args:  cobra.ExactArgs(4),
		Short: "Query ojo blob data on celestia",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tia.Query(args[0], args[1], args[2], args[3])
		},
	}

	return versionCmd
}
