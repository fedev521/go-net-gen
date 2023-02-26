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
  name        = "vpc-gng-hub-01"
  data_folder = "net/subnets/hub"
}

module "peering_hub_spoke1" {
  source = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-vpc-peering?ref=v19.0.0"
  prefix = "peering"

  local_network              = module.hub_vpc.self_link
  peer_network               = module.spoke1_vpc.self_link
  export_local_custom_routes = true # hub -> spoke1
  export_peer_custom_routes  = true # spoke1 -> hub
  peer_create_peering        = true # create peering in both directions
}

module "peering_hub_spoke2" {
  source = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-vpc-peering?ref=v19.0.0"
  prefix = "peering"

  local_network              = module.hub_vpc.self_link
  peer_network               = module.spoke2_vpc.self_link
  export_local_custom_routes = true # hub -> spoke2
  export_peer_custom_routes  = true # spoke2 -> hub
  peer_create_peering        = true # create peering in both directions
  depends_on                 = [module.peering_hub_spoke1]
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
    (module.spoke1_project.project_id) = [
      "roles/compute.networkViewer",
    ]
    (module.spoke2_project.project_id) = [
      "roles/compute.networkViewer",
    ]
  }
}
