package cmd

import (
	"github.com/learn-video/rtmp-live/api"
	"github.com/spf13/cobra"
)

func RunAuthServer() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Start authorizer server",
		Run: func(cmd *cobra.Command, args []string) {
			api.RunServer()
		},
	}
}
