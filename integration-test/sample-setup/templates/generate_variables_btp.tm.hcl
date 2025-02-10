generate_hcl "_terramate_generated_variables.tf" {
  condition = tm_contains(terramate.stack.tags, "btp")
  content {

    variable "globalaccount" {
      description = "Subdomain of the global account"
      type        = string
    }

    variable "project_name" {
      description = "Name of the project"
      type        = string
      default     = "Inttest Exporter"
    }

    variable "subaccount_region" {
      description = "Region of the subaccount"
      type        = string
      default     = "us10"
      validation {
        condition     = contains(["us10", "eu10"], var.subaccount_region)
        error_message = "Region must be one of us10 or eu10"
      }
    }

    variable "project_costcenter" {
      description = "Cost center of the project"
      type        = string
      default     = "12345"
      validation {
        condition     = can(regex("^[0-9]{5}$", var.project_costcenter))
        error_message = "Cost center must be a 5 digit number"
      }
    }

    variable "cf_landscape_label" {
      type        = string
      description = "The Cloud Foundry landscape (format example us10-001)."
      default     = ""
    }

  }
}
