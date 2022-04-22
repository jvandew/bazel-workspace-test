load("@rules_java//java:defs.bzl", "java_library")

java_library(
    name = "hello",
    srcs = glob(["*.java"]),
    deps = [
        "//hello/src/jvm/hello/util",
        "//shared/src/jvm/common/print",
    ],
)
