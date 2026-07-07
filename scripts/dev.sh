#!/bin/sh
echo "Starting Jiddat3D Dev Server..."

# Run tailwind and go in parallel
./tailwindcss-linux -i ui/static/css/input.css -o ui/static/css/output.css --watch &
TAILWIND_PID=$!

go run ./cmd/jiddat serve &
GO_PID=$!

# Wait for both processes
wait $TAILWIND_PID $GO_PID
