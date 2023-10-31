package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojotia/tia"
)

const (
	logLevelJSON = "json"
	logLevelText = "text"

	flagLogLevel          = "log-level"
	flagLogFormat         = "log-format"
	flagSkipProviderCheck = "skip-provider-check"
)

var rootCmd = &cobra.Command{
	Use:   "ojotia [auth-token] [celestia-rpc-addr] [ojo-grpc-addr]",
	Args:  cobra.ExactArgs(3),
	Short: "ojotia is a side-car process which takes information from ojo's price data and posts it to celestia.",
	RunE:  relayCmdHandler,
}

func init() {
	// We need to set our bech32 address prefix because it was moved
	// out of ojo's init function.
	// Ref: https://github.com/ojo-network/ojo/pull/63
	params.SetAddressPrefixes()
	rootCmd.PersistentFlags().String(flagLogLevel, zerolog.InfoLevel.String(), "logging level")
	rootCmd.PersistentFlags().String(flagLogFormat, logLevelText, "logging format; must be either json or text")
	rootCmd.PersistentFlags().Bool(flagSkipProviderCheck, false, "skip the coingecko API provider check")

	rootCmd.AddCommand(getQueryCmd())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func relayCmdHandler(cmd *cobra.Command, args []string) error {
	return tia.Submit(args[0], args[1], args[2], cmd.Context())
}
