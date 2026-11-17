# typed: true
# frozen_string_literal: true

# Btptf is a formula for installing BTPTFExporter CLI
class Btptf < Formula
  desc "Command-line tool for Exporting SAP BTP Resources to Terraform"
  homepage "https://sap.github.io/terraform-exporter-btp/"
  version "2.0.0"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v2.0.0/btptf_2.0.0_darwin_arm64"
      sha256 "c7069b879bf9907fb3daccefdc9f90188c0b9f25712ac2a77c69e0b5f66677c7"
    else
      url url "https://github.com/SAP/terraform-exporter-btp/releases/download/v2.0.0/btptf_2.0.0_darwin_amd64"
      sha256 "c48d4eb1b9e316d86a5d1098685e813f170503dd0d97bea383906d0591a39f4a"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v2.0.0/btptf_2.0.0_linux_arm64"
      sha256 "03f25b90463d168553ec642006172faea1ad3e24a7455f784e723831d416a7b9"
    else
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v2.0.0/btptf_2.0.0_linux_amd64"
      sha256 "230fd523046cf088ecf417664a913836df6038a6da3611c67361ba86649ef68a"
      depends_on arch: :x86_64
    end
  end

  def install
    bin.install stable.url.split("/")[-1] => "btptf"
  end

  def caveats
    <<~EOS
      [HINT]
      Please ensure you have Terraform or OpenTofu installed.
      Run:
         btptf --help for more information.
    EOS
  end

  test do
     system "#{bin}/btptf", "--version"
  end
end
