class Pvm < Formula
  desc "Pulumi Version Manager"
  homepage "https://github.com/tomski747/pvm"
  url "https://github.com/tomski747/pvm.git", 
      tag:      "v0.0.0",
      revision: "homebrew"
  head "https://github.com/tomski747/pvm.git", 
      branch: "homebrew"
  
  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"pvm", "./cmd/pvm"
  end

  test do
    system "#{bin}/pvm", "--version"
  end
end