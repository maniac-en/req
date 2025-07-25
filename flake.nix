{
  description = "Dev environment for req";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs, }: {
    devShells.x86_64-linux.default =
      let pkgs = import nixpkgs { system = "x86_64-linux"; };
      in pkgs.mkShell {
        buildInputs = [ pkgs.go pkgs.zsh ];

        shellHook = ''
          export GOPATH=$PWD/.gopath
          export GOBIN=$GOPATH/bin
          export PATH=$GOBIN:$PATH
          mkdir -p "$GOBIN"
          go mod tidy
          if [ -z "$ZSH_VERSION" ]; then
            exec zsh
          fi
        '';
      };
  };
}
