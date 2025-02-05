package rpkgengine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	mpstruct "github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v3"
)

type RpkgBuildFile struct {
	Name          string   `yaml:"name"`
	Version       string   `yaml:"version"`
	Revision      string   `yaml:"revision"`
	Authors       []string `yaml:"authors"`
	Deps          []string `yaml:"deps"`
	BuildDeps     []string `yaml:"build_deps"`
	BuildWith     string   `yaml:"build_with"`
	BuildCommands []string `yaml:"build_commands"`
}

func unmarshall(filename string) (int, RpkgBuildFile, error) {
	// Empty struct for return purposes
	var g RpkgBuildFile
	// Read the file
	f, err := os.ReadFile(filename)
	if err != nil {
		return 1, g, errors.New("wasn't able to read file")
	}
	// Create the RpkgBuildFile and interface variables (for unmarshalling)
	var F RpkgBuildFile
	var raw interface{}
	// Unmarshall the YAML
	if err := yaml.Unmarshal(f, &raw); err != nil {
		return 1, g, errors.New("YAML unmarshalling failed")
	}
	// Decode the unmarshalled YAML using mapstructure (shortened to mpstruct)
	decoder, _ := mpstruct.NewDecoder(&mpstruct.DecoderConfig{WeaklyTypedInput: true, Result: &F})
	if err := decoder.Decode(raw); err != nil {
		return 1, g, errors.New("YAML decoding failed")
	}
	return 0, F, nil
}

func build(project string) (int, error) {
	code, f, err := unmarshall(project + "/rpkg.build.yaml")
	if code == 1 && err != nil {
		return 1, errors.New("YAML unmarshalling failed")
	}
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
			fmt.Printf("Installing %s... ", []any{f.BuildDeps[i]}...)
			if strings.Contains(f.BuildDeps[i], "@latest") {
				dep := strings.Split(f.BuildDeps[i], "@")[0]
				cmd := exec.Command("python3.13", "-m", "pip", "install", "--upgrade", dep)
				cmd.Stdout = nil
				if _, err := cmd.Output(); err != nil {
					fmt.Printf("Could not install %s\n", []any{dep}...)
					return 1, fmt.Errorf("could not install dependency %s", dep)
				} else {
					fmt.Printf("Installed %s\n", []any{dep}...)
				}
			} else {
				dep := strings.Split(f.BuildDeps[i], "@")
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
			fmt.Printf("Installing %s... ", []any{f.Deps[i]}...)
			if strings.Contains(f.Deps[i], "@latest") {
				dep := strings.Split(f.Deps[i], "@")[0]
				cmd := exec.Command("python3.13", "-m", "pip", "install", "--upgrade", dep)
				cmd.Stdout = nil
				if _, err := cmd.Output(); err != nil {
					fmt.Printf("Could not install %s\n", []any{dep}...)
					return 1, fmt.Errorf("could not install dependency %s", dep)
				} else {
					fmt.Printf("Installed %s\n", []any{dep}...)
				}
			} else {
				dep := strings.Split(f.Deps[i], "@")
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
			cmds += f.BuildCommands[i] + " && "
		}
		cmds = cmds[:len(cmds)-4]
		Cmd := exec.Command(cmds)
		Cmd.Stdout = nil
		if _, err := Cmd.Output(); err != nil {
			return 1, errors.New("build commands could not be run")
		}
	}
	fmt.Println("Package built successfully.")
	return 0, nil
}
