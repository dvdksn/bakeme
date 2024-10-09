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

target "binaries" {
  matrix = {
    mode = ["release", "debug"]
  }
  name = "binaries-${mode}"
  target = "binaries"
  args = {
    BUILD_TAGS = mode
  }
  output = ["type=local,dest=./build/${mode}"]
  platforms = [
    "linux/amd64",
    "linux/arm64",
    "linux/riscv64",
  ]
}

target "localbin" {
  # Enable debug mode
  args = {
    BUILD_TAGS = "debug"
  }
  target = "binaries"
  output = ["type=local,dest=./build/local"]
  platforms = ["local"]
}
