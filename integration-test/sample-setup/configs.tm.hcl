// Configure default default Terraform providers
globals "terraform" "providers" "btp" {
  version_dev = "~> 1.19.0"
}

globals "terraform" "providers" "cloudfoundry" {
  version_dev = "~> 1.12.0"
}
