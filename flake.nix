{
  description = "renamatic - CLI tool for renaming functions in .gno files";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "renamatic";
          version = "0.1.0";
          src = ./.;

          vendorHash = null;

          meta = with pkgs.lib; {
            description = "CLI tool for renaming functions in .gno files";
            homepage = "https://github.com/notJoon/renamatic";
            license = licenses.apache20;
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            golangci-lint
            delve
            gofumpt
            git
            make
          ];

          shellHook = ''
            export GOPATH="$HOME/go"
            export PATH="$GOPATH/bin:$PATH"

            echo "Available tools:"
            echo "  - go: Go compiler"
            echo "  - gopls: Go language server"
            echo "  - golangci-lint: Go linter"
            echo "  - delve: Go debugger"
            echo "  - gofumpt: Go code formatter"
          '';
        };
      }
    );
}
