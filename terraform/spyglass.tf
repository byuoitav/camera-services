data "aws_ssm_parameter" "spyglass_opa_token" {
  name = "/env/camera-services/spyglass/opa-token"
}

data "aws_ssm_parameter" "spyglass_client_id" {
  name = "/env/camera-services/spyglass/client-id"
}

data "aws_ssm_parameter" "spyglass_client_secret" {
  name = "/env/camera-services/spyglass/client-secret"
}

module "spyglass" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-spyglass"
  image          = "docker.pkg.github.com/byuoitav/camera-services/camera-spyglass-dev"
  image_version  = "37098b3"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["spyglass.av.byu.edu"]
  container_env = {
    "GIN_MODE" = "release"
  }
  container_args = [
    "--port", "8080",
    "--db-address", data.aws_ssm_parameter.prd_db_addr.value,
    "--db-username", data.aws_ssm_parameter.prd_db_username.value,
    "--db-password", data.aws_ssm_parameter.prd_db_password.value,
    "--key-service", "control-keys",
    "--callback-url", "https://spyglass.av.byu.edu",
    "--client-id", data.aws_ssm_parameter.spyglass_client_id.value,
    "--client-secret", data.aws_ssm_parameter.spyglass_client_secret.value,
    "--gateway-url", data.aws_ssm_parameter.gateway_url.value,
    "--opa-url", data.aws_ssm_parameter.opa_url.value,
    "--opa-token", data.aws_ssm_parameter.spyglass_opa_token.value
  ]
  health_check = false
}
