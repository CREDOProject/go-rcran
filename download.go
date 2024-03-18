package gorcran

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/CREDOProject/sharedutils/files"
)

const defaultMirror = "http://cran.us.r-project.org"

type DownloadOptions struct {
	InstallOptions
	DestinationDirectory string
	Repository           string
}

func Download(options *DownloadOptions) (string, error) {
	const download = `download.packages(
	"%s", # package name
	"%s", # destination directory
	repos = "%s", # repository
)`
	if options.Repository == "" {
		options.Repository = defaultMirror
	}
	if !files.IsDir(options.DestinationDirectory) {
		return "", fmt.Errorf("Destination directory does not exist.")
	}
	return fmt.Sprintf(
		download,
		options.PackageName,
		options.DestinationDirectory,
		options.Repository,
	), nil
}

type InstallOptions struct {
	PackageName string
	Lib         string
}

const defaultLibrary = " .libPaths()[1L]"

func Install(options *InstallOptions) (string, error) {
	const install = `install.packages(
	"%s", # package name
	lib = "%s", #
)`
	if options.Lib == "" {
		options.Lib = defaultLibrary
	}
	_, err := os.Stat(options.PackageName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(install, options.PackageName, options.Lib), nil
}

var inquotes = regexp.MustCompile("\"(.*?)\"")

func ParsePath(output string) (string, error) {
	strings := inquotes.FindAllString(output, -1)
	if len(strings) < 1 {
		return "", fmt.Errorf("Could not find any strings.")
	}
	getLast := strings[len(strings)-1]
	getLast = path.Clean(getLast)
	getLast = path.Base(getLast)
	return getLast, nil
}
