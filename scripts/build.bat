@echo off
echo Building Jiddat3D...

echo Building Tailwind CSS...
.\tools\tailwindcss.exe -i ui\static\css\input.css -o ui\static\css\output.css --minify

echo Building Go Binary...
go build -ldflags="-s -w" -o jiddat3d.exe .\cmd\jiddat

echo Build complete!
