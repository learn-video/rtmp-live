package cmd

import (
	"github.com/learn-video/rtmp-live/auth"
	"github.com/spf13/cobra"
)

func RunAuthServer() *cobra.Command {
	return &cobra.Command{
		Use:   "auth_server",
		Short: "Start authorizer server",
		Run: func(cmd *cobra.Command, args []string) {
			auth.RunServer()
		},
	}
}
