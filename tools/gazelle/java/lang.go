package java

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/emirpasic/gods/sets/treeset"
	godsutils "github.com/emirpasic/gods/utils"
)


const (
	JavaName = "java"
	javaLibraryKind = "java_library"

	importPrefix = "import "
	staticPrefix = "static "

	thirdpartyMapOverrideFile = "3rdparty/jvm/thirdparty_map_overrides.json"
)

type javaLang struct {
	Configurer
	Resolver
}

func (*javaLang) Name() string {
	return JavaName
}

func NewLanguage() language.Language {
	return &javaLang{
		Configurer{},
		Resolver{
			thirdpartyMap: parseThirdpartyMap(),
		},
	}
}

var javaKinds = map[string]rule.KindInfo{
	javaLibraryKind: {
		MatchAny: true, // ?
		NonEmptyAttrs: map[string]bool{
			"deps": true,
			"srcs": true,
			"visibility": true,
		},
		SubstituteAttrs: map[string]bool{}, // ?
		MergeableAttrs: map[string]bool{
			"srcs": true,
		},
		ResolveAttrs: map[string]bool{
			"deps": true,
		},
	},
}

// Kinds returns a map of maps rule names (kinds) and information on how to
// match and merge attributes that may be found in rules of those kinds. All
// kinds of rules generated for this language may be found here.
//
// https://github.com/bazelbuild/bazel-gazelle/blob/master/language/go/kinds.go
// https://github.com/bazelbuild/rules_python/blob/main/gazelle/kinds.go
func (*javaLang) Kinds() map[string]rule.KindInfo {
	return javaKinds
}

var javaLoads = []rule.LoadInfo{
	{
		Name: "@rules_java//java:defs.bzl",
		Symbols: []string{
			javaLibraryKind,
		},
	},
}

// Loads returns .bzl files and symbols they define. Every rule generated by
// GenerateRules, now or in the past, should be loadable from one of these
// files.
//
// https://github.com/bazelbuild/bazel-gazelle/blob/master/language/go/kinds.go
// https://github.com/bazelbuild/rules_python/blob/main/gazelle/kinds.go
func (*javaLang) Loads() []rule.LoadInfo {
	return javaLoads
}

// GenerateRules extracts build metadata from source files in a directory.
// GenerateRules is called in each directory where an update is requested
// in depth-first post-order.
//
// args contains the arguments for GenerateRules. This is passed as a
// struct to avoid breaking implementations in the future when new
// fields are added.
//
// A GenerateResult struct is returned. Optional fields may be added to this
// type in the future.
//
// Any non-fatal errors this function encounters should be logged using
// log.Print.
//
// https://github.com/bazelbuild/bazel-gazelle/blob/master/language/go/generate.go
// https://github.com/bazelbuild/rules_python/blob/main/gazelle/generate.go
//
// NOTE(jacob): We take a pretty basic approach here: look for java source files and do
//		some brute force string parsing to search them for imports. These then get set as a
//		private attribute on the new rule we create, which is how they're passed to the
//		Resolver's Resolve function.
func (*javaLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	cfg := args.Config.Exts[JavaName].(JavaConfig)

	targetName := filepath.Base(args.Rel)

	javaSources := treeset.NewWith(godsutils.StringComparator)
	javaPackageDeps := treeset.NewWith(godsutils.StringComparator)
	for _, filename := range args.RegularFiles {
		if filepath.Ext(filename) == ".java" {
			log.Printf("found file: %s", filename)
			javaSources.Add(filename)

			absolutePath := filepath.Join(args.Dir, filename)
			importedPackages, err := parseImportedPackages(absolutePath)

			if err != nil {
				log.Printf("ERROR: failed to parse imports from %s: %w", absolutePath, err)
			} else {
				javaPackageDeps = javaPackageDeps.Union(importedPackages)
			}
		}
	}

	if javaSources.Empty() {
		return language.GenerateResult{}
	} else {
		javaLibrary := rule.NewRule(javaLibraryKind, targetName)
		javaLibrary.SetAttr("srcs", javaSources.Values())
		javaLibrary.SetAttr("visibility", []string{cfg.DefaultLibraryVisibility})
		javaLibrary.SetPrivateAttr(config.GazelleImportsKey, javaPackageDeps)

		return language.GenerateResult{
			Gen: []*rule.Rule{javaLibrary},
			Imports: []interface{}{javaLibrary.PrivateAttr(config.GazelleImportsKey)},
		}
	}
}

func parseImportedPackages(javaFilename string) (*treeset.Set, error) {
	importedPackages := treeset.NewWith(godsutils.StringComparator)

	file, err := os.Open(javaFilename)
	if err != nil {
		return importedPackages, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	foundImports := false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, importPrefix) {
			foundImports = true

			packageEndIndex := strings.LastIndex(line, ".")
			if packageEndIndex == -1 {
				log.Printf("WARN: possibly malformed import: '%s'", line)
				continue
			}

			importedPackage := line[len(importPrefix):packageEndIndex]
			if strings.HasPrefix(importedPackage, staticPrefix) {
				log.Printf("WARN: static imports not currently supported, skipping: '%s'", line)
				continue
			}

			log.Printf("found imported package: '%s'", importedPackage)
			importedPackages.Add(importedPackage)

		} else if foundImports {
			// We've previously encountered imports and this line doesn't seem to be one,
			// assume we're done.
			break

		} else {
			// still looking...
			continue
		}
	}

	err = scanner.Err()
	return importedPackages, err
}

// Fix repairs deprecated usage of language-specific rules in f. This is
// called before the file is indexed. Unless c.ShouldFix is true, fixes
// that delete or rename rules should not be performed.
//
// https://github.com/bazelbuild/bazel-gazelle/blob/master/language/go/fix.go
// https://github.com/bazelbuild/rules_python/blob/main/gazelle/fix.go
func (*javaLang) Fix(c *config.Config, f *rule.File) {
	// TODO
}
