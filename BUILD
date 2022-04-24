load("@bazel_gazelle//:def.bzl", "DEFAULT_LANGUAGES", "gazelle", "gazelle_binary")

gazelle_binary(
    name = "gazelle_runner",
    languages = [
        # "@bazel_gazelle//language/go", TODO(jacob): this is mucking up a bazel def for some reason
        "//tools/gazelle/java",
    ],
    visibility = ["//visibility:public"],
)

# Gazelle configuration options.
# See https://github.com/bazelbuild/bazel-gazelle#running-gazelle-with-bazel
# # gazelle:prefix github.com/foursquare/bazel-workspace-test
# gazelle:exclude bazel-out
gazelle(
    name = "gazelle",
    gazelle = ":gazelle_runner",
)

gazelle(
    name = "update_go_deps",
    args = [
        "-from_file=go.mod",
        "-to_macro=tools/gazelle/deps.bzl%internal_gazelle_deps",
        "-prune",
    ],
    command = "update-repos",
)
