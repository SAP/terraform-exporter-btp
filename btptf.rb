# typed: true
# frozen_string_literal: true

# Btptf is a formula for installing BTPTFExporter CLI
class Btptf < Formula
  desc "Command-line tool for Exporting SAP BTP Resources to Terraform"
  homepage "https://sap.github.io/terraform-exporter-btp/"
  version "1.4.0"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.4.0/btptf_1.4.0_darwin_arm64"
      sha256 "950754be14683b12951117e14828949896a3a4b834835723bbc5c2df7431e873"
    else
      url url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.4.0/btptf_1.4.0_darwin_amd64"
      sha256 "0671d71778c4ee2d88cde33c2b5f0dcc8c7f7f1d3203c051a12c334b4c9369b9"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.4.0/btptf_1.4.0_linux_arm64"
      sha256 "fc7d820eb28d60510da7dfb2e2cfb8d15b8cb7b1d43ef4571b3bafd89005f6bd"
    else
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.4.0/btptf_1.4.0_linux_amd64"
      sha256 "f5b3e6467f1cbc60a546ad51d28c9563a44570a4ca4741584ba2dc7cb28c1a74"
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
