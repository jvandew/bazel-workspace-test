package java

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	bzl "github.com/bazelbuild/buildtools/build"
	"github.com/emirpasic/gods/sets/treeset"
	godsutils "github.com/emirpasic/gods/utils"
)

func parseThirdpartyMapOverrides() map[string][]string {
	jsonBytes, err := os.ReadFile(thirdpartyMapOverrideFile)
	if err != nil {
		log.Fatalf("ERROR: unable to read 3rdparty map override file: \"%s\"", err)
	}

	var res map[string][]string
	if err := json.Unmarshal(jsonBytes, &res); err != nil {
		log.Fatalf("ERROR: unable to parse 3rdparty map override file: \"%s\"", err)
	}
	return res
}

func addTargetPackageToMap(
	thirdpartyMap map[string]interface{},
	targetPackage string,
	maven_coordinates string,
	target string,
) {
	subpackages := strings.Split(targetPackage, ".")
	submap := thirdpartyMap
	for _, subpackage := range(subpackages) {
		if nextmap, exists := submap[subpackage]; exists {
		  submap = nextmap.(map[string]interface{})
		} else {
			nextmap := make(map[string]interface{})
			submap[subpackage] = nextmap
			submap = nextmap
		}
	}

	if wildcard_entry, exists := submap["**"]; exists {
		log.Fatalf(
			"ERROR: package \"%s\" for coordinates \"%s\" already exists in 3rdparty map as \"%s\"",
			targetPackage,
			maven_coordinates,
			wildcard_entry,
		)
	} else {
		submap["**"] = target
	}
}

func parseThirdpartyMap() map[string]interface{} {
	log.Printf("parsing maven targets")
	cmd := exec.Command(
		"bazel",
		"query",
		"kind(jvm_import, @maven//:all)",
		"--output=build",
	)
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to query maven targets due to error: \"%s\"", err)
	}

	overrideMap := parseThirdpartyMapOverrides()

	name_prefix := "  name = \""
	tags_prefix := "  tags = ["
	tag_regex, _ := regexp.Compile("\"maven_coordinates=([a-z0-9-_.:]+)\"")

	scanner := bufio.NewScanner(bytes.NewReader(out))
	var current_target_name *string
	var current_maven_coordinates *string
	res := make(map[string]interface{})

	for scanner.Scan() {
		line := scanner.Text()
		// https://github.com/bazelbuild/rules_jvm_external#generated-targets
		if current_target_name == nil && strings.HasPrefix(line, name_prefix) {
			target_name := fmt.Sprintf("@maven//:%s", line[len(name_prefix):strings.LastIndex(line, "\"")])
			current_target_name = &target_name

		} else if strings.HasPrefix(line, tags_prefix) {
			maven_coordinates := tag_regex.FindStringSubmatch(line)[1]
			current_maven_coordinates = &maven_coordinates

		} else if line == ")" {
			if current_target_name == nil || current_maven_coordinates == nil {
				log.Fatalf(
					"Unparseable maven target? current_target_name: \"%s\", current_maven_coordinates: \"%s\"",
					current_target_name,
					current_maven_coordinates,
				)

			} else if overridePackages, exists := overrideMap[*current_target_name]; exists {
				for _, overridePackage := range(overridePackages) {
					addTargetPackageToMap(
						res,
						overridePackage,
						*current_maven_coordinates,
						*current_target_name,
					)
				}

			} else {
				maven_group := (*current_maven_coordinates)[:strings.Index(*current_maven_coordinates, ":")]
				addTargetPackageToMap(
					res,
					maven_group,
					*current_maven_coordinates,
					*current_target_name,
				)
			}

			current_target_name = nil
			current_maven_coordinates = nil
		}
	}

	return res
}

// Resolver satisfies the resolve.Resolver interface. It resolves dependencies
// in rules generated by this extension.
type Resolver struct{
	/* map from module to maven target providers and sub-modules, eg.
	 * {
	 *	 "com": {
	 *	   "fasterxml": {
	 *	   	 "jackson": {
	 *	   	 	 "core": {
	 *	   	 	 	 "**": "@maven//:com_fasterxml_jackson_core_jackson_core"
	 *	   	 	 },
	 *	   	 	 "databind": {
	 *	   	 	 	 "**": "@maven//:com_fasterxml_jackson_core_jackson_databind"
	 *	   	 	 }
	 *	   	 }
	 *	   }
	 *	 }
	 * }
	 */
	thirdpartyMap map[string]interface{};
}

// Name returns the name of the language. This should be a prefix of the
// kinds of rules generated by the language, e.g., "go" for the Go extension
// since it generates "go_library" rules.
func (*Resolver) Name() string {
	return JavaName
}

