script "teardown" {
  job {
    name        = "Terraform Teardown"
    description = "Destroy Terraform setup."

    commands = [
      ["terraform", "destroy", "-auto-approve", "-no-color", {
        enable_sharing = true
        mock_on_fail   = true
      }],
    ]
  }
}
