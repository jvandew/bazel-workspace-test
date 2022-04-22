load("@rules_java//java:defs.bzl", "java_binary", "java_library")

java_library(
    name = "hello",
    srcs = glob(["*.java"]),
    deps = [
        "//hello/src/jvm/hello/util",
        "//shared/src/jvm/common/print",
    ],
)

java_binary(
    name = "hello-bin",
    main_class = "hello.Hello",
    runtime_deps = [
        ":hello",
    ],
)
