load("@rules_java//java:defs.bzl", "java_binary", "java_library")

java_binary(
    name = "hello-bin",
    main_class = "hello.Hello",
    runtime_deps = [
        ":hello",
    ],
)

java_library(
    name = "hello",
    srcs = ["Hello.java"],
    visibility = ["//visibility:public"],
    deps = [
        "@maven//:com_fasterxml_jackson_core_jackson_core",
        "@maven//:com_fasterxml_jackson_core_jackson_databind",
        "//hello/src/jvm/hello/util",
        "//shared/src/jvm/common/print",
    ],
)
