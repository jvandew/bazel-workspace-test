package java

import (
	"flag"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type JavaConfig struct{
	SourceTreePrefix string
}

// Configurer satisfies the config.Configurer interface. It's the
// language-specific configuration extension.
type Configurer struct{}

// RegisterFlags registers command-line flags used by the extension. This
// method is called once with the root configuration when Gazelle
// starts. RegisterFlags may set an initial values in Config.Exts. When flags
// are set, they should modify these values.
func (*Configurer) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	javaConfig := JavaConfig{}

	fs.StringVar(
		&javaConfig.SourceTreePrefix,
		"java_source_tree_prefix",
		"src/jvm/",
		"filesystem prefix for java source files",
	)

	c.Exts[JavaName] = javaConfig
}

// CheckFlags validates the configuration after command line flags are parsed.
// This is called once with the root configuration when Gazelle starts.
// CheckFlags may set default values in flags or make implied changes.
func (*Configurer) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	return nil
}

// KnownDirectives returns a list of directive keys that this Configurer can
// interpret. Gazelle prints errors for directives that are not recoginized by
// any Configurer.
func (*Configurer) KnownDirectives() []string {
	// TODO(jacob): this is where command line arguments go
	return make([]string, 0)
}

// Configure modifies the configuration using directives and other information
// extracted from a build file. Configure is called in each directory.
//
// c is the configuration for the current directory. It starts out as a copy
// of the configuration for the parent directory.
//
// rel is the slash-separated relative path from the repository root to
// the current directory. It is "" for the root directory itself.
//
// f is the build file for the current directory or nil if there is no
// existing build file.
func (*Configurer) Configure(c *config.Config, rel string, f *rule.File) {}
