# typed: true
# frozen_string_literal: true

# Btptf is a formula for installing BTPTFExporter CLI
class Btptf < Formula
  desc "Command-line tool for Exporting SAP BTP Resources to Terraform"
  homepage "https://sap.github.io/terraform-exporter-btp/"
  version "2.0.1"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v2.0.1/btptf_2.0.1_darwin_arm64"
      sha256 "a211f3a6c258db6fe1a3a5ffe6f914503d5a57a48bf6a17e671c8f2a387ea1be"
    else
      url url "https://github.com/SAP/terraform-exporter-btp/releases/download/v2.0.1/btptf_2.0.1_darwin_amd64"
      sha256 "81908538de41b1512b7e3bca00d54782fcd726ecc69274717e77455544912007"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v2.0.1/btptf_2.0.1_linux_arm64"
      sha256 "c2deb940eb3307d96c29d44f70ea5ff89e4fe373b799d75fd14ab0a721141325"
    else
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v2.0.1/btptf_2.0.1_linux_amd64"
      sha256 "798002afad6b10c09dbf928b2f980cf9d58cde1bffcc7e48cf6afce047fc465c"
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