// Imports returns a list of ImportSpecs that can be used to import the rule
// r. This is used to populate RuleIndex.
//
// If nil is returned, the rule will not be indexed. If any non-nil slice is
// returned, including an empty slice, the rule will be indexed.
//
// NOTE(jacob): Doc translation: "Return the packages defined by this BUILD target." This
//		is a dead simple implementation currently. Assumptions made:
//		- we have no packages split across build targets
//		- package structure matches directory structure, ignoring the top-level module name
//			and source tree prefix
//		- targets only contain sources in their immediate directory, no sub-directories
func (*Resolver) Imports(
	c *config.Config,
	r *rule.Rule,
	f *rule.File,
) []resolve.ImportSpec {
	cfg := c.Exts[JavaName].(JavaConfig)
	sourceTreePrefix := cfg.SourceTreePrefix
	packageStartIndex := strings.Index(f.Pkg, sourceTreePrefix) + len(sourceTreePrefix)
	javaPackage := strings.ReplaceAll(f.Pkg[packageStartIndex:], "/", ".")

	return []resolve.ImportSpec{
		resolve.ImportSpec{
			Lang: JavaName,
			Imp: javaPackage,
		},
	}
}

// Embeds returns a list of labels of rules that the given rule embeds. If
// a rule is embedded by another importable rule of the same language, only
// the embedding rule will be indexed. The embedding rule will inherit
// the imports of the embedded rule.
func (*Resolver) Embeds(r *rule.Rule, from label.Label) []label.Label {
	// TODO(jacob): nothing to do here for java?
	return make([]label.Label, 0)
}

// Resolve translates imported libraries for a given rule into Bazel
// dependencies. Information about imported libraries is returned for each
// rule generated by language.GenerateRules in
// language.GenerateResult.Imports. Resolve generates a "deps" attribute (or
// the appropriate language-specific equivalent) for each import according to
// language-specific rules and heuristics.
func (self *Resolver) Resolve(
	c *config.Config,
	ruleIndex *resolve.RuleIndex,
	rc *repo.RemoteCache,
	r *rule.Rule,
	imports interface{},
	from label.Label,
) {
	deps := treeset.NewWith(godsutils.StringComparator)

	it := imports.(*treeset.Set).Iterator()
	for it.Next() {
		javaPackage := it.Value().(string)
		javaImportSpec := resolve.ImportSpec{
			Lang: JavaName,
			Imp: javaPackage,
		}
		foundRules := ruleIndex.FindRulesByImportWithConfig(c, javaImportSpec, JavaName)

		if len(foundRules) != 1 {
			if len(foundRules) == 0 {
				subpackages := strings.Split(javaPackage, ".")
				submap := self.thirdpartyMap
				for _, subpackage := range(subpackages) {
					if nextmap, exists := submap[subpackage]; exists {
					  submap = nextmap.(map[string]interface{})
					} else {
						break
					}
				}

				if target, exists := submap["**"]; exists {
					deps.Add(target)
				} else {
					log.Fatalf(
						"ERROR: failed to find a BUILD target containing the \"%s\" package",
						javaPackage,
					)
				}

			} else {
				targets := make([]string, len(foundRules))
				for i, foundRule := range foundRules {
					targets[i] = fmt.Sprintf("%s:%s", foundRule.Label.Pkg, foundRule.Label.Name)
				}
				log.Fatalf(
					"ERROR: multiple BUILD targets containing the \"%s\" package: %+q",
					javaPackage,
					targets,
				)
			}

		} else {
			// TODO(jacob): This assumes the target name matches the package/directory name,
			//		which is at least mostly true but possibly not always true? It does look
			//		cleaner though.
			deps.Add(fmt.Sprintf("//%s", foundRules[0].Label.Pkg))
		}
	}

	if deps.Size() > 0 {
		r.SetAttr("deps", convertDependencySetToExpr(deps))
	}
}

// convertDependencySetToExpr converts the given set of dependencies to an
// expression to be used in the deps attribute.
//
// from https://github.com/bazelbuild/rules_python/blob/27d0c7bb8e663dd2e2e9b295ecbfed680e641dfd/gazelle/resolve.go#L264-L274
func convertDependencySetToExpr(set *treeset.Set) bzl.Expr {
	deps := make([]bzl.Expr, set.Size())
	it := set.Iterator()
	for it.Next() {
		dep := it.Value().(string)
		deps[it.Index()] = &bzl.StringExpr{Value: dep}
	}
	return &bzl.ListExpr{List: deps}
}
