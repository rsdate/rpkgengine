package rpkgengine

// Total lines in this file: 96
import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	e "github.com/rsdate/utils/errors"
	t "github.com/rsdate/utils/types"
)

var (
	errChecker = e.ErrChecker{
		ErrPrefix: "build error",
		PanicMode: "true",
		EM:        e.EM["eMre"],
		TestMode:  false,
	}
	Em = errChecker.EM
)

// Build builds the package using the rpkg.build.yaml file as a struct. It also takes the project path and a boolean to check if the project folder should be removed after building.
func Build(project string, f RpkgBuildFile, removeProjectFolder bool) error {
	os.Chdir(project + "/Package")
	wd := project + "/Package"
	fmt.Printf("Building package in %v\n", wd)
	switch lang := f.BuildWith; lang {
	case "python3.13":
		// Check if python3.13 is installed
		fmt.Print("Scanning for Python... ")
		var val, _ = errChecker.CheckErr(Em[""], func() (any, error) {
			cmd := exec.Command("python3", "--version")
			return cmd.Output()
		})
		verStr := ""
		for i := range val.([]byte) {
			verStr += string(val.([]byte)[i])
		}
		fmt.Printf("Found version %s\n", t.Cast(verStr, false, "string").(string)[7:])
		// Upgrade pip
		fmt.Print("Upgrading pip... ")
		var _, _ = errChecker.CheckErr(Em[""], func() (any, error) {
			cmd2 := exec.Command("python3.13", "-m", "pip", "install", "--upgrade", "pip")
			return cmd2.Output()
		})
		fmt.Println("Pip upgraded successfully")
		// Install build dependencies
		fmt.Println("Installing build dependencies... ")
		var _, _ = errChecker.CheckErr(Em[""], func() (any, error) {
			return nil, installPythonDeps(f.BuildDeps, true)
		})
		fmt.Println("Build dependencies installed")
		// Install dependencies
		fmt.Println("Installing dependencies... ")
		var _, _ = errChecker.CheckErr(Em[""], func() (any, error) {
			return nil, installPythonDeps(f.Deps, false)
		})
		fmt.Println("Dependencies installed")
		// Run build commands
		fmt.Print("Running build commands... ")
		var _, _ = errChecker.CheckErr(Em[""], func() (any, error) {
			cmds := ""
			for i := range f.BuildCommands {
				if cmd, ok := f.BuildCommands[i].(string); ok {
					cmds += cmd + " && "
				} else {
					fmt.Printf("Build command %v is not a string\n", []any{f.BuildCommands[i]}...)
					return nil, fmt.Errorf("build command %v is not a string", []any{f.BuildCommands[i]}...)
				}
			}
			cmds = cmds[:len(cmds)-4]
			Cmd := exec.Command("sh", "-c", cmds)
			Cmd.Stdout = nil
			if _, err := Cmd.Output(); err != nil {
				fmt.Println("Could not run build commands")
				return 1, errors.New("build commands could not be run")
			}
			return nil, nil
		})
		fmt.Println("Build commands ran successfully.")
		// Clean up
		fmt.Print("Cleaning up... ")
		var _, _ = errChecker.CheckErr(Em[""], func() (any, error) {
			cmd := exec.Command("mv", "./dist/", "../dist/")
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
			return nil, nil
		})
	}

	// Success (exit code 0)
	fmt.Println("Package built successfully.")
	return nil
}
