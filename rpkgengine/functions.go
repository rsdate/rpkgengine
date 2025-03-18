package rpkgengine

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	e "github.com/rsdate/utils/errors"
	"github.com/spf13/viper"
)

// Description: Hello is a function to introduce rpkg
//
// Parameters: None
//
// Returns: It returns a string of a text art.
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

// Description: installDeps installs dependencies. It is a helper function for the Build function (also so that it wasn't too long).
//
// Parameters: It takes a list of dependencies and a boolean to check if the dependencies are build dependencies.
//
// Returns: It returns an error. The error is any error that occurred during the dependency installation process.
func installPythonDeps(deps []any, buildDeps bool) error {
	for i := range deps {
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
				return fmt.Errorf("build dependency %s is not a valid dependency", dep)
			} else {
				fmt.Printf("Dependency %s is not a valid dependency\n", []any{dep}...)
				return fmt.Errorf("dependency %s is not a valid dependency", dep)
			}
		} else if dep, ok := deps[i].(string); ok && strings.Contains(dep, "@latest") {
			fmt.Printf("Installing %s... ", []any{deps[i]}...)
			dep := strings.Split(deps[i].(string), "@")[0]
			cmd := exec.Command("python3.13", "-m", "pip", "install", "--upgrade", dep)
			cmd.Stdout = nil
			if _, err := cmd.Output(); err != nil {
				fmt.Printf("Could not install %s\n", []any{dep}...)
				if buildDeps {
					return fmt.Errorf("could not install build dependency %s", dep)
				} else {
					return fmt.Errorf("could not install dependency %s", dep)
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
					return fmt.Errorf("could not install build dependency %s", dep)
				} else {
					return fmt.Errorf("could not install dependency %s", dep)
				}
			} else {
				fmt.Printf("Installed %s\n", []any{dep[0]}...)
			}
		}
	}
	return nil
}

// InitVars initializes the variables from the rpkg.build.yaml file and fills the RpkgBuildFile struct with the data.
func InitVars(viper_instance *viper.Viper) RpkgBuildFile {
	name := viper_instance.Get("name").(string)
	version := viper_instance.Get("version").(string)
	revision := viper_instance.Get("revision").(int)
	authors := viper_instance.Get("authors").([]interface{})
	deps := viper_instance.Get("deps").([]interface{})
	buildDeps := viper_instance.Get("build_deps").([]interface{})
	buildWith := viper_instance.Get("build_with").(string)
	buildCommands := viper_instance.Get("build_commands").([]interface{})
	f := RpkgBuildFile{
		Name:          name,
		Version:       version,
		Revision:      revision,
		Authors:       authors,
		Deps:          deps,
		BuildDeps:     buildDeps,
		BuildWith:     buildWith,
		BuildCommands: buildCommands,
	}
	return f
}

