load("@rules_java//java:defs.bzl", "java_binary", "java_library")

java_library(
    name = "hello",
    srcs = ["Hello.java"],
    visibility = ["//visibility:public"],
)

java_binary(
    name = "hello-bin",
    main_class = "hello.Hello",
    runtime_deps = [
        ":hello",
    ],
)
