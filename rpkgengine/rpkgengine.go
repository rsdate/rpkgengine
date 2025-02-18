package rpkgengine

// RpkgBuildFile is the struct for the rpkg.build.yaml file
type RpkgBuildFile struct {
	Name          string
	Version       string
	Revision      int
	Authors       []interface{}
	Deps          []interface{}
	BuildDeps     []interface{}
	BuildWith     string
	BuildCommands []interface{}
}

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
