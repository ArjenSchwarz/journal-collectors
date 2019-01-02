workflow "Build and Publish" {
  on = "push"
  resolves = [ "Deploy" ]
}

action "Build" {
  uses = "apex/actions/go@master"
  runs = "make"
  args = "build"
}

action "Package" {
  uses = "ArjenSchwarz/aws/cli@master"
  needs = ["Build"]
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
