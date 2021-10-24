package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var app = &cobra.Command{
	Use:   "level-go",
	Short: "level is something alright",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := app.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
