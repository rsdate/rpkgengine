package rpkgengine

// Total lines in this file: 116
import (
	"fmt"
	"os/exec"
	"strings"
)

// RpkgBuildFile is the struct for the rpkg.build.yaml file
//
// Name: The name of the package (this should be the same as the name of the package folder) e.g. test
//
// Version: The version of the package (this should be in the format of major.minor.patch) e.g. 1.0.0
//
// Revision: The revision of the package (default is 0 for a new major, minor or patch version) e.g. 0fa5d3
//
// Authors: The authors of the package (this should be a list of authors) e.g. ["John Doe", "Jane Doe"]
//
// Deps: The dependencies of the package (this should be a list of dependencies) e.g. ["requests@latest", "flask@1.1.2"]
//
// BuildDeps: The build dependencies of the package (this should be a list of build dependencies) e.g. ["pytest@latest", "flake8@3.9.2"]
//
// BuildWith: The language the package is built with (this should be the language the package is built with with a version right after it) e.g. python3.13
//
// BuildCommands: The build commands of the package (this should be a list of build commands that build your package) e.g. ["python3.13 setup.py sdist"]
type RpkgBuildFile struct {
	Name          string
	Version       string
	Revision      int
	Authors       []interface{}
	Deps          []interface{}
	BuildDeps     []interface{}
	BuildWith     string
	BuildCommands []interface{}
}

// Description: Hello returns a string of a text art
//
// Parameters: None
//
// Returns: It returns a string of a text art
func Hello() string {
	text := `
\\  \\  \\  \\  \\                  .. .. .. .. .. ..  
\\                \\              ..                 ..
\\                \\              ..                 ..
\\                \\              ..                 ..
\\                \\              ..                 ..
\\  \\  \\  \\  \\                .. .. .. .. .. .. ..              
\\           \\                   ..
\\            \\                  ..
\\             \\                 ..
\\              \\                ..
\\               \\               ..
\\                \\              ..
	`
	return text
}

// Description: installDeps installs dependencies. It is a helper function for the Build function (also so that it wasn't too long)
//
// Parameters: It takes a list of dependencies and a boolean to check if the dependencies are build dependencies
//
// Returns: It returns an integer and an error. The integer is the exit code of the dependency installation process (1 or 0) and the error is any error that occurred during the dependency installation process
func installDeps(deps []interface{}, buildDeps bool) (int, error) {
	for i := 0; i < len(deps); i++ {
		if deps[0] == "none" {
			if buildDeps {
				fmt.Println("No build dependencies")
				break
			} else {
				fmt.Println("No dependencies")
				break
			}
		} else if dep, ok := deps[i].(string); ok && !strings.Contains(dep, "@") {
			if buildDeps {
				fmt.Printf("Build dependency %s is not a valid dependency\n", []any{dep}...)
				return 1, fmt.Errorf("build dependency %s is not a valid dependency", dep)
			} else {
				fmt.Printf("Dependency %s is not a valid dependency\n", []any{dep}...)
				return 1, fmt.Errorf("dependency %s is not a valid dependency", dep)
			}
		} else if dep, ok := deps[i].(string); ok && strings.Contains(dep, "@latest") {
			fmt.Printf("Installing %s... ", []any{deps[i]}...)
			dep := strings.Split(deps[i].(string), "@")[0]
			cmd := exec.Command("python3.13", "-m", "pip", "install", "--upgrade", dep)
			cmd.Stdout = nil
			if _, err := cmd.Output(); err != nil {
				fmt.Printf("Could not install %s\n", []any{dep}...)
				if buildDeps {
					return 1, fmt.Errorf("could not install build dependency %s", dep)
				} else {
					return 1, fmt.Errorf("could not install dependency %s", dep)
				}
			} else {
				fmt.Printf("Installed %s\n", []any{dep}...)
			}
		} else {
			dep := strings.Split(deps[i].(string), "@")
			cmd := exec.Command("python3.13", "-m", "pip", "install", dep[0], "==", dep[1])
			cmd.Stdout = nil
			if _, err := cmd.Output(); err != nil {
				fmt.Printf("Could not install %s\n", []any{dep}...)
				if buildDeps {
					return 1, fmt.Errorf("could not install build dependency %s", dep)
				} else {
					return 1, fmt.Errorf("could not install dependency %s", dep)
				}
			} else {
				fmt.Printf("Installed %s\n", []any{dep[0]}...)
			}
		}
	}
	return 0, nil
}
