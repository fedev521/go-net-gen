module "spoke2_project" {
  source          = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/project?ref=v19.0.0"
  name            = "prj-gonetgen-infra-spoke2"
  billing_account = var.billing_account

  services = [
    "compute.googleapis.com",
  ]
}

module "spoke2_vpc" {
  source      = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-vpc?ref=v19.0.0"
  project_id  = module.spoke2_project.project_id
  name        = "vpc-gng-spoke2-01"
  data_folder = "net/subnets/spoke2"

  # peered with hub_vpc
}
