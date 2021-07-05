package main

import (
	"errors"
	"grpc-upload/app"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "CLI tool written in Go that allows a user to upload an image file and returns a URL",
	RunE:  rootRunner,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a file name argument")
		}
		if len(args) > 1 {
			return errors.New("only file name argument is required")
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().String("grpcPort", "127.0.0.1:8899", "grpc port")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func rootRunner(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	grpcPort, err := flags.GetString("grpcPort")
	if err != nil {
		return err
	}
	a, err := app.New(grpcPort)
	if err != nil {
		return err
	}

	url, err := a.UploadImage(args[0])
	if err != nil {
		return err
	}
	log.Printf(url)
	return nil
}
