terraform {
  backend "s3" {
    bucket         = "terraform-state-storage-586877430255"
    dynamodb_table = "terraform-state-lock-586877430255"
    region         = "us-west-2"

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
  config_path = "~/.kube/config"
}

data "aws_ssm_parameter" "gateway_url" {
  name = "/env/gateway-url"
}

data "aws_ssm_parameter" "opa_url" {
  name = "/env/opa-url"
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

data "aws_ssm_parameter" "hub_address" {
  name = "/env/hub-address"
}
