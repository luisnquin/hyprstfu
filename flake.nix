{
  description = "Utility to mute Hyprland windows for PulseAudio and Pipewire";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
        };

        defaultPackage = pkgs.callPackage (builtins.path {
          name = "hyprstfu-package";
          path = ./default.nix;
        }) {};
      in {
        inherit defaultPackage;

        defaultApp = flake-utils.lib.mkApp {
          drv = defaultPackage;
        };

        devShell = pkgs.mkShell {
          buildInputs = [defaultPackage];
        };
      }
    );
}
