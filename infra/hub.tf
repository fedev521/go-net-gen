module "hub_project" {
  source          = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/project?ref=v19.0.0"
  name            = "prj-gonetgen-infra-hub"
  billing_account = var.billing_account

  services = [
    "compute.googleapis.com",
  ]
}

module "hub_vpc" {
  source      = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-vpc?ref=v19.0.0"
  project_id  = module.hub_project.project_id
  name        = "vpc-gonetgen-hub-01"
  data_folder = "net/subnets"
}

module "gonetgen_sa" {
  source       = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v19.0.0"
  project_id   = module.hub_project.project_id
  name         = "sa-gonetgen-app"
  display_name = "GoNetGen app SA"
  description  = "Service account used by GoNetGen application to read data it needs."
  generate_key = false
  # authoritative roles granted *on* the service accounts to other identities
  iam = {}
  # non-authoritative roles granted *to* the service accounts on other resources
  iam_project_roles = {
    (module.hub_project.project_id) = [
      "roles/compute.networkViewer",
    ]
  }
}
