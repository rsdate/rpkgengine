package rpkgengine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func Build(project string, f RpkgBuildFile) (int, error) {
	os.Chdir(project + "/Package")
	switch lang := f.BuildWith; lang {
	case "python3.13":
		fmt.Printf("Creating new venv for %s...\n", []any{project}...)
		c := exec.Command("python3.13", "-m", "venv", project)
		c.Stdout = nil
		if _, err := c.Output(); err != nil {
			return 1, errors.New("could not create venv")
		}
		fmt.Println("Venv created")
		fmt.Print("Scanning for Python... ")
		cmd := exec.Command("python3.13", "-v")
		cmd.Stdout = nil
		if _, err := cmd.Output(); err != nil {
			fmt.Print("Not Found")
			return 1, errors.New("python 3.13 was not found on your system")
		} else {
			fmt.Print("Found version 3.13\n")
			fmt.Println("Installing build dependencies... ")
		}
		if _, err := installDeps(f.BuildDeps, true); err != nil {
			return 1, err
		}
		fmt.Println("Build dependencies installed")
		fmt.Println("Installing dependencies... ")
		if _, err := installDeps(f.Deps, false); err != nil {
			return 1, err
		}
		fmt.Println("Dependencies installed")
		fmt.Print("Running build commands... ")
		cmds := ""
		for i := 0; i < len(f.BuildCommands); i++ {
			if cmd, ok := f.BuildCommands[i].(string); ok {
				cmds += cmd + " && "
			} else {
				fmt.Printf("Build command %v is not a string\n", []any{f.BuildCommands[i]}...)
				return 1, fmt.Errorf("build command %v is not a string", f.BuildCommands[i])
			}
		}
		cmds = cmds[:len(cmds)-4]
		Cmd := exec.Command("sh", "-c", cmds)
		Cmd.Stdout = nil
		if _, err := Cmd.Output(); err != nil {
			fmt.Println("Could not run build commands")
			return 1, errors.New("build commands could not be run")
		}
	}
	fmt.Println("Package built successfully.")
	return 0, nil
}
