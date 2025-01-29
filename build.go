package rpkgengine

import (
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
	BuildWith     string   `yaml:"build_with"`
	BuildCommands []string `yaml:"build_commands"`
}

func unmarshall(filename string) (int, RpkgBuildFile) {
	// Empty struct for return purposes
	var g RpkgBuildFile
	// Read the file
	f, err := os.ReadFile(filename)
	if err != nil {
		return 1, g
	}
	// Create the RpkgBuildFile and interface variables (for unmarshalling)
	var F RpkgBuildFile
	var raw interface{}
	// Unmarshall the YAML
	if err := yaml.Unmarshal(f, &raw); err != nil {
		return 1, g
	}
	// Decode the unmarshalled YAML using mapstructure (shortened to mpstruct)
	decoder, _ := mpstruct.NewDecoder(&mpstruct.DecoderConfig{WeaklyTypedInput: true, Result: &F})
	if err := decoder.Decode(raw); err != nil {
		return 1, g
	}
	return 0, F
}

func build() (int, string) {
	code, f := unmarshall("rpkg.build.yaml")
	if code == 1 {
		return 1, "YAML unmarshalling failed."
	}
	switch lang := f.BuildWith; lang {
	case "python3.13":
		cmd := exec.Command("python3.13", "-v")
		if _, err := cmd.Output(); err != nil {
			return 1, "Python 3.13 was not found on your system."
		} else {
			cmds := ""
			for i := 0; i < len(f.BuildCommands); i++ {
				cmds = cmds + f.BuildCommands[i] + " && "
			}
			cmds = cmds[:len(cmds)-4]
			cmd := exec.Command(cmds)
			if _, err := cmd.Output(); err != nil {
				return 1, "Build commands could not be run"
			}
		}
	}
	return 0, "Package built successfully"
}
