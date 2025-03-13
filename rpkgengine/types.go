package rpkgengine

// RpkgBuildFile is the struct for the rpkg.build.yaml file
type RpkgBuildFile struct {
	// Name: The name of the package (this should be the same as the name of the package folder) e.g. test
	Name string
	// Version: The version of the package (this should be in the format of major.minor.patch) e.g. 1.0.0
	Version string
	// Revision: The revision of the package (default is 0 for a new major, minor or patch version) e.g. 0fa5d3
	Revision int
	// Authors: The authors of the package (this should be a list of authors) e.g. ["John Doe", "Jane Doe"]
	Authors []any
	// Deps: The dependencies of the package (this should be a list of dependencies) e.g. ["requests@latest", "flask@1.1.2"]
	Deps []any
	// BuildDeps: The build dependencies of the package (this should be a list of build dependencies) e.g. ["pytest@latest", "flake8@3.9.2"]
	BuildDeps []any
	// BuildWith: The language the package is built with (this should be the language the package is built with with a version right after it) e.g. python3.13
	BuildWith string
	// BuildCommands: The build commands of the package (this should be a list of build commands that build your package) e.g. ["python3.13 setup.py sdist"]
	BuildCommands []any
}
