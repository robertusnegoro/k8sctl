class K8ctl < Formula
  desc "Enhanced kubectl with colored tables and built-in context/namespace management"
  homepage "https://github.com/robertusnegoro/k8ctl"
  url "https://github.com/robertusnegoro/k8ctl/releases/download/v0.1.0/k8ctl_darwin_amd64.tar.gz"
  sha256 "" # Will be filled by GoReleaser
  version "0.1.0"

  if Hardware::CPU.arm?
    url "https://github.com/robertusnegoro/k8ctl/releases/download/v0.1.0/k8ctl_darwin_arm64.tar.gz"
    sha256 "" # Will be filled by GoReleaser
  end

  def install
    bin.install "k8ctl"
    bash_completion.install "completions/bash/k8ctl.bash" => "k8ctl"
    zsh_completion.install "completions/zsh/_k8ctl" => "_k8ctl"
    fish_completion.install "completions/fish/k8ctl.fish" => "k8ctl.fish"
  end

  test do
    system "#{bin}/k8ctl", "version"
  end
end
