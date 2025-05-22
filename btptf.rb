# typed: true
# frozen_string_literal: true

# Btptf is a formula for installing BTPTFExporter CLI
class Btptf < Formula
  desc "Command-line tool for Exporting SAP BTP Resources to Terraform"
  homepage "https://sap.github.io/terraform-exporter-btp/"
  version "1.1.0"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.1.0/btptf_1.1.0_darwin_arm64"
      sha256 "9d9e2c888d12dc1bc8e78e53682e6dfc2c25b832185161a7fea6c5cd286ee06c"
    else
      url url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.1.0/btptf_1.1.0_darwin_amd64"
      sha256 "eedef62132a5da618bb26d009d71a1586b4cba6ed4004d2dd0bc4e90e4298589"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.1.0/btptf_1.1.0_linux_arm64"
      sha256 "044b15563c8920305bd37028d21d2de4103311a4e36d6d045eb67242cd12ce43"
    else
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.1.0/btptf_1.1.0_linux_amd64"
      sha256 "69d1e8f336d4f5c89d4dd1c9f346ed88d00b5e35936afbd1e3a163a9e7440ff9"
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
