script "deploy" {
  job {
    name        = "Terraform Deployment"
    description = "Apply Terraform changes."
    commands = [
      ["terraform", "apply", "-auto-approve", "-no-color", {
        enable_sharing = true
        mock_on_fail   = true
      }],
    ]
  }
}
