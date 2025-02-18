package rpkgengine

import (
	"fmt"
	"os/exec"
	"strings"
)

// RpkgBuildFile is the struct for the rpkg.build.yaml file
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
func installDeps(deps []interface{}, buildDeps bool) (int, error) {
	// Install dependencies
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
