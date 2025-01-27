package rpkgengine

import (
	"os"

	mpstruct "github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v3"
)

type Fields struct {
	Name     string   `yaml:"name"`
	Version  string   `yaml:"version"`
	Revision string   `yaml:"revision"`
	Authors  []string `yaml:"authors"`
}

func unmarshall(filename string) (int, Fields) {
	// empty struct for return purposes
	var g Fields
	// read the file
	f, err := os.ReadFile(filename)
	if err != nil {
		return 1, g
	}
	// create the Fields and interface variables (for unmarshalling)
	var F Fields
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
	const build_err = "package was not able to be built"
	code, f := unmarshall("rpkg.build.yaml")
	if code == 1 {
		return 1, build_err
	}
	if err := os.Mkdir(f.Name+"-"+f.Version+"rev"+f.Revision, 0755); err != nil {
		return 1, build_err
	}
	return 0, "package built successfully"
}
