workflow "Build and Publish" {
  on = "push"
  resolves = "Lint"
}

action "GoGet" {
  uses = "docker://golang:1.11"
  runs = "go"
  args = "get -u ./..."
  env = {
      GOPATH = "/github/workspace"
  }
}

action "Lint" {
  needs = "GoGet"
  uses = "docker://golang:1.11"
  runs = "golint"
  args = "
      GOPATH = "/github/workspace"
  }
}
