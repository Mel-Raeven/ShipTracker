{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [ pkgs.awscli pkgs.gcc pkgs.python3 pkgs.python3Packages.numpy pkgs.python3Packages.boto3 pkgs.python3Packages.folium pkgs.python3Packages.pandas ];

  shellHook = ''
    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:${pkgs.gcc}/lib
  '';
}
