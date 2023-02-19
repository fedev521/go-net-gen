variable "billing_account" {
  description = "Billing account id."
  type        = string
  sensitive   = true
  default     = null
}

variable "free_machine_type" {
  description = "VM machine type in the free tier. Don't change the value."
  type        = string
  default     = "e2-micro"
}
