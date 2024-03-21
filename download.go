package gorcran

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/CREDOProject/sharedutils/files"
)

const defaultMirror = "http://cran.us.r-project.org"

func GetDependencies(o *InstallOptions) (string, error) {
	const retrieve = `
	r <- getOption("repos")
	r["CRAN"] <- "%s"
	options(repos=r)
	pkgs = utils:::getDependencies(
		pkgs = "%s",
		available = available.packages(),
		lib = "%s",
	)
	cat(pkgs, sep="\n")
`

	if o.Repository == "" {
		o.Repository = defaultMirror
	}
	if o.Lib == "" {
		o.Lib = defaultLibrary
	}
	if o.PackageName == "" {
		return "", fmt.Errorf("Package name not specified")
	}

	return fmt.Sprintf(retrieve,
		o.Repository,
		o.PackageName,
		o.Lib,
	), nil
}

type DownloadOptions struct {
	PackageName          string
	DestinationDirectory string
	Repository           string
}

func Download(options *DownloadOptions) (string, error) {
	const download = `download.packages(
	pkgs    = "%s", # package name
	destdir = "%s", # destination directory
	repos   = "%s", # repository
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
	Repository  string
	DryRun      bool
}

const defaultLibrary = " .libPaths()[1L]"

func Install(options *InstallOptions) (string, error) {
	const install = `install.packages(
	pkgs  = "%s", # package name
	lib   = "%s", # Library
	repos = "%s", # Repository
)`
	if options.Lib == "" {
		options.Lib = defaultLibrary
	}
	if options.Repository == "" {
		options.Repository = defaultMirror
	}
	if !options.DryRun {
		_, err := os.Stat(options.PackageName)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf(install,
			options.PackageName,
			options.Lib,
			options.Repository),
		nil
}

var inquotes = regexp.MustCompile("(?:\").+?(?:\")")

func ParsePath(output string) (string, error) {
	outputStrings := inquotes.FindAllString(output, -1)
	if len(outputStrings) < 1 {
		return "", fmt.Errorf("Could not find any strings.")
	}
	getLast := outputStrings[len(outputStrings)-1]
	getLast = path.Clean(getLast)
	getLast = path.Base(getLast)
	getLast = strings.TrimPrefix(getLast, "\"")
	getLast = strings.TrimSuffix(getLast, "\"")
	return getLast, nil
}
