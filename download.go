package gorcran

import (
	"fmt"
	"os"

	"github.com/CREDOProject/sharedutils/files"
)

const defaultMirror = "http://cran.us.r-project.org"

type DownloadOptions struct {
	InstallOptions
	destinationDirectory string
	repository           string
}

func Download(options *DownloadOptions) (string, error) {
	const download = `download.packages(
	"%s", # package name
	"%s", # destination directory
	repos = "%s", # repository
)`
	if options.repository == "" {
		options.repository = defaultMirror
	}
	if !files.IsDir(options.destinationDirectory) {
		return "", fmt.Errorf("Destination directory does not exist.")
	}
	return fmt.Sprintf(
		download,
		options.packageName,
		options.destinationDirectory,
		options.repository,
	), nil
}

type InstallOptions struct {
	packageName string
}

func Install(options *InstallOptions) (string, error) {
	const install = `install.packages(
	"%s" # package name
)`
	_, err := os.Stat(options.packageName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(install, options.packageName), nil
}
