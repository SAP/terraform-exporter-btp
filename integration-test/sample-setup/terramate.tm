terramate {
  config {
    experiments = [
      "outputs-sharing",
      "scripts",
    ]
  }
}

sharing_backend "default" {
  type     = terraform
  filename = "sharing_generated.tf"
  command  = ["terraform", "output", "-json"]
}
