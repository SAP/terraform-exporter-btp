// Configure default default Terraform providers
globals "terraform" "providers" "btp" {
  version_dev = "~> 1.21.3"
}

globals "terraform" "providers" "cloudfoundry" {
  version_dev = "~> 1.14.0"
}
