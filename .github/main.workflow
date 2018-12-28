workflow "Build and Publish" {
  on = "push"
  resolves = "Build"
}

action "Build" {
  uses = "docker://golang:1.11"
  runs = "make"
  args = "github-actions"
}
