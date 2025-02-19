package rpkgengine

// Total lines in this file: 96
import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Description: Build builds the package
//
// Parameters: It takes the project folder, the build file, and a boolean to remove the project folder after building
//
// Returns: It returns an integer and an error. The integer is the exit code of the build process (1 or 0) and the error is any error that occurred during the build process
func Build(project string, f RpkgBuildFile, removeProjectFolder bool) (int, error) {
	os.Chdir(project + "/Package")
	wd, _ := os.Getwd()
	fmt.Printf("Building package in %v\n", wd)
	switch lang := f.BuildWith; lang {
	case "python3.13":
		// Check if python3.13 is installed
		fmt.Print("Scanning for Python... ")
		cmd := exec.Command("python3.13", "-v")
		cmd.Stdout = nil
		if _, err := cmd.Output(); err != nil {
			fmt.Print("Not Found")
			return 1, errors.New("python 3.13 was not found on your system")
		} else {
			fmt.Print("Found version 3.13\n")
		}
		// Upgrade pip
		fmt.Println("Upgrading pip... ")
		cmd2 := exec.Command("python3.13", "-m", "pip", "install", "--upgrade", "pip")
		cmd2.Stdout = nil
		if _, err := cmd2.Output(); err != nil {
			fmt.Println("Could not upgrade pip")
			return 1, errors.New("could not upgrade pip")
		} else {
			fmt.Println("Pip upgraded successfully")
		}
		// Install build dependencies
		fmt.Println("Installing build dependencies... ")
		if _, err := installDeps(f.BuildDeps, true); err != nil {
			return 1, err
		}
		fmt.Println("Build dependencies installed")
		// Install dependencies
		fmt.Println("Installing dependencies... ")
		if _, err := installDeps(f.Deps, false); err != nil {
			return 1, err
		}
		fmt.Println("Dependencies installed")
		// Run build commands
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
		Cmd := exec.Command("sh", "-c", "'"+cmds+"'")
		Cmd.Stdout = nil
		if _, err := Cmd.Output(); err != nil {
			fmt.Println("Could not run build commands")
			return 1, errors.New("build commands could not be run")
		}
	}
	fmt.Print("Build commands ran successfully.")
	// Clean up
	fmt.Println("Cleaning up... ")
	os.Chdir("../")
	cmd := exec.Command("mv", "./dist", "../dist")
	cmd.Stdout = nil
	if _, err := cmd.Output(); err != nil {
		fmt.Println("Could not move dist folder.")
		return 1, errors.New("could not move dist folder")
	}
	if removeProjectFolder {
		fmt.Println("Removing project folder... ")
		cmd2 := exec.Command("rm", "-rf", "./")
		cmd2.Stdout = nil
		if _, err := cmd2.Output(); err != nil {
			fmt.Println("Could not remove project folder.")
			return 1, errors.New("could not remove project folder")
		}
	}
	// Success (exit code 0)
	fmt.Println("Package built successfully.")
	return 0, nil
}
