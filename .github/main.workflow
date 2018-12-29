workflow "Build and Publish" {
  on = "push"
  resolves = ["GitHub Action for AWS"]
}

action "Build" {
  uses = "docker://golang:1.11"
  runs = "make"
  args = "github-actions"
}

action "GitHub Action for AWS" {
  uses = "actions/aws/cli@8d31870"
  needs = ["Build"]
  runs = "make"
  args = "aws"
  secrets = ["AWS_SECRET_ACCESS_KEY", "AWS_ACCESS_KEY_ID"]
  env = {
      AWS_DEFAULT_REGION = "us-east-1"
  }
}
