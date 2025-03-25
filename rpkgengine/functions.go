package rpkgengine

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	rec "github.com/rsdate/rpkgengine/rpkgengineconfig"
	e "github.com/rsdate/utils/errors"
	pbar "github.com/schollz/progressbar/v3"
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

// installDeps is a function to install the dependecies (or build dependencies) from deps
func InstallDeps(deps []any, buildDeps bool, name string, installCommands map[string][]string, versionDelimiter string) error {
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
			commands := []string{}
			commands = append(commands, installCommands["latest"]...)
			commands = append(commands, dep)
			cmd := exec.Command(name, commands...)
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
			fmt.Printf("Installing %s... ", []any{deps[i]}...)
			dep := strings.Split(deps[i].(string), "@")
			commands := []string{}
			commands = append(commands, installCommands["version"]...)
			commands = append(commands, dep[0]+versionDelimiter+dep[1])
			cmd := exec.Command(name, commands...)
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
	f := RpkgBuildFile{
		Name:          viper_instance.Get("name").(string),
		Version:       viper_instance.Get("version").(string),
		Revision:      viper_instance.Get("revision").(int),
		Authors:       viper_instance.Get("authors").([]any),
		Deps:          viper_instance.Get("deps").([]any),
		BuildDeps:     viper_instance.Get("build_deps").([]any),
		BuildWith:     viper_instance.Get("build_with").(string),
		BuildCommands: viper_instance.Get("build_commands").([]any),
	}
	return f
}

