{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
      };

      nativeBuildInputs = with pkgs; [
        go
        pkg-config
      ];

      buildInputs = with pkgs; [
        alsa-lib
        libGL
        xorg.libX11
        xorg.libXcursor
        xorg.libXext
        xorg.libXi
        xorg.libXinerama
        xorg.libXrandr
        xorg.libXxf86vm
      ];

      hardeningDisable = [
        "zerocallusedregs"
      ];
    in {
      devShells.default = pkgs.mkShell {
        inherit nativeBuildInputs buildInputs hardeningDisable;
        packages = with pkgs; [
          nil
          gopls
        ];
        shellHook = ''
          export LD_LIBRARY_PATH=${pkgs.wayland}/lib:${pkgs.lib.getLib pkgs.libGL}/lib:${pkgs.lib.getLib pkgs.libGL}/lib:$LD_LIBRARY_PATH
        '';
      };

      formatter = pkgs.alejandra;
    });
}
