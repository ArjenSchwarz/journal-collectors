workflow "Build and Publish" {
  on = "push"
  resolves = [
    "Deploy",
    "apex/actions/go@master-1",
  ]
}

action "Build" {
  uses = "docker://golang:1.11"
  runs = "make"
  args = "github-actions"
}

# Filter for master branch
action "Master" {
  needs = "Build"
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "Package" {
  uses = "ArjenSchwarz/aws/cli@master"
  needs = ["Master"]
  args = "cloudformation package --template-file ./template.yaml --s3-bucket public.ig.nore.me --output-template-file packaged-template.yaml"
  secrets = ["AWS_SECRET_ACCESS_KEY", "AWS_ACCESS_KEY_ID"]
  env = {
    ONLY_IN_BRANCH = "master"
  }
}

action "Deploy" {
  uses = "ArjenSchwarz/aws/cli@master"
  needs = ["Package"]
  args = "cloudformation deploy --template-file ./packaged-template.yaml --stack-name journal-collectors"
  secrets = ["AWS_SECRET_ACCESS_KEY", "AWS_ACCESS_KEY_ID"]
  env = {
    AWS_DEFAULT_REGION = "us-east-1"
    ONLY_IN_BRANCH = "master"
  }
}

action "apex/actions/go@master" {
  uses = "apex/actions/go@master"
  args = "get -v ./..."
}

action "apex/actions/go@master-1" {
  uses = "apex/actions/go@master"
  needs = ["apex/actions/go@master"]
  runs = "make"
  args = "test"
}
