{
  description = "geniveev dev shell";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.utils.follows = "utils";
    };
  };

  outputs = { self, nixpkgs, utils, gomod2nix }:
    utils.lib.eachDefaultSystem (system:
    let pkgs =  import nixpkgs {
      inherit system;
      overlays = [
        gomod2nix.overlays.default
      ];
    };
    in {
      packages.default = pkgs.buildGoApplication {
        pname = "geniveev";
        version = "0.1.1";
        pwd = ./.;
        src = ./.;
        modules = ./gomod2nix.toml;
      };

      devShells.default = pkgs.mkShell {
        packages = [
          pkgs.bashInteractive
          pkgs.go_1_20
          pkgs.goreleaser
          gomod2nix.packages.${system}.default
        ];
      };
    });
}