// Build builds the package using the rpkg.build.yaml file as a struct. It also takes the project path and a boolean to check if the project folder should be removed after building.
func Build(project string, f RpkgBuildFile, removeProjectFolder bool, errChecker e.ErrChecker) error {
	os.Chdir(project + "/Package")
	wd, _ := os.Getwd()
	fmt.Printf("Building package in %v\n", wd)
	switch lang := f.BuildWith; lang {
	case "python3.13":
		// Check if python3.13 is installed
		fmt.Print("Scanning for Python... ")
		var val, _ = errChecker.CheckErr(Emb[""], func() (any, error) {
			cmd := exec.Command("python3", "--version")
			return cmd.Output()
		})
		verStr := ""
		for i := range val.([]byte) {
			verStr += string(val.([]byte)[i])
		}
		fmt.Printf("Found version %s", verStr[7:])
		// Upgrade pip
		fmt.Print("Upgrading pip... ")
		var _, _ = errChecker.CheckErr(Emb[""], func() (any, error) {
			cmd2 := exec.Command("python3.13", "-m", "pip", "install", "--upgrade", "pip")
			return cmd2.Output()
		})
		fmt.Println("Pip upgraded successfully")
		// Install build dependencies
		fmt.Println("Installing build dependencies... ")
		var _, _ = errChecker.CheckErr(Emb[""], func() (any, error) {
			return nil, installPythonDeps(f.BuildDeps, true)
		})
		fmt.Println("Build dependencies installed")
		// Install dependencies
		fmt.Println("Installing dependencies... ")
		var _, _ = errChecker.CheckErr(Emb[""], func() (any, error) {
			return nil, installPythonDeps(f.Deps, false)
		})
		fmt.Println("Dependencies installed")
		// Run build commands
		fmt.Print("Running build commands... ")
		var _, _ = errChecker.CheckErr(Emb[""], func() (any, error) {
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
			fmt.Println(cmds)
			Cmd := exec.Command("sh", "-c", cmds)
			Cmd.Stdout = nil
			if _, err := Cmd.Output(); err != nil {
				return nil, errors.New("build commands could not be run")
			}
			return nil, nil
		})
		fmt.Println("Build commands ran successfully.")
		// Clean up
		fmt.Println("Cleaning up... ")
		var _, _ = errChecker.CheckErr(Emb[""], func() (any, error) {
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

func DownloadPackage(filepath string, url string, errChecker e.ErrChecker) error {
	// Create the file
	fmt.Fprint(os.Stdout, []any{"Creating the file... "}...)
	out, err := os.Create(filepath)
	errChecker.CheckErr(Emb["dwnerr1"], func() (any, error) {
		return nil, err
	})
	fmt.Fprintln(os.Stdout, []any{"File created successfully."}...)
	defer out.Close()

	// Get the data
	fmt.Fprint(os.Stdout, []any{"Downloading the file... "}...)
	resp, err := http.Get(url)
	errChecker.CheckErr(Emb["dwnerr2"], func() (any, error) {
		return nil, err
	})
	fmt.Fprintln(os.Stdout, []any{"File downloaded successfully."}...)
	defer resp.Body.Close()

	// Check server response
	switch code := resp.StatusCode; code {
	case http.StatusNotFound:
		errChecker.CheckErr(Emb["dwnerr3"], func() (any, error) {
			return nil, errors.New(Emb["dwnerr3"])
		})
	case http.StatusForbidden:
		errChecker.CheckErr(Emb["dwnerr4"], func() (any, error) {
			return nil, errors.New(Emb["dwnerr4"])
		})
	case http.StatusUnauthorized:
		errChecker.CheckErr(Emb["dwnerr5"], func() (any, error) {
			return nil, errors.New(Emb["dwnerr5"])
		})
	case http.StatusInternalServerError:
		errChecker.CheckErr(Emb["dwnerr6"], func() (any, error) {
			return nil, errors.New(Emb["dwnerr6"])
		})
	case http.StatusServiceUnavailable:
		errChecker.CheckErr(Emb["dwnerr7"], func() (any, error) {
			return nil, errors.New(Emb["dwnerr7"])
		})
	}
	// Write the body to file
	fmt.Fprint(os.Stdout, []any{"Writing the download to the file... "}...)
	_, err = io.Copy(out, resp.Body)
	errChecker.CheckErr(Emb["dwnerr8"], func() (any, error) {
		return nil, err
	})
	fmt.Fprintln(os.Stdout, []any{"File written to successfully."}...)
	return nil
}

func BuildPackage(projectPath string, errChecker e.ErrChecker) error {
	if viper_instance == nil {
		errChecker.CheckErr(Emb["blderr1"], func() (any, error) {
			return nil, errors.New(Emb["blderr1"])
		})
	}
	fmt.Fprint(os.Stdout, []any{"Building package... "}...)
	errChecker.CheckErr(Emb["blderr2"], func() (any, error) {
		f := InitVars(viper_instance)
		os.Chdir(projectPath + "/Package")
		err := Build(projectPath, f, false, ErrCheckerBuild)
		return nil, err
	})
	fmt.Fprintln(os.Stdout, []any{"Build successful."}...)
	return nil
}

func InstallPackage(downloadPath string, projectPath string, dirName string, noInstallConf bool, errChecker e.ErrChecker) error {
	fullName := "https://" + os.Getenv(mirror) + "/projects/" + projectPath
	fmt.Fprintf(os.Stdout, "The package path on the mirror is %s and it will download to %s.\nWould you like to proceed with the installation? [Y or n]", []any{projectPath, downloadPath}...)
	fmt.Fscan(os.Stdin, &conf)
	if conf == "Y" || noInstallConf {
		fmt.Fprint(os.Stdout, []any{"Downloading package... "}...)
		errChecker.CheckErr(Emb["insterr1"], func() (any, error) {
			err := DownloadPackage(downloadPath, fullName, errChecker)
			return nil, err
		})
		fmt.Fprintln(os.Stdout, []any{"Package downloaded successfully."}...)
		fmt.Fprint(os.Stdout, []any{"Unziping package... "}...)
		errChecker.CheckErr(Emb["insterr2"], func() (any, error) {
			cmd := exec.Command("tar", "-xzf", projectPath)
			cmd.Stdout = nil
			err := cmd.Run()
			return nil, err
		})
		fmt.Fprintln(os.Stdout, []any{"Package unziped successfully."}...)
		fmt.Fprint(os.Stdout, []any{"Building package... "}...)
		errChecker.CheckErr(Emb["insterr3"], func() (any, error) {
			os.Chdir(dirName)
			err := BuildPackage(".", errChecker)
			return nil, err
		})
		fmt.Fprintln(os.Stdout, []any{"Installation completed! 🎉"}...)
	} else if conf == "n" {
		fmt.Fprintln(os.Stdout, []any{"Installation aborted."}...)
		os.Exit(0)
	}
	return nil
}
