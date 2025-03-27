# typed: true
# frozen_string_literal: true

# Btptf is a formula for installing BTPTFExporter CLI
class Btptf < Formula
  desc "Command-line tool for Exporting SAP BTP Resources to Terraform"
  homepage "https://sap.github.io/terraform-exporter-btp/"
  version "1.0.0-rc1"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.0.0-rc1/btptf_1.0.0-rc1_darwin_arm64"
      sha256 "dcd0bad837211cea2cf81c94ccc2c5e4457365f6e7ef661523af15faa1aa3416"
    else
      url url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.0.0-rc1/btptf_1.0.0-rc1_darwin_amd64"
      sha256 "ec5fc58d89050fa4cb4cc4f95a9e77a5375ff9efe053422ea7f831a6d43762a8"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.0.0-rc1/btptf_1.0.0-rc1_linux_arm64"
      sha256 "b74ac46699640f1c54620385af445b36cdb21e1ebf32eaa1a092676d8ea32522"
    else
      url "https://github.com/SAP/terraform-exporter-btp/releases/download/v1.0.0-rc1/btptf_1.0.0-rc1_linux_amd64"
      sha256 "63943bee3514411184146335afb890f015571a6b66a3d005501af524f0126f76"
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
