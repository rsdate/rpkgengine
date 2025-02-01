package rpkgengine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

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

func build() (int, error) {
	code, f, err := unmarshall("./rpkg.build.yaml")
	if code == 1 && err != nil {
		return 1, errors.New("YAML unmarshalling failed")
	}
	switch lang := f.BuildWith; lang {
	case "python3.13":
		fmt.Print("Finding Python... ")
		cmd := exec.Command("python3.13", "-v")
		if _, err := cmd.Output(); err != nil {
			return 1, errors.New("python 3.13 was not found on your system")
		} else {
			fmt.Print("Found version 3.13")
			fmt.Println("Installing dependencies...")
		}

		cmds := ""
		for i := 0; i < len(f.BuildCommands); i++ {
			cmds = cmds + f.BuildCommands[i] + " && "
		}
		cmds = cmds[:len(cmds)-4]
		Cmd := exec.Command(cmds)
		if _, err := Cmd.Output(); err != nil {
			return 1, errors.New("build commands could not be run")
		}
	}
	return 0, nil
}
