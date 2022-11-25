/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weatherme",
	Short: "A basic cli weather app",
	Long:  `Type in weatherme and the name of a city to find out the weather in that city.`,

	Run: func(cmd *cobra.Command, args []string) {
		flagVar, err := cmd.Flags().GetBool("differentmessage")
		if err != nil {
			fmt.Println(err)
		}
		if flagVar {
			fmt.Println("This is a different message")
			return
		}
		fmt.Println("No flag addded")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("differentmessage", "d", false, "Toggle a different message")
}
