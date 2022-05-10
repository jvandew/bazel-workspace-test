load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "java",
    srcs = [
        "configure.go",
        "lang.go",
        "resolve.go",
    ],
    data = [
        "//3rdparty/jvm:thirdparty_map_overrides.json",
    ],
    importpath = "github.com/foursquare/bazel-workspace-test/tools/gazelle/java",
    visibility = ["//visibility:public"],
    deps = [
        "@bazel_gazelle//config:go_default_library",
        "@bazel_gazelle//label:go_default_library",
        "@bazel_gazelle//language:go_default_library",
        "@bazel_gazelle//repo:go_default_library",
        "@bazel_gazelle//resolve:go_default_library",
        "@bazel_gazelle//rule:go_default_library",
        "@com_github_bazelbuild_buildtools//build:go_default_library",
        "@com_github_emirpasic_gods//sets/treeset",
        "@com_github_emirpasic_gods//utils",
    ],
)
