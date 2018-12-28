workflow "Build and Publish" {
  on = "push"
  resolves = "Lint"
}

action "Lint" {
  uses = "docker://golang:1.11"
  runs = "GOPATH=`pwd`;go"
  args = "get -t ./..."
}
