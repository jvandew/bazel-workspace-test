module github.com/foursquare/bazel-workspace-test

go 1.18

// TODO(jacob): update everything
require (
  github.com/bazelbuild/bazel-gazelle v0.23.0
  github.com/bazelbuild/buildtools v0.0.0-20200718160251-b1667ff58f71
  // github.com/bazelbuild/rules_go v0.0.0-20190719190356-6dae44dc5cab
  // github.com/bmatcuk/doublestar v1.2.2
  github.com/emirpasic/gods v1.18.1
  // github.com/ghodss/yaml v1.0.0
  // github.com/google/uuid v1.3.0
  // gopkg.in/yaml.v2 v2.2.8

  // NOTE(jacob): https://github.com/bazelbuild/bazel-gazelle/issues/1217
  golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
)
