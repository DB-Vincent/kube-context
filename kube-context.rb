# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class KubeContext < Formula
  desc ""
  homepage "https://github.com/DB-Vincent/kube-context"
  version "0.4.0"

  on_macos do
    on_intel do
      url "https://github.com/DB-Vincent/kube-context/releases/download/v0.4.0/kube-context_Darwin_x86_64.tar.gz"
      sha256 "bf262bd21903d64363dfdc10243fab521b6a3619c77785dbadfda6137ee72291"

      def install
        bin.install "kube-context"
      end
    end
    on_arm do
      url "https://github.com/DB-Vincent/kube-context/releases/download/v0.4.0/kube-context_Darwin_arm64.tar.gz"
      sha256 "53b3b9768740f0cad6197fc8d75003ab1bf4221b66d5777e141e1cd5b51788b7"

      def install
        bin.install "kube-context"
      end
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/DB-Vincent/kube-context/releases/download/v0.4.0/kube-context_Linux_x86_64.tar.gz"
        sha256 "0029b017c9956d8f6c060d77280752359252286ecf9afe9b0f08171c6e192cda"

        def install
          bin.install "kube-context"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/DB-Vincent/kube-context/releases/download/v0.4.0/kube-context_Linux_arm64.tar.gz"
        sha256 "e088b787c06d5eeb7a82f28196729bb0a3f710fd1f170fd8ff66d79a4dcbfc48"

        def install
          bin.install "kube-context"
        end
      end
    end
  end
end
