/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flevin58/icns/icns"
	"github.com/flevin58/xlog"
	"github.com/spf13/cobra"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extracts the icons in the given file or folder",
	Long: `Extracts the icons in the given file or folder.
	The folder should be an absolute path or relative to the current working dir.
	If the -a flag is given, the app name will be searched in /Applications and user /Applications
	examples:
	extract -a Twitter
	extract /Applications`,
	// Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !verbose {
			xlog.Disable()
		}
		xlog.SetFlags(0)

		var icons []string
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

		// Create output folder and extract all icons there
		outputdir, _ := cmd.Flags().GetString("outputdir")
		os.MkdirAll(outputdir, 0775)
		for _, iconPath := range icons {
			ic, err := icns.NewIconFromFile(iconPath)
			if err != nil {
				continue
			}
			for i, data := range ic.Data {
				if data.IconData[0] == 0x89 && string(data.IconData[1:4]) == "PNG" {
					isize := icns.GetIconSize(data.IconType.String())
					minsize, err := cmd.Flags().GetInt("minsize")
					if err != nil || isize < minsize {
						continue
					}
					normpath := icns.NormalizeIconName(iconPath)
					fname := filepath.Join(outputdir, fmt.Sprintf("%s_icon%d (%dx%d).png", normpath, i+1, isize, isize))
					f, err := os.Create(fname)
					if err != nil {
						continue
					}
					defer f.Close()
					f.Write(data.IconData)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.Flags().StringP("application", "a", "", "Extracts the icon from the given app")
	extractCmd.Flags().StringP("outputdir", "d", "./extracted", "Extracts the icon(s) to the given folder")
	extractCmd.Flags().IntP("minsize", "s", 128, "Smallest size to be extracted")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// extractCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// extractCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
