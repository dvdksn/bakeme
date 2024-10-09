target "default" {
  matrix = {
    mode = ["release", "debug"]
  }
  name = "image-${mode}"
  target = "image"
  args = {
    BUILD_TAGS = mode
  }
  tags = [
    mode == "release" ? "bakeme:latest" : "bakeme:dev"
  ]
  attest = [
    "type=provenance,mode=max",
    "type=sbom",
  ]
  platforms = [
    "linux/amd64",
    "linux/arm64",
    "linux/riscv64",
  ]
}

target "test" {
  target = "test"
  output = ["type=cacheonly"]
}

target "lint" {
  target = "lint"
  output = ["type=cacheonly"]
}

group "validate" {
  targets = ["test", "lint"]
}

target "bin" {
  target = "bin"
  output = ["build/bin"]
  platforms = ["local"]
}

target "bin-cross" {
  inherits = ["bin"]
  platforms = [
    "linux/amd64",
    "linux/arm64",
    "linux/riscv64",
  ]
}

target "bin-all" {
  inherits = ["bin-cross"]
  matrix = {
    mode = ["release", "debug"]
  }
  name = "bin-${mode}"
  args = {
    BUILD_TAGS = mode
  }
  output = ["build/bin/${mode}"]
}
