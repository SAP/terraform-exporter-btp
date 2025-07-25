# typed: true
# frozen_string_literal: true

# Btptf is a formula for installing BTPTFExporter CLI
class Btptf < Formula
  desc "Command-line tool for Exporting SAP BTP Resources to Terraform"
  homepage "https://sap.github.io/terraform-exporter-btp/"
  version "1.3.1"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.3.1/btptf_1.3.1_darwin_arm64"
      sha256 "582b66939e61a9a7d3232193a54a4214ddc534743287fd7c7e46c9000e6157dd"
    else
      url url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.3.1/btptf_1.3.1_darwin_amd64"
      sha256 "6c2b1b2190c69bf7b95a6c58ece00b99c4dada312ae85e6834b6181a34fba0cc"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.3.1/btptf_1.3.1_linux_arm64"
      sha256 "4215c493d0ab85dea7ba7bcc0a94a9ef3f853b9a80a35c86c2c5f3c98f161764"
    else
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.3.1/btptf_1.3.1_linux_amd64"
      sha256 "51737f0ec10c1fe61c31412b6acd6eaefa52dce7c4ca8c693b8e5eb4ef724316"
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
