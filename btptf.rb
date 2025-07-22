# typed: true
# frozen_string_literal: true

# Btptf is a formula for installing BTPTFExporter CLI
class Btptf < Formula
  desc "Command-line tool for Exporting SAP BTP Resources to Terraform"
  homepage "https://sap.github.io/terraform-exporter-btp/"
  version "1.3.0"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.3.0/btptf_1.3.0_darwin_arm64"
      sha256 "f8aa38064f58801cdce120368b4c66ffbb9a3c248fc6991171f0d10ee2e86b21"
    else
      url url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.3.0/btptf_1.3.0_darwin_amd64"
      sha256 "6e63db44d28e300a7588c27cafc9cfdba5c9e684e9ea811c5a9bf4b41e4b9c86"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.3.0/btptf_1.3.0_linux_arm64"
      sha256 "09960a53d4b7d7a405b616a91ad4e32b2b00b9935404ca838a1bd40c94becb63"
    else
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.3.0/btptf_1.3.0_linux_amd64"
      sha256 "b55ef80a19641f3c7f28a931cca03d24d1b73cda895c7c4ecf540b7179274b89"
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
