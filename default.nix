{
  fetchFromGitHub,
  buildGoModule,
  ...
}:
buildGoModule rec {
  pname = "hyprstfu";
  version = "1.2.4";
  src = fetchFromGitHub {
    owner = "luisnquin";
    repo = pname;
    rev = "v${version}";
    hash = "sha256-g7m95uz4wWqZyHY3lulGLQiKKF6ACpxIqUVRhIg3ndU=";
  };

  ldflags = ["-X main.version=v${version}"];
  buildTarget = ".";

  vendorHash = "sha256-eiTGunM+U4AnBpkl1SymHOrY+Uij/ss0/BEtoZBfXB0=";
  doCheck = false;

  meta.mainProgram = "hyprstfu";
}
