script "plan" {
  job {
    name        = "Terraform Deployment"
    description = "Plan Terraform changes."
    commands = [
      ["terraform", "plan", "-no-color", {
        enable_sharing = true
        mock_on_fail   = true
      }],
    ]
  }
}
