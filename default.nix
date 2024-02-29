{
  fetchFromGitHub,
  buildGo122Module,
  ...
}:
buildGo122Module rec {
  pname = "hyprstfu";
  version = "1.1.0";
  src = fetchFromGitHub {
    owner = "luisnquin";
    repo = pname;
    rev = "v${version}";
    hash = "sha256-H56T8g9evppa4FZR54CBlchLyCxgTNX12SIJo+RDiCY=";
  };

  ldflags = ["-X main.version=v${version}"];
  buildTarget = ".";

  vendorHash = "sha256-ENrDNIVGnUZ/APhYQaQZCIb/UXLG71yFEqzW9ZK+PPo=";
  doCheck = false;
}
