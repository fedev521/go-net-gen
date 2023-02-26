locals {
  hub_subnet     = module.hub_vpc.subnets["us-central1/subnet-gng-hub-usc1-01"]
  spoke1_subnet  = module.spoke1_vpc.subnets["us-central1/subnet-gng-spoke1-usc1-01"]
  spoke2_subnet1 = module.spoke2_vpc.subnets["us-west1/subnet-gng-spoke2-usw1-01"]
  spoke2_subnet2 = module.spoke2_vpc.subnets["us-west1/subnet-gng-spoke2-usw1-02"]
}

# ------------------------------------------------------------------------------

module "hub_vm_01" {
  count = var.enable_vm_creation ? 1 : 0

  source        = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/compute-vm?ref=v19.0.0"
  project_id    = module.hub_project.project_id
  name          = "vm-hub-usc1b-01"
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

# ------------------------------------------------------------------------------

module "spoke1_vm_01" {
  count = var.enable_vm_creation ? 1 : 0

  source        = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/compute-vm?ref=v19.0.0"
  project_id    = module.spoke1_project.project_id
  name          = "vm-spoke1-usc1a-01"
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

module "spoke1_vm_02" {
  count = var.enable_vm_creation ? 1 : 0

  source        = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/compute-vm?ref=v19.0.0"
  project_id    = module.spoke1_project.project_id
  name          = "vm-spoke1-usc1a-02"
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

# ------------------------------------------------------------------------------

module "spoke2_subnet1_vms" {
  count = var.enable_vm_creation ? 2 : 0

  source        = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/compute-vm?ref=v19.0.0"
  project_id    = module.spoke2_project.project_id
  name          = "vm-spoke2-usw1a-0${count.index + 1}"
  zone          = "us-west1-a"
  instance_type = var.free_machine_type

  network_interfaces = [{
    network    = module.spoke2_vpc.self_link
    subnetwork = local.spoke2_subnet1.self_link
  }]

  boot_disk = {
    type = "pd-standard"
  }
}

module "spoke2_subnet2_vms" {
  count = var.enable_vm_creation ? 3 : 0

  source        = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/compute-vm?ref=v19.0.0"
  project_id    = module.spoke2_project.project_id
  name          = "vm-spoke2-usw1b-0${count.index + 1}"
  zone          = "us-west1-b"
  instance_type = var.free_machine_type

  network_interfaces = [{
    network    = module.spoke2_vpc.self_link
    subnetwork = local.spoke2_subnet2.self_link
  }]

  boot_disk = {
    type = "pd-standard"
  }
}
