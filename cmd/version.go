package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	PROGRAM_NAME    = "miningopt"
	PROGRAM_TITLE   = "Ultimate Pit Optimization"
	PROGRAM_VERSION = "v6.0.0"
	COPYRIGHT       = "Copyright (C) 2019 Robert Wright"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the program version",
	Long:  "Output the program version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", PROGRAM_NAME, PROGRAM_VERSION)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