// Build builds the package using the rpkg.build.yaml file as a struct. It also takes the project path and a boolean to check if the project folder should be removed after building.
func Build(project string, f RpkgBuildFile, removeProjectFolder bool, errChecker e.ErrChecker) error {
	os.Chdir(project + "/Package")
	wd, _ := os.Getwd()
	fmt.Printf("Building package in %v\n", wd[:len(wd)-8])
	switch lang := f.BuildWith; lang {
	case "python3.13":
		// Check if python3.13 is installed
		fmt.Print("Scanning for Python... ")
		var val, _ = errChecker.CheckErr("blderr1", func() (any, error) {
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
		errChecker.CheckErr("blderr2", func() (any, error) {
			cmd2 := exec.Command("python3.13", "-m", "pip", "install", "--upgrade", "pip")
			return nil, cmd2.Run()
		})
		fmt.Println("Pip upgraded successfully")
		// Install build dependencies
		fmt.Println("Installing build dependencies... ")
		errChecker.CheckErr("blderr3", func() (any, error) {
			return nil, InstallDeps(f.BuildDeps, true, "pip", map[string][]string{
				"latest":  {"install", "--upgrade"},
				"version": {"install"},
			}, "==")
		})
		fmt.Println("Build dependencies installed")
		// Install dependencies
		fmt.Println("Installing dependencies... ")
		errChecker.CheckErr("blderr4", func() (any, error) {
			return nil, InstallDeps(f.Deps, false, "pip", map[string][]string{
				"latest":  {"install", "--upgrade"},
				"version": {"install"},
			}, "==")
		})
		fmt.Println("Dependencies installed")
		// Run build commands
		fmt.Print("Running build commands... ")
		errChecker.CheckErr("blderr5", func() (any, error) {
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
		fmt.Print("Cleaning up... ")
		var _, _ = errChecker.CheckErr("blderr6", func() (any, error) {
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
	fmt.Println("Clean up successful.")
	// Success (exit code 0)
	fmt.Println("Package built successfully.")
	return nil
}

// DownloadPackage downloads a package to the filepath parameter from the url.
func DownloadPackage(filepath string, url string, errChecker e.ErrChecker) error {
	// Create the file
	fmt.Fprint(os.Stdout, []any{"Creating the file... "}...)
	out, err := os.Create(filepath)
	errChecker.CheckErr("dwnerr1", func() (any, error) {
		return nil, err
	})
	fmt.Fprintln(os.Stdout, []any{"File created successfully."}...)
	defer out.Close()
	// Get the data
	resp, err := http.Get(url)
	errChecker.CheckErr("dwnerr2", func() (any, error) {
		return nil, err
	})
	defer resp.Body.Close()
	// Check server response
	switch code := resp.StatusCode; code {
	case http.StatusNotFound:
		errChecker.CheckErr("dwnerr3", func() (any, error) {
			return nil, errors.New(Emre["dwnerr3"])
		})
	case http.StatusForbidden:
		errChecker.CheckErr("dwnerr4", func() (any, error) {
			return nil, errors.New(Emre["dwnerr4"])
		})
	case http.StatusUnauthorized:
		errChecker.CheckErr("dwnerr5", func() (any, error) {
			return nil, errors.New(Emre["dwnerr5"])
		})
	case http.StatusInternalServerError:
		errChecker.CheckErr("dwnerr6", func() (any, error) {
			return nil, errors.New(Emre["dwnerr6"])
		})
	case http.StatusServiceUnavailable:
		errChecker.CheckErr("dwnerr7", func() (any, error) {
			return nil, errors.New(Emre["dwnerr7"])
		})
	}
	// Create a progress bar
	bar := pbar.DefaultBytes(
		resp.ContentLength,
		"Downloading the file...",
	)
	errChecker.CheckErr("", func() (any, error) {
		_, err := io.Copy(io.MultiWriter(out, bar), resp.Body)
		return nil, err
	})
	return nil
}

// BuildPackage builds the package using the Build() function
func BuildPackage(projectPath string, errChecker e.ErrChecker, viper_instance *viper.Viper) error {
	// Check if viper_instance is nil.
	if viper_instance == nil {
		errChecker.CheckErr("bldperr1", func() (any, error) {
			return nil, errors.New(Emre["bldperr1"])
		})
	}
	// Build the package.
	fmt.Fprint(os.Stdout, []any{"Building package... "}...)
	errChecker.CheckErr("bldperr2", func() (any, error) {
		rec.InitConfig(projectPath)
		f := InitVars(viper_instance)
		os.Chdir(projectPath + "/Package")
		err := Build(projectPath, f, false, ErrCheckerBuild)
		return nil, err
	})
	fmt.Fprintln(os.Stdout, []any{"Build successful."}...)
	return nil
}

// InstallPackage installs the package from the given mirror.
func InstallPackage(dirName string, noInstallConf bool, errChecker e.ErrChecker, viper_instance *viper.Viper) error {
	projectPath := dirName + ".tar.gz"
	downloadPath := "./" + projectPath
	fullName := "https://" + os.Getenv(mirror) + "/projects/" + projectPath
	fmt.Fprintf(os.Stdout, "The package path on the mirror is %s and it will download to %s.\nWould you like to proceed with the installation? [Y or n]", []any{projectPath, downloadPath}...)
	fmt.Fscan(os.Stdin, &conf)
	if conf == "Y" || noInstallConf {
		fmt.Fprint(os.Stdout, []any{"Downloading package... "}...)
		errChecker.CheckErr("insterr1", func() (any, error) {
			err := DownloadPackage(downloadPath, fullName, errChecker)
			return nil, err
		})
		fmt.Fprintln(os.Stdout, []any{"Package downloaded successfully."}...)
		fmt.Fprint(os.Stdout, []any{"Unziping package... "}...)
		errChecker.CheckErr("insterr2", func() (any, error) {
			cmd := exec.Command("tar", "-xzf", projectPath)
			cmd.Stdout = nil
			err := cmd.Run()
			return nil, err
		})
		fmt.Fprintln(os.Stdout, []any{"Package unziped successfully."}...)
		fmt.Fprint(os.Stdout, []any{"Building package... "}...)
		errChecker.CheckErr("insterr3", func() (any, error) {
			os.Chdir(dirName)
			err := BuildPackage(".", errChecker, viper_instance)
			return nil, err
		})
		fmt.Fprintln(os.Stdout, []any{"Installation completed! 🎉"}...)
	} else if conf == "n" {
		fmt.Fprintln(os.Stdout, []any{"Installation aborted."}...)
		os.Exit(0)
	}
	return nil
}
