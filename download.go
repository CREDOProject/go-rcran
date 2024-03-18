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
	lib         string
}

const defaultLibrary = " .libPaths()[1L]"

func Install(options *InstallOptions) (string, error) {
	const install = `install.packages(
	"%s", # package name
	lib = "%s", #
)`
	if options.lib == "" {
		options.lib = defaultLibrary
	}
	_, err := os.Stat(options.packageName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(install, options.packageName, options.lib), nil
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
