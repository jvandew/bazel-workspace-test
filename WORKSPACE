load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

SKYLIB_VERSION = "1.0.3"
http_archive(
    name = "bazel_skylib",
    sha256 = "1c531376ac7e5a180e0237938a2536de0c54d93f5c278634818e0efc952dd56c",
    urls = [
        "https://github.com/bazelbuild/bazel-skylib/releases/download/{version}/bazel-skylib-{version}.tar.gz".format(
            version = SKYLIB_VERSION,
        ),
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/{version}/bazel-skylib-{version}.tar.gz".format(
            version = SKYLIB_VERSION,
        ),
    ],
)

RULES_JVM_EXTERNAL_VERSION = "4.2"

http_archive(
    name = "rules_jvm_external",
    sha256 = "cd1a77b7b02e8e008439ca76fd34f5b07aecb8c752961f9640dea15e9e5ba1ca",
    strip_prefix = "rules_jvm_external-{version}".format(version = RULES_JVM_EXTERNAL_VERSION),
    url = "https://github.com/bazelbuild/rules_jvm_external/archive/{version}.zip".format(
        version = RULES_JVM_EXTERNAL_VERSION,
    ),
)

load("@rules_jvm_external//:defs.bzl", "maven_install")

JACKSON_REV = "2.13.2"

maven_install(
    artifacts = [
        "com.fasterxml.jackson.core:jackson-core:{}".format(JACKSON_REV),
        "com.fasterxml.jackson.core:jackson-databind:{}".format(JACKSON_REV),
    ],
    fail_if_repin_required = True,
    maven_install_json = "//3rdparty/jvm:maven_install.json",
    repositories = [
        "https://repo1.maven.org/maven2",
    ],
    strict_visibility = True,
    version_conflict_policy = "pinned",
)

load("@maven//:defs.bzl", "pinned_maven_install")

pinned_maven_install()

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f2dcd210c7095febe54b804bb1cd3a58fe8435a909db2ec04e31542631cf715c",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "5982e5463f171da99e3bdaeff8c0f48283a7a5f396ec5282910b9e8a49c0dd7e",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.25.0/bazel-gazelle-v0.25.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.25.0/bazel-gazelle-v0.25.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
go_rules_dependencies()
go_register_toolchains(version = "1.18")

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

load("//tools/gazelle:deps.bzl", "internal_gazelle_deps")
# gazelle:repository_macro tools/gazelle/deps.bzl%internal_gazelle_deps
internal_gazelle_deps()


