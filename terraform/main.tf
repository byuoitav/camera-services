terraform {
  backend "s3" {
    bucket     = "terraform-state-storage-586877430255"
    lock_table = "terraform-state-lock-586877430255"
    region     = "us-west-2"

    // THIS MUST BE UNIQUE
    key = "camera-services.tfstate"
  }
}

provider "aws" {
  region = "us-west-2"
}

data "aws_ssm_parameter" "eks_cluster_endpoint" {
  name = "/eks/av-cluster-endpoint"
}

provider "kubernetes" {
  host = data.aws_ssm_parameter.eks_cluster_endpoint.value
}

data "aws_ssm_parameter" "gateway_url" {
  name = "/env/gateway-url"
}

data "aws_ssm_parameter" "opa_url" {
  name = "/env/opa-url"
}

data "aws_ssm_parameter" "control_opa_token" {
  name = "/env/camera-services/control/opa-token"
}

data "aws_ssm_parameter" "prd_db_addr" {
  name = "/env/couch-new-address"
}

data "aws_ssm_parameter" "prd_db_username" {
  name = "/env/couch-username"
}

data "aws_ssm_parameter" "prd_db_password" {
  name = "/env/couch-password"
}

data "aws_ssm_parameter" "event_url" {
  name = "/env/camera-services/event-url"
}

data "aws_ssm_parameter" "dns_addr" {
  name = "/env/camera-services/dns-addr"
}

data "aws_ssm_parameter" "aver_username" {
  name = "/env/camera-services/aver/username"
}

data "aws_ssm_parameter" "aver_password" {
  name = "/env/camera-services/aver/password"
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

data "aws_ssm_parameter" "slack_token" {
  name = "/env/camera-services/slack/slack-token"
}

data "aws_ssm_parameter" "slack_channel" {
  name = "/env/camera-services/slack/slack-channel"
}

data "aws_ssm_parameter" "hub_address" {
  name = "/env/hub-address"
}

module "aver_dev" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-aver-dev"
  image          = "docker.pkg.github.com/byuoitav/camera-services/aver-dev"
  image_version  = "9cca408"
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
  image_version  = "be8cb3e"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["aver.av.byu.edu"]
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

module "axis_dev" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-axis-dev"
  image          = "docker.pkg.github.com/byuoitav/camera-services/axis-dev"
  image_version  = "9cca408"
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
  image_version  = "be8cb3e"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["axis.av.byu.edu"]
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

module "slack" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-slack"
  image          = "docker.pkg.github.com/byuoitav/camera-services/camera-slack-dev"
  image_version  = "9cca408"
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
    "--snapshot-delay", "5s",
    "--slack-token", data.aws_ssm_parameter.slack_token.value,
    "--channel-id", data.aws_ssm_parameter.slack_channel.value,
  ]
  health_check = false
}

module "control_dev" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "camera-services-control-dev"
  image          = "docker.pkg.github.com/byuoitav/camera-services/control-dev"
  image_version  = "be8cb3e"
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
  image_version  = "7a46887"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/camera-services"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["cameras.av.byu.edu"]
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
