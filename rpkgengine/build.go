package rpkgengine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Build(project string, f RpkgBuildFile) (int, error) {
	os.Chdir(project + "/Package")
	switch lang := f.BuildWith; lang {
	case "python3.13":
		fmt.Printf("Creating new venv for %s...", []any{project}...)
		c := exec.Command("python3.13", "-m", "venv", project)
		c.Stdout = nil
		if _, err := c.Output(); err != nil {
			return 1, errors.New("could not create venv")
		}
		fmt.Print("Scanning for Python... ")
		cmd := exec.Command("python3.13", "-v")
		cmd.Stdout = nil
		if _, err := cmd.Output(); err != nil {
			fmt.Print("Not Found")
			return 1, errors.New("python 3.13 was not found on your system")
		} else {
			fmt.Print("Found version 3.13")
			fmt.Println("Installing build dependencies... ")
		}
		for i := 0; i < len(f.BuildDeps); i++ {
			if len(f.BuildDeps) == 0 {
				fmt.Println("No build dependencies")
				break
			} else if dep, ok := f.BuildDeps[i].(string); ok && !strings.Contains(dep, "@") {
				return 1, fmt.Errorf("build dependency %s is not a valid dependency", dep)
			} else if dep, ok := f.BuildDeps[i].(string); ok && strings.Contains(dep, "@latest") {
				fmt.Printf("Installing %s... ", []any{f.BuildDeps[i]}...)
				dep := strings.Split(f.BuildDeps[i].(string), "@")[0]
				cmd := exec.Command("python3.13", "-m", "pip", "install", "--upgrade", dep)
				cmd.Stdout = nil
				if _, err := cmd.Output(); err != nil {
					fmt.Printf("Could not install %s\n", []any{dep}...)
					return 1, fmt.Errorf("could not install dependency %s", dep)
				} else {
					fmt.Printf("Installed %s\n", []any{dep}...)
				}
			} else {
				dep := strings.Split(f.BuildDeps[i].(string), "@")
				cmd := exec.Command("python3.13", "-m", "pip", "install", dep[0], "==", dep[1])
				cmd.Stdout = nil
				if _, err := cmd.Output(); err != nil {
					fmt.Printf("Could not install %s\n", []any{dep}...)
					return 1, fmt.Errorf("could not install dependency %s", dep)
				} else {
					fmt.Printf("Installed %s\n", []any{dep[0]}...)
				}
			}
		}
		fmt.Println("Build dependencies installed")
		fmt.Println("Installing dependencies... ")
		for i := 0; i < len(f.Deps); i++ {
			if f.Deps[0] == "none" {
				fmt.Println("No dependencies")
				break
			} else if dep, ok := f.Deps[i].(string); ok && !strings.Contains(dep, "@") {
				return 1, fmt.Errorf("dependency %s is not a valid dependency", dep)
			} else if dep, ok := f.Deps[i].(string); ok && strings.Contains(dep, "@latest") {
				fmt.Printf("Installing %s... ", []any{f.Deps[i]}...)
				dep := strings.Split(f.Deps[i].(string), "@")[0]
				cmd := exec.Command("python3.13", "-m", "pip", "install", "--upgrade", dep)
				cmd.Stdout = nil
				if _, err := cmd.Output(); err != nil {
					fmt.Printf("Could not install %s\n", []any{dep}...)
					return 1, fmt.Errorf("could not install dependency %s", dep)
				} else {
					fmt.Printf("Installed %s\n", []any{dep}...)
				}
			} else {
				dep := strings.Split(f.Deps[i].(string), "@")
				cmd := exec.Command("python3.13", "-m", "pip", "install", dep[0], "==", dep[1])
				cmd.Stdout = nil
				if _, err := cmd.Output(); err != nil {
					fmt.Printf("Could not install %s\n", []any{dep}...)
					return 1, fmt.Errorf("could not install dependency %s", dep)
				} else {
					fmt.Printf("Installed %s\n", []any{dep[0]}...)
				}
			}
		}
		fmt.Println("Dependencies installed")
		fmt.Print("Running build commands... ")
		cmds := ""
		for i := 0; i < len(f.BuildCommands); i++ {
			if cmd, ok := f.BuildCommands[i].(string); ok {
				cmds += cmd + " && "
			} else {
				return 1, fmt.Errorf("build command %v is not a string", f.BuildCommands[i])
			}
		}
		cmds = cmds[:len(cmds)-4]
		Cmd := exec.Command("sh", "-c", cmds)
		Cmd.Stdout = nil
		if _, err := Cmd.Output(); err != nil {
			return 1, errors.New("build commands could not be run")
		}
	}
	fmt.Println("Package built successfully.")
	return 0, nil
}
