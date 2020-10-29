package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	PROGRAM_NAME    = "whattle"
	PROGRAM_TITLE   = "Ultimate Pit Optimization"
	PROGRAM_VERSION = "v6.1.0"
	COPYRIGHT       = "Copyright (C) 2020 Robert Wright"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the program version",
	Long:  "Output the program version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", PROGRAM_NAME, PROGRAM_VERSION, COPYRIGHT)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
