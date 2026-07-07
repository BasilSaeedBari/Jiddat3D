package compiler

import (
	"fmt"
	"path/filepath"
	"strings"
	"os/exec"
)

// ProcessImage converts the uploaded image to responsive WebP variants.
// Note: Requires cwebp installed on the host. If cwebp is not present, it gracefully returns without generating variants.
func ProcessImage(uploadPath, filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return nil
	}

	// We simply use cwebp here to convert and resize.
	// 480, 960, 1600 are the responsive widths requested.
	widths := []int{480, 960, 1600}
	
	for _, w := range widths {
		outFilename := fmt.Sprintf("%s_%dw.webp", strings.TrimSuffix(filename, ext), w)
		outPath := filepath.Join(uploadPath, outFilename)
		
		cmd := exec.Command("cwebp", "-resize", fmt.Sprintf("%d", w), "0", filepath.Join(uploadPath, filename), "-o", outPath)
		// If cwebp doesn't exist on host, this will fail silently and not crash the app.
		if err := cmd.Run(); err != nil {
			fmt.Printf("cwebp failed or missing: %v\n", err)
			return err
		}
	}
	
	return nil
}
