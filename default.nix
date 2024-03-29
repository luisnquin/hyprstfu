{
  fetchFromGitHub,
  buildGo122Module,
  ...
}:
buildGo122Module rec {
  pname = "hyprstfu";
  version = "1.2.1";
  src = fetchFromGitHub {
    owner = "luisnquin";
    repo = pname;
    rev = "v${version}";
    hash = "sha256-hVldMjaDfZND9DduhLEDK3ga1KyjSgf+12d6+75Mhzo=";
  };

  ldflags = ["-X main.version=v${version}"];
  buildTarget = ".";

  vendorHash = "sha256-5Ahu8N1hV05QkT2y28e6EPHKQq1+YcD6E6mp1b3duEo=";
  doCheck = false;
}
