package icns

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func FindInFolder(folder string) ([]string, error) {
	var icons []string
	filepath.WalkDir(folder, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !file.IsDir() {
			if strings.ToLower(filepath.Ext(path)) == ".icns" {
				icons = append(icons, path)
			}
		}
		return nil
	})

	if _, err := os.Stat(folder); err != nil {
		return icons, err
	}
	return icons, nil
}

func FindInApp(app string) (icons []string, err error) {
	if strings.ToLower(filepath.Ext(app)) != ".app" {
		app += ".app"
	}
	appfolder := path.Join("/Applications", app)
	icons, err = FindInFolder(appfolder)
	if err != nil {
		return icons, err
	}

	userapps, err := os.UserHomeDir()
	if err != nil {
		return icons, nil
	}

	appfolder = path.Join(userapps, app)
	usericons, err := FindInFolder(appfolder)
	if err != nil {
		icons = append(icons, usericons...)
	}
	return icons, nil
}
