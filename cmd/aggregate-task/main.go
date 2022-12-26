package main

import (
	"aggregate-task/internal/busi"
	"context"

	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// @title spacescope aggregate task
// @version 1.0
// @description spacescope aggregate task
// @termsOfService http://swagger.io/terms/

// @contact.name xueyouchen
// @contact.email xueyou@starboardventures.io

// @host aggregate-api.spacescope.io
// @BasePath /

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-task",
		Short: "at",
		Run: func(cmd *cobra.Command, args []string) {
			if err := entry(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.PersistentFlags().StringVar(&busi.Flags.Config, "conf", "", "path of the configuration file")

	return cmd
}

func entry() error {
	busi.NewServer(context.Background()).Start()
	return nil
}

func main() {
	if err := NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
