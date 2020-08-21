data "aws_ssm_parameter" "slack_token" {
  name = "/env/camera-services/slack/slack-token"
}

data "aws_ssm_parameter" "slack_channel" {
  name = "/env/camera-services/slack/slack-channel"
}

module "slack" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-slack"
  image          = "docker.pkg.github.com/byuoitav/camera-services/camera-slack-dev"
  image_version  = "26c4c7c"
  container_port = 8080 // doesn't actually have a port...
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  container_env     = {}
  container_args = [
    "--db-address", data.aws_ssm_parameter.prd_db_addr.value,
    "--db-username", data.aws_ssm_parameter.prd_db_username.value,
    "--db-password", data.aws_ssm_parameter.prd_db_password.value,
    "--hub-address", data.aws_ssm_parameter.hub_address.value,
    "--aver-username", data.aws_ssm_parameter.aver_username.value,
    "--aver-password", data.aws_ssm_parameter.aver_password.value,
    "--snapshot-delay", "3s",
    "--slack-token", data.aws_ssm_parameter.slack_token.value,
    "--channel-id", data.aws_ssm_parameter.slack_channel.value,
  ]
  health_check = false
}
