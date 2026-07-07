#!/bin/sh
echo "Building Jiddat3D..."

# Assuming tailwindcss linux binary is installed or npm is used in docker
# In Docker, we will install tailwindcss standalone for linux
./tailwindcss-linux -i ui/static/css/input.css -o ui/static/css/output.css --minify

CGO_ENABLED=0 go build -ldflags="-s -w" -o jiddat3d ./cmd/jiddat

echo "Build complete!"
