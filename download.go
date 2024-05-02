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

func GetBioconductorDepenencies(o *InstallOptions) (string, error) {
	const retrieve = `
	require("BiocManager")
	r <- getOption("repos")
	r <- BiocManager::repositories()
	r["CRAN"] <- "%s"
	options(repos=r)
	pkgs = tools::package_dependencies(
		packages = "%s",
		recursive = TRUE,
	)$%s
	cat(pkgs, sep="\n")
`
	return _getDependencies(retrieve, o)
}

func _getDependencies(template string, o *InstallOptions) (string, error) {
	if o.Repository == "" {
		o.Repository = defaultMirror
	}
	if o.PackageName == "" {
		return "", fmt.Errorf("Package name not specified")
	}
	return fmt.Sprintf(template,
		o.Repository,
		o.PackageName,
		o.PackageName,
	), nil
}

func GetDependencies(o *InstallOptions) (string, error) {
	const retrieve = `
	r <- getOption("repos")
	r["CRAN"] <- "%s"
	options(repos=r)
	pkgs = tools::package_dependencies(
		packages = "%s",
		recursive = TRUE,
	)$%s
	cat(pkgs, sep="\n")
`
	return _getDependencies(retrieve, o)
}

type DownloadOptions struct {
	PackageName          string
	DestinationDirectory string
	Repository           string
}

func DownloadBioconductor(options *DownloadOptions) (string, error) {
	const download = `
	require("BiocManager")
	r <- getOption("repos")
	r <- BiocManager::repositories()
	r["CRAN"] <- "%s"
	options(repos=r)
	withCallingHandlers(
		download.packages(
			repos   = r, # repository
			pkgs    = "%s", # package name
			destdir = "%s", # destination directory
		),
		warning = function(w) quit(status=1)
	)`
	return _download(download, options)
}

func _download(template string, options *DownloadOptions) (string, error) {
	if options.Repository == "" {
		options.Repository = defaultMirror
	}
	if !files.IsDir(options.DestinationDirectory) {
		return "", fmt.Errorf("Destination directory does not exist.")
	}
	return fmt.Sprintf(
		template,
		options.Repository,
		options.PackageName,
		options.DestinationDirectory,
	), nil
}

func Download(options *DownloadOptions) (string, error) {
	const download = `
	withCallingHandlers(
		download.packages(
			repos   = "%s", # repository
			pkgs    = "%s", # package name
			destdir = "%s", # destination directory
		),
		warning = function(w) quit(status=1)
	)`

	return _download(download, options)
}

type InstallOptions struct {
	PackageName string
	Lib         string
	Repository  string
	DryRun      bool
}

func _install(template string, o *InstallOptions) (string, error) {
	if o.Lib == "" {
		o.Lib = defaultLibrary
	}
	if o.Repository == "" {
		o.Repository = defaultMirror
	}
	if !o.DryRun {
		_, err := os.Stat(o.PackageName)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf(template,
			o.Repository,
			o.PackageName,
			o.Lib,
		),
		nil
}

func InstallBioconductor(o *InstallOptions) (string, error) {
	const install = `
	require("BiocManager")
	r <- getOption("repos")
	r <- BiocManager::repositories()
	r["CRAN"] <- "%s"
	withCallingHandlers(
		install.packages(
			pkgs  = "%s",
			lib   = "%s",
			repos = r,
		),
		warning = function(w) quit(status=1),
	)`

	return _install(install, o)
}

const defaultLibrary = ".libPaths()[1L]"

func Install(options *InstallOptions) (string, error) {
	const install = `install.packages(
		repos = "%s", # Repository
		pkgs  = "%s", # package name
		lib   = "%s", # Library
	)`
	return _install(install, options)
}

func InstallLocal(o *InstallOptions) (string, error) {
	const install = `install.packages(
		repos = NULL, # Repository
		pkgs  = "%s", # package name
		lib   = "%s", # Library
	)`
	if o.Lib == "" {
		o.Lib = defaultLibrary
	}
	if !o.DryRun {
		_, err := os.Stat(o.PackageName)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf(install,
		o.PackageName,
		o.Lib,
	), nil
}

var inquotes = regexp.MustCompile("(?:\").+?(?:\")")

func ParsePath(output string) (string, error) {
	getLast, err := GetPath(output)
	if err != nil {
		return "", err
	}
	getLast = path.Base(getLast)
	return getLast, nil
}

func GetPath(output string) (string, error) {
	outputStrings := inquotes.FindAllString(output, -1)
	if len(outputStrings) < 1 {
		return "", fmt.Errorf("Could not find any strings.")
	}
	filepath := outputStrings[len(outputStrings)-1]
	filepath = path.Clean(filepath)
	filepath = strings.TrimPrefix(filepath, "\"")
	filepath = strings.TrimSuffix(filepath, "\"")
	return filepath, nil
}
