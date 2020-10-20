data "aws_ssm_parameter" "aver_username" {
  name = "/env/camera-services/aver/username"
}

data "aws_ssm_parameter" "aver_password" {
  name = "/env/camera-services/aver/password"
}

module "aver_dev" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-aver-dev"
  image          = "docker.pkg.github.com/byuoitav/camera-services/aver-dev"
  image_version  = "2fa8a90"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["aver-dev.av.byu.edu"]
  private           = true
  container_env = {
    "GIN_MODE" = "release"
  }
  container_args = [
    "--port", "8080",
    "--log-level", "info",
    "--name", "k8s-camera-services-aver-dev",
    "--event-url", data.aws_ssm_parameter.event_url.value,
    "--dns-addr", data.aws_ssm_parameter.dns_addr.value,
    "--cam-username", data.aws_ssm_parameter.aver_username.value,
    "--cam-password", data.aws_ssm_parameter.aver_password.value,
  ]
  health_check = false
}

module "aver" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-aver"
  image          = "docker.pkg.github.com/byuoitav/camera-services/aver-dev"
  image_version  = "2fa8a90"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["aver.av.byu.edu"]
  replicas          = 3
  private           = true
  container_env = {
    "GIN_MODE" = "release"
  }
  container_args = [
    "--port", "8080",
    "--log-level", "info",
    "--name", "k8s-camera-services-aver",
    "--event-url", data.aws_ssm_parameter.event_url.value,
    "--dns-addr", data.aws_ssm_parameter.dns_addr.value,
    "--cam-username", data.aws_ssm_parameter.aver_username.value,
    "--cam-password", data.aws_ssm_parameter.aver_password.value,
  ]
  health_check = false
}
