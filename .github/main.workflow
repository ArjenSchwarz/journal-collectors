workflow "Build and Publish" {
  on = "push"
  resolves = "Lint"
}

action "Lint" {
  uses = "docker://golang:1.11"
  runs = "go"
  args = "get -t ./..."
  env = "GOPATH=/github/workspace"
}
