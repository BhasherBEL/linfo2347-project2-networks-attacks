let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-23.11";
  pkgs = import nixpkgs {
    config = { };
    overlay = { };
  };
in
pkgs.mkShellNoCC {
  packages = with pkgs; [
    inetutils
    go
  ];
}
