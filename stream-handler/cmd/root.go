package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "stream_handler",
		Short:         "Handle incoming RTMP streamings",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	rootCmd.AddCommand(RunAuthServer())
	return rootCmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
