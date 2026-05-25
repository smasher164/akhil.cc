{
  description = "A flake that sets up the devShell";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
     mexdown-src = {
      url = "github:smasher164/mexdown";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, mexdown-src }:
    let supportedSystems = [ "aarch64-linux" "i686-linux" "x86_64-linux" "aarch64-darwin" ];
    in flake-utils.lib.eachSystem supportedSystems (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ ];
        };

	mexdown = pkgs.buildGoModule {
	  pname = "mexdown";
	  version = "0.0.1";
	  src = mexdown-src;
	  overrideModAttrs = (_: {
	    outputHashAlgo = "sha256";
	    outputHash = "sha256-eqXE4fvDBch6G7fA5TRDEf+IBLpQPR0iwCasoO/+yUo=";
	    preBuild = ''
	      export HOME=$TMPDIR
	      go mod tidy
	    '';
	  });
	  vendorHash = "sha256-cU59wdqf5z1KbhQwgNuLwpDxspazsINI5dbLHx123Tc=";
	};
      in { devShell = pkgs.mkShell { buildInputs = [ pkgs.go mexdown ]; }; });
}
