class Pvm < Formula
  desc "Pulumi Version Manager"
  homepage "https://github.com/tomski747/pvm"
  head "https://github.com/tomski747/pvm.git", branch: "main"
  
  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"pvm", "./cmd/pvm"
  end

  test do
    system "#{bin}/pvm", "--version"
  end
end