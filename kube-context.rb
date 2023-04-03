# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class KubeContext < Formula
  desc ""
  homepage "https://github.com/DB-Vincent/kube-context"
  version "0.2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/DB-Vincent/kube-context/releases/download/v0.2.0/kube-context_Darwin_arm64.tar.gz"
      sha256 "412a19a43f239e7ca04b38127cd0c598cb59f9431cc4a9beb1a01d16217bbae9"

      def install
        bin.install "kube-context"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/DB-Vincent/kube-context/releases/download/v0.2.0/kube-context_Darwin_x86_64.tar.gz"
      sha256 "1bc7aa9fdcde464bd4e090cbc061e8b5409e55b1e2ed8e32cd647219cc750e43"

      def install
        bin.install "kube-context"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/DB-Vincent/kube-context/releases/download/v0.2.0/kube-context_Linux_arm64.tar.gz"
      sha256 "fc7469635589c76391c4770c9282210454c43e7e48877ba1aa6b052d07bced49"

      def install
        bin.install "kube-context"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/DB-Vincent/kube-context/releases/download/v0.2.0/kube-context_Linux_x86_64.tar.gz"
      sha256 "93536f992a13a5681ab5d0d9db55507c838196d5421370cd5edd9a1e6233409a"

      def install
        bin.install "kube-context"
      end
    end
  end
end
