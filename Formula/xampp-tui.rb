class XamppTui < Formula
  desc "Terminal UI for managing XAMPP stack on Linux"
  homepage "https://github.com/MarcosLesca/xamp-tui"
  url "https://github.com/MarcosLesca/xamp-tui/archive/v0.1.0.tar.gz"
  sha256 "66f204de0e1b109f8049a7129c86ab43c88ca287e19a18ec5dafef322663d46e"
  license "MIT"
  head "https://github.com/MarcosLesca/xamp-tui.git"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", "xampp-tui"
    bin.install "xampp-tui"
  end

  test do
    # Binary exists and is executable
    assert_predicate bin/"xampp-tui", :exist?
    assert_predicate bin/"xampp-tui", :executable?
  end
end