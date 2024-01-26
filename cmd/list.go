/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/flevin58/icns/icns"
	"github.com/flevin58/xlog"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var icons []string
		if !verbose {
			xlog.Disable()
		}
		xlog.SetFlags(0)
		start, err := cmd.Flags().GetString("application")
		if err != nil {
			xlog.Fatalf("error getting flag 'a': %v\n", err)
		}
		if start != "" {
			icons, err = icns.FindInApp(start)
			if err != nil {
				xlog.Fatalf("error searching app: %v\n", err)
			}
		} else {
			if len(args) == 0 {
				xlog.Fatalln("missing folder name")
			}
			icons, err = icns.FindInFolder(args[0])
			if err != nil {
				xlog.Fatalf("error searching app: %v\n", err)
			}
		}
		if len(icons) == 0 {
			xlog.Infoln("no icons found")
			os.Exit(0)
		}

		for _, iconPath := range icons {
			ic, err := icns.NewIconFromFile(iconPath)
			if err != nil {
				continue
			}
			list := false
			for _, data := range ic.Data {
				if data.IconData[0] == 0x89 && string(data.IconData[1:4]) == "PNG" {
					isize := icns.GetIconSize(data.IconType.String())
					minsize, err := cmd.Flags().GetInt("minsize")
					if err != nil || isize < minsize {
						continue
					}
					list = true
				}
			}
			if list {
				fmt.Println(iconPath)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("application", "a", "", "Lists the icons from the given app")
	listCmd.Flags().IntP("minsize", "s", 128, "Smallest size to be listed")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
