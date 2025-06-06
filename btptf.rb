# typed: true
# frozen_string_literal: true

# Btptf is a formula for installing BTPTFExporter CLI
class Btptf < Formula
  desc "Command-line tool for Exporting SAP BTP Resources to Terraform"
  homepage "https://sap.github.io/terraform-exporter-btp/"
  version "1.2.0"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.2.0/btptf_1.2.0_darwin_arm64"
      sha256 "71a7b4a66fc211328e5b387ee4da41dc6ab0f3c8042a534d089f82a5d5277c58"
    else
      url url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.2.0/btptf_1.2.0_darwin_amd64"
      sha256 "4ff2a2cdd223d2ef2d82134340b760bd5000187ee46c9ffaaae36fcaaa93d889"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.2.0/btptf_1.2.0_linux_arm64"
      sha256 "91997d8db970db18c9c5f72dc4bd67b66381932a767e947a021a348043649e9c"
    else
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.2.0/btptf_1.2.0_linux_amd64"
      sha256 "f2ee859621ffa5ec08ceb816224eac3720bab95256f6b598386f4b7637116c23"
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
