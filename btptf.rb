# typed: true
# frozen_string_literal: true

# Btptf is a formula for installing BTPTFExporter CLI
class Btptf < Formula
  desc "Command-line tool for Exporting SAP BTP Resources to Terraform"
  homepage "https://sap.github.io/terraform-exporter-btp/"
  version "1.5.0"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.5.0/btptf_1.5.0_darwin_arm64"
      sha256 "46aff1bdf83bc1e7fa7ebf5985f25133f30545c5853db14a1f8a900c13d1e770"
    else
      url url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.5.0/btptf_1.5.0_darwin_amd64"
      sha256 "0554da38617d7febafc65580df771e29641c536185203a7ab88d886ab369b643"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.5.0/btptf_1.5.0_linux_arm64"
      sha256 "7503de150af01d1cc457215bb02eaef41614c34ecbdaed31959e9ec1c2cfaee7"
    else
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.5.0/btptf_1.5.0_linux_amd64"
      sha256 "97bae73ecb7cf2311ccef4aebde7859c30478928c1846acac84f8ed597170d4a"
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
