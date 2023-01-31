module "spoke2_project" {
  source          = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/project?ref=v19.0.0"
  name            = "prj-gonetgen-infra-spoke2"
  billing_account = var.billing_account

  services = [
    "compute.googleapis.com",
  ]
}
