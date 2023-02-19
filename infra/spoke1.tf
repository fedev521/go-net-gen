locals {
  spoke1_subnet = module.spoke1_vpc.subnets["us-central1/subnet-gonetgen-spoke1-usc1-01"]
}

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
  name        = "vpc-gonetgen-spoke1-01"
  data_folder = "net/subnets/spoke1"

  # peered with hub_vpc
}

module "spoke1_vm_01" {
  source        = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/compute-vm?ref=v19.0.0"
  project_id    = module.spoke1_project.project_id
  name          = "vm-gonetgen-spoke1-usc1a-01"
  zone          = "us-central1-a"
  instance_type = var.free_machine_type

  network_interfaces = [{
    network    = module.spoke1_vpc.self_link
    subnetwork = local.spoke1_subnet.self_link
  }]

  boot_disk = {
    type = "pd-standard"
  }
}
