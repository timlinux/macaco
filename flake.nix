{
  description = "MoCaCo - Motion Capture Combatant - Vim motion training game";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        version = "1.0.0";
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "macaco";
            inherit version;
            src = ./.;

            vendorHash = null; # Uses vendor directory

            ldflags = [
              "-s"
              "-w"
              "-X main.Version=${version}"
              "-X main.BuildTime=nix-build"
            ];

            meta = with pkgs.lib; {
              description = "MoCaCo - Vim motion training game";
              homepage = "https://github.com/timlinux/macaco";
              license = licenses.mit;
              maintainers = [ ];
              mainProgram = "macaco";
            };
          };

          macaco = self.packages.${system}.default;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            golangci-lint
            python3
            python3Packages.mkdocs
            python3Packages.mkdocs-material
          ];

          shellHook = ''
            echo "MoCaCo development environment"
            echo "Go version: $(go version)"
            echo ""
            echo "Available commands:"
            echo "  make build     - Build the application"
            echo "  make run       - Run the application"
            echo "  make test      - Run tests"
            echo "  make docs      - Serve documentation"
          '';
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/macaco";
        };
      }
    );
}
