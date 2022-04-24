package java

import (
	"log"
  "path/filepath"

  "github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
  "github.com/bazelbuild/bazel-gazelle/rule"
  "github.com/emirpasic/gods/sets/treeset"
  godsutils "github.com/emirpasic/gods/utils"
)


const (
	JavaName = "java"
	javaLibraryKind = "java_library"
)

type javaLang struct {
  Configurer
  Resolver
}

func (*javaLang) Name() string {
  return JavaName
}

func NewLanguage() language.Language {
  return &javaLang{}
}

var javaKinds = map[string]rule.KindInfo{
  javaLibraryKind: {
    MatchAny: true, // ?
    NonEmptyAttrs: map[string]bool{
      "deps":       true,
      "srcs":       true,
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
func (*javaLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
  // cfgs := args.Config.Exts[languageName].(pythonconfig.Configs)
  // cfg := cfgs[args.Rel]

  // if !cfg.ExtensionEnabled() {
  //   return language.GenerateResult{}
  // }

  javaLibraryFilenames := treeset.NewWith(godsutils.StringComparator)

  for _, f := range args.RegularFiles {
  	if filepath.Ext(f) == ".java" {
  		log.Printf("found file: %s", f)
      javaLibraryFilenames.Add(f)
    }
  }

  return language.GenerateResult{}
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
