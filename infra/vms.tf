locals {
  hub_subnet    = module.hub_vpc.subnets["us-central1/subnet-gonetgen-hub-usc1-01"]
  spoke1_subnet = module.spoke1_vpc.subnets["us-central1/subnet-gonetgen-spoke1-usc1-01"]
}

module "hub_vm_01" {
  source        = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/compute-vm?ref=v19.0.0"
  project_id    = module.hub_project.project_id
  name          = "vm-gonetgen-hub-usc1b-01"
  zone          = "us-central1-b"
  instance_type = var.free_machine_type

  network_interfaces = [{
    network    = module.hub_vpc.self_link
    subnetwork = local.hub_subnet.self_link
  }]

  boot_disk = {
    type = "pd-standard"
  }
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
