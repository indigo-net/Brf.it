class Brfit < Formula
  desc "AI 코딩 어시스턴트를 위한 코드 브리핑 도구"
  homepage "https://github.com/indigo-net/Brf.it"
  version "0.10.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/indigo-net/Brf.it/releases/download/v0.10.0/Brf.it_0.10.0_darwin_arm64.tar.gz"
      sha256 "0355f166dd1ffd46cac5b08dbc9f419a9543c416bdc16a0f57e9447c09ac9925"
    else
      url "https://github.com/indigo-net/Brf.it/releases/download/v0.10.0/Brf.it_0.10.0_darwin_amd64.tar.gz"
      sha256 "d7caea6efa9b3eeaa28098b0827fe57f8cf056455e0c134fd118aefcbf9903c5"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/indigo-net/Brf.it/releases/download/v0.10.0/Brf.it_0.10.0_linux_arm64.tar.gz"
      sha256 "a44478c3ca676d2709b138e492500ae99ab3ea3318276725ea5bc349d2a1e1bb"
    else
      url "https://github.com/indigo-net/Brf.it/releases/download/v0.10.0/Brf.it_0.10.0_linux_amd64.tar.gz"
      sha256 "4fae556f0575f96adbc47d8087902db0f05c725702b9324fe653a15d41867775"
    end
  end

  def install
    bin.install "brfit"
  end

  test do
    system "#{bin}/brfit", "--version"
  end
end
