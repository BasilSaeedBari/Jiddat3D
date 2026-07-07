@echo off
echo Starting Jiddat3D Dev Server...

start "Tailwind Watcher" cmd /k ".\tools\tailwindcss.exe -i ui\static\css\input.css -o ui\static\css\output.css --watch"
start "PocketBase" cmd /k "go run .\cmd\jiddat serve"

echo Dev servers started. Press Ctrl+C in their respective windows to stop.
