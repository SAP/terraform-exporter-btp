# typed: true
# frozen_string_literal: true

# Btptf is a formula for installing BTPTFExporter CLI
class Btptf < Formula
  desc "Command-line tool for Exporting SAP BTP Resources to Terraform"
  homepage "https://sap.github.io/terraform-exporter-btp/"
  version "1.6.0"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.6.0/btptf_1.6.0_darwin_arm64"
      sha256 "74b54d17af61ce45e3ddd3ac3ab5d63e87c6b57fec4de4ae735d251563b841cd"
    else
      url url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.6.0/btptf_1.6.0_darwin_amd64"
      sha256 "d927ddb618764d8a2207502666e813954e14b104ecb44c2b77dea57647d61bad"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.6.0/btptf_1.6.0_linux_arm64"
      sha256 "0a30d9f4358b5a0bf3e0f88af1a71068d0b77ed3faec50ea2e189c82daaba1b3"
    else
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.6.0/btptf_1.6.0_linux_amd64"
      sha256 "07a5a79bf2638bfe8ecc01bd17e2685dd5d018090691bf4ac6a3c3d0123fc648"
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
