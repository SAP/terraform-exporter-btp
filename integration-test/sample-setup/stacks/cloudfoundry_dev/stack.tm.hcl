stack {
  name        = "cloudfoundry-dev"
  description = "CF for BTP setup (DEV)"
  tags        = ["cloudfoundry", "dev"]
  id          = "ce1c14ee-ea28-412f-8aa5-4d5893a42c36"
  after       = ["tag:btp:dev"]
}
