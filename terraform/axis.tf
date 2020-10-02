module "axis_dev" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-axis-dev"
  image          = "docker.pkg.github.com/byuoitav/camera-services/axis-dev"
  image_version  = "629a75b"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["axis-dev.av.byu.edu"]
  private           = true
  container_env = {
    "GIN_MODE" = "release"
  }
  container_args = [
    "--port", "8080",
    "--log-level", "info",
    "--name", "k8s-camera-services-axis-dev",
    "--event-url", data.aws_ssm_parameter.event_url.value,
    "--dns-addr", data.aws_ssm_parameter.dns_addr.value,
  ]
  health_check = false
}

module "axis" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-axis"
  image          = "docker.pkg.github.com/byuoitav/camera-services/axis-dev"
  image_version  = "e4e4048"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["axis.av.byu.edu"]
  replicas          = 2
  private           = true
  container_env = {
    "GIN_MODE" = "release"
  }
  container_args = [
    "--port", "8080",
    "--log-level", "info",
    "--name", "k8s-camera-services-axis",
    "--event-url", data.aws_ssm_parameter.event_url.value,
    "--dns-addr", data.aws_ssm_parameter.dns_addr.value,
  ]
  health_check = false
}

