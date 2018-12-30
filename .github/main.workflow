workflow "Build and Publish" {
  on = "push"
  resolves = ["Package"]
}

action "Build" {
  uses = "docker://golang:1.11"
  runs = "make"
  args = "github-actions"
}

action "Package" {
  uses = "ArjenSchwarz/aws/cli@master"
  needs = ["Build"]
  args = "cloudformation deploy --template-file ./packaged-template.yaml --stack-name journal-collectors"
  secrets = ["AWS_SECRET_ACCESS_KEY", "AWS_ACCESS_KEY_ID"]
  env = {
      AWS_DEFAULT_REGION = "us-east-1",
      ONLY_IN_BRANCH = "master"
  }
}
