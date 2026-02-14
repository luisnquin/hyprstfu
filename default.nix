{
  fetchFromGitHub,
  buildGoModule,
  ...
}:
buildGoModule rec {
  pname = "hyprstfu";
  version = "1.2.5";
  src = fetchFromGitHub {
    owner = "luisnquin";
    repo = pname;
    rev = "v${version}";
    hash = "sha256-WLrSh9j+HVB3LPTEr/UnUJF9JtQSrmI5Kx/xUhs6RE8=";
  };

  ldflags = ["-X main.version=v${version}"];
  buildTarget = ".";

  vendorHash = "sha256-eiTGunM+U4AnBpkl1SymHOrY+Uij/ss0/BEtoZBfXB0=";
  doCheck = false;

  meta.mainProgram = "hyprstfu";
}
