module "spoke1_project" {
  source          = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/project?ref=v19.0.0"
  name            = "prj-gonetgen-infra-spoke1"
  billing_account = var.billing_account

  services = [
    "compute.googleapis.com",
  ]
}
