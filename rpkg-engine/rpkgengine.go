package rpkgengine

// RpkgBuildFile is the struct for the rpkg.build.yaml file
type RpkgBuildFile struct {
	Name          string
	Version       string
	Revision      int
	Authors       []string
	Deps          []string
	BuildDeps     []string
	BuildWith     string
	BuildCommands []string
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
