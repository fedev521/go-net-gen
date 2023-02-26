module "spoke1_project" {
  source          = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/project?ref=v19.0.0"
  name            = "prj-gonetgen-infra-spoke1"
  billing_account = var.billing_account

  services = [
    "compute.googleapis.com",
  ]
}

module "spoke1_vpc" {
  source      = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-vpc?ref=v19.0.0"
  project_id  = module.spoke1_project.project_id
  name        = "vpc-gng-spoke1-01"
  data_folder = "net/subnets/spoke1"

  # peered with hub_vpc
}
