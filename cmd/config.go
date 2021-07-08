package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "commands for managaing reporting configuration",
}

type Config struct {
	Networks map[string]*NetworkDetails `json:"networks" yaml:"networks" mapstructure:"networks"`
}
