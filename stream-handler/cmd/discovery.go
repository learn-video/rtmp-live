package cmd

import (
	"log"

	"github.com/learn-video/rtmp-live/discovery"
	"github.com/spf13/cobra"
)

func RunDiscovery() *cobra.Command {
	return &cobra.Command{
		Use:   "discovery",
		Short: "Discover running streams",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := discovery.LoadConfig()
			if err != nil {
				log.Fatal(err)
			}
			discovery.Watch(cfg)
		},
	}
}
