package main

import (
	"log"
	"upload-server/app"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "CLI tool written in Go that allows a user to upload an image file and returns a URL",
	RunE:  rootRunner,
}

func init() {
	rootCmd.PersistentFlags().String("httpPort", ":8888", "http target")
	rootCmd.PersistentFlags().String("grpcPort", ":8899", "grpc target")
	rootCmd.PersistentFlags().String("filePath", "/app/data", "file folder path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func rootRunner(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	httpPort, err := flags.GetString("httpPort")
	if err != nil {
		return err
	}
	grpcPort, err := flags.GetString("grpcPort")
	if err != nil {
		return err
	}
	a, err := app.New(grpcPort, httpPort)
	if err != nil {
		return err
	}
	return a.Run()
}
