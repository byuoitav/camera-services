data "aws_ssm_parameter" "control_opa_token" {
  name = "/env/camera-services/control/opa-token"
}

data "aws_ssm_parameter" "dev_control_client_id" {
  name = "/env/camera-services/control/dev-client-id"
}

data "aws_ssm_parameter" "dev_control_client_secret" {
  name = "/env/camera-services/control/dev-client-secret"
}

data "aws_ssm_parameter" "control_client_id" {
  name = "/env/camera-services/control/client-id"
}

data "aws_ssm_parameter" "control_client_secret" {
  name = "/env/camera-services/control/client-secret"
}

data "aws_ssm_parameter" "control_signing_secret" {
  name = "/env/camera-services/control/signing-secret"
}

module "control_dev" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-control-dev"
  image          = "docker.pkg.github.com/byuoitav/camera-services/control-dev"
  image_version  = "980d0f2"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["cameras-dev.av.byu.edu"]
  private           = true
  container_env = {
    "GIN_MODE" = "release"
  }
  container_args = [
    "--port", "8080",
    "--log-level", "info",
    "--db-address", data.aws_ssm_parameter.prd_db_addr.value,
    "--db-username", data.aws_ssm_parameter.prd_db_username.value,
    "--db-password", data.aws_ssm_parameter.prd_db_password.value,
    "--key-service", "control-keys",
    "--callback-url", "https://cameras-dev.av.byu.edu",
    "--client-id", data.aws_ssm_parameter.dev_control_client_id.value,
    "--client-secret", data.aws_ssm_parameter.dev_control_client_secret.value,
    "--gateway-url", data.aws_ssm_parameter.gateway_url.value,
    "--opa-url", data.aws_ssm_parameter.opa_url.value,
    "--opa-token", data.aws_ssm_parameter.control_opa_token.value,
    "--signing-secret", data.aws_ssm_parameter.control_signing_secret,
    "--aver-proxy", "http://camera-services-aver-dev",
    "--axis-proxy", "http://camera-services-axis-dev",
  ]
  health_check = false
}

module "control" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-control"
  image          = "docker.pkg.github.com/byuoitav/camera-services/control-dev"
  image_version  = "e3cc934"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["cameras.av.byu.edu"]
  replicas          = 1
  container_env = {
    "GIN_MODE" = "release"
  }
  container_args = [
    "--port", "8080",
    "--log-level", "info",
    "--db-address", data.aws_ssm_parameter.prd_db_addr.value,
    "--db-username", data.aws_ssm_parameter.prd_db_username.value,
    "--db-password", data.aws_ssm_parameter.prd_db_password.value,
    "--key-service", "control-keys",
    "--callback-url", "https://cameras.av.byu.edu",
    "--client-id", data.aws_ssm_parameter.control_client_id.value,
    "--client-secret", data.aws_ssm_parameter.control_client_secret.value,
    "--gateway-url", data.aws_ssm_parameter.gateway_url.value,
    "--opa-url", data.aws_ssm_parameter.opa_url.value,
    "--opa-token", data.aws_ssm_parameter.control_opa_token.value,
    "--aver-proxy", "http://camera-services-aver",
    "--axis-proxy", "http://camera-services-axis",
  ]
  health_check = false
}
