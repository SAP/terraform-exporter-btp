// TERRAMATE: GENERATED AUTOMATICALLY DO NOT EDIT

variable "globalaccount" {
  description = "Subdomain of the global account"
  type        = string
}
variable "project_name" {
  default     = "Inttest Exporter"
  description = "Name of the project"
  type        = string
}
variable "subaccount_region" {
  default     = "us10"
  description = "Region of the subaccount"
  type        = string
  validation {
    condition = contains([
      "us10",
      "eu10",
    ], var.subaccount_region)
    error_message = "Region must be one of us10 or eu10"
  }
}
variable "project_costcenter" {
  default     = "12345"
  description = "Cost center of the project"
  type        = string
  validation {
    condition     = can(regex("^[0-9]{5}$", var.project_costcenter))
    error_message = "Cost center must be a 5 digit number"
  }
}
variable "cf_landscape_label" {
  default     = ""
  description = "The Cloud Foundry landscape (format example us10-001)."
  type        = string
}
