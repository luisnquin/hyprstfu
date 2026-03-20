{
  description = "Utility to mute Hyprland windows for PulseAudio and Pipewire";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default-linux";
  };

  outputs = {
    self,
    nixpkgs,
    systems,
    ...
  }: let
    inherit (nixpkgs) lib;
    eachSystem = lib.genAttrs (import systems);
    pkgsFor = eachSystem (system:
      import nixpkgs {
        localSystem = system;
      });
  in {
    packages = eachSystem (system: {
      default = pkgsFor.${system}.callPackage (builtins.path {
        name = "hyprstfu-package";
        path = ./default.nix;
      }) {};
    });

    overlays.default = final: prev: {
      hyprstfu = self.packages.${final.system}.default;
    };

    devShells = eachSystem (system: {
      default = pkgsFor.${system}.mkShell {
        buildInputs = [self.packages.${system}.default];
      };
    });
  };
}
